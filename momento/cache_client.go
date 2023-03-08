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

	Close()
}

// defaultScsClient represents all information needed for momento client to enable cache control and data operations.
type defaultScsClient struct {
	credentialProvider auth.CredentialProvider
	controlClient      *services.ScsControlClient
	dataClient         *scsDataClient
}

type CacheClientProps struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
	DefaultTtl         time.Duration
}

// NewCacheClient returns a new CacheClient with provided configuration, credential provider, and default TTL seconds arguments.
func NewCacheClient(configuration config.Configuration, credentialProvider auth.CredentialProvider, defaultTtl time.Duration) (CacheClient, error) {
	props := CacheClientProps{
		Configuration:      configuration,
		CredentialProvider: credentialProvider,
		DefaultTtl:         defaultTtl,
	}
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

	client.dataClient = dataClient
	client.controlClient = controlClient

	return client, nil
}

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

func (c defaultScsClient) ListCaches(ctx context.Context, request *ListCachesRequest) (responses.ListCachesResponse, error) {
	rsp, err := c.controlClient.ListCaches(ctx, &models.ListCachesRequest{
		NextToken: request.NextToken,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return responses.NewListCachesSuccess(rsp.NextToken, rsp.Caches), nil
}

func (c defaultScsClient) Set(ctx context.Context, r *SetRequest) (responses.SetResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) Get(ctx context.Context, r *GetRequest) (responses.GetResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) Delete(ctx context.Context, r *DeleteRequest) (responses.DeleteResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetFetch(ctx context.Context, r *SortedSetFetchRequest) (responses.SortedSetFetchResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetPut(ctx context.Context, r *SortedSetPutRequest) (responses.SortedSetPutResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetGetScores(ctx context.Context, r *SortedSetGetScoresRequest) (responses.SortedSetGetScoresResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetRemove(ctx context.Context, r *SortedSetRemoveRequest) (responses.SortedSetRemoveResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetGetRank(ctx context.Context, r *SortedSetGetRankRequest) (responses.SortedSetGetRankResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SortedSetIncrementScore(ctx context.Context, r *SortedSetIncrementScoreRequest) (responses.SortedSetIncrementScoreResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetAddElements(ctx context.Context, r *SetAddElementsRequest) (responses.SetAddElementsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

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

func (c defaultScsClient) SetFetch(ctx context.Context, r *SetFetchRequest) (responses.SetFetchResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetRemoveElements(ctx context.Context, r *SetRemoveElementsRequest) (responses.SetRemoveElementsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

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

func (c defaultScsClient) SetContainsElements(ctx context.Context, r *SetContainsElementsRequest) (responses.SetContainsElementsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListPushFront(ctx context.Context, r *ListPushFrontRequest) (responses.ListPushFrontResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListPushBack(ctx context.Context, r *ListPushBackRequest) (responses.ListPushBackResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListPopFront(ctx context.Context, r *ListPopFrontRequest) (responses.ListPopFrontResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListPopBack(ctx context.Context, r *ListPopBackRequest) (responses.ListPopBackResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListConcatenateFront(ctx context.Context, r *ListConcatenateFrontRequest) (responses.ListConcatenateFrontResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListConcatenateBack(ctx context.Context, r *ListConcatenateBackRequest) (responses.ListConcatenateBackResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListFetch(ctx context.Context, r *ListFetchRequest) (responses.ListFetchResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListLength(ctx context.Context, r *ListLengthRequest) (responses.ListLengthResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) ListRemoveValue(ctx context.Context, r *ListRemoveValueRequest) (responses.ListRemoveValueResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
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

func (c defaultScsClient) DictionarySetFields(ctx context.Context, r *DictionarySetFieldsRequest) (responses.DictionarySetFieldsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) DictionaryFetch(ctx context.Context, r *DictionaryFetchRequest) (responses.DictionaryFetchResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

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

func (c defaultScsClient) DictionaryGetFields(ctx context.Context, r *DictionaryGetFieldsRequest) (responses.DictionaryGetFieldsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) DictionaryIncrement(ctx context.Context, r *DictionaryIncrementRequest) (responses.DictionaryIncrementResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
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

func (c defaultScsClient) DictionaryRemoveFields(ctx context.Context, r *DictionaryRemoveFieldsRequest) (responses.DictionaryRemoveFieldsResponse, error) {
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

func isCacheNameValid(cacheName string) momentoerrors.MomentoSvcErr {
	if len(strings.TrimSpace(cacheName)) < 1 {
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "Cache name cannot be empty", nil)
	}
	return nil
}
