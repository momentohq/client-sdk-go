// Package incubating represents experimental packages and clients for Momento

package incubating

import (
	"context"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/internal/services"
	"github.com/momentohq/client-sdk-go/momento"
)

type ScsClient interface {
	momento.ScsClient

	SubscribeTopic(ctx context.Context, request *TopicSubscribeRequest) (SubscriptionIFace, error)
	PublishTopic(ctx context.Context, request *TopicPublishRequest) error

	Close()
}

// DefaultScsClient default implementation of the Momento incubating ScsClient interface
type DefaultScsClient struct {
	controlClient  *services.ScsControlClient
	pubSubClient   *services.PubSubClient
	internalClient momento.ScsClient
}

// NewScsClient returns a new ScsClient with provided authToken, defaultTtl,, and opts arguments.
func NewScsClient(credentialProvider auth.CredentialProvider, defaultTtlSeconds uint32, opts ...momento.Option) (ScsClient, error) {

	controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
		CredentialProvider: credentialProvider,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	pubSubClient, err := services.NewPubSubClient(&models.PubSubClientRequest{
		CredentialProvider: credentialProvider,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}
	internalClient, mErr := momento.NewSimpleCacheClient(credentialProvider, defaultTtlSeconds, opts...)
	if mErr != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}
	client := &DefaultScsClient{
		controlClient:  controlClient,
		pubSubClient:   pubSubClient,
		internalClient: internalClient,
	}

	return client, nil
}

func newLocalScsClient(port int) (ScsClient, error) {
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

	client := &DefaultScsClient{
		//controlClient: controlClient,
		pubSubClient: pubSubClient,
	}
	return client, nil
}
func (c *DefaultScsClient) SubscribeTopic(ctx context.Context, request *TopicSubscribeRequest) (SubscriptionIFace, error) {
	clientStream, err := c.pubSubClient.Subscribe(ctx, &models.TopicSubscribeRequest{
		CacheName: request.CacheName,
		TopicName: request.TopicName,
	})
	if err != nil {
		return nil, err
	}
	return &Subscription{grpcClient: clientStream}, err
}

func (c *DefaultScsClient) PublishTopic(ctx context.Context, request *TopicPublishRequest) error {
	return c.pubSubClient.Publish(ctx, &models.TopicPublishRequest{
		CacheName: request.CacheName,
		TopicName: request.TopicName,
		Value:     request.Value,
	})
}

// Close shutdown the client.
func (c *DefaultScsClient) Close() {
	defer c.internalClient.Close()
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

func (c *DefaultScsClient) CreateCache(ctx context.Context, request *momento.CreateCacheRequest) error {
	return c.internalClient.CreateCache(ctx, request)
}

func (c *DefaultScsClient) DeleteCache(ctx context.Context, request *momento.DeleteCacheRequest) error {
	return c.internalClient.DeleteCache(ctx, request)
}

func (c *DefaultScsClient) ListCaches(ctx context.Context, request *momento.ListCachesRequest) (*momento.ListCachesResponse, error) {
	return c.internalClient.ListCaches(ctx, request)
}

func (c *DefaultScsClient) CreateSigningKey(ctx context.Context, request *momento.CreateSigningKeyRequest) (*momento.CreateSigningKeyResponse, error) {
	return c.internalClient.CreateSigningKey(ctx, request)
}

func (c *DefaultScsClient) RevokeSigningKey(ctx context.Context, request *momento.RevokeSigningKeyRequest) error {
	return c.internalClient.RevokeSigningKey(ctx, request)
}

func (c *DefaultScsClient) ListSigningKeys(ctx context.Context, request *momento.ListSigningKeysRequest) (*momento.ListSigningKeysResponse, error) {
	return c.internalClient.ListSigningKeys(ctx, request)
}

func (c *DefaultScsClient) Set(ctx context.Context, request *momento.CacheSetRequest) (*momento.SetCacheResponse, error) {
	return c.internalClient.Set(ctx, request)
}

func (c *DefaultScsClient) Get(ctx context.Context, request *momento.CacheGetRequest) (*momento.GetCacheResponse, error) {
	return c.internalClient.Get(ctx, request)
}

func (c *DefaultScsClient) Delete(ctx context.Context, request *momento.CacheDeleteRequest) error {
	return c.internalClient.Delete(ctx, request)
}
