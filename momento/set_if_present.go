package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetIfPresentRequest struct {
	// Name of the cache to store the item in.
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// string ot byte value to be stored.
	Value Value
	// Optional Time to live in cache in seconds.
	// If not provided, then default TTL for the cache client instance is used.
	Ttl time.Duration
}

func (r *SetIfPresentRequest) cacheName() string { return r.CacheName }

func (r *SetIfPresentRequest) key() Key { return r.Key }

func (r *SetIfPresentRequest) value() Value { return r.Value }

func (r *SetIfPresentRequest) ttl() time.Duration { return r.Ttl }

func (r *SetIfPresentRequest) requestName() string { return "SetIfPresent" }

func (r *SetIfPresentRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
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

	grpcRequest := &pb.XSetIfRequest{
		CacheKey:        key,
		CacheBody:       value,
		TtlMilliseconds: ttl,
		Condition:       &pb.XSetIfRequest_Present{},
	}
	return grpcRequest, nil
}

func (r *SetIfPresentRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetIf(requestMetadata, grpcRequest.(*pb.XSetIfRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SetIfPresentRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XSetIfResponse)
	switch myResp.Result.(type) {
	case *pb.XSetIfResponse_Stored:
		return &responses.SetIfPresentStored{}, nil
	case *pb.XSetIfResponse_NotStored:
		return &responses.SetIfPresentNotStored{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}
