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
	ListPushFront(ctx context.Context, request *ListPushFrontRequest) (ListPushFrontResponse, error)
	ListPushBack(ctx context.Context, request *ListPushBackRequest) (ListPushBackResponse, error)
	ListPopFront(ctx context.Context, request *ListPopFrontRequest) (ListPopFrontResponse, error)
	ListPopBack(ctx context.Context, request *ListPopBackRequest) (ListPopBackResponse, error)

	SortedSetPut(ctx context.Context, request *SortedSetPutRequest) error
	SortedSetFetch(ctx context.Context, request *SortedSetFetchRequest) (SortedSetFetchResponse, error)
	SortedSetGetScore(ctx context.Context, request *SortedSetGetScoreRequest) (SortedSetGetScoreResponse, error)
	SortedSetRemove(ctx context.Context, request *SortedSetRemoveRequest) error
	SortedSetGetRank(ctx context.Context, request *SortedSetGetRankRequest) (SortedSetGetRankResponse, error)
	// TODO need to impl sortedset increment still
	//SortedSetIncrement(ctx context.Context, request *SortedSetIncrementRequest)

	Close()
}

// DefaultScsClient default implementation of the Momento incubating ScsClient interface
type DefaultScsClient struct {
	controlClient  *services.ScsControlClient
	dataClient     *services.ScsDataClient
	pubSubClient   *services.PubSubClient
	internalClient momento.ScsClient
}

// NewScsClient returns a new ScsClient with provided authToken, defaultTTL,, and opts arguments.
func NewScsClient(props *momento.SimpleCacheClientProps) (ScsClient, error) {

	controlClient, err := services.NewScsControlClient(&models.ControlClientRequest{
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

	dataClient, err := services.NewScsDataClient(&models.DataClientRequest{
		CredentialProvider: props.CredentialProvider,
		Configuration:      props.Configuration,
		DefaultTtl:         props.DefaultTTL,
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
	// TODO: validate list name
	rsp, err := c.dataClient.ListFetch(ctx, &models.ListFetchRequest{
		CacheName: request.CacheName,
		ListName:  request.ListName,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return convertListFetchResponse(rsp)
}

func (c *DefaultScsClient) ListLength(ctx context.Context, request *ListLengthRequest) (ListLengthResponse, error) {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}
	// TODO: validate list name
	rsp, err := c.dataClient.ListLength(ctx, &models.ListLengthRequest{
		CacheName: request.CacheName,
		ListName:  request.ListName,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return convertListLengthResponse(rsp)
}

func (c *DefaultScsClient) ListPushFront(ctx context.Context, request *ListPushFrontRequest) (ListPushFrontResponse, error) {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}
	// TODO: validate list name
	rsp, err := c.dataClient.ListPushFront(ctx, &models.ListPushFrontRequest{
		CacheName:          request.CacheName,
		ListName:           request.ListName,
		Value:              request.Value,
		TruncateBackToSize: request.TruncateBackToSize,
		CollectionTtl:      request.CollectionTTL,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return convertListPushFrontResponse(rsp)
}

func (c *DefaultScsClient) ListPushBack(ctx context.Context, request *ListPushBackRequest) (ListPushBackResponse, error) {
	if err := isCacheNameValid(request.CacheName); err != nil {
		return nil, err
	}
	// TODO: validate list name
	rsp, err := c.dataClient.ListPushBack(ctx, &models.ListPushBackRequest{
		CacheName:           request.CacheName,
		ListName:            request.ListName,
		Value:               request.Value,
		TruncateFrontToSize: request.TruncateFrontToSize,
		CollectionTtl:       request.CollectionTTL,
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	return convertListPushBackResponse(rsp)
}

func (c *DefaultScsClient) SortedSetPut(ctx context.Context, request *SortedSetPutRequest) error {
	setName, err := isSetNameValid([]byte(request.SetName))
	if err != nil {
		return convertMomentoSvcErrorToCustomerError(err)
	}

	err = c.dataClient.SortedSetPut(ctx, &models.SortedSetPutRequest{
		CacheName:     request.CacheName,
		SetName:       setName,
		Elements:      convertSortedSetScoreRequestElement(request.Elements),
		CollectionTTL: request.CollectionTTL,
	})
	if err != nil {
		return convertMomentoSvcErrorToCustomerError(err)
	}
	return nil
}

func (c *DefaultScsClient) SortedSetFetch(ctx context.Context, request *SortedSetFetchRequest) (SortedSetFetchResponse, error) {
	setName, err := isSetNameValid([]byte(request.SetName))
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}

	rsp, err := c.dataClient.SortedSetFetch(ctx, &models.SortedSetFetchRequest{
		CacheName:       request.CacheName,
		SetName:         setName,
		Order:           models.SortedSetOrder(request.Order),
		NumberOfResults: convertSortedSetFetchNumResultsRequest(request.NumberOfResults),
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	switch r := rsp.(type) {
	case *models.SortedSetFetchMissing:
		return SortedSetFetchMiss{}, nil
	case *models.SortedSetFetchFound:
		return &SortedSetFetchHit{
			Elements: convertInternalSortedSetElement(r.Elements),
		}, nil
	default:
		return nil, momentoerrors.NewMomentoSvcErr(
			momento.ClientSdkError,
			fmt.Sprintf("unexpected sortedset fetch status returned %+v", r),
			nil,
		)
	}
}

func (c *DefaultScsClient) SortedSetGetScore(ctx context.Context, request *SortedSetGetScoreRequest) (SortedSetGetScoreResponse, error) {
	setName, err := isSetNameValid([]byte(request.SetName))
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}

	rsp, err := c.dataClient.SortedSetGetScore(ctx, &models.SortedSetGetScoreRequest{
		CacheName:    request.CacheName,
		SetName:      setName,
		ElementNames: momentoBytesListToPrimitiveByteList(request.ElementNames),
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	switch r := rsp.(type) {
	case *models.SortedSetGetScoreMiss:
		return &SortedSetGetScoreMiss{}, nil
	case *models.SortedSetGetScoreHit:
		return &SortedSetGetScoreHit{
			Elements: convertSortedSetScoreElement(r.Elements),
		}, nil
	default:
		return nil, momentoerrors.NewMomentoSvcErr(
			momento.ClientSdkError,
			fmt.Sprintf("unexpected sortedset getscore status returned %+v", r),
			nil,
		)
	}

}
func (c *DefaultScsClient) SortedSetRemove(ctx context.Context, request *SortedSetRemoveRequest) error {
	setName, err := isSetNameValid([]byte(request.SetName))
	if err != nil {
		return convertMomentoSvcErrorToCustomerError(err)
	}

	err = c.dataClient.SortedSetRemove(ctx, &models.SortedSetRemoveRequest{
		CacheName:        request.CacheName,
		SetName:          setName,
		ElementsToRemove: convertSortedSetRemoveNumItemsRequest(request.ElementsToRemove),
	})
	if err != nil {
		return convertMomentoSvcErrorToCustomerError(err)
	}
	return nil
}

func (c *DefaultScsClient) SortedSetGetRank(ctx context.Context, request *SortedSetGetRankRequest) (SortedSetGetRankResponse, error) {
	setName, err := isSetNameValid([]byte(request.SetName))
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	// TODO validate element name

	rsp, err := c.dataClient.SortedSetGetRank(ctx, &models.SortedSetGetRankRequest{
		CacheName:   request.CacheName,
		SetName:     setName,
		ElementName: request.ElementName.AsBytes(),
	})
	if err != nil {
		return nil, convertMomentoSvcErrorToCustomerError(err)
	}
	switch r := rsp.(type) {
	case *models.SortedSetGetRankMiss:
		return &SortedSetGetRankMiss{}, nil
	case *models.SortedSetGetRankHit:
		if r.Status == models.Hit {
			return &SortedSetGetRankHit{
				Element: &SortedSetRankHit{
					Rank: r.Rank,
				},
			}, nil
		} else if r.Status == models.Miss {
			return &SortedSetGetRankHit{
				Element: &SortedSetRankMiss{},
			}, nil
		} else {
			return nil, momentoerrors.NewMomentoSvcErr(
				momento.ClientSdkError,
				fmt.Sprintf("unexpected sortedset getRank status returned %+v", r.Status),
				nil,
			)
		}
	default:
		return nil, momentoerrors.NewMomentoSvcErr(
			momento.ClientSdkError,
			fmt.Sprintf("unexpected sortedset getscore status returned %+v", r),
			nil,
		)
	}
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

func convertListPushFrontResponse(r models.ListPushFrontResponse) (ListPushFrontResponse, momento.MomentoError) {
	switch response := r.(type) {
	case *models.ListPushFrontSuccess:
		return &ListPushFrontSuccess{value: response.Value}, nil
	default:
		return nil, convertMomentoSvcErrorToCustomerError(momentoerrors.NewMomentoSvcErr(
			momento.ClientSdkError,
			fmt.Sprintf("unexpected list push front status returned %+v", response),
			nil,
		))
	}
}

func convertListPushBackResponse(r models.ListPushBackResponse) (ListPushBackResponse, momento.MomentoError) {
	switch response := r.(type) {
	case *models.ListPushBackSuccess:
		return &ListPushBackSuccess{value: response.Value}, nil
	default:
		return nil, momentoerrors.NewMomentoSvcErr(
			momento.ClientSdkError,
			fmt.Sprintf("unexpected list push back status returned %+v", response),
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

func convertSortedSetFetchNumResultsRequest(results SortedSetFetchNumResults) models.SortedSetFetchNumResults {
	switch r := results.(type) {
	case FetchLimitedItems:
		return &models.FetchLimitedItems{Limit: r.Limit}
	default:
		return &models.FetchAllItems{}
	}
}
func isSetNameValid(key []byte) ([]byte, momentoerrors.MomentoSvcErr) {
	if len(key) == 0 {
		return key, momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, "key cannot be empty", nil)
	}
	return key, nil
}

func convertInternalSortedSetElement(e []*models.SortedSetElement) []*SortedSetElement {
	var rList []*SortedSetElement
	for _, el := range e {
		rList = append(rList, &SortedSetElement{
			Name:  el.Name,
			Score: el.Score,
		})
	}
	return rList
}

func convertSortedSetScoreRequestElement(e []*SortedSetScoreRequestElement) []*models.SortedSetElement {
	var rList []*models.SortedSetElement
	for _, el := range e {
		rList = append(rList, &models.SortedSetElement{
			Name:  el.Name.AsBytes(),
			Score: el.Score,
		})
	}
	return rList
}

func convertSortedSetScoreElement(e []*models.SortedSetScore) []SortedSetScoreElement {
	var rList []SortedSetScoreElement
	for _, el := range e {
		if el.Result == models.Hit {
			rList = append(rList, &SortedSetScoreHit{
				Score: el.Score,
			})
		} else if el.Result == models.Miss {
			rList = append(rList, &SortedSetScoreMiss{})
		} else {
			rList = append(rList, &SortedSetScoreInvalid{})
		}
	}
	return rList
}

func convertSortedSetRemoveNumItemsRequest(results SortedSetRemoveNumItems) models.SortedSetRemoveNumItems {
	switch r := results.(type) {
	case RemoveSomeItems:
		return &models.RemoveSomeItems{
			ElementsToRemove: momentoBytesListToPrimitiveByteList(r.elementsToRemove),
		}
	default:
		return &models.RemoveAllItems{}
	}
}

func momentoBytesListToPrimitiveByteList(i []momento.Bytes) [][]byte {
	var rList [][]byte
	for _, mb := range i {
		rList = append(rList, mb.AsBytes())
	}
	return rList
}
