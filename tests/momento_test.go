package tests

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/cacheclient"
)

const DefaultTtlSeconds = 60

func setUp(t *testing.T) (*cacheclient.ScsClient, error) {
	authToken := os.Getenv("TEST_AUTH_TOKEN")
	if authToken == "" {
		log.Fatal("Integration tests require TEST_AUTH_TOKEN env var.")
	} else {
		client, err := cacheclient.SimpleCacheClient(authToken, DefaultTtlSeconds)
		if err != nil {
			return nil, err
		} else {
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

func TestCreateCacheGetSetValueAndDeleteCache(t *testing.T) {
	cacheName := uuid.NewString()
	key := []byte(uuid.NewString())
	value := []byte(uuid.NewString())
	client, err := setUp(&testing.T{})
	if err != nil {
		log.Fatal("setUp error: " + err.Error())
	}
	createCacheErr := client.CreateCache(cacheName)
	if createCacheErr != nil {
		log.Fatal(createCacheErr.Error())
	}
	_, setErr := client.Set(cacheName, key, value, DefaultTtlSeconds)
	if setErr != nil {
		log.Fatal(setErr.Error())
	}
	getResp, getErr := client.Get(cacheName, key)
	if getErr != nil {
		log.Fatal(getErr.Error())
	}
	if !bytes.Equal(getResp.ByteValue(), value) {
		log.Fatal("Set byte value and returned byte value are not equal")
	}
	deleteCacheErr := client.DeleteCache(cacheName)
	if deleteCacheErr != nil {
		log.Fatal(deleteCacheErr.Error())
	}
	cleanUp(client)
}
