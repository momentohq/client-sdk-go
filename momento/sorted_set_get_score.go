package momento

type SortedSetGetScoreRequest struct {
	CacheName string
	SetName   string
	Value     Value
}
