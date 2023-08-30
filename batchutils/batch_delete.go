package batchutils

import (
	"context"
	"sync"

	"github.com/momentohq/client-sdk-go/momento"
)

const maxConcurrentDeletes = 5

func deleteWorker(
	ctx context.Context,
	client momento.CacheClient,
	cacheName string,
	keyChan chan momento.Key,
	errChan chan *errKeyVal,
) {
	for {
		myKey := <-keyChan
		if myKey == nil {
			return
		}
		_, err := client.Delete(ctx, &momento.DeleteRequest{
			CacheName: cacheName,
			Key:       myKey,
		})
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
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	// stop the key distributor when we return
	defer cancelFunc()
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
			deleteWorker(ctx, props.Client, props.CacheName, keyChan, errChan)
		}()
	}

	go keyDistributor(cancelCtx, props.Keys, keyChan)

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
