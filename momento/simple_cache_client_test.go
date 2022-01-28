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
		simpleCacheClientRequest := requests.SimpleCacheClientRequest{
			AuthToken:         TestAuthToken,
			DefaultTtlSeconds: DefaultTtlSeconds,
		}
		client, err := SimpleCacheClient(simpleCacheClientRequest)
		if err != nil {
			return nil, err
		} else {
			// Check if TestCacheName exists
			createCacheRequest := requests.CreateCacheRequest{
				CacheName: TestCacheName,
			}
			err := client.CreateCache(createCacheRequest)
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

	createCacheRequest := requests.CreateCacheRequest{
		CacheName: cacheName,
	}
	err = client.CreateCache(createCacheRequest)
	if err != nil {
		t.Error(err.Error())
	}

	setRequest := requests.CacheSetRequest{
		CacheName: cacheName,
		Key:       key,
		Value:     value,
	}
	_, err = client.Set(setRequest)
	if err != nil {
		t.Error(err.Error())
	}

	getRequest := requests.CacheGetRequest{
		CacheName: cacheName,
		Key:       key,
	}
	getResp, err := client.Get(getRequest)
	if err != nil {
		t.Error(err.Error())
	}
	if getResp.Result() != responses.HIT {
		t.Error("Cache miss")
	}
	if !bytes.Equal(getResp.ByteValue(), value) {
		t.Error("Set byte value and returned byte value are not equal")
	}

	getRequest2 := requests.CacheGetRequest{
		CacheName: TestCacheName,
		Key:       key,
	}
	existingCacheResp, _ := client.Get(getRequest2)
	if existingCacheResp.Result() != responses.MISS {
		t.Errorf("key: %s shouldn't exist in %s since it's never set.", string(key), TestCacheName)
	}

	deleteCacheRequest := requests.DeleteCacheRequest{
		CacheName: cacheName,
	}
	deleteCacheErr := client.DeleteCache(deleteCacheRequest)
	if deleteCacheErr != nil {
		t.Error(deleteCacheErr.Error())
	}
	cleanUp(client)
}
