package momento

import (
	"context"
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
}

func newTestTopicGrpcConnectionPool(streamClient pb.PubsubClient) *testTopicGrpcConnectionPool {
	return &testTopicGrpcConnectionPool{
		manager: &grpcmanagers.TopicGrpcManager{
			StreamClient: streamClient,
		},
	}
}

func (p *testTopicGrpcConnectionPool) GetNextTopicGrpcManager() (*topic_manager_lists.Reservation, momentoerrors.MomentoSvcErr) {
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

// fixedBackoffRetryStrategy always retries after a fixed backoff.
type fixedBackoffRetryStrategy struct {
	backoffMs int
}

func (f fixedBackoffRetryStrategy) DetermineWhenToRetry(retry.StrategyProps) *int {
	return &f.backoffMs
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

// TestTopicSubscriptionReconnectAbortsWhenContextCanceled pins the retry
// loop's response to user-context cancellation: Event must return promptly
// with balanced counters instead of retrying (here, sleeping a 10-minute
// backoff) against a context that can never produce a live stream again.
func TestTopicSubscriptionReconnectAbortsWhenContextCanceled(t *testing.T) {
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
			// Reconnect attempts keep failing so Event stays in the retry loop
			// (sleeping the fixed backoff) until the context is canceled.
			return nil, disconnectErr
		},
	}
	client, pool := newAccountingTestClientWithStrategy(
		streamClient,
		time.Second,
		fixedBackoffRetryStrategy{backoffMs: 10 * 60 * 1000},
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sub, err := client.Subscribe(ctx, testSubscribeRequest())
	if err != nil {
		t.Fatalf("Subscribe returned error: %v", err)
	}
	assertAccounting(t, pool, 1)

	eventDone := make(chan error, 1)
	go func() {
		_, eventErr := sub.Event(ctx)
		eventDone <- eventErr
	}()

	// Give Event a moment to hit the disconnect and enter the backoff sleep,
	// then cancel. Whichever point the loop has reached, cancellation must
	// surface promptly; without the ctx checks this would block ~10 minutes.
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case eventErr := <-eventDone:
		if eventErr == nil {
			t.Fatal("expected Event to return an error after context cancellation")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Event did not return within 5s of context cancellation")
	}
	assertAccounting(t, pool, 0)
}

// TestTopicSubscriptionEventContextCancelTearsDownStream verifies a ctx error
// from Event is terminal: the slot is released AND the stream context is
// canceled, so the pool's freed slot isn't backed by a still-open stream.
func TestTopicSubscriptionEventContextCancelTearsDownStream(t *testing.T) {
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

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, eventErr := sub.Event(canceledCtx); eventErr == nil {
		t.Fatal("expected Event to return an error for a canceled context")
	}
	assertAccounting(t, pool, 0)

	state := sub.(*topicSubscription).state.Load()
	if state.cancelContext.Err() == nil {
		t.Fatal("expected the stream context to be torn down with the released slot")
	}
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
	client, pool := newAccountingTestClientWithStrategy(
		streamClient,
		time.Second,
		fixedBackoffRetryStrategy{backoffMs: 10 * 60 * 1000},
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

	// Let Event hit the disconnect and enter the 10-minute backoff, then Close.
	time.Sleep(50 * time.Millisecond)
	sub.Close()

	select {
	case eventErr := <-eventDone:
		if eventErr == nil {
			t.Fatal("expected Event to return an error after Close")
		}
		assertErrorCode(t, eventErr, momentoerrors.CanceledError)
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
