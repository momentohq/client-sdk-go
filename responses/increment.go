package responses

// IncrementResponse is the base response type for a  increment request.
type IncrementResponse interface {
	isIncrementResponse()
}

// IncrementSuccess indicates a successful  increment success.
type IncrementSuccess struct {
	value int64
}

func (IncrementSuccess) isIncrementResponse() {}

// Value returns the new value of the element after incrementing.
func (resp IncrementSuccess) Value() int64 {
	return resp.value
}

// NewIncrementSuccess returns a new IncrementSuccess containing the supplied value.
func NewIncrementSuccess(value int64) *IncrementSuccess {
	return &IncrementSuccess{value: value}
}
