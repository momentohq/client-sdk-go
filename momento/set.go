package momento

import (
	"context"
	"time"

	client_sdk_go "github.com/momentohq/client-sdk-go/internal/protos"
)

//////////// SetResponse /////////////

type SetResponse interface {
	isSetResponse()
}

type SetSuccess struct{}

func (SetSuccess) isSetResponse() {}

///////////// SetRequest /////////////

type SetRequest struct {
	// Name of the cache to store the item in.
	CacheName string
	// string or byte key to be used to store item.
	Key Bytes
	// string ot byte value to be stored.
	Value Bytes
	// Optional Time to live in cache in seconds.
	// If not provided, then default TTL for the cache client instance is used.
	TTL time.Duration

	grpcRequest  *client_sdk_go.XSetRequest
	grpcResponse *client_sdk_go.XSetResponse
	response     SetResponse
}

func (r SetRequest) cacheName() string { return r.CacheName }

func (r SetRequest) key() Bytes { return r.Key }

func (r SetRequest) value() Bytes { return r.Value }

func (r SetRequest) ttl() time.Duration { return r.TTL }

func (r SetRequest) requestName() string { return "Set" }

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
	if ttl, err = prepareTTL(r, client.defaultTtl); err != nil {
		return err
	}

	r.grpcRequest = &client_sdk_go.XSetRequest{
		CacheKey:        key,
		CacheBody:       value,
		TtlMilliseconds: ttl,
	}

	return nil
}

func (r *SetRequest) makeGrpcRequest(client scsDataClient, metadata context.Context) (grpcResponse, error) {
	resp, err := client.grpcClient.Set(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SetRequest) interpretGrpcResponse() error {
	r.response = SetSuccess{}
	return nil
}
