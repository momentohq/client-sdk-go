package momento

type TopicItem struct {
	Message             Value
	PublisherId         String
	TopicSequenceNumber uint64
}

func (m TopicItem) GetValue() Value {
	return m.Message
}

func (m TopicItem) GetValueString() string {
	return m.Message.asString()
}

func (m TopicItem) GetValueBytes() []byte {
	return m.Message.asBytes()
}

func (m TopicItem) GetPublisherId() String {
	return m.PublisherId
}

func (m TopicItem) GetTopicSequenceNumber() uint64 {
	return m.TopicSequenceNumber
}
