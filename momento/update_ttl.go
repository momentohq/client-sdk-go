package momento

import (
	"context"
	"time"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
)

type UpdateTtlRequest struct {
	// Name of the cache to get the item from to be deleted
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// Time to live that you want to update in cache in seconds.
	UpdateTtl time.Duration

	grpcRequest  *pb.XUpdateTtlRequest
	grpcResponse *pb.XUpdateTtlResponse
	response     responses.UpdateTtlResponse
}

func (r *UpdateTtlRequest) cacheName() string { return r.CacheName }

func (r *UpdateTtlRequest) key() Key { return r.Key }

func (r *UpdateTtlRequest) updateTtl() time.Duration { return r.UpdateTtl }

func (r *UpdateTtlRequest) requestName() string { return "UpdateTtl" }

func (r *UpdateTtlRequest) initGrpcRequest(client scsDataClient) error {
	var err error
	var updateTtl uint64
	defaultTtlMills := uint64(client.defaultTtl.Milliseconds())

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return err
	}

	if updateTtl, err = prepareUpdateTtl(r); err != nil {
		return err
	}
	var request *pb.XUpdateTtlRequest
	if updateTtl > defaultTtlMills {
		request = &pb.XUpdateTtlRequest{CacheKey: key, UpdateTtl: &pb.XUpdateTtlRequest_IncreaseToMilliseconds{IncreaseToMilliseconds: updateTtl}}
	} else if updateTtl < defaultTtlMills {
		request = &pb.XUpdateTtlRequest{CacheKey: key, UpdateTtl: &pb.XUpdateTtlRequest_DecreaseToMilliseconds{DecreaseToMilliseconds: updateTtl}}
	} else {
		request = &pb.XUpdateTtlRequest{CacheKey: key, UpdateTtl: &pb.XUpdateTtlRequest_OverwriteToMilliseconds{OverwriteToMilliseconds: defaultTtlMills}}
	}
	r.grpcRequest = request

	return nil
}

func (r *UpdateTtlRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.UpdateTtl(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *UpdateTtlRequest) interpretGrpcResponse() error {
	grpcResp := r.grpcResponse

	var resp responses.UpdateTtlResponse
	switch grpcResp.Result.(type) {
	case *pb.XUpdateTtlResponse_NotSet:
		resp = &responses.UpdateTtlNotSet{}
	case *pb.XUpdateTtlResponse_Missing:
		resp = &responses.UpdateTtlMiss{}
	case *pb.XUpdateTtlResponse_Set:
		resp = &responses.UpdateTtlSet{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}

	r.response = resp

	return nil
}
