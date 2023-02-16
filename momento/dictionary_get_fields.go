package momento

import (
	"context"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// DictionaryGetFieldsResponse

type DictionaryGetFieldsResponse interface {
	isDictionaryGetFieldsResponse()
}

type DictionaryGetFieldsHit struct {
	items     []*pb.XDictionaryGetResponse_XDictionaryGetResponsePart
	fields    [][]byte
	responses []DictionaryGetFieldResponse
}

func (DictionaryGetFieldsHit) isDictionaryGetFieldsResponse() {}

func (resp DictionaryGetFieldsHit) ValueMap() map[string]string {
	return resp.ValueMapStringString()
}

func (resp DictionaryGetFieldsHit) ValueMapStringString() map[string]string {
	ret := make(map[string]string)
	for idx, item := range resp.items {
		if item.Result == pb.ECacheResult_Hit {
			ret[string(resp.fields[idx])] = string(item.CacheBody)
		}
	}
	return ret
}

func (resp DictionaryGetFieldsHit) ValueMapStringBytes() map[string][]byte {
	ret := make(map[string][]byte)
	for idx, item := range resp.items {
		if item.Result == pb.ECacheResult_Hit {
			ret[string(resp.fields[idx])] = item.CacheBody
		}
	}
	return ret
}

type DictionaryGetFieldsMiss struct{}

func (DictionaryGetFieldsMiss) isDictionaryGetFieldsResponse() {}

// DictionaryGetFieldsRequest

type DictionaryGetFieldsRequest struct {
	CacheName      string
	DictionaryName string
	Fields         []Value

	grpcRequest  *pb.XDictionaryGetRequest
	grpcResponse *pb.XDictionaryGetResponse
	response     DictionaryGetFieldsResponse
}

func (r *DictionaryGetFieldsRequest) cacheName() string { return r.CacheName }

func (r *DictionaryGetFieldsRequest) fields() []Value { return r.Fields }

func (r *DictionaryGetFieldsRequest) requestName() string { return "DictionaryGetFields" }

func (r *DictionaryGetFieldsRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.DictionaryName, "Dictionary name"); err != nil {
		return err
	}

	var fields [][]byte
	if fields, err = prepareFields(r); err != nil {
		return err
	}

	r.grpcRequest = &pb.XDictionaryGetRequest{
		DictionaryName: []byte(r.DictionaryName),
		Fields:         fields,
	}

	return nil
}

func (r *DictionaryGetFieldsRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.DictionaryGet(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}
	r.grpcResponse = resp
	return resp, nil
}

func (r *DictionaryGetFieldsRequest) interpretGrpcResponse() error {
	switch rtype := r.grpcResponse.Dictionary.(type) {
	case *pb.XDictionaryGetResponse_Missing:
		r.response = &DictionaryGetFieldsMiss{}
	case *pb.XDictionaryGetResponse_Found:
		var responses []DictionaryGetFieldResponse
		var fields [][]byte
		for idx, val := range rtype.Found.Items {
			var field []byte
			if val.Result == pb.ECacheResult_Hit {
				field = val.CacheBody
				responses = append(responses, &DictionaryGetFieldHit{field: field})
			} else if val.Result == pb.ECacheResult_Miss {
				field = r.Fields[idx].asBytes()
				responses = append(responses, &DictionaryGetFieldMiss{field: field})
			} else {
				field = r.Fields[idx].asBytes()
				responses = append(responses, nil)
			}
			fields = append(fields, field)
		}
		r.response = &DictionaryGetFieldsHit{fields: fields, items: rtype.Found.Items, responses: responses}
	default:
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
	return nil
}
