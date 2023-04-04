package batchutils

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/momentohq/client-sdk-go/momento"
)

const (
	defaultNumWorkers          = 50
	defaultMaxRequestPerSecond = 100
)

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
	id int,
	client momento.CacheClient,
	cacheName string,
	workerDelayBetweenRequests float64,
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
			time.Sleep(time.Millisecond * time.Duration(int64(workerDelayBetweenRequests*1000)))
		}
	}
}

type BatchDeleteRequest struct {
	Client               momento.CacheClient
	CacheName            string
	Keys                 []momento.Key
	NumWorkers           int
	MaxRequestsPerSecond int
}

// BatchDelete deletes a slice of keys from the cache, returning an array containing messages from any delete errors
func BatchDelete(ctx context.Context, props *BatchDeleteRequest) []string {
	// initialize return value
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	// stop the key distributor when we return
	defer cancelFunc()
	var wg sync.WaitGroup

	if props.NumWorkers == 0 {
		props.NumWorkers = defaultNumWorkers
	}
	if props.MaxRequestsPerSecond == 0 {
		props.MaxRequestsPerSecond = defaultMaxRequestPerSecond
	}
	if len(props.Keys) < props.NumWorkers {
		props.NumWorkers = len(props.Keys)
	}
	if props.NumWorkers > props.MaxRequestsPerSecond {
		props.NumWorkers = props.MaxRequestsPerSecond
	}

	workerDelayBetweenRequests := float64(props.NumWorkers) / float64(props.MaxRequestsPerSecond)
	keyChan := make(chan momento.Key, props.NumWorkers)
	errChan := make(chan string, len(props.Keys))

	for i := 0; i < props.NumWorkers; i++ {
		wg.Add(1)

		// avoid reuse of the same "i" value in each closure by passing it in to the goroutine
		go func(i int) {
			defer wg.Done()
			deleteWorker(ctx, i, props.Client, props.CacheName, workerDelayBetweenRequests, keyChan, errChan)
		}(i)
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
