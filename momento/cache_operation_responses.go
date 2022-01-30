package momento

type ListCachesResponse struct {
	nextToken string
	caches    []CacheInfo
}

func (resp *ListCachesResponse) NextToken() string {
	return resp.nextToken
}

func (resp *ListCachesResponse) Caches() []CacheInfo {
	return resp.caches
}

type CacheInfo struct {
	name string
}

func (ci CacheInfo) Name() string {
	return ci.name
}

const (
	HIT  string = "HIT"
	MISS string = "MISS"
)

type GetCacheResponse struct {
	value  []byte
	result string
}

func (resp *GetCacheResponse) StringValue() string {
	if resp.result == HIT {
		return string(resp.value)
	}
	return ""
}

func (resp *GetCacheResponse) ByteValue() []byte {
	if resp.result == HIT {
		return resp.value
	}
	return nil
}

func (resp *GetCacheResponse) Result() string {
	return resp.result
}

type SetCacheResponse struct {
	value []byte
}

func (resp *SetCacheResponse) StringValue() string {
	return string(resp.value)
}

func (resp *SetCacheResponse) ByteValue() []byte {
	return resp.value
}
