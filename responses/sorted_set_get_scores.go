package responses

// High level responses.

type SortedSetGetScoresResponse interface {
	isSortedSetGetScoresResponse()
}

// SortedSetGetScoresMiss Miss Response to a cache SortedSetScore api request.
type SortedSetGetScoresMiss struct{}

func (SortedSetGetScoresMiss) isSortedSetGetScoresResponse() {}

// SortedSetGetScoresHit Hit Response to a cache SortedSetScore api request.
type SortedSetGetScoresHit struct {
	scores []SortedSetGetScore
}

func (SortedSetGetScoresHit) isSortedSetGetScoresResponse() {}

func NewSortedSetGetScoresHit(scores []SortedSetGetScore) *SortedSetGetScoresHit {
	return &SortedSetGetScoresHit{scores: scores}
}

func (r SortedSetGetScoresHit) Scores() []SortedSetGetScore {
	return r.scores
}

// Responses for individual scores.

type SortedSetGetScore interface {
	isSortedSetScoreElement()
}

type SortedSetScore float64

func (SortedSetScore) isSortedSetScoreElement() {}

type SortedSetScoreMiss struct{}

func (SortedSetScoreMiss) isSortedSetScoreElement() {}

type SortedSetScoreInvalid struct{}

func (SortedSetScoreInvalid) isSortedSetScoreElement() {}
