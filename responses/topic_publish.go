package responses

// TopicPublishResponse is a base response type for a publish request.
type TopicPublishResponse interface {
	isTopicPublishResponse()
}

// TopicPublishSuccess indicates a successful publish request.
type TopicPublishSuccess struct{}

func (TopicPublishSuccess) isTopicPublishResponse() {}
