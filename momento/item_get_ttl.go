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



	response responses.ItemGetTtlResponse
}

func (r *ItemGetTtlRequest) cacheName() string { return r.CacheName }

func (r *ItemGetTtlRequest) key() Key { return r.Key }

func (r *ItemGetTtlRequest) requestName() string { return "ItemGetTypeTL" }

func (r *ItemGetTtlRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
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

func (r *ItemGetTtlRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XItemGetTtlResponse)

	switch myResp.Result.(type) {
	case *pb.XItemGetTtlResponse_Found:
		r.response = responses.NewItemGetTtlHit(myResp.GetFound().GetRemainingTtlMillis())
		return nil
	case *pb.XItemGetTtlResponse_Missing:
		r.response = &responses.ItemGetTtlMiss{}
		return nil
	default:
		return errUnexpectedGrpcResponse(r, myResp)
	}
}

func (r *ItemGetTtlRequest) validateResponseType(resp grpcResponse) error {
	_, ok := resp.(*pb.XItemGetTtlResponse)
	if !ok {
		return errUnexpectedGrpcResponse(nil, resp)
	}
	return nil
}
