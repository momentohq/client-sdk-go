package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ListLengthRequest struct {
	CacheName string
	ListName  string



	response responses.ListLengthResponse
}

func (r *ListLengthRequest) cacheName() string { return r.CacheName }

func (r *ListLengthRequest) requestName() string { return "ListLength" }

func (r *ListLengthRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	if _, err := prepareName(r.ListName, "List name"); err != nil {
		return nil, err
	}

	grpcRequest := &pb.XListLengthRequest{
		ListName: []byte(r.ListName),
	}

	return grpcRequest, nil
}

func (r *ListLengthRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ListLength(requestMetadata, grpcRequest.(*pb.XListLengthRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *ListLengthRequest) interpretGrpcResponse(resp interface{}) error {
	myResp := resp.(*pb.XListLengthResponse)
	switch rtype := myResp.List.(type) {
	case *pb.XListLengthResponse_Found:
		r.response = responses.NewListLengthHit(rtype.Found.Length)
	case *pb.XListLengthResponse_Missing:
		r.response = &responses.ListLengthMiss{}
	default:
		return errUnexpectedGrpcResponse(r, myResp)
	}
	return nil
}

func (r *ListLengthRequest) validateResponseType(resp grpcResponse) error {
	_, ok := resp.(*pb.XListLengthResponse)
	if !ok {
		return errUnexpectedGrpcResponse(nil, resp)
	}
	return nil
}

