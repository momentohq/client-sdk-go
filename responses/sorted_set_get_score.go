package responses

type SortedSetScoreElement interface {
	isSortedSetScoreElement()
}

type SortedSetGetScoreResponse interface {
	isSortedSetGetScoreResponse()
}

// SortedSetGetScoreMiss Miss Response to a cache SortedSetScore api request.
type SortedSetGetScoreMiss struct{}

func (SortedSetGetScoreMiss) isSortedSetGetScoreResponse() {}

// SortedSetGetScoreHit Hit Response to a cache SortedSetScore api request.
type SortedSetGetScoreHit struct {
	Elements []SortedSetScoreElement
}

func (SortedSetGetScoreHit) isSortedSetGetScoreResponse() {}

type SortedSetScore float64

func (SortedSetScore) isSortedSetScoreElement() {}

type SortedSetScoreMiss struct{}

func (SortedSetScoreMiss) isSortedSetScoreElement() {}

type SortedSetScoreInvalid struct{}

func (SortedSetScoreInvalid) isSortedSetScoreElement() {}
