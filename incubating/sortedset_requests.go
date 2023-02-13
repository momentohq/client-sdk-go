package incubating

import (
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/utils"
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
	SetName       string
	Elements      []*SortedSetScoreRequestElement
	CollectionTTL utils.CollectionTTL
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
	SetName         string
	Order           SortedSetOrder
	NumberOfResults SortedSetFetchNumResults
}

type SortedSetGetScoreRequest struct {
	CacheName    string
	SetName      string
	ElementNames []momento.Bytes
}

type SortedSetRemoveRequest struct {
	CacheName        string
	SetName          string
	ElementsToRemove SortedSetRemoveNumElements
}

type SortedSetRemoveNumElements interface {
	isSortedSetRemoveNumElement()
}

type RemoveAllElements struct{}

func (RemoveAllElements) isSortedSetRemoveNumElement() {}

type RemoveSomeElements struct {
	elementsToRemove []momento.Bytes
}

func (RemoveSomeElements) isSortedSetRemoveNumElement() {}

type SortedSetGetRankRequest struct {
	CacheName   string
	SetName     string
	ElementName momento.Bytes
}

type SortedSetIncrementRequest struct {
	CacheName     string
	SetName       string
	ElementName   momento.Bytes
	Amount        uint64
	CollectionTTL utils.CollectionTTL
}
