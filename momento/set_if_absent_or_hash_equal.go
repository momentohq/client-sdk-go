package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetIfAbsentOrHashEqualRequest struct {
	// Name of the cache to store the item in.
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// string ot byte value to be stored.
	Value Value
	// string or byte value to compare with the existing value in the cache.
	HashEqual Value
	// Optional Time to live in cache in seconds.
	// If not provided, then default TTL for the cache client instance is used.
	Ttl time.Duration
}

func (r *SetIfAbsentOrHashEqualRequest) cacheName() string { return r.CacheName }

func (r *SetIfAbsentOrHashEqualRequest) key() Key { return r.Key }

func (r *SetIfAbsentOrHashEqualRequest) value() Value { return r.Value }

func (r *SetIfAbsentOrHashEqualRequest) equal() Value { return r.HashEqual }

func (r *SetIfAbsentOrHashEqualRequest) ttl() time.Duration { return r.Ttl }

func (r *SetIfAbsentOrHashEqualRequest) requestName() string { return "SetIfAbsentOrHashEqual" }

func (r *SetIfAbsentOrHashEqualRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return nil, err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return nil, err
	}

	var hashEqual []byte
	if hashEqual, err = prepareEqual(r); err != nil {
		return nil, err
	}

	var ttl uint64
	if ttl, err = prepareTtl(r, client.defaultTtl); err != nil {
		return nil, err
	}

	var condition = &pb.XSetIfHashRequest_AbsentOrHashEqual{
		AbsentOrHashEqual: &pb.AbsentOrHashEqual{
			HashToCheck: hashEqual,
		},
	}
	grpcRequest := &pb.XSetIfHashRequest{
		CacheKey:        key,
		CacheBody:       value,
		TtlMilliseconds: ttl,
		Condition:       condition,
	}

	return grpcRequest, nil
}

func (r *SetIfAbsentOrHashEqualRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetIfHash(
		requestMetadata,
		grpcRequest.(*pb.XSetIfHashRequest),
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SetIfAbsentOrHashEqualRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	grpcResponse := resp.(*pb.XSetIfHashResponse)
	switch response := grpcResponse.Result.(type) {
	case *pb.XSetIfHashResponse_Stored:
		return responses.NewSetIfAbsentOrHashEqualStored(response.Stored.NewHash), nil
	case *pb.XSetIfHashResponse_NotStored:
		return &responses.SetIfAbsentOrHashEqualNotStored{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, grpcResponse)
	}
}
