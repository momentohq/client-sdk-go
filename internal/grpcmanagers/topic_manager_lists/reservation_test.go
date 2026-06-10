package topic_manager_lists

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
)

// newCountingReleaser returns a release function that increments a counter
// each time it is invoked, plus the counter so tests can observe invocations.
func newCountingReleaser() (func(*grpcmanagers.TopicGrpcManager) int64, *atomic.Int64) {
	var calls atomic.Int64
	release := func(mgr *grpcmanagers.TopicGrpcManager) int64 {
		calls.Add(1)
		return mgr.NumActiveSubscriptions.Add(-1)
	}
	return release, &calls
}

func TestReservation_ReleaseIsIdempotent(t *testing.T) {
	t.Parallel()

	mgr := &grpcmanagers.TopicGrpcManager{}
	mgr.NumActiveSubscriptions.Store(1)

	release, calls := newCountingReleaser()
	r := NewReservation(mgr, release)

	first := r.Release()
	if first != 0 {
		t.Fatalf("first Release returned %d, want 0 (manager counter after decrement)", first)
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("release function called %d times after first Release, want 1", got)
	}

	// Redundant calls must not double-decrement.
	for i := 0; i < 5; i++ {
		r.Release()
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("release function called %d times after repeated Release, want 1", got)
	}
	if got := mgr.NumActiveSubscriptions.Load(); got != 0 {
		t.Fatalf("NumActiveSubscriptions = %d, want 0 (only one decrement should have occurred)", got)
	}
}

func TestReservation_ReleaseReportsCurrentCountOnSubsequentCalls(t *testing.T) {
	t.Parallel()

	mgr := &grpcmanagers.TopicGrpcManager{}
	mgr.NumActiveSubscriptions.Store(3)

	release, _ := newCountingReleaser()
	r := NewReservation(mgr, release)

	first := r.Release()
	if first != 2 {
		t.Fatalf("first Release returned %d, want 2 (started at 3, decremented once)", first)
	}

	// External actor adjusts the counter to simulate other subscriptions. The
	// second Release should report whatever the live value is, without
	// performing another decrement.
	mgr.NumActiveSubscriptions.Add(5)
	second := r.Release()
	if second != 7 {
		t.Fatalf("second Release returned %d, want 7 (live counter value, no decrement)", second)
	}
}

func TestReservation_ManagerReturnsUnderlying(t *testing.T) {
	t.Parallel()

	mgr := &grpcmanagers.TopicGrpcManager{}
	release, _ := newCountingReleaser()
	r := NewReservation(mgr, release)

	if r.Manager() != mgr {
		t.Fatal("Manager() should return the manager passed to NewReservation")
	}

	r.Release()
	if r.Manager() != mgr {
		t.Fatal("Manager() should remain accessible after Release for in-flight callers")
	}
}

func TestReservation_ConcurrentReleaseDecrementsOnce(t *testing.T) {
	t.Parallel()

	// Run many concurrent Release calls on the same Reservation; only the
	// first should perform the decrement. Validates the CAS guarantee under
	// contention and exercises the race detector when -race is enabled.
	const goroutines = 256

	mgr := &grpcmanagers.TopicGrpcManager{}
	mgr.NumActiveSubscriptions.Store(int64(goroutines + 10))

	release, calls := newCountingReleaser()
	r := NewReservation(mgr, release)

	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			<-start
			r.Release()
		}()
	}
	close(start)
	wg.Wait()

	if got := calls.Load(); got != 1 {
		t.Fatalf("release function called %d times, want 1 under concurrent Release", got)
	}
	if got := mgr.NumActiveSubscriptions.Load(); got != int64(goroutines+9) {
		t.Fatalf("NumActiveSubscriptions = %d, want %d (started at goroutines+10, exactly one decrement)", got, goroutines+9)
	}
}
