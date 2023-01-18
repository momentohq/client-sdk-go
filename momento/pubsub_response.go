package momento

type CreateTopicResponse struct{}

type TopicMessageReceiveResponse struct {
	// TODO talk about user experience for bytes/strings
	value string
}

// StringValue Decodes and returns byte value sent in topic to string.
func (resp *TopicMessageReceiveResponse) StringValue() string {
	return resp.value
}
