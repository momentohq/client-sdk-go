package topic_manager_lists

import (
	"fmt"
	"math"

	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

// dynamicStreamGrpcManagerPool manages a dynamic pool of gRPC channels for
// stream pubsub requests. Allocation, release, and Close live in the embedded
// streamPoolCore; this type supplies the grow-on-demand capacity policy.
//
// grpcManagers and currentMaxConcurrentStreams are only mutated by addManager
// under the core's mutex.
type dynamicStreamGrpcManagerPool struct {
	streamPoolCore
	grpcManagers []*grpcmanagers.TopicGrpcManager
	// managerIndex is only touched by getNextManager under the core's mutex.
	managerIndex                uint64
	maxManagerCount             int    // max grpc channels
	currentMaxConcurrentStreams uint32 // grpc channels * MAX_CONCURRENT_STREAMS_PER_CHANNEL
	logger                      logger.MomentoLogger
	newTopicManagerProps        *models.TopicStreamGrpcManagerRequest
}

// NewDynamicStreamGrpcManagerPool creates a new pool with a dynamic number of grpc managers for stream pubsub requests.
//
// Defaults to one grpc manager to start with and will dynamically add more as needed up until a max of maxSubscriptions/100
// since each grpc manager can handle 100 concurrent streams.
func NewDynamicStreamGrpcManagerPool(request *models.TopicStreamGrpcManagerRequest, maxSubscriptions uint32, logger logger.MomentoLogger) (*dynamicStreamGrpcManagerPool, momentoerrors.MomentoSvcErr) {
	// make just one manager to start with
	streamTopicManagers := make([]*grpcmanagers.TopicGrpcManager, 0)
	streamTopicManager, err := grpcmanagers.NewStreamTopicGrpcManager(request)
	if err != nil {
		return nil, err
	}
	streamTopicManagers = append(streamTopicManagers, streamTopicManager)
	logger.Debug("Max subscriptions: %d, max manager count: %d", maxSubscriptions, int(math.Ceil(float64(maxSubscriptions)/float64(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL))))

	pool := &dynamicStreamGrpcManagerPool{
		grpcManagers:                streamTopicManagers,
		maxManagerCount:             int(math.Ceil(float64(maxSubscriptions) / float64(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL))),
		currentMaxConcurrentStreams: uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL), // for one channel
		logger:                      logger,
		newTopicManagerProps:        request,
	}
	pool.initStreamPoolCore(pool.getNextManager, func() {
		closeAllManagers(pool.grpcManagers, logger)
	})
	return pool, nil
}

// GetCurrentNumberOfGrpcManagers returns the current number of grpc managers in the pool.
func (d *dynamicStreamGrpcManagerPool) GetCurrentNumberOfGrpcManagers() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.grpcManagers)
}

// getNextManager returns the next available stream manager from the pool,
// reserving a slot on it and growing the pool if needed. Called by
// streamPoolCore under its mutex.
func (d *dynamicStreamGrpcManagerPool) getNextManager() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	// First check if there is enough grpc stream capacity to make a new subscription
	err := d.checkNumConcurrentStreams()
	if err != nil {
		return nil, err
	}

	// The attempt bound is generous: allocation is serialized by the core
	// mutex, so one pass over every channel would suffice.
	if topicManager := d.reserveSlot(d.grpcManagers, &d.managerIndex, d.currentMaxConcurrentStreams, d.logger); topicManager != nil {
		return topicManager, nil
	}

	// If there are no more streams available, return an error
	errorMessage := fmt.Sprintf("Cannot start new subscription, all grpc channels may be at maximum capacity. There are %d total subscriptions allowed across %d grpc channels. Please use the WithMaxSubscriptions configuration if you wish to start more subscriptions.\n", d.currentMaxConcurrentStreams, d.maxManagerCount)
	return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.ClientResourceExhaustedError, errorMessage, nil)
}

// checkNumConcurrentStreams checks the number of concurrent streams before starting a new subscription.
// If the current maximum number of concurrent streams is reached but the maximum number of grpc managers
// has not been reached, it will add a new grpc manager to the pool.
// If both maximums have been reached, it will return a ClientResourceExhaustedError.
func (d *dynamicStreamGrpcManagerPool) checkNumConcurrentStreams() momentoerrors.MomentoSvcErr {
	numActiveStreams := d.currentActiveStreamsCount.Load()
	d.logger.Debug("Current number of active subscriptions: %d", numActiveStreams)

	numStreamManagers := len(d.grpcManagers)

	if numActiveStreams >= uint64(d.currentMaxConcurrentStreams) && numStreamManagers >= d.maxManagerCount {
		errorMessage := fmt.Sprintf(
			"Already at maximum number of concurrent grpc streams (%d), cannot make new subscribe requests\n",
			d.currentMaxConcurrentStreams,
		)
		return momentoerrors.NewMomentoSvcErr(momentoerrors.ClientResourceExhaustedError, errorMessage, nil)
	} else if numActiveStreams >= uint64(d.currentMaxConcurrentStreams) && numStreamManagers < d.maxManagerCount {
		// otherwise we can try to add a new manager
		err := d.addManager()
		if err != nil {
			return err
		}
		d.logger.Debug("Added new manager, current number of managers: %d", len(d.grpcManagers))
	}

	// If we are approaching the grpc maximum concurrent stream limit, log a warning
	if numStreamManagers == d.maxManagerCount {
		remainingStreams := uint64(d.currentMaxConcurrentStreams) - numActiveStreams
		if remainingStreams < 10 {
			d.logger.Warn(
				"WARNING: approaching grpc maximum concurrent stream limit, %d remaining of total %d streams\n",
				remainingStreams, d.currentMaxConcurrentStreams,
			)
		}
	}
	return nil
}

// addManager is called by checkNumConcurrentStreams to add more stream capacity by
// adding a new grpc manager to the pool.
func (d *dynamicStreamGrpcManagerPool) addManager() momentoerrors.MomentoSvcErr {
	streamTopicManager, err := grpcmanagers.NewStreamTopicGrpcManager(d.newTopicManagerProps)
	if err != nil {
		return err
	}
	d.grpcManagers = append(d.grpcManagers, streamTopicManager)
	d.currentMaxConcurrentStreams = uint32(len(d.grpcManagers)) * uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
	return nil
}
