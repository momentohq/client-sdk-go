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
}

func (r *ItemGetTypeRequest) cacheName() string { return r.CacheName }

func (r *ItemGetTypeRequest) key() Key { return r.Key }

func (r *ItemGetTypeRequest) requestName() string { return "ItemGetType" }

func (r *ItemGetTypeRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error
	var key []byte

	if key, err = prepareKey(r); err != nil {
		return nil, err
	}
	grpcRequest := &pb.XItemGetTypeRequest{CacheKey: key}

	return grpcRequest, nil
}

func (r *ItemGetTypeRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ItemGetType(requestMetadata, grpcRequest.(*pb.XItemGetTypeRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *ItemGetTypeRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XItemGetTypeResponse)

	switch myResp.Result.(type) {
	case *pb.XItemGetTypeResponse_Found:
		return responses.NewItemGetTypeHit(myResp.GetFound().ItemType), nil
	case *pb.XItemGetTypeResponse_Missing:
		return &responses.ItemGetTypeMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}

func (c ItemGetTypeRequest) GetRequestName() string {
	return "ItemGetTypeRequest"
}
