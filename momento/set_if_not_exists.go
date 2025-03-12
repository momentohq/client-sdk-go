package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetIfNotExistsRequest struct {
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

	response     responses.SetIfNotExistsResponse
}

func (r *SetIfNotExistsRequest) cacheName() string { return r.CacheName }

func (r *SetIfNotExistsRequest) key() Key { return r.Key }

func (r *SetIfNotExistsRequest) value() Value { return r.Value }

func (r *SetIfNotExistsRequest) ttl() time.Duration { return r.Ttl }

func (r *SetIfNotExistsRequest) requestName() string { return "SetIfNotExists" }

func (r *SetIfNotExistsRequest) initGrpcRequest(client scsDataClient) error {
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

	condition := &pb.XSetIfRequest_Absent{
		Absent: &pb.Absent{},
	}
	r.grpcRequest = &pb.XSetIfRequest{
		CacheKey:        key,
		CacheBody:       value,
		TtlMilliseconds: ttl,
		Condition:       condition,
	}

	return nil
}

func (r *SetIfNotExistsRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.SetIf(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *SetIfNotExistsRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XSetIfResponse)

	switch myResp.Result.(type) {
	case *pb.XSetIfResponse_Stored:
		r.response = &responses.SetIfNotExistsStored{}
	case *pb.XSetIfResponse_NotStored:
		r.response = &responses.SetIfNotExistsNotStored{}
	default:
		return errUnexpectedGrpcResponse(r, myResp)
	}
	return nil
}

func (r *SetIfNotExistsRequest) getResponse() interface{} {
	return r.response
}
