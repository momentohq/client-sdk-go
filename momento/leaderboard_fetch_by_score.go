package momento

type LeaderboardFetchByScoreRequest struct {
	MinScore *float64
	MaxScore *float64
	Order    *LeaderboardOrder
	Offset   *uint32
	Count    *uint32
}

type LeaderboardInternalFetchByScoreRequest struct {
	CacheName       string
	LeaderboardName string
	MinScore        *float64
	MaxScore        *float64
	Order           *LeaderboardOrder
	Offset          *uint32
	Count           *uint32
}
