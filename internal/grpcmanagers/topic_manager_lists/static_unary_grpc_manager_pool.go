package topic_manager_lists

import (
	"sync"
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
	closeOnce    sync.Once
}

// GetNextTopicGrpcManager returns the next manager via round-robin.
// Reservation.Release is a no-op since unary requests don't hold slots.
func (list *staticUnaryGrpcManagerPool) GetNextTopicGrpcManager() (*Reservation, momentoerrors.MomentoSvcErr) {
	nextManagerIndex := list.managerIndex.Add(1)
	return NewReservation(list.grpcManagers[nextManagerIndex%uint64(len(list.grpcManagers))], list.releaseManager), nil
}

func (list *staticUnaryGrpcManagerPool) releaseManager(_ *grpcmanagers.TopicGrpcManager) int64 {
	return 0
}

// Close shuts down all gRPC connections. Safe to call multiple times.
func (list *staticUnaryGrpcManagerPool) Close() {
	list.closeOnce.Do(func() {
		closeAllManagers(list.grpcManagers, list.logger)
	})
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
