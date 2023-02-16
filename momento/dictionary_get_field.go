package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// DictionaryGetFieldResponse

type DictionaryGetFieldResponse interface {
	isDictionaryGetFieldResponse()
}

type DictionaryGetFieldHit struct {
	field []byte
}

func (DictionaryGetFieldHit) isDictionaryGetFieldResponse() {}

func (resp DictionaryGetFieldHit) FieldString() string {
	return string(resp.field)
}

func (resp DictionaryGetFieldHit) FieldByte() []byte {
	return resp.field
}

type DictionaryGetFieldMiss struct {
	field []byte
}

func (DictionaryGetFieldMiss) isDictionaryGetFieldResponse() {}

func (resp DictionaryGetFieldMiss) FieldString() string {
	return string(resp.field)
}

func (resp DictionaryGetFieldMiss) FieldByte() []byte {
	return resp.field
}

// DictionaryGetFieldRequest

type DictionaryGetFieldRequest struct {
	CacheName      string
	DictionaryName string
	Field          Value

	grpcRequest  *pb.XDictionaryGetRequest
	grpcResponse *pb.XDictionaryGetResponse
	response     DictionaryGetFieldResponse
}

func (r *DictionaryGetFieldRequest) cacheName() string { return r.CacheName }

func (r *DictionaryGetFieldRequest) field() Value { return r.Field }

func (r *DictionaryGetFieldRequest) requestName() string { return "DictionaryGetField" }

func (r *DictionaryGetFieldRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return err
	}

	var field []byte
	field, err = prepareField(r)
	if err != nil {
		return err
	}
	var fields [][]byte
	fields = append(fields, field)

	r.grpcRequest = &pb.XDictionaryGetRequest{
		DictionaryName: []byte(r.DictionaryName),
		Fields:         fields,
	}

	return nil
}

func (r *DictionaryGetFieldRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.DictionaryGet(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *DictionaryGetFieldRequest) interpretGrpcResponse() error {
	switch rtype := r.grpcResponse.Dictionary.(type) {
	case *pb.XDictionaryGetResponse_Missing:
		r.response = &DictionaryGetFieldMiss{field: r.Field.asBytes()}
	case *pb.XDictionaryGetResponse_Found:
		if rtype.Found.Items[0].Result == pb.ECacheResult_Miss {
			r.response = &DictionaryGetFieldMiss{field: r.Field.asBytes()}
		} else {
			r.response = &DictionaryGetFieldHit{field: rtype.Found.Items[0].CacheBody}
		}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}
