package grpcmanagers

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/interceptor"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc"
)

type TopicGrpcManager struct {
	Conn                   *grpc.ClientConn
	StreamClient           pb.PubsubClient
	NumActiveSubscriptions atomic.Int64
}

func NewStreamTopicGrpcManager(request *models.TopicStreamGrpcManagerRequest) (*TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	endpoint := request.CredentialProvider.GetCacheEndpoint()
	authToken := request.CredentialProvider.GetAuthToken()

	headerInterceptors := []grpc.StreamClientInterceptor{
		interceptor.AddStreamHeaderInterceptor(authToken),
	}

	conn, err := grpc.NewClient(
		endpoint,
		AllDialOptions(
			request.GrpcConfiguration,
			request.CredentialProvider.IsCacheEndpointSecure(),
			grpc.WithChainStreamInterceptor(headerInterceptors...),
			grpc.WithChainUnaryInterceptor(interceptor.AddAuthHeadersInterceptor(authToken)),
		)...,
	)

	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &TopicGrpcManager{
		Conn:         conn,
		StreamClient: pb.NewPubsubClient(conn),
	}, nil
}

func (topicManager *TopicGrpcManager) Close() momentoerrors.MomentoSvcErr {
	topicManager.NumActiveSubscriptions.Store(0)
	return momentoerrors.ConvertSvcErr(topicManager.Conn.Close())
}

// Implemented by StaticUnaryManagerList
type TopicManagerList interface {
	GetNextManager() (*TopicGrpcManager, momentoerrors.MomentoSvcErr)
	Close()
}

// Implemented by StaticStreamManagerList and DynamicStreamManagerList
type TopicStreamManagerListWithBookkeeping interface {
	TopicManagerList
	CountNumberOfActiveSubscriptions() int64
	checkNumConcurrentStreams() momentoerrors.MomentoSvcErr
}

// Most basic case: static list of unary grpc managers
type StaticUnaryManagerList struct {
	grpcManagers []*TopicGrpcManager
	managerIndex atomic.Uint64
	logger       logger.MomentoLogger
}

// No special checks needed for getting next manager, no possible error returned.
//
// Each grpc connection can multiplex 100 subscribe/publish requests.
// Publish requests will queue up on client while waiting for in-flight requests to complete if
// the number of concurrent requests exceeds numUnaryChannels*100, but will eventually complete.
// Therefore we can just round-robin the unaryTopicManagers, no need to keep track of how many
// publish requests are in flight on each one.
func (list *StaticUnaryManagerList) GetNextManager() (*TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	nextManagerIndex := list.managerIndex.Add(1)
	return list.grpcManagers[nextManagerIndex%uint64(len(list.grpcManagers))], nil
}

func (list *StaticUnaryManagerList) Close() {
	for _, topicManager := range list.grpcManagers {
		err := topicManager.Close()
		if err != nil {
			list.logger.Error("Error closing topic manager: %v", err)
		}
	}
}

func NewStaticUnaryManagerList(request *models.TopicStreamGrpcManagerRequest, numUnaryChannels uint32, logger logger.MomentoLogger) (*StaticUnaryManagerList, momentoerrors.MomentoSvcErr) {
	unaryTopicManagers := make([]*TopicGrpcManager, 0)
	for i := 0; uint32(i) < numUnaryChannels; i++ {
		unaryTopicManager, err := NewStreamTopicGrpcManager(request)
		if err != nil {
			return nil, err
		}
		unaryTopicManagers = append(unaryTopicManagers, unaryTopicManager)
	}
	return &StaticUnaryManagerList{
		grpcManagers: unaryTopicManagers,
		logger:       logger,
	}, nil
}

// Middle case: static list of stream grpc managers
type StaticStreamManagerList struct {
	grpcManagers         []*TopicGrpcManager
	managerIndex         atomic.Uint64
	maxConcurrentStreams uint32
	logger               logger.MomentoLogger
}

// Must check that there is a stream available to return, otherwise return an error.
//
// Each grpc connection can multiplex 100 subscribe/publish requests.
// Grpc channels that already have 100 subscriptions will silently queue up subsequent requests.
// We must prevent subscription requests from queuing up if there are already numStreamChannels*100
// concurrent streams as it causes the program to hang indefinitely with no error.
func (list *StaticStreamManagerList) GetNextManager() (*TopicGrpcManager, momentoerrors.MomentoSvcErr) {
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
		if newCount <= int64(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL) {
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
	for _, topicManager := range list.grpcManagers {
		err := topicManager.Close()
		if err != nil {
			list.logger.Error("Error closing topic manager: %v", err)
		}
	}
}

// Helper to count number of active subscriptions on a static stream manager list
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

func NewStaticStreamManagerList(request *models.TopicStreamGrpcManagerRequest, numStreamChannels uint32, logger logger.MomentoLogger) (*StaticStreamManagerList, momentoerrors.MomentoSvcErr) {
	streamTopicManagers := make([]*TopicGrpcManager, 0)
	for i := 0; uint32(i) < numStreamChannels; i++ {
		streamTopicManager, err := NewStreamTopicGrpcManager(request)
		if err != nil {
			return nil, err
		}
		streamTopicManagers = append(streamTopicManagers, streamTopicManager)
	}
	return &StaticStreamManagerList{
		grpcManagers:         streamTopicManagers,
		maxConcurrentStreams: numStreamChannels * config.MAX_CONCURRENT_STREAMS_PER_CHANNEL,
		logger:               logger,
	}, nil
}

// Most complex case: dynamic list of stream grpc managers
type DynamicStreamManagerList struct {
	grpcManagers                []*TopicGrpcManager
	managerIndex                atomic.Uint64
	maxManagerCount             int    // maximum number of grpc channels that can be created
	currentMaxConcurrentStreams uint32 // current number of grpc channels * MAX_CONCURRENT_STREAMS_PER_CHANNEL
	logger                      logger.MomentoLogger
	rwLock                      sync.RWMutex
	newTopicManagerProps        *models.TopicStreamGrpcManagerRequest
}

// If current max streams is reached, check if we can add a new channel and return a stream, else error.
func (list *DynamicStreamManagerList) GetNextManager() (*TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	// First check if there is enough grpc stream capacity to make a new subscription
	err := list.checkNumConcurrentStreams()
	if err != nil {
		return nil, err
	}

	list.rwLock.RLock()
	defer list.rwLock.RUnlock()

	// Max number of attempts is set to the max number of concurrent streams in order to preserve
	// the round-robin system (incrementing nextManagerIndex) but to not cut short the number
	//  of attempts in case there are many subscriptions starting up at the same time.
	for i := 0; uint32(i) < list.currentMaxConcurrentStreams; i++ {
		nextManagerIndex := list.managerIndex.Add(1)
		topicManager := list.grpcManagers[nextManagerIndex%uint64(len(list.grpcManagers))]
		newCount := topicManager.NumActiveSubscriptions.Add(1)
		if newCount <= int64(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL) {
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
	for _, topicManager := range list.grpcManagers {
		err := topicManager.Close()
		if err != nil {
			list.logger.Error("Error closing topic manager: %v", err)
		}
	}
}

// Called by checkNumConcurrentStreams, which holds a write lock.
func (list *DynamicStreamManagerList) CountNumberOfActiveSubscriptions() int64 {
	count := int64(0)
	for _, topicManager := range list.grpcManagers {
		count += topicManager.NumActiveSubscriptions.Load()
	}
	return count
}

// Helper function to help sanity check number of concurrent streams before starting a new subscription
func (list *DynamicStreamManagerList) checkNumConcurrentStreams() momentoerrors.MomentoSvcErr {
	list.rwLock.Lock()
	defer list.rwLock.Unlock()

	numActiveStreams := list.CountNumberOfActiveSubscriptions()
	list.logger.Debug("Current number of active subscriptions: %d", numActiveStreams)

	if numActiveStreams >= int64(list.currentMaxConcurrentStreams) && len(list.grpcManagers) == list.maxManagerCount {
		errorMessage := fmt.Sprintf(
			"Number of grpc streams: %d; number of channels: %d; max concurrent streams: %d; Already at maximum number of concurrent grpc streams, cannot make new subscribe requests\n",
			numActiveStreams, list.maxManagerCount, list.currentMaxConcurrentStreams,
		)
		return momentoerrors.NewMomentoSvcErr(momentoerrors.LimitExceededError, errorMessage, nil)
	} else if numActiveStreams >= int64(list.currentMaxConcurrentStreams) && len(list.grpcManagers) < list.maxManagerCount {
		// otherwise we can try to add a new manager
		err := list.addManager()
		if err != nil {
			return err
		}
		list.logger.Debug("Added new manager, current number of managers: %d", len(list.grpcManagers))
	}

	// If we are approaching the grpc maximum concurrent stream limit, log a warning
	if len(list.grpcManagers) == list.maxManagerCount {
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

// called by checkNumConcurrentStreams, which holds a write lock
func (list *DynamicStreamManagerList) addManager() momentoerrors.MomentoSvcErr {
	streamTopicManager, err := NewStreamTopicGrpcManager(list.newTopicManagerProps)
	if err != nil {
		return err
	}
	list.grpcManagers = append(list.grpcManagers, streamTopicManager)
	list.currentMaxConcurrentStreams = uint32(len(list.grpcManagers)) * config.MAX_CONCURRENT_STREAMS_PER_CHANNEL
	return nil
}

func NewDynamicStreamManagerList(request *models.TopicStreamGrpcManagerRequest, maxSubscriptions uint32, logger logger.MomentoLogger) (*DynamicStreamManagerList, momentoerrors.MomentoSvcErr) {
	// make just one manager to start with
	streamTopicManagers := make([]*TopicGrpcManager, 0)
	streamTopicManager, err := NewStreamTopicGrpcManager(request)
	if err != nil {
		return nil, err
	}
	streamTopicManagers = append(streamTopicManagers, streamTopicManager)
	logger.Debug("Max subscriptions: %d, max manager count: %d", maxSubscriptions, int(math.Ceil(float64(maxSubscriptions)/float64(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL))))
	return &DynamicStreamManagerList{
		grpcManagers:                streamTopicManagers,
		maxManagerCount:             int(math.Ceil(float64(maxSubscriptions) / float64(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL))),
		currentMaxConcurrentStreams: config.MAX_CONCURRENT_STREAMS_PER_CHANNEL, // for one channel
		logger:                      logger,
		newTopicManagerProps:        request,
	}, nil
}
