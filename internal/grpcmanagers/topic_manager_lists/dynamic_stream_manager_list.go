package topic_manager_lists

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/momento"
)

// DynamicStreamManagerList manages a dynamic pool of gRPC channels for stream pubsub requests.
//
// The DynamicStreamManagerList will continually place the next available stream manager
// (or an error if no TopicGrpcManagers are available) on the streamManagerRequestQueue channel
// for the pubSubClient to pull from.
// Without this channel, a mutex and additional logic would be needed to protect the dynamic list
// of grpc managers from race conditions due to concurrent access.
type DynamicStreamManagerList struct {
	grpcManagers                []*grpcmanagers.TopicGrpcManager
	managerIndex                atomic.Uint64
	maxManagerCount             int    // maximum number of grpc channels that can be created
	currentMaxConcurrentStreams uint32 // current number of grpc channels * MAX_CONCURRENT_STREAMS_PER_CHANNEL
	logger                      logger.MomentoLogger
	newTopicManagerProps        *models.TopicStreamGrpcManagerRequest
	streamManagerRequestQueue   chan *StreamManagerRequest
	ctx                         context.Context
	cancel                      context.CancelFunc
}

func (list *DynamicStreamManagerList) GetCurrentNumberOfGrpcManagers() int {
	return len(list.grpcManagers)
}

// GetNextManager checks if there is a stream available to return, otherwise returns an error.
//
// Each grpc connection can multiplex 100 subscribe/publish requests. Grpc channels that already
// have 100 subscriptions will silently queue up subsequent requests. We must prevent subscription
// requests from queuing up if there are already numStreamChannels*100 concurrent streams as it
// causes the program to hang indefinitely with no error.
// If the current max concurrent streams is reached, checks if a new channel can be added and
// return a stream from the new channel, else return an error.
func (list *DynamicStreamManagerList) GetNextManager() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	// First check if there is enough grpc stream capacity to make a new subscription
	err := list.checkNumConcurrentStreams()
	if err != nil {
		return nil, err
	}

	// Max number of attempts is set to the max number of concurrent streams in order to preserve
	// the round-robin system (incrementing nextManagerIndex) but to not cut short the number
	//  of attempts in case there are many subscriptions starting up at the same time.
	for i := 0; uint32(i) < list.currentMaxConcurrentStreams; i++ {
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
	errorMessage := fmt.Sprintf("Cannot start new subscription, all grpc channels may be at maximum capacity. There are %d total subscriptions allowed across %d grpc channels. Please use the WithNumStreamGrpcChannels configuration if you wish to start more subscriptions.\n", list.currentMaxConcurrentStreams, list.maxManagerCount)
	return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.LimitExceededError, errorMessage, nil)
}

func (list *DynamicStreamManagerList) Close() {
	list.cancel() // Cancel context first to stop goroutines
	close(list.streamManagerRequestQueue)
	for _, topicManager := range list.grpcManagers {
		err := topicManager.Close()
		if err != nil {
			list.logger.Error("Error closing topic manager: %s", err.Error())
		}
	}
}

// Called by checkNumConcurrentStreams.
func (list *DynamicStreamManagerList) CountNumberOfActiveSubscriptions() int64 {
	count := int64(0)
	for _, topicManager := range list.grpcManagers {
		count += topicManager.NumActiveSubscriptions.Load()
	}
	return count
}

// Helper function to help sanity check number of concurrent streams before starting a new subscription
func (list *DynamicStreamManagerList) checkNumConcurrentStreams() momentoerrors.MomentoSvcErr {
	numActiveStreams := list.CountNumberOfActiveSubscriptions()
	list.logger.Debug("Current number of active subscriptions: %d", numActiveStreams)

	numStreamManagers := len(list.grpcManagers)

	if numActiveStreams >= int64(list.currentMaxConcurrentStreams) && numStreamManagers >= list.maxManagerCount {
		errorMessage := fmt.Sprintf(
			"Number of grpc streams: %d; number of channels: %d; max concurrent streams: %d; Already at maximum number of concurrent grpc streams, cannot make new subscribe requests\n",
			numActiveStreams, list.maxManagerCount, list.currentMaxConcurrentStreams,
		)
		return momentoerrors.NewMomentoSvcErr(momentoerrors.LimitExceededError, errorMessage, nil)
	} else if numActiveStreams >= int64(list.currentMaxConcurrentStreams) && numStreamManagers < list.maxManagerCount {
		// otherwise we can try to add a new manager
		err := list.addManager()
		if err != nil {
			return err
		}
		list.logger.Debug("Added new manager, current number of managers: %d", len(list.grpcManagers))
	}

	// If we are approaching the grpc maximum concurrent stream limit, log a warning
	if numStreamManagers == list.maxManagerCount {
		remainingStreams := int64(list.currentMaxConcurrentStreams) - numActiveStreams
		if remainingStreams < 10 {
			list.logger.Warn(
				"WARNING: approaching grpc maximum concurrent stream limit, %d remaining of total %d streams\n",
				remainingStreams, list.currentMaxConcurrentStreams,
			)
		}
	}
	return nil
}

// Called by checkNumConcurrentStreams to add more stream capacity if needed.
func (list *DynamicStreamManagerList) addManager() momentoerrors.MomentoSvcErr {
	streamTopicManager, err := grpcmanagers.NewStreamTopicGrpcManager(list.newTopicManagerProps)
	if err != nil {
		return err
	}
	list.grpcManagers = append(list.grpcManagers, streamTopicManager)
	list.currentMaxConcurrentStreams = uint32(len(list.grpcManagers)) * uint32(momento.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
	return nil
}

func NewDynamicStreamManagerList(request *models.TopicStreamGrpcManagerRequest, maxSubscriptions uint32, logger logger.MomentoLogger) (*DynamicStreamManagerList, chan *StreamManagerRequest, momentoerrors.MomentoSvcErr) {
	// make just one manager to start with
	streamTopicManagers := make([]*grpcmanagers.TopicGrpcManager, 0)
	streamTopicManager, err := grpcmanagers.NewStreamTopicGrpcManager(request)
	if err != nil {
		return nil, nil, err
	}
	streamTopicManagers = append(streamTopicManagers, streamTopicManager)
	logger.Debug("Max subscriptions: %d, max manager count: %d", maxSubscriptions, int(math.Ceil(float64(maxSubscriptions)/float64(momento.MAX_CONCURRENT_STREAMS_PER_CHANNEL))))

	// Unbuffered channel so the stream manager list will block on sending the next
	// available stream manager until the most recent request is processed.
	streamManagerRequestQueue := make(chan *StreamManagerRequest)
	ctx, cancel := context.WithCancel(context.Background())

	list := &DynamicStreamManagerList{
		grpcManagers:                streamTopicManagers,
		maxManagerCount:             int(math.Ceil(float64(maxSubscriptions) / float64(momento.MAX_CONCURRENT_STREAMS_PER_CHANNEL))),
		currentMaxConcurrentStreams: uint32(momento.MAX_CONCURRENT_STREAMS_PER_CHANNEL), // for one channel
		logger:                      logger,
		newTopicManagerProps:        request,
		streamManagerRequestQueue:   streamManagerRequestQueue,
		ctx:                         ctx,
		cancel:                      cancel,
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
