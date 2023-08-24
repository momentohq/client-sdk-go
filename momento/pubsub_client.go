package momento

import (
	"context"
	"math"
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/internal/retry"

	"google.golang.org/grpc"
)

type pubSubClient struct {
	streamTopicManagers []*grpcmanagers.TopicGrpcManager
	unaryDataManager    *grpcmanagers.DataGrpcManager
	unaryGrpcClient     pb.PubsubClient
	endpoint            string
}

var streamTopicManagerCount uint64

func newPubSubClient(request *models.PubSubClientRequest) (*pubSubClient, momentoerrors.MomentoSvcErr) {
	var numChannels uint32
	numSubscriptions := float64(request.TopicsConfiguration.GetMaxSubscriptions())
	if numSubscriptions > 0 {
		// a single channel can support 100 streams, so we need to create enough
		// channels to handle the maximum number of subscriptions
		numChannels = uint32(math.Ceil(numSubscriptions / 100.0))
	} else {
		numChannels = 1
	}
	streamTopicManagers := make([]*grpcmanagers.TopicGrpcManager, 0)

	for i := 0; uint32(i) < numChannels; i++ {
		streamTopicManager, err := grpcmanagers.NewStreamTopicGrpcManager(&models.TopicStreamGrpcManagerRequest{
			CredentialProvider: request.CredentialProvider,
		})
		if err != nil {
			return nil, err
		}
		streamTopicManagers = append(streamTopicManagers, streamTopicManager)
	}

	unaryDataManager, err := grpcmanagers.NewUnaryDataGrpcManager(&models.DataGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
		RetryStrategy:      retry.NewNeverRetryStrategy(),
	})
	if err != nil {
		return nil, err
	}

	return &pubSubClient{
		streamTopicManagers: streamTopicManagers,
		unaryDataManager:    unaryDataManager,
		unaryGrpcClient:     pb.NewPubsubClient(unaryDataManager.Conn),
		endpoint:            request.CredentialProvider.GetCacheEndpoint(),
	}, nil
}

func (client *pubSubClient) getNextStreamTopicManager() *grpcmanagers.TopicGrpcManager {
	nextMangerIndex := atomic.AddUint64(&streamTopicManagerCount, 1)
	topicManager := client.streamTopicManagers[nextMangerIndex%uint64(len(client.streamTopicManagers))]
	return topicManager
}

func (client *pubSubClient) topicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (*grpcmanagers.TopicGrpcManager, grpc.ClientStream, error) {
	topicManager := client.getNextStreamTopicManager()
	clientStream, err := topicManager.StreamClient.Subscribe(ctx, &pb.XSubscriptionRequest{
		CacheName:                   request.CacheName,
		Topic:                       request.TopicName,
		ResumeAtTopicSequenceNumber: request.ResumeAtTopicSequenceNumber,
	})
	return topicManager, clientStream, err
}

func (client *pubSubClient) topicPublish(ctx context.Context, request *TopicPublishRequest) error {
	switch value := request.Value.(type) {
	case String:
		_, err := client.unaryGrpcClient.Publish(ctx, &pb.XPublishRequest{
			CacheName: request.CacheName,
			Topic:     request.TopicName,
			Value: &pb.XTopicValue{
				Kind: &pb.XTopicValue_Text{
					Text: value.asString(),
				},
			},
		})
		return err
	case Bytes:
		_, err := client.unaryGrpcClient.Publish(ctx, &pb.XPublishRequest{
			CacheName: request.CacheName,
			Topic:     request.TopicName,
			Value: &pb.XTopicValue{
				Kind: &pb.XTopicValue_Binary{
					Binary: value.asBytes(),
				},
			},
		})
		return err
	default:
		return momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"error encoding topic value only support []byte or string currently", nil,
		)
	}
}

func (client *pubSubClient) close() {
	for clientIndex := range client.streamTopicManagers {
		defer client.streamTopicManagers[clientIndex].Close()
	}
	defer client.unaryDataManager.Close()
}
