package momento

type ListCachesRequest struct{}

func (c ListCachesRequest) GetRequestName() string {
	return "ListCachesRequest"
}
