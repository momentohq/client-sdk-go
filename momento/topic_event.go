package momento

// TopicEvent is an interface that represents all possible
// events that can be received from a topic subscription.
// This includes messages, heartbeats, and discontinuities.
type TopicEvent interface {
	isTopicEvent()
}

type TopicHeartbeat struct{}

func (TopicHeartbeat) isTopicEvent() {}

type TopicDiscontinuity struct {
	lastKnownSequenceNumber uint64
	newSequenceNumber       uint64
}

func (d TopicDiscontinuity) GetLastKnownSequenceNumber() uint64 {
	return d.lastKnownSequenceNumber
}

func (d TopicDiscontinuity) GetNewSequenceNumber() uint64 {
	return d.newSequenceNumber
}

func NewTopicDiscontinuity(lastKnownSequenceNumber uint64, newSequenceNumber uint64) TopicDiscontinuity {
	return TopicDiscontinuity{
		lastKnownSequenceNumber: lastKnownSequenceNumber,
		newSequenceNumber:       newSequenceNumber,
	}
}

func (TopicDiscontinuity) isTopicEvent() {}

type TopicItem struct {
	message             TopicValue
	publisherId         String
	topicSequenceNumber uint64
}

func (m TopicItem) isTopicEvent() {}

func (m TopicItem) GetValue() TopicValue {
	return m.message
}

func (m TopicItem) GetPublisherId() String {
	return m.publisherId
}

func (m TopicItem) GetTopicSequenceNumber() uint64 {
	return m.topicSequenceNumber
}

func NewTopicItem(message TopicValue, publisherId String, topicSequenceNumber uint64) TopicItem {
	return TopicItem{
		message:             message,
		publisherId:         publisherId,
		topicSequenceNumber: topicSequenceNumber,
	}
}
