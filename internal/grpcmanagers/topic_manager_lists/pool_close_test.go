package topic_manager_lists

import (
	"testing"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/models"
)

// newTestManagerRequest builds a manager request against momento-local
// defaults. grpc.NewClient dials lazily, so these tests never need a running
// server: no RPC is ever issued.
func newTestManagerRequest(t *testing.T) *models.TopicStreamGrpcManagerRequest {
	t.Helper()
	credProvider, err := auth.NewMomentoLocalProvider(&auth.MomentoLocalConfig{})
	if err != nil {
		t.Fatalf("NewMomentoLocalProvider returned error: %v", err)
	}
	return &models.TopicStreamGrpcManagerRequest{
		GrpcConfiguration:  config.NewTopicsStaticGrpcConfiguration(&config.TopicsGrpcConfigurationProps{}),
		CredentialProvider: credProvider,
	}
}

func waitForActiveStreamsCount(t *testing.T, got func() uint64, want uint64) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if got() == want {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("active streams count = %d, want %d after waiting", got(), want)
}

// TestTopicStreamPoolCloseReleasesPrefetchedSlotStatic pins the dispatcher's
// shutdown accounting: makeNextManagerAvailable eagerly reserves a slot for
// the envelope it parks on the unbuffered channel, and Close interrupting that
// parked send must return the slot so the counters end exact.
func TestTopicStreamPoolCloseReleasesPrefetchedSlotStatic(t *testing.T) {
	t.Parallel()

	pool, err := NewStaticStreamGrpcManagerPool(
		newTestManagerRequest(t),
		1,
		logger.NewNoopMomentoLoggerFactory().GetLogger("pool-close-test"),
	)
	if err != nil {
		t.Fatalf("NewStaticStreamGrpcManagerPool returned error: %v", err)
	}

	// The dispatcher prefetches one manager and parks holding its slot.
	waitForActiveStreamsCount(t, pool.GetCurrentActiveStreamsCount, 1)

	// Close waits on dispatcherDone, so the release is visible once it returns.
	pool.Close()
	if got := pool.GetCurrentActiveStreamsCount(); got != 0 {
		t.Fatalf("active streams count after Close = %d, want 0", got)
	}
	if got := pool.grpcManagers[0].NumActiveSubscriptions.Load(); got != 0 {
		t.Fatalf("manager subscription count after Close = %d, want 0", got)
	}
}

// TestTopicStreamPoolCloseReleasesPrefetchedSlotDynamic is the dynamic-pool
// twin of the static prefetch-release test.
func TestTopicStreamPoolCloseReleasesPrefetchedSlotDynamic(t *testing.T) {
	t.Parallel()

	pool, err := NewDynamicStreamGrpcManagerPool(
		newTestManagerRequest(t),
		100,
		logger.NewNoopMomentoLoggerFactory().GetLogger("pool-close-test"),
	)
	if err != nil {
		t.Fatalf("NewDynamicStreamGrpcManagerPool returned error: %v", err)
	}

	waitForActiveStreamsCount(t, pool.GetCurrentActiveStreamsCount, 1)

	pool.Close()
	if got := pool.GetCurrentActiveStreamsCount(); got != 0 {
		t.Fatalf("active streams count after Close = %d, want 0", got)
	}
	for i, manager := range pool.grpcManagers {
		if got := manager.NumActiveSubscriptions.Load(); got != 0 {
			t.Fatalf("manager %d subscription count after Close = %d, want 0", i, got)
		}
	}
}

// TestTopicStreamPoolCloseIsIdempotent guards the closeOnce hardening: double
// and concurrent Close must not panic (send on closed channel, double close)
// or hang.
func TestTopicStreamPoolCloseIsIdempotent(t *testing.T) {
	t.Parallel()

	pool, err := NewStaticStreamGrpcManagerPool(
		newTestManagerRequest(t),
		1,
		logger.NewNoopMomentoLoggerFactory().GetLogger("pool-close-test"),
	)
	if err != nil {
		t.Fatalf("NewStaticStreamGrpcManagerPool returned error: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		pool.Close()
	}()
	pool.Close()
	<-done
	pool.Close()
}
