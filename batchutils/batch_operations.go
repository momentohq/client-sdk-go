package batchutils

import (
	"context"
	"github.com/momentohq/client-sdk-go/momento"
)

const maxConcurrentGets = 5

func keyDistributor(ctx context.Context, keys []momento.Key, keyChan chan momento.Key) {
	for _, k := range keys {
		keyChan <- k
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			keyChan <- nil
		}
	}
}

type errKeyVal struct {
	key   momento.Value
	error error
}
