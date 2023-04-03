package momento

type SortedSetGetScoreRequest struct {
	CacheName string
	SetName   string
	Value     Value
}

func (r *SortedSetGetScoreRequest) cacheName() string { return r.CacheName }
