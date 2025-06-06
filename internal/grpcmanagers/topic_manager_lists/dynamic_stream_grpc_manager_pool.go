package topic_manager_lists

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

// dynamicStreamGrpcManagerPool manages a dynamic pool of gRPC channels for stream pubsub requests.
//
// The dynamicStreamGrpcManagerPool will continually place the next available stream manager
// (or an error if no TopicGrpcManagers are available) on the streamManagerRequestQueue channel
// for the pubSubClient to pull from.
// Without this channel, a mutex and additional logic would be needed to protect the dynamic list
// of grpc managers from race conditions due to concurrent access.
type dynamicStreamGrpcManagerPool struct {
	grpcManagers                    []*grpcmanagers.TopicGrpcManager
	managerIndex                    atomic.Uint64
	maxManagerCount                 int    // maximum number of grpc channels that can be created
	currentMaxConcurrentStreams     uint32 // current number of grpc channels * MAX_CONCURRENT_STREAMS_PER_CHANNEL
	currentActiveStreamsCount       atomic.Uint64
	logger                          logger.MomentoLogger
	newTopicManagerProps            *models.TopicStreamGrpcManagerRequest
	nextAvailableGrpcManagerChannel chan *StreamGrpcManagerRequest
	ctx                             context.Context
	cancel                          context.CancelFunc
}

func (d *dynamicStreamGrpcManagerPool) GetNextTopicGrpcManager() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	topicManagerRequest := <-d.nextAvailableGrpcManagerChannel
	if topicManagerRequest.Err != nil {
		return nil, topicManagerRequest.Err
	}
	return topicManagerRequest.TopicManager, nil
}

func (d *dynamicStreamGrpcManagerPool) Close() {
	d.cancel() // Cancel context first to stop goroutines
	close(d.nextAvailableGrpcManagerChannel)
	for _, topicManager := range d.grpcManagers {
		err := topicManager.Close()
		if err != nil {
			d.logger.Error("Error closing topic manager: %s", err.Error())
		}
	}
}

func (d *dynamicStreamGrpcManagerPool) GetCurrentActiveStreamsCount() uint64 {
	return d.currentActiveStreamsCount.Load()
}

func (d *dynamicStreamGrpcManagerPool) GetCurrentNumberOfGrpcManagers() int {
	return len(d.grpcManagers)
}

func NewDynamicStreamGrpcManagerPool(request *models.TopicStreamGrpcManagerRequest, maxSubscriptions uint32, logger logger.MomentoLogger) (*dynamicStreamGrpcManagerPool, momentoerrors.MomentoSvcErr) {
	// make just one manager to start with
	streamTopicManagers := make([]*grpcmanagers.TopicGrpcManager, 0)
	streamTopicManager, err := grpcmanagers.NewStreamTopicGrpcManager(request)
	if err != nil {
		return nil, err
	}
	streamTopicManagers = append(streamTopicManagers, streamTopicManager)
	logger.Debug("Max subscriptions: %d, max manager count: %d", maxSubscriptions, int(math.Ceil(float64(maxSubscriptions)/float64(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL))))

	// Unbuffered channel so the stream manager list will block on sending the next
	// available stream manager until the most recent request is processed.
	nextAvailableGrpcManagerChannel := make(chan *StreamGrpcManagerRequest)
	ctx, cancel := context.WithCancel(context.Background())

	pool := &dynamicStreamGrpcManagerPool{
		grpcManagers:                    streamTopicManagers,
		maxManagerCount:                 int(math.Ceil(float64(maxSubscriptions) / float64(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL))),
		currentMaxConcurrentStreams:     uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL), // for one channel
		logger:                          logger,
		newTopicManagerProps:            request,
		nextAvailableGrpcManagerChannel: nextAvailableGrpcManagerChannel,
		ctx:                             ctx,
		cancel:                          cancel,
	}

	// Start goroutine to continually make the next available stream manager
	// available on the streamManagerRequestQueue
	go pool.makeNextManagerAvailable()

	return pool, nil
}

func (d *dynamicStreamGrpcManagerPool) makeNextManagerAvailable() {
	for {
		select {
		case <-d.ctx.Done():
			return
		default:
			topicManager, err := d.getNextManager()
			select {
			case <-d.ctx.Done():
				return
			case d.nextAvailableGrpcManagerChannel <- &StreamGrpcManagerRequest{
				TopicManager: topicManager,
				Err:          err,
			}:
			}
		}
	}
}

func (d *dynamicStreamGrpcManagerPool) getNextManager() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	// First check if there is enough grpc stream capacity to make a new subscription
	err := d.checkNumConcurrentStreams()
	if err != nil {
		return nil, err
	}

	// Max number of attempts is set to the max number of concurrent streams in order to preserve
	// the round-robin system (incrementing nextManagerIndex) but to not cut short the number
	//  of attempts in case there are many subscriptions starting up at the same time.
	for i := 0; uint32(i) < d.currentMaxConcurrentStreams; i++ {
		nextManagerIndex := d.managerIndex.Add(1)
		topicManager := d.grpcManagers[nextManagerIndex%uint64(len(d.grpcManagers))]
		newCount := topicManager.NumActiveSubscriptions.Add(1)
		if newCount <= int64(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL) {
			d.logger.Debug("Starting new subscription on grpc channel %d which now has %d streams", nextManagerIndex%uint64(len(d.grpcManagers)), newCount)
			d.currentActiveStreamsCount.Add(1)
			return topicManager, nil
		}
		topicManager.NumActiveSubscriptions.Add(-1)
	}

	// If there are no more streams available, return an error
	errorMessage := fmt.Sprintf("Cannot start new subscription, all grpc channels may be at maximum capacity. There are %d total subscriptions allowed across %d grpc channels. Please use the WithMaxSubscriptions configuration if you wish to start more subscriptions.\n", d.currentMaxConcurrentStreams, d.maxManagerCount)
	return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.ClientResourceExhaustedError, errorMessage, nil)
}

// Helper function to help sanity check number of concurrent streams before starting a new subscription
func (d *dynamicStreamGrpcManagerPool) checkNumConcurrentStreams() momentoerrors.MomentoSvcErr {
	numActiveStreams := d.currentActiveStreamsCount.Load()
	d.logger.Debug("Current number of active subscriptions: %d", d.currentActiveStreamsCount.Load())

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

// Called by checkNumConcurrentStreams to add more stream capacity if needed.
func (d *dynamicStreamGrpcManagerPool) addManager() momentoerrors.MomentoSvcErr {
	streamTopicManager, err := grpcmanagers.NewStreamTopicGrpcManager(d.newTopicManagerProps)
	if err != nil {
		return err
	}
	d.grpcManagers = append(d.grpcManagers, streamTopicManager)
	d.currentMaxConcurrentStreams = uint32(len(d.grpcManagers)) * uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
	return nil
}
