package momento

import (
	"context"
	"time"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type IncreaseTtlRequest struct {
	// Name of the cache to get the item from to be deleted
	CacheName string
	// string or byte key to be used to store item.
	Key Key
	// Time to live that you want to increase to.
	Ttl time.Duration

	response responses.IncreaseTtlResponse
}

func (r *IncreaseTtlRequest) cacheName() string { return r.CacheName }

func (r *IncreaseTtlRequest) key() Key { return r.Key }

func (r *IncreaseTtlRequest) updateTtl() time.Duration { return r.Ttl }

func (r *IncreaseTtlRequest) requestName() string { return "IncreaseTtl" }

func (r *IncreaseTtlRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error
	var ttl uint64

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return nil, err
	}

	if ttl, err = prepareUpdateTtl(r); err != nil {
		return nil, err
	}
	grpcRequest := &pb.XUpdateTtlRequest{CacheKey: key, UpdateTtl: &pb.XUpdateTtlRequest_IncreaseToMilliseconds{IncreaseToMilliseconds: ttl}}

	return grpcRequest, nil
}

func (r *IncreaseTtlRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.UpdateTtl(requestMetadata, grpcRequest.(*pb.XUpdateTtlRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *IncreaseTtlRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XUpdateTtlResponse)

	var theResponse responses.IncreaseTtlResponse
	switch myResp.Result.(type) {
	case *pb.XUpdateTtlResponse_NotSet:
		theResponse = &responses.IncreaseTtlNotSet{}
	case *pb.XUpdateTtlResponse_Missing:
		theResponse = &responses.IncreaseTtlMiss{}
	case *pb.XUpdateTtlResponse_Set:
		theResponse = &responses.IncreaseTtlSet{}
	default:
		return errUnexpectedGrpcResponse(r, myResp)
	}

	r.response = theResponse

	return nil
}

func (r *IncreaseTtlRequest) validateResponseType(resp grpcResponse) error {
	_, ok := resp.(*pb.XUpdateTtlResponse)
	if !ok {
		return errUnexpectedGrpcResponse(nil, resp)
	}
	return nil
}
