package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/utils"
)

type SortedSetPutElementRequest struct {
	CacheName string
	SetName   string
	Value     Value
	Score     float64
	Ttl       *utils.CollectionTtl

	grpcRequest  *pb.XSortedSetPutRequest
	grpcResponse *pb.XSortedSetPutResponse
	response     responses.SortedSetPutElementResponse
}

func (r *SortedSetPutElementRequest) cacheName() string { return r.CacheName }

func (r *SortedSetPutElementRequest) requestName() string { return "Sorted set put" }

func (r *SortedSetPutElementRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *SortedSetPutElementRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *SortedSetPutElementRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	var ttlMilliseconds uint64
	var refreshTtl bool
	if ttlMilliseconds, refreshTtl, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return err
	}

	elements := convertSortedSetElementToGrpc(r.Value, r.Score)

	r.grpcRequest = &pb.XSortedSetPutRequest{
		SetName:         []byte(r.SetName),
		Elements:        elements,
		TtlMilliseconds: ttlMilliseconds,
		RefreshTtl:      refreshTtl,
	}
	return nil
}

func (r *SortedSetPutElementRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetPut(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SortedSetPutElementRequest) interpretGrpcResponse() error {
	r.response = &responses.SortedSetPutElementSuccess{}
	return nil
}

func convertSortedSetElementToGrpc(value Value, score float64) []*pb.XSortedSetElement {
	var returnList []*pb.XSortedSetElement
	return append(returnList, &pb.XSortedSetElement{
		Value: value.asBytes(),
		Score: score,
	})
}
