package momento

import (
	"context"
	"time"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
)

type IncreaseTtlRequest struct {
	// Name of the cache to get the item from to be deleted
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// Time to live that you want to increase to.
	Ttl time.Duration

	grpcRequest  *pb.XUpdateTtlRequest
	grpcResponse *pb.XUpdateTtlResponse
	response     responses.IncreaseTtlResponse
}

func (r *IncreaseTtlRequest) cacheName() string { return r.CacheName }

func (r *IncreaseTtlRequest) key() Key { return r.Key }

func (r *IncreaseTtlRequest) updateTtl() time.Duration { return r.Ttl }

func (r *IncreaseTtlRequest) requestName() string { return "IncreaseTtl" }

func (r *IncreaseTtlRequest) initGrpcRequest(client scsDataClient) error {
	var err error
	var ttl uint64

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return err
	}

	if ttl, err = prepareUpdateTtl(r); err != nil {
		return err
	}
	r.grpcRequest = &pb.XUpdateTtlRequest{CacheKey: key, UpdateTtl: &pb.XUpdateTtlRequest_IncreaseToMilliseconds{IncreaseToMilliseconds: ttl}}

	return nil
}

func (r *IncreaseTtlRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.UpdateTtl(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *IncreaseTtlRequest) interpretGrpcResponse() error {
	grpcResp := r.grpcResponse

	var resp responses.IncreaseTtlResponse
	switch grpcResp.Result.(type) {
	case *pb.XUpdateTtlResponse_NotSet:
		resp = &responses.IncreaseTtlNotSet{}
	case *pb.XUpdateTtlResponse_Missing:
		resp = &responses.IncreaseTtlMiss{}
	case *pb.XUpdateTtlResponse_Set:
		resp = &responses.IncreaseTtlSet{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}

	r.response = resp

	return nil
}
