package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

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

func (r *IncrementRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.Increment(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	r.grpcResponse = resp
	return resp, nil, nil
}

func (r *IncrementRequest) interpretGrpcResponse(_ interface{}) error {
	r.response = responses.NewIncrementSuccess(r.grpcResponse.Value)
	return nil
}

func (r *IncrementRequest) getResponse() interface{} { return r.response }
