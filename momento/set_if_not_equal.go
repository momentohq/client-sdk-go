package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetIfNotEqualRequest struct {
	// Name of the cache to store the item in.
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// string ot byte value to be stored.
	Value Value
	// string or byte value to compare with the existing value in the cache.
	NotEqual Value
	// Optional Time to live in cache in seconds.
	// If not provided, then default TTL for the cache client instance is used.
	Ttl time.Duration
}

func (r *SetIfNotEqualRequest) cacheName() string { return r.CacheName }

func (r *SetIfNotEqualRequest) key() Key { return r.Key }

func (r *SetIfNotEqualRequest) value() Value { return r.Value }

func (r *SetIfNotEqualRequest) notEqual() Value { return r.NotEqual }

func (r *SetIfNotEqualRequest) ttl() time.Duration { return r.Ttl }

func (r *SetIfNotEqualRequest) requestName() string { return "SetIfNotEqual" }

func (r *SetIfNotEqualRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return nil, err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return nil, err
	}

	var notEqual []byte
	if notEqual, err = prepareNotEqual(r); err != nil {
		return nil, err
	}

	var ttl uint64
	if ttl, err = prepareTtl(r, client.defaultTtl); err != nil {
		return nil, err
	}

	condition := &pb.XSetIfRequest_NotEqual{
		NotEqual: &pb.NotEqual{
			ValueToCheck: notEqual,
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

func (r *SetIfNotEqualRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetIf(requestMetadata, grpcRequest.(*pb.XSetIfRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SetIfNotEqualRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XSetIfResponse)

	switch myResp.Result.(type) {
	case *pb.XSetIfResponse_Stored:
		return &responses.SetIfNotEqualStored{}, nil
	case *pb.XSetIfResponse_NotStored:
		return &responses.SetIfNotEqualNotStored{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}

func (c SetIfNotEqualRequest) GetRequestName() string {
	return "SetIfNotEqualRequest"
}
