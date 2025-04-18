package momento

import (
	"github.com/momentohq/client-sdk-go/config/retry"
)

type TopicSubscribeRequest struct {
	CacheName                   string
	TopicName                   string
	ResumeAtTopicSequenceNumber uint64
	SequencePage                uint64
	RetryStrategy               retry.Strategy
}

func (r TopicSubscribeRequest) cacheName() string { return r.CacheName }
