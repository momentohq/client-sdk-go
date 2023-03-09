package responses

// ListPopFrontResponse is a base response type for a list pop front request.
type ListPopFrontResponse interface {
	isListPopFrontResponse()
}

// ListPopFrontHit indicates a successful list pop front request.
type ListPopFrontHit struct {
	value []byte
}

func (ListPopFrontHit) isListPopFrontResponse() {}

// ValueByte returns the data as a byte array.
func (resp ListPopFrontHit) ValueByte() []byte {
	return resp.value
}

// ValueString returns the data as an utf-8 string, decoded from the underlying byte array.
func (resp ListPopFrontHit) ValueString() string {
	return string(resp.value)
}

// ListPopFrontMiss indicates a list pop front request was a miss.
type ListPopFrontMiss struct{}

func (ListPopFrontMiss) isListPopFrontResponse() {}

// NewListPopFrontHit returns a new ListPopFrontHit contains value.
func NewListPopFrontHit(value []byte) *ListPopFrontHit {
	return &ListPopFrontHit{value: value}
}
