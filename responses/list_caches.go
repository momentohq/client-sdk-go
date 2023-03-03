package responses

type ListCachesResponse interface {
	isListCachesResponse()
}

// ListCachesSuccess Output of the List caches operation.
type ListCachesSuccess struct {
	nextToken string
	caches    []CacheInfo
}

func (ListCachesSuccess) isListCachesResponse() {}

func NewListCachesSuccess(nextToken string, caches []CacheInfo) *ListCachesSuccess {
	return &ListCachesSuccess{
		nextToken: nextToken,
		caches:    caches, //convertCacheInfo(caches),
	}
}

// NextToken Next Page Token returned by Simple Cache Service along with the list of caches.
// If nextToken is present, then this token must be provided in the next call to continue paginating through the list.
// This is done by setting this value in ListCachesRequest.
func (resp ListCachesSuccess) NextToken() string {
	return resp.nextToken
}

// Caches Returns all caches.
func (resp ListCachesSuccess) Caches() []CacheInfo {
	return resp.caches
}

// CacheInfo Information about a Cache.
type CacheInfo struct {
	name string
}

// Name Returns cache's name.
func (ci CacheInfo) Name() string {
	return ci.name
}

func NewCacheInfo(name string) CacheInfo {
	return CacheInfo{name: name}
}

//func convertCacheInfo(i []CacheInfo) []CacheInfo {
//	var convertedList []CacheInfo
//	for _, c := range i {
//		convertedList = append(convertedList, CacheInfo{
//			name: c.Name,
//		})
//	}
//	return convertedList
//}
