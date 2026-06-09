package topic_manager_lists

import (
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
)

// Reservation is an acquired stream slot from a TopicGrpcConnectionPool.
// Call Release when done; redundant calls are safe. Decrement and leak
// semantics come from the pool's release function (the unary pool's is a
// no-op). Dropping a stream-pool Reservation without calling Release leaks
// the slot for the lifetime of the pool.
type Reservation struct {
	manager  *grpcmanagers.TopicGrpcManager
	release  func(*grpcmanagers.TopicGrpcManager) int64
	released atomic.Bool
}

// NewReservation builds a Reservation backed by release. Both args must be non-nil.
func NewReservation(
	manager *grpcmanagers.TopicGrpcManager,
	release func(*grpcmanagers.TopicGrpcManager) int64,
) *Reservation {
	return &Reservation{
		manager: manager,
		release: release,
	}
}

func (r *Reservation) Manager() *grpcmanagers.TopicGrpcManager {
	return r.manager
}

// Release returns the slot to the pool. Only the first call invokes the
// pool's release; subsequent calls return the live per-channel count without
// effect. That value is advisory — a concurrent winning Release may not have
// performed its decrement yet.
func (r *Reservation) Release() int64 {
	if !r.released.CompareAndSwap(false, true) {
		return r.manager.NumActiveSubscriptions.Load()
	}
	return r.release(r.manager)
}
