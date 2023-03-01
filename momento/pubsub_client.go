package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"google.golang.org/grpc"
)

type pubSubClient struct {
	streamDataManager *grpcmanagers.DataGrpcManager
	unaryDataManager  *grpcmanagers.DataGrpcManager
	streamGrpcClient  pb.PubsubClient
	unaryGrpcClient   pb.PubsubClient
	endpoint          string
}

func newPubSubClient(request *models.PubSubClientRequest) (*pubSubClient, momentoerrors.MomentoSvcErr) {
	streamDataManager, err := grpcmanagers.NewStreamDataGrpcManager(&models.DataGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
	})
	if err != nil {
		return nil, err
	}
	unaryDataManager, err := grpcmanagers.NewUnaryDataGrpcManager(&models.DataGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
	})
	if err != nil {
		return nil, err
	}
	return &pubSubClient{
		streamDataManager: streamDataManager,
		unaryDataManager:  unaryDataManager,
		streamGrpcClient:  pb.NewPubsubClient(streamDataManager.Conn),
		unaryGrpcClient:   pb.NewPubsubClient(unaryDataManager.Conn),
		endpoint:          request.CredentialProvider.GetCacheEndpoint(),
	}, nil
}

func (client *pubSubClient) TopicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (grpc.ClientStream, error) {
	streamClient, err := client.streamGrpcClient.Subscribe(ctx, &pb.XSubscriptionRequest{
		CacheName: request.CacheName,
		Topic:     request.TopicName,
		//ResumeAtTopicSequenceNumber: 0, TODO think about re-establish topic case
	})
	return streamClient, err
}
func (client *pubSubClient) TopicPublish(ctx context.Context, request *TopicPublishRequest) error {
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

func (client *pubSubClient) Endpoint() string {
	return client.endpoint
}

func (client *pubSubClient) Close() {
	defer client.streamDataManager.Close()
	defer client.unaryDataManager.Close()
}
