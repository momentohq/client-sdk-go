package momento

type SortedSetRemoveElementRequest struct {
	CacheName string
	SetName   string
	Value     Value
}
