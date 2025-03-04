package momento

import (
	"context"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type TopicSubscription interface {
	// Item returns only subscription events that contain a string or byte message.
	// Example:
	//
	//	item, err := sub.Item(ctx)
	//	if err != nil {
	//		panic(err)
	//	}
	//	switch msg := item.(type) {
	//	case momento.String:
	//		fmt.Printf("received message as string: '%v'\n", msg)
	//	case momento.Bytes:
	//		fmt.Printf("received message as bytes: '%v'\n", msg)
	//	}
	Item(ctx context.Context) (TopicValue, error)

	// Event returns all possible topics subscription events, such as messages,
	// discontinuities, and heartbeats.
	//
	// Example:
	//
	//	event, err := sub.Event(ctx)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	switch e := event.(type) {
	//	case momento.TopicItem:
	//		fmt.Printf("received item with sequence number %d\n", e.GetTopicSequenceNumber())
	//		fmt.Printf("received item with publisher Id %s\n", e.GetPublisherId())
	//		switch msg := e.GetValue().(type) {
	//		case momento.String:
	//			fmt.Printf("received message as string: '%v'\n", msg)
	//		case momento.Bytes:
	//			fmt.Printf("received message as bytes: '%v'\n", msg)
	//		}
	//	case momento.TopicHeartbeat:
	//		fmt.Printf("received heartbeat\n")
	//	case momento.TopicDiscontinuity:
	//			fmt.Printf("received discontinuity\n")
	//	}
	Event(ctx context.Context) (TopicEvent, error)

	// Close closes the subscription stream.
	Close()
}

type topicSubscription struct {
	topicManager            *grpcmanagers.TopicGrpcManager
	subscribeClient         pb.Pubsub_SubscribeClient
	momentoTopicClient      *pubSubClient
	cacheName               string
	topicName               string
	log                     logger.MomentoLogger
	lastKnownSequenceNumber uint64
	lastKnownSequencePage   uint64
	cancelContext           context.Context
	cancelFunction          context.CancelFunc
}

func (s *topicSubscription) Item(ctx context.Context) (TopicValue, error) {
	for {
		item, err := s.Event(ctx)
		if err != nil {
			return nil, err
		}

		switch item := item.(type) {
		case TopicItem:
			return item.GetValue(), nil
		case TopicHeartbeat:
			continue
		case TopicDiscontinuity:
			continue
		}
	}
}

func (s *topicSubscription) Event(ctx context.Context) (TopicEvent, error) {
	for {
		// Its totally possible a client just calls `cancel` on the `context` immediately after subscribing to an
		// item, so we should check that here.
		select {
		case <-ctx.Done():
			// Context has been canceled, return an error
			s.decrementSubscriptionCount()
			s.log.Debug(
				"[Event] Context done, number of active streams: %d",
				s.momentoTopicClient.countNumberOfActiveSubscriptions(),
			)
			return nil, ctx.Err()
		case <-s.cancelContext.Done():
			// Context has been canceled, return an error
			s.decrementSubscriptionCount()
			s.log.Debug(
				"[Event] Context cancelled, number of active streams: %d",
				s.momentoTopicClient.countNumberOfActiveSubscriptions(),
			)
			return nil, s.cancelContext.Err()
		default:
			// Proceed as is
		}

		rawMsg, err := s.subscribeClient.Recv()
		if err != nil {
			select {
			case <-ctx.Done():
				{
					s.decrementSubscriptionCount()
					s.log.Debug(
						"[Event RecvMsg] Context done, number of active streams: %d",
						s.momentoTopicClient.countNumberOfActiveSubscriptions(),
					)
					return nil, ctx.Err()
				}
			case <-s.cancelContext.Done():
				{
					s.decrementSubscriptionCount()
					s.log.Debug(
						"[Event RecvMsg] Context cancelled, number of active streams: %d",
						s.momentoTopicClient.countNumberOfActiveSubscriptions(),
					)
					return nil, s.cancelContext.Err()
				}
			default:
				{
					// Disconnected, decrement and explicitly close the stream, then attempt to reconnect
					s.log.Error("Stream disconnected due to error: %s", err.Error())
					s.cancelFunction()
					s.decrementSubscriptionCount()
					s.log.Debug(
						"[Event RecvMsg] Default case, attempting to reconnect, number of active streams: %d",
						s.momentoTopicClient.countNumberOfActiveSubscriptions(),
					)
					s.attemptReconnect(ctx)
				}
			}

			// retry getting the latest item
			continue
		}

		switch typedMsg := rawMsg.Kind.(type) {
		case *pb.XSubscriptionItem_Discontinuity:
			s.log.Debug("received discontinuity item: %+v", typedMsg.Discontinuity)
			return NewTopicDiscontinuity(
				typedMsg.Discontinuity.LastTopicSequence,
				typedMsg.Discontinuity.NewTopicSequence,
				typedMsg.Discontinuity.NewSequencePage,
			), nil
		case *pb.XSubscriptionItem_Item:
			s.lastKnownSequenceNumber = typedMsg.Item.GetTopicSequenceNumber()
			s.lastKnownSequencePage = typedMsg.Item.GetSequencePage()
			publisherId := typedMsg.Item.GetPublisherId()

			s.log.Trace(
				"received item with sequence number %d, sequence page %d, and publisher Id %s",
				s.lastKnownSequenceNumber, s.lastKnownSequencePage, publisherId,
			)

			switch subscriptionItem := typedMsg.Item.Value.Kind.(type) {
			case *pb.XTopicValue_Text:
				return NewTopicItem(String(subscriptionItem.Text), String(publisherId), s.lastKnownSequenceNumber, s.lastKnownSequencePage), nil
			case *pb.XTopicValue_Binary:
				return NewTopicItem(Bytes(subscriptionItem.Binary), String(publisherId), s.lastKnownSequenceNumber, s.lastKnownSequencePage), nil
			}
		case *pb.XSubscriptionItem_Heartbeat:
			s.log.Trace("received heartbeat item")
			return TopicHeartbeat{}, nil
		default:
			s.log.Warn("Unrecognized response detected.",
				"response", fmt.Sprint(typedMsg))
			continue
		}
	}
}

func (s *topicSubscription) decrementSubscriptionCount() {
	s.topicManager.NumActiveSubscriptions.Add(-1)
}

func (s *topicSubscription) attemptReconnect(ctx context.Context) {
	// This will attempt to reconnect indefinetly
	reconnectDelay := 500 * time.Millisecond
	for {
		s.log.Info("Attempting reconnecting to client stream")
		time.Sleep(reconnectDelay)
		topicManager, subscribeClient, cancelContext, cancelFunction, err := s.momentoTopicClient.topicSubscribe(ctx, &TopicSubscribeRequest{
			CacheName:                   s.cacheName,
			TopicName:                   s.topicName,
			ResumeAtTopicSequenceNumber: s.lastKnownSequenceNumber,
			SequencePage:                s.lastKnownSequencePage,
		})

		if err != nil {
			s.log.Warn("failed to reconnect to stream, will continue to try in %s milliseconds", fmt.Sprint(reconnectDelay))
		} else {
			s.log.Info("successfully reconnected to subscription stream")
			s.topicManager = topicManager
			s.subscribeClient = subscribeClient
			s.cancelContext = cancelContext
			s.cancelFunction = cancelFunction
			return
		}
	}
}

// Note: number of active subscriptions is decremented in the `Event` method
// for each of the cases when a stream is closed. We do not decrement here
// to avoid double counting.
func (s *topicSubscription) Close() {
	s.cancelFunction()
}
