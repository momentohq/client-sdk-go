package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// DictionaryRemoveFieldsResponse

type DictionaryRemoveFieldsResponse interface {
	isDictionaryRemoveFieldsResponse()
}

type DictionaryRemoveFieldsSuccess struct{}

func (DictionaryRemoveFieldsSuccess) isDictionaryRemoveFieldsResponse() {}

// DictionaryRemoveFieldsRequest

type DictionaryRemoveFieldsRequest struct {
	CacheName      string
	DictionaryName string
	Fields         []Bytes

	grpcRequest  *pb.XDictionaryDeleteRequest
	grpcResponse *pb.XDictionaryDeleteResponse
	response     DictionaryRemoveFieldsResponse
}

func (r *DictionaryRemoveFieldsRequest) cacheName() string { return r.CacheName }

func (r *DictionaryRemoveFieldsRequest) fields() []Bytes { return r.Fields }

func (r *DictionaryRemoveFieldsRequest) requestName() string { return "DictionaryRemoveFields" }

func (r *DictionaryRemoveFieldsRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return err
	}

	fields := momentoBytesListToPrimitiveByteList(r.Fields)
	if fields, err = prepareFields(r); err != nil {
		return err
	}
	r.grpcRequest = &pb.XDictionaryDeleteRequest{
		DictionaryName: []byte(r.DictionaryName),
		Delete:         &pb.XDictionaryDeleteRequest_Some_{Some: &pb.XDictionaryDeleteRequest_Some{Fields: fields}},
	}

	return nil
}

func (r *DictionaryRemoveFieldsRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.DictionaryDelete(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *DictionaryRemoveFieldsRequest) interpretGrpcResponse() error {
	r.response = &DictionaryRemoveFieldsSuccess{}
	return nil
}
