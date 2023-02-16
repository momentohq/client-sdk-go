package momento

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"

	"github.com/google/uuid"
)

var (
	testCacheName             = os.Getenv("TEST_CACHE_NAME")
	testCredentialProvider, _ = newCredentialProvider("TEST_AUTH_TOKEN")
)

const (
	defaultTTL = time.Second * 3
)

func TestSet(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		cacheName   string
		key         string
		value       string
		expectedErr string
	}{
		"test set on non existent cache": {
			cacheName:   uuid.NewString(),
			key:         uuid.NewString(),
			value:       uuid.NewString(),
			expectedErr: NotFoundError,
		},
		"test set on empty cache name": {
			cacheName:   "",
			key:         uuid.NewString(),
			value:       uuid.NewString(),
			expectedErr: InvalidArgumentError,
		},
		"test set on nil key": {
			cacheName:   testCacheName,
			key:         "",
			value:       uuid.NewString(),
			expectedErr: InvalidArgumentError,
		},
		"test set on nil value": {
			cacheName:   testCacheName,
			key:         uuid.NewString(),
			value:       "",
			expectedErr: InvalidArgumentError,
		},
	}
	for name, tt := range tests {
		client, err := newTestClient(testCredentialProvider)
		defer teardown(client)

		if err != nil {
			t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
		}
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			_, err := client.Set(ctx, &SetRequest{
				CacheName: tt.cacheName,
				Key:       String(tt.key),
				Value:     String(tt.value),
			})
			if tt.expectedErr != "" && err == nil {
				t.Errorf("expected error but got none expected=%+v got=%+v", tt.expectedErr, err)
			}

			if tt.expectedErr != "" && err != nil {
				var momentoErr MomentoError
				if errors.As(err, &momentoErr) {
					if momentoErr.Code() == tt.expectedErr {
						return // Success end test we expected this
					}
				}
				t.Errorf(
					"unexpected error occurred setting cache got=%+v expected=%+v",
					err, tt.expectedErr,
				)
			}
		})
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		cacheName   string
		key         string
		expectedErr string
	}{
		"test get on non existent cache": {
			cacheName:   uuid.NewString(),
			key:         uuid.NewString(),
			expectedErr: NotFoundError,
		},
		"test get on empty cache name": {
			cacheName:   "",
			key:         uuid.NewString(),
			expectedErr: InvalidArgumentError,
		},
		"test get on empty key": {
			cacheName:   uuid.NewString(),
			key:         "",
			expectedErr: InvalidArgumentError,
		},
	}
	for name, tt := range tests {
		client, err := newTestClient(testCredentialProvider)
		defer teardown(client)

		if err != nil {
			t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
		}
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			_, err := client.Get(ctx, &GetRequest{
				CacheName: tt.cacheName,
				Key:       String(tt.key),
			})
			if tt.expectedErr != "" && err == nil {
				t.Errorf("expected error but got none expected=%+v got=%+v", tt.expectedErr, err)
			}

			if tt.expectedErr != "" && err != nil {
				var momentoErr MomentoError
				if errors.As(err, &momentoErr) {
					if momentoErr.Code() == tt.expectedErr {
						return // Success end test we expected this
					}
				}
				t.Errorf(
					"unexpected error occurred setting cache got=%+v expected=%+v",
					err, tt.expectedErr,
				)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		cacheName   string
		key         string
		expectedErr string
	}{
		"test delete on non existent cache": {
			cacheName:   uuid.NewString(),
			key:         uuid.NewString(),
			expectedErr: NotFoundError,
		},
		"test delete on empty cache name": {
			cacheName:   "",
			key:         uuid.NewString(),
			expectedErr: InvalidArgumentError,
		},
		"test delete on bad key": {
			cacheName:   testCacheName,
			key:         "",
			expectedErr: InvalidArgumentError,
		},
	}
	for name, tt := range tests {
		client, err := newTestClient(testCredentialProvider)
		defer teardown(client)

		if err != nil {
			t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
		}
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			_, err = client.Delete(ctx, &DeleteRequest{
				CacheName: tt.cacheName,
				Key:       String(tt.key),
			})
			if tt.expectedErr != "" && err == nil {
				t.Errorf("expected error but got none expected=%+v got=%+v", tt.expectedErr, err)
			}

			if tt.expectedErr != "" && err != nil {
				var momentoErr MomentoError
				if errors.As(err, &momentoErr) {
					if momentoErr.Code() == tt.expectedErr {
						return // Success end test we expected this
					}
				}
				t.Errorf(
					"unexpected error occurred calling cache delete got=%+v expected=%+v",
					err, tt.expectedErr,
				)
			}
		})
	}
}

func newCredentialProvider(envVarName string) (auth.CredentialProvider, error) {
	credentialProvider, err := auth.NewEnvMomentoTokenProvider(envVarName)
	if err != nil {
		return nil, err
	}
	return credentialProvider, nil
}

func newTestClient(credentialProvider auth.CredentialProvider) (SimpleCacheClient, error) {
	ctx := context.Background()
	if testCacheName == "" {
		return nil, errors.New("integration tests require TEST_CACHE_NAME env var")
	}

	client, err := NewSimpleCacheClient(&SimpleCacheClientProps{
		Configuration:      config.LatestLaptopConfig(),
		CredentialProvider: credentialProvider,
		DefaultTTL:         defaultTTL,
	})
	if err != nil {
		return nil, err
	}

	// Check if testCacheName exists
	_, err = client.CreateCache(ctx, &CreateCacheRequest{
		CacheName: testCacheName,
	})
	var momentoErr MomentoError
	if errors.As(err, &momentoErr) {
		if momentoErr.Code() != AlreadyExistsError {
			return nil, err
		}
	}
	return client, nil
}

func teardown(client SimpleCacheClient, cacheNames ...string) {
	ctx := context.Background()

	for _, cacheName := range cacheNames {
		_, err := client.DeleteCache(ctx, &DeleteCacheRequest{
			CacheName: cacheName,
		})

		// It's ok if the cache doesn't exist by the time
		// we're tearing it down. Makes tearing down safer
		// to just throw all possible caches at it.
		var momentoErr MomentoError
		if errors.As(err, &momentoErr) {
			if momentoErr.Code() != NotFoundError {
				panic(err)
			}
		}
	}

	if client != nil {
		client.Close()
	}
}
