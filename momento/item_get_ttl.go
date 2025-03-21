package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ItemGetTtlRequest struct {
	CacheName string
	Key       Key

}

func (r *ItemGetTtlRequest) cacheName() string { return r.CacheName }

func (r *ItemGetTtlRequest) key() Key { return r.Key }

func (r *ItemGetTtlRequest) requestName() string { return "ItemGetTtl" }

func (r *ItemGetTtlRequest) initGrpcRequest(scsDataClient) (interface{}, error) {
	var err error
	var key []byte

	if key, err = prepareKey(r); err != nil {
		return nil, err
	}
	grpcRequest := &pb.XItemGetTtlRequest{CacheKey: key}

	return grpcRequest, nil
}

func (r *ItemGetTtlRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ItemGetTtl(requestMetadata, grpcRequest.(*pb.XItemGetTtlRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *ItemGetTtlRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XItemGetTtlResponse)

	switch myResp.Result.(type) {
	case *pb.XItemGetTtlResponse_Found:
		return responses.NewItemGetTtlHit(myResp.GetFound().GetRemainingTtlMillis()), nil
	case *pb.XItemGetTtlResponse_Missing:
		return &responses.ItemGetTtlMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}
