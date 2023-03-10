package responses

// SortedSetGetScoresResponse is the base response type for a sorted set get scores request.
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

// NewSortedSetGetScoresHit returns a new SortedSetGetScoresHit containing the supplied scores.
func NewSortedSetGetScoresHit(scores []SortedSetGetScore) *SortedSetGetScoresHit {
	return &SortedSetGetScoresHit{scores: scores}
}

// Scores returns an array of SortedSetGetScore.
func (r SortedSetGetScoresHit) Scores() []SortedSetGetScore {
	return r.scores
}

// SortedSetGetScore is the base response for individual scores.
type SortedSetGetScore interface {
	isSortedSetScoreElement()
}

// SortedSetScore indicates a sorted set score request was a hit.
type SortedSetScore float64

func (SortedSetScore) isSortedSetScoreElement() {}

// SortedSetScoreMiss indicates a sorted set score request was a miss.
type SortedSetScoreMiss struct{}

func (SortedSetScoreMiss) isSortedSetScoreElement() {}

// SortedSetScoreInvalid indicates an unknown response was returned for a sorted set score request.
type SortedSetScoreInvalid struct{}

func (SortedSetScoreInvalid) isSortedSetScoreElement() {}
