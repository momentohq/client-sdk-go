package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"

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

	grpcRequest  *pb.XSetIfNotExistsRequest
	grpcResponse *pb.XSetIfNotExistsResponse
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

	r.grpcRequest = &pb.XSetIfNotExistsRequest{
		CacheKey:        key,
		CacheBody:       value,
		TtlMilliseconds: ttl,
	}

	return nil
}

func (r *SetIfNotExistsRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SetIfNotExists(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SetIfNotExistsRequest) interpretGrpcResponse() error {
	grpcResp := r.grpcResponse
	var resp responses.SetIfNotExistsResponse

	switch grpcResp.Result.(type) {
	case *pb.XSetIfNotExistsResponse_Stored:
		resp = &responses.SetIfNotExistsStored{}
	case *pb.XSetIfNotExistsResponse_NotStored:
		resp = &responses.SetIfNotExistsNotStored{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}

	r.response = resp
	return nil
}
