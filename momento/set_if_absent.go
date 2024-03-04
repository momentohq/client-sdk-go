package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SetIfAbsentRequest struct {
	// Name of the cache to store the item in.
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// string ot byte value to be stored.
	Value Value
	// Optional Time to live in cache in seconds.
	// If not provided, then default TTL for the cache client instance is used.
	Ttl time.Duration

	// We issue a SetIfNotExists request to the server instead of a SetIf request because
	// the backend implementation of SetIfNotExists is more efficient than SetIf.
	grpcRequest  *pb.XSetIfRequest
	grpcResponse *pb.XSetIfResponse
	response     responses.SetIfAbsentResponse
}

func (r *SetIfAbsentRequest) cacheName() string { return r.CacheName }

func (r *SetIfAbsentRequest) key() Key { return r.Key }

func (r *SetIfAbsentRequest) value() Value { return r.Value }

func (r *SetIfAbsentRequest) ttl() time.Duration { return r.Ttl }

func (r *SetIfAbsentRequest) requestName() string { return "SetIfNotExists" }

func (r *SetIfAbsentRequest) initGrpcRequest(client scsDataClient) error {
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

func (r *SetIfAbsentRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SetIf(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *SetIfAbsentRequest) interpretGrpcResponse() error {
	grpcResp := r.grpcResponse
	var resp responses.SetIfAbsentResponse

	switch grpcResp.Result.(type) {
	case *pb.XSetIfResponse_Stored:
		resp = &responses.SetIfAbsentStored{}
	case *pb.XSetIfResponse_NotStored:
		resp = &responses.SetIfAbsentNotStored{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}

	r.response = resp
	return nil
}
