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
}

func (r *UpdateTtlRequest) cacheName() string { return r.CacheName }

func (r *UpdateTtlRequest) key() Key { return r.Key }

func (r *UpdateTtlRequest) updateTtl() time.Duration { return r.Ttl }

func (r *UpdateTtlRequest) requestName() string { return "UpdateTtl" }

func (r *UpdateTtlRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error
	var ttl uint64

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return nil, err
	}

	if ttl, err = prepareUpdateTtl(r); err != nil {
		return nil, err
	}
	grpcRequest := &pb.XUpdateTtlRequest{CacheKey: key, UpdateTtl: &pb.XUpdateTtlRequest_OverwriteToMilliseconds{OverwriteToMilliseconds: ttl}}

	return grpcRequest, nil
}

func (r *UpdateTtlRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.UpdateTtl(requestMetadata, grpcRequest.(*pb.XUpdateTtlRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *UpdateTtlRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XUpdateTtlResponse)

	switch myResp.Result.(type) {
	case *pb.XUpdateTtlResponse_Missing:
		return &responses.UpdateTtlMiss{}, nil
	case *pb.XUpdateTtlResponse_Set:
		return &responses.UpdateTtlSet{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}
