package incubating

// TopicValue Base type for possible messages received on a topic
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
