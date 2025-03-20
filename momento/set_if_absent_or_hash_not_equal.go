package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetIfAbsentOrHashNotEqualRequest struct {
	// Name of the cache to store the item in.
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// string ot byte value to be stored.
	Value Value
	// string or byte value to compare with the existing value in the cache.
	HashNotEqual Value
	// Optional Time to live in cache in seconds.
	// If not provided, then default TTL for the cache client instance is used.
	Ttl time.Duration

	grpcRequest  *pb.XSetIfHashRequest
	grpcResponse *pb.XSetIfHashResponse
	response     responses.SetIfAbsentOrHashNotEqualResponse
}

func (r *SetIfAbsentOrHashNotEqualRequest) cacheName() string { return r.CacheName }

func (r *SetIfAbsentOrHashNotEqualRequest) key() Key { return r.Key }

func (r *SetIfAbsentOrHashNotEqualRequest) value() Value { return r.Value }

func (r *SetIfAbsentOrHashNotEqualRequest) notEqual() Value { return r.HashNotEqual }

func (r *SetIfAbsentOrHashNotEqualRequest) ttl() time.Duration { return r.Ttl }

func (r *SetIfAbsentOrHashNotEqualRequest) requestName() string { return "SetIfAbsentOrHashNotEqual" }

func (r *SetIfAbsentOrHashNotEqualRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return err
	}

	var hashNotEqual []byte
	if hashNotEqual, err = prepareNotEqual(r); err != nil {
		return err
	}

	var ttl uint64
	if ttl, err = prepareTtl(r, client.defaultTtl); err != nil {
		return err
	}

	var condition = &pb.XSetIfHashRequest_AbsentOrNotHashEqual{
		AbsentOrEqual: &pb.AbsentOrNotHashEqual{
			HashToCheck: hashNotEqual,
		},
	}
	r.grpcRequest = &pb.XSetIfHashRequest{
		CacheKey:        key,
		CacheBody:       value,
		TtlMilliseconds: ttl,
		Condition:       condition,
	}

	return nil
}

func (r *SetIfAbsentOrHashNotEqualRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetIf(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	r.grpcResponse = resp
	return resp, nil, nil
}

func (r *SetIfAbsentOrHashNotEqualRequest) interpretGrpcResponse() error {
	grpcResp := r.grpcResponse
	var resp responses.SetIfAbsentOrHashNotEqualResponse

	switch grpcResp.Result.(type) {
	case *pb.XSetIfResponse_Stored:
		resp = &responses.SetIfAbsentOrHashNotEqualStored{}
	case *pb.XSetIfResponse_NotStored:
		resp = &responses.SetIfAbsentOrHashNotEqualNotStored{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}

	r.response = resp
	return nil
}
