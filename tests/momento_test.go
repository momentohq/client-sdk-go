package tests

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/cacheclient"
	"github.com/momentohq/client-sdk-go/responses"
)

var TestAuthToken = os.Getenv("TEST_AUTH_TOKEN")
var TestCacheName = os.Getenv("TEST_CACHE_NAME")

const BadToken = "eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJpbnRlZ3JhdGlvbiIsImNwIjoiY29udHJvbC5jZWxsLWFscGhhLWRldi5wcmVwcm9kLmEubW9tZW50b2hxLmNvbSIsImMiOiJjYWNoZS5jZWxsLWFscGhhLWRldi5wcmVwcm9kLmEubW9tZW50b2hxLmNvbSJ9.gdghdjjfjyehhdkkkskskmmls76573jnajhjjjhjdhnndy"
const DefaultTtlSeconds = 60

func setUp(t *testing.T) (*cacheclient.ScsClient, error) {
	if TestAuthToken == "" {
		log.Fatal("Integration tests require TEST_AUTH_TOKEN env var.")
	} else if TestCacheName == "" {
		log.Fatal("Integration tests require TEST_CACHE_NAME env var.")
	} else {
		client, err := cacheclient.SimpleCacheClient(TestAuthToken, DefaultTtlSeconds)
		if err != nil {
			return nil, err
		} else {
			// Check if TestCacheName exists
			createErr := client.CreateCache(TestCacheName)
			if !strings.Contains(createErr.Error(), "AlreadyExists") {
				log.Fatal(createErr.Error())
			}
			return client, nil
		}
	}
	return nil, nil
}

func cleanUp(client *cacheclient.ScsClient) {
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

	createCacheErr := client.CreateCache(cacheName)
	if createCacheErr != nil {
		t.Error(createCacheErr.Error())
	}

	_, setErr := client.Set(cacheName, key, value, DefaultTtlSeconds)
	if setErr != nil {
		t.Error(setErr.Error())
	}

	getResp, getErr := client.Get(cacheName, key)
	if getErr != nil {
		t.Error(getErr.Error())
	}
	if getResp.Result() != responses.HIT {
		t.Error("Cache miss")
	}
	if !bytes.Equal(getResp.ByteValue(), value) {
		t.Error("Set byte value and returned byte value are not equal")
	}

	existingCacheResp, _ := client.Get(TestCacheName, key)
	if existingCacheResp.Result() != responses.MISS {
		t.Errorf("key: %s shouldn't exist in %s since it's never set.", string(key), TestCacheName)
	}

	deleteCacheErr := client.DeleteCache(cacheName)
	if deleteCacheErr != nil {
		t.Error(deleteCacheErr.Error())
	}
	cleanUp(client)
}
