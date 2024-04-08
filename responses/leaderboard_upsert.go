package responses

type LeaderboardUpsertResponse interface {
	isLeaderboardUpsertResponse()
}

type LeaderboardUpsertSuccess struct{}

func (LeaderboardUpsertSuccess) isLeaderboardUpsertResponse() {}
