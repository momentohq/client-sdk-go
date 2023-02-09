package incubating

import (
	"github.com/momentohq/client-sdk-go/momento"
	incubating "github.com/momentohq/client-sdk-go/utils"
)

type SortedSetOrder int

const (
	ASCENDING  SortedSetOrder = 0
	DESCENDING SortedSetOrder = 1
)

type SortedSetScoreRequestElement struct {
	Name  momento.Bytes
	Score float64
}

type SortedSetRemoveRequestElement struct {
	Name momento.Bytes
}

type SortedSetPutRequest struct {
	CacheName     string
	SetName       momento.Bytes
	Elements      []*SortedSetScoreRequestElement
	CollectionTTL incubating.CollectionTtl
}

type SortedSetFetchNumResults interface {
	isSortedSetFetchNumResults()
}

type FetchAllItems struct{}

func (_ FetchAllItems) isSortedSetFetchNumResults() {}

type FetchLimitedItems struct {
	Limit uint32
}

func (_ FetchLimitedItems) isSortedSetFetchNumResults() {}

type SortedSetFetchRequest struct {
	CacheName       string
	SetName         momento.Bytes
	Order           SortedSetOrder
	NumberOfResults SortedSetFetchNumResults
}

type SortedSetGetScoreRequest struct {
	CacheName    string
	SetName      momento.Bytes
	ElementNames []momento.Bytes
}

type SortedSetRemoveRequest struct {
	CacheName        string
	SetName          momento.Bytes
	ElementsToRemove SortedSetRemoveNumItems
}

type SortedSetRemoveNumItems interface {
	isSortedSetRemoveNumItem()
}

type RemoveAllItems struct{}

func (_ RemoveAllItems) isSortedSetRemoveNumItem() {}

type RemoveSomeItems struct {
	elementsToRemove []momento.Bytes
}

func (_ RemoveSomeItems) isSortedSetRemoveNumItem() {}

type SortedSetGetRankRequest struct {
	CacheName   string
	SetName     momento.Bytes
	ElementName momento.Bytes
}

type SortedSetIncrementRequest struct {
	CacheName     string
	SetName       momento.Bytes
	ElementName   momento.Bytes
	Amount        uint64
	CollectionTTL incubating.CollectionTtl
}
