package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/internal"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers/topic_manager_lists"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const DEFAULT_NUM_STREAM_GRPC_CHANNELS uint32 = 4
const DEFAULT_NUM_UNARY_GRPC_CHANNELS uint32 = 4

// Deprecated: use config.MAX_CONCURRENT_STREAMS_PER_CHANNEL instead
const MAX_CONCURRENT_STREAMS_PER_CHANNEL int = 100

type pubSubClient struct {
	endpoint                 string
	log                      logger.MomentoLogger
	requestTimeout           time.Duration
	middleware               []middleware.TopicMiddleware
	unaryGrpcConnectionPool  topic_manager_lists.TopicGrpcConnectionPool
	streamGrpcConnectionPool topic_manager_lists.TopicGrpcConnectionPool
}

func newPubSubClient(request *models.PubSubClientRequest) (*pubSubClient, momentoerrors.MomentoSvcErr) {
	var timeout time.Duration
	if request.TopicsConfiguration.GetClientSideTimeout() < 1 {
		timeout = defaultRequestTimeout
	} else {
		timeout = request.TopicsConfiguration.GetClientSideTimeout()
	}

	grpcConfig := request.TopicsConfiguration.GetTransportStrategy().GetGrpcConfig()
	topicManagerProps := &models.TopicStreamGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
		GrpcConfiguration:  grpcConfig,
	}

	// Create pool of grpc channels for unary operations
	numUnaryChannels := DEFAULT_NUM_UNARY_GRPC_CHANNELS
	if request.TopicsConfiguration.GetNumUnaryGrpcChannels() > 0 {
		numUnaryChannels = request.TopicsConfiguration.GetNumUnaryGrpcChannels()
	}
	unaryPool, err := topic_manager_lists.NewStaticUnaryGrpcManagerPool(
		topicManagerProps,
		numUnaryChannels,
		request.Log,
	)
	if err != nil {
		return nil, err
	}

	// Create pool of grpc channels for stream operations depending on static vs dynamic transport strategy
	var streamPool topic_manager_lists.TopicGrpcConnectionPool
	if request.TopicsConfiguration.GetMaxSubscriptions() > 0 {
		request.Log.Debug("Creating dynamic stream manager list with max subscriptions: %d", request.TopicsConfiguration.GetMaxSubscriptions())
		streamPool, err = topic_manager_lists.NewDynamicStreamGrpcManagerPool(
			topicManagerProps,
			request.TopicsConfiguration.GetMaxSubscriptions(),
			request.Log,
		)
	} else {
		numStreamChannels := DEFAULT_NUM_STREAM_GRPC_CHANNELS
		if request.TopicsConfiguration.GetNumStreamGrpcChannels() > 0 {
			numStreamChannels = request.TopicsConfiguration.GetNumStreamGrpcChannels()
		}
		request.Log.Debug("Creating static stream manager list with num stream channels: %d", numStreamChannels)
		streamPool, err = topic_manager_lists.NewStaticStreamGrpcManagerPool(
			topicManagerProps,
			numStreamChannels,
			request.Log,
		)
	}
	if err != nil {
		return nil, err
	}

	return &pubSubClient{
		endpoint:                 request.CredentialProvider.GetCacheEndpoint(),
		log:                      request.Log,
		requestTimeout:           timeout,
		middleware:               request.TopicsConfiguration.GetMiddleware(),
		unaryGrpcConnectionPool:  unaryPool,
		streamGrpcConnectionPool: streamPool,
	}, nil
}

func (client *pubSubClient) topicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (*grpcmanagers.TopicGrpcManager, pb.Pubsub_SubscribeClient, context.Context, context.CancelFunc, error) {
	// Get the next available grpc manager
	topicManager, topicManagerErr := client.streamGrpcConnectionPool.GetNextTopicGrpcManager()
	if topicManagerErr != nil {
		return nil, nil, nil, nil, topicManagerErr
	}

	subscriptionRequest := &pb.XSubscriptionRequest{
		CacheName:                   request.CacheName,
		Topic:                       request.TopicName,
		ResumeAtTopicSequenceNumber: request.ResumeAtTopicSequenceNumber,
		SequencePage:                request.SequencePage,
	}

	requestMetadata := make(map[string]string)
	for _, mw := range client.middleware {
		newMd := mw.OnSubscribeMetadata(deepCopyMap(requestMetadata))
		if newMd != nil {
			requestMetadata = newMd
		}
	}

	// add withCancel to context
	cancelContext, cancelFunction := context.WithCancel(ctx)
	requestContext := internal.CreateTopicRequestContextFromMetadataMap(cancelContext, request.CacheName, requestMetadata)

	var header, trailer metadata.MD
	subscribeClient, err := topicManager.StreamClient.Subscribe(requestContext, subscriptionRequest)

	if err != nil {
		topicManager.NumActiveSubscriptions.Add(-1)
		cancelFunction()
		if subscribeClient != nil {
			header, _ = subscribeClient.Header()
			trailer = subscribeClient.Trailer()
		}
		return nil, nil, nil, nil, momentoerrors.ConvertSvcErr(err, header, trailer)
	}

	return topicManager, subscribeClient, cancelContext, cancelFunction, err
}

func (client *pubSubClient) topicPublish(ctx context.Context, request *TopicPublishRequest) error {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := make(map[string]string)
	for _, mw := range client.middleware {
		newMd := mw.OnPublishMetadata(deepCopyMap(requestMetadata))
		if newMd != nil {
			requestMetadata = newMd
		}
	}

	requestContext := internal.CreateTopicRequestContextFromMetadataMap(ctx, request.CacheName, requestMetadata)
	topicManager, err := client.unaryGrpcConnectionPool.GetNextTopicGrpcManager()
	if err != nil {
		return err
	}
	var header, trailer metadata.MD
	switch value := request.Value.(type) {
	case String:
		_, err := topicManager.StreamClient.Publish(requestContext, &pb.XPublishRequest{
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
		_, err := topicManager.StreamClient.Publish(requestContext, &pb.XPublishRequest{
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
	client.streamGrpcConnectionPool.Close()
	client.unaryGrpcConnectionPool.Close()
}
