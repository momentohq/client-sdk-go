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
// type pubSubClient interface {
// 	topicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (*grpcmanagers.TopicGrpcManager, pb.Pubsub_SubscribeClient, context.Context, context.CancelFunc, error)
// 	topicPublish(ctx context.Context, request *TopicPublishRequest) error
// 	close()
// 	getMiddlewares() []middleware.TopicMiddleware
// 	countNumberOfActiveSubscriptions() int64
// }

type pubSubClient struct {
	endpoint          string
	log               logger.MomentoLogger
	requestTimeout    time.Duration
	middleware        []middleware.TopicMiddleware
	unaryManagerList  grpcmanagers.TopicManagerList
	streamManagerList grpcmanagers.TopicManagerList
}

func newPubSubClient(request *models.PubSubClientRequest) (*pubSubClient, momentoerrors.MomentoSvcErr) {
	var timeout time.Duration
	if request.TopicsConfiguration.GetClientSideTimeout() < 1 {
		timeout = defaultRequestTimeout
	} else {
		timeout = request.TopicsConfiguration.GetClientSideTimeout()
	}

	grpcConfig := request.TopicsConfiguration.GetGrpcConfig()
	topicManagerProps := &models.TopicStreamGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
		GrpcConfiguration:  grpcConfig,
	}

	// Create pool of grpc channels for unary operations
	numUnaryChannels := DEFAULT_NUM_UNARY_GRPC_CHANNELS
	if request.TopicsConfiguration.GetNumUnaryGrpcChannels() > 0 {
		numUnaryChannels = request.TopicsConfiguration.GetNumUnaryGrpcChannels()
	}
	unaryManagerList, err := grpcmanagers.NewStaticUnaryManagerList(topicManagerProps, numUnaryChannels, request.Log)
	if err != nil {
		return nil, err
	}

	// Create pool of grpc channels for stream operations depending on static vs dynamic transport strategy
	var streamManagerList grpcmanagers.TopicManagerList
	switch request.TopicsConfiguration.GetTransportStrategy().(type) {
	case *config.TopicsStaticTransportStrategy:
		numStreamChannels := DEFAULT_NUM_STREAM_GRPC_CHANNELS
		if request.TopicsConfiguration.GetNumStreamGrpcChannels() > 0 {
			numStreamChannels = request.TopicsConfiguration.GetNumStreamGrpcChannels()
		}
		streamManagerList, err = grpcmanagers.NewStaticStreamManagerList(topicManagerProps, numStreamChannels, request.Log)
		if err != nil {
			return nil, err
		}
	case *config.TopicsDynamicTransportStrategy:
		maxSubscriptions := DEFAULT_NUM_STREAM_GRPC_CHANNELS * config.MAX_CONCURRENT_STREAMS_PER_CHANNEL
		if request.TopicsConfiguration.GetMaxSubscriptions() > 0 {
			maxSubscriptions = request.TopicsConfiguration.GetMaxSubscriptions()
		}
		streamManagerList, err = grpcmanagers.NewDynamicStreamManagerList(topicManagerProps, maxSubscriptions, request.Log)
		if err != nil {
			return nil, err
		}
	}

	return &pubSubClient{
		endpoint:          request.CredentialProvider.GetCacheEndpoint(),
		log:               request.Log,
		requestTimeout:    timeout,
		middleware:        request.TopicsConfiguration.GetMiddleware(),
		unaryManagerList:  unaryManagerList,
		streamManagerList: streamManagerList,
	}, nil
}

func (client *pubSubClient) topicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (*grpcmanagers.TopicGrpcManager, pb.Pubsub_SubscribeClient, context.Context, context.CancelFunc, error) {
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
	topicManager, err := client.unaryManagerList.GetNextManager()
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
	client.streamManagerList.Close()
	client.unaryManagerList.Close()
}

func (client *pubSubClient) countNumberOfActiveSubscriptions() int64 {
	switch streamManagerList := client.streamManagerList.(type) {
	case *grpcmanagers.StaticStreamManagerList:
		return streamManagerList.CountNumberOfActiveSubscriptions()
	case *grpcmanagers.DynamicStreamManagerList:
		return streamManagerList.CountNumberOfActiveSubscriptions()
	}
	return 0
}
