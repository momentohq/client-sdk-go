package responses

type SetContainsElementsResponse interface {
	isSetContainsElementsResponse()
}

type SetContainsElementsMiss struct{}

func (*SetContainsElementsMiss) isSetContainsElementsResponse() {}

type SetContainsElementsHit struct {
	values []bool
}

func (*SetContainsElementsHit) isSetContainsElementsResponse() {}

func (r *SetContainsElementsHit) ContainsElements() []bool {
	return r.values
}

func NewSetContainsElementsHit(values []bool) *SetContainsElementsHit {
	return &SetContainsElementsHit{values: values}
}
