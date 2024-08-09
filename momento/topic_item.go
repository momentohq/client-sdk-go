package momento

type TopicItem struct {
	Message             TopicValue
	PublisherId         String
	TopicSequenceNumber uint64
}

func (m TopicItem) GetValue() TopicValue {
	return m.Message
}

func (m TopicItem) GetPublisherId() String {
	return m.PublisherId
}

func (m TopicItem) GetTopicSequenceNumber() uint64 {
	return m.TopicSequenceNumber
}
