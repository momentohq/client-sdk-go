package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ListPopBackRequest struct {
	CacheName string
	ListName  string

	grpcRequest  *pb.XListPopBackRequest
	grpcResponse *pb.XListPopBackResponse
	response     responses.ListPopBackResponse
}

func (r *ListPopBackRequest) cacheName() string { return r.CacheName }

func (r *ListPopBackRequest) requestName() string { return "ListPopBack" }

func (r *ListPopBackRequest) initGrpcRequest(scsDataClient) error {
	if _, err := prepareName(r.ListName, "List name"); err != nil {
		return err
	}
	r.grpcRequest = &pb.XListPopBackRequest{
		ListName: []byte(r.ListName),
	}
	return nil
}

func (r *ListPopBackRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.ListPopBack(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *ListPopBackRequest) interpretGrpcResponse() error {
	switch rtype := r.grpcResponse.List.(type) {
	case *pb.XListPopBackResponse_Found:
		r.response = responses.NewListPopBackHit(rtype.Found.Back)
	case *pb.XListPopBackResponse_Missing:
		r.response = &responses.ListPopBackMiss{}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}

func (r *ListPopBackRequest) getResponse() map[string]string {
	return getMomentoResponseData(r.response)
}
