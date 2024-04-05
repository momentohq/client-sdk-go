package responses

type LeaderboardDeleteResponse interface {
	isLeaderboardDeleteResponse()
}

type LeaderboardDeleteSuccess struct{}

func (LeaderboardDeleteSuccess) isLeaderboardDeleteResponse() {}
