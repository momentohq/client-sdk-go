package incubating

type CreateTopicRequest struct {
	// string used to create a topic.
	TopicName string
}

type TopicSubscribeRequest struct {
	// string used to create a topic.
	TopicName string
}

type TopicPublishRequest struct {
	TopicName string
	Value     string // TODO think about string vs byte more
}
