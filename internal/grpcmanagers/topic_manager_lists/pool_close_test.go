package topic_manager_lists

import (
	"strings"
	"sync"
	"testing"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
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

func newStaticTestPool(t *testing.T, numChannels uint32) *staticStreamGrpcManagerPool {
	t.Helper()
	pool, err := NewStaticStreamGrpcManagerPool(
		newTestManagerRequest(t),
		numChannels,
		logger.NewNoopMomentoLoggerFactory().GetLogger("pool-test"),
	)
	if err != nil {
		t.Fatalf("NewStaticStreamGrpcManagerPool returned error: %v", err)
	}
	return pool
}

// TestTopicStreamPoolHoldsNoSlotsWhileIdle pins the exact-counting contract:
// an idle pool reports zero active streams, each reservation counts exactly
// one, and a release is immediately visible.
func TestTopicStreamPoolHoldsNoSlotsWhileIdle(t *testing.T) {
	t.Parallel()

	pool := newStaticTestPool(t, 1)
	defer pool.Close()

	if got := pool.GetCurrentActiveStreamsCount(); got != 0 {
		t.Fatalf("idle active streams count = %d, want 0", got)
	}

	reservation, err := pool.GetNextTopicGrpcManager()
	if err != nil {
		t.Fatalf("GetNextTopicGrpcManager returned error: %v", err)
	}
	if got := pool.GetCurrentActiveStreamsCount(); got != 1 {
		t.Fatalf("active streams count with one reservation = %d, want 1", got)
	}

	reservation.Release()
	if got := pool.GetCurrentActiveStreamsCount(); got != 0 {
		t.Fatalf("active streams count after release = %d, want 0", got)
	}
}

// TestTopicStreamPoolExhaustionErrorIsNeverStale verifies capacity is
// evaluated per call: a pool that just returned ResourceExhausted hands out a
// manager on the very next call once a slot frees.
func TestTopicStreamPoolExhaustionErrorIsNeverStale(t *testing.T) {
	t.Parallel()

	pool := newStaticTestPool(t, 1)
	defer pool.Close()

	reservations := make([]*Reservation, 0, config.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
	for i := 0; i < config.MAX_CONCURRENT_STREAMS_PER_CHANNEL; i++ {
		reservation, err := pool.GetNextTopicGrpcManager()
		if err != nil {
			t.Fatalf("GetNextTopicGrpcManager %d returned error: %v", i, err)
		}
		reservations = append(reservations, reservation)
	}

	if _, err := pool.GetNextTopicGrpcManager(); err == nil {
		t.Fatal("expected ResourceExhausted error from a full pool")
	} else if err.Code() != momentoerrors.ClientResourceExhaustedError {
		t.Fatalf("error code = %s, want %s", err.Code(), momentoerrors.ClientResourceExhaustedError)
	}

	// Free one slot; the next call must succeed immediately.
	reservations[0].Release()
	reservation, err := pool.GetNextTopicGrpcManager()
	if err != nil {
		t.Fatalf("GetNextTopicGrpcManager after release returned error: %v", err)
	}
	reservation.Release()
	for _, r := range reservations[1:] {
		r.Release()
	}
	if got := pool.GetCurrentActiveStreamsCount(); got != 0 {
		t.Fatalf("active streams count after releasing all = %d, want 0", got)
	}
}

// TestTopicStreamPoolGetAfterCloseReturnsCanceled verifies shutdown behavior:
// Close is idempotent (including concurrently) and subsequent reservations
// are refused with CanceledError.
func TestTopicStreamPoolGetAfterCloseReturnsCanceled(t *testing.T) {
	t.Parallel()

	pool := newStaticTestPool(t, 1)

	done := make(chan struct{})
	go func() {
		defer close(done)
		pool.Close()
	}()
	pool.Close()
	<-done
	pool.Close()

	reservation, err := pool.GetNextTopicGrpcManager()
	if err == nil {
		t.Fatal("expected GetNextTopicGrpcManager to fail after Close")
	}
	if err.Code() != momentoerrors.CanceledError {
		t.Fatalf("error code = %s, want %s", err.Code(), momentoerrors.CanceledError)
	}
	if reservation != nil {
		t.Fatal("expected nil reservation after Close")
	}
	if got := pool.GetCurrentActiveStreamsCount(); got != 0 {
		t.Fatalf("active streams count after Close = %d, want 0", got)
	}
}

// TestTopicStreamPoolGetRacingCloseIsSafe hammers GetNextTopicGrpcManager
// while Close runs: every call must return either a usable Reservation or a
// shutdown/capacity error — never panic or hang — and releasing everything
// drains the counters to zero.
func TestTopicStreamPoolGetRacingCloseIsSafe(t *testing.T) {
	t.Parallel()

	const goroutines = 32
	const perGoroutine = 25

	pool := newStaticTestPool(t, 2)

	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(goroutines + 1)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			<-start
			for i := 0; i < perGoroutine; i++ {
				reservation, err := pool.GetNextTopicGrpcManager()
				if err != nil {
					if code := err.Code(); code != momentoerrors.CanceledError && code != momentoerrors.ClientResourceExhaustedError {
						t.Errorf("unexpected error code %s: %v", code, err)
					}
					continue
				}
				reservation.Release()
			}
		}()
	}
	go func() {
		defer wg.Done()
		<-start
		pool.Close()
	}()
	close(start)
	wg.Wait()

	if got := pool.GetCurrentActiveStreamsCount(); got != 0 {
		t.Fatalf("active streams count after race = %d, want 0", got)
	}
}

// TestTopicStreamPoolRespectsPerChannelCap fills a multi-channel pool to
// capacity and checks the invariant the round-robin reservation must hold:
// no channel ever exceeds MAX_CONCURRENT_STREAMS_PER_CHANNEL, and a full pool
// lands exactly at the cap on every channel.
func TestTopicStreamPoolRespectsPerChannelCap(t *testing.T) {
	t.Parallel()

	const numChannels = 3
	pool := newStaticTestPool(t, numChannels)
	defer pool.Close()

	perChannel := int64(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
	total := numChannels * config.MAX_CONCURRENT_STREAMS_PER_CHANNEL
	reservations := make([]*Reservation, 0, total)
	for i := 0; i < total; i++ {
		reservation, err := pool.GetNextTopicGrpcManager()
		if err != nil {
			t.Fatalf("GetNextTopicGrpcManager %d returned error: %v", i, err)
		}
		reservations = append(reservations, reservation)
		for channel, manager := range pool.grpcManagers {
			if got := manager.NumActiveSubscriptions.Load(); got > perChannel {
				t.Fatalf("channel %d exceeded the per-channel cap after %d reservations: %d > %d", channel, i+1, got, perChannel)
			}
		}
	}
	for channel, manager := range pool.grpcManagers {
		if got := manager.NumActiveSubscriptions.Load(); got != perChannel {
			t.Fatalf("channel %d subscription count at full pool = %d, want exactly %d", channel, got, perChannel)
		}
	}

	for _, reservation := range reservations {
		reservation.Release()
	}
	if got := pool.GetCurrentActiveStreamsCount(); got != 0 {
		t.Fatalf("active streams count after releasing all = %d, want 0", got)
	}
}

// TestTopicDynamicStreamPoolGrowsOnDemand verifies the dynamic pool adds a
// channel when the current capacity is consumed, with exact counters and no
// RPCs (slots are reserved without opening streams).
func TestTopicDynamicStreamPoolGrowsOnDemand(t *testing.T) {
	t.Parallel()

	perChannel := config.MAX_CONCURRENT_STREAMS_PER_CHANNEL
	pool, err := NewDynamicStreamGrpcManagerPool(
		newTestManagerRequest(t),
		uint32(perChannel+50), // maxManagerCount = 2
		logger.NewNoopMomentoLoggerFactory().GetLogger("pool-test"),
	)
	if err != nil {
		t.Fatalf("NewDynamicStreamGrpcManagerPool returned error: %v", err)
	}
	defer pool.Close()

	if got := pool.GetCurrentNumberOfGrpcManagers(); got != 1 {
		t.Fatalf("initial manager count = %d, want 1", got)
	}

	reservations := make([]*Reservation, 0, perChannel+1)
	for i := 0; i < perChannel; i++ {
		reservation, getErr := pool.GetNextTopicGrpcManager()
		if getErr != nil {
			t.Fatalf("GetNextTopicGrpcManager %d returned error: %v", i, getErr)
		}
		reservations = append(reservations, reservation)
	}
	if got := pool.GetCurrentNumberOfGrpcManagers(); got != 1 {
		t.Fatalf("manager count at exactly one channel's capacity = %d, want 1", got)
	}

	// One past the first channel's capacity forces growth.
	reservation, getErr := pool.GetNextTopicGrpcManager()
	if getErr != nil {
		t.Fatalf("GetNextTopicGrpcManager past capacity returned error: %v", getErr)
	}
	reservations = append(reservations, reservation)
	if got := pool.GetCurrentNumberOfGrpcManagers(); got != 2 {
		t.Fatalf("manager count after growth = %d, want 2", got)
	}
	if got := pool.GetCurrentActiveStreamsCount(); got != uint64(perChannel+1) {
		t.Fatalf("active streams count after growth = %d, want %d", got, perChannel+1)
	}

	for _, r := range reservations {
		r.Release()
	}
	if got := pool.GetCurrentActiveStreamsCount(); got != 0 {
		t.Fatalf("active streams count after releasing all = %d, want 0", got)
	}
}

// TestTopicDynamicStreamPoolStopsAtMaxSubscriptions verifies the dynamic pool
// refuses reservations past maxSubscriptions with a fresh error message.
func TestTopicDynamicStreamPoolStopsAtMaxSubscriptions(t *testing.T) {
	t.Parallel()

	perChannel := config.MAX_CONCURRENT_STREAMS_PER_CHANNEL
	pool, err := NewDynamicStreamGrpcManagerPool(
		newTestManagerRequest(t),
		uint32(perChannel), // a single channel, no growth allowed
		logger.NewNoopMomentoLoggerFactory().GetLogger("pool-test"),
	)
	if err != nil {
		t.Fatalf("NewDynamicStreamGrpcManagerPool returned error: %v", err)
	}
	defer pool.Close()

	for i := 0; i < perChannel; i++ {
		if _, getErr := pool.GetNextTopicGrpcManager(); getErr != nil {
			t.Fatalf("GetNextTopicGrpcManager %d returned error: %v", i, getErr)
		}
	}
	_, getErr := pool.GetNextTopicGrpcManager()
	if getErr == nil {
		t.Fatal("expected error past maxSubscriptions")
	}
	if getErr.Code() != momentoerrors.ClientResourceExhaustedError {
		t.Fatalf("error code = %s, want %s", getErr.Code(), momentoerrors.ClientResourceExhaustedError)
	}
	if !strings.Contains(getErr.Error(), "maximum number of concurrent grpc streams") {
		t.Fatalf("unexpected error message: %v", getErr)
	}
	if got := pool.GetCurrentNumberOfGrpcManagers(); got != 1 {
		t.Fatalf("manager count = %d, want 1 (growth not allowed)", got)
	}
}
