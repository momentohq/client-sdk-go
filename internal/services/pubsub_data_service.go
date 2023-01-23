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
	grpcManager *grpcmanagers.DataGrpcManager
	grpcClient  pb.PubsubClient
	endpoint    string
}

func NewPubSubClient(request *models.PubSubClientRequest) (*PubSubClient, momentoerrors.MomentoSvcErr) {
	dataManager, err := grpcmanagers.NewStreamDataGrpcManager(&models.DataGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
	})

	if err != nil {
		return nil, err
	}
	return &PubSubClient{
		grpcManager: dataManager,
		grpcClient:  pb.NewPubsubClient(dataManager.Conn),
		endpoint:    request.CredentialProvider.GetCacheEndpoint(),
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
		grpcManager: dataManager,
		grpcClient:  pb.NewPubsubClient(dataManager.Conn),
		endpoint:    "localhost",
	}, nil
}

func (client *PubSubClient) Subscribe(ctx context.Context, request *models.TopicSubscribeRequest) (grpc.ClientStream, error) {
	streamClient, err := client.grpcClient.Subscribe(ctx, &pb.XSubscriptionRequest{
		CacheName: "topic-" + request.TopicName,
		Topic:     request.TopicName,
		//ResumeAtTopicSequenceNumber: 0, TODO think about re-establish topic case
	})
	return streamClient, err
}
func (client *PubSubClient) Publish(ctx context.Context, request *models.TopicPublishRequest) error {
	_, err := client.grpcClient.Publish(ctx, &pb.XPublishRequest{
		CacheName: "topic-" + request.TopicName,
		Topic:     request.TopicName,
		Value: &pb.XTopicValue{
			Kind: &pb.XTopicValue_Text{
				Text: request.Value,
			},
		},
	})
	return err
}

func (client *PubSubClient) Endpoint() string {
	return client.endpoint
}

func (client *PubSubClient) Close() momentoerrors.MomentoSvcErr {
	return client.grpcManager.Close()
}
