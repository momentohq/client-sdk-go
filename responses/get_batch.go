package responses

// GetBatchResponse is the base response type for a batch get request.
type GetBatchResponse interface {
	isGetBatchResponse()
}

// GetBatchSuccess is the successful response to a batch get api request.
type GetBatchSuccess struct {
	responses []GetResponse
	keys      [][]byte
}

func (GetBatchSuccess) isGetBatchResponse() {}

// NewGetBatchSuccess returns a new GetBatchSuccess containing the supplied results.
func NewGetBatchSuccess(responses []GetResponse, keys [][]byte) *GetBatchSuccess {
	return &GetBatchSuccess{responses: responses, keys: keys}
}

// ValueMap returns the data as a Map whose keys and values are utf-8 strings,
// decoded from the underlying byte arrays.
// This is a convenience alias for ValueMapStringString.
func (resp GetBatchSuccess) ValueMap() map[string]string {
	return resp.ValueMapStringString()
}

// ValueMapStringString returns the data as a Map whose keys and values are utf-8 strings,
// decoded from the underlying byte arrays. Misses are represented as empty strings.
func (resp GetBatchSuccess) ValueMapStringString() map[string]string {
	ret := make(map[string]string)
	for idx, element := range resp.responses {
		switch e := element.(type) {
		case *GetHit:
			ret[string(resp.keys[idx])] = e.ValueString()
		}
	}
	return ret
}

// ValueMapStringBytes returns the data as a Map whose keys are utf-8 strings,
// decoded from the underlying byte array, and whose values are byte arrays.
// Misses are represented as nil.
func (resp GetBatchSuccess) ValueMapStringBytes() map[string][]byte {
	ret := make(map[string][]byte)
	for idx, element := range resp.responses {
		switch e := element.(type) {
		case *GetHit:
			ret[string(resp.keys[idx])] = e.ValueByte()
		}
	}
	return ret
}

// Results returns the data as a list of GetResponse objects.
func (resp GetBatchSuccess) Results() []GetResponse {
	return resp.responses
}
