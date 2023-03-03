package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ListLengthRequest struct {
	CacheName string
	ListName  string

	grpcRequest  *pb.XListLengthRequest
	grpcResponse *pb.XListLengthResponse
	response     responses.ListLengthResponse
}

func (r *ListLengthRequest) cacheName() string { return r.CacheName }

func (r *ListLengthRequest) requestName() string { return "ListLength" }

func (r *ListLengthRequest) initGrpcRequest(scsDataClient) error {
	if _, err := prepareName(r.ListName, "List name"); err != nil {
		return err
	}

	r.grpcRequest = &pb.XListLengthRequest{
		ListName: []byte(r.ListName),
	}

	return nil
}

func (r *ListLengthRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.ListLength(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *ListLengthRequest) interpretGrpcResponse() error {
	resp := r.grpcResponse
	switch rtype := resp.List.(type) {
	case *pb.XListLengthResponse_Found:
		r.response = responses.NewListLengthHit(rtype.Found.Length)
	case *pb.XListLengthResponse_Missing:
		r.response = &responses.ListLengthMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}
