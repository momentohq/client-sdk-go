package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetIfPresentAndHashNotEqualRequest struct {
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
}

func (r *SetIfPresentAndHashNotEqualRequest) cacheName() string { return r.CacheName }

func (r *SetIfPresentAndHashNotEqualRequest) key() Key { return r.Key }

func (r *SetIfPresentAndHashNotEqualRequest) value() Value { return r.Value }

func (r *SetIfPresentAndHashNotEqualRequest) notEqual() Value { return r.HashNotEqual }

func (r *SetIfPresentAndHashNotEqualRequest) ttl() time.Duration { return r.Ttl }

func (r *SetIfPresentAndHashNotEqualRequest) requestName() string {
	return "SetIfPresentAndHashNotEqual"
}

func (r *SetIfPresentAndHashNotEqualRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return nil, err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return nil, err
	}

	var hashNotEqual []byte
	if hashNotEqual, err = prepareNotEqual(r); err != nil {
		return nil, err
	}

	var ttl uint64
	if ttl, err = prepareTtl(r, client.defaultTtl); err != nil {
		return nil, err
	}

	var condition = &pb.XSetIfHashRequest_PresentAndNotHashEqual{
		PresentAndNotHashEqual: &pb.PresentAndNotHashEqual{
			HashToCheck: hashNotEqual,
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

func (r *SetIfPresentAndHashNotEqualRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
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

func (r *SetIfPresentAndHashNotEqualRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	grpcResponse := resp.(*pb.XSetIfHashResponse)
	switch response := grpcResponse.Result.(type) {
	case *pb.XSetIfHashResponse_Stored:
		return responses.NewSetIfPresentAndHashNotEqualStored(response.Stored.NewHash), nil
	case *pb.XSetIfHashResponse_NotStored:
		return &responses.SetIfPresentAndHashNotEqualNotStored{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, grpcResponse)
	}
}
