package incubating

// TopicMessage Base type for possible messages received on a topic
type TopicMessage interface {
	isTopicMsg()
}

type TopicMessageBytes struct {
	Value []byte
}

func (_ TopicMessageBytes) isTopicMsg() {}

type TopicMessageString struct {
	Value string
}

func (_ TopicMessageString) isTopicMsg() {}
