package momento

type LeaderboardRemoveElementsRequest struct {
	Ids []uint32
}

type LeaderboardInternalRemoveElementsRequest struct {
	CacheName       string
	LeaderboardName string
	Ids             []uint32
}
