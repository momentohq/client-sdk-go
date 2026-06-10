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

	_, err := client.topicSubscribe(context.Background(), testSubscribeRequest())
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
	// Generous timeout: the fake handshake is synchronous, so the watchdog
	// must never fire even on an oversubscribed -race CI runner.
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

	if got := pool.releaseCount.Load(); got != publishCount {
		t.Fatalf("releaseCount = %d, want %d (exactly one Release per publish)", got, publishCount)
	}
}
