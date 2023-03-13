package responses

// SortedSetGetScoreResponse is the base response type for a sorted set get score request.
type SortedSetGetScoreResponse interface {
	isSortedSetGetScoreResponse()
}

// SortedSetGetScoreMiss Miss Response to a cache SortedSetGetScore api request.
type SortedSetGetScoreMiss struct{}

func (SortedSetGetScoreMiss) isSortedSetGetScoreResponse() {}

// SortedSetGetScoreInvalid Invalid Response to a cache SortedSetGetScore api request.
type SortedSetGetScoreInvalid struct{}

func (SortedSetGetScoreInvalid) isSortedSetGetScoreResponse() {}

// SortedSetGetScoreHit Hit Response to a cache SortedSetGetScore api request.
type SortedSetGetScoreHit struct {
	score float64
}

func (SortedSetGetScoreHit) isSortedSetGetScoreResponse() {}

// NewSortedSetGetScoreHit returns a new SortedSetGetScoreHit containing the supplied score.
func NewSortedSetGetScoreHit(score float64) *SortedSetGetScoreHit {
	return &SortedSetGetScoreHit{score: score}
}

// Score returns the float64 score value.
func (r SortedSetGetScoreHit) Score() float64 {
	return r.score
}
