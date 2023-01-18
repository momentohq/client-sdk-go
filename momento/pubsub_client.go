// Package momento represents API ScsClient interface accessors including control/data operations, errors, operation requests and responses for the SDK.

package momento

import (
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/internal/resolver"
	"github.com/momentohq/client-sdk-go/internal/services"
)

type PubSubClient interface {
	CreateTopic(request *CreateTopicRequest) error
	SubscribeTopic(request *TopicSubscribeRequest) (SubscriptionIFace, error)

	Close()
}

// DefaultPubSubClient represents all information needed for momento client to enable pubsub control and data operations.
type DefaultPubSubClient struct {
	authToken             string
	controlClient         *services.ScsControlClient
	pubSubClient          *services.PubSubClient
	defaultRequestTimeout uint32
}

// NewPubSubClient returns a new PubSubClient with provided authToken, and opts arguments.
func NewPubSubClient(authToken string, opts ...Option) (PubSubClient, error) {
	endpoints, err := resolver.Resolve(&models.ResolveRequest{
		AuthToken:        authToken,
		EndpointOverride: "localhost", // FIXME remove this just testing quick

	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}

	client := &DefaultPubSubClient{
		authToken: authToken,
	}

	// Loop through all user passed options before building up internal clients
	// No options for now FIXME refactor how we do SDK options so not tied to just SCSClient
	//for _, opt := range opts {
	//	// Call the option giving the instantiated
	//	// *House as the argument
	//	err := opt(client)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
		AuthToken: authToken,
		Endpoint:  endpoints.ControlEndpoint,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	pubSubClient, err := services.NewPubSubClient(&models.PubSubClientRequest{
		AuthToken: authToken,
		//Endpoint:  endpoints.CacheEndpoint,
		Endpoint: "localhost:3000", // FIXME dont hard code here
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	client.pubSubClient = pubSubClient
	client.controlClient = controlClient

	return client, nil
}

func (c *DefaultPubSubClient) CreateTopic(request *CreateTopicRequest) error {
	return c.controlClient.CreateTopic(&models.CreateTopicRequest{TopicName: request.TopicName})
}

func (c *DefaultPubSubClient) SubscribeTopic(request *TopicSubscribeRequest) (SubscriptionIFace, error) {
	clientStream, err := c.pubSubClient.Subscribe(&models.TopicSubscribeRequest{
		TopicName: request.TopicName,
	})
	if err != nil {
		return nil, err
	}
	return &Subscription{grpcClient: clientStream}, err
}

// Close shutdown the client.
func (c *DefaultPubSubClient) Close() {
	defer c.controlClient.Close()
	defer c.pubSubClient.Close()
}
