package momento

type TopicPublishRequest struct {
	CacheName string
	TopicName string
	Value     TopicValue
}
