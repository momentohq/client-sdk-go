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

	grpcRequest  *pb.XItemGetTtlRequest

	response     responses.ItemGetTtlResponse
}

func (r *ItemGetTtlRequest) cacheName() string { return r.CacheName }

func (r *ItemGetTtlRequest) key() Key { return r.Key }

func (r *ItemGetTtlRequest) requestName() string { return "ItemGetTypeTL" }

func (r *ItemGetTtlRequest) initGrpcRequest(scsDataClient) error {
	var err error
	var key []byte

	if key, err = prepareKey(r); err != nil {
		return err
	}
	r.grpcRequest = &pb.XItemGetTtlRequest{CacheKey: key}

	return nil
}

func (r *ItemGetTtlRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ItemGetTtl(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
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

func (r *ItemGetTtlRequest) getResponse() interface{} {
	return r.response
}
