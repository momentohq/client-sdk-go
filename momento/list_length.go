package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// ListLengthResponse

type ListLengthResponse interface {
	isListLengthResponse()
}

type ListLengthHit struct {
	value uint32
}

func (ListLengthHit) isListLengthResponse() {}

func (resp ListLengthHit) Length() uint32 {
	return resp.value
}

type ListLengthMiss struct{}

func (ListLengthMiss) isListLengthResponse() {}

// ListLengthRequest

type ListLengthRequest struct {
	CacheName string
	ListName  string

	grpcRequest  *pb.XListLengthRequest
	grpcResponse *pb.XListLengthResponse
	response     ListLengthResponse
}

func (r ListLengthRequest) cacheName() string { return r.CacheName }

func (r ListLengthRequest) requestName() string { return "ListLength" }

func (r *ListLengthRequest) initGrpcRequest(client scsDataClient) error {
	if _, err := prepareName(r.ListName, "List name"); err != nil {
		return err
	}

	r.grpcRequest = &pb.XListLengthRequest{
		ListName: []byte(r.ListName),
	}

	return nil
}

func (r *ListLengthRequest) makeGrpcRequest(client scsDataClient, ctx context.Context) (grpcResponse, error) {
	resp, err := client.grpcClient.ListLength(ctx, r.grpcRequest)
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
		r.response = ListLengthHit{value: rtype.Found.Length}
	case *pb.XListLengthResponse_Missing:
		r.response = ListLengthMiss{}
	default:
		return errUnexpectedGrpcResponse
	}
	return nil
}
