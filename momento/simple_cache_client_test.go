package momento

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	testAuthToken = os.Getenv("TEST_AUTH_TOKEN")
	testCacheName = os.Getenv("TEST_CACHE_NAME")
)

const (
	defaultTtlSeconds = 3
)

// Basic happy path test - create a cache, operate set/get, and delete the cache
func TestBasicHappyPathSDKFlow(t *testing.T) {
	ctx := context.Background()
	cacheName := uuid.NewString()
	key := []byte(uuid.NewString())
	value := []byte(uuid.NewString())
	client, err := newTestClient()
	if err != nil {
		t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
	}
	err = client.CreateCache(ctx, &CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		t.Error(fmt.Errorf("error occurred creating cache err=%+v", err))
	}

	_, err = client.Set(ctx, &CacheSetRequest{
		CacheName: cacheName,
		Key:       key,
		Value:     value,
	})
	if err != nil {
		t.Errorf("error occurred setting key err=%+v", err)
	}

	_, err = client.Set(ctx, &CacheSetRequest{
		CacheName:  cacheName,
		Key:        uuid.NewString(),
		Value:      value,
		TtlSeconds: TTL(1),
	})
	if err != nil {
		t.Errorf("error occurred setting key with custom ttl err=%+v", err)
	}

	getResp, err := client.Get(ctx, &CacheGetRequest{
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

	existingCacheResp, err := client.Get(ctx, &CacheGetRequest{
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

	err = client.DeleteCache(ctx, &DeleteCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		t.Error(fmt.Errorf("error occurred deleting cache=%s err=%+v", cacheName, err))
	}

	cleanUpClient(client)
}

func TestBasicHappyPathDelete(t *testing.T) {
	ctx := context.Background()ctx
	cacheName := uuid.NewString()
	key := []byte(uuid.NewString())
	value := []byte(uuid.NewString())
	client, err := newTestClient()
	if err != nil {
		t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
	}
	err = client.CreateCache(ctx, &CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		t.Error(fmt.Errorf("error occurred creating cache err=%+v", err))
	}

	_, err = client.Set(ctx, &CacheSetRequest{
		CacheName: cacheName,
		Key:       key,
		Value:     value,
	})
	if err != nil {
		t.Errorf("error occurred setting key err=%+v", err)
	}

	getResp, err := client.Get(ctx, &CacheGetRequest{
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
	err = client.Delete(ctx, &CacheDeleteRequest{
		CacheName: cacheName,
		Key:       key,
	})
	if err != nil {
		t.Errorf("error occurred deleting key err=%+v", err)
	}
	existingCacheResp, err := client.Get(ctx, &CacheGetRequest{
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

	err = client.DeleteCache(ctx, &DeleteCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		t.Error(fmt.Errorf("error occurred deleting cache=%s err=%+v", cacheName, err))
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

func TestCreateCache(t *testing.T) {
	ctx:= context.Background()
	const correctCacheName = "correct-cache-name"
	tests := map[string]struct {
		expectedErr string
		cacheName   string
	}{
		"happy path": {
			cacheName: correctCacheName,
		},
		"test creating already existing cache name": {
			expectedErr: AlreadyExistsError,
			cacheName:   testCacheName,
		},
		"test creating empty cache name": {
			expectedErr: InvalidArgumentError,
			cacheName:   "",
		},
	}
	for name, tt := range tests {
		client, err := newTestClient()
		if err != nil {
			t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
		}
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err = client.CreateCache(ctx, &CreateCacheRequest{CacheName: tt.cacheName})
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
					"unexpected error occurred creating cache got=%+v expected=%+v",
					err, tt.expectedErr,
				)
			}

			if tt.expectedErr == "" && err != nil {
				t.Errorf("unexpected error occurred on creating cache expected=%+v got=%+v", tt.expectedErr, err)
			}

			// delete happy path cache for TestCreateCache
			if tt.cacheName == correctCacheName {
				err = client.DeleteCache(ctx, &DeleteCacheRequest{CacheName: tt.cacheName})
				if err != nil {
					t.Error(fmt.Errorf("error occurred deleting cache=%s err=%+v", tt.cacheName, err))
				}
			}
			cleanUpClient(client)
		})
	}
}

func TestDeleteCache(t *testing.T) {
	ctx := context.Background()
	var unknownCache = uuid.NewString()
	tests := map[string]struct {
		expectedErr string
		cacheName   string
	}{
		"happy path": {
			cacheName: testCacheName,
		},
		"test deleteing unknown cache name": {
			expectedErr: NotFoundError,
			cacheName:   unknownCache,
		},
		"test deleting empty cache name": {
			expectedErr: InvalidArgumentError,
			cacheName:   "",
		},
	}
	for name, tt := range tests {
		client, err := newTestClient()
		if err != nil {
			t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
		}
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := client.DeleteCache(ctx, &DeleteCacheRequest{CacheName: tt.cacheName})
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
					"unexpected error occurred deleteing cache got=%+v expected=%+v",
					err, tt.expectedErr,
				)
			}

			if tt.expectedErr == "" && err != nil {
				t.Errorf("unexpected error occurred on deleteing cache expected=%+v got=%+v", tt.expectedErr, err)
			}
			cleanUpClient(client)
		})
	}
}

func TestListCache(t *testing.T) {
	ctx := context.Background()
	var unknownCache = uuid.NewString()
	tests := map[string]struct {
		cacheName string
		inList    bool
		notInList bool
	}{
		"happy path": {
			cacheName: testCacheName,
		},
	}
	for name, tt := range tests {
		client, err := newTestClient()
		if err != nil {
			t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
		}
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			resp, err := client.ListCaches(ctx, &ListCachesRequest{})
			if err != nil {
				t.Errorf("unexpected error occurred on listing caches err=%+v", err)
			}
			var cacheNameInList = false
			var unknownCacheInList = false
			for _, cache := range resp.Caches() {
				if cache.Name() == tt.cacheName {
					cacheNameInList = true
				}
				if cache.Name() == unknownCache {
					unknownCacheInList = true
				}
			}
			if cacheNameInList == false {
				t.Errorf("cache=%s was not found in cache list", tt.cacheName)
			}
			if unknownCacheInList == true {
				t.Errorf("unexpected cache=%s was found in cache list", unknownCache)
			}
			cleanUpClient(client)
		})
	}
}

func TestCreateListRevokeSigningKeys(t *testing.T) {
	ctx := context.Background()
	client, err := newTestClient()
	if err != nil {
		t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
	}
	createSigningKeyResponse, err := client.CreateSigningKey(ctx, &CreateSigningKeyRequest{TtlMinutes: 30})
	if err != nil {
		t.Errorf("unexpected error occurred on create signing key err=%+v", err)
	}
	listSigningKeysResponse, err := client.ListSigningKeys(ctx, &ListSigningKeysRequest{})
	if err != nil {
		t.Errorf("unexpected error occurred on list signing keys err=%+v", err)
	}
	var signingKeyFound = false
	for _, signingKey := range listSigningKeysResponse.SigningKeys() {
		if signingKey.KeyId() == createSigningKeyResponse.KeyId() {
			signingKeyFound = true
			err = client.RevokeSigningKey(ctx, &RevokeSigningKeyRequest{KeyId: createSigningKeyResponse.KeyId()})
			if err != nil {
				t.Errorf("unexpected error on revoke signing key err=%+v", err)
			}
		}
	}
	if !signingKeyFound {
		t.Errorf("expected to find %s keyId in ListSigningKeysResponse, never found", createSigningKeyResponse.KeyId())
	}
	cleanUpClient(client)
}

func TestSetGet(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		key               string
		value             string
		expectedGetResult string
		ttl               uint32
	}{
		"happy path with HIT": {
			key:               uuid.NewString(),
			value:             uuid.NewString(),
			expectedGetResult: "HIT",
		},
		"test cache miss after ttl expired": {
			key:               uuid.NewString(),
			value:             uuid.NewString(),
			expectedGetResult: "MISS",
		},
		"test set with different ttl and HIT": {
			key:               uuid.NewString(),
			value:             uuid.NewString(),
			expectedGetResult: "HIT",
			ttl:               2,
		},
		"test set with different ttl and MISS": {
			key:               uuid.NewString(),
			value:             uuid.NewString(),
			expectedGetResult: "MISS",
			ttl:               2,
		},
	}
	for name, tt := range tests {
		client, err := newTestClient()
		if err != nil {
			t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
		}
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.ttl == 0 {
				// set string key/value with default ttl
				_, err := client.Set(ctx, &CacheSetRequest{CacheName: testCacheName, Key: tt.key, Value: tt.value})
				if err != nil {
					t.Errorf("unexpected error occurred on setting cache err=%+v", err)
				}
			} else {
				// set string key/value with different ttl
				_, err := client.Set(ctx, &CacheSetRequest{CacheName: testCacheName, Key: tt.key, Value: tt.value, TtlSeconds: TTL(tt.ttl)})
				if err != nil {
					t.Errorf("unexpected error occurred on setting cache err=%+v", err)
				}
			}
			if tt.expectedGetResult == "HIT" {
				resp, err := client.Get(ctx, &CacheGetRequest{CacheName: testCacheName, Key: tt.key})
				if err != nil {
					t.Errorf("unexpected error occurred on getting cache err=%+v", err)
				}
				if tt.value != resp.StringValue() {
					t.Errorf("set string value=%s is not the same as returned string value=%s", tt.value, resp.StringValue())
				}
				if tt.expectedGetResult != resp.Result() {
					t.Errorf("expected result=%s but got result=%s", tt.expectedGetResult, resp.Result())
				}
			} else {
				// make sure result it cache miss after ttl is expired
				time.Sleep(5 * time.Second)
				resp, err := client.Get(ctx, &CacheGetRequest{CacheName: testCacheName, Key: tt.key})
				if err != nil {
					t.Errorf("unexpected error occurred on getting cache err=%+v", err)
				}
				if tt.expectedGetResult != resp.Result() {
					t.Errorf("expected result=%s but got result=%s", tt.expectedGetResult, resp.Result())
				}
			}
			// set byte key/value
			_, err = client.Set(ctx, &CacheSetRequest{CacheName: testCacheName, Key: []byte(tt.key), Value: []byte(tt.value)})
			if err != nil {
				t.Errorf("unexpected error occurred on setting cache err=%+v", err)
			}
			if tt.expectedGetResult == "HIT" {
				resp, err := client.Get(ctx, &CacheGetRequest{CacheName: testCacheName, Key: []byte(tt.key)})
				if err != nil {
					t.Errorf("unexpected error occurred on getting cache err=%+v", err)
				}
				if tt.value != string(resp.ByteValue()) {
					t.Errorf("set byte value=%s is not the same as returned byte value=%s", tt.value, resp.ByteValue())
				}
				if tt.expectedGetResult != resp.Result() {
					t.Errorf("expected result=%s but got result=%s", tt.expectedGetResult, resp.Result())
				}
			}
			cleanUpClient(client)
		})
	}
}

func TestSet(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		cacheName   string
		key         interface{}
		value       interface{}
		expectedErr string
	}{
		"test set on non existant cache": {
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
			key:         nil,
			value:       uuid.NewString(),
			expectedErr: InvalidArgumentError,
		},
		"test set on nil value": {
			cacheName:   testCacheName,
			key:         uuid.NewString(),
			value:       nil,
			expectedErr: InvalidArgumentError,
		},
		"test set on bad key": {
			cacheName:   testCacheName,
			key:         1,
			value:       uuid.NewString(),
			expectedErr: InvalidArgumentError,
		},
		"test set on bad value": {
			cacheName:   testCacheName,
			key:         uuid.NewString(),
			value:       1,
			expectedErr: InvalidArgumentError,
		},
	}
	for name, tt := range tests {
		client, err := newTestClient()
		if err != nil {
			t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
		}
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := client.Set(ctx, &CacheSetRequest{CacheName: tt.cacheName, Key: tt.key, Value: tt.value})
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
					"unexpected error occurred setting cache got=%+v expected=%+v",
					err, tt.expectedErr,
				)
			}
			cleanUpClient(client)
		})
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		cacheName   string
		key         interface{}
		expectedErr string
	}{
		"test get on non existant cache": {
			cacheName:   uuid.NewString(),
			key:         uuid.NewString(),
			expectedErr: NotFoundError,
		},
		"test get on empty cache name": {
			cacheName:   "",
			key:         uuid.NewString(),
			expectedErr: InvalidArgumentError,
		},
		"test get on nil key": {
			cacheName:   testCacheName,
			key:         nil,
			expectedErr: InvalidArgumentError,
		},
		"test get on bad key": {
			cacheName:   testCacheName,
			key:         1,
			expectedErr: InvalidArgumentError,
		},
	}
	for name, tt := range tests {
		client, err := newTestClient()
		if err != nil {
			t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
		}
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := client.Get(ctx, &CacheGetRequest{CacheName: tt.cacheName, Key: tt.key})
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
					"unexpected error occurred setting cache got=%+v expected=%+v",
					err, tt.expectedErr,
				)
			}
			cleanUpClient(client)
		})
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		cacheName   string
		key         interface{}
		expectedErr string
	}{
		"test delete on non existant cache": {
			cacheName:   uuid.NewString(),
			key:         uuid.NewString(),
			expectedErr: NotFoundError,
		},
		"test delete on empty cache name": {
			cacheName:   "",
			key:         uuid.NewString(),
			expectedErr: InvalidArgumentError,
		},
		"test delete on nil key": {
			cacheName:   testCacheName,
			key:         nil,
			expectedErr: InvalidArgumentError,
		},
		"test delete on bad key": {
			cacheName:   testCacheName,
			key:         1,
			expectedErr: InvalidArgumentError,
		},
	}
	for name, tt := range tests {
		client, err := newTestClient()
		if err != nil {
			t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
		}
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err = client.Delete(ctx, &CacheDeleteRequest{CacheName: tt.cacheName, Key: tt.key})
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
					"unexpected error occurred calling cache delete got=%+v expected=%+v",
					err, tt.expectedErr,
				)
			}
			cleanUpClient(client)
		})
	}
}

func newTestClient() (ScsClient, error) {
	ctx := context.Background()
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
	err = client.CreateCache(ctx, &CreateCacheRequest{
		CacheName: testCacheName,
	})
	if momentoErr, ok := err.(MomentoError); ok {
		if momentoErr.Code() != AlreadyExistsError {
			return nil, err
		}
	}
	return client, nil
}

func cleanUpClient(client ScsClient) {
	client.Close()
}
