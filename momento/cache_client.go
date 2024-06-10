// Package momento represents API CacheClient interface accessors including control/data operations, errors, operation requests and responses for the SDK.
package momento

import (
	"context"
	"strings"
	"sync/atomic"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"

	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/internal/services"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/responses"
)

var dataClientCount uint64

type CacheClient interface {
	Logger() logger.MomentoLogger

	// CreateCache Creates a cache if it does not exist.
	CreateCache(ctx context.Context, request *CreateCacheRequest) (responses.CreateCacheResponse, error)
	// DeleteCache deletes a cache and all the items within it.
	DeleteCache(ctx context.Context, request *DeleteCacheRequest) (responses.DeleteCacheResponse, error)
	// ListCaches lists all caches.
	ListCaches(ctx context.Context, request *ListCachesRequest) (responses.ListCachesResponse, error)

	// Increment adds an integer quantity to a field value.
	Increment(ctx context.Context, r *IncrementRequest) (responses.IncrementResponse, error)
	// Set sets the value in cache with a given time to live (TTL)
	Set(ctx context.Context, r *SetRequest) (responses.SetResponse, error)
	// SetIfNotExists sets the value in cache with a given time to live (TTL) if key is not already present
	//
	// Deprecated: Use SetIfAbsent instead.
	SetIfNotExists(ctx context.Context, r *SetIfNotExistsRequest) (responses.SetIfNotExistsResponse, error)
	// SetIfAbsent sets the value in cache with a given time to live (TTL) if key is not already present
	SetIfAbsent(ctx context.Context, r *SetIfAbsentRequest) (responses.SetIfAbsentResponse, error)
	// SetIfPresent sets the value in cache with a given time to live (TTL) if key is present
	SetIfPresent(ctx context.Context, r *SetIfPresentRequest) (responses.SetIfPresentResponse, error)
	// SetIfPresentAndNotEqual sets the value in cache with a given time to live (TTL) if key is present and its value is not equal to the given value
	SetIfPresentAndNotEqual(ctx context.Context, r *SetIfPresentAndNotEqualRequest) (responses.SetIfPresentAndNotEqualResponse, error)
	// SetIfEqual sets the value in cache with a given time to live (TTL) if key is present and its value is equal to the given value
	SetIfEqual(ctx context.Context, r *SetIfEqualRequest) (responses.SetIfEqualResponse, error)
	// SetIfAbsentOrEqual sets the value in cache with a given time to live (TTL) if key is present and its value is equal to the given value
	SetIfAbsentOrEqual(ctx context.Context, r *SetIfAbsentOrEqualRequest) (responses.SetIfAbsentOrEqualResponse, error)
	// SetIfNotEqual sets the value in cache with a given time to live (TTL) if key is present and its value is not equal to the given value
	SetIfNotEqual(ctx context.Context, r *SetIfNotEqualRequest) (responses.SetIfNotEqualResponse, error)
	// SetBatch sets multiple values in cache with a given time to live (TTL)
	SetBatch(ctx context.Context, r *SetBatchRequest) (responses.SetBatchResponse, error)
	// Get gets the cache value stored for the given key.
	Get(ctx context.Context, r *GetRequest) (responses.GetResponse, error)
	// GetBatch gets the cache values stored for the given keys.
	GetBatch(ctx context.Context, r *GetBatchRequest) (responses.GetBatchResponse, error)
	// Delete removes the key from the cache.
	Delete(ctx context.Context, r *DeleteRequest) (responses.DeleteResponse, error)
	// KeysExist checks if provided keys exist in the cache.
	KeysExist(ctx context.Context, r *KeysExistRequest) (responses.KeysExistResponse, error)
	// ItemGetType returns the type of the key in the cache
	ItemGetType(ctx context.Context, r *ItemGetTypeRequest) (responses.ItemGetTypeResponse, error)
	// ItemGetTtl returns the TTL for a key in the cache
	ItemGetTtl(ctx context.Context, r *ItemGetTtlRequest) (responses.ItemGetTtlResponse, error)

	// SortedSetFetchByRank fetches the elements in the given sorted set by rank.
	SortedSetFetchByRank(ctx context.Context, r *SortedSetFetchByRankRequest) (responses.SortedSetFetchResponse, error)
	// SortedSetFetchByScore fetches the elements in the given sorted set by score.
	SortedSetFetchByScore(ctx context.Context, r *SortedSetFetchByScoreRequest) (responses.SortedSetFetchResponse, error)
	// SortedSetPutElement adds an element to the given sorted set. If the element already exists,
	// its score is updated. Creates the sorted set if it does not exist.
	SortedSetPutElement(ctx context.Context, r *SortedSetPutElementRequest) (responses.SortedSetPutElementResponse, error)
	// SortedSetPutElements adds elements to the given sorted set. If an element already exists,
	// its score is updated. Creates the sorted set if it does not exist.
	SortedSetPutElements(ctx context.Context, r *SortedSetPutElementsRequest) (responses.SortedSetPutElementsResponse, error)
	// SortedSetGetScore looks up the score of an element in the sorted set, by the value of the elements.
	SortedSetGetScore(ctx context.Context, r *SortedSetGetScoreRequest) (responses.SortedSetGetScoreResponse, error)
	// SortedSetGetScores looks up the scores of multiple elements in the sorted set, by the value of the elements.
	SortedSetGetScores(ctx context.Context, r *SortedSetGetScoresRequest) (responses.SortedSetGetScoresResponse, error)
	// SortedSetRemoveElement removes an element from the sorted set.
	SortedSetRemoveElement(ctx context.Context, r *SortedSetRemoveElementRequest) (responses.SortedSetRemoveElementResponse, error)
	// SortedSetRemoveElements removes elements from the sorted set.
	SortedSetRemoveElements(ctx context.Context, r *SortedSetRemoveElementsRequest) (responses.SortedSetRemoveElementsResponse, error)
	// SortedSetGetRank looks up the rank of an element in the sorted set, by the value of the element.
	SortedSetGetRank(ctx context.Context, r *SortedSetGetRankRequest) (responses.SortedSetGetRankResponse, error)
	// SortedSetLength gets the number of elements in the sorted set.
	SortedSetLength(ctx context.Context, r *SortedSetLengthRequest) (responses.SortedSetLengthResponse, error)
	// SortedSetLengthByScore gets the number of elements in the sorted set by an optional score range.
	SortedSetLengthByScore(ctx context.Context, r *SortedSetLengthByScoreRequest) (responses.SortedSetLengthByScoreResponse, error)
	// SortedSetIncrementScore increments the score of an element in the sorted set.
	SortedSetIncrementScore(ctx context.Context, r *SortedSetIncrementScoreRequest) (responses.SortedSetIncrementScoreResponse, error)

	// SetAddElement adds an element to the given set. Creates the set if it does not already exist.
	SetAddElement(ctx context.Context, r *SetAddElementRequest) (responses.SetAddElementResponse, error)
	// SetAddElements adds multiple elements to the given set. Creates the set if it does not already exist.
	SetAddElements(ctx context.Context, r *SetAddElementsRequest) (responses.SetAddElementsResponse, error)
	// SetFetch fetches the requested set.
	SetFetch(ctx context.Context, r *SetFetchRequest) (responses.SetFetchResponse, error)
	// SetLength gets the number of elements in the set.
	SetLength(ctx context.Context, r *SetLengthRequest) (responses.SetLengthResponse, error)
	// SetRemoveElement removes an element from the given set.
	SetRemoveElement(ctx context.Context, r *SetRemoveElementRequest) (responses.SetRemoveElementResponse, error)
	// SetRemoveElements removes multiple elements from the set.
	SetRemoveElements(ctx context.Context, r *SetRemoveElementsRequest) (responses.SetRemoveElementsResponse, error)
	// SetContainsElements checks if provided elements are in the given set.
	SetContainsElements(ctx context.Context, r *SetContainsElementsRequest) (responses.SetContainsElementsResponse, error)
	// SetPop removes and returns a given number of elements from the given set.
	SetPop(ctx context.Context, r *SetPopRequest) (responses.SetPopResponse, error)

	// ListPushFront adds an element to the front of the given list. Creates the list if it does not already exist.
	ListPushFront(ctx context.Context, r *ListPushFrontRequest) (responses.ListPushFrontResponse, error)
	// ListPushBack adds an element to the back of the given list. Creates the list if it does not already exist.
	ListPushBack(ctx context.Context, r *ListPushBackRequest) (responses.ListPushBackResponse, error)
	// ListPopFront gets and removes the first value from the given list.
	ListPopFront(ctx context.Context, r *ListPopFrontRequest) (responses.ListPopFrontResponse, error)
	// ListPopBack gets and removes the last value from the given list.
	ListPopBack(ctx context.Context, r *ListPopBackRequest) (responses.ListPopBackResponse, error)
	// ListConcatenateFront adds multiple elements to the front of the given list. Creates the list if it does not already exist.
	ListConcatenateFront(ctx context.Context, r *ListConcatenateFrontRequest) (responses.ListConcatenateFrontResponse, error)
	// ListConcatenateBack adds multiple elements to the back of the given list. Creates the list if it does not already exist.
	ListConcatenateBack(ctx context.Context, r *ListConcatenateBackRequest) (responses.ListConcatenateBackResponse, error)
	// ListFetch fetches all elements of the given list.
	ListFetch(ctx context.Context, r *ListFetchRequest) (responses.ListFetchResponse, error)
	// ListLength gets the number of elements in the given list.
	ListLength(ctx context.Context, r *ListLengthRequest) (responses.ListLengthResponse, error)
	// ListRemoveValue removes all elements from the given list equal to the given value.
	ListRemoveValue(ctx context.Context, r *ListRemoveValueRequest) (responses.ListRemoveValueResponse, error)

	// DictionarySetField adds an element to the given dictionary. Creates the dictionary if it does not already exist.
	DictionarySetField(ctx context.Context, r *DictionarySetFieldRequest) (responses.DictionarySetFieldResponse, error)
	// DictionarySetFields adds multiple elements to the given dictionary. Creates the dictionary if it does not already exist.
	//  Use momento.DictionaryElementsFromMap to help construct the Request from a map object.
	DictionarySetFields(ctx context.Context, r *DictionarySetFieldsRequest) (responses.DictionarySetFieldsResponse, error)
	// DictionaryFetch fetches all elements of the given dictionary.
	DictionaryFetch(ctx context.Context, r *DictionaryFetchRequest) (responses.DictionaryFetchResponse, error)
	// DictionaryLength gets the number of items in the given dictionary.
	DictionaryLength(ctx context.Context, r *DictionaryLengthRequest) (responses.DictionaryLengthResponse, error)
	// DictionaryGetField gets the value stored for the given dictionary and field.
	DictionaryGetField(ctx context.Context, r *DictionaryGetFieldRequest) (responses.DictionaryGetFieldResponse, error)
	// DictionaryGetFields gets multiple values from the given dictionary.
	DictionaryGetFields(ctx context.Context, r *DictionaryGetFieldsRequest) (responses.DictionaryGetFieldsResponse, error)
	// DictionaryIncrement adds an integer quantity to a dictionary value.
	// Incrementing the value of a missing field sets the value to amount.
	DictionaryIncrement(ctx context.Context, r *DictionaryIncrementRequest) (responses.DictionaryIncrementResponse, error)
	// DictionaryRemoveField removes an element from the given dictionary.
	// Performs a no-op if the dictionary or field does not exist.
	DictionaryRemoveField(ctx context.Context, r *DictionaryRemoveFieldRequest) (responses.DictionaryRemoveFieldResponse, error)
	// DictionaryRemoveFields removes multiple fields from the given dictionary.
	// Performs a no-op if the dictionary or fields do not exist.
	DictionaryRemoveFields(ctx context.Context, r *DictionaryRemoveFieldsRequest) (responses.DictionaryRemoveFieldsResponse, error)

	// UpdateTtl overwrites the TTL for key to the provided value.
	UpdateTtl(ctx context.Context, r *UpdateTtlRequest) (responses.UpdateTtlResponse, error)
	// IncreaseTtl sets the TTL for a key to the provided value only if it would increase the existing TTL.
	IncreaseTtl(ctx context.Context, r *IncreaseTtlRequest) (responses.IncreaseTtlResponse, error)
	// DecreaseTtl sets the TTL for a key to the provided value only if it would decrease the existing TTL.
	DecreaseTtl(ctx context.Context, r *DecreaseTtlRequest) (responses.DecreaseTtlResponse, error)

	// Ping pings the cache endpoint to check if the service is up and running.
	Ping(ctx context.Context) (responses.PingResponse, error)

	Close()
}

// defaultScsClient represents all information needed for momento client to enable cache control and data operations.
type defaultScsClient struct {
	defaultCache       string
	logger             logger.MomentoLogger
	credentialProvider auth.CredentialProvider
	controlClient      *services.ScsControlClient
	dataClients        []*scsDataClient
	pingClient         *services.ScsPingClient
}

type CacheClientProps struct {
	CacheName string
	// Configuration to use for logging, transport, retries, and middlewares.
	Configuration config.Configuration
	// CredentialProvider Momento credential provider.
	CredentialProvider  auth.CredentialProvider
	DefaultTtl          time.Duration
	EagerConnectTimeout time.Duration
}

func validateEagerConnectionTimeout(timeout time.Duration) momentoerrors.MomentoSvcErr {
	if timeout < 0*time.Second {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "eager connection timeout must be greater than 0", nil)
	}
	return nil
}

func commonCacheClient(props CacheClientProps) (CacheClient, error) {
	if props.Configuration.GetClientSideTimeout() < 1 {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "request timeout must be greater than 0", nil)
	}
	client := &defaultScsClient{
		logger:             props.Configuration.GetLoggerFactory().GetLogger("CacheClient"),
		credentialProvider: props.CredentialProvider,
	}

	controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
		CredentialProvider: props.CredentialProvider,
		Configuration:      props.Configuration,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	if props.DefaultTtl == 0 {
		return nil, convertMomentoSvcErrorToCustomerError(
			momentoerrors.NewMomentoSvcErr(
				momentoerrors.InvalidArgumentError,
				"Must Define a non zero Default TTL", nil),
		)
	}

	err = validateEagerConnectionTimeout(props.EagerConnectTimeout)
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	numChannels := props.Configuration.GetNumGrpcChannels()
	dataClients := make([]*scsDataClient, 0)

	for i := 0; uint32(i) < numChannels; i++ {
		dataClient, err := newScsDataClient(&models.DataClientRequest{
			CredentialProvider: props.CredentialProvider,
			Configuration:      props.Configuration,
			DefaultTtl:         props.DefaultTtl,
		}, props.EagerConnectTimeout)
		if err != nil {
			return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
		}
		dataClients = append(dataClients, dataClient)

		if props.EagerConnectTimeout > 0 {
			err := dataClient.Connect()
			if err != nil {
				logger := props.Configuration.GetLoggerFactory().GetLogger("CacheClient")
				logger.Debug("Failed to connect to the server within the given eager connection timeout:", err.Error())
				return nil, momentoerrors.NewConnectionError(err)
			}
		}
	}

	pingClient, err := services.NewScsPingClient(&models.PingClientRequest{
		Configuration:      props.Configuration,
		CredentialProvider: props.CredentialProvider,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	client.defaultCache = props.CacheName
	client.dataClients = dataClients
	client.controlClient = controlClient
	client.pingClient = pingClient

	return client, nil
}

// NewCacheClient returns a new CacheClient with provided configuration, credential provider, and default TTL seconds arguments.
func NewCacheClient(configuration config.Configuration, credentialProvider auth.CredentialProvider, defaultTtl time.Duration) (CacheClient, error) {
	props := CacheClientProps{
		Configuration:       configuration,
		CredentialProvider:  credentialProvider,
		DefaultTtl:          defaultTtl,
		EagerConnectTimeout: 30 * time.Second,
	}
	return commonCacheClient(props)
}

// NewCacheClientWithEagerConnectTimeout returns a new CacheClient with
// provided configuration, credential provider, and default TTL seconds
// arguments, as well as eagerly attempting to establish gRPC connections.
// A value of 0 for eagerConnectTimeout indicates no eager connections.
func NewCacheClientWithEagerConnectTimeout(configuration config.Configuration, credentialProvider auth.CredentialProvider, defaultTtl time.Duration, eagerConnectTimeout time.Duration) (CacheClient, error) {
	props := CacheClientProps{
		Configuration:       configuration,
		CredentialProvider:  credentialProvider,
		DefaultTtl:          defaultTtl,
		EagerConnectTimeout: eagerConnectTimeout,
	}
	return commonCacheClient(props)
}

func NewCacheClientWithDefaultCache(configuration config.Configuration, credentialProvider auth.CredentialProvider, defaultTtl time.Duration, cacheName string) (CacheClient, error) {
	props := CacheClientProps{
		CacheName:           cacheName,
		Configuration:       configuration,
		CredentialProvider:  credentialProvider,
		DefaultTtl:          defaultTtl,
		EagerConnectTimeout: 30 * time.Second,
	}
	return commonCacheClient(props)
}

func (c defaultScsClient) Logger() logger.MomentoLogger {
	return c.logger
}

func (c defaultScsClient) getCacheNameForRequest(request hasCacheName) string {
	if request.cacheName() != "" {
		return request.cacheName()
	}
	return c.defaultCache
}

func (c defaultScsClient) getNextDataClient() *scsDataClient {
	nextClientIndex := atomic.AddUint64(&dataClientCount, 1)
	dataClient := c.dataClients[nextClientIndex%uint64(len(c.dataClients))]
	return dataClient
}

func (c defaultScsClient) CreateCache(ctx context.Context, request *CreateCacheRequest) (responses.CreateCacheResponse, error) {
	request.CacheName = c.getCacheNameForRequest(request)
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}
	c.logger.Info("Creating cache with name: %s", request.CacheName)
	err := c.controlClient.CreateCache(ctx, &models.CreateCacheRequest{
		CacheName: request.CacheName,
	})
	if err != nil {
		if err.Code() == AlreadyExistsError {
			c.logger.Info("Cache with name '%s' already exists, skipping", request.CacheName)
			return &responses.CreateCacheAlreadyExists{}, nil
		}
		c.logger.Warn("Error creating cache '%s': %s", request.CacheName, err.Message())
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	c.logger.Info("Cache '%s' created successfully", request.CacheName)
	return &responses.CreateCacheSuccess{}, nil
}

func (c defaultScsClient) DeleteCache(ctx context.Context, request *DeleteCacheRequest) (responses.DeleteCacheResponse, error) {
	request.CacheName = c.getCacheNameForRequest(request)
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}
	c.logger.Info("Deleting cache with name: %s", request.CacheName)
	err := c.controlClient.DeleteCache(ctx, &models.DeleteCacheRequest{
		CacheName: request.CacheName,
	})
	if err != nil {
		if err.Code() == NotFoundError {
			c.logger.Info("Cache with name '%s' does not exist, skipping", request.CacheName)
			return &responses.DeleteCacheSuccess{}, nil
		}
		c.logger.Warn("Error deleting cache '%s': %s", request.CacheName, err.Message())
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	c.logger.Info("Cache '%s' deleted successfully", request.CacheName)
	return &responses.DeleteCacheSuccess{}, nil
}

func (c defaultScsClient) ListCaches(ctx context.Context, request *ListCachesRequest) (responses.ListCachesResponse, error) {
	rsp, err := c.controlClient.ListCaches(ctx, &models.ListCachesRequest{})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return responses.NewListCachesSuccess(rsp.NextToken, rsp.Caches), nil
}

func (c defaultScsClient) Increment(ctx context.Context, r *IncrementRequest) (responses.IncrementResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) Set(ctx context.Context, r *SetRequest) (responses.SetResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetIfNotExists(ctx context.Context, r *SetIfNotExistsRequest) (responses.SetIfNotExistsResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetIfAbsent(ctx context.Context, r *SetIfAbsentRequest) (responses.SetIfAbsentResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetIfPresent(ctx context.Context, r *SetIfPresentRequest) (responses.SetIfPresentResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetIfPresentAndNotEqual(ctx context.Context, r *SetIfPresentAndNotEqualRequest) (responses.SetIfPresentAndNotEqualResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetIfEqual(ctx context.Context, r *SetIfEqualRequest) (responses.SetIfEqualResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetIfAbsentOrEqual(ctx context.Context, r *SetIfAbsentOrEqualRequest) (responses.SetIfAbsentOrEqualResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetIfNotEqual(ctx context.Context, r *SetIfNotEqualRequest) (responses.SetIfNotEqualResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetBatch(ctx context.Context, r *SetBatchRequest) (responses.SetBatchResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) Get(ctx context.Context, r *GetRequest) (responses.GetResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) GetBatch(ctx context.Context, r *GetBatchRequest) (responses.GetBatchResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) Delete(ctx context.Context, r *DeleteRequest) (responses.DeleteResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) KeysExist(ctx context.Context, r *KeysExistRequest) (responses.KeysExistResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ItemGetType(ctx context.Context, r *ItemGetTypeRequest) (responses.ItemGetTypeResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ItemGetTtl(ctx context.Context, r *ItemGetTtlRequest) (responses.ItemGetTtlResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetFetchByRank(ctx context.Context, r *SortedSetFetchByRankRequest) (responses.SortedSetFetchResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetFetchByScore(ctx context.Context, r *SortedSetFetchByScoreRequest) (responses.SortedSetFetchResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetPutElement(ctx context.Context, r *SortedSetPutElementRequest) (responses.SortedSetPutElementResponse, error) {
	if r.Value == nil {
		return nil, convertMomentoSvcErrorToCustomerError(
			momentoerrors.NewMomentoSvcErr(
				momentoerrors.InvalidArgumentError, "value cannot be nil", nil,
			),
		)
	}
	r.CacheName = c.getCacheNameForRequest(r)
	newRequest := &SortedSetPutElementsRequest{
		CacheName: r.CacheName,
		SetName:   r.SetName,
		Elements:  []SortedSetElement{{Value: r.Value, Score: r.Score}},
		Ttl:       r.Ttl,
	}
	if err := c.getNextDataClient().makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}

	return &responses.SortedSetPutElementSuccess{}, nil
}

func (c defaultScsClient) SortedSetPutElements(ctx context.Context, r *SortedSetPutElementsRequest) (responses.SortedSetPutElementsResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetGetScores(ctx context.Context, r *SortedSetGetScoresRequest) (responses.SortedSetGetScoresResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetGetScore(ctx context.Context, r *SortedSetGetScoreRequest) (responses.SortedSetGetScoreResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	newRequest := &SortedSetGetScoresRequest{
		CacheName: r.CacheName,
		SetName:   r.SetName,
		Values:    []Value{r.Value},
	}
	if err := c.getNextDataClient().makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	switch result := newRequest.response.(type) {
	case *responses.SortedSetGetScoresHit:
		return result.Responses()[0], nil
	case *responses.SortedSetGetScoresMiss:
		return &responses.SortedSetGetScoreMiss{}, nil
	}
	return nil, errUnexpectedGrpcResponse(newRequest, newRequest.grpcResponse)
}

func (c defaultScsClient) SortedSetRemoveElement(ctx context.Context, r *SortedSetRemoveElementRequest) (responses.SortedSetRemoveElementResponse, error) {
	if r.Value == nil {
		return nil, convertMomentoSvcErrorToCustomerError(
			momentoerrors.NewMomentoSvcErr(
				momentoerrors.InvalidArgumentError, "value cannot be nil", nil,
			),
		)
	}
	r.CacheName = c.getCacheNameForRequest(r)
	newRequest := &SortedSetRemoveElementsRequest{
		CacheName: r.CacheName,
		SetName:   r.SetName,
		Values:    []Value{r.Value},
	}
	if err := c.getNextDataClient().makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	return &responses.SortedSetRemoveElementSuccess{}, nil
}

func (c defaultScsClient) SortedSetRemoveElements(ctx context.Context, r *SortedSetRemoveElementsRequest) (responses.SortedSetRemoveElementsResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetGetRank(ctx context.Context, r *SortedSetGetRankRequest) (responses.SortedSetGetRankResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetLength(ctx context.Context, r *SortedSetLengthRequest) (responses.SortedSetLengthResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetLengthByScore(ctx context.Context, r *SortedSetLengthByScoreRequest) (responses.SortedSetLengthByScoreResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetIncrementScore(ctx context.Context, r *SortedSetIncrementScoreRequest) (responses.SortedSetIncrementScoreResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetAddElement(ctx context.Context, r *SetAddElementRequest) (responses.SetAddElementResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	newRequest := &SetAddElementsRequest{
		CacheName: r.CacheName,
		SetName:   r.SetName,
		Elements:  []Value{r.Element},
		Ttl:       r.Ttl,
	}
	if err := c.getNextDataClient().makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	return &responses.SetAddElementSuccess{}, nil
}

func (c defaultScsClient) SetAddElements(ctx context.Context, r *SetAddElementsRequest) (responses.SetAddElementsResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetFetch(ctx context.Context, r *SetFetchRequest) (responses.SetFetchResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetLength(ctx context.Context, r *SetLengthRequest) (responses.SetLengthResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetRemoveElement(ctx context.Context, r *SetRemoveElementRequest) (responses.SetRemoveElementResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	newRequest := &SetRemoveElementsRequest{
		CacheName: r.CacheName,
		SetName:   r.SetName,
		Elements:  []Value{r.Element},
	}
	if err := c.getNextDataClient().makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	return &responses.SetRemoveElementSuccess{}, nil
}

func (c defaultScsClient) SetRemoveElements(ctx context.Context, r *SetRemoveElementsRequest) (responses.SetRemoveElementsResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetContainsElements(ctx context.Context, r *SetContainsElementsRequest) (responses.SetContainsElementsResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetPop(ctx context.Context, r *SetPopRequest) (responses.SetPopResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListPushFront(ctx context.Context, r *ListPushFrontRequest) (responses.ListPushFrontResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListPushBack(ctx context.Context, r *ListPushBackRequest) (responses.ListPushBackResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListPopFront(ctx context.Context, r *ListPopFrontRequest) (responses.ListPopFrontResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListPopBack(ctx context.Context, r *ListPopBackRequest) (responses.ListPopBackResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListConcatenateFront(ctx context.Context, r *ListConcatenateFrontRequest) (responses.ListConcatenateFrontResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListConcatenateBack(ctx context.Context, r *ListConcatenateBackRequest) (responses.ListConcatenateBackResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListFetch(ctx context.Context, r *ListFetchRequest) (responses.ListFetchResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListLength(ctx context.Context, r *ListLengthRequest) (responses.ListLengthResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListRemoveValue(ctx context.Context, r *ListRemoveValueRequest) (responses.ListRemoveValueResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) DictionarySetField(ctx context.Context, r *DictionarySetFieldRequest) (responses.DictionarySetFieldResponse, error) {
	if r.Field == nil {
		return nil, convertMomentoSvcErrorToCustomerError(
			momentoerrors.NewMomentoSvcErr(
				momentoerrors.InvalidArgumentError, "field cannot be nil", nil,
			),
		)
	}

	r.CacheName = c.getCacheNameForRequest(r)
	var elements []DictionaryElement
	elements = append(elements, DictionaryElement{
		Field: r.Field,
		Value: r.Value,
	})
	newRequest := &DictionarySetFieldsRequest{
		CacheName:      r.CacheName,
		DictionaryName: r.DictionaryName,
		Elements:       elements,
		Ttl:            r.Ttl,
	}
	if err := c.getNextDataClient().makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	return &responses.DictionarySetFieldSuccess{}, nil
}

func (c defaultScsClient) DictionarySetFields(ctx context.Context, r *DictionarySetFieldsRequest) (responses.DictionarySetFieldsResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) DictionaryFetch(ctx context.Context, r *DictionaryFetchRequest) (responses.DictionaryFetchResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) DictionaryLength(ctx context.Context, r *DictionaryLengthRequest) (responses.DictionaryLengthResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) DictionaryGetField(ctx context.Context, r *DictionaryGetFieldRequest) (responses.DictionaryGetFieldResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	newRequest := &DictionaryGetFieldsRequest{
		CacheName:      r.CacheName,
		DictionaryName: r.DictionaryName,
		Fields:         []Value{r.Field},
	}
	if err := c.getNextDataClient().makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	switch rtype := newRequest.response.(type) {
	case *responses.DictionaryGetFieldsMiss:
		return &responses.DictionaryGetFieldMiss{}, nil
	case *responses.DictionaryGetFieldsHit:
		switch rtype.Responses()[0].(type) {
		case *responses.DictionaryGetFieldHit:
			return responses.NewDictionaryGetFieldHitFromFieldsHit(rtype), nil
		case *responses.DictionaryGetFieldMiss:
			return &responses.DictionaryGetFieldMiss{}, nil
		default:
			return nil, errUnexpectedGrpcResponse(newRequest, newRequest.grpcResponse)
		}
	default:
		return nil, errUnexpectedGrpcResponse(newRequest, newRequest.grpcResponse)
	}
}

func (c defaultScsClient) DictionaryGetFields(ctx context.Context, r *DictionaryGetFieldsRequest) (responses.DictionaryGetFieldsResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) DictionaryIncrement(ctx context.Context, r *DictionaryIncrementRequest) (responses.DictionaryIncrementResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) DictionaryRemoveField(ctx context.Context, r *DictionaryRemoveFieldRequest) (responses.DictionaryRemoveFieldResponse, error) {
	if r.Field == nil {
		return nil, convertMomentoSvcErrorToCustomerError(
			momentoerrors.NewMomentoSvcErr(
				momentoerrors.InvalidArgumentError, "field cannot be nil", nil,
			),
		)
	}
	r.CacheName = c.getCacheNameForRequest(r)
	newRequest := &DictionaryRemoveFieldsRequest{
		CacheName:      r.CacheName,
		DictionaryName: r.DictionaryName,
		Fields:         []Value{r.Field},
	}
	if err := c.getNextDataClient().makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	return &responses.DictionaryRemoveFieldSuccess{}, nil
}

func (c defaultScsClient) DictionaryRemoveFields(ctx context.Context, r *DictionaryRemoveFieldsRequest) (responses.DictionaryRemoveFieldsResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) UpdateTtl(ctx context.Context, r *UpdateTtlRequest) (responses.UpdateTtlResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) IncreaseTtl(ctx context.Context, r *IncreaseTtlRequest) (responses.IncreaseTtlResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) DecreaseTtl(ctx context.Context, r *DecreaseTtlRequest) (responses.DecreaseTtlResponse, error) {
	r.CacheName = c.getCacheNameForRequest(r)
	if err := c.getNextDataClient().makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) Ping(ctx context.Context) (responses.PingResponse, error) {
	if err := c.pingClient.Ping(ctx); err != nil {
		return nil, err
	}
	return &responses.PingSuccess{}, nil
}

func (c defaultScsClient) Close() {
	defer c.controlClient.Close()
	defer c.getNextDataClient().Close()
}

func convertMomentoSvcErrorToCustomerError(e momentoerrors.MomentoSvcErr) MomentoError {
	if e == nil {
		return nil
	}
	return NewMomentoError(e.Code(), e.Message(), e.OriginalErr())
}

func isCacheNameValid(cacheName string) momentoerrors.MomentoSvcErr {
	if len(strings.TrimSpace(cacheName)) < 1 {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Cache name cannot be empty", nil)
	}
	return nil
}
