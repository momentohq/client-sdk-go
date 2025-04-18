package responses

type LeaderboardDeleteResponse interface {
	MomentoLeaderboardResponse
	isLeaderboardDeleteResponse()
}

type LeaderboardDeleteSuccess struct{}

func (LeaderboardDeleteSuccess) isLeaderboardDeleteResponse() {}
