package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ListRemoveValueRequest struct {
	CacheName string
	ListName  string
	Value     Value
}

func (r *ListRemoveValueRequest) cacheName() string { return r.CacheName }

func (r *ListRemoveValueRequest) value() Value { return r.Value }

func (r *ListRemoveValueRequest) requestName() string { return "ListRemoveValue" }

func (r *ListRemoveValueRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.ListName, "List name"); err != nil {
		return nil, err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return nil, err
	}

	grpcRequest := &pb.XListRemoveRequest{
		ListName: []byte(r.ListName),
		Remove:   &pb.XListRemoveRequest_AllElementsWithValue{AllElementsWithValue: value},
	}

	return grpcRequest, nil
}

func (r *ListRemoveValueRequest) makeGrpcRequest(grpcRequest interface{}, requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.ListRemove(requestMetadata, grpcRequest.(*pb.XListRemoveRequest), grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *ListRemoveValueRequest) interpretGrpcResponse(_ interface{}) (interface{}, error) {
	return &responses.ListRemoveValueSuccess{}, nil
}
