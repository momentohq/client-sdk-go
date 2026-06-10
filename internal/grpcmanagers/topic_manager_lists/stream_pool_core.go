package topic_manager_lists

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

// streamPoolCore implements the reserve/release/Close machinery shared by the
// static and dynamic stream pools. The mutex serializes allocation, so every
// capacity decision sees live counters: no slot is held while the pool is
// idle, and a capacity error can never be stale — a Release that frees a slot
// is visible to the very next GetNextTopicGrpcManager call.
//
// Pools supply their capacity policy via getNextManager and their connection
// teardown via closeManagers; both run under the mutex.
type streamPoolCore struct {
	mu                        sync.Mutex
	closed                    bool
	currentActiveStreamsCount atomic.Uint64
	// getNextManager reserves a slot on success. Called only under mu.
	getNextManager func() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr)
	// closeManagers tears down the pool's connections. Called only under mu.
	closeManagers func()
}

func (c *streamPoolCore) initStreamPoolCore(
	getNextManager func() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr),
	closeManagers func(),
) {
	c.getNextManager = getNextManager
	c.closeManagers = closeManagers
}

// GetNextTopicGrpcManager reserves a manager from the pool. After Close it
// returns CanceledError.
func (c *streamPoolCore) GetNextTopicGrpcManager() (*Reservation, momentoerrors.MomentoSvcErr) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.CanceledError, "connection pool is shutting down", context.Canceled)
	}
	topicManager, err := c.getNextManager()
	if err != nil {
		return nil, err
	}
	return NewReservation(topicManager, c.releaseManager), nil
}

// Close shuts down all gRPC connections. Safe to call multiple times;
// concurrent calls block on the mutex until the first completes.
func (c *streamPoolCore) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	c.closeManagers()
}

// GetCurrentActiveStreamsCount returns the current number of active streams in the pool.
func (c *streamPoolCore) GetCurrentActiveStreamsCount() uint64 {
	return c.currentActiveStreamsCount.Load()
}

// releaseManager backs Reservation.Release for the stream pools. Decrements
// the per-channel and pool-wide counters; the pool-wide CAS bottoms at zero
// so it can't go negative. Deliberately not guarded by mu so releases can't
// contend with Close.
func (c *streamPoolCore) releaseManager(manager *grpcmanagers.TopicGrpcManager) int64 {
	newCount := manager.NumActiveSubscriptions.Add(-1)
	for {
		current := c.currentActiveStreamsCount.Load()
		if current == 0 {
			break
		}
		if c.currentActiveStreamsCount.CompareAndSwap(current, current-1) {
			break
		}
	}
	return newCount
}

// reserveSlot round-robins over managers until one has spare per-channel
// capacity, reserving a slot on it and bumping the pool-wide counter. index
// is owned by the calling pool and only advances under the core mutex.
// Returns nil when no channel has capacity within attempts probes.
func (c *streamPoolCore) reserveSlot(
	managers []*grpcmanagers.TopicGrpcManager,
	index *uint64,
	attempts uint32,
	log logger.MomentoLogger,
) *grpcmanagers.TopicGrpcManager {
	for i := uint32(0); i < attempts; i++ {
		(*index)++
		topicManager := managers[*index%uint64(len(managers))]
		newCount := topicManager.NumActiveSubscriptions.Add(1)
		if newCount <= int64(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL) {
			log.Debug("Starting new subscription on grpc channel %d which now has %d streams", *index%uint64(len(managers)), newCount)
			c.currentActiveStreamsCount.Add(1)
			return topicManager
		}
		topicManager.NumActiveSubscriptions.Add(-1)
	}
	return nil
}

// closeAllManagers closes every connection in managers, logging failures.
func closeAllManagers(managers []*grpcmanagers.TopicGrpcManager, log logger.MomentoLogger) {
	for _, topicManager := range managers {
		if err := topicManager.Close(); err != nil {
			log.Error("Error closing topic manager: %s", err.Error())
		}
	}
}
