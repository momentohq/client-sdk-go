package batchutils

import (
	"context"
	"errors"
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

// BatchSetIfNotExists will set the key-value pairs ONLY if all keys don't already exist
func BatchSetIfNotExists(ctx context.Context, props *BatchSetIfNotExistsRequest) (*BatchSetResponse, *BatchSetError, error) {
	// First check if all keys exist or not
	resp, err := props.Client.KeysExist(ctx, &momento.KeysExistRequest{
		CacheName: props.CacheName,
		Keys:      ExtractKeys(props.Items),
	})
	if err != nil {
		return nil, nil, err
	}

	switch result := resp.(type) {
	case *responses.KeysExistSuccess:
		// Check if any of the keys already exist
		for _, keyExists := range result.Exists() {
			if keyExists {
				return nil, nil, momento.NewMomentoError(momento.AlreadyExistsError, "At least one key already exists", errors.New("at least one key already exists"))
			}
		}
	default:
		return nil, nil, err
	}

	// If none of the keys exist, set the items using BatchSet
	setBatchResponse, setBatchError := BatchSet(ctx, &BatchSetRequest{
		Client:    props.Client,
		CacheName: props.CacheName,
		Items:     props.Items,
	})
	return setBatchResponse, setBatchError, nil
}
