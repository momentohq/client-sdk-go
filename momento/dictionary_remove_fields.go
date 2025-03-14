package momento

import (
	"context"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type DictionaryRemoveFieldsRequest struct {
	CacheName      string
	DictionaryName string
	Fields         []Value

	grpcRequest *pb.XDictionaryDeleteRequest

	response responses.DictionaryRemoveFieldsResponse
}

func (r *DictionaryRemoveFieldsRequest) cacheName() string { return r.CacheName }

func (r *DictionaryRemoveFieldsRequest) fields() []Value { return r.Fields }

func (r *DictionaryRemoveFieldsRequest) requestName() string { return "DictionaryRemoveFields" }

func (r *DictionaryRemoveFieldsRequest) initGrpcRequest(client scsDataClient) (interface{}, error) {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return nil, err
	}

	var fields [][]byte
	if fields, err = prepareFields(r); err != nil {
		return nil, err
	}
	r.grpcRequest = &pb.XDictionaryDeleteRequest{
		DictionaryName: []byte(r.DictionaryName),
		Delete:         &pb.XDictionaryDeleteRequest_Some_{Some: &pb.XDictionaryDeleteRequest_Some{Fields: fields}},
	}

	return r.grpcRequest, nil
}

func (r *DictionaryRemoveFieldsRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	resp, err := client.grpcClient.DictionaryDelete(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))
	responseMetadata := []metadata.MD{header, trailer}
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *DictionaryRemoveFieldsRequest) interpretGrpcResponse(_ interface{}) error {
	r.response = &responses.DictionaryRemoveFieldsSuccess{}
	return nil
}

func (r *DictionaryRemoveFieldsRequest) validateResponseType(resp grpcResponse) error {
	_, ok := resp.(*pb.XDictionaryDeleteResponse)
	if !ok {
		return errUnexpectedGrpcResponse(nil, resp)
	}
	return nil
}
