package services

import (
	"context"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc"
)

type PubSubClient struct {
	grpcManager *grpcmanagers.ScsDataGrpcManager
	grpcClient  pb.PubsubClient
	endpoint    string
}

func NewPubSubClient(request *models.PubSubClientRequest) (*PubSubClient, momentoerrors.MomentoSvcErr) {
	//dataManager, err := grpcmanagers.NewScsDataGrpcManager(&models.DataGrpcManagerRequest{ // FIXME
	dataManager, err := grpcmanagers.NewLocalScsDataGrpcManager(&models.DataGrpcManagerRequest{
		AuthToken: request.AuthToken,
		//Endpoint:  fmt.Sprint(request.Endpoint, cachePort), // FIXME
		Endpoint: request.Endpoint,
	})

	if err != nil {
		return nil, err
	}
	return &PubSubClient{
		grpcManager: dataManager,
		grpcClient:  pb.NewPubsubClient(dataManager.Conn),
		endpoint:    request.Endpoint,
	}, nil
}

type PubSubSubscriptionWrapper struct {
	grpcClient grpc.ClientStream
}

func (client *PubSubClient) Subscribe(ctx context.Context, request *models.TopicSubscribeRequest) (grpc.ClientStream, error) {
	streamClient, err := client.grpcClient.Subscribe(context.Background(), &pb.XSubscriptionRequest{
		CacheName: "topic-" + request.TopicName,
		Topic:     request.TopicName,
		//ResumeAtTopicSequenceNumber: 0, TODO think about re-establish case may want to expose this
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
