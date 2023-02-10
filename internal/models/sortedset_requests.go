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

type FetchAllItems struct{}

func (FetchAllItems) isSortedSetFetchNumResults() {}

type FetchLimitedItems struct {
	Limit uint32
}

func (FetchLimitedItems) isSortedSetFetchNumResults() {}

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
	ElementsToRemove SortedSetRemoveNumItems
}
type SortedSetRemoveNumItems interface {
	isSortedSetRemoveNumItem()
}

type RemoveAllItems struct{}

func (RemoveAllItems) isSortedSetRemoveNumItem() {}

type RemoveSomeItems struct {
	ElementsToRemove [][]byte
}

func (RemoveSomeItems) isSortedSetRemoveNumItem() {}

type SortedSetGetRankRequest struct {
	CacheName   string
	SetName     []byte
	ElementName []byte
}

type SortedSetIncrementRequest struct {
	CacheName     string
	SetName       []byte
	ElementName   []byte
	Amount        uint64
	CollectionTTL incubating.CollectionTTL
}
