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
	fields    []Bytes
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
			ret[string(resp.fields[idx].AsBytes())] = string(item.CacheBody)
		}
	}
	return ret
}

func (resp DictionaryGetFieldsHit) ValueMapStringBytes() map[string][]byte {
	ret := make(map[string][]byte)
	for idx, item := range resp.items {
		if item.Result == pb.ECacheResult_Hit {
			ret[string(resp.fields[idx].AsBytes())] = item.CacheBody
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
	Fields         []Bytes

	grpcRequest  *pb.XDictionaryGetRequest
	grpcResponse *pb.XDictionaryGetResponse
	response     DictionaryGetFieldsResponse
}

func (r *DictionaryGetFieldsRequest) cacheName() string { return r.CacheName }

func (r *DictionaryGetFieldsRequest) fields() []Bytes { return r.Fields }

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
		for idx, val := range rtype.Found.Items {
			if val.Result == pb.ECacheResult_Hit {
				responses = append(responses, &DictionaryGetFieldHit{field: RawBytes{val.CacheBody}})
			} else if val.Result == pb.ECacheResult_Miss {
				responses = append(responses, &DictionaryGetFieldMiss{field: r.Fields[idx]})
			} else {
				responses = append(responses, nil)
			}
		}
		r.response = &DictionaryGetFieldsHit{fields: r.Fields, items: rtype.Found.Items}
	default:
		return errUnexpectedGrpcResponse
	}
	return nil
}
