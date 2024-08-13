package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetRequest struct {
	// Name of the cache to store the item in.
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// string ot byte value to be stored.
	Value Value
	// Optional Time to live in cache in seconds.
	// If not provided, then default TTL for the cache client instance is used.
	Ttl time.Duration

	grpcRequest  *pb.XSetRequest
	grpcResponse *pb.XSetResponse
	response     responses.SetResponse
}

func (r *SetRequest) cacheName() string { return r.CacheName }

func (r *SetRequest) key() Key { return r.Key }

func (r *SetRequest) value() Value { return r.Value }

func (r *SetRequest) ttl() time.Duration { return r.Ttl }

func (r *SetRequest) requestName() string { return "Set" }

func (r *SetRequest) initGrpcRequest(client scsDataClient) error {
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
	// If ttl is not provided, then default TTL for the cache client instance is used.
	// Make sure that the updated ttl is represented in the request.
	r.Ttl = time.Duration(ttl)
	r.grpcRequest = &pb.XSetRequest{
		CacheKey:        key,
		CacheBody:       value,
		TtlMilliseconds: ttl,
	}

	return nil
}

func (r *SetRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.Set(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SetRequest) interpretGrpcResponse() error {
	r.response = &responses.SetSuccess{}
	return nil
}

func (r *SetRequest) getResponse() interface{} {
	return r.response
}
