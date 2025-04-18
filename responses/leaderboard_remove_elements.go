package responses

type LeaderboardRemoveElementsResponse interface {
	MomentoLeaderboardResponse
	isLeaderboardRemoveElementsResponse()
}

type LeaderboardRemoveElementsSuccess struct{}

func (LeaderboardRemoveElementsSuccess) isLeaderboardRemoveElementsResponse() {}
