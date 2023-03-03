package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

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

func (r *ListPopFrontRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.ListPopFront(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *ListPopFrontRequest) interpretGrpcResponse() error {
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
