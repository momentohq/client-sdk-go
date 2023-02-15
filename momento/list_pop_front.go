package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// ListPopFrontResponse

type ListPopFrontResponse interface {
	isListPopFrontResponse()
}

type ListPopFrontHit struct {
	value Bytes
}

func (ListPopFrontHit) isListPopFrontResponse() {}

func (resp ListPopFrontHit) ValueByte() []byte {
	return resp.value.AsBytes()
}

func (resp ListPopFrontHit) ValueString() string {
	return string(resp.value.AsBytes())
}

type ListPopFrontMiss struct{}

func (ListPopFrontMiss) isListPopFrontResponse() {}

// ListPopFrontRequest

type ListPopFrontRequest struct {
	CacheName string
	ListName  string

	grpcRequest  *pb.XListPopFrontRequest
	grpcResponse *pb.XListPopFrontResponse
	response     ListPopFrontResponse
}

func (r ListPopFrontRequest) cacheName() string { return r.CacheName }

func (r ListPopFrontRequest) requestName() string { return "ListPopFront" }

func (r *ListPopFrontRequest) initGrpcRequest(client scsDataClient) error {
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
		r.response = ListPopFrontHit{value: RawBytes{rtype.Found.Front}}
	case *pb.XListPopFrontResponse_Missing:
		r.response = ListPopFrontMiss{}
	default:
		return errUnexpectedGrpcResponse
	}
	return nil
}
