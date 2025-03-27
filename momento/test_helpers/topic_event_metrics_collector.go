package helpers

import (
	"fmt"
	"github.com/momentohq/client-sdk-go/config/middleware"
)

type topicEventPayload struct {
	cacheName string
	requestName string
	eventType middleware.TopicSubscriptionEventType
}

type TopicEventCounter struct {
	Heartbeats int
	Discontinuities int
	Items int
	Errors int
	Reconnects int
}

type topicEventMetrics struct {
	data map[string]map[string]TopicEventCounter
}

type TopicEventMetricsCollector interface {
	AddEvent(cacheName string, requestName string, event middleware.TopicSubscriptionEventType)
	GetEventCounter(cacheName string, requestName string) (*TopicEventCounter, error)
}

func NewTopicEventMetricsCollector() TopicEventMetricsCollector {
	return &topicEventMetrics{
		data: make(map[string]map[string]TopicEventCounter),
	}
}

func (t *topicEventMetrics) GetEventCounter(cacheName string, requestName string) (*TopicEventCounter, error) {
	if _, ok := t.data[cacheName]; !ok {
		return nil, fmt.Errorf("cache name '%s' is not valid", cacheName)
	}
	if _, ok := t.data[cacheName][requestName]; !ok {
		return nil, fmt.Errorf("request name '%s' is not valid", requestName)
	}
	counter := t.data[cacheName][requestName]
	return &counter, nil
}

func (t *topicEventMetrics) initializeEventMap(cacheName string, requestName string) {
	if _, ok := t.data[cacheName]; !ok {
		t.data[cacheName] = make(map[string]TopicEventCounter)
	}
	if _, ok := t.data[cacheName][requestName]; !ok {
		t.data[cacheName][requestName] = TopicEventCounter{
			Heartbeats: 0,
			Discontinuities: 0,
			Items: 0,
			Errors: 0,
			Reconnects: 0,
		}
	}
}

func (t *topicEventMetrics) AddEvent(cacheName string, requestName string, event middleware.TopicSubscriptionEventType) {
	t.initializeEventMap(cacheName, requestName)
	eventCounter := t.data[cacheName][requestName]

	switch event {
	case middleware.HEARTBEAT:
		eventCounter.Heartbeats++
	case middleware.DISCONTINUITY:
		eventCounter.Discontinuities++
	case middleware.ITEM:
		eventCounter.Items++
	case middleware.RECONNECT:
		eventCounter.Reconnects++
	case middleware.ERROR:
		eventCounter.Errors++
	}

	t.data[cacheName][requestName] = eventCounter
}

