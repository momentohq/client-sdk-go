package responses

// TopicPublishResponse is the base response type for a publish request.
type TopicPublishResponse interface {
	MomentoTopicResponse
	isTopicPublishResponse()
}

// TopicPublishSuccess indicates a successful publish request.
type TopicPublishSuccess struct{}

func (TopicPublishSuccess) isTopicPublishResponse() {}
