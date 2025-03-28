package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type GetWithHashRequest struct {
	// Name of the cache to get the item from
	CacheName string
	// string or byte key to be used to store item
	Key Key
}

func (r *GetWithHashRequest) cacheName() string { return r.CacheName }

func (r *GetWithHashRequest) key() Key { return r.Key }

func (r *GetWithHashRequest) requestName() string { return "GetWithHash" }

func (r *GetWithHashRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return nil, err
	}

	grpcRequest := &pb.XGetWithHashRequest{
		CacheKey: key,
	}

	return grpcRequest, nil
}

func (r *GetWithHashRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.GetWithHash(
		requestMetadata,
		grpcRequest.(*pb.XGetWithHashRequest),
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *GetWithHashRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	grpcResponse := resp.(*pb.XGetWithHashResponse)
	switch response := grpcResponse.Result.(type) {
	case *pb.XGetWithHashResponse_Found:
		return responses.NewGetWithHashHit(response.Found.Value, response.Found.Hash), nil
	case *pb.XGetWithHashResponse_Missing:
		return &responses.GetWithHashMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, grpcResponse)
	}
}
