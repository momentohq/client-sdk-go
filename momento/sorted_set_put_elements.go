package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/utils"
)

type SortedSetPutElement struct {
	Value Value
	Score float64
}

type SortedSetPutElementsRequest struct {
	CacheName string
	SetName   string
	Elements  []*SortedSetPutElement
	Ttl       *utils.CollectionTtl

	grpcRequest  *pb.XSortedSetPutRequest
	grpcResponse *pb.XSortedSetPutResponse
	response     responses.SortedSetPutElementsResponse
}

func (r *SortedSetPutElementsRequest) cacheName() string { return r.CacheName }

func (r *SortedSetPutElementsRequest) requestName() string { return "Sorted set put" }

func (r *SortedSetPutElementsRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *SortedSetPutElementsRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *SortedSetPutElementsRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return err
	}

	elements := convertSortedSetElementToGrpc(r.Elements)

	r.grpcRequest = &pb.XSortedSetPutRequest{
		SetName:         []byte(r.SetName),
		Elements:        elements,
		TtlMilliseconds: ttlMilliseconds,
		RefreshTtl:      refreshTtl,
	}
	return nil
}

func (r *SortedSetPutElementsRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetPut(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SortedSetPutElementsRequest) interpretGrpcResponse() error {
	r.response = &responses.SortedSetPutSuccess{}
	return nil
}

func convertSortedSetElementToGrpc(modelSetElements []*SortedSetPutElement) []*pb.XSortedSetElement {
	var returnList []*pb.XSortedSetElement
	for _, el := range modelSetElements {
		returnList = append(returnList, &pb.XSortedSetElement{
			Value: el.Value.asBytes(),
			Score: el.Score,
		})
	}
	return returnList
}
