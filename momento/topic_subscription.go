package momento

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"

	"github.com/momentohq/client-sdk-go/config/middleware"
	"github.com/momentohq/client-sdk-go/config/retry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers/topic_manager_lists"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

// TopicSubscription is the consumer-side handle for an active topic subscription.
//
// Concurrency contract:
//   - Item and Event must be called from a single goroutine at a time. The
//     underlying gRPC stream's Recv is not safe for concurrent use, and the
//     subscription tracks per-stream sequence state without locking on the
//     hot path.
//   - Close is safe to call concurrently with an in-flight Item/Event call,
//     and is safe to call more than once. Close unblocks the consumer
//     goroutine by cancelling the stream context; the consumer will then
//     observe a cancelled-context error from Item/Event.
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

// subscribeMethodName is the request name reported to middleware for Subscribe
// events.
const subscribeMethodName = "Subscribe"

// streamState bundles the fields swapped together on every (re)connect. Held
// in an atomic.Pointer so readers see a consistent snapshot.
type streamState struct {
	reservation     *topic_manager_lists.Reservation
	subscribeClient pb.Pubsub_SubscribeClient
	cancelContext   context.Context
	cancelFunction  context.CancelFunc
}

type topicSubscription struct {
	state atomic.Pointer[streamState]

	topicEventCallback      func(cacheName string, requestName string, event middleware.TopicSubscriptionEventType)
	momentoTopicClient      *pubSubClient
	cacheName               string
	topicName               string
	log                     logger.MomentoLogger
	lastKnownSequenceNumber uint64
	lastKnownSequencePage   uint64
	retryStrategy           retry.Strategy
	// closed is set by Close. attemptReconnect bails out if it observes this set.
	closed atomic.Bool
}

// grpcStatusCode returns the gRPC status code for err, unwrapping
// MomentoSvcErr if the error came back through topicSubscribe.
func grpcStatusCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	var svcErr momentoerrors.MomentoSvcErr
	if errors.As(err, &svcErr) {
		if original := svcErr.OriginalErr(); original != nil {
			return status.Code(original)
		}
		return codes.Unknown
	}
	return status.Code(err)
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
		// Snapshot state once per iteration so every read in this loop body
		// references the same Reservation/cancelContext/cancelFunction triple.
		// attemptReconnect (called below, same goroutine) Stores a new state
		// that the next iteration picks up; the snapshot also keeps a
		// concurrent Close, which Loads and cancels, from seeing a torn view.
		state := s.state.Load()

		// Callers can cancel ctx right after subscribing; check before Recv.
		select {
		case <-ctx.Done():
			decremented := state.reservation.Release()
			s.log.Debug(
				"[Event] Context done, number of active streams on current grpc channel: %d",
				decremented,
			)
			return nil, ctx.Err()
		case <-state.cancelContext.Done():
			decremented := state.reservation.Release()
			s.log.Debug(
				"[Event] Context cancelled, number of active streams on current grpc channel: %d",
				decremented,
			)
			return nil, state.cancelContext.Err()
		default:
		}

		rawMsg, err := state.subscribeClient.Recv()
		if err != nil {
			select {
			case <-ctx.Done():
				{
					decremented := state.reservation.Release()
					s.log.Debug(
						"[Event RecvMsg] Context done, number of active streams on current grpc channel: %d",
						decremented,
					)
					return nil, ctx.Err()
				}
			case <-state.cancelContext.Done():
				{
					decremented := state.reservation.Release()
					s.log.Debug(
						"[Event RecvMsg] Context cancelled, number of active streams on current grpc channel: %d",
						decremented,
					)
					return nil, state.cancelContext.Err()
				}
			default:
				{
					s.onTopicEvent(subscribeMethodName, middleware.ERROR)
					// Cancel the dead stream before releasing the slot so the
					// server-side resources tear down promptly.
					s.log.Error("Stream disconnected due to error: %s", err.Error())
					state.cancelFunction()
					decremented := state.reservation.Release()
					s.log.Debug(
						"[Event RecvMsg] Default case, attempting to reconnect, number of active streams on current grpc channel: %d",
						decremented,
					)

					err := s.attemptReconnect(ctx, err)
					if err != nil {
						return nil, momentoerrors.ConvertSvcErr(err)
					}
				}
			}

			continue
		}

		switch typedMsg := rawMsg.Kind.(type) {
		case *pb.XSubscriptionItem_Discontinuity:
			s.log.Debug("received discontinuity item: %+v", typedMsg.Discontinuity)
			s.onTopicEvent(subscribeMethodName, middleware.DISCONTINUITY)
			return NewTopicDiscontinuity(
				typedMsg.Discontinuity.LastTopicSequence,
				typedMsg.Discontinuity.NewTopicSequence,
				typedMsg.Discontinuity.NewSequencePage,
			), nil
		case *pb.XSubscriptionItem_Item:
			s.onTopicEvent(subscribeMethodName, middleware.ITEM)
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
			s.onTopicEvent(subscribeMethodName, middleware.HEARTBEAT)
			return TopicHeartbeat{}, nil
		default:
			s.log.Warn("Unrecognized response detected.",
				"response", fmt.Sprint(typedMsg))
			continue
		}
	}
}

func (s *topicSubscription) onTopicEvent(method string, event middleware.TopicSubscriptionEventType) {
	if s.topicEventCallback != nil {
		s.topicEventCallback(s.cacheName, method, event)
	}
}

func (s *topicSubscription) attemptReconnect(ctx context.Context, err error) error {
	if s.retryStrategy == nil {
		s.log.Info("No retry strategy provided, returning error")
		return err
	}
	attempt := 1
	for {
		if s.closed.Load() {
			s.log.Info("Subscription closed during reconnect; aborting retry loop")
			return momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "subscription closed", nil)
		}
		// A reconnected stream would be a child of ctx, so once ctx is dead a
		// reconnect can never succeed; stop instead of retrying forever.
		if ctx.Err() != nil {
			s.log.Info("Context canceled during reconnect; aborting retry loop")
			return momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "subscribe context canceled during reconnect", ctx.Err())
		}

		retryBackoffTime := s.retryStrategy.DetermineWhenToRetry(retry.StrategyProps{
			GrpcStatusCode: grpcStatusCode(err),
			GrpcMethod:     "/cache_client.pubsub.Pubsub/Subscribe",
			AttemptNumber:  attempt,
		})

		if retryBackoffTime == nil {
			s.log.Warn("Retry strategy determined that we should not retry, returning error")
			return err
		}

		s.onTopicEvent(subscribeMethodName, middleware.RECONNECT)

		if *retryBackoffTime > 0 {
			s.log.Info("Waiting %s milliseconds before attempting to reconnect", fmt.Sprint(*retryBackoffTime))
			select {
			case <-time.After(time.Duration(*retryBackoffTime) * time.Millisecond):
			case <-ctx.Done():
				return momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "subscribe context canceled during reconnect", ctx.Err())
			}
		}

		s.log.Info("Attempting reconnecting to client stream")
		newState, reconnectErr := s.momentoTopicClient.topicSubscribe(ctx, &TopicSubscribeRequest{
			CacheName:                   s.cacheName,
			TopicName:                   s.topicName,
			ResumeAtTopicSequenceNumber: s.lastKnownSequenceNumber,
			SequencePage:                s.lastKnownSequencePage,
		})

		if reconnectErr != nil {
			s.log.Warn("Failed to reconnect to stream: %s", reconnectErr.Error())
			err = reconnectErr
			attempt++
			continue
		}

		s.log.Info("Successfully reconnected to subscription stream")
		s.state.Store(newState)

		// Store before reading closed; Close stores closed before reading
		// state. Either side observes the other's write, so a Close racing
		// this reconnect always tears down the new stream.
		if s.closed.Load() {
			s.log.Info("Subscription closed during reconnect; tearing down newly established stream")
			newState.cancelFunction()
			newState.reservation.Release()
			return momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "subscription closed", nil)
		}
		return nil
	}
}

// Close cancels the stream and releases the slot. Safe to call concurrently
// with Event or another Close; Reservation.Release is idempotent.
func (s *topicSubscription) Close() {
	// Set closed before reading state — see attemptReconnect for the inverse.
	s.closed.Store(true)
	state := s.state.Load()
	state.cancelFunction()
	state.reservation.Release()
}
