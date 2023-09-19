package batchutils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

type BatchSetIfNotExistsRequest struct {
	Client            momento.CacheClient
	CacheName         string
	Items             []BatchSetItem
	MaxConcurrentSets int
	// timeout for individual requests, defaults to 10 seconds
	RequestTimeout *time.Duration
}

type Items []BatchSetItem

func ExtractKeys(items Items) []momento.Key {
	var list []momento.Key
	for _, item := range items {
		list = append(list, item.Key)
	}
	return list
}

// BatchSetIfNotExistsError contains a map associating failing cache keys with their specific errors.
// Errors may be related to keys already existing or batch set failing.
// It may be necessary to use a type assertion to access the errors:
//
// errors := (err.(*BatchSetIfNotExistsError)).Errors()
type BatchSetIfNotExistsError struct {
	errors map[momento.Key]error
}

func (e *BatchSetIfNotExistsError) Error() string {
	return "Errors occurred during BatchSetIfNotExists; call Errors() to get a map of key -> errorType"
}

// Errors contains a map associating errors with their cache keys.
func (e *BatchSetIfNotExistsError) Errors() map[momento.Key]error {
	return e.errors
}

func CreateErrorMapping(keys []momento.Key, err error) map[momento.Key]error {
	var errors = make(map[momento.Key]error, len(keys))
	for i := 0; i < len(keys); i++ {
		errors[keys[i]] = err
	}
	return errors
}

// BatchSetIfNotExists will set the key-value pairs ONLY if all keys don't already exist
func BatchSetIfNotExists(ctx context.Context, props *BatchSetIfNotExistsRequest) (*BatchSetResponse, *BatchSetIfNotExistsError) {
	// First check if all keys exist or not
	keys := ExtractKeys(props.Items)
	resp, err := props.Client.KeysExist(ctx, &momento.KeysExistRequest{
		CacheName: props.CacheName,
		Keys:      keys,
	})
	if err != nil {
		var keysExistsError = &BatchSetIfNotExistsError{errors: CreateErrorMapping(keys, err)}
		return nil, keysExistsError
	}

	switch result := resp.(type) {
	case *responses.KeysExistSuccess:
		// Check if any of the keys already exist
		for _, keyExists := range result.Exists() {
			if keyExists {
				var momentoError = momento.NewMomentoError(momento.AlreadyExistsError, "At least one key already exists", errors.New("at least one key already exists"))
				var atLeastOneKeyExistsError = &BatchSetIfNotExistsError{errors: CreateErrorMapping(keys, momentoError)}
				return nil, atLeastOneKeyExistsError
			}
		}
	default:
		var message = fmt.Sprintf("Unexpected KeysExistResponse type: %T\n", resp)
		var unexpectedError = momento.NewMomentoError(momento.ClientSdkError, message, errors.New(message))
		var unexpectedResponseError = &BatchSetIfNotExistsError{errors: CreateErrorMapping(keys, unexpectedError)}
		return nil, unexpectedResponseError
	}

	// If none of the keys exist, set the items using BatchSet
	setBatchResponse, setBatchError := BatchSet(ctx, &BatchSetRequest{
		Client:    props.Client,
		CacheName: props.CacheName,
		Items:     props.Items,
	})
	if setBatchError != nil {
		return nil, &BatchSetIfNotExistsError{errors: setBatchError.Errors()}
	}
	return setBatchResponse, nil
}
