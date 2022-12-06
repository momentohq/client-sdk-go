<head>
  <meta name="Momento Go Client Library Documentation" content="Go client software development kit for Momento Serverless Cache">
</head>
<img src="https://docs.momentohq.com/img/logo.svg" alt="logo" width="400"/>

[![project status](https://momentohq.github.io/standards-and-practices/badges/project-status-official.svg)](https://github.com/momentohq/standards-and-practices/blob/main/docs/momento-on-github.md)
[![project stability](https://momentohq.github.io/standards-and-practices/badges/project-stability-experimental.svg)](https://github.com/momentohq/standards-and-practices/blob/main/docs/momento-on-github.md) 

# Momento Go Client Library


:warning: Experimental SDK :warning:

This is an official Momento SDK, but the API is in an early experimental stage and subject to backward-incompatible
changes.  For more info, click on the experimental badge above.


Go client SDK for Momento Serverless Cache: a fast, simple, pay-as-you-go caching solution without
any of the operational overhead required by traditional caching solutions!



## Getting Started :running:

### Requirements

- [Go version 1.18.\*](https://go.dev/dl/)
- A Momento Auth Token is required, you can generate one using
  the [Momento CLI](https://github.com/momentohq/momento-cli)
- golang
  - `brew install go`
- golint
  - `go get -u golang.org/x/lint/golint`

### Examples

Ready to dive right in? Just check out the [examples](./examples/README.md) directory for complete, working examples of
how to use the SDK.

### Installation

```bash
go get github.com/momentohq/client-sdk-go
```

### Usage

Checkout our [examples](./examples/README.md) directory for complete examples of how to use the SDK.

Here is a quickstart you can use in your own project:

```go
package main

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/momento"
)

func main() {
	var authToken = os.Getenv("MOMENTO_AUTH_TOKEN")
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
		CacheName: cacheName,
	})
	if err != nil && err.Error() != momento.AlreadyExistsError {
		panic(err)
	}
	log.Printf("Cache named %s is created\n", cacheName)

	// List caches
	token := ""
	for {
		listCacheResp, err := client.ListCaches(&momento.ListCachesRequest{NextToken: token})
		if err != nil {
			panic(err)
		}
		for _, cacheInfo := range listCacheResp.Caches() {
			log.Printf("%s\n", cacheInfo.Name())
		}
		token = listCacheResp.NextToken()
		if token == "" {
			break
		}
	}

	// Sets key with default TTL and gets value with that key
	key := []byte(uuid.NewString())
	value := []byte(uuid.NewString())
	log.Printf("Setting key: %s, value: %s\n", key, value)
	_, err = client.Set(&momento.CacheSetRequest{
		CacheName: cacheName,
		Key:       key,
		Value:     value,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("Getting key: %s\n", key)
	resp, err := client.Get(&momento.CacheGetRequest{
		CacheName: cacheName,
		Key:       key,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("Lookup resulted in a : %s\n", resp.Result())
	log.Printf("Looked up value: %s\n", resp.StringValue())

	// Permanently delete the cache
	err = client.DeleteCache(&momento.DeleteCacheRequest{CacheName: cacheName})
	if err != nil {
		panic(err)
	}
	log.Printf("Cache named %s is deleted\n", cacheName)
}

```

### Error Handling

Coming soon...

### Tuning

Coming soon...

----------------------------------------------------------------------------------------
For more info, visit our website at [https://gomomento.com](https://gomomento.com)!
