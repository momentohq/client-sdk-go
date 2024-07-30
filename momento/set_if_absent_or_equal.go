package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetIfAbsentOrEqualRequest struct {
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

	grpcRequest  *pb.XSetIfRequest
	grpcResponse *pb.XSetIfResponse
	response     responses.SetIfAbsentOrEqualResponse
}

func (r *SetIfAbsentOrEqualRequest) cacheName() string { return r.CacheName }

func (r *SetIfAbsentOrEqualRequest) key() Key { return r.Key }

func (r *SetIfAbsentOrEqualRequest) value() Value { return r.Value }

func (r *SetIfAbsentOrEqualRequest) equal() Value { return r.Equal }

func (r *SetIfAbsentOrEqualRequest) ttl() time.Duration { return r.Ttl }

func (r *SetIfAbsentOrEqualRequest) requestName() string { return "SetIfNotExists" }

func (r *SetIfAbsentOrEqualRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return err
	}

	var equal []byte
	if equal, err = prepareEqual(r); err != nil {
		return err
	}

	var ttl uint64
	if ttl, err = prepareTtl(r, client.defaultTtl); err != nil {
		return err
	}

	var condition = &pb.XSetIfRequest_AbsentOrEqual{
		AbsentOrEqual: &pb.AbsentOrEqual{
			ValueToCheck: equal,
		},
	}
	r.grpcRequest = &pb.XSetIfRequest{
		CacheKey:        key,
		CacheBody:       value,
		TtlMilliseconds: ttl,
		Condition:       condition,
	}

	return nil
}

func (r *SetIfAbsentOrEqualRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetIf(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	r.grpcResponse = resp
	return resp, nil, nil
}

func (r *SetIfAbsentOrEqualRequest) interpretGrpcResponse() error {
	grpcResp := r.grpcResponse
	var resp responses.SetIfAbsentOrEqualResponse

	switch grpcResp.Result.(type) {
	case *pb.XSetIfResponse_Stored:
		resp = &responses.SetIfAbsentOrEqualStored{}
	case *pb.XSetIfResponse_NotStored:
		resp = &responses.SetIfAbsentOrEqualNotStored{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}

	r.response = resp
	return nil
}

func (r *SetIfAbsentOrEqualRequest) getResponse() map[string]string {
	return getMomentoResponseData(r.response)
}
