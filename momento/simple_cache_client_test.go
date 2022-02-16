package momento

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

var (
	TestAuthToken = os.Getenv("TEST_AUTH_TOKEN")
	TestCacheName = os.Getenv("TEST_CACHE_NAME")
)

const (
	defaultTtlSeconds = 60
)

func setUp() (ScsClient, error) {
	if TestAuthToken == "" {
		return nil, errors.New("integration tests require TEST_CACHE_NAME env var")
	}
	if TestCacheName == "" {
		return nil, errors.New("integration tests require TEST_CACHE_NAME env var")
	}

	client, err := NewSimpleCacheClient(&SimpleCacheClientRequest{
		AuthToken:         TestAuthToken,
		DefaultTtlSeconds: defaultTtlSeconds,
	})
	if err != nil {
		return nil, err
	}

	// Check if TestCacheName exists
	err = client.CreateCache(&CreateCacheRequest{
		CacheName: TestCacheName,
	})
	if momentoErr, ok := err.(MomentoError); ok {
		if momentoErr.Code() != AlreadyExists {
			return nil, err
		}
	}
	return client, nil
}

func cleanUp(client ScsClient) {
	_ = client.Close()
}

// Basic happy path test - create a cache, operate set/get, and delete the cache
func TestBasicHappyPathSDKFlow(t *testing.T) {
	cacheName := uuid.NewString()
	key := []byte(uuid.NewString())
	value := []byte(uuid.NewString())
	client, err := setUp()
	if err != nil {
		t.Error(fmt.Errorf("error occured setting up client err=%+v", err))
	}
	err = client.CreateCache(&CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		t.Error(fmt.Errorf("error occured creating cache err=%+v", err))
	}

	_, err = client.Set(&CacheSetRequest{
		CacheName: cacheName,
		Key:       key,
		Value:     value,
	})
	if err != nil {
		t.Errorf("error occured setting key err=%+v", err)
	}

	getResp, err := client.Get(&CacheGetRequest{
		CacheName: cacheName,
		Key:       key,
	})
	if err != nil {
		t.Errorf("error occured getting key err=%+v", err)
		return
	}

	if getResp.Result() != HIT {
		t.Errorf("unexpected result when getting test key got=%+v expected=%+v", getResp.Result(), HIT)
	}
	if !bytes.Equal(getResp.ByteValue(), value) {
		t.Errorf(
			"set byte value and returned byte value are not equal "+
				"got=%+v expected=%+v", getResp.ByteValue(), value,
		)
	}

	existingCacheResp, err := client.Get(&CacheGetRequest{
		CacheName: TestCacheName,
		Key:       key,
	})
	if err != nil {
		t.Error(err.Error())
	}

	if existingCacheResp.Result() != MISS {
		t.Errorf(
			"key: %s shouldn't exist in %s since it's never set. got=%s", string(key),
			TestCacheName, existingCacheResp.StringValue(),
		)
	}

	err = client.DeleteCache(&DeleteCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		t.Errorf("error occured deleting cache err=%+v", err)
	}

	cleanUp(client)
}

func TestClientInitialization(t *testing.T) {

	tests := map[string]struct {
		expectedErr string
		req         *SimpleCacheClientRequest
	}{
		"happy path": {
			req: &SimpleCacheClientRequest{
				AuthToken: TestAuthToken,
			},
		},
		"happy path custom timeout": {
			req: &SimpleCacheClientRequest{
				AuthToken:             TestAuthToken,
				RequestTimeoutSeconds: 100,
			},
		},
		"happy path custom timeout and ttl": {
			req: &SimpleCacheClientRequest{
				AuthToken:             TestAuthToken,
				RequestTimeoutSeconds: 100,
				DefaultTtlSeconds:     defaultTtlSeconds,
			},
		},
		"test invalid auth token": {
			expectedErr: ClientSdkError,
			req: &SimpleCacheClientRequest{
				AuthToken:         "NOT_A_VALID_JWT",
				DefaultTtlSeconds: defaultTtlSeconds,
			},
		},
	}
	for name, tt := range tests {
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := NewSimpleCacheClient(tt.req)
			if tt.expectedErr != "" && err == nil {
				t.Errorf("expected error but got none expected=%+v got=%+v", tt.expectedErr, err)
			}

			if tt.expectedErr != "" && err != nil {
				if momentoErr, ok := err.(MomentoError); ok {
					if momentoErr.Code() == tt.expectedErr {
						return // Success end test we expected this
					}
				}
				t.Errorf(
					"unexpected error occured initilizing client got=%+v expected=%+v",
					err, tt.expectedErr,
				)
			}

			if tt.expectedErr == "" && err != nil {
				t.Errorf("unexpected error occured on init expected=%+v got=%+v", tt.expectedErr, err)
			}
		})
	}
}
