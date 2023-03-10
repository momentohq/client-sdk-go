package responses

// ListLengthResponse is the base response type for a list length request.
type ListLengthResponse interface {
	isListLengthResponse()
}

// ListLengthHit indicates a list length request was a hit.
type ListLengthHit struct {
	value uint32
}

func (ListLengthHit) isListLengthResponse() {}

// Length returns the length of the list.
func (resp ListLengthHit) Length() uint32 {
	return resp.value
}

// ListLengthMiss indicates a list length request was a miss.
type ListLengthMiss struct{}

func (ListLengthMiss) isListLengthResponse() {}

// NewListLengthHit returns a new ListLengthHit containing the supplied value.
func NewListLengthHit(value uint32) *ListLengthHit {
	return &ListLengthHit{value: value}
}
