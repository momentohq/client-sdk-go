// Package incubating represents experimental packages and clients for Momento

package incubating

import (
	"context"

	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/internal/resolver"
	"github.com/momentohq/client-sdk-go/internal/services"
	"github.com/momentohq/client-sdk-go/momento"
)

type PubSubClient interface {
	CreateTopic(ctx context.Context, request *CreateTopicRequest) error
	SubscribeTopic(ctx context.Context, request *TopicSubscribeRequest) (SubscriptionIFace, error)
	PublishTopic(ctx context.Context, request *TopicPublishRequest) error

	Close()
}

// DefaultPubSubClient represents all information needed for momento client to enable pubsub control and data operations.
type DefaultPubSubClient struct {
	controlClient         *services.ScsControlClient
	pubSubClient          *services.PubSubClient
	defaultRequestTimeout uint32
}

// NewPubSubClient returns a new PubSubClient with provided authToken, and opts arguments.
func NewPubSubClient(authToken string) (PubSubClient, error) {
	endpoints, err := resolver.Resolve(&models.ResolveRequest{
		AuthToken: authToken,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}

	controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
		AuthToken: authToken,
		Endpoint:  endpoints.ControlEndpoint,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	pubSubClient, err := services.NewPubSubClient(&models.PubSubClientRequest{
		AuthToken: authToken,
		Endpoint:  endpoints.CacheEndpoint,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}
	client := &DefaultPubSubClient{
		controlClient: controlClient,
		pubSubClient:  pubSubClient,
	}

	return client, nil
}

func NewLocalPubSubClient(port int) (PubSubClient, error) {
	// TODO impl basic local control plane for pubsub topics
	//controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
	//	AuthToken: authToken,
	//	Endpoint:  endpoints.ControlEndpoint,
	//})
	//if err != nil {
	//	return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	//}

	pubSubClient, err := services.NewLocalPubSubClient(port)
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	client := &DefaultPubSubClient{
		//controlClient: controlClient,
		pubSubClient: pubSubClient,
	}
	return client, nil
}

func (c *DefaultPubSubClient) CreateTopic(ctx context.Context, request *CreateTopicRequest) error {
	return c.controlClient.CreateTopic(ctx, &models.CreateTopicRequest{TopicName: request.TopicName})
}

func (c *DefaultPubSubClient) SubscribeTopic(ctx context.Context, request *TopicSubscribeRequest) (SubscriptionIFace, error) {
	clientStream, err := c.pubSubClient.Subscribe(ctx, &models.TopicSubscribeRequest{
		TopicName: request.TopicName,
	})
	if err != nil {
		return nil, err
	}
	return &Subscription{grpcClient: clientStream}, err
}

func (c *DefaultPubSubClient) PublishTopic(ctx context.Context, request *TopicPublishRequest) error {
	return c.pubSubClient.Publish(ctx, &models.TopicPublishRequest{
		TopicName: request.TopicName,
		Value:     request.Value,
	})
}

// Close shutdown the client.
func (c *DefaultPubSubClient) Close() {
	defer c.controlClient.Close()
	defer c.pubSubClient.Close()
}

// TODO figure out better way to dry this up is copy pasta from simple cache client
func convertMomentoSvcErrorToCustomerError(e momentoerrors.MomentoSvcErr) momento.MomentoError {
	if e == nil {
		return nil
	}
	return momento.NewMomentoError(e.Code(), e.Message(), e.OriginalErr())
}
