package topic_manager_lists

import (
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

// TopicGrpcConnectionPool hands out the next available manager to the pubsub
// client on request.
type TopicGrpcConnectionPool interface {
	// GetNextTopicGrpcManager reserves a manager from the pool. Callers must
	// call Reservation.Release when done; Release is idempotent so overlapping
	// cleanup paths can each call it. Dropping a Reservation without releasing
	// leaks the slot for the lifetime of the pool; Close shuts the pool down
	// but does not reclaim leaked slots.
	//
	// The unary pool's Release is a no-op since unary requests don't hold
	// long-lived slots; callers invoke it uniformly anyway.
	GetNextTopicGrpcManager() (*Reservation, momentoerrors.MomentoSvcErr)

	// Close shuts down all gRPC connections in the pool.
	Close()
}

// StreamGrpcManagerRequest is the dispatcher-to-caller envelope: either the
// next manager or the allocation error.
type StreamGrpcManagerRequest struct {
	TopicManager *grpcmanagers.TopicGrpcManager
	Err          momentoerrors.MomentoSvcErr
}
