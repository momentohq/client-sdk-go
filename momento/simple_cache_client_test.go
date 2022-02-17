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
	testAuthToken = os.Getenv("TEST_AUTH_TOKEN")
	testCacheName = os.Getenv("TEST_CACHE_NAME")
)

const (
	defaultTtlSeconds = 60
)

// Basic happy path test - create a cache, operate set/get, and delete the cache
func TestBasicHappyPathSDKFlow(t *testing.T) {
	cacheName := uuid.NewString()
	key := []byte(uuid.NewString())
	value := []byte(uuid.NewString())
	client, err := newTestClient()
	if err != nil {
		t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
	}
	err = client.CreateCache(&CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		t.Error(fmt.Errorf("error occurred creating cache err=%+v", err))
	}

	_, err = client.Set(&CacheSetRequest{
		CacheName: cacheName,
		Key:       key,
		Value:     value,
	})
	if err != nil {
		t.Errorf("error occurred setting key err=%+v", err)
	}

	_, err = client.Set(&CacheSetRequest{
		CacheName:  cacheName,
		Key:        uuid.NewString(),
		Value:      value,
		TtlSeconds: TTL(1),
	})
	if err != nil {
		t.Errorf("error occurred setting key with custom ttl err=%+v", err)
	}

	getResp, err := client.Get(&CacheGetRequest{
		CacheName: cacheName,
		Key:       key,
	})
	if err != nil {
		t.Errorf("error occurred getting key err=%+v", err)
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
		CacheName: testCacheName,
		Key:       key,
	})
	if err != nil {
		t.Error(err.Error())
	}

	if existingCacheResp.Result() != MISS {
		t.Errorf(
			"key: %s shouldn't exist in %s since it's never set. got=%s", string(key),
			testCacheName, existingCacheResp.StringValue(),
		)
	}

	err = client.DeleteCache(&DeleteCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		t.Errorf("error occurred deleting cache err=%+v", err)
	}

	cleanUpClient(client)
}

func TestClientInitialization(t *testing.T) {

	testRequestTimeout := uint32(100)
	badRequestTimeout := uint32(0)
	tests := map[string]struct {
		expectedErr           string
		authToken             string
		defaultTtlSeconds     uint32
		requestTimeoutSeconds *uint32
	}{
		"happy path": {
			authToken:         testAuthToken,
			defaultTtlSeconds: defaultTtlSeconds,
		},
		"happy path custom timeout": {
			authToken:             testAuthToken,
			defaultTtlSeconds:     defaultTtlSeconds,
			requestTimeoutSeconds: &testRequestTimeout,
		},
		"test invalid auth token": {
			expectedErr:       ClientSdkError,
			authToken:         "NOT_A_VALID_JWT",
			defaultTtlSeconds: defaultTtlSeconds,
		},
		"test invalid default ttl": {
			expectedErr:       InvalidArgumentError,
			authToken:         testAuthToken,
			defaultTtlSeconds: 0,
		},
		"test invalid request timeout": {
			expectedErr:           InvalidArgumentError,
			authToken:             testAuthToken,
			defaultTtlSeconds:     defaultTtlSeconds,
			requestTimeoutSeconds: &badRequestTimeout,
		},
	}
	for name, tt := range tests {
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c, err := NewSimpleCacheClient(tt.authToken, tt.defaultTtlSeconds)
			if tt.requestTimeoutSeconds != nil {
				c, err = NewSimpleCacheClient(tt.authToken, tt.defaultTtlSeconds, WithRequestTimeout(*tt.requestTimeoutSeconds))
			}
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
					"unexpected error occurred initializing client got=%+v expected=%+v",
					err, tt.expectedErr,
				)
			}

			if tt.expectedErr == "" && err != nil {
				t.Errorf("unexpected error occurred on init expected=%+v got=%+v", tt.expectedErr, err)
			}
			cleanUpClient(c)
		})
	}
}

func newTestClient() (ScsClient, error) {
	if testAuthToken == "" {
		return nil, errors.New("integration tests require TEST_CACHE_NAME env var")
	}
	if testCacheName == "" {
		return nil, errors.New("integration tests require TEST_CACHE_NAME env var")
	}

	client, err := NewSimpleCacheClient(testAuthToken, defaultTtlSeconds)
	if err != nil {
		return nil, err
	}

	// Check if testCacheName exists
	err = client.CreateCache(&CreateCacheRequest{
		CacheName: testCacheName,
	})
	if momentoErr, ok := err.(MomentoError); ok {
		if momentoErr.Code() != AlreadyExists {
			return nil, err
		}
	}
	return client, nil
}

func cleanUpClient(client ScsClient) {
	client.Close()
}
