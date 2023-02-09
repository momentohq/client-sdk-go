package models

type ListFetchResponse interface {
	isListFetchResponse()
}

type ListFetchMiss struct{}

func (_ ListFetchMiss) isListFetchResponse() {}

type ListFetchHit struct {
	Value [][]byte
}

func (_ ListFetchHit) isListFetchResponse() {}

type ListLengthResponse interface {
	isListLengthResponse()
}

type ListLengthSuccess struct {
	Value uint32
}

func (_ ListLengthSuccess) isListLengthResponse() {}

type ListPushFrontResponse interface {
	isListPushFrontResponse()
}

type ListPushFrontSuccess struct {
	Value uint32
}

func (_ ListPushFrontSuccess) isListPushFrontResponse() {}

type ListPushBackResponse interface {
	isListPushBackResponse()
}

type ListPushBackSuccess struct {
	Value uint32
}

func (_ ListPushBackSuccess) isListPushBackResponse() {}
