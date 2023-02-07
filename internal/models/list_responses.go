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
