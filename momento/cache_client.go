// Package momento represents API CacheClient interface accessors including control/data operations, errors, operation requests and responses for the SDK.
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
	"github.com/momentohq/client-sdk-go/responses"
)

type CacheClient interface {
	CreateCache(ctx context.Context, request *CreateCacheRequest) (responses.CreateCacheResponse, error)
	DeleteCache(ctx context.Context, request *DeleteCacheRequest) (responses.DeleteCacheResponse, error)
	ListCaches(ctx context.Context, request *ListCachesRequest) (responses.ListCachesResponse, error)

	Set(ctx context.Context, r *SetRequest) (responses.SetResponse, error)
	Get(ctx context.Context, r *GetRequest) (responses.GetResponse, error)
	Delete(ctx context.Context, r *DeleteRequest) (responses.DeleteResponse, error)
	KeysExist(ctx context.Context, r *KeysExistRequest) (responses.KeysExistResponse, error)

	SortedSetFetch(ctx context.Context, r *SortedSetFetchRequest) (responses.SortedSetFetchResponse, error)
	SortedSetPut(ctx context.Context, r *SortedSetPutRequest) (responses.SortedSetPutResponse, error)
	SortedSetGetScores(ctx context.Context, r *SortedSetGetScoresRequest) (responses.SortedSetGetScoresResponse, error)
	SortedSetRemove(ctx context.Context, r *SortedSetRemoveRequest) (responses.SortedSetRemoveResponse, error)
	SortedSetGetRank(ctx context.Context, r *SortedSetGetRankRequest) (responses.SortedSetGetRankResponse, error)
	SortedSetIncrementScore(ctx context.Context, r *SortedSetIncrementScoreRequest) (responses.SortedSetIncrementScoreResponse, error)

	SetAddElement(ctx context.Context, r *SetAddElementRequest) (responses.SetAddElementResponse, error)
	SetAddElements(ctx context.Context, r *SetAddElementsRequest) (responses.SetAddElementsResponse, error)
	SetFetch(ctx context.Context, r *SetFetchRequest) (responses.SetFetchResponse, error)
	SetRemoveElement(ctx context.Context, r *SetRemoveElementRequest) (responses.SetRemoveElementResponse, error)
	SetRemoveElements(ctx context.Context, r *SetRemoveElementsRequest) (responses.SetRemoveElementsResponse, error)
	SetContainsElements(ctx context.Context, r *SetContainsElementsRequest) (responses.SetContainsElementsResponse, error)

	ListPushFront(ctx context.Context, r *ListPushFrontRequest) (responses.ListPushFrontResponse, error)
	ListPushBack(ctx context.Context, r *ListPushBackRequest) (responses.ListPushBackResponse, error)
	ListPopFront(ctx context.Context, r *ListPopFrontRequest) (responses.ListPopFrontResponse, error)
	ListPopBack(ctx context.Context, r *ListPopBackRequest) (responses.ListPopBackResponse, error)
	ListConcatenateFront(ctx context.Context, r *ListConcatenateFrontRequest) (responses.ListConcatenateFrontResponse, error)
	ListConcatenateBack(ctx context.Context, r *ListConcatenateBackRequest) (responses.ListConcatenateBackResponse, error)
	ListFetch(ctx context.Context, r *ListFetchRequest) (responses.ListFetchResponse, error)
	ListLength(ctx context.Context, r *ListLengthRequest) (responses.ListLengthResponse, error)
	ListRemoveValue(ctx context.Context, r *ListRemoveValueRequest) (responses.ListRemoveValueResponse, error)

	DictionarySetField(ctx context.Context, r *DictionarySetFieldRequest) (responses.DictionarySetFieldResponse, error)
	DictionarySetFields(ctx context.Context, r *DictionarySetFieldsRequest) (responses.DictionarySetFieldsResponse, error)
	DictionaryFetch(ctx context.Context, r *DictionaryFetchRequest) (responses.DictionaryFetchResponse, error)
	DictionaryGetField(ctx context.Context, r *DictionaryGetFieldRequest) (responses.DictionaryGetFieldResponse, error)
	DictionaryGetFields(ctx context.Context, r *DictionaryGetFieldsRequest) (responses.DictionaryGetFieldsResponse, error)
	DictionaryIncrement(ctx context.Context, r *DictionaryIncrementRequest) (responses.DictionaryIncrementResponse, error)
	DictionaryRemoveField(ctx context.Context, r *DictionaryRemoveFieldRequest) (responses.DictionaryRemoveFieldResponse, error)
	DictionaryRemoveFields(ctx context.Context, r *DictionaryRemoveFieldsRequest) (responses.DictionaryRemoveFieldsResponse, error)

	UpdateTtl(ctx context.Context, r *UpdateTtlRequest) (responses.UpdateTtlResponse, error)
	IncreaseTtl(ctx context.Context, r *IncreaseTtlRequest) (responses.IncreaseTtlResponse, error)
	DecreaseTtl(ctx context.Context, r *DecreaseTtlRequest) (responses.DecreaseTtlResponse, error)

	Ping(ctx context.Context) (responses.PingResponse, error)

	Close()
}

// defaultScsClient represents all information needed for momento client to enable cache control and data operations.
type defaultScsClient struct {
	credentialProvider auth.CredentialProvider
	controlClient      *services.ScsControlClient
	dataClient         *scsDataClient
	pingClient         *services.ScsPingClient
}

type CacheClientProps struct {
	// Configuration to use for the transport, retries, middlewares.
	Configuration config.Configuration
	// CredentialProvider Momento credential provider.
	CredentialProvider auth.CredentialProvider
	// DefaultTtl is default time to live for the item in cache.
	DefaultTtl time.Duration
}

// NewCacheClient returns a new CacheClient with provided configuration, credential provider, and default TTL seconds arguments.
func NewCacheClient(configuration config.Configuration, credentialProvider auth.CredentialProvider, defaultTtl time.Duration) (CacheClient, error) {
	props := CacheClientProps{
		Configuration:      configuration,
		CredentialProvider: credentialProvider,
		DefaultTtl:         defaultTtl,
	}
	if props.Configuration.GetClientSideTimeout() < 1 {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "request timeout must be greater than 0", nil)
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

	if props.DefaultTtl == 0 {
		return nil, convertMomentoSvcErrorToCustomerError(
			momentoerrors.NewMomentoSvcErr(
				momentoerrors.InvalidArgumentError,
				"Must Define a non zero Default TTL", nil),
		)
	}

	dataClient, err := newScsDataClient(&models.DataClientRequest{
		CredentialProvider: props.CredentialProvider,
		Configuration:      props.Configuration,
		DefaultTtl:         props.DefaultTtl,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	pingClient, err := services.NewScsPingClient(&models.PingClientRequest{
		Configuration:      props.Configuration,
		CredentialProvider: props.CredentialProvider,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.ConvertSvcErr(err))
	}

	client.dataClient = dataClient
	client.controlClient = controlClient
	client.pingClient = pingClient

	return client, nil
}

// CreateCache Creates a cache if it does not exist.
func (c defaultScsClient) CreateCache(ctx context.Context, request *CreateCacheRequest) (responses.CreateCacheResponse, error) {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}
	err := c.controlClient.CreateCache(ctx, &models.CreateCacheRequest{
		CacheName: request.CacheName,
	})
	if err != nil {
		if err.Code() == AlreadyExistsError {
			return &responses.CreateCacheAlreadyExists{}, nil
		}
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return &responses.CreateCacheSuccess{}, nil
}

// DeleteCache deletes a cache and all of the items within it.
func (c defaultScsClient) DeleteCache(ctx context.Context, request *DeleteCacheRequest) (responses.DeleteCacheResponse, error) {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}
	err := c.controlClient.DeleteCache(ctx, &models.DeleteCacheRequest{
		CacheName: request.CacheName,
	})
	if err != nil {
		if err.Code() == NotFoundError {
			return &responses.DeleteCacheSuccess{}, nil
		}
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return &responses.DeleteCacheSuccess{}, nil
}

// ListCaches lists all caches.
func (c defaultScsClient) ListCaches(ctx context.Context, request *ListCachesRequest) (responses.ListCachesResponse, error) {
	rsp, err := c.controlClient.ListCaches(ctx, &models.ListCachesRequest{
		NextToken: request.NextToken,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return responses.NewListCachesSuccess(rsp.NextToken, rsp.Caches), nil
}

// Set sets the value in cache with a given time to live (TTL) seconds
func (c defaultScsClient) Set(ctx context.Context, r *SetRequest) (responses.SetResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// Get gets the cache value stored for the given key.
func (c defaultScsClient) Get(ctx context.Context, r *GetRequest) (responses.GetResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// Delete removes the key from the cache.
func (c defaultScsClient) Delete(ctx context.Context, r *DeleteRequest) (responses.DeleteResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// KeysExist checks if provided keys exist in the cache.
func (c defaultScsClient) KeysExist(ctx context.Context, r *KeysExistRequest) (responses.KeysExistResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// SortedSetFetch fetches the elements in the given sorted set by index (rank).
func (c defaultScsClient) SortedSetFetch(ctx context.Context, r *SortedSetFetchRequest) (responses.SortedSetFetchResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// SortedSetPut adds elements to the given sorted set. If the element already exists,
// its score is updated. Creates the sorted set if it does not exist.
func (c defaultScsClient) SortedSetPut(ctx context.Context, r *SortedSetPutRequest) (responses.SortedSetPutResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// SortedSetGetScores looks up the scores of multiple elements in the sorted set, by the value of the elements.
func (c defaultScsClient) SortedSetGetScores(ctx context.Context, r *SortedSetGetScoresRequest) (responses.SortedSetGetScoresResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// SortedSetRemove removes elements from the sorted set
func (c defaultScsClient) SortedSetRemove(ctx context.Context, r *SortedSetRemoveRequest) (responses.SortedSetRemoveResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// SortedSetGetRank looks up the rank of an element in the sorted set, by the value of the element.
func (c defaultScsClient) SortedSetGetRank(ctx context.Context, r *SortedSetGetRankRequest) (responses.SortedSetGetRankResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// SortedSetIncrementScore increments the score of an element in the sorted set.
func (c defaultScsClient) SortedSetIncrementScore(ctx context.Context, r *SortedSetIncrementScoreRequest) (responses.SortedSetIncrementScoreResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// SetAddElements adds multiple elements to the given set. Creates the set if it does not already exist.
// After this operation the set will contain the union of the element passed in and the original elements of the set.
func (c defaultScsClient) SetAddElements(ctx context.Context, r *SetAddElementsRequest) (responses.SetAddElementsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// SetAddElement adds an element to the given set. Creates the set if it does not already exist.
// After this operation the set will contain the union of the element passed in and the original elements of the set.
func (c defaultScsClient) SetAddElement(ctx context.Context, r *SetAddElementRequest) (responses.SetAddElementResponse, error) {
	newRequest := &SetAddElementsRequest{
		CacheName: r.CacheName,
		SetName:   r.SetName,
		Elements:  []Value{r.Element},
		Ttl:       r.Ttl,
	}
	if err := c.dataClient.makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	return &responses.SetAddElementSuccess{}, nil
}

// SetFetch fetches the elements in the given sorted set by index (rank).
func (c defaultScsClient) SetFetch(ctx context.Context, r *SetFetchRequest) (responses.SetFetchResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// SetRemoveElements removes multiple elements from the sorted set.
func (c defaultScsClient) SetRemoveElements(ctx context.Context, r *SetRemoveElementsRequest) (responses.SetRemoveElementsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// SetRemoveElement removes an element from the given set.
func (c defaultScsClient) SetRemoveElement(ctx context.Context, r *SetRemoveElementRequest) (responses.SetRemoveElementResponse, error) {
	newRequest := &SetRemoveElementsRequest{
		CacheName: r.CacheName,
		SetName:   r.SetName,
		Elements:  []Value{r.Element},
	}
	if err := c.dataClient.makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	return &responses.SetRemoveElementSuccess{}, nil
}

// SetContainsElements checks if provided elements are in the given set.
func (c defaultScsClient) SetContainsElements(ctx context.Context, r *SetContainsElementsRequest) (responses.SetContainsElementsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// ListPushFront adds an element to the front of the given list. Creates the list if it does not already exist.
func (c defaultScsClient) ListPushFront(ctx context.Context, r *ListPushFrontRequest) (responses.ListPushFrontResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// ListPushBack adds an element to the back of the given list. Creates the list if it does not already exist.
func (c defaultScsClient) ListPushBack(ctx context.Context, r *ListPushBackRequest) (responses.ListPushBackResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// ListPopFront gets and removes the first value from the given list.
func (c defaultScsClient) ListPopFront(ctx context.Context, r *ListPopFrontRequest) (responses.ListPopFrontResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// ListPopBack gets and removes the last value from the given list.
func (c defaultScsClient) ListPopBack(ctx context.Context, r *ListPopBackRequest) (responses.ListPopBackResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// ListConcatenateFront adds multiple elements to the front of the given list. Creates the list if it does not already exist.
func (c defaultScsClient) ListConcatenateFront(ctx context.Context, r *ListConcatenateFrontRequest) (responses.ListConcatenateFrontResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// ListConcatenateBack adds multiple elements to the back of the given list. Creates the list if it does not already exist.
func (c defaultScsClient) ListConcatenateBack(ctx context.Context, r *ListConcatenateBackRequest) (responses.ListConcatenateBackResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// ListFetch fetched all elements of the given list.
func (c defaultScsClient) ListFetch(ctx context.Context, r *ListFetchRequest) (responses.ListFetchResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// ListLength gets the number of elements in the given list.
func (c defaultScsClient) ListLength(ctx context.Context, r *ListLengthRequest) (responses.ListLengthResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// ListRemoveValue removes all elements from the given list equal to the given value.
func (c defaultScsClient) ListRemoveValue(ctx context.Context, r *ListRemoveValueRequest) (responses.ListRemoveValueResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// DictionarySetField adds an element to the given dictionary. Creates the dictionary if it does not already exist.
func (c defaultScsClient) DictionarySetField(ctx context.Context, r *DictionarySetFieldRequest) (responses.DictionarySetFieldResponse, error) {
	if r.Field == nil {
		return nil, convertMomentoSvcErrorToCustomerError(
			momentoerrors.NewMomentoSvcErr(
				momentoerrors.InvalidArgumentError, "field cannot be nil", nil,
			),
		)
	}

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
	if err := c.dataClient.makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	return &responses.DictionarySetFieldSuccess{}, nil
}

// DictionarySetFields adds multiple elements to the given dictionary. Creates the dictionary if it does not already exist.
func (c defaultScsClient) DictionarySetFields(ctx context.Context, r *DictionarySetFieldsRequest) (responses.DictionarySetFieldsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// DictionaryFetch fetches all elements of the given dictionary.
func (c defaultScsClient) DictionaryFetch(ctx context.Context, r *DictionaryFetchRequest) (responses.DictionaryFetchResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// DictionaryGetField gets the value stored for the given dictionary and field.
func (c defaultScsClient) DictionaryGetField(ctx context.Context, r *DictionaryGetFieldRequest) (responses.DictionaryGetFieldResponse, error) {
	newRequest := &DictionaryGetFieldsRequest{
		CacheName:      r.CacheName,
		DictionaryName: r.DictionaryName,
		Fields:         []Value{r.Field},
	}
	if err := c.dataClient.makeRequest(ctx, newRequest); err != nil {
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

// DictionaryGetFields gets multiple values from the given dictionary.
func (c defaultScsClient) DictionaryGetFields(ctx context.Context, r *DictionaryGetFieldsRequest) (responses.DictionaryGetFieldsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// DictionaryIncrement adds an integer quantity to a dictionary value.
// Incrementing the value of a missing field sets the value to amount.
func (c defaultScsClient) DictionaryIncrement(ctx context.Context, r *DictionaryIncrementRequest) (responses.DictionaryIncrementResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// DictionaryRemoveField removes an element from the given dictionary.
// Performs a no-op if the dictionary or field does not exist.
func (c defaultScsClient) DictionaryRemoveField(ctx context.Context, r *DictionaryRemoveFieldRequest) (responses.DictionaryRemoveFieldResponse, error) {
	if r.Field == nil {
		return nil, convertMomentoSvcErrorToCustomerError(
			momentoerrors.NewMomentoSvcErr(
				momentoerrors.InvalidArgumentError, "field cannot be nil", nil,
			),
		)
	}
	newRequest := &DictionaryRemoveFieldsRequest{
		CacheName:      r.CacheName,
		DictionaryName: r.DictionaryName,
		Fields:         []Value{r.Field},
	}
	if err := c.dataClient.makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	return &responses.DictionaryRemoveFieldSuccess{}, nil
}

// DictionaryRemoveFields removes multiple fields from the given dictionary.
// Performs a no-op if the dictionary or fields do not exist.
func (c defaultScsClient) DictionaryRemoveFields(ctx context.Context, r *DictionaryRemoveFieldsRequest) (responses.DictionaryRemoveFieldsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// UpdateTtl overwrites the TTL for key to the provided value.
func (c defaultScsClient) UpdateTtl(ctx context.Context, r *UpdateTtlRequest) (responses.UpdateTtlResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// IncreaseTtl sets the TTL for a key to the provided value only if it would increase the existing TTL.
func (c defaultScsClient) IncreaseTtl(ctx context.Context, r *IncreaseTtlRequest) (responses.IncreaseTtlResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// DecreaseTtl sets the TTL for a key to the provided value only if it would decrease the existing TTL.
func (c defaultScsClient) DecreaseTtl(ctx context.Context, r *DecreaseTtlRequest) (responses.DecreaseTtlResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

// Ping pings the cache endpoint to check if the service is up and running.
func (c defaultScsClient) Ping(ctx context.Context) (responses.PingResponse, error) {
	if err := c.pingClient.Ping(ctx); err != nil {
		return nil, err
	}
	return &responses.PingSuccess{}, nil
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

func isCacheNameValid(cacheName string) momentoerrors.MomentoSvcErr {
	if len(strings.TrimSpace(cacheName)) < 1 {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Cache name cannot be empty", nil)
	}
	return nil
}
