package momento

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/retry"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers/topic_manager_lists"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type testTopicGrpcConnectionPool struct {
	manager      *grpcmanagers.TopicGrpcManager
	activeCount  atomic.Int64
	releaseCount atomic.Int64
	// closed makes GetNextTopicGrpcManager fail like a real pool after Close.
	closed atomic.Bool
}

func newTestTopicGrpcConnectionPool(streamClient pb.PubsubClient) *testTopicGrpcConnectionPool {
	return &testTopicGrpcConnectionPool{
		manager: &grpcmanagers.TopicGrpcManager{
			StreamClient: streamClient,
		},
	}
}

func (p *testTopicGrpcConnectionPool) GetNextTopicGrpcManager() (*topic_manager_lists.Reservation, momentoerrors.MomentoSvcErr) {
	if p.closed.Load() {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "connection pool is shutting down", context.Canceled)
	}
	p.manager.NumActiveSubscriptions.Add(1)
	p.activeCount.Add(1)
	return topic_manager_lists.NewReservation(p.manager, p.releaseManager), nil
}

// releaseManager mirrors the real stream pool's release for accounting.
func (p *testTopicGrpcConnectionPool) releaseManager(manager *grpcmanagers.TopicGrpcManager) int64 {
	p.releaseCount.Add(1)
	p.activeCount.Add(-1)
	return manager.NumActiveSubscriptions.Add(-1)
}

func (p *testTopicGrpcConnectionPool) Close() {}

type testPubsubClient struct {
	subscribe func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error)
}

func (c *testPubsubClient) Publish(context.Context, *pb.XPublishRequest, ...grpc.CallOption) (*pb.XEmpty, error) {
	return &pb.XEmpty{}, nil
}

func (c *testPubsubClient) Subscribe(ctx context.Context, request *pb.XSubscriptionRequest, opts ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
	return c.subscribe(ctx, request, opts...)
}

type recvResult struct {
	item *pb.XSubscriptionItem
	err  error
}

type testSubscribeClient struct {
	grpc.ClientStream
	results []recvResult
	recv    func() (*pb.XSubscriptionItem, error)
}

func (c *testSubscribeClient) Recv() (*pb.XSubscriptionItem, error) {
	if c.recv != nil {
		return c.recv()
	}
	if len(c.results) == 0 {
		return nil, status.Error(codes.Unavailable, "test stream exhausted")
	}
	result := c.results[0]
	c.results = c.results[1:]
	return result.item, result.err
}

type alwaysRetryStrategy struct{}

func (alwaysRetryStrategy) DetermineWhenToRetry(retry.StrategyProps) *int {
	delay := 0
	return &delay
}

// giveUpAfterRetryStrategy retries up to `attempts` times, then returns nil
// to signal give-up.
type giveUpAfterRetryStrategy struct {
	attempts int
}

func (g giveUpAfterRetryStrategy) DetermineWhenToRetry(props retry.StrategyProps) *int {
	if props.AttemptNumber > g.attempts {
		return nil
	}
	delay := 0
	return &delay
}

// signalingBackoffRetryStrategy always retries after a fixed backoff and
// closes entered the first time the retry loop commits to that backoff, so
// tests can deterministically race Close/cancel against the backoff wait.
type signalingBackoffRetryStrategy struct {
	backoffMs int
	entered   chan struct{}
	once      sync.Once
}

// firstCallBackoffRetryStrategy returns a long backoff on the first retry
// consultation (closing entered), then zero for all subsequent ones, so a
// paused reconnect can resume promptly.
type firstCallBackoffRetryStrategy struct {
	firstBackoffMs int
	calls          atomic.Int64
	entered        chan struct{}
}

func (s *firstCallBackoffRetryStrategy) DetermineWhenToRetry(retry.StrategyProps) *int {
	if s.calls.Add(1) == 1 {
		close(s.entered)
		return &s.firstBackoffMs
	}
	zero := 0
	return &zero
}

func newSignalingBackoffRetryStrategy(backoffMs int) *signalingBackoffRetryStrategy {
	return &signalingBackoffRetryStrategy{
		backoffMs: backoffMs,
		entered:   make(chan struct{}),
	}
}

func (s *signalingBackoffRetryStrategy) DetermineWhenToRetry(retry.StrategyProps) *int {
	s.once.Do(func() { close(s.entered) })
	return &s.backoffMs
}

func newAccountingTestClient(streamClient pb.PubsubClient, timeout time.Duration) (defaultTopicClient, *testTopicGrpcConnectionPool) {
	return newAccountingTestClientWithStrategy(streamClient, timeout, alwaysRetryStrategy{})
}

func newAccountingTestClientWithStrategy(
	streamClient pb.PubsubClient,
	timeout time.Duration,
	strategy retry.Strategy,
) (defaultTopicClient, *testTopicGrpcConnectionPool) {
	pool := newTestTopicGrpcConnectionPool(streamClient)
	client := defaultTopicClient{
		pubSubClient: &pubSubClient{
			streamGrpcConnectionPool: pool,
		},
		log:            logger.NewNoopMomentoLoggerFactory().GetLogger("topic-subscription-accounting-test"),
		requestTimeout: timeout,
		retryStrategy:  strategy,
	}
	return client, pool
}

func testSubscribeRequest() *TopicSubscribeRequest {
	return &TopicSubscribeRequest{
		CacheName: "cache",
		TopicName: "topic",
	}
}

func heartbeatItem() *pb.XSubscriptionItem {
	return &pb.XSubscriptionItem{
		Kind: &pb.XSubscriptionItem_Heartbeat{
			Heartbeat: &pb.XHeartbeat{},
		},
	}
}

func topicItem(sequenceNumber uint64) *pb.XSubscriptionItem {
	return &pb.XSubscriptionItem{
		Kind: &pb.XSubscriptionItem_Item{
			Item: &pb.XTopicItem{
				TopicSequenceNumber: sequenceNumber,
				SequencePage:        1,
				Value: &pb.XTopicValue{
					Kind: &pb.XTopicValue_Text{Text: "value"},
				},
			},
		},
	}
}

// assertErrorCode fails unless err is a MomentoSvcErr carrying the given code.
func assertErrorCode(t *testing.T, err error, code string) {
	t.Helper()
	svcErr, ok := err.(momentoerrors.MomentoSvcErr)
	if !ok {
		t.Fatalf("error %v (%T) is not a MomentoSvcErr", err, err)
	}
	if svcErr.Code() != code {
		t.Fatalf("error code = %s, want %s (err: %v)", svcErr.Code(), code, err)
	}
}

func discontinuityItem() *pb.XSubscriptionItem {
	return &pb.XSubscriptionItem{
		Kind: &pb.XSubscriptionItem_Discontinuity{
			Discontinuity: &pb.XDiscontinuity{},
		},
	}
}

func assertAccounting(t *testing.T, pool *testTopicGrpcConnectionPool, activeCount int64) {
	t.Helper()
	if got := pool.activeCount.Load(); got != activeCount {
		t.Fatalf("active stream count = %d, want %d", got, activeCount)
	}
	if got := pool.manager.NumActiveSubscriptions.Load(); got != activeCount {
		t.Fatalf("manager subscription count = %d, want %d", got, activeCount)
	}
}

func TestTopicSubscribeReleasesPoolCountersWhenStreamOpenFails(t *testing.T) {
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return nil, status.Error(codes.Unavailable, "stream open failed")
		},
	}
	_, pool := newAccountingTestClient(streamClient, time.Second)
	client := &pubSubClient{streamGrpcConnectionPool: pool}

	cancelCtx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	_, err := client.topicSubscribe(cancelCtx, cancelFn, testSubscribeRequest())
	if err == nil {
		t.Fatal("expected topicSubscribe to return an error")
	}
	assertAccounting(t, pool, 0)
	if got := pool.releaseCount.Load(); got != 1 {
		t.Fatalf("release count = %d, want 1", got)
	}
}

func TestSubscribeReleasesPoolCountersWhenFirstMessageRecvFails(t *testing.T) {
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				results: []recvResult{{err: status.Error(codes.Unavailable, "first message failed")}},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, time.Second)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err == nil {
		t.Fatal("expected Subscribe to return an error")
	}
	assertErrorCode(t, err, momentoerrors.ServerUnavailableError)
	if sub != nil {
		t.Fatal("expected subscription to be nil")
	}
	assertAccounting(t, pool, 0)
}

func TestSubscribeReleasesPoolCountersWhenFirstMessageIsNotHeartbeat(t *testing.T) {
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				results: []recvResult{{item: topicItem(1)}},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, time.Second)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err == nil {
		t.Fatal("expected Subscribe to return an error")
	}
	assertErrorCode(t, err, momentoerrors.InternalServerError)
	if sub != nil {
		t.Fatal("expected subscription to be nil")
	}
	assertAccounting(t, pool, 0)
}

func TestTopicSubscriptionCloseReleasesPoolCountersWithoutEvent(t *testing.T) {
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				results: []recvResult{{item: heartbeatItem()}},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, time.Second)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	sub.Close()
	assertAccounting(t, pool, 0)
}

// TestTopicSubscriptionCloseDuringReconnectTearsDownNewStream covers the
// close-during-reconnect race: Close sets closed=true while attemptReconnect
// holds a new Reservation. The post-Store check must release that new
// Reservation rather than leaking it. The subscribe function is gated so
// Close runs deterministically inside the race window.
func TestTopicSubscriptionCloseDuringReconnectTearsDownNewStream(t *testing.T) {
	reconnectEntered := make(chan struct{})
	unblockReconnect := make(chan struct{})
	var subscribeCalls atomic.Uint64

	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			call := subscribeCalls.Add(1)
			switch call {
			case 1:
				// Handshake heartbeat, then a disconnect to drive Event into
				// attemptReconnect.
				return &testSubscribeClient{
					results: []recvResult{
						{item: heartbeatItem()},
						{err: status.Error(codes.Unavailable, "disconnected")},
					},
				}, nil
			case 2:
				// Reconnect: signal the test that a new Reservation is held,
				// then block so Close can race with this in-flight reconnect.
				close(reconnectEntered)
				<-unblockReconnect
				return &testSubscribeClient{
					results: []recvResult{{item: heartbeatItem()}},
				}, nil
			default:
				return nil, status.Error(codes.Canceled, "test: subscribe should not be called after Close")
			}
		},
	}

	client, pool := newAccountingTestClient(streamClient, time.Second)
	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	eventDone := make(chan error, 1)
	go func() {
		_, err := sub.Event(context.Background())
		eventDone <- err
	}()

	<-reconnectEntered

	// Close while reconnect is mid-flight. Sets closed=true; the old
	// Reservation was already released by Event before attemptReconnect ran.
	sub.Close()

	// Let reconnect proceed. The post-Store check observes closed=true and
	// releases the new Reservation.
	close(unblockReconnect)

	eventErr := <-eventDone
	if eventErr == nil {
		t.Fatal("expected Event to return an error after close-during-reconnect")
	}
	assertErrorCode(t, eventErr, momentoerrors.CanceledError)
	if !errors.Is(eventErr, context.Canceled) {
		t.Fatalf("terminal error should unwrap to context.Canceled, got %v", eventErr)
	}
	assertAccounting(t, pool, 0)
}

// TestManyConcurrentSubscriptionsKeepsPoolBalanced stresses the pool counter
// with parallel Subscribe and Close calls.
func TestManyConcurrentSubscriptionsKeepsPoolBalanced(t *testing.T) {
	const numSubs = 100

	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				results: []recvResult{{item: heartbeatItem()}},
			}, nil
		},
	}
	// Generous timeout: the fake handshake is synchronous, so the handshake
	// timeout must never fire even on an oversubscribed -race CI runner.
	client, pool := newAccountingTestClient(streamClient, 30*time.Second)

	subs := make([]TopicSubscription, numSubs)
	var wg sync.WaitGroup
	wg.Add(numSubs)
	for i := 0; i < numSubs; i++ {
		go func(i int) {
			defer wg.Done()
			sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
			if err != nil {
				t.Errorf("Subscribe %d returned error: %v", i, err)
				return
			}
			subs[i] = sub
		}(i)
	}
	wg.Wait()
	assertAccounting(t, pool, numSubs)

	wg.Add(numSubs)
	for i := 0; i < numSubs; i++ {
		go func(i int) {
			defer wg.Done()
			if subs[i] != nil {
				subs[i].Close()
			}
		}(i)
	}
	wg.Wait()
	assertAccounting(t, pool, 0)
}

// TestConcurrentCloseOnSameSubscription verifies many goroutines calling
// Close on the same subscription release the pool slot exactly once.
func TestConcurrentCloseOnSameSubscription(t *testing.T) {
	const goroutines = 64

	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				results: []recvResult{{item: heartbeatItem()}},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, 30*time.Second)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			<-start
			sub.Close()
		}()
	}
	close(start)
	wg.Wait()

	assertAccounting(t, pool, 0)
	if got := pool.releaseCount.Load(); got != 1 {
		t.Fatalf("releaseCount = %d, want 1 (concurrent Close should release exactly once)", got)
	}
}

func TestTopicSubscriptionReconnectKeepsPoolCountersBalanced(t *testing.T) {
	disconnectErr := status.Error(codes.Unavailable, "stream disconnected")
	var subscribeCalls atomic.Uint64
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			call := subscribeCalls.Add(1)
			if call == 1 {
				return &testSubscribeClient{
					results: []recvResult{
						{item: heartbeatItem()},
						{err: disconnectErr},
					},
				}, nil
			}
			return &testSubscribeClient{
				results: []recvResult{
					{item: topicItem(call)},
					{err: disconnectErr},
				},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, time.Second)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	for i := 0; i < 5; i++ {
		event, err := sub.Event(context.Background())
		if err != nil {
			t.Fatalf("Event %d returned error: %v", i, err)
		}
		if _, ok := event.(TopicItem); !ok {
			t.Fatalf("Event %d = %T, want TopicItem", i, event)
		}
		assertAccounting(t, pool, 1)
	}

	sub.Close()
	assertAccounting(t, pool, 0)
}

// TestTopicSubscriptionReconnectGiveUpReleasesPoolCounters drives
// attemptReconnect to give up after the retry strategy returns nil. Counters
// must be balanced when Event surfaces the error, and a subsequent Close
// must be a no-op on the already-released Reservation.
func TestTopicSubscriptionReconnectGiveUpReleasesPoolCounters(t *testing.T) {
	const allowedAttempts = 3
	disconnectErr := status.Error(codes.Unavailable, "stream disconnected")
	var subscribeCalls atomic.Uint64

	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			call := subscribeCalls.Add(1)
			if call == 1 {
				return &testSubscribeClient{
					results: []recvResult{
						{item: heartbeatItem()},
						{err: disconnectErr},
					},
				}, nil
			}
			// Every reconnect attempt fails so the retry strategy is exhausted.
			return nil, disconnectErr
		},
	}
	client, pool := newAccountingTestClientWithStrategy(
		streamClient,
		time.Second,
		giveUpAfterRetryStrategy{attempts: allowedAttempts},
	)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	_, eventErr := sub.Event(context.Background())
	if eventErr == nil {
		t.Fatal("expected Event to surface the give-up error from attemptReconnect")
	}

	// Original slot released by Event before reconnect; each failed reconnect
	// attempt released its slot on the early-return path: 1 + allowedAttempts.
	assertAccounting(t, pool, 0)
	if got := pool.releaseCount.Load(); got != int64(1+allowedAttempts) {
		t.Fatalf("releaseCount = %d, want %d (one per failed attempt plus the original slot)", got, 1+allowedAttempts)
	}

	// Close after give-up must be a no-op and must not panic.
	sub.Close()
	assertAccounting(t, pool, 0)
}

// TestTopicSubscriptionEventContextErrorKeepsSubscriptionAlive pins the
// non-terminal per-call ctx contract: an expired or canceled Event ctx stops
// only that call. The stream and its pool slot stay intact and a later call
// with a live ctx receives messages, matching the timeout semantics of
// comparable clients (NATS NextMsg, Kafka ReadMessage).
func TestTopicSubscriptionEventContextErrorKeepsSubscriptionAlive(t *testing.T) {
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				results: []recvResult{
					{item: heartbeatItem()},
					{item: topicItem(2)},
				},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, time.Second)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, eventErr := sub.Event(canceledCtx); eventErr == nil {
		t.Fatal("expected Event to return an error for a canceled context")
	}

	// The slot stays held and the stream stays alive.
	assertAccounting(t, pool, 1)
	state := sub.(*topicSubscription).state.Load()
	if state.cancelContext.Err() != nil {
		t.Fatal("a per-call ctx error must not cancel the stream context")
	}

	// A later call with a live ctx picks up where we left off.
	event, eventErr := sub.Event(context.Background())
	if eventErr != nil {
		t.Fatalf("Event after a per-call ctx error returned error: %v", eventErr)
	}
	if _, ok := event.(TopicItem); !ok {
		t.Fatalf("Event = %T, want TopicItem from the surviving stream", event)
	}
	assertAccounting(t, pool, 1)

	sub.Close()
	assertAccounting(t, pool, 0)
}

// TestTopicSubscriptionEventAfterCloseReturnsTypedCanceled pins the terminal
// error contract: after Close, Event returns a CanceledError that still
// unwraps to context.Canceled, so callers can distinguish a dead subscription
// from their own ctx expiring.
func TestTopicSubscriptionEventAfterCloseReturnsTypedCanceled(t *testing.T) {
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				results: []recvResult{{item: heartbeatItem()}},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, time.Second)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	sub.Close()
	assertAccounting(t, pool, 0)

	_, eventErr := sub.Event(context.Background())
	if eventErr == nil {
		t.Fatal("expected Event after Close to return an error")
	}
	assertErrorCode(t, eventErr, momentoerrors.CanceledError)
	if !errors.Is(eventErr, context.Canceled) {
		t.Fatalf("terminal error should unwrap to context.Canceled, got %v", eventErr)
	}
}

// TestTopicSubscriptionSubscribeCtxCancelEndsSubscription pins the other
// terminal flavor: cancelling the ctx passed to Subscribe tears down an
// established subscription. The blocked-Recv fake also pins stream lineage —
// the stream context must be a child of the Subscribe ctx, so cancellation
// unblocks a waiting Event, which returns the typed terminal error and
// releases the slot.
func TestTopicSubscriptionSubscribeCtxCancelEndsSubscription(t *testing.T) {
	recvBlocked := make(chan struct{})
	var recvBlockedOnce sync.Once
	streamClient := &testPubsubClient{
		subscribe: func(ctx context.Context, _ *pb.XSubscriptionRequest, _ ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			heartbeatSent := false
			return &testSubscribeClient{
				recv: func() (*pb.XSubscriptionItem, error) {
					if !heartbeatSent {
						heartbeatSent = true
						return heartbeatItem(), nil
					}
					recvBlockedOnce.Do(func() { close(recvBlocked) })
					<-ctx.Done()
					return nil, status.FromContextError(ctx.Err()).Err()
				},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, time.Second)

	subscribeCtx, cancelSubscribeCtx := context.WithCancel(context.Background())
	defer cancelSubscribeCtx()
	sub, err := client.Subscribe(subscribeCtx, testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	eventDone := make(chan error, 1)
	go func() {
		_, eventErr := sub.Event(context.Background())
		eventDone <- eventErr
	}()

	// Cancel the Subscribe-time ctx while Event is blocked in Recv. The
	// stream context is its child, so the blocked call must return.
	<-recvBlocked
	cancelSubscribeCtx()

	state := sub.(*topicSubscription).state.Load()
	if state.cancelContext.Err() == nil {
		t.Fatal("the stream context must be a child of the Subscribe ctx")
	}

	select {
	case eventErr := <-eventDone:
		if eventErr == nil {
			t.Fatal("expected Event to return an error after Subscribe-ctx cancellation")
		}
		assertErrorCode(t, eventErr, momentoerrors.CanceledError)
		if !errors.Is(eventErr, context.Canceled) {
			t.Fatalf("terminal error should unwrap to context.Canceled, got %v", eventErr)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Event did not return within 5s of Subscribe-ctx cancellation")
	}
	assertAccounting(t, pool, 0)
}

// TestTopicSubscriptionItemSkipsNonMessageEvents pins Item's contract: it
// swallows heartbeats and discontinuities and returns the next message value.
func TestTopicSubscriptionItemSkipsNonMessageEvents(t *testing.T) {
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				results: []recvResult{
					{item: heartbeatItem()}, // handshake
					{item: heartbeatItem()},
					{item: discontinuityItem()},
					{item: topicItem(2)},
				},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, time.Second)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}

	value, err := sub.Item(context.Background())
	if err != nil {
		t.Fatalf("Item returned error: %v", err)
	}
	if got, ok := value.(String); !ok || got != String("value") {
		t.Fatalf("Item = %v (%T), want String(\"value\")", value, value)
	}

	sub.Close()
	assertAccounting(t, pool, 0)
}

// TestTopicSubscriptionSubscribeCtxCancelDuringBackoffIsTerminal pins the
// asymmetry between the two context flavors mid-backoff: the Subscribe-time
// ctx dying ends the subscription (typed CanceledError, no pause), unlike the
// per-call ctx, which pauses recovery for the next call.
func TestTopicSubscriptionSubscribeCtxCancelDuringBackoffIsTerminal(t *testing.T) {
	disconnectErr := status.Error(codes.Unavailable, "stream disconnected")
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				results: []recvResult{
					{item: heartbeatItem()},
					{err: disconnectErr},
				},
			}, nil
		},
	}
	strategy := newSignalingBackoffRetryStrategy(10 * 60 * 1000)
	client, pool := newAccountingTestClientWithStrategy(streamClient, time.Second, strategy)

	subscribeCtx, cancelSubscribeCtx := context.WithCancel(context.Background())
	defer cancelSubscribeCtx()
	sub, err := client.Subscribe(subscribeCtx, testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	eventDone := make(chan error, 1)
	go func() {
		_, eventErr := sub.Event(context.Background())
		eventDone <- eventErr
	}()

	<-strategy.entered
	cancelSubscribeCtx()

	select {
	case eventErr := <-eventDone:
		if eventErr == nil {
			t.Fatal("expected Event to return an error after Subscribe-ctx cancellation")
		}
		assertErrorCode(t, eventErr, momentoerrors.CanceledError)
		if !errors.Is(eventErr, context.Canceled) {
			t.Fatalf("terminal error should unwrap to context.Canceled, got %v", eventErr)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Event did not return within 5s of Subscribe-ctx cancellation")
	}
	assertAccounting(t, pool, 0)

	// Terminal, not paused: a later call must report the dead subscription,
	// not retry the reconnect.
	_, eventErr := sub.Event(context.Background())
	if eventErr == nil {
		t.Fatal("expected the subscription to stay terminal")
	}
	assertErrorCode(t, eventErr, momentoerrors.CanceledError)
}

// TestTopicSubscriptionReconnectPausesOnContextCancelAndResumes drives the
// retry loop into a long backoff, cancels the caller's ctx (Event must return
// promptly), then verifies the next call with a live ctx resumes the
// interrupted reconnect instead of the subscription dying.
func TestTopicSubscriptionReconnectPausesOnContextCancelAndResumes(t *testing.T) {
	disconnectErr := status.Error(codes.Unavailable, "stream disconnected")
	var subscribeCalls atomic.Uint64
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			if subscribeCalls.Add(1) == 1 {
				return &testSubscribeClient{
					results: []recvResult{
						{item: heartbeatItem()},
						{err: disconnectErr},
					},
				}, nil
			}
			return &testSubscribeClient{
				results: []recvResult{{item: topicItem(2)}},
			}, nil
		},
	}
	strategy := &firstCallBackoffRetryStrategy{
		firstBackoffMs: 10 * 60 * 1000,
		entered:        make(chan struct{}),
	}
	client, pool := newAccountingTestClientWithStrategy(streamClient, time.Second, strategy)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	eventCtx, cancelEventCtx := context.WithCancel(context.Background())
	eventDone := make(chan error, 1)
	go func() {
		_, eventErr := sub.Event(eventCtx)
		eventDone <- eventErr
	}()

	// Wait until the retry loop has committed to the 10-minute backoff, then
	// cancel. Event must return promptly; without the ctx-aware wait this
	// would block ~10 minutes.
	<-strategy.entered
	cancelEventCtx()

	select {
	case eventErr := <-eventDone:
		if eventErr == nil {
			t.Fatal("expected Event to return an error after context cancellation")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Event did not return within 5s of context cancellation")
	}
	// The failed stream's slot was released; recovery is paused, not dead.
	assertAccounting(t, pool, 0)

	// A call with another dead ctx stays paused rather than resuming.
	deadCtx, cancelDead := context.WithCancel(context.Background())
	cancelDead()
	if _, pausedErr := sub.Event(deadCtx); pausedErr == nil {
		t.Fatal("expected a dead-ctx call on a paused subscription to return an error")
	}
	assertAccounting(t, pool, 0)

	// A live ctx resumes the reconnect (zero backoff on the second
	// consultation) and delivers from the new stream.
	event, eventErr := sub.Event(context.Background())
	if eventErr != nil {
		t.Fatalf("Event after paused reconnect returned error: %v", eventErr)
	}
	if _, ok := event.(TopicItem); !ok {
		t.Fatalf("Event = %T, want TopicItem from the resumed stream", event)
	}
	assertAccounting(t, pool, 1)

	sub.Close()
	assertAccounting(t, pool, 0)
}

// TestTopicSubscriptionStreamErrorWithDeadContextResumesNextCall covers the
// other pause entry: the stream fails while the caller's ctx is already dead.
// The dead stream's slot is released immediately and the next call with a
// live ctx reconnects.
func TestTopicSubscriptionStreamErrorWithDeadContextResumesNextCall(t *testing.T) {
	disconnectErr := status.Error(codes.Unavailable, "stream disconnected")
	gate := make(chan struct{})
	recvBlocked := make(chan struct{})
	var recvBlockedOnce sync.Once
	var subscribeCalls atomic.Uint64
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			if subscribeCalls.Add(1) == 1 {
				heartbeatSent := false
				return &testSubscribeClient{
					recv: func() (*pb.XSubscriptionItem, error) {
						if !heartbeatSent {
							heartbeatSent = true
							return heartbeatItem(), nil
						}
						recvBlockedOnce.Do(func() { close(recvBlocked) })
						<-gate
						return nil, disconnectErr
					},
				}, nil
			}
			return &testSubscribeClient{
				results: []recvResult{{item: topicItem(2)}},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, time.Second)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	eventCtx, cancelEventCtx := context.WithCancel(context.Background())
	eventDone := make(chan error, 1)
	go func() {
		_, eventErr := sub.Event(eventCtx)
		eventDone <- eventErr
	}()

	// Kill the caller's ctx while Event is blocked in Recv, then fail the
	// stream: Event observes the dead ctx on the error path.
	<-recvBlocked
	cancelEventCtx()
	close(gate)

	select {
	case eventErr := <-eventDone:
		if eventErr == nil {
			t.Fatal("expected Event to return an error")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Event did not return within 5s")
	}
	assertAccounting(t, pool, 0)

	event, eventErr := sub.Event(context.Background())
	if eventErr != nil {
		t.Fatalf("Event after deferred reconnect returned error: %v", eventErr)
	}
	if _, ok := event.(TopicItem); !ok {
		t.Fatalf("Event = %T, want TopicItem from the reconnected stream", event)
	}
	assertAccounting(t, pool, 1)

	sub.Close()
	assertAccounting(t, pool, 0)
}

// TestTopicSubscriptionCloseInterruptsReconnectBackoff verifies Close wakes a
// reconnect backoff wait instead of letting it sleep out the full interval.
func TestTopicSubscriptionCloseInterruptsReconnectBackoff(t *testing.T) {
	disconnectErr := status.Error(codes.Unavailable, "stream disconnected")
	var subscribeCalls atomic.Uint64
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			if subscribeCalls.Add(1) == 1 {
				return &testSubscribeClient{
					results: []recvResult{
						{item: heartbeatItem()},
						{err: disconnectErr},
					},
				}, nil
			}
			return nil, disconnectErr
		},
	}
	strategy := newSignalingBackoffRetryStrategy(10 * 60 * 1000)
	client, pool := newAccountingTestClientWithStrategy(
		streamClient,
		time.Second,
		strategy,
	)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	eventDone := make(chan error, 1)
	go func() {
		_, eventErr := sub.Event(context.Background())
		eventDone <- eventErr
	}()

	// Wait until the retry loop has committed to the 10-minute backoff, then
	// Close. The closedSignal case must interrupt the wait promptly.
	<-strategy.entered
	sub.Close()

	select {
	case eventErr := <-eventDone:
		if eventErr == nil {
			t.Fatal("expected Event to return an error after Close")
		}
		assertErrorCode(t, eventErr, momentoerrors.CanceledError)
		if !errors.Is(eventErr, context.Canceled) {
			t.Fatalf("terminal error should unwrap to context.Canceled, got %v", eventErr)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Event did not return within 5s of Close (backoff not interrupted)")
	}
	assertAccounting(t, pool, 0)
}

// TestTopicSubscriptionReconnectUsesSubscribeContext pins stream-context
// lineage: reconnected streams are parented to the ctx passed to Subscribe,
// so canceling the Event ctx that happened to drive the reconnect must not
// kill the replacement stream.
func TestTopicSubscriptionReconnectUsesSubscribeContext(t *testing.T) {
	disconnectErr := status.Error(codes.Unavailable, "stream disconnected")
	var subscribeCalls atomic.Uint64
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			if subscribeCalls.Add(1) == 1 {
				return &testSubscribeClient{
					results: []recvResult{
						{item: heartbeatItem()},
						{err: disconnectErr},
					},
				}, nil
			}
			return &testSubscribeClient{
				results: []recvResult{
					{item: topicItem(2)},
					{item: topicItem(3)},
				},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, time.Second)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	// This Event call drives the disconnect -> reconnect -> item(2) sequence.
	eventCtx, cancelEventCtx := context.WithCancel(context.Background())
	event, eventErr := sub.Event(eventCtx)
	if eventErr != nil {
		t.Fatalf("Event returned error: %v", eventErr)
	}
	if _, ok := event.(TopicItem); !ok {
		t.Fatalf("Event = %T, want TopicItem after reconnect", event)
	}
	assertAccounting(t, pool, 1)

	// Cancel the ctx that drove the reconnect; the reconnected stream must
	// survive because it is parented to the Subscribe-time ctx.
	cancelEventCtx()

	event, eventErr = sub.Event(context.Background())
	if eventErr != nil {
		t.Fatalf("Event after canceling the reconnect-driving ctx returned error: %v", eventErr)
	}
	if _, ok := event.(TopicItem); !ok {
		t.Fatalf("Event = %T, want TopicItem from the surviving stream", event)
	}
	assertAccounting(t, pool, 1)

	sub.Close()
	assertAccounting(t, pool, 0)
}

// TestTopicSubscriptionCloseUnblocksEventBlockedInRecv pins the documented
// Close contract: an Event call blocked in Recv returns promptly with an
// error when Close cancels the stream context, with the slot released.
func TestTopicSubscriptionCloseUnblocksEventBlockedInRecv(t *testing.T) {
	recvBlocked := make(chan struct{})
	var recvBlockedOnce sync.Once
	streamClient := &testPubsubClient{
		subscribe: func(ctx context.Context, _ *pb.XSubscriptionRequest, _ ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			heartbeatSent := false
			return &testSubscribeClient{
				recv: func() (*pb.XSubscriptionItem, error) {
					if !heartbeatSent {
						heartbeatSent = true
						return heartbeatItem(), nil
					}
					recvBlockedOnce.Do(func() { close(recvBlocked) })
					<-ctx.Done()
					return nil, status.FromContextError(ctx.Err()).Err()
				},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, time.Second)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	eventDone := make(chan error, 1)
	go func() {
		_, eventErr := sub.Event(context.Background())
		eventDone <- eventErr
	}()

	<-recvBlocked // Event is now blocked inside Recv
	sub.Close()

	select {
	case eventErr := <-eventDone:
		if eventErr == nil {
			t.Fatal("expected Event to return an error after Close")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Event blocked in Recv did not return within 5s of Close")
	}
	assertAccounting(t, pool, 0)
}

// TestTopicSubscriptionReconnectStopsWhenPoolIsClosed verifies the retry loop
// treats the pool's shutdown CanceledError as terminal: even an always-retry
// strategy must not spin against a closed pool.
func TestTopicSubscriptionReconnectStopsWhenPoolIsClosed(t *testing.T) {
	disconnectErr := status.Error(codes.Unavailable, "stream disconnected")
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				results: []recvResult{
					{item: heartbeatItem()},
					{err: disconnectErr},
				},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, time.Second)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	// Shut the pool down, then drive Event into the disconnect. Every
	// reconnect attempt now fails with the pool's CanceledError; the
	// always-retry strategy must not loop on it.
	pool.closed.Store(true)

	eventDone := make(chan error, 1)
	go func() {
		_, eventErr := sub.Event(context.Background())
		eventDone <- eventErr
	}()

	select {
	case eventErr := <-eventDone:
		if eventErr == nil {
			t.Fatal("expected Event to return an error after the pool closed")
		}
		assertErrorCode(t, eventErr, momentoerrors.CanceledError)
		if !errors.Is(eventErr, context.Canceled) {
			t.Fatalf("terminal error should unwrap to context.Canceled, got %v", eventErr)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Event did not return within 5s; reconnect loop is spinning against a closed pool")
	}
	assertAccounting(t, pool, 0)
}

// TestTopicPublishUnaryReleasesReservationOncePerCall verifies the unary
// path's defer reservation.Release() runs exactly once per publish.
func TestTopicPublishUnaryReleasesReservationOncePerCall(t *testing.T) {
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return nil, status.Error(codes.Internal, "Subscribe should not be called from a publish-only test")
		},
	}
	pool := newTestTopicGrpcConnectionPool(streamClient)
	client := &pubSubClient{
		unaryGrpcConnectionPool: pool,
		log:                     logger.NewNoopMomentoLoggerFactory().GetLogger("unary-publish-test"),
		requestTimeout:          time.Second,
	}

	const publishCount = 50
	for i := 0; i < publishCount; i++ {
		err := client.topicPublish(context.Background(), &TopicPublishRequest{
			CacheName: "cache",
			TopicName: "topic",
			Value:     String("payload"),
		})
		if err != nil {
			t.Fatalf("topicPublish %d returned error: %v", i, err)
		}
	}

	// The Bytes branch and the unsupported-value branch release too.
	if err := client.topicPublish(context.Background(), &TopicPublishRequest{
		CacheName: "cache",
		TopicName: "topic",
		Value:     Bytes("payload"),
	}); err != nil {
		t.Fatalf("topicPublish with Bytes returned error: %v", err)
	}
	err := client.topicPublish(context.Background(), &TopicPublishRequest{
		CacheName: "cache",
		TopicName: "topic",
		Value:     unsupportedTopicValue{},
	})
	if err == nil {
		t.Fatal("expected an error for an unsupported topic value type")
	}
	assertErrorCode(t, err, momentoerrors.InvalidArgumentError)

	if got := pool.releaseCount.Load(); got != publishCount+2 {
		t.Fatalf("releaseCount = %d, want %d (exactly one Release per publish)", got, publishCount+2)
	}
}

// unsupportedTopicValue drives topicPublish's default branch.
type unsupportedTopicValue struct{}

func (unsupportedTopicValue) isTopicValue() {}

// TestTopicSubscriptionReconnectEstablishmentInterruptedByEventCtx pins the
// pause contract against a dial that blocks: the reconnect stream is parented
// to the Subscribe ctx for lineage, but the caller's per-call ctx must still
// interrupt a blocked establishment promptly, pausing recovery for the next
// call instead of terminating the subscription.
func TestTopicSubscriptionReconnectEstablishmentInterruptedByEventCtx(t *testing.T) {
	disconnectErr := status.Error(codes.Unavailable, "stream disconnected")
	establishmentEntered := make(chan struct{})
	var establishmentOnce sync.Once
	var subscribeCalls atomic.Uint64
	streamClient := &testPubsubClient{
		subscribe: func(ctx context.Context, _ *pb.XSubscriptionRequest, _ ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			switch subscribeCalls.Add(1) {
			case 1:
				return &testSubscribeClient{
					results: []recvResult{
						{item: heartbeatItem()},
						{err: disconnectErr},
					},
				}, nil
			case 2:
				// A dial that blocks mid-outage: completes only when the
				// stream context dies.
				establishmentOnce.Do(func() { close(establishmentEntered) })
				<-ctx.Done()
				return nil, status.FromContextError(ctx.Err()).Err()
			default:
				return &testSubscribeClient{
					results: []recvResult{{item: topicItem(2)}},
				}, nil
			}
		},
	}
	client, pool := newAccountingTestClient(streamClient, time.Second)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	eventCtx, cancelEventCtx := context.WithCancel(context.Background())
	eventDone := make(chan error, 1)
	go func() {
		_, eventErr := sub.Event(eventCtx)
		eventDone <- eventErr
	}()

	// Cancel the caller's ctx while the reconnect dial is blocked. Event must
	// return promptly even though the stream is parented to the Subscribe ctx.
	<-establishmentEntered
	cancelEventCtx()

	select {
	case eventErr := <-eventDone:
		if eventErr == nil {
			t.Fatal("expected Event to return an error after context cancellation")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Event blocked in reconnect establishment did not return within 5s of ctx cancellation")
	}
	assertAccounting(t, pool, 0)

	// Paused, not dead: a live ctx resumes and delivers from a fresh stream.
	event, eventErr := sub.Event(context.Background())
	if eventErr != nil {
		t.Fatalf("Event after interrupted establishment returned error: %v", eventErr)
	}
	if _, ok := event.(TopicItem); !ok {
		t.Fatalf("Event = %T, want TopicItem from the resumed stream", event)
	}
	assertAccounting(t, pool, 1)

	sub.Close()
	assertAccounting(t, pool, 0)
}

// TestSubscribeWatchdogDoesNotCancelLiveSubscription guards the sendSubscribe
// watchdog race: when both firstMessageDone and firstMessageCtx are ready by
// the time the watchdog wakes, Go's select picks randomly. Without the
// handshakeDecided arbiter, ~half of those wakeups would cancel a
// subscription that was already handed back to the caller. Concurrent Subscribe calls
// force the watchdog to be descheduled across the race window.
func TestSubscribeWatchdogDoesNotCancelLiveSubscription(t *testing.T) {
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				results: []recvResult{{item: heartbeatItem()}},
			}, nil
		},
	}
	// 50ms gives sendSubscribe time to populate subChan before the timeout
	// fires, while still leaving defer cancel() racing the watchdog.
	client, _ := newAccountingTestClient(streamClient, 50*time.Millisecond)

	const goroutines = 64
	const perGoroutine = 50

	var wg sync.WaitGroup
	var successes atomic.Int64
	var cancelled atomic.Int64

	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
				if err != nil {
					continue
				}
				successes.Add(1)
				ts := sub.(*topicSubscription)
				state := ts.state.Load()
				select {
				case <-state.cancelContext.Done():
					cancelled.Add(1)
				default:
				}
				sub.Close()
			}
		}()
	}
	wg.Wait()

	if cancelled.Load() > 0 {
		t.Fatalf("watchdog cancelled %d of %d returned subscriptions (race condition)",
			cancelled.Load(), successes.Load())
	}
	if successes.Load() == 0 {
		t.Fatal("no Subscribe calls succeeded; test is vacuous")
	}
}

// TestSubscribeDeadlineBoundaryNeverLeaksOrDeliversDeadSubscriptions hammers
// the handshake-deadline boundary, where an instant heartbeat and a tiny
// requestTimeout race inside Subscribe's select. Every returned subscription
// must be live (the deadline tie prefers a delivered subscription), every
// timeout must leave no slot behind, and the counters must drain to zero.
func TestSubscribeDeadlineBoundaryNeverLeaksOrDeliversDeadSubscriptions(t *testing.T) {
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				results: []recvResult{{item: heartbeatItem()}},
			}, nil
		},
	}
	// 1ms keeps the deadline and the instant heartbeat permanently racing.
	client, pool := newAccountingTestClient(streamClient, time.Millisecond)

	const goroutines = 16
	const perGoroutine = 100

	var wg sync.WaitGroup
	var successes, timeouts atomic.Int64
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
				if err != nil {
					timeouts.Add(1)
					continue
				}
				successes.Add(1)
				state := sub.(*topicSubscription).state.Load()
				select {
				case <-state.cancelContext.Done():
					t.Error("Subscribe returned a subscription with a dead stream")
				default:
				}
				sub.Close()
			}
		}()
	}
	wg.Wait()

	if successes.Load() == 0 && timeouts.Load() == 0 {
		t.Fatal("no outcomes recorded; test is vacuous")
	}
	// Timed-out handshakes release in sendSubscribe's goroutine; poll for the
	// stragglers.
	waitForReleasedAccounting(t, pool)
}

// TestTopicSubscribeReleasesSlotOnFirstMessageTimeout exercises the
// firstMessageCtx timeout path. Recv blocks on the stream context (matching
// production gRPC behavior); the watchdog cancels it after the timeout fires
// and the error path releases the slot.
func TestTopicSubscribeReleasesSlotOnFirstMessageTimeout(t *testing.T) {
	streamClient := &testPubsubClient{
		subscribe: func(ctx context.Context, _ *pb.XSubscriptionRequest, _ ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				recv: func() (*pb.XSubscriptionItem, error) {
					<-ctx.Done()
					return nil, ctx.Err()
				},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, 50*time.Millisecond)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err == nil {
		t.Fatal("expected Subscribe to return a timeout error")
	}
	if sub != nil {
		t.Fatal("expected Subscribe to return a nil subscription on timeout")
	}
	assertErrorCode(t, err, momentoerrors.TimeoutError)

	waitForReleasedAccounting(t, pool)
	if got := pool.releaseCount.Load(); got != 1 {
		t.Fatalf("releaseCount = %d, want 1 (exactly one Release for the cleanup path)", got)
	}
}

// waitForReleasedAccounting polls until the pool counters drain to zero;
// cleanup runs in sendSubscribe's goroutine after Subscribe has returned, so
// the release lands asynchronously. Generous deadline for loaded CI runners.
func waitForReleasedAccounting(t *testing.T, pool *testTopicGrpcConnectionPool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if pool.activeCount.Load() == 0 && pool.manager.NumActiveSubscriptions.Load() == 0 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	assertAccounting(t, pool, 0)
}

// TestSubscribeReleasesSlotWhenStreamEstablishmentHangs pins that the
// handshake watchdog bounds stream ESTABLISHMENT too: a stream-open call that
// blocks (black-holed endpoint) must not hold the reserved slot past the
// handshake deadline, even though Subscribe already returned TimeoutError.
func TestSubscribeReleasesSlotWhenStreamEstablishmentHangs(t *testing.T) {
	streamClient := &testPubsubClient{
		subscribe: func(ctx context.Context, _ *pb.XSubscriptionRequest, _ ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			// Establishment completes only when the stream context is
			// cancelled, as with a black-holed endpoint.
			<-ctx.Done()
			return nil, ctx.Err()
		},
	}
	client, pool := newAccountingTestClient(streamClient, 50*time.Millisecond)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err == nil {
		t.Fatal("expected Subscribe to return a timeout error")
	}
	if sub != nil {
		t.Fatal("expected Subscribe to return a nil subscription on timeout")
	}
	assertErrorCode(t, err, momentoerrors.TimeoutError)

	waitForReleasedAccounting(t, pool)
	if got := pool.releaseCount.Load(); got != 1 {
		t.Fatalf("releaseCount = %d, want 1 (exactly one Release for the hung establishment)", got)
	}
}

// TestSubscribeLateFirstMessageAfterTimeoutReleasesSlot pins the
// dead-on-arrival guard: the first heartbeat lands only after Subscribe has
// returned via the timeout path, so the handshake arbiter must refuse to
// deliver the already-cancelled stream and release its slot instead.
func TestSubscribeLateFirstMessageAfterTimeoutReleasesSlot(t *testing.T) {
	releaseHeartbeat := make(chan struct{})
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				recv: func() (*pb.XSubscriptionItem, error) {
					<-releaseHeartbeat
					return heartbeatItem(), nil
				},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, 50*time.Millisecond)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err == nil {
		t.Fatal("expected Subscribe to return a timeout error")
	}
	if sub != nil {
		t.Fatal("expected Subscribe to return a nil subscription on timeout")
	}
	assertErrorCode(t, err, momentoerrors.TimeoutError)

	// Deliver the heartbeat only now, after the timeout return. The watchdog
	// already cancelled the stream, so sendSubscribe must fail the handshake
	// and release the slot instead of delivering a dead subscription.
	close(releaseHeartbeat)

	waitForReleasedAccounting(t, pool)
	if got := pool.releaseCount.Load(); got != 1 {
		t.Fatalf("releaseCount = %d, want 1 (exactly one Release for the late-message path)", got)
	}
}

// TestSubscribeUserContextCancelReturnsCanceledError verifies a caller
// cancellation is reported as CanceledError — never as a handshake timeout —
// even though firstMessageCtx (a child of the caller's ctx) fires at the same
// instant. Repeated because the wrong outcome was a random select pick.
func TestSubscribeUserContextCancelReturnsCanceledError(t *testing.T) {
	// One gate per Subscribe call; recv holds the handshake open until the
	// test closes the gate, keeping errChan empty until Subscribe has returned.
	gates := make(chan chan struct{}, 16)
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			gate := make(chan struct{})
			gates <- gate
			return &testSubscribeClient{
				recv: func() (*pb.XSubscriptionItem, error) {
					<-gate
					return heartbeatItem(), nil
				},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, 10*time.Second)

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		subscribeDone := make(chan error, 1)
		go func() {
			_, err := client.Subscribe(ctx, testSubscribeRequest())
			subscribeDone <- err
		}()

		gate := <-gates // the stream is open and the handshake is in flight
		cancel()

		var err error
		select {
		case err = <-subscribeDone:
		case <-time.After(5 * time.Second):
			t.Fatal("Subscribe did not return after context cancellation")
		}
		if err == nil {
			t.Fatal("expected Subscribe to return an error after context cancellation")
		}
		assertErrorCode(t, err, momentoerrors.CanceledError)

		// Let the held-open handshake finish; the drainer must release the slot.
		close(gate)
		waitForReleasedAccounting(t, pool)
	}
}

// TestDrainerClosesLateDeliveredSubscription pins drainAndCloseSubscription's
// subscription branch directly: a live subscription that lands on subChan
// after Subscribe has returned must be closed and its slot released.
func TestDrainerClosesLateDeliveredSubscription(t *testing.T) {
	streamClient := &testPubsubClient{
		subscribe: func(context.Context, *pb.XSubscriptionRequest, ...grpc.CallOption) (pb.Pubsub_SubscribeClient, error) {
			return &testSubscribeClient{
				results: []recvResult{{item: heartbeatItem()}},
			}, nil
		},
	}
	client, pool := newAccountingTestClient(streamClient, time.Second)

	sub, err := client.Subscribe(context.Background(), testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	subChan := make(chan *topicSubscription, 1)
	errChan := make(chan error, 1)
	subChan <- sub.(*topicSubscription)
	drainAndCloseSubscription(subChan, errChan)
	assertAccounting(t, pool, 0)
}
