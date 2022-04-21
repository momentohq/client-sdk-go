# Momento client-sdk-go

:warning: Experimental SDK :warning:

Go SDK for Momento is experimental and under active development. There could be non-backward compatible changes or
removal in the future. Please be aware that you may need to update your source code with the current version of the SDK
when its version gets upgraded.

---

<br />

<div align="center">
    <img src="images/gopher.png" alt="Logo" width="200" height="150">
</div>

Go SDK for Momento, a serverless cache that automatically scales without any of the operational overhead required by
traditional caching solutions.

<br/>

# Getting Started :running:

## Requirements

- [Go version 1.17.\*](https://go.dev/dl/)
- A Momento Auth Token is required, you can generate one using
  the [Momento CLI](https://github.com/momentohq/momento-cli)
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
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	. "github.com/momentohq/client-sdk-go/momento"
)

func main() {
	authToken := os.Getenv("MOMENTO_AUTH_TOKEN")
	const (
		cacheName             = "my-first-cache-😊"
		itemDefaultTtlSeconds = 60
	)

	if authToken == "" {
		log.Fatal("Missing required environment variable MOMENTO_AUTH_TOKEN")
	}

	// Initializes Momento
	client, err := NewSimpleCacheClient(authToken, itemDefaultTtlSeconds)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed initilizing cache client with err %+v", err))
	}

	// Create Cache Ignore if Cache already exists
	err = client.CreateCache(&CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		if momentoErr, ok := err.(MomentoError); ok {
			if momentoErr.Code() != AlreadyExists {
				log.Fatal(fmt.Sprintf(
					"failed creating cache %s with err %+v",
					cacheName, momentoErr,
				))
			}
		}
	}
	log.Println(fmt.Sprintf("Cache %s is created", cacheName))

	// Sets key with default TTL and custom ttl and then retrieves the items from cache
	key1 := []byte(uuid.NewString())
	key2 := []byte(uuid.NewString())
	value1 := []byte(uuid.NewString())
	value2 := []byte(uuid.NewString())

	log.Println(fmt.Sprintf("Setting key: %s, value: %s", key1, value1))
	_, err = client.Set(&CacheSetRequest{
		CacheName: cacheName,
		Key:       key1,
		Value:     value2,
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("failed setting key %s with err %+v", key1, err))
	}
	log.Println(fmt.Sprintf("Setting key with custom ttl key: %s, value: %s", key2, value2))
	_, err = client.Set(&CacheSetRequest{
		CacheName:  cacheName,
		Key:        key2,
		Value:      value2,
		TtlSeconds: TTL(60),
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("failed setting key %s with err %+v", key2, err))
	}
	log.Println(fmt.Sprintf("Getting key: %s", key1))
	resp, err := client.Get(&CacheGetRequest{
		CacheName: cacheName,
		Key:       key1,
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("failed getting key %s with err %+v", key1, err))
	}
	log.Println(fmt.Sprintf(
		"Get request succeded result: %s value: %s",
		resp.Result(), resp.StringValue(),
	))

	// Delete the cache
	err = client.DeleteCache(&DeleteCacheRequest{CacheName: cacheName})
	if err != nil {
		log.Fatal(fmt.Sprintf("failed deleting cache %s with err %+v", cacheName, err))
	}
	log.Printf(fmt.Sprintf("Cache %s is deleted", cacheName))
}
```

<br />

You can also specify request timeout for Momento client

```golang
var authToken = os.Getenv("MOMENTO_AUTH_TOKEN")
const (
		cacheName             = "cache"
		itemDefaultTtlSeconds = 60
		requestTimeoutSeconds = 10
	)
client, err = NewSimpleCacheClient(authToken, itemDefaultTtlSeconds, WithRequestTimeout(requestTimeoutSeconds))
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

## Updating GRPC protos
1. Follow the [quick-start instructions](https://grpc.io/docs/languages/go/quickstart/) to set up your environment
2. Checkout the latest changes from https://github.com/momentohq/client_protos
3. Copy the `.proto` files from `client_protos/proto` to `client-sdk-go/internal/protos/`
4. `cd` to the `internal` directory and run:

```protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative protos/*.proto
```

You should now have updated auto-generated files with the updated protos