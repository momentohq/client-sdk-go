package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ListPopBackRequest struct {
	CacheName string
	ListName  string

}

func (r *ListPopBackRequest) cacheName() string { return r.CacheName }

func (r *ListPopBackRequest) requestName() string { return "ListPopBack" }

func (r *ListPopBackRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	if _, err := prepareName(r.ListName, "List name"); err != nil {
		return nil, err
	}
	grpcRequest := &pb.XListPopBackRequest{
		ListName: []byte(r.ListName),
	}
	return grpcRequest, nil
}

func (r *ListPopBackRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ListPopBack(requestMetadata, grpcRequest.(*pb.XListPopBackRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *ListPopBackRequest) interpretGrpcResponse(resp interface{}) (interface{}, error) {
	myResp := resp.(*pb.XListPopBackResponse)
	switch rtype := myResp.List.(type) {
	case *pb.XListPopBackResponse_Found:
		return responses.NewListPopBackHit(rtype.Found.Back), nil
	case *pb.XListPopBackResponse_Missing:
		return &responses.ListPopBackMiss{}, nil
	default:
		return nil, errUnexpectedGrpcResponse(r, myResp)
	}
}
