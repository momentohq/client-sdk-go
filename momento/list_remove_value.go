package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// ListRemoveValueResponse

type ListRemoveValueResponse interface {
	isListRemoveValueResponse()
}

type ListRemoveValueSuccess struct{}

func (ListRemoveValueSuccess) isListRemoveValueResponse() {}

// ListRemoveValueRequest

type ListRemoveValueRequest struct {
	CacheName string
	ListName  string
	Value     Bytes

	grpcRequest  *pb.XListRemoveRequest
	grpcResponse *pb.XListRemoveResponse
	response     ListRemoveValueResponse
}

func (r ListRemoveValueRequest) cacheName() string { return r.CacheName }

func (r ListRemoveValueRequest) value() Bytes { return r.Value }

func (r ListRemoveValueRequest) requestName() string { return "ListRemoveValue" }

func (r *ListRemoveValueRequest) initGrpcRequest(client scsDataClient) error {
	var err error

	if _, err = prepareName(r.ListName, "List name"); err != nil {
		return err
	}

	var value []byte
	if value, err = prepareValue(r); err != nil {
		return err
	}

	r.grpcRequest = &pb.XListRemoveRequest{
		ListName: []byte(r.ListName),
		Remove:   &pb.XListRemoveRequest_AllElementsWithValue{AllElementsWithValue: value},
	}

	return nil
}

func (r *ListRemoveValueRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.ListRemove(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *ListRemoveValueRequest) interpretGrpcResponse() error {
	r.response = ListRemoveValueSuccess{}
	return nil
}
