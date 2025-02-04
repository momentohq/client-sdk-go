package momento

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type pubSubClient struct {
	numUnaryChannels        uint32
	unaryTopicManagers      []*grpcmanagers.TopicGrpcManager
	unaryTopicManagerCount  atomic.Uint64
	numStreamChannels       uint32
	streamTopicManagers     []*grpcmanagers.TopicGrpcManager
	streamTopicManagerCount atomic.Uint64
	endpoint                string
	log                     logger.MomentoLogger
	numGrpcStreams          atomic.Int64 // TODO: encapsulate in grpc manager
}

func newPubSubClient(request *models.PubSubClientRequest) (*pubSubClient, momentoerrors.MomentoSvcErr) {
	// NOTE: This is hard-coded for now but we may want to expose it via TopicConfiguration in the future,
	// as we do with some of the other clients. Defaults to keep-alive pings enabled.
	grpcConfig := config.NewStaticGrpcConfiguration(&config.GrpcConfigurationProps{})

	// Default to using 4 grpc channels for subscriptions
	numStreamChannels := uint32(4)
	if request.TopicsConfiguration.GetNumStreamGrpcChannels() > 0 {
		numStreamChannels = request.TopicsConfiguration.GetNumStreamGrpcChannels()
	} else if request.TopicsConfiguration.GetNumGrpcChannels() > 0 {
		// numGrpcChannels is deprecated, but we'll use it to set both numUnaryChannels and numStreamChannels
		// in case there are customers still using it.
		numStreamChannels = request.TopicsConfiguration.GetNumGrpcChannels()
	} else if request.TopicsConfiguration.GetMaxSubscriptions() > 0 {
		// maxSubscriptions is deprecated, but we'll use it to set numStreamChannels
		// in case there are customers still using it.
		numStreamChannels = uint32(math.Ceil(float64(request.TopicsConfiguration.GetMaxSubscriptions()) / 100.0))
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

	// Default to using 4 grpc channels for publishes
	numUnaryChannels := uint32(4)
	if request.TopicsConfiguration.GetNumUnaryGrpcChannels() > 0 {
		numUnaryChannels = request.TopicsConfiguration.GetNumUnaryGrpcChannels()
	} else if request.TopicsConfiguration.GetNumGrpcChannels() > 0 {
		// numGrpcChannels is deprecated, but we'll use it to set both numUnaryChannels and numStreamChannels
		// in case there are customers still using it.
		numUnaryChannels = request.TopicsConfiguration.GetNumGrpcChannels()
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

	return &pubSubClient{
		numUnaryChannels:    numUnaryChannels,
		unaryTopicManagers:  unaryTopicManagers,
		numStreamChannels:   numStreamChannels,
		streamTopicManagers: streamTopicManagers,
		endpoint:            request.CredentialProvider.GetCacheEndpoint(),
		log:                 request.Log,
	}, nil
}

func (client *pubSubClient) getNextStreamTopicManager() *grpcmanagers.TopicGrpcManager {
	nextManagerIndex := client.streamTopicManagerCount.Add(1)
	topicManager := client.streamTopicManagers[nextManagerIndex%uint64(client.numStreamChannels)]
	return topicManager
}

func (client *pubSubClient) getNextUnaryTopicManager() *grpcmanagers.TopicGrpcManager {
	nextManagerIndex := client.unaryTopicManagerCount.Add(1)
	topicManager := client.unaryTopicManagers[nextManagerIndex%uint64(client.numUnaryChannels)]
	return topicManager
}

// Each grpc connection can multiplex 100 subscribe/publish requests.
//
// Publish requests will queue up on client while waiting for in-flight requests to complete if
// the number of concurrent requests exceeds numUnaryChannels*100, but will eventually complete.
//
// We must prevent subscription requests from similarly queuing up if there are already numStreamChannels*100
// concurrent streams as it causes the program to hang indefinitely with no error.
func (client *pubSubClient) checkNumConcurrentStreams(log logger.MomentoLogger) error {
	if client.numGrpcStreams.Load() >= int64(client.numStreamChannels*100) {
		errorMessage := fmt.Sprintf(
			"Number of grpc streams: %d; number of channels: %d; max concurrent streams: %d; Already at maximum number of concurrent grpc streams, cannot make new subscribe requests\n",
			client.numGrpcStreams.Load(), client.numStreamChannels, client.numStreamChannels*100,
		)
		return momentoerrors.NewMomentoSvcErr(LimitExceededError, errorMessage, nil)
	}
	return nil
}

func (client *pubSubClient) topicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (*grpcmanagers.TopicGrpcManager, grpc.ClientStream, context.Context, context.CancelFunc, error) {
	// First check if there is enough grpc stream capabity to make a new subscription
	err := client.checkNumConcurrentStreams(client.log)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// add metadata to context
	requestMetadata := internal.CreateMetadata(ctx, internal.Topic)

	// add withCancel to context
	cancelContext, cancelFunction := context.WithCancel(requestMetadata)

	var header, trailer metadata.MD
	client.numGrpcStreams.Add(1)
	topicManager := client.getNextStreamTopicManager()
	clientStream, err := topicManager.StreamClient.Subscribe(cancelContext, &pb.XSubscriptionRequest{
		CacheName:                   request.CacheName,
		Topic:                       request.TopicName,
		ResumeAtTopicSequenceNumber: request.ResumeAtTopicSequenceNumber,
		SequencePage:                request.SequencePage,
	})

	if err != nil {
		client.numGrpcStreams.Add(-1)
		cancelFunction()
		if clientStream != nil {
			header, _ = clientStream.Header()
			trailer = clientStream.Trailer()
		}
		return nil, nil, nil, nil, momentoerrors.ConvertSvcErr(err, header, trailer)
	}

	// If we are approaching the grpc maximum concurrent stream limit, log a warning
	numStreamsLimit := int64(client.numStreamChannels * 100)
	remainingStreams := numStreamsLimit - client.numGrpcStreams.Load()
	if remainingStreams < 10 {
		client.log.Warn(
			"WARNING: approaching grpc maximum concurrent stream limit, %d remaining of total %d streams\n",
			remainingStreams, numStreamsLimit,
		)
	}

	client.log.Debug("Started new subscription, number of concurrent streams: %d", client.numGrpcStreams.Load())
	return topicManager, clientStream, cancelContext, cancelFunction, err
}

func (client *pubSubClient) topicPublish(ctx context.Context, request *TopicPublishRequest) error {
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
	client.numStreamChannels = 0
	for clientIndex := range client.streamTopicManagers {
		defer client.streamTopicManagers[clientIndex].Close()
	}
}
