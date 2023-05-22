package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/utils"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type IncrementRequest struct {
	CacheName string
	Field     Field
	Amount    int64
	Ttl       *utils.CollectionTtl

	grpcRequest  *pb.XIncrementRequest
	grpcResponse *pb.XIncrementResponse
	response     responses.IncrementResponse
}

func (r *IncrementRequest) cacheName() string { return r.CacheName }

func (r *IncrementRequest) field() Field { return r.Field }

func (r *IncrementRequest) ttl() time.Duration { return r.Ttl.Ttl }

func (r *IncrementRequest) collectionTtl() *utils.CollectionTtl { return r.Ttl }

func (r *IncrementRequest) requestName() string { return "Increment" }

func (r *IncrementRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	var field []byte
	if field, err = prepareField(r); err != nil {
		return err
	}

	var ttlMilliseconds uint64
	if ttlMilliseconds, _, err = prepareCollectionTtl(r, client.defaultTtl); err != nil {
		return err
	}

	r.grpcRequest = &pb.XIncrementRequest{
		CacheKey:        field,
		Amount:          r.Amount,
		TtlMilliseconds: ttlMilliseconds,
	}

	return nil
}

func (r *IncrementRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.Increment(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *IncrementRequest) interpretGrpcResponse() error {
	r.response = responses.NewIncrementSuccess(r.grpcResponse.Value)
	return nil
}
