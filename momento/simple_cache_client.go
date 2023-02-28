// Package momento represents API CacheClient interface accessors including control/data operations, errors, operation requests and responses for the SDK.
package momento

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/internal/services"
)

type CacheClient interface {
	CreateCache(ctx context.Context, request *CreateCacheRequest) (CreateCacheResponse, error)
	DeleteCache(ctx context.Context, request *DeleteCacheRequest) (DeleteCacheResponse, error)
	ListCaches(ctx context.Context, request *ListCachesRequest) (ListCachesResponse, error)

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
	SortedSetIncrementScore(ctx context.Context, r *SortedSetIncrementScoreRequest) (SortedSetIncrementScoreResponse, error)

	SetAddElement(ctx context.Context, r *SetAddElementRequest) (SetAddElementResponse, error)
	SetAddElements(ctx context.Context, r *SetAddElementsRequest) (SetAddElementsResponse, error)
	SetFetch(ctx context.Context, r *SetFetchRequest) (SetFetchResponse, error)
	SetRemoveElement(ctx context.Context, r *SetRemoveElementRequest) (SetRemoveElementResponse, error)
	SetRemoveElements(ctx context.Context, r *SetRemoveElementsRequest) (SetRemoveElementsResponse, error)

	ListPushFront(ctx context.Context, r *ListPushFrontRequest) (ListPushFrontResponse, error)
	ListPushBack(ctx context.Context, r *ListPushBackRequest) (ListPushBackResponse, error)
	ListPopFront(ctx context.Context, r *ListPopFrontRequest) (ListPopFrontResponse, error)
	ListPopBack(ctx context.Context, r *ListPopBackRequest) (ListPopBackResponse, error)
	ListConcatenateFront(ctx context.Context, r *ListConcatenateFrontRequest) (ListConcatenateFrontResponse, error)
	ListConcatenateBack(ctx context.Context, r *ListConcatenateBackRequest) (ListConcatenateBackResponse, error)
	ListFetch(ctx context.Context, r *ListFetchRequest) (ListFetchResponse, error)
	ListLength(ctx context.Context, r *ListLengthRequest) (ListLengthResponse, error)
	ListRemoveValue(ctx context.Context, r *ListRemoveValueRequest) (ListRemoveValueResponse, error)

	DictionarySetField(ctx context.Context, r *DictionarySetFieldRequest) (DictionarySetFieldResponse, error)
	DictionarySetFields(ctx context.Context, r *DictionarySetFieldsRequest) (DictionarySetFieldsResponse, error)
	DictionaryFetch(ctx context.Context, r *DictionaryFetchRequest) (DictionaryFetchResponse, error)
	DictionaryGetField(ctx context.Context, r *DictionaryGetFieldRequest) (DictionaryGetFieldResponse, error)
	DictionaryGetFields(ctx context.Context, r *DictionaryGetFieldsRequest) (DictionaryGetFieldsResponse, error)
	DictionaryIncrement(ctx context.Context, r *DictionaryIncrementRequest) (DictionaryIncrementResponse, error)
	DictionaryRemoveField(ctx context.Context, r *DictionaryRemoveFieldRequest) (DictionaryRemoveFieldResponse, error)
	DictionaryRemoveFields(ctx context.Context, r *DictionaryRemoveFieldsRequest) (DictionaryRemoveFieldsResponse, error)

	Close()
}

// defaultScsClient represents all information needed for momento client to enable cache control and data operations.
type defaultScsClient struct {
	credentialProvider auth.CredentialProvider
	controlClient      *services.ScsControlClient
	dataClient         *scsDataClient
	pubSubClient       *pubSubClient
}

type CacheClientProps struct {
	Configuration      config.Configuration
	CredentialProvider auth.CredentialProvider
	DefaultTTL         time.Duration
}

// NewCacheClient returns a new CacheClient with provided authToken, DefaultTTLSeconds, and opts arguments.
func NewCacheClient(props *CacheClientProps) (CacheClient, error) {
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

func (c defaultScsClient) CreateCache(ctx context.Context, request *CreateCacheRequest) (CreateCacheResponse, error) {
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

func (c defaultScsClient) DeleteCache(ctx context.Context, request *DeleteCacheRequest) (DeleteCacheResponse, error) {
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

func (c defaultScsClient) ListCaches(ctx context.Context, request *ListCachesRequest) (ListCachesResponse, error) {
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

func (c defaultScsClient) TopicPublish(ctx context.Context, request *TopicPublishRequest) (TopicPublishResponse, error) {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}

	if _, err := prepareName(request.TopicName, "Topic name"); err != nil {
		return nil, err
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

func (c defaultScsClient) SortedSetIncrementScore(ctx context.Context, r *SortedSetIncrementScoreRequest) (SortedSetIncrementScoreResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetAddElements(ctx context.Context, r *SetAddElementsRequest) (SetAddElementsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetAddElement(ctx context.Context, r *SetAddElementRequest) (SetAddElementResponse, error) {
	newRequest := &SetAddElementsRequest{
		CacheName: r.CacheName,
		SetName:   r.SetName,
		Elements:  []Value{r.Element},
		Ttl:       r.Ttl,
	}
	if err := c.dataClient.makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	return &SetAddElementSuccess{}, nil
}

func (c defaultScsClient) SetFetch(ctx context.Context, r *SetFetchRequest) (SetFetchResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetRemoveElements(ctx context.Context, r *SetRemoveElementsRequest) (SetRemoveElementsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) SetRemoveElement(ctx context.Context, r *SetRemoveElementRequest) (SetRemoveElementResponse, error) {
	newRequest := &SetRemoveElementsRequest{
		CacheName: r.CacheName,
		SetName:   r.SetName,
		Elements:  []Value{r.Element},
	}
	if err := c.dataClient.makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	return &SetRemoveElementSuccess{}, nil
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

func (c defaultScsClient) DictionarySetField(ctx context.Context, r *DictionarySetFieldRequest) (DictionarySetFieldResponse, error) {
	elements := make(map[string]Value)
	elements[string(r.Field.asBytes())] = r.Value
	newRequest := &DictionarySetFieldsRequest{
		CacheName:      r.CacheName,
		DictionaryName: r.DictionaryName,
		Elements:       elements,
		Ttl:            r.Ttl,
	}
	if err := c.dataClient.makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	return &DictionarySetFieldSuccess{}, nil
}

func (c defaultScsClient) DictionarySetFields(ctx context.Context, r *DictionarySetFieldsRequest) (DictionarySetFieldsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) DictionaryFetch(ctx context.Context, r *DictionaryFetchRequest) (DictionaryFetchResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) DictionaryGetField(ctx context.Context, r *DictionaryGetFieldRequest) (DictionaryGetFieldResponse, error) {
	newRequest := &DictionaryGetFieldsRequest{
		CacheName:      r.CacheName,
		DictionaryName: r.DictionaryName,
		Fields:         []Value{r.Field},
	}
	if err := c.dataClient.makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	switch rtype := newRequest.response.(type) {
	case *DictionaryGetFieldsMiss:
		return &DictionaryGetFieldMiss{}, nil
	case *DictionaryGetFieldsHit:
		switch rtype.responses[0].(type) {
		case *DictionaryGetFieldHit:
			return &DictionaryGetFieldHit{
				field: rtype.fields[0],
				body:  rtype.elements[0].CacheBody,
			}, nil
		case *DictionaryGetFieldMiss:
			return &DictionaryGetFieldMiss{}, nil
		default:
			return nil, errUnexpectedGrpcResponse(newRequest, newRequest.grpcResponse)
		}
	default:
		return nil, errUnexpectedGrpcResponse(newRequest, newRequest.grpcResponse)
	}
}

func (c defaultScsClient) DictionaryGetFields(ctx context.Context, r *DictionaryGetFieldsRequest) (DictionaryGetFieldsResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) DictionaryIncrement(ctx context.Context, r *DictionaryIncrementRequest) (DictionaryIncrementResponse, error) {
	if err := c.dataClient.makeRequest(ctx, r); err != nil {
		return nil, err
	}
	return r.response, nil
}

func (c defaultScsClient) DictionaryRemoveField(ctx context.Context, r *DictionaryRemoveFieldRequest) (DictionaryRemoveFieldResponse, error) {
	newRequest := &DictionaryRemoveFieldsRequest{
		CacheName:      r.CacheName,
		DictionaryName: r.DictionaryName,
		Fields:         []Value{r.Field},
	}
	if err := c.dataClient.makeRequest(ctx, newRequest); err != nil {
		return nil, err
	}
	return &DictionaryRemoveFieldSuccess{}, nil
}

func (c defaultScsClient) DictionaryRemoveFields(ctx context.Context, r *DictionaryRemoveFieldsRequest) (DictionaryRemoveFieldsResponse, error) {
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
