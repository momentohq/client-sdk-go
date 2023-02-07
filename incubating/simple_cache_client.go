// Package incubating represents experimental packages and clients for Momento
package incubating

import (
	"context"
	"fmt"
	"strings"

	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/internal/services"
	"github.com/momentohq/client-sdk-go/momento"
)

type ScsClient interface {
	momento.ScsClient

	SubscribeTopic(ctx context.Context, request *TopicSubscribeRequest) (SubscriptionIFace, error)
	PublishTopic(ctx context.Context, request *TopicPublishRequest) error

	ListFetch(ctx context.Context, request *ListFetchRequest) (ListFetchResponse, error)
	ListLength(ctx context.Context, request *ListLengthRequest) (ListLengthResponse, error)

	Close()
}

// DefaultScsClient default implementation of the Momento incubating ScsClient interface
type DefaultScsClient struct {
	controlClient  *services.ScsControlClient
	dataClient     *services.ScsDataClient
	pubSubClient   *services.PubSubClient
	internalClient momento.ScsClient
}

// NewScsClient returns a new ScsClient with provided authToken, defaultTtl,, and opts arguments.
func NewScsClient(props *momento.SimpleCacheClientProps) (ScsClient, error) {

	controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
		CredentialProvider: props.CredentialProvider,
		Configuration:      props.Configuration,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	dataClient, err := services.NewScsDataClient(&models.DataClientRequest{
		CredentialProvider: props.CredentialProvider,
		Configuration:      props.Configuration,
		DefaultTtlSeconds:  props.DefaultTTLSeconds,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	pubSubClient, err := services.NewPubSubClient(&models.PubSubClientRequest{
		CredentialProvider: props.CredentialProvider,
		Configuration:      props.Configuration,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}
	internalClient, mErr := momento.NewSimpleCacheClient(props)
	if mErr != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}
	client := &DefaultScsClient{
		controlClient:  controlClient,
		dataClient:     dataClient,
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
	switch value := request.Value.(type) {
	case *TopicValueBytes:
		return c.pubSubClient.Publish(ctx, &models.TopicPublishRequest{
			CacheName: request.CacheName,
			TopicName: request.TopicName,
			Value: &models.TopicValueBytes{
				Bytes: value.Bytes,
			},
		})
	case *TopicValueString:
		return c.pubSubClient.Publish(ctx, &models.TopicPublishRequest{
			CacheName: request.CacheName,
			TopicName: request.TopicName,
			Value: &models.TopicValueString{
				Text: value.Text,
			},
		})
	default:
		return momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			fmt.Sprintf("unexpected TopicPublishRequest type passed %+v", value),
			nil,
		)
	}
}

func (c *DefaultScsClient) ListFetch(ctx context.Context, request *ListFetchRequest) (ListFetchResponse, error) {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}
	rsp, err := c.dataClient.ListFetch(ctx, &models.ListFetchRequest{
		CacheName: request.CacheName,
		ListName:  request.ListName,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	// TODO: Um, yeah
	rsp2, _ := convertListFetchResponse(rsp)
	return rsp2, nil
}

func (c *DefaultScsClient) ListLength(ctx context.Context, request *ListLengthRequest) (ListLengthResponse, error) {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}
	rsp, err := c.dataClient.ListLength(ctx, &models.ListLengthRequest{
		CacheName: request.CacheName,
		ListName:  request.ListName,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	rsp2, _ := convertListLengthResponse(rsp)
	return rsp2, nil
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
func (c *DefaultScsClient) Set(ctx context.Context, request *momento.CacheSetRequest) error {
	return c.internalClient.Set(ctx, request)
}
func (c *DefaultScsClient) Get(ctx context.Context, request *momento.CacheGetRequest) (momento.CacheGetResponse, error) {
	return c.internalClient.Get(ctx, request)
}
func (c *DefaultScsClient) Delete(ctx context.Context, request *momento.CacheDeleteRequest) error {
	return c.internalClient.Delete(ctx, request)
}

func convertListFetchResponse(r models.ListFetchResponse) (ListFetchResponse, momento.MomentoError) {
	switch response := r.(type) {
	case *models.ListFetchMiss:
		return &ListFetchMiss{}, nil
	case *models.ListFetchHit:
		return &ListFetchHit{
			value: response.Value,
		}, nil
	default:
		return nil, momentoerrors.NewMomentoSvcErr(
			momento.ClientSdkError,
			fmt.Sprintf("unexpected list fetch status returned %+v", response),
			nil,
		)
	}
}

func convertListLengthResponse(r models.ListLengthResponse) (ListLengthResponse, momento.MomentoError) {
	switch response := r.(type) {
	case *models.ListLengthSuccess:
		return &ListLengthSuccess{value: response.Value}, nil
	default:
		return nil, momentoerrors.NewMomentoSvcErr(
			momento.ClientSdkError,
			fmt.Sprintf("unexpected list fetch status returned %+v", response),
			nil,
		)
	}
}

// TODO: refactor these for sharing with momento module
func isCacheNameValid(cacheName string) momentoerrors.MomentoSvcErr {
	if len(strings.TrimSpace(cacheName)) < 1 {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Cache name cannot be empty", nil)
	}
	return nil
}
