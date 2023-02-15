// Package momento represents API SimpleCacheClient interface accessors including control/data operations, errors, operation requests and responses for the SDK.
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

type SimpleCacheClient interface {
	CreateCache(ctx context.Context, request *CreateCacheRequest) error
	DeleteCache(ctx context.Context, request *DeleteCacheRequest) error
	ListCaches(ctx context.Context, request *ListCachesRequest) (*ListCachesResponse, error)

	Set(ctx context.Context, r *SetRequest) (SetResponse, error)
	Get(ctx context.Context, r *GetRequest) (GetResponse, error)
	Delete(ctx context.Context, r *DeleteRequest) (DeleteResponse, error)

	TopicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (TopicSubscription, error)
	TopicPublish(ctx context.Context, request *TopicPublishRequest) (TopicPublishResponse, error)

	SortedSetFetch(ctx context.Context, r *SortedSetFetchRequest) (SortedSetFetchResponse, error)
	SortedSetPut(ctx context.Context, r *SortedSetPutRequest) (SortedSetPutResponse, error)
	SortedSetGetScore(ctx context.Context, r *SortedSetGetScoreRequest) (SortedSetGetScoreResponse, error)
	SortedSetRemove(ctx context.Context, r *SortedSetRemoveRequest) (SortedSetRemoveResponse, error)
	SortedSetGetRank(ctx context.Context, r *SortedSetGetRankRequest) (SortedSetGetRankResponse, error)
	SortedSetIncrement(ctx context.Context, r *SortedSetIncrementRequest) (SortedSetIncrementResponse, error)

	ListPushFront(ctx context.Context, r *ListPushFrontRequest) (ListPushFrontResponse, error)
	ListPushBack(ctx context.Context, r *ListPushBackRequest) (ListPushBackResponse, error)
	ListPopFront(ctx context.Context, r *ListPopFrontRequest) (ListPopFrontResponse, error)
	ListPopBack(ctx context.Context, r *ListPopBackRequest) (ListPopBackResponse, error)
	ListConcatenateFront(ctx context.Context, r *ListConcatenateFrontRequest) (ListConcatenateFrontResponse, error)
	ListConcatenateBack(ctx context.Context, r *ListConcatenateBackRequest) (ListConcatenateBackResponse, error)
	ListFetch(ctx context.Context, r *ListFetchRequest) (ListFetchResponse, error)
	ListLength(ctx context.Context, r *ListLengthRequest) (ListLengthResponse, error)
	ListRemoveValue(ctx context.Context, r *ListRemoveValueRequest) (ListRemoveValueResponse, error)

	Close()
}

// defaultScsClient represents all information needed for momento client to enable cache control and data operations.
type defaultScsClient struct {
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

// NewSimpleCacheClient returns a new defaultScsClient with provided authToken, DefaultTTLSeconds, and opts arguments.
func NewSimpleCacheClient(props *SimpleCacheClientProps) (SimpleCacheClient, error) {
	if props.Configuration.GetClientSideTimeout() < 1 {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "request timeout must not be 0", nil)
	}
	client := &defaultScsClient{
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

func (c defaultScsClient) CreateCache(ctx context.Context, request *CreateCacheRequest) error {
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

func (c defaultScsClient) DeleteCache(ctx context.Context, request *DeleteCacheRequest) error {
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

func (c defaultScsClient) ListCaches(ctx context.Context, request *ListCachesRequest) (*ListCachesResponse, error) {
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

func (c defaultScsClient) Set(ctx context.Context, r *SetRequest) (SetResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) Get(ctx context.Context, r *GetRequest) (GetResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) Delete(ctx context.Context, r *DeleteRequest) (DeleteResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) TopicSubscribe(ctx context.Context, request *TopicSubscribeRequest) (TopicSubscription, error) {
	clientStream, err := c.pubSubClient.TopicSubscribe(ctx, &TopicSubscribeRequest{
		CacheName: request.CacheName,
		TopicName: request.TopicName,
	})
	if err != nil {
		return nil, err
	}
	return topicSubscription{grpcClient: clientStream}, err
}

func (c defaultScsClient) TopicPublish(ctx context.Context, request *TopicPublishRequest) (TopicPublishResponse, error) {
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

func (c defaultScsClient) SortedSetFetch(ctx context.Context, r *SortedSetFetchRequest) (SortedSetFetchResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetPut(ctx context.Context, r *SortedSetPutRequest) (SortedSetPutResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetGetScore(ctx context.Context, r *SortedSetGetScoreRequest) (SortedSetGetScoreResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetRemove(ctx context.Context, r *SortedSetRemoveRequest) (SortedSetRemoveResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetGetRank(ctx context.Context, r *SortedSetGetRankRequest) (SortedSetGetRankResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetIncrement(ctx context.Context, r *SortedSetIncrementRequest) (SortedSetIncrementResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListPushFront(ctx context.Context, r *ListPushFrontRequest) (ListPushFrontResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListPushBack(ctx context.Context, r *ListPushBackRequest) (ListPushBackResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListPopFront(ctx context.Context, r *ListPopFrontRequest) (ListPopFrontResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListPopBack(ctx context.Context, r *ListPopBackRequest) (ListPopBackResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListConcatenateFront(ctx context.Context, r *ListConcatenateFrontRequest) (ListConcatenateFrontResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListConcatenateBack(ctx context.Context, r *ListConcatenateBackRequest) (ListConcatenateBackResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListFetch(ctx context.Context, r *ListFetchRequest) (ListFetchResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListLength(ctx context.Context, r *ListLengthRequest) (ListLengthResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListRemoveValue(ctx context.Context, r *ListRemoveValueRequest) (ListRemoveValueResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) Close() {
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
