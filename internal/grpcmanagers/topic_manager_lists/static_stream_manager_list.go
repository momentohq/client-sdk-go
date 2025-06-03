package topic_manager_lists

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/momento"
)

// TopicStreamManagerListWithBookkeeping extends the TopicManagerList interface to
// include methods used by StaticStreamManagerList and DynamicStreamManagerList.
type TopicStreamManagerListWithBookkeeping interface {
	TopicManagerList
	CountNumberOfActiveSubscriptions() int64
	checkNumConcurrentStreams() momentoerrors.MomentoSvcErr
}

// StreamManagerRequest is used for putting the next available stream manager on a channel for the
// pubSubClient to pull from, or an error that specifies why no stream manager is available.
type StreamManagerRequest struct {
	TopicManager *grpcmanagers.TopicGrpcManager
	Err          error
}

// StaticStreamManagerList manages a static pool of gRPC channels for stream pubsub requests.
//
// The StaticStreamManagerList will continually place the next available stream manager
// (or an error if no TopicGrpcManagers are available) on the streamManagerRequestQueue channel
// for the pubSubClient to pull from.
// Without this channel, a mutex and additional logic would be needed to protect the static list
// of grpc managers from race conditions due to concurrent access.
type StaticStreamManagerList struct {
	grpcManagers              []*grpcmanagers.TopicGrpcManager
	managerIndex              atomic.Uint64
	maxConcurrentStreams      uint32
	logger                    logger.MomentoLogger
	streamManagerRequestQueue chan *StreamManagerRequest
	ctx                       context.Context
	cancel                    context.CancelFunc
}

// GetNextManager checks if there is a stream available to return, otherwise returns an error.
//
// Each grpc connection can multiplex 100 subscribe/publish requests. Grpc channels that already
// have 100 subscriptions will silently queue up subsequent requests. We must prevent subscription
// requests from queuing up if there are already numStreamChannels*100 concurrent streams as it
// causes the program to hang indefinitely with no error.
func (list *StaticStreamManagerList) GetNextManager() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	// First check if there is enough grpc stream capacity to make a new subscription
	err := list.checkNumConcurrentStreams()
	if err != nil {
		return nil, err
	}

	// Max number of attempts is set to the max number of concurrent streams in order to preserve
	// the round-robin system (incrementing nextManagerIndex) but to not cut short the number
	//  of attempts in case there are many subscriptions starting up at the same time.
	for i := 0; uint32(i) < list.maxConcurrentStreams; i++ {
		nextManagerIndex := list.managerIndex.Add(1)
		topicManager := list.grpcManagers[nextManagerIndex%uint64(len(list.grpcManagers))]
		newCount := topicManager.NumActiveSubscriptions.Add(1)
		if newCount <= int64(momento.MAX_CONCURRENT_STREAMS_PER_CHANNEL) {
			list.logger.Debug("Starting new subscription on grpc channel %d which now has %d streams", nextManagerIndex%uint64(len(list.grpcManagers)), newCount)
			return topicManager, nil
		}
		topicManager.NumActiveSubscriptions.Add(-1)
	}

	// If there are no more streams available, return an error
	errorMessage := fmt.Sprintf("Cannot start new subscription, all grpc channels may be at maximum capacity. There are %d total subscriptions allowed across %d grpc channels. Please use the WithNumStreamGrpcChannels configuration if you wish to start more subscriptions.\n", list.maxConcurrentStreams, len(list.grpcManagers))
	return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.LimitExceededError, errorMessage, nil)
}

func (list *StaticStreamManagerList) Close() {
	list.cancel() // Cancel context first to stop goroutines
	close(list.streamManagerRequestQueue)
	for _, topicManager := range list.grpcManagers {
		err := topicManager.Close()
		if err != nil {
			list.logger.Error("Error closing topic manager: %v", err)
		}
	}
}

// Helper to count number of active subscriptions in a StaticStreamManagerList
func (list *StaticStreamManagerList) CountNumberOfActiveSubscriptions() int64 {
	count := int64(0)
	for _, topicManager := range list.grpcManagers {
		count += topicManager.NumActiveSubscriptions.Load()
	}
	return count
}

// Helper function to help sanity check number of concurrent streams before starting a new subscription
func (list *StaticStreamManagerList) checkNumConcurrentStreams() momentoerrors.MomentoSvcErr {
	numActiveStreams := list.CountNumberOfActiveSubscriptions()
	if numActiveStreams >= int64(list.maxConcurrentStreams) {
		errorMessage := fmt.Sprintf(
			"Number of grpc streams: %d; number of channels: %d; max concurrent streams: %d; Already at maximum number of concurrent grpc streams, cannot make new subscribe requests\n",
			numActiveStreams, len(list.grpcManagers), list.maxConcurrentStreams,
		)
		return momentoerrors.NewMomentoSvcErr(momentoerrors.LimitExceededError, errorMessage, nil)
	}
	// If we are approaching the grpc maximum concurrent stream limit, log a warning
	remainingStreams := int64(list.maxConcurrentStreams) - numActiveStreams
	if remainingStreams < 10 {
		list.logger.Warn(
			"WARNING: approaching grpc maximum concurrent stream limit, %d remaining of total %d streams\n",
			remainingStreams, list.maxConcurrentStreams,
		)
	}
	return nil
}

func NewStaticStreamManagerList(
	request *models.TopicStreamGrpcManagerRequest,
	numStreamChannels uint32,
	logger logger.MomentoLogger,
) (*StaticStreamManagerList, chan *StreamManagerRequest, momentoerrors.MomentoSvcErr) {
	streamTopicManagers := make([]*grpcmanagers.TopicGrpcManager, 0)
	for i := 0; uint32(i) < numStreamChannels; i++ {
		streamTopicManager, err := grpcmanagers.NewStreamTopicGrpcManager(request)
		if err != nil {
			return nil, nil, err
		}
		streamTopicManagers = append(streamTopicManagers, streamTopicManager)
	}

	// Use an unbuffered channel here so the StaticStreamManagerList will block on sending
	// the next available stream manager until the most recent request is processed.
	streamManagerRequestQueue := make(chan *StreamManagerRequest)
	ctx, cancel := context.WithCancel(context.Background())

	list := &StaticStreamManagerList{
		grpcManagers:              streamTopicManagers,
		maxConcurrentStreams:      numStreamChannels * uint32(momento.MAX_CONCURRENT_STREAMS_PER_CHANNEL),
		logger:                    logger,
		streamManagerRequestQueue: streamManagerRequestQueue,
		ctx:                       ctx,
		cancel:                    cancel,
	}

	// Start goroutine to continually make the next available stream manager
	// available on the streamManagerRequestQueue
	go func() {
		for {
			select {
			case <-list.ctx.Done():
				return
			default:
				topicManager, err := list.GetNextManager()
				select {
				case <-list.ctx.Done():
					return
				case streamManagerRequestQueue <- &StreamManagerRequest{
					TopicManager: topicManager,
					Err:          err,
				}:
				}
			}
		}
	}()

	return list, streamManagerRequestQueue, nil
}
