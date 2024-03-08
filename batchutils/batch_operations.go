package batchutils

import (
	"context"
	"github.com/momentohq/client-sdk-go/config/logger"
	"time"

	"github.com/momentohq/client-sdk-go/momento"
)

const defaultRequestTimeout = 10 * time.Second

func getRequestTimeout(propsTimeout *time.Duration) (requestTimeout time.Duration) {
	if propsTimeout == nil {
		requestTimeout = defaultRequestTimeout
	} else {
		requestTimeout = *propsTimeout
	}
	return
}

func keyDistributor(ctx context.Context, logger logger.MomentoLogger, numWorkers int, keys []momento.Key, keyChan chan momento.Key) {
	for _, k := range keys {
		keyChan <- k
	}

	logger.Trace("keyDistributor has put all of the keys on the channel")

	// after we have put all the keys onto the channel, we add one nil for each worker to signal that they should exit
	for i := 0; i < numWorkers; i++ {
		keyChan <- nil
	}

	logger.Trace("keyDistributor has put a nil on the channel for each worker")

	for range ctx.Done() {
		logger.Trace("keyDistributor context done, exiting for loop")
		return
	}
}

type errKeyVal struct {
	key   momento.Value
	error error
}
