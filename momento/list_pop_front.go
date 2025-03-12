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

	grpcRequest  *pb.XListPopFrontRequest
	grpcResponse *pb.XListPopFrontResponse
	response     responses.ListPopFrontResponse
}

func (r *ListPopFrontRequest) cacheName() string { return r.CacheName }

func (r *ListPopFrontRequest) requestName() string { return "ListPopFront" }

func (r *ListPopFrontRequest) initGrpcRequest(scsDataClient) error {
	if _, err := prepareName(r.ListName, "List name"); err != nil {
		return err
	}
	r.grpcRequest = &pb.XListPopFrontRequest{
		ListName: []byte(r.ListName),
	}
	return nil
}

func (r *ListPopFrontRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ListPopFront(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	r.grpcResponse = resp
	return resp, nil, nil
}

func (r *ListPopFrontRequest) interpretGrpcResponse(_ interface{}) error {
	switch rtype := r.grpcResponse.List.(type) {
	case *pb.XListPopFrontResponse_Found:
		r.response = responses.NewListPopFrontHit(rtype.Found.Front)
	case *pb.XListPopFrontResponse_Missing:
		r.response = &responses.ListPopFrontMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}

func (r *ListPopFrontRequest) getResponse() interface{} {
	return r.response
}
