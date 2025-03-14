package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetRequest struct {
	// Name of the cache to store the item in.
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// string ot byte value to be stored.
	Value Value
	// Optional Time to live in cache in seconds.
	// If not provided, then default TTL for the cache client instance is used.
	Ttl time.Duration

	grpcRequest *pb.XSetRequest

	response responses.SetResponse
}

func (r *SetRequest) cacheName() string { return r.CacheName }

func (r *SetRequest) key() Key { return r.Key }

func (r *SetRequest) value() Value { return r.Value }

func (r *SetRequest) ttl() time.Duration { return r.Ttl }

func (r *SetRequest) requestName() string { return "Set" }

func (r *SetRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return nil, err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return nil, err
	}

	var ttl uint64
	if ttl, err = prepareTtl(r, client.defaultTtl); err != nil {
		return nil, err
	}

	r.grpcRequest = &pb.XSetRequest{
		CacheKey:        key,
		CacheBody:       value,
		TtlMilliseconds: ttl,
	}

	return r.grpcRequest, nil
}

func (r *SetRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.Set(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	//r.grpcResponse = resp
	return resp, nil, nil
}

func (r *SetRequest) interpretGrpcResponse(_ interface{}) error {
	r.response = &responses.SetSuccess{}
	return nil
}

func (r *SetRequest) validateResponseType(resp grpcResponse) error {
	_, ok := resp.(*pb.XSetResponse)
	if !ok {
		return errUnexpectedGrpcResponse(nil, resp)
	}
	return nil
}
