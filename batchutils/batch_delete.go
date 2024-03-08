package batchutils

import (
	"context"
	"sync"
	"time"

	"github.com/momentohq/client-sdk-go/momento"
)

const maxConcurrentDeletes = 5

func deleteWorker(
	ctx context.Context,
	client momento.CacheClient,
	cacheName string,
	keyChan chan momento.Key,
	errChan chan *errKeyVal,
	timeout time.Duration,
) {
	for {
		myKey := <-keyChan
		if myKey == nil {
			return
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)

		_, err := client.Delete(timeoutCtx, &momento.DeleteRequest{
			CacheName: cacheName,
			Key:       myKey,
		})

		cancel()

		if err != nil {
			errChan <- &errKeyVal{
				key:   myKey,
				error: err,
			}
		} else {
			errChan <- nil
		}
	}
}

type BatchDeleteRequest struct {
	Client               momento.CacheClient
	CacheName            string
	Keys                 []momento.Key
	MaxConcurrentDeletes int
	// timeout for individual requests, defaults to 10 seconds
	RequestTimeout *time.Duration
}

// BatchDeleteError contains a map associating failing cache keys with their specific errors.
// It may be necessary to use a type assertion to access the errors:
//
// errors := (err.(*BatchDeleteError)).Errors()
type BatchDeleteError struct {
	errors map[momento.Value]error
}

func (e *BatchDeleteError) Error() string {
	return "errors occurred during batch delete"
}

func (e *BatchDeleteError) Errors() map[momento.Value]error {
	return e.errors
}

// BatchDelete deletes a slice of keys from the cache, returning a map from failing cache keys to their specific errors.
func BatchDelete(ctx context.Context, props *BatchDeleteRequest) *BatchDeleteError {
	// initialize return value
	var wg sync.WaitGroup

	if props.MaxConcurrentDeletes == 0 {
		props.MaxConcurrentDeletes = maxConcurrentDeletes
	}
	if len(props.Keys) < props.MaxConcurrentDeletes {
		props.MaxConcurrentDeletes = len(props.Keys)
	}
	keyChan := make(chan momento.Key, props.MaxConcurrentDeletes)
	errChan := make(chan *errKeyVal, len(props.Keys))

	for i := 0; i < props.MaxConcurrentDeletes; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			deleteWorker(ctx, props.Client, props.CacheName, keyChan, errChan, getRequestTimeout(props.RequestTimeout))
		}()
	}

	for _, k := range props.Keys {
		keyChan <- k
	}

	props.Client.Logger().Trace("BatchDelete: put all of the keys on the channel")

	// after we have put all the keys onto the channel, we add one nil for each worker to signal that they should exit
	for i := 0; i < props.MaxConcurrentDeletes; i++ {
		keyChan <- nil
	}

	props.Client.Logger().Trace("BatchDelete: put a nil on the channel for each worker")

	// wait for the workers to return
	wg.Wait()

	var errors = make(map[momento.Value]error, 0)
	for i := 0; i < len(props.Keys); i++ {
		msg := <-errChan
		if msg != nil {
			errors[msg.key] = msg.error
		}
	}

	if len(errors) == 0 {
		return nil
	}
	return &BatchDeleteError{errors: errors}
}
