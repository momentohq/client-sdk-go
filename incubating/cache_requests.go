package incubating

type TopicSubscribeRequest struct {
	// string used to create a topic.
	CacheName string
	TopicName string
}

type TopicPublishRequest struct {
	CacheName string
	TopicName string
	Value     TopicValue
}
