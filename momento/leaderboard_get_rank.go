package momento

type LeaderboardGetRankRequest struct {
	Ids   []uint32
	Order *LeaderboardOrder
}

type LeaderboardInternalGetRankRequest struct {
	CacheName       string
	LeaderboardName string
	Ids             []uint32
	Order           *LeaderboardOrder
}
