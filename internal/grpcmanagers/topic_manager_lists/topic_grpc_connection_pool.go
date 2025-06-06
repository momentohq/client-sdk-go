package topic_manager_lists

import (
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

// TopicGrpcConnectionPool is the base interface for all topic grpc connection pool structs,
// which manage a pool of grpc connections and continually provide the next available grpc stub
// for the pubsub client to use.
type TopicGrpcConnectionPool interface {
	// GetNextTopicGrpcManager returns the next available TopicGrpcManager from the pool.
	GetNextTopicGrpcManager() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr)

	// Close shuts down all the grpc connections in the pool.
	Close()
}

// StreamGrpcManagerRequest is used for putting the next available stream manager on a channel for the
// pubSubClient to pull from, or an error that specifies why no stream manager is available.
type StreamGrpcManagerRequest struct {
	TopicManager *grpcmanagers.TopicGrpcManager
	Err          momentoerrors.MomentoSvcErr
}
