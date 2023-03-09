package momento

import (
	"context"
	"time"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
)

type DecreaseTtlRequest struct {
	// Name of the cache to get the item from to be deleted
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// Time to live that you want to decrease to.
	Ttl time.Duration

	grpcRequest  *pb.XUpdateTtlRequest
	grpcResponse *pb.XUpdateTtlResponse
	response     responses.DecreaseTtlResponse
}

func (r *DecreaseTtlRequest) cacheName() string { return r.CacheName }

func (r *DecreaseTtlRequest) key() Key { return r.Key }

func (r *DecreaseTtlRequest) updateTtl() time.Duration { return r.Ttl }

func (r *DecreaseTtlRequest) requestName() string { return "DecreaseTtl" }

func (r *DecreaseTtlRequest) initGrpcRequest(client scsDataClient) error {
	var err error
	var ttl uint64

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return err
	}

	if ttl, err = prepareUpdateTtl(r); err != nil {
		return err
	}
	r.grpcRequest = &pb.XUpdateTtlRequest{CacheKey: key, UpdateTtl: &pb.XUpdateTtlRequest_DecreaseToMilliseconds{DecreaseToMilliseconds: ttl}}

	return nil
}

func (r *DecreaseTtlRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.UpdateTtl(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *DecreaseTtlRequest) interpretGrpcResponse() error {
	grpcResp := r.grpcResponse

	var resp responses.UpdateTtlResponse
	switch grpcResp.Result.(type) {
	case *pb.XUpdateTtlResponse_NotSet:
		resp = &responses.DecreaseTtlNotSet{}
	case *pb.XUpdateTtlResponse_Missing:
		resp = &responses.DecreaseTtlMiss{}
	case *pb.XUpdateTtlResponse_Set:
		resp = &responses.DecreaseTtlSet{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}

	r.response = resp

	return nil
}
