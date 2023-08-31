package batchutils

import (
	"context"
	"sync"
	"time"

	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

const maxConcurrentSets = 5

type setKeyResp struct {
	key  momento.Value
	resp responses.SetResponse
}

func setWorker(
	ctx context.Context,
	client momento.CacheClient,
	cacheName string,
	itemChan chan BatchSetItem,
	resultChan chan *setResultOrError,
) {
	for {
		item := <-itemChan
		if item.Key == nil {
			return
		}
		setResponse, err := client.Set(ctx, &momento.SetRequest{
			CacheName: cacheName,
			Key:       item.Key,
			Value:     item.Value,
			Ttl:       item.Ttl,
		})

		if err != nil {
			resultChan <- &setResultOrError{err: &errKeyVal{
				key:   item.Key,
				error: err,
			}}
		} else {
			resultChan <- &setResultOrError{result: &setKeyResp{
				key:  item.Key,
				resp: setResponse,
			}}
		}
	}
}

type BatchSetRequest struct {
	Client            momento.CacheClient
	CacheName         string
	Items             []BatchSetItem
	MaxConcurrentSets int
}

type BatchSetItem struct {
	Key   momento.Key
	Value momento.Value
	Ttl   time.Duration
}

type setResultOrError struct {
	result *setKeyResp
	err    *errKeyVal
}

// BatchSetError contains a map associating failing cache keys with their specific errors.
// It may be necessary to use a type assertion to access the errors:
//
// errors := (err.(*BatchSetError)).Errors()
type BatchSetError struct {
	errors map[momento.Value]error
}

func (e *BatchSetError) Error() string {
	return "Errors occurred during batch set; call Errors() to get a map of key -> errorType"
}

// Errors contains a map associating unsuccessful set errors with their cache keys.
func (e *BatchSetError) Errors() map[momento.Value]error {
	return e.errors
}

// BatchSetResponse contains a map associating successful set responses with their cache keys.
type BatchSetResponse struct {
	responses map[momento.Value]responses.SetResponse
}

func (e *BatchSetResponse) Responses() map[momento.Value]responses.SetResponse {
	return e.responses
}

func itemDistributor(ctx context.Context, items []BatchSetItem, itemChan chan BatchSetItem) {
	for _, item := range items {
		itemChan <- item
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			itemChan <- BatchSetItem{}
		}
	}
}

// BatchSet sets a slice of keys to the cache, returning a map from failing cache keys to their specific errors.
func BatchSet(ctx context.Context, props *BatchSetRequest) (*BatchSetResponse, *BatchSetError) {
	// initialize return value
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	// stop the key distributor when we return
	defer cancelFunc()
	var wg sync.WaitGroup

	if props.MaxConcurrentSets == 0 {
		props.MaxConcurrentSets = maxConcurrentSets
	}
	if len(props.Items) < props.MaxConcurrentSets {
		props.MaxConcurrentSets = len(props.Items)
	}
	itemChan := make(chan BatchSetItem, props.MaxConcurrentSets)
	resultChan := make(chan *setResultOrError, len(props.Items))

	for i := 0; i < props.MaxConcurrentSets; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			setWorker(ctx, props.Client, props.CacheName, itemChan, resultChan)
		}()
	}

	go itemDistributor(cancelCtx, props.Items, itemChan)

	// wait for the workers to return
	wg.Wait()

	var errors = make(map[momento.Value]error, 0)
	var results = make(map[momento.Value]responses.SetResponse, 0)

	for i := 0; i < len(props.Items); i++ {
		resOrErr := <-resultChan
		if resOrErr.result != nil {
			results[resOrErr.result.key] = resOrErr.result.resp
		} else if resOrErr.err != nil {
			errors[resOrErr.err.key] = resOrErr.err.error
		}
	}

	var batchSetResponses *BatchSetResponse
	var batchSetErrors *BatchSetError

	if len(results) == 0 {
		batchSetResponses = nil
	} else {
		batchSetResponses = &BatchSetResponse{responses: results}
	}

	if len(errors) == 0 {
		batchSetErrors = nil
	} else {
		batchSetErrors = &BatchSetError{errors: errors}
	}
	return batchSetResponses, batchSetErrors
}
