package models

type ListFetchResponse interface {
	isListFetchResponse()
}

type ListFetchMiss struct{}

func (ListFetchMiss) isListFetchResponse() {}

type ListFetchHit struct {
	Value [][]byte
}

func (ListFetchHit) isListFetchResponse() {}

type ListLengthResponse interface {
	isListLengthResponse()
}

type ListLengthSuccess struct {
	Value uint32
}

func (ListLengthSuccess) isListLengthResponse() {}

type ListPushFrontResponse interface {
	isListPushFrontResponse()
}

type ListPushFrontSuccess struct {
	Value uint32
}

func (ListPushFrontSuccess) isListPushFrontResponse() {}

type ListPushBackResponse interface {
	isListPushBackResponse()
}

type ListPushBackSuccess struct {
	Value uint32
}

func (ListPushBackSuccess) isListPushBackResponse() {}

type ListPopFrontResponse interface {
	isListPopFrontResponse()
}

type ListPopFrontHit struct {
	Value []byte
}

func (ListPopFrontHit) isListPopFrontResponse() {}

type ListPopFrontMiss struct{}

func (ListPopFrontMiss) isListPopFrontResponse() {}

type ListPopBackResponse interface {
	isListPopBackResponse()
}

type ListPopBackHit struct {
	Value []byte
}

func (ListPopBackHit) isListPopBackResponse() {}

type ListPopBackMiss struct{}

func (ListPopBackMiss) isListPopBackResponse() {}
