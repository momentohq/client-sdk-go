package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/utils"
)

//////// Response

type SortedSetIncrementResponse interface {
	isSortedSetIncrementResponse()
}
type SortedSetIncrementSuccess struct {
	Value float64
}

func (SortedSetIncrementSuccess) isSortedSetIncrementResponse() {}

////// Request

type SortedSetIncrementRequest struct {
	CacheName     string
	SetName       string
	ElementName   Bytes
	Amount        float64
	CollectionTTL utils.CollectionTTL

	grpcRequest  *pb.XSortedSetIncrementRequest
	grpcResponse *pb.XSortedSetIncrementResponse
	response     SortedSetIncrementResponse
}

func (r SortedSetIncrementRequest) cacheName() string { return r.CacheName }

func (r SortedSetIncrementRequest) requestName() string { return "Sorted set increment" }

func (r *SortedSetIncrementRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	ttlMills, refreshTTL := prepareCollectionTtl(r.CollectionTTL, client.defaultTtl)

	r.grpcRequest = &pb.XSortedSetIncrementRequest{
		SetName:         []byte(r.SetName),
		ElementName:     r.ElementName.AsBytes(),
		Amount:          r.Amount,
		TtlMilliseconds: ttlMills,
		RefreshTtl:      refreshTTL,
	}
	return nil
}

func (r *SortedSetIncrementRequest) makeGrpcRequest(client scsDataClient, metadata context.Context) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetIncrement(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SortedSetIncrementRequest) interpretGrpcResponse() error {
	r.response = SortedSetIncrementSuccess{}
	return nil
}
