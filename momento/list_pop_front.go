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



	response responses.ListPopFrontResponse
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

func (r *ListPopFrontRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XListPopFrontResponse)
	switch rtype := myResp.List.(type) {
	case *pb.XListPopFrontResponse_Found:
		r.response = responses.NewListPopFrontHit(rtype.Found.Front)
	case *pb.XListPopFrontResponse_Missing:
		r.response = &responses.ListPopFrontMiss{}
	default:
		return errUnexpectedGrpcResponse(r, myResp)
	}
	return nil
}

func (r *ListPopFrontRequest) validateResponseType(resp grpcResponse) error {
	_, ok := resp.(*pb.XListPopFrontResponse)
	if !ok {
		return errUnexpectedGrpcResponse(nil, resp)
	}
	return nil
}
