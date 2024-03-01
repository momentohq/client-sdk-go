package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"

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

	grpcRequest  *pb.XSetIfRequest
	grpcResponse *pb.XSetIfResponse
	response     responses.SetIfNotEqualResponse
}

func (r *SetIfNotEqualRequest) cacheName() string { return r.CacheName }

func (r *SetIfNotEqualRequest) key() Key { return r.Key }

func (r *SetIfNotEqualRequest) value() Value { return r.Value }

func (r *SetIfNotEqualRequest) notEqual() Value { return r.NotEqual }

func (r *SetIfNotEqualRequest) ttl() time.Duration { return r.Ttl }

func (r *SetIfNotEqualRequest) requestName() string { return "SetIfNotExists" }

func (r *SetIfNotEqualRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return err
	}

	var notEqual []byte
	if notEqual, err = prepareNotEqual(r); err != nil {
		return err
	}

	var ttl uint64
	if ttl, err = prepareTtl(r, client.defaultTtl); err != nil {
		return err
	}

	condition := &pb.XSetIfRequest_NotEqual{
		NotEqual: &pb.NotEqual{
			ValueToCheck: notEqual,
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

func (r *SetIfNotEqualRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SetIf(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SetIfNotEqualRequest) interpretGrpcResponse() error {
	grpcResp := r.grpcResponse
	var resp responses.SetIfNotEqualResponse

	switch grpcResp.Result.(type) {
	case *pb.XSetIfResponse_Stored:
		resp = &responses.SetIfNotEqualStored{}
	case *pb.XSetIfResponse_NotStored:
		resp = &responses.SetIfNotEqualNotStored{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}

	r.response = resp
	return nil
}
