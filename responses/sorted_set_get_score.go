package responses

// SortedSetGetScoreResponse is the base response type for a sorted set get score request.
type SortedSetGetScoreResponse interface {
	isSortedSetGetScoreResponse()
}

// SortedSetGetScoreMiss Miss Response to a cache SortedSetScore api request.
type SortedSetGetScoreMiss struct{}

func (SortedSetGetScoreMiss) isSortedSetGetScoreResponse() {}

// SortedSetGetScoreHit Hit Response to a cache SortedSetScore api request.
type SortedSetGetScoreHit struct {
	score SortedSetGetScore
}

func (SortedSetGetScoreHit) isSortedSetGetScoreResponse() {}

// NewSortedSetGetScoreHit returns a new SortedSetGetScoreHit containing the supplied score.
func NewSortedSetGetScoreHit(score SortedSetGetScore) *SortedSetGetScoreHit {
	return &SortedSetGetScoreHit{score: score}
}

// Score returns a SortedSetGetScore.
func (r SortedSetGetScoreHit) Score() SortedSetGetScore {
	return r.score
}
