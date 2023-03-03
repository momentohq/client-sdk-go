package responses

type ListFetchResponse interface {
	isListFetchResponse()
}

type ListFetchHit struct {
	value       [][]byte
	stringValue []string
}

func (ListFetchHit) isListFetchResponse() {}

func (resp ListFetchHit) ValueListByte() [][]byte {
	return resp.value
}

func (resp ListFetchHit) ValueListString() []string {
	if resp.stringValue == nil {
		for _, element := range resp.value {
			resp.stringValue = append(resp.stringValue, string(element))
		}
	}
	return resp.stringValue
}

func (resp ListFetchHit) ValueList() []string {
	return resp.ValueListString()
}

type ListFetchMiss struct{}

func (ListFetchMiss) isListFetchResponse() {}

func NewListFetchHit(value [][]byte) *ListFetchHit {
	return &ListFetchHit{value: value}
}
