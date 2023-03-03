package responses

import pb "github.com/momentohq/client-sdk-go/internal/protos"

type DictionaryGetFieldsResponse interface {
	isDictionaryGetFieldsResponse()
}

type DictionaryGetFieldsHit struct {
	elements  []*pb.XDictionaryGetResponse_XDictionaryGetResponsePart
	fields    [][]byte
	responses []DictionaryGetFieldResponse
}

func (DictionaryGetFieldsHit) isDictionaryGetFieldsResponse() {}

func (resp DictionaryGetFieldsHit) ValueMap() map[string]string {
	return resp.ValueMapStringString()
}

func (resp DictionaryGetFieldsHit) ValueMapStringString() map[string]string {
	ret := make(map[string]string)
	for idx, element := range resp.elements {
		if element.Result == pb.ECacheResult_Hit {
			ret[string(resp.fields[idx])] = string(element.CacheBody)
		}
	}
	return ret
}

func (resp DictionaryGetFieldsHit) ValueMapStringBytes() map[string][]byte {
	ret := make(map[string][]byte)
	for idx, element := range resp.elements {
		if element.Result == pb.ECacheResult_Hit {
			ret[string(resp.fields[idx])] = element.CacheBody
		}
	}
	return ret
}

func (resp DictionaryGetFieldsHit) Responses() []DictionaryGetFieldResponse {
	return resp.responses
}

type DictionaryGetFieldsMiss struct{}

func (DictionaryGetFieldsMiss) isDictionaryGetFieldsResponse() {}

func NewDictionaryGetFieldsHit(
	fields [][]byte, elements []*pb.XDictionaryGetResponse_XDictionaryGetResponsePart, responses []DictionaryGetFieldResponse,
) *DictionaryGetFieldsHit {
	return &DictionaryGetFieldsHit{
		elements:  elements,
		fields:    fields,
		responses: responses,
	}
}
