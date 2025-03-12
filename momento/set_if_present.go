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

	grpcRequest  *pb.XSetIfRequest

	response     responses.SetIfPresentResponse
}

func (r *SetIfPresentRequest) cacheName() string { return r.CacheName }

func (r *SetIfPresentRequest) key() Key { return r.Key }

func (r *SetIfPresentRequest) value() Value { return r.Value }

func (r *SetIfPresentRequest) ttl() time.Duration { return r.Ttl }

func (r *SetIfPresentRequest) requestName() string { return "SetIfNotExists" }

func (r *SetIfPresentRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return err
	}

	var ttl uint64
	if ttl, err = prepareTtl(r, client.defaultTtl); err != nil {
		return err
	}

	r.grpcRequest = &pb.XSetIfRequest{
		CacheKey:        key,
		CacheBody:       value,
		TtlMilliseconds: ttl,
		Condition:       &pb.XSetIfRequest_Present{},
	}

	return nil
}

func (r *SetIfPresentRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetIf(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SetIfPresentRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XSetIfResponse)

	switch myResp.Result.(type) {
	case *pb.XSetIfResponse_Stored:
		r.response = &responses.SetIfPresentStored{}
	case *pb.XSetIfResponse_NotStored:
		r.response = &responses.SetIfPresentNotStored{}
	default:
		return errUnexpectedGrpcResponse(r, myResp)
	}
	return nil
}

func (r *SetIfPresentRequest) getResponse() interface{} {
	return r.response
}
