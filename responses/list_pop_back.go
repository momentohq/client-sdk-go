package responses

// ListPopBackResponse is the base response type for a list pop back request.
type ListPopBackResponse interface {
	isListPopBackResponse()
}

// ListPopBackHit indicates a successful list pop back request.
type ListPopBackHit struct {
	value []byte
}

func (ListPopBackHit) isListPopBackResponse() {}

// ValueByte returns the data as a byte array.
func (resp ListPopBackHit) ValueByte() []byte {
	return resp.value
}

// ValueString returns the data as an utf-8 string, decoded from the underlying byte array.
func (resp ListPopBackHit) ValueString() string {
	return string(resp.value)
}

// ListPopBackMiss indicates a list pop back request was a miss.
type ListPopBackMiss struct{}

func (ListPopBackMiss) isListPopBackResponse() {}

// NewListPopBackHit returns a new ListPopBackHit containing the supplied value.
func NewListPopBackHit(value []byte) *ListPopBackHit {
	return &ListPopBackHit{value: value}
}
