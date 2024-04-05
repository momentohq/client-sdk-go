package momento

type LeaderboardFetchByRankRequest struct {
	StartRank uint32
	EndRank   uint32
	Order     *LeaderboardOrder
}

type LeaderboardInternalFetchByRankRequest struct {
	CacheName       string
	LeaderboardName string
	StartRank       uint32
	EndRank         uint32
	Order           *LeaderboardOrder
}
