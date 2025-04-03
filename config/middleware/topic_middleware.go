package middleware

import "github.com/momentohq/client-sdk-go/config/logger"

type topicMiddleware struct {
	logger logger.MomentoLogger
}

type TopicMiddlewareProps struct {
	Logger logger.MomentoLogger
}

type TopicMiddleware interface {
	GetLogger() logger.MomentoLogger
	OnSubscribeMetadata(map[string]string) map[string]string
	OnPublishMetadata(map[string]string) map[string]string
}

func (mw *topicMiddleware) GetLogger() logger.MomentoLogger {
	return mw.logger
}

func (mw *topicMiddleware) OnSubscribeMetadata(map[string]string) map[string]string {
	return nil
}

func (mw *topicMiddleware) OnPublishMetadata(map[string]string) map[string]string {
	return nil
}

func NewTopicMiddleware(props Props) TopicMiddleware {
	if props.Logger == nil {
		props.Logger = logger.NewNoopMomentoLoggerFactory().GetLogger("noop")
	}
	return &topicMiddleware{logger: props.Logger}
}

type TopicSubscriptionEventType string

const (
	HEARTBEAT     TopicSubscriptionEventType = "heartbeat"
	ITEM          TopicSubscriptionEventType = "item"
	DISCONTINUITY TopicSubscriptionEventType = "discontinuity"
	RECONNECT     TopicSubscriptionEventType = "reconnect"
	ERROR         TopicSubscriptionEventType = "error"
)

type TopicEventCallbackMiddleware interface {
	TopicMiddleware
	OnTopicEvent(cacheName string, method string, event TopicSubscriptionEventType)
}
