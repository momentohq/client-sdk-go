// Package momento represents API ScsClient interface accessors including control/data operations, errors, operation requests and responses for the SDK.
package momento

import (
	"context"
	"strings"
	"time"

	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/internal/services"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
)

// ScsClient represents all information needed for momento client to enable cache control and data operations.
type ScsClient struct {
	credentialProvider auth.CredentialProvider
	controlClient      *services.ScsControlClient
	dataClient         *scsDataClient
	pubSubClient       *pubSubClient
}

type SimpleCacheClientProps struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
	DefaultTTL         time.Duration
}

// NewSimpleCacheClient returns a new ScsClient with provided authToken, DefaultTTLSeconds, and opts arguments.
func NewSimpleCacheClient(props *SimpleCacheClientProps) (*ScsClient, error) {
	if props.Configuration.GetClientSideTimeout() < 1 {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "request timeout must not be 0", nil)
	}
	client := &ScsClient{
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

	if props.DefaultTTL == 0 {
		return nil, convertMomentoSvcErrorToCustomerError(
			momentoerrors.NewMomentoSvcErr(
				momentoerrors.InvalidArgumentError,
				"Must Define a non zero Default TTL", nil),
		)
	}

	dataClient, err := newScsDataClient(&models.DataClientRequest{
		CredentialProvider: props.CredentialProvider,
		Configuration:      props.Configuration,
		DefaultTtl:         props.DefaultTTL,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	client.dataClient = dataClient
	client.controlClient = controlClient
	client.pubSubClient = pubSubClient

	return client, nil
}

func (c ScsClient) CreateCache(ctx context.Context, request *CreateCacheRequest) error {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return err
	}
	err := c.controlClient.CreateCache(ctx, &models.CreateCacheRequest{
		CacheName: request.CacheName,
	})
	if err != nil {
		return convertMomentoSvcErrorToCustomerError(err)
	}
	return nil
}

func (c ScsClient) DeleteCache(ctx context.Context, request *DeleteCacheRequest) error {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return err
	}
	err := c.controlClient.DeleteCache(ctx, &models.DeleteCacheRequest{
		CacheName: request.CacheName,
	})
	if err != nil {
		return convertMomentoSvcErrorToCustomerError(err)
	}
	return nil
}

func (c ScsClient) ListCaches(ctx context.Context, request *ListCachesRequest) (*ListCachesResponse, error) {
	rsp, err := c.controlClient.ListCaches(ctx, &models.ListCachesRequest{
		NextToken: request.NextToken,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return &ListCachesResponse{
		nextToken: rsp.NextToken,
		caches:    convertCacheInfo(rsp.Caches),
	}, nil
}

func (c ScsClient) Set(ctx context.Context, r *SetRequest) (SetResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c ScsClient) Get(ctx context.Context, r *GetRequest) (GetResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c ScsClient) Delete(ctx context.Context, r *DeleteRequest) (DeleteResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c ScsClient) TopicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (TopicSubscription, error) {
	clientStream, err := c.pubSubClient.TopicSubscribe(ctx, &TopicSubscribeRequest{
		CacheName: request.CacheName,
		TopicName: request.TopicName,
	})
	if err != nil {
		return nil, err
	}
	return topicSubscription{grpcClient: clientStream}, err
}

func (c ScsClient) TopicPublish(ctx context.Context, request *TopicPublishRequest) (TopicPublishResponse, error) {
	err := c.pubSubClient.TopicPublish(ctx, &TopicPublishRequest{
		CacheName: request.CacheName,
		TopicName: request.TopicName,
		Value:     request.Value,
	})

	if err != nil {
		return nil, err
	}

	return TopicPublishSuccess{}, err
}

func (c ScsClient) SortedSetFetch(ctx context.Context, r *SortedSetFetchRequest) (SortedSetFetchResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c ScsClient) SortedSetPut(ctx context.Context, r *SortedSetPutRequest) (SortedSetPutResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c ScsClient) SortedSetGetScore(ctx context.Context, r *SortedSetGetScoreRequest) (SortedSetGetScoreResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c ScsClient) SortedSetRemove(ctx context.Context, r *SortedSetRemoveRequest) (SortedSetRemoveResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c ScsClient) SortedSetGetRank(ctx context.Context, r *SortedSetGetRankRequest) (SortedSetGetRankResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c ScsClient) SortedSetIncrement(ctx context.Context, r *SortedSetIncrementRequest) (SortedSetIncrementResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c *ScsClient) Close() {
	defer c.controlClient.Close()
	defer c.dataClient.Close()
}

func convertMomentoSvcErrorToCustomerError(e momentoerrors.MomentoSvcErr) MomentoError {
	if e == nil {
		return nil
	}
	return NewMomentoError(e.Code(), e.Message(), e.OriginalErr())
}

func convertCacheInfo(i []models.CacheInfo) []CacheInfo {
	var convertedList []CacheInfo
	for _, c := range i {
		convertedList = append(convertedList, CacheInfo{
			name: c.Name,
		})
	}
	return convertedList
}

func isCacheNameValid(cacheName string) momentoerrors.MomentoSvcErr {
	if len(strings.TrimSpace(cacheName)) < 1 {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Cache name cannot be empty", nil)
	}
	return nil
}
