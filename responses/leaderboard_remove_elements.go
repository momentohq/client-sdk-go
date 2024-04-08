package responses

type LeaderboardRemoveElementsResponse interface {
	isLeaderboardRemoveElementsResponse()
}

type LeaderboardRemoveElementsSuccess struct{}

func (LeaderboardRemoveElementsSuccess) isLeaderboardRemoveElementsResponse() {}
