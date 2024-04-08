package momento

type LeaderboardUpsertRequest struct {
	Elements []LeaderboardUpsertElement
}

type LeaderboardInternalUpsertRequest struct {
	CacheName       string
	LeaderboardName string
	Elements        []LeaderboardUpsertElement
}
