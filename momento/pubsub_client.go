package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/middleware"
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

// Internal interface for pubsub client.
type pubSubClient interface {
	topicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (*grpcmanagers.TopicGrpcManager, pb.Pubsub_SubscribeClient, context.Context, context.CancelFunc, error)
	topicPublish(ctx context.Context, request *TopicPublishRequest) error
	close()
	getMiddlewares() []middleware.TopicMiddleware
	countNumberOfActiveSubscriptions() int64
}

// Pubsub client implementation that uses a static number of GRPC channels for subscriptions and unary operations.
type staticPubSubClient struct {
	endpoint          string
	log               logger.MomentoLogger
	requestTimeout    time.Duration
	middleware        []middleware.TopicMiddleware
	unaryManagerList  *grpcmanagers.StaticUnaryManagerList
	streamManagerList *grpcmanagers.StaticStreamManagerList
}

func newStaticPubSubClient(request *models.PubSubClientRequest) (*staticPubSubClient, momentoerrors.MomentoSvcErr) {
	grpcConfig := request.TopicsConfiguration.GetGrpcConfig()

	numStreamChannels := DEFAULT_NUM_STREAM_GRPC_CHANNELS
	if request.TopicsConfiguration.GetNumStreamGrpcChannels() > 0 {
		numStreamChannels = request.TopicsConfiguration.GetNumStreamGrpcChannels()
	}

	numUnaryChannels := DEFAULT_NUM_UNARY_GRPC_CHANNELS
	if request.TopicsConfiguration.GetNumUnaryGrpcChannels() > 0 {
		numUnaryChannels = request.TopicsConfiguration.GetNumUnaryGrpcChannels()
	}

	var timeout time.Duration
	if request.TopicsConfiguration.GetClientSideTimeout() < 1 {
		timeout = defaultRequestTimeout
	} else {
		timeout = request.TopicsConfiguration.GetClientSideTimeout()
	}

	topicManagerProps := &models.TopicStreamGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
		GrpcConfiguration:  grpcConfig,
	}
	unaryManagerList, err := grpcmanagers.NewStaticUnaryManagerList(topicManagerProps, numUnaryChannels, request.Log)
	if err != nil {
		return nil, err
	}
	streamManagerList, err := grpcmanagers.NewStaticStreamManagerList(topicManagerProps, numStreamChannels, request.Log)
	if err != nil {
		return nil, err
	}

	return &staticPubSubClient{
		endpoint:          request.CredentialProvider.GetCacheEndpoint(),
		log:               request.Log,
		requestTimeout:    timeout,
		middleware:        request.TopicsConfiguration.GetMiddleware(),
		unaryManagerList:  unaryManagerList,
		streamManagerList: streamManagerList,
	}, nil
}

func (client *staticPubSubClient) topicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (*grpcmanagers.TopicGrpcManager, pb.Pubsub_SubscribeClient, context.Context, context.CancelFunc, error) {
	topicManager, getManagerErr := client.streamManagerList.GetNextManager()
	if getManagerErr != nil {
		return nil, nil, nil, nil, getManagerErr
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

	client.log.Debug("Starting new subscription")
	return topicManager, subscribeClient, cancelContext, cancelFunction, err
}

func (client *staticPubSubClient) topicPublish(ctx context.Context, request *TopicPublishRequest) error {
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
	topicManager := client.unaryManagerList.GetNextManager()
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

func (client *staticPubSubClient) close() {
	client.streamManagerList.Close()
	client.unaryManagerList.Close()
}

func (client *staticPubSubClient) getMiddlewares() []middleware.TopicMiddleware {
	return client.middleware
}

func (client *staticPubSubClient) countNumberOfActiveSubscriptions() int64 {
	return client.streamManagerList.CountNumberOfActiveSubscriptions()
}

// Pubsub client implementation that uses a dynamic number of GRPC channels for subscriptions and unary operations.
type dynamicPubSubClient struct {
	endpoint          string
	log               logger.MomentoLogger
	requestTimeout    time.Duration
	middleware        []middleware.TopicMiddleware
	unaryManagerList  *grpcmanagers.StaticUnaryManagerList
	streamManagerList *grpcmanagers.DynamicStreamManagerList
}

func newDynamicPubSubClient(request *models.PubSubClientRequest) (*dynamicPubSubClient, momentoerrors.MomentoSvcErr) {
	grpcConfig := request.TopicsConfiguration.GetGrpcConfig()

	maxSubscriptions := DEFAULT_NUM_STREAM_GRPC_CHANNELS * config.MAX_CONCURRENT_STREAMS_PER_CHANNEL
	if request.TopicsConfiguration.GetMaxSubscriptions() > 0 {
		maxSubscriptions = request.TopicsConfiguration.GetMaxSubscriptions()
	}

	numUnaryChannels := DEFAULT_NUM_UNARY_GRPC_CHANNELS
	if request.TopicsConfiguration.GetNumUnaryGrpcChannels() > 0 {
		numUnaryChannels = request.TopicsConfiguration.GetNumUnaryGrpcChannels()
	}

	var timeout time.Duration
	if request.TopicsConfiguration.GetClientSideTimeout() < 1 {
		timeout = defaultRequestTimeout
	} else {
		timeout = request.TopicsConfiguration.GetClientSideTimeout()
	}

	topicManagerProps := &models.TopicStreamGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
		GrpcConfiguration:  grpcConfig,
	}
	unaryManagerList, err := grpcmanagers.NewStaticUnaryManagerList(topicManagerProps, numUnaryChannels, request.Log)
	if err != nil {
		return nil, err
	}
	streamManagerList, err := grpcmanagers.NewDynamicStreamManagerList(topicManagerProps, maxSubscriptions, request.Log)
	if err != nil {
		return nil, err
	}

	return &dynamicPubSubClient{
		endpoint:          request.CredentialProvider.GetCacheEndpoint(),
		log:               request.Log,
		requestTimeout:    timeout,
		middleware:        request.TopicsConfiguration.GetMiddleware(),
		unaryManagerList:  unaryManagerList,
		streamManagerList: streamManagerList,
	}, nil
}

func (client *dynamicPubSubClient) topicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (*grpcmanagers.TopicGrpcManager, pb.Pubsub_SubscribeClient, context.Context, context.CancelFunc, error) {
	topicManager, getStreamErr := client.streamManagerList.GetNextManager()
	if getStreamErr != nil {
		return nil, nil, nil, nil, getStreamErr
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

func (client *dynamicPubSubClient) topicPublish(ctx context.Context, request *TopicPublishRequest) error {
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
	topicManager := client.unaryManagerList.GetNextManager()
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

func (client *dynamicPubSubClient) close() {
	client.streamManagerList.Close()
	client.unaryManagerList.Close()
}

func (client *dynamicPubSubClient) getMiddlewares() []middleware.TopicMiddleware {
	return client.middleware
}

func (client *dynamicPubSubClient) countNumberOfActiveSubscriptions() int64 {
	return client.streamManagerList.CountNumberOfActiveSubscriptions()
}
