package responses

type SortedSetScoreElement interface {
	isSortedSetScoreElement()
}

type SortedSetGetScoresResponse interface {
	isSortedSetGetScoresResponse()
}

// SortedSetGetScoresMiss Miss Response to a cache SortedSetScore api request.
type SortedSetGetScoresMiss struct{}

func (SortedSetGetScoresMiss) isSortedSetGetScoresResponse() {}

// SortedSetGetScoresHit Hit Response to a cache SortedSetScore api request.
type SortedSetGetScoresHit struct {
	Elements []SortedSetScoreElement
}

func (SortedSetGetScoresHit) isSortedSetGetScoresResponse() {}

type SortedSetScore float64

func (SortedSetScore) isSortedSetScoreElement() {}

type SortedSetScoreMiss struct{}

func (SortedSetScoreMiss) isSortedSetScoreElement() {}

type SortedSetScoreInvalid struct{}

func (SortedSetScoreInvalid) isSortedSetScoreElement() {}
