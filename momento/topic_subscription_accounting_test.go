package momento

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/retry"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
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

func (p *testTopicGrpcConnectionPool) GetNextTopicGrpcManager() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	p.manager.NumActiveSubscriptions.Add(1)
	p.activeCount.Add(1)
	return p.manager, nil
}

func (p *testTopicGrpcConnectionPool) ReleaseTopicGrpcManager(manager *grpcmanagers.TopicGrpcManager) int64 {
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

func newAccountingTestClient(streamClient pb.PubsubClient, timeout time.Duration) (defaultTopicClient, *testTopicGrpcConnectionPool) {
	pool := newTestTopicGrpcConnectionPool(streamClient)
	client := defaultTopicClient{
		pubSubClient: &pubSubClient{
			streamGrpcConnectionPool: pool,
		},
		log:            logger.NewNoopMomentoLoggerFactory().GetLogger("topic-subscription-accounting-test"),
		requestTimeout: timeout,
		retryStrategy:  alwaysRetryStrategy{},
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

	_, _, _, _, err := client.topicSubscribe(context.Background(), testSubscribeRequest())
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
