package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ListPopFrontRequest struct {
	CacheName string
	ListName  string
}

func (r *ListPopFrontRequest) cacheName() string { return r.CacheName }

func (r *ListPopFrontRequest) requestName() string { return "ListPopFront" }

func (r *ListPopFrontRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	if _, err := prepareName(r.ListName, "List name"); err != nil {
		return nil, err
	}
	grpcRequest := &pb.XListPopFrontRequest{
		ListName: []byte(r.ListName),
	}
	return grpcRequest, nil
}

func (r *ListPopFrontRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ListPopFront(requestMetadata, grpcRequest.(*pb.XListPopFrontRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *ListPopFrontRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XListPopFrontResponse)
	switch rtype := myResp.List.(type) {
	case *pb.XListPopFrontResponse_Found:
		return responses.NewListPopFrontHit(rtype.Found.Front), nil
	case *pb.XListPopFrontResponse_Missing:
		return &responses.ListPopFrontMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}

func (c ListPopFrontRequest) GetRequestName() string {
	return "ListPopFrontRequest"
}
