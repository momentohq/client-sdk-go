package momento

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const DEFAULT_NUM_STREAM_GRPC_CHANNELS uint32 = 4
const DEFAULT_NUM_UNARY_GRPC_CHANNELS uint32 = 4
const MAX_CONCURRENT_STREAMS_PER_CHANNEL int = 100

type pubSubClient struct {
	numUnaryChannels        uint32
	unaryTopicManagers      []*grpcmanagers.TopicGrpcManager
	unaryTopicManagerCount  atomic.Uint64
	numStreamChannels       uint32
	streamTopicManagers     []*grpcmanagers.TopicGrpcManager
	streamTopicManagerCount atomic.Uint64
	endpoint                string
	log                     logger.MomentoLogger
	requestTimeout          time.Duration
	maxConcurrentStreams    int
}

func newPubSubClient(request *models.PubSubClientRequest) (*pubSubClient, momentoerrors.MomentoSvcErr) {
	grpcConfig := request.TopicsConfiguration.GetTransportStrategy().GetGrpcConfig()

	numStreamChannels := uint32(DEFAULT_NUM_STREAM_GRPC_CHANNELS)
	if request.TopicsConfiguration.GetNumStreamGrpcChannels() > 0 {
		numStreamChannels = request.TopicsConfiguration.GetNumStreamGrpcChannels()
	}
	streamTopicManagers := make([]*grpcmanagers.TopicGrpcManager, 0)
	for i := 0; uint32(i) < numStreamChannels; i++ {
		streamTopicManager, err := grpcmanagers.NewStreamTopicGrpcManager(&models.TopicStreamGrpcManagerRequest{
			CredentialProvider: request.CredentialProvider,
			GrpcConfiguration:  grpcConfig,
		})
		if err != nil {
			return nil, err
		}
		streamTopicManagers = append(streamTopicManagers, streamTopicManager)
	}

	numUnaryChannels := uint32(DEFAULT_NUM_UNARY_GRPC_CHANNELS)
	if request.TopicsConfiguration.GetNumUnaryGrpcChannels() > 0 {
		numUnaryChannels = request.TopicsConfiguration.GetNumUnaryGrpcChannels()
	}
	unaryTopicManagers := make([]*grpcmanagers.TopicGrpcManager, 0)
	for i := 0; uint32(i) < numUnaryChannels; i++ {
		unaryTopicManager, err := grpcmanagers.NewStreamTopicGrpcManager(&models.TopicStreamGrpcManagerRequest{
			CredentialProvider: request.CredentialProvider,
			GrpcConfiguration:  grpcConfig,
		})
		if err != nil {
			return nil, err
		}
		unaryTopicManagers = append(unaryTopicManagers, unaryTopicManager)
	}

	var timeout time.Duration
	if request.TopicsConfiguration.GetClientSideTimeout() < 1 {
		timeout = defaultRequestTimeout
	} else {
		timeout = request.TopicsConfiguration.GetClientSideTimeout()
	}

	return &pubSubClient{
		numUnaryChannels:     numUnaryChannels,
		unaryTopicManagers:   unaryTopicManagers,
		numStreamChannels:    numStreamChannels,
		streamTopicManagers:  streamTopicManagers,
		endpoint:             request.CredentialProvider.GetCacheEndpoint(),
		log:                  request.Log,
		requestTimeout:       timeout,
		maxConcurrentStreams: int(numStreamChannels) * MAX_CONCURRENT_STREAMS_PER_CHANNEL,
	}, nil
}

// Each grpc connection can multiplex 100 subscribe/publish requests.
// Grpc channels that already have 100 subscriptions will silently queue up subsequent requests.
// We must prevent subscription requests from queuing up if there are already numStreamChannels*100
// concurrent streams as it causes the program to hang indefinitely with no error.
func (client *pubSubClient) getNextStreamTopicManager() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	// Max number of attempts is set to the max number of concurrent streams in order to preserve
	// the round-robin system (incrementing nextManagerIndex) but to not cut short the number
	//  of attempts in case there are many subscriptions starting up at the same time.
	for i := 0; i < client.maxConcurrentStreams; i++ {
		nextManagerIndex := client.streamTopicManagerCount.Add(1)
		topicManager := client.streamTopicManagers[nextManagerIndex%uint64(len(client.streamTopicManagers))]
		newCount := topicManager.NumActiveSubscriptions.Add(1)
		if newCount <= int64(MAX_CONCURRENT_STREAMS_PER_CHANNEL) {
			client.log.Debug("Starting new subscription on grpc channel %d which now has %d streams", nextManagerIndex%uint64(len(client.streamTopicManagers)), newCount)
			return topicManager, nil
		}
		topicManager.NumActiveSubscriptions.Add(-1)
	}

	// If there are no more streams available, return an error
	errorMessage := fmt.Sprintf("Cannot start new subscription, all grpc channels may be at maximum capacity. There are %d total subscriptions allowed across %d grpc channels. Please use the WithNumStreamGrpcChannels configuration if you wish to start more subscriptions.\n", client.maxConcurrentStreams, client.numStreamChannels)
	return nil, momentoerrors.NewMomentoSvcErr(LimitExceededError, errorMessage, nil)
}

// Each grpc connection can multiplex 100 subscribe/publish requests.
// Publish requests will queue up on client while waiting for in-flight requests to complete if
// the number of concurrent requests exceeds numUnaryChannels*100, but will eventually complete.
// Therefore we can just round-robin the unaryTopicManagers, no need to keep track of how many
// publish requests are in flight on each one.
func (client *pubSubClient) getNextUnaryTopicManager() *grpcmanagers.TopicGrpcManager {
	nextManagerIndex := client.unaryTopicManagerCount.Add(1)
	topicManager := client.unaryTopicManagers[nextManagerIndex%uint64(client.numUnaryChannels)]
	return topicManager
}

func (client *pubSubClient) countNumberOfActiveSubscriptions() int64 {
	count := int64(0)
	for _, topicManager := range client.streamTopicManagers {
		count += topicManager.NumActiveSubscriptions.Load()
	}
	return count
}

// Helper function to help sanity check number of concurrent streams before starting a new subscription
func (client *pubSubClient) checkNumConcurrentStreams() error {
	numActiveStreams := client.countNumberOfActiveSubscriptions()
	if numActiveStreams >= int64(client.maxConcurrentStreams) {
		errorMessage := fmt.Sprintf(
			"Number of grpc streams: %d; number of channels: %d; max concurrent streams: %d; Already at maximum number of concurrent grpc streams, cannot make new subscribe requests\n",
			numActiveStreams, client.numStreamChannels, client.maxConcurrentStreams,
		)
		return momentoerrors.NewMomentoSvcErr(LimitExceededError, errorMessage, nil)
	}
	return nil
}

func (client *pubSubClient) topicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (*grpcmanagers.TopicGrpcManager, pb.Pubsub_SubscribeClient, context.Context, context.CancelFunc, error) {
	// First check if there is enough grpc stream capabity to make a new subscription
	err := client.checkNumConcurrentStreams()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Then actually attempt to get a topic manager
	topicManager, err := client.getNextStreamTopicManager()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// add metadata to context
	requestMetadata := internal.CreateMetadata(ctx, internal.Topic)

	// add withCancel to context
	cancelContext, cancelFunction := context.WithCancel(requestMetadata)

	var header, trailer metadata.MD
	subscribeClient, err := topicManager.StreamClient.Subscribe(cancelContext, &pb.XSubscriptionRequest{
		CacheName:                   request.CacheName,
		Topic:                       request.TopicName,
		ResumeAtTopicSequenceNumber: request.ResumeAtTopicSequenceNumber,
		SequencePage:                request.SequencePage,
	})

	if err != nil {
		topicManager.NumActiveSubscriptions.Add(-1)
		cancelFunction()
		if subscribeClient != nil {
			header, _ = subscribeClient.Header()
			trailer = subscribeClient.Trailer()
		}
		return nil, nil, nil, nil, momentoerrors.ConvertSvcErr(err, header, trailer)
	}

	// If we are approaching the grpc maximum concurrent stream limit, log a warning
	numActiveStreams := client.countNumberOfActiveSubscriptions()
	remainingStreams := int64(client.maxConcurrentStreams) - numActiveStreams
	if remainingStreams < 10 {
		client.log.Warn(
			"WARNING: approaching grpc maximum concurrent stream limit, %d remaining of total %d streams\n",
			remainingStreams, client.maxConcurrentStreams,
		)
	}
	client.log.Debug("Starting new subscription, total number of streams now: %d", numActiveStreams)
	return topicManager, subscribeClient, cancelContext, cancelFunction, err
}

func (client *pubSubClient) topicPublish(ctx context.Context, request *TopicPublishRequest) error {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := internal.CreateMetadata(ctx, internal.Topic)
	topicManager := client.getNextUnaryTopicManager()
	var header, trailer metadata.MD
	switch value := request.Value.(type) {
	case String:
		_, err := topicManager.StreamClient.Publish(requestMetadata, &pb.XPublishRequest{
			CacheName: request.CacheName,
			Topic:     request.TopicName,
			Value: &pb.XTopicValue{
				Kind: &pb.XTopicValue_Text{
					Text: value.asString(),
				},
			},
		}, grpc.Header(&header), grpc.Trailer(&trailer))
		if err != nil {
			return momentoerrors.ConvertSvcErr(err, header, trailer)
		}
		return err
	case Bytes:
		_, err := topicManager.StreamClient.Publish(requestMetadata, &pb.XPublishRequest{
			CacheName: request.CacheName,
			Topic:     request.TopicName,
			Value: &pb.XTopicValue{
				Kind: &pb.XTopicValue_Binary{
					Binary: value.asBytes(),
				},
			},
		}, grpc.Header(&header), grpc.Trailer(&trailer))
		if err != nil {
			return momentoerrors.ConvertSvcErr(err, header, trailer)
		}
		return err
	default:
		return momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"error encoding topic value only support []byte or string currently", nil,
		)
	}
}

func (client *pubSubClient) close() {
	// Close all stream grpc channels
	client.numStreamChannels = 0
	for clientIndex := range client.streamTopicManagers {
		defer client.streamTopicManagers[clientIndex].Close()
	}

	// Close all unary grpc channels
	client.numUnaryChannels = 0
	for clientIndex := range client.unaryTopicManagers {
		defer client.unaryTopicManagers[clientIndex].Close()
	}
}
