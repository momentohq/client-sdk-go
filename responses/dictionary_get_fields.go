package responses

import pb "github.com/momentohq/client-sdk-go/internal/protos"

// DictionaryGetFieldsResponse is a base response type for a dictionary fields request.
type DictionaryGetFieldsResponse interface {
	isDictionaryGetFieldsResponse()
}

// DictionaryGetFieldsHit Indicates that the requested data was successfully retrieved from the cache.  Provides
// `Value*` accessors to retrieve the data in the appropriate format.
type DictionaryGetFieldsHit struct {
	elements  []*pb.XDictionaryGetResponse_XDictionaryGetResponsePart
	fields    [][]byte
	responses []DictionaryGetFieldResponse
}

func (DictionaryGetFieldsHit) isDictionaryGetFieldsResponse() {}

// ValueMap returns the data as a Map whose keys and values are utf-8 strings, decoded from the underlying byte arrays.
// This is a convenience alias for ValueMapStringString.
func (resp DictionaryGetFieldsHit) ValueMap() map[string]string {
	return resp.ValueMapStringString()
}

// ValueMapStringString returns the data as a Map whose keys and values are utf-8 strings, decoded from the underlying byte arrays.
func (resp DictionaryGetFieldsHit) ValueMapStringString() map[string]string {
	ret := make(map[string]string)
	for idx, element := range resp.elements {
		if element.Result == pb.ECacheResult_Hit {
			ret[string(resp.fields[idx])] = string(element.CacheBody)
		}
	}
	return ret
}

// ValueMapStringBytes returns the data as a Map whose keys are utf-8 strings, decoded from the underlying byte array, and whose values are byte arrays.
func (resp DictionaryGetFieldsHit) ValueMapStringBytes() map[string][]byte {
	ret := make(map[string][]byte)
	for idx, element := range resp.elements {
		if element.Result == pb.ECacheResult_Hit {
			ret[string(resp.fields[idx])] = element.CacheBody
		}
	}
	return ret
}

// Responses returns an array of DictionaryGetFieldResponse.
func (resp DictionaryGetFieldsHit) Responses() []DictionaryGetFieldResponse {
	return resp.responses
}

// DictionaryGetFieldsMiss indicates that the requested data was not available in the cache.
type DictionaryGetFieldsMiss struct{}

func (DictionaryGetFieldsMiss) isDictionaryGetFieldsResponse() {}

// NewDictionaryGetFieldsHit returns a new DictionaryGetFieldsHit contains elements, fields, and responses.
func NewDictionaryGetFieldsHit(
	fields [][]byte, elements []*pb.XDictionaryGetResponse_XDictionaryGetResponsePart, responses []DictionaryGetFieldResponse,
) *DictionaryGetFieldsHit {
	return &DictionaryGetFieldsHit{
		elements:  elements,
		fields:    fields,
		responses: responses,
	}
}
