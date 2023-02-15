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

type TopicValue interface {
	isTopicValue()
}

type TopicValueBytes struct {
	Bytes []byte
}

func (TopicValueBytes) isTopicValue() {}

type TopicValueString struct {
	Text string
}

func (TopicValueString) isTopicValue() {}
