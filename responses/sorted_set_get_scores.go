package responses

// SortedSetGetScoresResponse is the base response type for a sorted set get responses request.
type SortedSetGetScoresResponse interface {
	isSortedSetGetScoresResponse()
}

// SortedSetGetScoresMiss Miss Response to a cache SortedSetGetScores api request.
type SortedSetGetScoresMiss struct{}

func (SortedSetGetScoresMiss) isSortedSetGetScoresResponse() {}

// SortedSetGetScoresHit Hit Response to a cache SortedSetGetScores api request.
type SortedSetGetScoresHit struct {
	responses []SortedSetGetScoreResponse
	values    [][]byte
}

func (SortedSetGetScoresHit) isSortedSetGetScoresResponse() {}

// NewSortedSetGetScoresHit returns a new SortedSetGetScoresHit containing the supplied responses.
func NewSortedSetGetScoresHit(responses []SortedSetGetScoreResponse, values [][]byte) *SortedSetGetScoresHit {
	return &SortedSetGetScoresHit{responses: responses, values: values}
}

// Responses returns an array of SortedSetGetScoreResponse which will either be
// of type SortedSetGetScoreHit, SortedSetGetScoreMiss, or SortedSetGetScoreInvalid
func (r SortedSetGetScoresHit) Responses() []SortedSetGetScoreResponse {
	return r.responses
}

// ScoresArray returns an array of float64 values that represent the hit responses. Misses
// are not represented the array.
func (r SortedSetGetScoresHit) ScoresArray() []float64 {
	var hits []float64
	for _, v := range r.responses {
		switch vType := v.(type) {
		case *SortedSetGetScoreHit:
			hits = append(hits, vType.score)
		}
	}
	return hits
}

// ScoresMap returns a map with string keys representing the originally supplied values
// and float64 values representing the corresponding score. Misses are not represented
// in the map.
func (r SortedSetGetScoresHit) ScoresMap() map[string]float64 {
	hits := make(map[string]float64)
	for index, v := range r.responses {
		switch vType := v.(type) {
		case *SortedSetGetScoreHit:
			hits[string(r.values[index])] = vType.score
		}
	}
	return hits
}
