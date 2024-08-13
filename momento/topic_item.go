package momento

type DetailedTopicItem interface {
	isDetailedTopicItem()
}

// type TopicHeartbeat struct{}

// func (TopicHeartbeat) isDetailedTopicItem() {}

// type TopicDiscontinuity struct {
// 	lastKnownSequenceNumber uint64
// 	newSequenceNumber       uint64
// }

// func (d TopicDiscontinuity) GetLastKnownSequenceNumber() uint64 {
// 	return d.lastKnownSequenceNumber
// }

// func (d TopicDiscontinuity) GetNewSequenceNumber() uint64 {
// 	return d.newSequenceNumber
// }

// func NewTopicDiscontinuity(lastKnownSequenceNumber uint64, newSequenceNumber uint64) TopicDiscontinuity {
// 	return TopicDiscontinuity{
// 		lastKnownSequenceNumber: lastKnownSequenceNumber,
// 		newSequenceNumber:       newSequenceNumber,
// 	}
// }

// func (TopicDiscontinuity) isDetailedTopicItem() {}

type TopicMessage struct {
	message             TopicValue
	publisherId         String
	topicSequenceNumber uint64
}

func (m TopicMessage) isDetailedTopicItem() {}

func (m TopicMessage) GetValue() TopicValue {
	return m.message
}

func (m TopicMessage) GetPublisherId() String {
	return m.publisherId
}

func (m TopicMessage) GetTopicSequenceNumber() uint64 {
	return m.topicSequenceNumber
}

func NewTopicMessage(message TopicValue, publisherId String, topicSequenceNumber uint64) TopicMessage {
	return TopicMessage{
		message:             message,
		publisherId:         publisherId,
		topicSequenceNumber: topicSequenceNumber,
	}
}
