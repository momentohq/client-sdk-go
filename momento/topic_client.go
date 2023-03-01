// Package momento represents API CacheClient interface accessors including control/data operations, errors, operation requests and responses for the SDK.
package momento

import (
	"context"
	"fmt"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type TopicClient interface {
	Subscribe(ctx context.Context, request *TopicSubscribeRequest) (TopicSubscription, error)
	Publish(ctx context.Context, request *TopicPublishRequest) (TopicPublishResponse, error)

	Close()
}

// defaultScsClient represents all information needed for momento client to enable cache control and data operations.
type defaultTopicClient struct {
	credentialProvider auth.CredentialProvider
	pubSubClient       *pubSubClient
}

type TopicClientProps struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
}

// NewTopicClient returns a new TopicClient with provided authToken, DefaultTTLSeconds, and opts arguments.
func NewTopicClient(props *TopicClientProps) (TopicClient, error) {
	if props.Configuration.GetClientSideTimeout() < 1 {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "request timeout must not be 0", nil)
	}
	client := &defaultTopicClient{
		credentialProvider: props.CredentialProvider,
	}

	pubSubClient, err := newPubSubClient(&models.PubSubClientRequest{
		CredentialProvider: props.CredentialProvider,
		Configuration:      props.Configuration,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	client.pubSubClient = pubSubClient

	return client, nil
}

func (c defaultTopicClient) Subscribe(ctx context.Context, request *TopicSubscribeRequest) (TopicSubscription, error) {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}

	if _, err := prepareName(request.TopicName, "Topic name"); err != nil {
		return nil, err
	}

	clientStream, err := c.pubSubClient.TopicSubscribe(ctx, &TopicSubscribeRequest{
		CacheName: request.CacheName,
		TopicName: request.TopicName,
	})
	if err != nil {
		return nil, err
	}

	// Ping the stream to provide a nice error message if the cache does not exist.
	rawMsg := new(pb.XSubscriptionItem)
	err = clientStream.RecvMsg(rawMsg)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	switch rawMsg.Kind.(type) {
	case *pb.XSubscriptionItem_Heartbeat:
		// The first message to a new subscription will always be a heartbeat.
	default:
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.InternalServerError,
			fmt.Sprintf("expected a heartbeat message, got: %T", rawMsg.Kind),
			err,
		)
	}

	return &topicSubscription{
		grpcClient:         clientStream,
		momentoTopicClient: c.pubSubClient,
		cacheName:          request.CacheName,
		topicName:          request.TopicName,
	}, nil
}

func (c defaultTopicClient) Publish(ctx context.Context, request *TopicPublishRequest) (TopicPublishResponse, error) {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}

	if _, err := prepareName(request.TopicName, "Topic name"); err != nil {
		return nil, err
	}

	if request.Value == nil {
		return nil, convertMomentoSvcErrorToCustomerError(
			momentoerrors.NewMomentoSvcErr(
				momentoerrors.InvalidArgumentError, "value cannot be nil", nil,
			),
		)
	}

	err := c.pubSubClient.TopicPublish(ctx, &TopicPublishRequest{
		CacheName: request.CacheName,
		TopicName: request.TopicName,
		Value:     request.Value,
	})

	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	return &TopicPublishSuccess{}, err
}
func (c defaultTopicClient) Close() {
	defer c.pubSubClient.Close()
}
