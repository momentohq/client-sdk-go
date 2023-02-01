package services

import (
	"context"
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"

	"google.golang.org/grpc"
)

type PubSubClient struct {
	streamDataManager *grpcmanagers.DataGrpcManager
	unaryDataManager  *grpcmanagers.DataGrpcManager
	streamGrpcClient  pb.PubsubClient
	unaryGrpcClient   pb.PubsubClient
	endpoint          string
}

func NewPubSubClient(request *models.PubSubClientRequest) (*PubSubClient, momentoerrors.MomentoSvcErr) {
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
	return &PubSubClient{
		streamDataManager: streamDataManager,
		unaryDataManager:  unaryDataManager,
		streamGrpcClient:  pb.NewPubsubClient(streamDataManager.Conn),
		unaryGrpcClient:   pb.NewPubsubClient(unaryDataManager.Conn),
		endpoint:          request.CredentialProvider.GetCacheEndpoint(),
	}, nil
}

func NewLocalPubSubClient(port int) (*PubSubClient, momentoerrors.MomentoSvcErr) {
	dataManager, err := grpcmanagers.NewLocalDataGrpcManager(&models.LocalDataGrpcManagerRequest{
		Endpoint: fmt.Sprintf("localhost:%d", port),
	})

	if err != nil {
		return nil, err
	}
	return &PubSubClient{
		unaryGrpcClient:  pb.NewPubsubClient(dataManager.Conn),
		streamGrpcClient: pb.NewPubsubClient(dataManager.Conn),
		endpoint:         "localhost",
	}, nil
}

func (client *PubSubClient) Subscribe(ctx context.Context, request *models.TopicSubscribeRequest) (grpc.ClientStream, error) {
	streamClient, err := client.streamGrpcClient.Subscribe(ctx, &pb.XSubscriptionRequest{
		CacheName: request.CacheName,
		Topic:     request.TopicName,
		//ResumeAtTopicSequenceNumber: 0, TODO think about re-establish topic case
	})
	return streamClient, err
}
func (client *PubSubClient) Publish(ctx context.Context, request *models.TopicPublishRequest) error {
	switch value := request.Value.(type) {
	case *models.TopicValueString:
		_, err := client.unaryGrpcClient.Publish(ctx, &pb.XPublishRequest{
			CacheName: request.CacheName,
			Topic:     request.TopicName,
			Value: &pb.XTopicValue{
				Kind: &pb.XTopicValue_Text{
					Text: value.Text,
				},
			},
		})
		return err
	case models.TopicValueBytes:
		_, err := client.unaryGrpcClient.Publish(ctx, &pb.XPublishRequest{
			CacheName: request.CacheName,
			Topic:     request.TopicName,
			Value: &pb.XTopicValue{
				Kind: &pb.XTopicValue_Binary{
					Binary: value.Bytes,
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

func (client *PubSubClient) Endpoint() string {
	return client.endpoint
}

func (client *PubSubClient) Close() {
	defer client.streamDataManager.Close()
	defer client.unaryDataManager.Close()
}
