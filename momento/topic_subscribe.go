package momento

type TopicSubscribeRequest struct {
	CacheName string
	TopicName string
}

func (r TopicSubscribeRequest) cacheName() string { return r.CacheName }
