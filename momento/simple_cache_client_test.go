package momento

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/momentohq/client-sdk-go/internal/scserrors"
	"github.com/momentohq/client-sdk-go/momento/requests"
	"github.com/momentohq/client-sdk-go/momento/responses"

	"github.com/google/uuid"
)

var TestAuthToken = os.Getenv("TEST_AUTH_TOKEN")
var TestCacheName = os.Getenv("TEST_CACHE_NAME")

const BadToken = "eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJpbnRlZ3JhdGlvbiIsImNwIjoiY29udHJvbC5jZWxsLWFscGhhLWRldi5wcmVwcm9kLmEubW9tZW50b2hxLmNvbSIsImMiOiJjYWNoZS5jZWxsLWFscGhhLWRldi5wcmVwcm9kLmEubW9tZW50b2hxLmNvbSJ9.gdghdjjfjyehhdkkkskskmmls76573jnajhjjjhjdhnndy"
const DefaultTtlSeconds = 60

func setUp(t *testing.T) (*ScsClient, error) {
	if TestAuthToken == "" {
		t.Error("Integration tests require TEST_AUTH_TOKEN env var.")
	} else if TestCacheName == "" {
		t.Error("Integration tests require TEST_CACHE_NAME env var.")
	} else {
		client, err := SimpleCacheClient(&requests.SimpleCacheClientRequest{
			AuthToken:         TestAuthToken,
			DefaultTtlSeconds: DefaultTtlSeconds,
		})
		if err != nil {
			return nil, err
		} else {
			// Check if TestCacheName exists
			err := client.CreateCache(&requests.CreateCacheRequest{
				CacheName: TestCacheName,
			})
			if !strings.Contains(err.Error(), scserrors.AlreadyExists) {
				t.Error(err.Error())
			}
			return client, nil
		}
	}
	return nil, nil
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

	client, err := setUp(&testing.T{})
	if err != nil {
		t.Error("Set up error: " + err.Error())
	}

	err = client.CreateCache(&requests.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.Set(&requests.CacheSetRequest{
		CacheName: cacheName,
		Key:       key,
		Value:     value,
	})
	if err != nil {
		t.Error(err.Error())
	}

	getResp, err := client.Get(&requests.CacheGetRequest{
		CacheName: cacheName,
		Key:       key,
	})
	if err != nil {
		t.Error(err.Error())
	}
	if getResp.Result() != responses.HIT {
		t.Error("Cache miss")
	}
	if !bytes.Equal(getResp.ByteValue(), value) {
		t.Error("Set byte value and returned byte value are not equal")
	}

	existingCacheResp, _ := client.Get(&requests.CacheGetRequest{
		CacheName: TestCacheName,
		Key:       key,
	})
	if existingCacheResp.Result() != responses.MISS {
		t.Errorf("key: %s shouldn't exist in %s since it's never set.", string(key), TestCacheName)
	}

	deleteCacheErr := client.DeleteCache(&requests.DeleteCacheRequest{
		CacheName: cacheName,
	})
	if deleteCacheErr != nil {
		t.Error(deleteCacheErr.Error())
	}
	cleanUp(client)
}
