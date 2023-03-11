package responses

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

func NewSortedSetGetScoreHit(score SortedSetGetScore) *SortedSetGetScoreHit {
	return &SortedSetGetScoreHit{score: score}
}

func (r SortedSetGetScoreHit) Score() SortedSetGetScore {
	return r.score
}
