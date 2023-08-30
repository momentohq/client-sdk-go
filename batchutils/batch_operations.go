package batchutils

import (
	"context"
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
