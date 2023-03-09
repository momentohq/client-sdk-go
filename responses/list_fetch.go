package responses

// ListFetchResponse is a base response type for a list fetch request.
type ListFetchResponse interface {
	isListFetchResponse()
}

type ListFetchHit struct {
	value       [][]byte
	stringValue []string
}

func (ListFetchHit) isListFetchResponse() {}

// ValueListByte returns the data as an array of byte arrays.
func (resp ListFetchHit) ValueListByte() [][]byte {
	return resp.value
}

// ValueListString returns the data as an array of strings, decoded from the underlying byte array.
func (resp ListFetchHit) ValueListString() []string {
	if resp.stringValue == nil {
		for _, element := range resp.value {
			resp.stringValue = append(resp.stringValue, string(element))
		}
	}
	return resp.stringValue
}

// ValueList returns the data as an array of strings, decoded from the underlying byte array.
// This is a convenience alias ValueListString.
func (resp ListFetchHit) ValueList() []string {
	return resp.ValueListString()
}

// ListFetchMiss indicates a list fetch was a miss.
type ListFetchMiss struct{}

func (ListFetchMiss) isListFetchResponse() {}

// NewListFetchHit returns a new ListFetchHit contains value.
func NewListFetchHit(value [][]byte) *ListFetchHit {
	return &ListFetchHit{value: value}
}
