package responses

// SetContainsElementsResponse is a base response type for a set contains elements request.
type SetContainsElementsResponse interface {
	isSetContainsElementsResponse()
}

// SetContainsElementsMiss indicates a set contains elements request was a miss.
type SetContainsElementsMiss struct{}

func (*SetContainsElementsMiss) isSetContainsElementsResponse() {}

// SetContainsElementsHit indicates a set contains element request was a hit.
type SetContainsElementsHit struct {
	values []bool
}

func (*SetContainsElementsHit) isSetContainsElementsResponse() {}

// ContainsElements returns an array of bool indicating the existence of each element in the given set.
func (r *SetContainsElementsHit) ContainsElements() []bool {
	return r.values
}

// NewSetContainsElementsHit returns a new SetContainsElementsHit contains values.
func NewSetContainsElementsHit(values []bool) *SetContainsElementsHit {
	return &SetContainsElementsHit{values: values}
}
