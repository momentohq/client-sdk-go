package momento

import (
	"context"
	"time"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/utils"
)

///////// Response

type SortedSetPutResponse interface {
	isSortedSetPutResponse()
}

type SortedSetPutSuccess struct{}

func (SortedSetPutSuccess) isSortedSetPutResponse() {}

///////// Request

type SortedSetPutElement struct {
	Value Value
	Score float64
}

type SortedSetPutRequest struct {
	CacheName string
	SetName   string
	Elements  []*SortedSetPutElement
	Ttl       *utils.CollectionTtl

	grpcRequest  *pb.XSortedSetPutRequest
	grpcResponse *pb.XSortedSetPutResponse
	response     SortedSetPutResponse
}

func (r *SortedSetPutRequest) cacheName() string { return r.CacheName }

func (r *SortedSetPutRequest) requestName() string { return "Sorted set put" }

func (r *SortedSetPutRequest) ttl() *time.Duration { return r.Ttl.Ttl }

func (r *SortedSetPutRequest) refreshTtl() *bool { return r.Ttl.RefreshTtl }

func (r *SortedSetPutRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	var ttlMills uint64
	if ttlMills, err = prepareTTL(r, client.defaultTtl); err != nil {
		return err
	}

	elements := convertSortedSetElementToGrpc(r.Elements)

	r.grpcRequest = &pb.XSortedSetPutRequest{
		SetName:         []byte(r.SetName),
		Elements:        elements,
		TtlMilliseconds: ttlMills,
		RefreshTtl:      *prepareRefreshTtl(r),
	}
	return nil
}

func (r *SortedSetPutRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetPut(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SortedSetPutRequest) interpretGrpcResponse() error {
	r.response = &SortedSetPutSuccess{}
	return nil
}

func convertSortedSetElementToGrpc(modelSetElements []*SortedSetPutElement) []*pb.XSortedSetElement {
	var returnList []*pb.XSortedSetElement
	for _, el := range modelSetElements {
		returnList = append(returnList, &pb.XSortedSetElement{
			Name:  el.Value.asBytes(),
			Score: el.Score,
		})
	}
	return returnList
}
