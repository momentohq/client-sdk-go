package momento

type TopicPublishRequest struct {
	CacheName string
	TopicName string
	Value     TopicValue
}

type TopicPublishResponse interface {
	isTopicPublishResponse()
}

type TopicPublishSuccess struct{}

func (TopicPublishSuccess) isTopicPublishResponse() {}

type TopicValue = Value

type TopicValueBytes = Bytes

type TopicValueString = String
