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

type ListPopFrontResponse interface {
	isListPopFrontResponse()
}

type ListPopFrontHit struct {
	value []byte
}

func (ListPopFrontHit) isListPopFrontResponse() {}

func (resp ListPopFrontHit) ValueByteArray() []byte {
	return resp.value
}

func (resp ListPopFrontHit) ValueString() string {
	return string(resp.value)
}

type ListPopFrontMiss struct{}

func (ListPopFrontMiss) isListPopFrontResponse() {}

type ListPopBackResponse interface {
	isListPopBackResponse()
}

type ListPopBackHit struct {
	value []byte
}

func (ListPopBackHit) isListPopBackResponse() {}

func (resp ListPopBackHit) ValueByteArray() []byte {
	return resp.value
}

func (resp ListPopBackHit) ValueString() string {
	return string(resp.value)
}

type ListPopBackMiss struct{}

func (ListPopBackMiss) isListPopBackResponse() {}
