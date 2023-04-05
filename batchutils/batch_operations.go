package batchutils

import (
	"context"
	"fmt"
	"sync"

	"github.com/momentohq/client-sdk-go/momento"
)

const maxConcurrentDeletes = 5

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

func deleteWorker(
	ctx context.Context,
	client momento.CacheClient,
	cacheName string,
	keyChan chan momento.Key,
	errChan chan string,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case myKey := <-keyChan:
			if myKey == nil {
				return
			}
			_, err := client.Delete(ctx, &momento.DeleteRequest{
				CacheName: cacheName,
				Key:       myKey,
			})
			if err != nil {
				errChan <- fmt.Sprintf("error deleting key %s: %s", myKey, err.Error())
			} else {
				errChan <- ""
			}
		}
	}
}

type BatchDeleteRequest struct {
	Client               momento.CacheClient
	CacheName            string
	Keys                 []momento.Key
	MaxConcurrentDeletes int
}

// BatchDelete deletes a slice of keys from the cache, returning an array containing messages from any delete errors
func BatchDelete(ctx context.Context, props *BatchDeleteRequest) []string {
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
	errChan := make(chan string, len(props.Keys))

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

	var errors = make([]string, 0)
	for i := 0; i < len(props.Keys); i++ {
		msg := <-errChan
		if msg != "" {
			errors = append(errors, msg)
		}
	}
	return errors
}
