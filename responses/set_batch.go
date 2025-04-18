package responses

// SetBatchResponse is the base response type for a batch set request.
type SetBatchResponse interface {
	MomentoCacheResponse
	isSetBatchResponse()
}

// SetBatchSuccess is the successful response to a batch set api request.
type SetBatchSuccess struct {
	responses []SetResponse
}

func (SetBatchSuccess) isSetBatchResponse() {}

// NewSetBatchSuccess returns a new SetBatchSuccess containing the supplied results.
func NewSetBatchSuccess(responses []SetResponse) *SetBatchSuccess {
	return &SetBatchSuccess{responses: responses}
}

// Results returns the data as a list of SetResponse objects.
func (resp SetBatchSuccess) Results() []SetResponse {
	return resp.responses
}
