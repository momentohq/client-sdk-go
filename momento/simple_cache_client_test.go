package momento

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	"github.com/stretchr/testify/assert"
)

var (
	TestAuthToken = os.Getenv("TEST_AUTH_TOKEN")
	TestCacheName = os.Getenv("TEST_CACHE_NAME")
	client        *ScsClient
	err           error
)

const (
	DefaultTtlSeconds = 60
)

func TestMain(m *testing.M) {
	client, err := setUp()
	if err != nil {
		panic(err)
	}
	exitVal := m.Run()
	cleanUp(client)

	os.Exit(exitVal)
}

func setUp() (*ScsClient, error) {
	if TestAuthToken == "" {
		return nil, fmt.Errorf("Integration tests require TEST_AUTH_TOKEN env var.")
	} else if TestCacheName == "" {
		return nil, fmt.Errorf("Integration tests require TEST_CACHE_NAME env var.")
	} else {
		client, err = SimpleCacheClient(&SimpleCacheClientRequest{
			AuthToken:         TestAuthToken,
			DefaultTtlSeconds: DefaultTtlSeconds,
		})
		if err != nil {
			return nil, err
		} else {
			// Check if TestCacheName exists
			err := client.CreateCache(&CreateCacheRequest{
				CacheName: TestCacheName,
			})
			if err != nil && err.Code() != "AlreadyExists" {
				return nil, err
			}
			return client, nil
		}
	}
}

func cleanUp(client *ScsClient) {
	err := client.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Basic happy path test - create a cache, operate set/get, and delete the cache
func TestCreateCacheGetSetValueAndDeleteCache(t *testing.T) {
	cacheName := uuid.NewString()
	key := []byte(uuid.NewString())
	value := []byte(uuid.NewString())

	err = client.CreateCache(&CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		t.Error(err.Error())
	}

	_, err := client.Set(&CacheSetRequest{
		CacheName: cacheName,
		Key:       key,
		Value:     value,
	})
	if err != nil {
		t.Error(err.Error())
	}

	getResp, err := client.Get(&CacheGetRequest{
		CacheName: cacheName,
		Key:       key,
	})
	if err != nil {
		t.Error(err.Error())
	}
	if getResp.Result() != HIT {
		t.Error("Cache miss")
	}
	if !bytes.Equal(getResp.ByteValue(), value) {
		t.Error("Set byte value and returned byte value are not equal")
	}

	existingCacheResp, err := client.Get(&CacheGetRequest{
		CacheName: TestCacheName,
		Key:       key,
	})
	if err != nil {
		t.Error(err.Error())
	}
	if existingCacheResp.Result() != MISS {
		t.Errorf("key: %s shouldn't exist in %s since it's never set.", string(key), TestCacheName)
	}

	err = client.DeleteCache(&DeleteCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		t.Error(err.Error())
	}
}

func TestZeroRequestTimeout(t *testing.T) {
	timeout := uint32(0)
	_, err := SimpleCacheClient(&SimpleCacheClientRequest{
		AuthToken:             TestAuthToken,
		DefaultTtlSeconds:     DefaultTtlSeconds,
		RequestTimeoutSeconds: &timeout,
	})
	if assert.Error(t, err) {
		assert.Equal(t, momentoerrors.InvalidArgumentError, err.Code())
	}
}
