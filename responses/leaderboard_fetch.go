package responses

type LeaderboardFetchResponse interface {
	isLeaderboardFetchResponse()
}

type LeaderboardElement struct {
	Id    uint32
	Score float64
	Rank  uint32
}

type LeaderboardFetchSuccess struct {
	elements []LeaderboardElement
}

func (LeaderboardFetchSuccess) isLeaderboardFetchResponse() {}

// NewLeaderboardFetchSuccess returns a new LeaderboardFetchSuccess containing the supplied elements.
func NewLeaderboardFetchSuccess(elements []LeaderboardElement) *LeaderboardFetchSuccess {
	if elements == nil {
		elements = []LeaderboardElement{}
	}
	return &LeaderboardFetchSuccess{elements: elements}
}

func (s LeaderboardFetchSuccess) Values() []LeaderboardElement {
	return s.elements
}
