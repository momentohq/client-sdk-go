package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ItemGetTypeRequest struct {
	CacheName string
	Key       Key

	grpcRequest *pb.XItemGetTypeRequest

	response responses.ItemGetTypeResponse
}

func (r *ItemGetTypeRequest) cacheName() string { return r.CacheName }

func (r *ItemGetTypeRequest) key() Key { return r.Key }

func (r *ItemGetTypeRequest) requestName() string { return "ItemGetType" }

func (r *ItemGetTypeRequest) initGrpcRequest(scsDataClient) error {
	var err error
	var key []byte

	if key, err = prepareKey(r); err != nil {
		return err
	}
	r.grpcRequest = &pb.XItemGetTypeRequest{CacheKey: key}

	return nil
}

func (r *ItemGetTypeRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ItemGetType(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *ItemGetTypeRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XItemGetTypeResponse)

	switch myResp.Result.(type) {
	case *pb.XItemGetTypeResponse_Found:
		r.response = responses.NewItemGetTypeHit(myResp.GetFound().ItemType)
		return nil
	case *pb.XItemGetTypeResponse_Missing:
		r.response = &responses.ItemGetTypeMiss{}
		return nil
	default:
		return errUnexpectedGrpcResponse(r, myResp)
	}
}
