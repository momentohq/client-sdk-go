package incubating

type topicResponseTypes string

const (
	stringType topicResponseTypes = "string"
	byteType   topicResponseTypes = "bytes"
)

// TopicMessage Base type for possible messages received on a topic
type TopicMessage struct {
	responseType topicResponseTypes
	value        []byte
}

type TopicMessageByte struct {
	value []byte
}

func (t *TopicMessageByte) Value() []byte {
	return t.value
}

type TopicMessageStringResponse struct {
	value string
}

func (t *TopicMessageStringResponse) Value() string {
	return t.value
}

// IsByteMessage returns true if the message received on a topic was of type bytes
func (r *TopicMessage) IsByteMessage() bool {
	return r.responseType == byteType
}

// AsByteMessage returns TopicMessageByte pointer if message is of type Bytes item otherwise returns nil
func (r *TopicMessage) AsByteMessage() *TopicMessageByte {
	if r.IsByteMessage() {
		return &TopicMessageByte{
			value: r.value,
		}
	}
	return nil
}

func (r *TopicMessage) IsStringMessage() bool {
	return r.responseType == stringType
}
func (r *TopicMessage) AsStringMessage() *TopicMessageStringResponse {
	if r.IsStringMessage() {
		return &TopicMessageStringResponse{
			value: string(r.value),
		}
	}
	return nil
}
