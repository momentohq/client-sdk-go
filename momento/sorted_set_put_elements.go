package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/utils"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetPutElementsRequest struct {
	CacheName string
	SetName   string
	Elements  []SortedSetElement
	Ttl       *utils.CollectionTtl

	grpcRequest  *pb.XSortedSetPutRequest
	grpcResponse *pb.XSortedSetPutResponse
	response     responses.SortedSetPutElementsResponse
}

func (r *SortedSetPutElementsRequest) cacheName() string { return r.CacheName }

func (r *SortedSetPutElementsRequest) requestName() string { return "Sorted set put elements" }

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

	elements, err := convertSortedSetElementsToGrpc(r.Elements)
	if err != nil {
		return err
	}

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
	r.response = &responses.SortedSetPutElementsSuccess{}
	return nil
}

func convertSortedSetElementsToGrpc(modelSetElements []SortedSetElement) ([]*pb.XSortedSetElement, error) {
	if modelSetElements == nil {
		return nil, buildError(
			momentoerrors.InvalidArgumentError, "elements cannot be nil", nil,
		)
	}
	var returnList []*pb.XSortedSetElement
	for _, el := range modelSetElements {
		if el.Value == nil {
			return nil, buildError(
				momentoerrors.InvalidArgumentError, "element value cannot be nil", nil,
			)
		}
		returnList = append(returnList, &pb.XSortedSetElement{
			Value: el.Value.asBytes(),
			Score: el.Score,
		})
	}
	return returnList, nil
}

func (r *SortedSetPutElementsRequest) getResponse() interface{} {
	return r.response
}
