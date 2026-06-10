package topic_manager_lists

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

// staticStreamGrpcManagerPool manages a static pool of gRPC channels for
// stream pubsub requests. Allocation, release, and Close live in the embedded
// streamPoolCore; this type supplies the fixed-size capacity policy.
type staticStreamGrpcManagerPool struct {
	streamPoolCore
	grpcManagers []*grpcmanagers.TopicGrpcManager
	// managerIndex is only touched by getNextManager under the core's mutex.
	managerIndex         uint64
	maxConcurrentStreams uint32
	logger               logger.MomentoLogger
}

// NewStaticStreamGrpcManagerPool creates a new pool with a fixed number of grpc managers for stream pubsub requests.
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

	pool := &staticStreamGrpcManagerPool{
		grpcManagers:         streamTopicManagers,
		maxConcurrentStreams: numStreamChannels * uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL),
		logger:               logger,
	}
	pool.initStreamPoolCore(pool.getNextManager, func() {
		closeAllManagers(pool.grpcManagers, logger)
	})
	return pool, nil
}

// checkNumConcurrentStreams checks the number of concurrent streams before starting a new subscription
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

// getNextManager returns the next available stream manager from the pool,
// reserving a slot on it. Called by streamPoolCore under its mutex.
func (s *staticStreamGrpcManagerPool) getNextManager() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	// First check if there is enough grpc stream capacity to make a new subscription
	err := s.checkNumConcurrentStreams()
	if err != nil {
		return nil, err
	}

	// The attempt bound is generous: allocation is serialized by the core
	// mutex, so one pass over every channel would suffice.
	if topicManager := s.reserveSlot(s.grpcManagers, &s.managerIndex, s.maxConcurrentStreams, s.logger); topicManager != nil {
		return topicManager, nil
	}

	// If there are no more streams available, return an error
	errorMessage := fmt.Sprintf("Cannot start new subscription, all grpc channels may be at maximum capacity. There are %d total subscriptions allowed across %d grpc channels. Please use the WithNumStreamGrpcChannels configuration if you wish to start more subscriptions.\n", s.maxConcurrentStreams, len(s.grpcManagers))
	return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.ClientResourceExhaustedError, errorMessage, nil)
}
