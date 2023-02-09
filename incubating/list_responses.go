package incubating

type ListFetchResponse interface {
	isListFetchResponse()
}

type ListFetchMiss struct{}

func (ListFetchMiss) isListFetchResponse() {}

type ListFetchHit struct {
	value       [][]byte
	stringValue []string
}

func (ListFetchHit) isListFetchResponse() {}

func (resp ListFetchHit) ValueListString() []string {
	if resp.stringValue == nil {
		for _, element := range resp.value {
			resp.stringValue = append(resp.stringValue, string(element))
		}
	}
	return resp.stringValue
}

func (resp ListFetchHit) ValueListByteArray() [][]byte {
	return resp.value
}

func (resp ListFetchHit) ValueList() []string {
	return resp.ValueListString()
}

type ListLengthResponse interface {
	isListLengthResponse()
}

type ListLengthSuccess struct {
	value uint32
}

func (ListLengthSuccess) isListLengthResponse() {}

func (resp ListLengthSuccess) Length() uint32 {
	return resp.value
}

type ListPushFrontResponse interface {
	isListPushFrontResponse()
}

type ListPushFrontSuccess struct {
	value uint32
}

func (ListPushFrontSuccess) isListPushFrontResponse() {}

func (resp ListPushFrontSuccess) ListLength() uint32 {
	return resp.value
}

type ListPushBackResponse interface {
	isListPushBackResponse()
}

type ListPushBackSuccess struct {
	value uint32
}

func (ListPushBackSuccess) isListPushBackResponse() {}

func (resp ListPushBackSuccess) ListLength() uint32 {
	return resp.value
}
