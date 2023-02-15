package momento

import (
	"bytes"
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

// Basic happy path test - create a cache, operate set/get, and delete the cache
func TestBasicHappyPathSDKFlow(t *testing.T) {
	ctx := context.Background()
	randomCacheName := uuid.NewString()
	key := []byte(uuid.NewString())
	value := []byte(uuid.NewString())
	client, err := newTestClient(testCredentialProvider)
	defer teardown(client, cacheName, randomCacheName)

	if err != nil {
		t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
	}
	err = client.CreateCache(ctx, &CreateCacheRequest{
		CacheName: randomCacheName,
	})
	if err != nil {
		t.Error(fmt.Errorf("error occurred creating cache err=%+v", err))
	}

	_, err = client.Set(ctx, &SetRequest{
		CacheName: randomCacheName,
		Key:       RawBytes{Bytes: key},
		Value:     RawBytes{Bytes: value},
	})
	if err != nil {
		t.Errorf("error occurred setting key err=%+v", err)
	}

	_, err = client.Set(ctx, &SetRequest{
		CacheName: randomCacheName,
		Key:       StringBytes{Text: uuid.NewString()},
		Value:     RawBytes{Bytes: value},
		TTL:       time.Second * 60,
	})
	if err != nil {
		t.Errorf("error occurred setting key with custom ttl err=%+v", err)
	}

	getResp, err := client.Get(ctx, &GetRequest{
		CacheName: randomCacheName,
		Key:       RawBytes{Bytes: key},
	})
	if err != nil {
		t.Errorf("error occurred getting key err=%+v", err)
		return
	}

	switch result := getResp.(type) {
	case GetHit:
		if !bytes.Equal(result.ValueByte(), value) {
			t.Errorf(
				"set byte value and returned byte value are not equal "+
					"got=%+v expected=%+v", result.ValueByte(), value,
			)
		}
	default:
		t.Errorf("unexpected responseType when getting test key got=%+v expected=%+v", getResp, GetHit{})
	}

	existingCacheResp, err := client.Get(ctx, &GetRequest{
		CacheName: testCacheName,
		Key:       RawBytes{Bytes: key},
	})
	if err != nil {
		t.Error(err.Error())
	}

	if r, ok := existingCacheResp.(GetHit); ok {
		t.Errorf(
			"key: %s shouldn't exist in %s since it's got deleted. got=%s",
			string(key), testCacheName, r.ValueString(),
		)
	}

	err = client.DeleteCache(ctx, &DeleteCacheRequest{
		CacheName: randomCacheName,
	})
	if err != nil {
		t.Error(fmt.Errorf("error occurred deleting cache=%s err=%+v", randomCacheName, err))
	}
}

func TestBasicHappyPathDelete(t *testing.T) {
	ctx := context.Background()
	cacheName := uuid.NewString()
	key := []byte(uuid.NewString())
	value := []byte(uuid.NewString())
	client, err := newTestClient(testCredentialProvider)
	defer teardown(client, testCacheName, cacheName)

	if err != nil {
		t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
	}
	err = client.CreateCache(ctx, &CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		t.Error(fmt.Errorf("error occurred creating cache err=%+v", err))
	}

	_, err = client.Set(ctx, &SetRequest{
		CacheName: cacheName,
		Key:       RawBytes{Bytes: key},
		Value:     RawBytes{Bytes: value},
	})
	if err != nil {
		t.Errorf("error occurred setting key err=%+v", err)
	}

	getResp, err := client.Get(ctx, &GetRequest{
		CacheName: cacheName,
		Key:       RawBytes{Bytes: key},
	})
	if err != nil {
		t.Errorf("error occurred getting key err=%+v", err)
		return
	}

	switch result := getResp.(type) {
	case GetHit:
		if !bytes.Equal(result.ValueByte(), value) {
			t.Errorf(
				"set byte value and returned byte value are not equal "+
					"got=%+v expected=%+v", result.ValueByte(), value,
			)
		}
	default:
		t.Errorf("unexpected responseType when getting test key got=%+v expected=%+v", getResp, GetHit{})
	}

	_, err = client.Delete(ctx, &DeleteRequest{
		CacheName: cacheName,
		Key:       RawBytes{Bytes: key},
	})
	if err != nil {
		t.Errorf("error occurred deleting key err=%+v", err)
	}
	existingCacheResp, err := client.Get(ctx, &GetRequest{
		CacheName: testCacheName,
		Key:       RawBytes{Bytes: key},
	})
	if err != nil {
		t.Error(err.Error())
	}

	if r, ok := existingCacheResp.(GetHit); ok {
		t.Errorf(
			"key: %s shouldn't exist in %s since it's got deleted. got=%s",
			string(key), testCacheName, r.ValueString(),
		)
	}

	err = client.DeleteCache(ctx, &DeleteCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		t.Error(fmt.Errorf("error occurred deleting cache=%s err=%+v", cacheName, err))
	}
}

func TestCredentialProvider(t *testing.T) {
	err := os.Setenv("BAD_TEST_AUTH_TOKEN", "Iamnotanauthtoken")
	if err != nil {
		return
	}
	_, err = newCredentialProvider("BAD_TEST_AUTH_TOKEN")
	if err == nil {
		t.Fatal("missing expected error for bad auth token")
	}
	var momentoErr MomentoError
	if errors.As(err, &momentoErr) {
		if momentoErr.Code() != ClientSdkError {
			t.Error("missing expected ClientSdkError")
		}
	}
}

func TestClientInitialization(t *testing.T) {
	testRequestTimeout := 100 * time.Second
	badRequestTimeout := 0 * time.Second
	tests := map[string]struct {
		expectedErr    string
		defaultTTL     time.Duration
		requestTimeout *time.Duration
	}{
		"happy path": {
			defaultTTL: defaultTTL,
		},
		"happy path custom timeout": {
			defaultTTL:     defaultTTL,
			requestTimeout: &testRequestTimeout,
		},
		"test invalid request timeout": {
			expectedErr:    InvalidArgumentError,
			defaultTTL:     defaultTTL,
			requestTimeout: &badRequestTimeout,
		},
	}
	for name, tt := range tests {
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			clientProps := &SimpleCacheClientProps{
				CredentialProvider: testCredentialProvider,
				DefaultTTL:         tt.defaultTTL,
			}
			if tt.requestTimeout != nil {
				clientProps.Configuration = config.LatestLaptopConfig().WithClientTimeout(*tt.requestTimeout)
			} else {
				clientProps.Configuration = config.LatestLaptopConfig()
			}
			c, err := NewSimpleCacheClient(clientProps)
			defer teardown(c)

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
					"unexpected error occurred initializing client got=%+v expected=%+v",
					err, tt.expectedErr,
				)
			}

			if tt.expectedErr == "" && err != nil {
				t.Errorf("unexpected error occurred on init expected=%+v got=%+v", tt.expectedErr, err)
			}
		})
	}
}

func TestCreateCache(t *testing.T) {
	ctx := context.Background()
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
		client, err := newTestClient(testCredentialProvider)
		defer teardown(client)

		if err != nil {
			t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
		}
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			err = client.CreateCache(ctx, &CreateCacheRequest{CacheName: tt.cacheName})
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
		"test deleting unknown cache name": {
			expectedErr: NotFoundError,
			cacheName:   unknownCache,
		},
		"test deleting empty cache name": {
			expectedErr: InvalidArgumentError,
			cacheName:   "",
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
			err := client.DeleteCache(ctx, &DeleteCacheRequest{CacheName: tt.cacheName})
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
					"unexpected error occurred deleteing cache got=%+v expected=%+v",
					err, tt.expectedErr,
				)
			}

			if tt.expectedErr == "" && err != nil {
				t.Errorf("unexpected error occurred on deleteing cache expected=%+v got=%+v", tt.expectedErr, err)
			}
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
		client, err := newTestClient(testCredentialProvider)
		defer teardown(client)

		if err != nil {
			t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
		}
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
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
		})
	}
}

func TestSetGet(t *testing.T) {
	ctx := context.Background()
	tests := map[string]struct {
		key               string
		value             string
		expectedGetResult string
		ttl               time.Duration
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
			ttl:               time.Duration(time.Second * 2),
		},
		"test set with different ttl and MISS": {
			key:               uuid.NewString(),
			value:             uuid.NewString(),
			expectedGetResult: "MISS",
			ttl:               time.Duration(time.Second * 2),
		},
	}
	for name, tt := range tests {
		client, err := newTestClient(testCredentialProvider)
		defer teardown(client, testCacheName)

		if err != nil {
			t.Error(fmt.Errorf("error occurred setting up client err=%+v", err))
		}
		tt := tt // for t.Parallel()
		t.Run(name, func(t *testing.T) {
			if tt.ttl == 0 {
				// set string key/value with default ttl
				_, err := client.Set(ctx, &SetRequest{
					CacheName: testCacheName,
					Key:       &StringBytes{Text: tt.key},
					Value:     &StringBytes{Text: tt.value},
				})
				if err != nil {
					t.Errorf("unexpected error occurred on setting cache err=%+v", err)
				}
			} else {
				// set string key/value with different ttl
				_, err := client.Set(ctx, &SetRequest{
					CacheName: testCacheName,
					Key:       &StringBytes{Text: tt.key},
					Value:     &StringBytes{Text: tt.value},
					TTL:       tt.ttl,
				})
				if err != nil {
					t.Errorf("unexpected error occurred on setting cache err=%+v", err)
				}
			}

			if tt.expectedGetResult == "HIT" {
				resp, err := client.Get(ctx, &GetRequest{
					CacheName: testCacheName,
					Key:       &StringBytes{Text: tt.key},
				})
				if err != nil {
					t.Errorf("unexpected error occurred on getting cache err=%+v", err)
				}
				switch result := resp.(type) {
				case GetHit:
					if tt.value != result.ValueString() {
						t.Errorf(
							"set string value=%s is not the same as returned string value=%s",
							tt.value, result.ValueString(),
						)
					}
				default:
					t.Errorf("expected hit but got responseType=%+v", resp)
				}

			} else {
				// make sure responseType it cache miss after ttl is expired
				time.Sleep(5 * time.Second)
				resp, err := client.Get(ctx, &GetRequest{
					CacheName: testCacheName,
					Key:       &StringBytes{Text: tt.key},
				})
				if err != nil {
					t.Errorf("unexpected error occurred on getting cache err=%+v", err)
				}
				switch result := resp.(type) {
				case GetMiss:
					// We expect miss
				default:
					t.Errorf("expected miss but got responseType=%+v", result)
				}

			}
		})
	}
}

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
				Key:       &StringBytes{Text: tt.key},
				Value:     &StringBytes{Text: tt.value},
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
				Key:       &StringBytes{Text: tt.key},
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
				Key:       StringBytes{Text: tt.key},
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

func newTestClient(credentialProvider auth.CredentialProvider) (*ScsClient, error) {
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
	err = client.CreateCache(ctx, &CreateCacheRequest{
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

func teardown(client *ScsClient, cacheNames ...string) {
	ctx := context.Background()

	for _, cacheName := range cacheNames {
		err := client.DeleteCache(ctx, &DeleteCacheRequest{
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
