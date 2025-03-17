package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetIfEqualRequest struct {
	// Name of the cache to store the item in.
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// string ot byte value to be stored.
	Value Value
	// string or byte value to compare with the existing value in the cache.
	Equal Value
	// Optional Time to live in cache in seconds.
	// If not provided, then default TTL for the cache client instance is used.
	Ttl time.Duration

	response responses.SetIfEqualResponse
}

func (r *SetIfEqualRequest) cacheName() string { return r.CacheName }

func (r *SetIfEqualRequest) key() Key { return r.Key }

func (r *SetIfEqualRequest) value() Value { return r.Value }

func (r *SetIfEqualRequest) equal() Value { return r.Equal }

func (r *SetIfEqualRequest) ttl() time.Duration { return r.Ttl }

func (r *SetIfEqualRequest) requestName() string { return "SetIfEqual" }

func (r *SetIfEqualRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return nil, err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return nil, err
	}

	var equal []byte
	if equal, err = prepareEqual(r); err != nil {
		return nil, err
	}

	var ttl uint64
	if ttl, err = prepareTtl(r, client.defaultTtl); err != nil {
		return nil, err
	}

	var condition = &pb.XSetIfRequest_Equal{
		Equal: &pb.Equal{
			ValueToCheck: equal,
		},
	}
	grpcRequest := &pb.XSetIfRequest{
		CacheKey:        key,
		CacheBody:       value,
		TtlMilliseconds: ttl,
		Condition:       condition,
	}

	return grpcRequest, nil
}

func (r *SetIfEqualRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetIf(requestMetadata, grpcRequest.(*pb.XSetIfRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SetIfEqualRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XSetIfResponse)
	switch myResp.Result.(type) {
	case *pb.XSetIfResponse_Stored:
		r.response = &responses.SetIfEqualStored{}
	case *pb.XSetIfResponse_NotStored:
		r.response = &responses.SetIfEqualNotStored{}
	default:
		return errUnexpectedGrpcResponse(r, myResp)
	}
	return nil
}

func (r *SetIfEqualRequest) validateResponseType(resp grpcResponse) error {
	_, ok := resp.(*pb.XSetIfResponse)
	if !ok {
		return errUnexpectedGrpcResponse(nil, resp)
	}
	return nil
}
