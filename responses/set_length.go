package responses

// SetLengthResponse is the base response type for a set length request.
type SetLengthResponse interface {
	isSetLengthResponse()
}

// SetLengthHit indicates a set length request was a hit.
type SetLengthHit struct {
	value uint32
}

func (SetLengthHit) isSetLengthResponse() {}

// Length returns the length of the set.
func (resp SetLengthHit) Length() uint32 {
	return resp.value
}

// SetLengthMiss indicates a set length request was a miss.
type SetLengthMiss struct{}

func (SetLengthMiss) isSetLengthResponse() {}

// NewSetLengthHit returns a new SetLengthHit containing the supplied value.
func NewSetLengthHit(value uint32) *SetLengthHit {
	return &SetLengthHit{value: value}
}
