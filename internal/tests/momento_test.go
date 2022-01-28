package tests

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/momento/requests"
	"github.com/momentohq/client-sdk-go/momento/responses"
)

var TestAuthToken = os.Getenv("TEST_AUTH_TOKEN")
var TestCacheName = os.Getenv("TEST_CACHE_NAME")

const BadToken = "eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJpbnRlZ3JhdGlvbiIsImNwIjoiY29udHJvbC5jZWxsLWFscGhhLWRldi5wcmVwcm9kLmEubW9tZW50b2hxLmNvbSIsImMiOiJjYWNoZS5jZWxsLWFscGhhLWRldi5wcmVwcm9kLmEubW9tZW50b2hxLmNvbSJ9.gdghdjjfjyehhdkkkskskmmls76573jnajhjjjhjdhnndy"
const DefaultTtlSeconds = 60

func setUp(t *testing.T) (*momento.ScsClient, error) {
	if TestAuthToken == "" {
		t.Error("Integration tests require TEST_AUTH_TOKEN env var.")
	} else if TestCacheName == "" {
		t.Error("Integration tests require TEST_CACHE_NAME env var.")
	} else {
		simpleCacheClientRequest := requests.SimpleCacheClientRequest{
			AuthToken:         TestAuthToken,
			DefaultTtlSeconds: DefaultTtlSeconds,
		}
		client, err := momento.SimpleCacheClient(simpleCacheClientRequest)
		if err != nil {
			return nil, err
		} else {
			// Check if TestCacheName exists
			createCacheRequest := requests.CreateCacheRequest{
				CacheName: TestCacheName,
			}
			createErr := client.CreateCache(createCacheRequest)
			if !strings.Contains(createErr.Error(), "AlreadyExists") {
				t.Error(createErr.Error())
			}
			return client, nil
		}
	}
	return nil, nil
}

func cleanUp(client *momento.ScsClient) {
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
	createCacheErr := client.CreateCache(createCacheRequest)
	if createCacheErr != nil {
		t.Error(createCacheErr.Error())
	}

	setRequest := requests.CacheSetRequest{
		CacheName: cacheName,
		Key:       key,
		Value:     value,
	}
	_, setErr := client.Set(setRequest)
	if setErr != nil {
		t.Error(setErr.Error())
	}

	getRequest := requests.CacheGetRequest{
		CacheName: cacheName,
		Key:       key,
	}
	getResp, getErr := client.Get(getRequest)
	if getErr != nil {
		t.Error(getErr.Error())
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
