package momento

import (
	"context"
	"time"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type DecreaseTtlRequest struct {
	// Name of the cache to get the item from to be deleted
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// Time to live that you want to decrease to.
	Ttl time.Duration
}

func (r *DecreaseTtlRequest) cacheName() string { return r.CacheName }

func (r *DecreaseTtlRequest) key() Key { return r.Key }

func (r *DecreaseTtlRequest) updateTtl() time.Duration { return r.Ttl }

func (r *DecreaseTtlRequest) requestName() string { return "DecreaseTtl" }

func (r *DecreaseTtlRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error
	var ttl uint64

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return nil, err
	}

	if ttl, err = prepareUpdateTtl(r); err != nil {
		return nil, err
	}
	grpcRequest := &pb.XUpdateTtlRequest{CacheKey: key, UpdateTtl: &pb.XUpdateTtlRequest_DecreaseToMilliseconds{DecreaseToMilliseconds: ttl}}

	return grpcRequest, nil
}

func (r *DecreaseTtlRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.UpdateTtl(requestMetadata, grpcRequest.(*pb.XUpdateTtlRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}

	return resp, nil, nil
}

func (r *DecreaseTtlRequest) interpretGrpcResponse(theResponse interface{}) (interface{}, error) {
	myResp := theResponse.(*pb.XUpdateTtlResponse)

	var resp responses.DecreaseTtlResponse
	switch myResp.Result.(type) {
	case *pb.XUpdateTtlResponse_NotSet:
		resp = &responses.DecreaseTtlNotSet{}
	case *pb.XUpdateTtlResponse_Missing:
		resp = &responses.DecreaseTtlMiss{}
	case *pb.XUpdateTtlResponse_Set:
		resp = &responses.DecreaseTtlSet{}
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
	return resp, nil
}
