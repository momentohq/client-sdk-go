package topic_manager_lists

import (
	"sync/atomic"

	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

// staticUnaryGrpcManagerPool manages a static pool of gRPC channels for unary pubsub requests.
type staticUnaryGrpcManagerPool struct {
	grpcManagers []*grpcmanagers.TopicGrpcManager
	managerIndex atomic.Uint64
	logger       logger.MomentoLogger
}

// GetNextTopicGrpcManager returns the next available TopicGrpcManager from the pool
// using a round-robin approach.
//
// Publish requests can queue up if there are >100 concurrent requests on a grpc connection,
// but unary requests will eventually complete so no need for the same level of bookkeeping
// as for the stream grpc manager pools.
func (list *staticUnaryGrpcManagerPool) GetNextTopicGrpcManager() (*grpcmanagers.TopicGrpcManager, momentoerrors.MomentoSvcErr) {
	nextManagerIndex := list.managerIndex.Add(1)
	return list.grpcManagers[nextManagerIndex%uint64(len(list.grpcManagers))], nil
}

// Close shuts down all the grpc connections in the pool.
func (list *staticUnaryGrpcManagerPool) Close() {
	for _, topicManager := range list.grpcManagers {
		err := topicManager.Close()
		if err != nil {
			list.logger.Error("Error closing topic manager: %v", err)
		}
	}
}

// NewStaticUnaryGrpcManagerPool creates a new pool with a fixed number of grpc managers for unary pubsub requests.
func NewStaticUnaryGrpcManagerPool(request *models.TopicStreamGrpcManagerRequest, numUnaryChannels uint32, logger logger.MomentoLogger) (*staticUnaryGrpcManagerPool, momentoerrors.MomentoSvcErr) {
	unaryTopicManagers := make([]*grpcmanagers.TopicGrpcManager, 0)
	for i := 0; uint32(i) < numUnaryChannels; i++ {
		unaryTopicManager, err := grpcmanagers.NewStreamTopicGrpcManager(request)
		if err != nil {
			return nil, err
		}
		unaryTopicManagers = append(unaryTopicManagers, unaryTopicManager)
	}
	return &staticUnaryGrpcManagerPool{
		grpcManagers: unaryTopicManagers,
		logger:       logger,
	}, nil
}
