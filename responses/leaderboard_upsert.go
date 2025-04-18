package responses

type LeaderboardUpsertResponse interface {
	MomentoLeaderboardResponse
	isLeaderboardUpsertResponse()
}

type LeaderboardUpsertSuccess struct{}

func (LeaderboardUpsertSuccess) isLeaderboardUpsertResponse() {}
