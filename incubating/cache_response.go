package incubating

type topicResponseTypes string

// TopicMessage Base type for possible messages received on a topic
type TopicMessage struct {
	value interface{}
}

type TopicMessageBytes struct {
	value []byte
}

func (t *TopicMessageBytes) Value() []byte {
	return t.value
}

type TopicMessageString struct {
	value string
}

func (t *TopicMessageString) Value() string {
	return t.value
}

// IsByteMessage returns true if the message received on a topic was of type bytes
func (r *TopicMessage) IsByteMessage() bool {
	if _, ok := r.value.([]byte); ok {
		return true
	}
	return false
}

// AsByteMessage returns TopicMessageBytes pointer if message is of type Bytes item otherwise returns nil
func (r *TopicMessage) AsByteMessage() *TopicMessageBytes {
	if r.IsByteMessage() {
		return &TopicMessageBytes{
			value: r.value.([]byte),
		}
	}
	return nil
}

func (r *TopicMessage) IsStringMessage() bool {
	if _, ok := r.value.([]byte); ok {
		return true
	}
	return false
}
func (r *TopicMessage) AsStringMessage() *TopicMessageString {
	if r.IsStringMessage() {
		return &TopicMessageString{
			value: r.value.(string),
		}
	}
	return nil
}
