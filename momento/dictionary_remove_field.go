package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// DictionaryRemoveFieldResponse

type DictionaryRemoveFieldResponse interface {
	isDictionaryRemoveFieldResponse()
}

type DictionaryRemoveFieldSuccess struct{}

func (DictionaryRemoveFieldSuccess) isDictionaryRemoveFieldResponse() {}

// DictionaryRemoveFieldRequest

type DictionaryRemoveFieldRequest struct {
	CacheName      string
	DictionaryName string
	Field          Bytes

	grpcRequest  *pb.XDictionaryDeleteRequest
	grpcResponse *pb.XDictionaryDeleteResponse
	response     DictionaryRemoveFieldResponse
}

func (r *DictionaryRemoveFieldRequest) cacheName() string { return r.CacheName }

func (r *DictionaryRemoveFieldRequest) field() Bytes { return r.Field }

func (r *DictionaryRemoveFieldRequest) requestName() string { return "DictionaryRemoveField" }

func (r *DictionaryRemoveFieldRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return err
	}

	var fields [][]byte
	var field []byte
	if field, err = prepareField(r); err != nil {
		return err
	}
	fields = append(fields, field)
	r.grpcRequest = &pb.XDictionaryDeleteRequest{
		DictionaryName: []byte(r.DictionaryName),
		Delete:         &pb.XDictionaryDeleteRequest_Some_{Some: &pb.XDictionaryDeleteRequest_Some{Fields: fields}},
	}

	return nil
}

func (r *DictionaryRemoveFieldRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.DictionaryDelete(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *DictionaryRemoveFieldRequest) interpretGrpcResponse() error {
	r.response = &DictionaryRemoveFieldSuccess{}
	return nil
}
