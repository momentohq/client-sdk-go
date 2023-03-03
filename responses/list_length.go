package responses

type ListLengthResponse interface {
	isListLengthResponse()
}

type ListLengthHit struct {
	value uint32
}

func (ListLengthHit) isListLengthResponse() {}

func (resp ListLengthHit) Length() uint32 {
	return resp.value
}

type ListLengthMiss struct{}

func (ListLengthMiss) isListLengthResponse() {}

func NewListLengthHit(value uint32) *ListLengthHit {
	return &ListLengthHit{value: value}
}
