package models

import (
	incubating "github.com/momentohq/client-sdk-go/utils"
)

type SortedSetOrder int

const (
	ASCENDING  SortedSetOrder = 1
	DESCENDING SortedSetOrder = 2
)

type SortedSetElement struct {
	Name  []byte
	Score float64
}

type SortedSetPutRequest struct {
	CacheName     string
	SetName       []byte
	Elements      []*SortedSetElement
	CollectionTTL incubating.CollectionTTL
}

type SortedSetFetchNumResults interface {
	isSortedSetFetchNumResults()
}

type FetchAllElements struct{}

func (FetchAllElements) isSortedSetFetchNumResults() {}

type FetchLimitedElements struct {
	Limit uint32
}

func (FetchLimitedElements) isSortedSetFetchNumResults() {}

type SortedSetFetchRequest struct {
	CacheName       string
	SetName         []byte
	Order           SortedSetOrder
	NumberOfResults SortedSetFetchNumResults
}

type SortedSetGetScoreRequest struct {
	CacheName    string
	SetName      []byte
	ElementNames [][]byte
}

type SortedSetRemoveRequest struct {
	CacheName        string
	SetName          []byte
	ElementsToRemove SortedSetRemoveNumElements
}
type SortedSetRemoveNumElements interface {
	isSortedSetRemoveNumElement()
}

type RemoveAllElements struct{}

func (RemoveAllElements) isSortedSetRemoveNumElement() {}

type RemoveSomeElements struct {
	ElementsToRemove [][]byte
}

func (RemoveSomeElements) isSortedSetRemoveNumElement() {}

type SortedSetGetRankRequest struct {
	CacheName   string
	SetName     []byte
	ElementName []byte
}

type SortedSetIncrementRequest struct {
	CacheName     string
	SetName       []byte
	ElementName   []byte
	Amount        float64
	CollectionTTL incubating.CollectionTTL
}
