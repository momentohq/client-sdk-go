package incubating

// TopicValue Base type for possible messages received on a topic
type TopicValue interface {
	isTopicValue()
}

type TopicValueBytes struct {
	Bytes []byte
}

func (_ TopicValueBytes) isTopicValue() {}

type TopicValueString struct {
	Text string
}

func (_ TopicValueString) isTopicValue() {}
