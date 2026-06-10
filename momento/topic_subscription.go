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
//     and is safe to call more than once. Close cancels the stream context;
//     an in-flight Item/Event call returns an error — immediately if blocked
//     receiving, or at the next retry check if a reconnect is in flight
//     (Close also interrupts a reconnect backoff wait).
//   - An error from the ctx passed to Item/Event means only that the call
//     stopped waiting: the subscription stays active and a later call with
//     a live ctx resumes it, including resuming an interrupted reconnect.
//     A call blocked receiving notices the ctx at the next
//     message/heartbeat boundary or stream event.
//   - Terminal errors (Close was called, or the ctx passed to Subscribe
//     ended) carry the CanceledError code and unwrap to the underlying
//     context error; check your own ctx first to tell the two apart. The
//     subscription holds a pool slot while it has a live stream; during a
//     reconnect the slot is released and re-acquired, and an abandoned
//     subscription holds its slot until Close.
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

// errEventCtxInterruptedReconnect marks attemptReconnect aborts caused by the
// caller's per-call ctx dying. Event converts exactly these into a paused
// recovery that the next call resumes; it wraps context.Canceled so the error
// stays errors.Is-compatible on the terminal edges where it can escape.
var errEventCtxInterruptedReconnect = fmt.Errorf("event context canceled during reconnect: %w", context.Canceled)

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
	// streamParentCtx is the ctx the caller passed to Subscribe. Reconnected
	// streams are parented to it so their lifetime matches the original
	// stream's, not the ctx of whichever Event call triggered the reconnect.
	streamParentCtx context.Context
	// needsReconnect and lastStreamErr carry an interrupted recovery across
	// Event calls: when the caller's ctx dies after a stream failure, the
	// next call resumes the reconnect instead of the subscription dying.
	// Only the Event goroutine touches them (see the concurrency contract).
	needsReconnect bool
	lastStreamErr  error
	// closed is set by Close. attemptReconnect bails out if it observes this set.
	closed atomic.Bool
	// closedSignal wakes a reconnect backoff wait when Close is called. Only
	// the Close that wins the closed CAS closes it.
	closedSignal chan struct{}
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
		// A previous call's ctx may have died mid-recovery; resume the
		// reconnect before reading so the subscription survives.
		if s.needsReconnect {
			if reconnectErr := s.resumeReconnect(ctx); reconnectErr != nil {
				return nil, reconnectErr
			}
			continue // re-snapshot the freshly stored state
		}

		// Snapshot state once per iteration so every read in this loop body
		// references the same Reservation/cancelContext/cancelFunction triple.
		// attemptReconnect (called below, same goroutine) Stores a new state
		// that the next iteration picks up; the snapshot also keeps a
		// concurrent Close, which Loads and cancels, from seeing a torn view.
		state := s.state.Load()

		// Callers can cancel ctx right after subscribing; check before Recv.
		select {
		case <-ctx.Done():
			// If the stream itself is already dead (Close or the Subscribe
			// ctx ended), prefer the terminal path so the slot is released
			// even when both contexts died together.
			select {
			case <-state.cancelContext.Done():
				return nil, s.terminate(state)
			default:
			}
			// The caller stopped waiting; the subscription and its slot stay
			// intact, and the next Item/Event call picks up where we left off.
			return nil, ctx.Err()
		case <-state.cancelContext.Done():
			return nil, s.terminate(state)
		default:
		}

		rawMsg, err := state.subscribeClient.Recv()
		if err != nil {
			select {
			case <-ctx.Done():
				{
					// The stream failed while the caller's ctx is dead: tear
					// the dead stream down now (its slot must not stay
					// counted), and let the next call resume via reconnect.
					state.cancelFunction()
					decremented := state.reservation.Release()
					s.needsReconnect = true
					s.lastStreamErr = err
					s.log.Debug(
						"[Event RecvMsg] Context done, reconnect deferred to the next call, number of active streams on current grpc channel: %d",
						decremented,
					)
					return nil, ctx.Err()
				}
			case <-state.cancelContext.Done():
				{
					return nil, s.terminate(state)
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

					streamErr := err
					reconnectErr := s.attemptReconnect(ctx, streamErr)
					if reconnectErr != nil {
						return nil, s.reconnectFailure(ctx, streamErr, reconnectErr)
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

// terminate releases the slot of a subscription whose stream context has
// ended (Close, or the Subscribe-time ctx died) and returns the typed
// terminal error. It wraps the context error so errors.Is still matches,
// while the CanceledError code lets callers distinguish a dead subscription
// from their own ctx expiring.
func (s *topicSubscription) terminate(state *streamState) error {
	decremented := state.reservation.Release()
	s.log.Debug(
		"[Event] Subscription stream context ended, number of active streams on current grpc channel: %d",
		decremented,
	)
	return momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "subscription ended", state.cancelContext.Err())
}

// resumeReconnect re-enters an interrupted recovery at the top of Event.
// Returns nil when the stream is re-established (needsReconnect cleared).
func (s *topicSubscription) resumeReconnect(ctx context.Context) error {
	if ctx.Err() != nil {
		// Still no live ctx; stay paused.
		return ctx.Err()
	}
	if reconnectErr := s.attemptReconnect(ctx, s.lastStreamErr); reconnectErr != nil {
		return s.reconnectFailure(ctx, s.lastStreamErr, reconnectErr)
	}
	s.needsReconnect = false
	s.lastStreamErr = nil
	return nil
}

// reconnectFailure classifies a failed attemptReconnect: a dead caller ctx
// pauses the recovery (the next call resumes it), anything else surfaces with
// its typed code intact.
func (s *topicSubscription) reconnectFailure(ctx context.Context, streamErr error, reconnectErr error) error {
	// Pause only when the reconnect was interrupted by the caller's ctx; a
	// give-up or terminal verdict that merely coincides with a dead ctx must
	// surface as-is rather than being resurrected by the next call.
	if errors.Is(reconnectErr, errEventCtxInterruptedReconnect) && !s.closed.Load() && s.streamParentCtx.Err() == nil {
		s.needsReconnect = true
		s.lastStreamErr = streamErr
		return ctx.Err()
	}
	// attemptReconnect returns typed MomentoSvcErrs (e.g. CanceledError on
	// Close); ConvertSvcErr would demote them to ClientSdkError since they
	// carry no gRPC status, so pass them through unchanged.
	var svcErr momentoerrors.MomentoSvcErr
	if errors.As(reconnectErr, &svcErr) {
		return svcErr
	}
	return momentoerrors.ConvertSvcErr(reconnectErr)
}

func (s *topicSubscription) onTopicEvent(method string, event middleware.TopicSubscriptionEventType) {
	if s.topicEventCallback != nil {
		s.topicEventCallback(s.cacheName, method, event)
	}
}

// waitBackoff sleeps for the retry backoff, waking early when the caller's
// ctx dies, the Subscribe-time ctx dies, or Close is called. The timer is
// stopped on early exit so interrupted waits don't leave timers running
// during reconnect storms with long backoffs.
func (s *topicSubscription) waitBackoff(ctx context.Context, backoff time.Duration) error {
	timer := time.NewTimer(backoff)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		// A terminal signal that fired in the same instant wins the tie so
		// the error carries the terminal cause. (reconnectFailure re-checks
		// the terminal conditions either way, so classification never
		// depends on this select.)
		select {
		case <-s.streamParentCtx.Done():
			return momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "subscribe context canceled during reconnect", s.streamParentCtx.Err())
		default:
		}
		select {
		case <-s.closedSignal:
			return momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "subscription closed", context.Canceled)
		default:
		}
		return momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "event context canceled during reconnect", errEventCtxInterruptedReconnect)
	case <-s.streamParentCtx.Done():
		return momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "subscribe context canceled during reconnect", s.streamParentCtx.Err())
	case <-s.closedSignal:
		s.log.Info("Subscription closed during reconnect backoff; aborting retry loop")
		return momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "subscription closed", context.Canceled)
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
			return momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "subscription closed", context.Canceled)
		}
		// A dead streamParentCtx is terminal (a stream parented to it can
		// never come up) and is checked before the caller's ctx so a tie
		// carries the terminal cause; a dead caller ctx pauses recovery
		// (Event resumes it on the next call).
		if s.streamParentCtx.Err() != nil {
			s.log.Info("Subscribe context canceled during reconnect; aborting retry loop")
			return momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "subscribe context canceled during reconnect", s.streamParentCtx.Err())
		}
		if ctx.Err() != nil {
			s.log.Info("Context canceled during reconnect; aborting retry loop")
			return momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "event context canceled during reconnect", errEventCtxInterruptedReconnect)
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
			if waitErr := s.waitBackoff(ctx, time.Duration(*retryBackoffTime)*time.Millisecond); waitErr != nil {
				return waitErr
			}
		}

		s.log.Info("Attempting reconnecting to client stream")
		// Parent the replacement stream to the Subscribe-time ctx, not this
		// Event call's ctx, so the stream's lifetime is stable across Events.
		newState, reconnectErr := s.momentoTopicClient.topicSubscribe(s.streamParentCtx, &TopicSubscribeRequest{
			CacheName:                   s.cacheName,
			TopicName:                   s.topicName,
			ResumeAtTopicSequenceNumber: s.lastKnownSequenceNumber,
			SequencePage:                s.lastKnownSequencePage,
		})

		if reconnectErr != nil {
			// A canceled reconnect can never succeed: the pool is shutting
			// down or the stream's parent context is dead. Stop instead of
			// burning retries (the legacy default strategy retries forever).
			var svcErr momentoerrors.MomentoSvcErr
			if errors.As(reconnectErr, &svcErr) && svcErr.Code() == momentoerrors.CanceledError {
				s.log.Info("Reconnect attempt was canceled; aborting retry loop")
				return reconnectErr
			}
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
			return momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "subscription closed", context.Canceled)
		}
		return nil
	}
}

// Close cancels the stream and releases the slot. Safe to call concurrently
// with Event or another Close; Reservation.Release is idempotent.
func (s *topicSubscription) Close() {
	// Set closed before reading state — see attemptReconnect for the inverse.
	// The CAS doubles as the once-guard for closing the signal channel.
	if s.closed.CompareAndSwap(false, true) {
		close(s.closedSignal)
	}
	state := s.state.Load()
	state.cancelFunction()
	state.reservation.Release()
}
