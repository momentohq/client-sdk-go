package momento

type TopicSubscribeRequest struct {
	CacheName string
	TopicName string
	ResumeAtTopicSequenceNumber uint64
}

func (r TopicSubscribeRequest) cacheName() string { return r.CacheName }
