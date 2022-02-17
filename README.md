# Momento client-sdk-go

:warning: Experimental SDK :warning:

Go SDK for Momento is experimental and under active development.
There could be non-backward compatible changes or removal in the future.
Please be aware that you may need to update your source code with the current version of the SDK when its version gets upgraded.

---

<br />

<div align="center">
    <img src="images/gopher.png" alt="Logo" width="200" height="150">
</div>

Go SDK for Momento, a serverless cache that automatically scales without any of the operational overhead required by traditional caching solutions.

<br/>

# Getting Started :running:

## Requirements

- [Go version 1.17.\*](https://go.dev/dl/)
- A Momento Auth Token is required, you can generate one using the [Momento CLI](https://github.com/momentohq/momento-cli)
- golang
  - `brew install go`
- golint
  - `go get -u golang.org/x/lint/golint`


<br/>

## Installing Momento and Running the Example

Check out our [Go SDK example repo](https://github.com/momentohq/client-sdk-examples/tree/main/golang)!

<br />

## Using Momento

```go
import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/momento"
)

func main() {
	authToken := os.Getenv("MOMENTO_AUTH_TOKEN")
	const (
		cacheName             = "cache"
		itemDefaultTtlSeconds = 60
	)

	if authToken == "" {
		log.Fatal("Missing required environment variable MOMENTO_AUTH_TOKEN")
	}

	// Initializes Momento
	client, err := momento.NewSimpleCacheClient(authToken, itemDefaultTtlSeconds)
	if err != nil {
		panic(err)
	}

	// Create Cache and check if CacheName exists
	err = client.CreateCache(&momento.CreateCacheRequest{
		CacheName: CacheName,
	})
	if momentoErr, ok := err.(momento.MomentoError); ok {
        if momentoErr.Code() != AlreadyExists {
            return nil, err
        }
	}
	log.Printf("Cache named %s is created\n", CacheName)

	// Sets key with default TTL and gets value with that key
	key := []byte(uuid.NewString())
	value := []byte(uuid.NewString())
	log.Printf("Setting key: %s, value: %s\n", key, value)
	_, err = client.Set(&momento.CacheSetRequest{
		CacheName: CacheName,
		Key:       key,
		Value:     value,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("Getting key: %s\n", key)
	resp, err := client.Get(&momento.CacheGetRequest{
		CacheName: CacheName,
		Key:       key,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("Lookup resulted in a : %s\n", resp.Result())
	log.Printf("Looked up value: %s\n", resp.StringValue())

	// Permanently delete the cache
	err = client.DeleteCache(&momento.DeleteCacheRequest{CacheName: CacheName})
	if err != nil {
		panic(err)
	}
	log.Printf("Cache named %s is deleted\n", CacheName)
}
```

<br />

# Running Tests :zap:

## Requirements

- `TEST_AUTH_TOKEN` - an auth token for testing
- `TEST_CACHE_NAME` - any string value would work

## How to Run Test

```bash
TEST_AUTH_TOKEN=<auth token> TEST_CACHE_NAME=<cache name> go test -v ./momento
```
