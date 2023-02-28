package momento

// For now, topics take the same types as normal methods.
// In the future this might not be so, topics might take
// different types than methods.
//
// The TopicValue interface alias future proofs us. Topics
// still take the normal String and Bytes, but in the future
// we can add types that other methods might not take.
type TopicValue interface {
	isTopicValue()
}

func (String) isTopicValue() {}

func (Bytes) isTopicValue() {}
