package responses

type TopicPublishResponse interface {
	isTopicPublishResponse()
}

type TopicPublishSuccess struct{}

func (TopicPublishSuccess) isTopicPublishResponse() {}
