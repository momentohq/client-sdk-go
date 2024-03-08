package batchutils

import (
	"context"
	"sync"
	"time"

	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

const maxConcurrentGets = 5

type getKeyResp struct {
	key  momento.Value
	resp responses.GetResponse
}

func getWorker(
	ctx context.Context,
	client momento.CacheClient,
	cacheName string,
	keyChan chan momento.Key,
	resultChan chan *getResultOrError,
	timeout time.Duration,
) {
	for {
		myKey := <-keyChan
		if myKey == nil {
			return
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)

		getResponse, err := client.Get(timeoutCtx, &momento.GetRequest{
			CacheName: cacheName,
			Key:       myKey,
		})

		cancel()

		if err != nil {
			resultChan <- &getResultOrError{err: &errKeyVal{
				key:   myKey,
				error: err,
			}}
		} else {
			resultChan <- &getResultOrError{result: &getKeyResp{
				key:  myKey,
				resp: getResponse,
			}}
		}
	}
}

type BatchGetRequest struct {
	Client            momento.CacheClient
	CacheName         string
	Keys              []momento.Key
	MaxConcurrentGets int
	// timeout for individual requests, defaults to 10 seconds
	RequestTimeout *time.Duration
}

// BatchGetError contains a map associating failing cache keys with their specific errors.
// It may be necessary to use a type assertion to access the errors:
//
// errors := (err.(*BatchGetError)).Errors()
type BatchGetError struct {
	errors map[momento.Value]error
}

type getResultOrError struct {
	result *getKeyResp
	err    *errKeyVal
}

func (e *BatchGetError) Error() string {
	return "errors occurred during batch delete"
}

func (e *BatchGetError) Errors() map[momento.Value]error {
	return e.errors
}

// BatchGetResponse contains a map associating successful get responses with their cache keys.
type BatchGetResponse struct {
	responses map[momento.Value]responses.GetResponse
}

func (e *BatchGetResponse) Responses() map[momento.Value]responses.GetResponse {
	return e.responses
}

// BatchGet gets a slice of keys from the cache, returning a map from failing cache keys to their specific errors.
func BatchGet(ctx context.Context, props *BatchGetRequest) (*BatchGetResponse, *BatchGetError) {
	// initialize return value
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	// stop the key distributor when we return
	defer cancelFunc()
	var wg sync.WaitGroup

	if props.MaxConcurrentGets == 0 {
		props.MaxConcurrentGets = maxConcurrentGets
	}
	if len(props.Keys) < props.MaxConcurrentGets {
		props.MaxConcurrentGets = len(props.Keys)
	}
	keyChan := make(chan momento.Key, props.MaxConcurrentGets)
	resultChan := make(chan *getResultOrError, len(props.Keys))

	for i := 0; i < props.MaxConcurrentGets; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			getWorker(ctx, props.Client, props.CacheName, keyChan, resultChan, getRequestTimeout(props.RequestTimeout))
		}()
	}

	go keyDistributor(cancelCtx, props.Client.Logger(), props.MaxConcurrentGets, props.Keys, keyChan)

	// wait for the workers to return
	wg.Wait()

	var errors = make(map[momento.Value]error, 0)
	var results = make(map[momento.Value]responses.GetResponse, 0)
	for i := 0; i < len(props.Keys); i++ {
		resOrError := <-resultChan
		if resOrError.result != nil {
			results[resOrError.result.key] = resOrError.result.resp
		} else if resOrError.err != nil {
			errors[resOrError.err.key] = resOrError.err.error
		}
	}

	var batchGetResponses *BatchGetResponse
	var batchGetErrors *BatchGetError

	if len(results) == 0 {
		batchGetResponses = nil
	} else {
		batchGetResponses = &BatchGetResponse{responses: results}
	}

	if len(errors) == 0 {
		batchGetErrors = nil
	} else {
		batchGetErrors = &BatchGetError{errors: errors}
	}
	return batchGetResponses, batchGetErrors
}
