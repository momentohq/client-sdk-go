package momento

import (
	"context"
	"time"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type UpdateTtlRequest struct {
	// Name of the cache to get the item from to be deleted
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// Time to live that you want to update in cache in seconds.
	Ttl time.Duration

	grpcRequest  *pb.XUpdateTtlRequest
	grpcResponse *pb.XUpdateTtlResponse
	response     responses.UpdateTtlResponse
}

func (r *UpdateTtlRequest) cacheName() string { return r.CacheName }

func (r *UpdateTtlRequest) key() Key { return r.Key }

func (r *UpdateTtlRequest) updateTtl() time.Duration { return r.Ttl }

func (r *UpdateTtlRequest) requestName() string { return "UpdateTtl" }

func (r *UpdateTtlRequest) initGrpcRequest(client scsDataClient) error {
	var err error
	var ttl uint64

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return err
	}

	if ttl, err = prepareUpdateTtl(r); err != nil {
		return err
	}
	r.grpcRequest = &pb.XUpdateTtlRequest{CacheKey: key, UpdateTtl: &pb.XUpdateTtlRequest_OverwriteToMilliseconds{OverwriteToMilliseconds: ttl}}

	return nil
}

func (r *UpdateTtlRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.UpdateTtl(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}

	r.grpcResponse = resp

	return resp, nil, nil
}

func (r *UpdateTtlRequest) interpretGrpcResponse() error {
	grpcResp := r.grpcResponse

	var resp responses.UpdateTtlResponse
	switch grpcResp.Result.(type) {
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

func (r *UpdateTtlRequest) getResponse() interface{} { return r.response }
