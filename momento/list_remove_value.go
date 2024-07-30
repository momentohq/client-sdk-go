package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type ListRemoveValueRequest struct {
	CacheName string
	ListName  string
	Value     Value

	grpcRequest  *pb.XListRemoveRequest
	grpcResponse *pb.XListRemoveResponse
	response     responses.ListRemoveValueResponse
}

func (r *ListRemoveValueRequest) cacheName() string { return r.CacheName }

func (r *ListRemoveValueRequest) value() Value { return r.Value }

func (r *ListRemoveValueRequest) requestName() string { return "ListRemoveValue" }

func (r *ListRemoveValueRequest) initGrpcRequest(scsDataClient) error {
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
	r.response = &responses.ListRemoveValueSuccess{}
	return nil
}

func (r *ListRemoveValueRequest) getResponse() map[string]string {
	return getMomentoResponseData(r.response)
}
