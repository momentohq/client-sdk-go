package topic_manager_lists

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

type staticStreamGrpcManagerPool struct {
	grpcManagers                    []*grpcmanagers.TopicGrpcManager
	managerIndex                    atomic.Uint64
	currentActiveStreamsCount       atomic.Uint64
	maxConcurrentStreams            uint32
	logger                          logger.MomentoLogger
	nextAvailableGrpcManagerChannel chan *StreamGrpcManagerRequest
	ctx                             context.Context
	cancel                          context.CancelFunc
}

func (s *staticStreamGrpcManagerPool) GetNextTopicGrpcManager() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	topicManagerRequest := <-s.nextAvailableGrpcManagerChannel
	if topicManagerRequest.Err != nil {
		return nil, topicManagerRequest.Err
	}
	return topicManagerRequest.TopicManager, nil
}

func (s *staticStreamGrpcManagerPool) Close() {
	s.cancel() // Cancel context first to stop goroutine
	close(s.nextAvailableGrpcManagerChannel)
	for _, topicManager := range s.grpcManagers {
		err := topicManager.Close()
		if err != nil {
			s.logger.Error("Error closing topic manager: %v", err)
		}
	}
}

func (s *staticStreamGrpcManagerPool) GetCurrentActiveStreamsCount() uint64 {
	return s.currentActiveStreamsCount.Load()
}

func NewStaticStreamGrpcManagerPool(
	request *models.TopicStreamGrpcManagerRequest,
	numStreamChannels uint32,
	logger logger.MomentoLogger,
) (*staticStreamGrpcManagerPool, momentoerrors.MomentoSvcErr) {
	streamTopicManagers := make([]*grpcmanagers.TopicGrpcManager, 0)
	for i := 0; uint32(i) < numStreamChannels; i++ {
		streamTopicManager, err := grpcmanagers.NewStreamTopicGrpcManager(request)
		if err != nil {
			return nil, err
		}
		streamTopicManagers = append(streamTopicManagers, streamTopicManager)
	}

	// Use an unbuffered channel here so the staticStreamGrpcManagerPool will block on sending
	// the next available stream manager until the most recent request is processed.
	nextAvailableGrpcManagerChannel := make(chan *StreamGrpcManagerRequest)
	ctx, cancel := context.WithCancel(context.Background())

	pool := &staticStreamGrpcManagerPool{
		grpcManagers:                    streamTopicManagers,
		maxConcurrentStreams:            numStreamChannels * uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL),
		logger:                          logger,
		nextAvailableGrpcManagerChannel: nextAvailableGrpcManagerChannel,
		ctx:                             ctx,
		cancel:                          cancel,
	}

	go pool.makeNextManagerAvailable()
	return pool, nil
}

func (s *staticStreamGrpcManagerPool) makeNextManagerAvailable() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			topicManager, err := s.getNextManager()
			select {
			case <-s.ctx.Done():
				return
			case s.nextAvailableGrpcManagerChannel <- &StreamGrpcManagerRequest{
				TopicManager: topicManager,
				Err:          err,
			}:
			}
		}
	}
}

// StreamGrpcManagerRequest is used for putting the next available stream manager on a channel for the
// pubSubClient to pull from, or an error that specifies why no stream manager is available.
type StreamGrpcManagerRequest struct {
	TopicManager *grpcmanagers.TopicGrpcManager
	Err          momentoerrors.MomentoSvcErr
}

// Helper function to help sanity check number of concurrent streams before starting a new subscription
func (s *staticStreamGrpcManagerPool) checkNumConcurrentStreams() momentoerrors.MomentoSvcErr {
	if s.currentActiveStreamsCount.Load() >= uint64(s.maxConcurrentStreams) {
		errorMessage := fmt.Sprintf(
			"Already at maximum number of concurrent grpc streams (%d), cannot make new subscribe requests\n",
			s.maxConcurrentStreams,
		)
		return momentoerrors.NewMomentoSvcErr(momentoerrors.ClientResourceExhaustedError, errorMessage, nil)
	}
	// If we are approaching the grpc maximum concurrent stream limit, log a warning
	remainingStreams := uint64(s.maxConcurrentStreams) - s.currentActiveStreamsCount.Load()
	if remainingStreams < 10 {
		s.logger.Warn(
			"WARNING: approaching grpc maximum concurrent stream limit, %d remaining of total %d streams\n",
			remainingStreams, s.maxConcurrentStreams,
		)
	}
	return nil
}

func (s *staticStreamGrpcManagerPool) getNextManager() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	// First check if there is enough grpc stream capacity to make a new subscription
	err := s.checkNumConcurrentStreams()
	if err != nil {
		return nil, err
	}

	// Max number of attempts is set to the max number of concurrent streams in order to preserve
	// the round-robin system (incrementing nextManagerIndex) but to not cut short the number
	//  of attempts in case there are many subscriptions starting up at the same time.
	for i := 0; uint32(i) < s.maxConcurrentStreams; i++ {
		nextManagerIndex := s.managerIndex.Add(1)
		topicManager := s.grpcManagers[nextManagerIndex%uint64(len(s.grpcManagers))]
		newCount := topicManager.NumActiveSubscriptions.Add(1)
		if newCount <= int64(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL) {
			s.logger.Debug("Starting new subscription on grpc channel %d which now has %d streams", nextManagerIndex%uint64(len(s.grpcManagers)), newCount)
			s.currentActiveStreamsCount.Add(1)
			return topicManager, nil
		}
		topicManager.NumActiveSubscriptions.Add(-1)
	}

	// If there are no more streams available, return an error
	errorMessage := fmt.Sprintf("Cannot start new subscription, all grpc channels may be at maximum capacity. There are %d total subscriptions allowed across %d grpc channels. Please use the WithNumStreamGrpcChannels configuration if you wish to start more subscriptions.\n", s.maxConcurrentStreams, len(s.grpcManagers))
	return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.ClientResourceExhaustedError, errorMessage, nil)
}
