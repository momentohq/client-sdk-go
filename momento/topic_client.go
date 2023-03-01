// Package momento represents API CacheClient interface accessors including control/data operations, errors, operation requests and responses for the SDK.
package momento

import (
	"context"
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/internal/services"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
)

type TopicClient interface {
	CreateCache(ctx context.Context, request *CreateCacheRequest) (CreateCacheResponse, error)
	DeleteCache(ctx context.Context, request *DeleteCacheRequest) (DeleteCacheResponse, error)
	ListCaches(ctx context.Context, request *ListCachesRequest) (ListCachesResponse, error)

	TopicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (TopicSubscription, error)
	TopicPublish(ctx context.Context, request *TopicPublishRequest) (TopicPublishResponse, error)

	Close()
}

// defaultScsClient represents all information needed for momento client to enable cache control and data operations.
type defaultTopicClient struct {
	credentialProvider auth.CredentialProvider
	controlClient      *services.ScsControlClient
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

	controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
		CredentialProvider: props.CredentialProvider,
		Configuration:      props.Configuration,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	pubSubClient, err := newPubSubClient(&models.PubSubClientRequest{
		CredentialProvider: props.CredentialProvider,
		Configuration:      props.Configuration,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	client.controlClient = controlClient
	client.pubSubClient = pubSubClient

	return client, nil
}

func (c defaultTopicClient) CreateCache(ctx context.Context, request *CreateCacheRequest) (CreateCacheResponse, error) {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}
	err := c.controlClient.CreateCache(ctx, &models.CreateCacheRequest{
		CacheName: request.CacheName,
	})
	if err != nil {
		if err.Code() == AlreadyExistsError {
			return &CreateCacheAlreadyExists{}, nil
		}
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return &CreateCacheSuccess{}, nil
}

func (c defaultTopicClient) DeleteCache(ctx context.Context, request *DeleteCacheRequest) (DeleteCacheResponse, error) {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}
	err := c.controlClient.DeleteCache(ctx, &models.DeleteCacheRequest{
		CacheName: request.CacheName,
	})
	if err != nil {
		if err.Code() == NotFoundError {
			return &DeleteCacheSuccess{}, nil
		}
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return &DeleteCacheSuccess{}, nil
}

func (c defaultTopicClient) ListCaches(ctx context.Context, request *ListCachesRequest) (ListCachesResponse, error) {
	rsp, err := c.controlClient.ListCaches(ctx, &models.ListCachesRequest{
		NextToken: request.NextToken,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return &ListCachesSuccess{
		nextToken: rsp.NextToken,
		caches:    convertCacheInfo(rsp.Caches),
	}, nil
}

func (c defaultTopicClient) TopicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (TopicSubscription, error) {
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
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.NotFoundError,
			fmt.Sprintf("Did not get a heartbeat from topic %v in cache %v", request.TopicName, request.CacheName),
			err,
		)
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

func (c defaultTopicClient) TopicPublish(ctx context.Context, request *TopicPublishRequest) (TopicPublishResponse, error) {
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
	defer c.controlClient.Close()
	defer c.pubSubClient.Close()
}
