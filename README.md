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
	"context"
	"errors"
	"log"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"

	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}

	const (
		cacheName             = "cache"
		itemDefaultTTLSeconds = 60
	)

	// Initializes Momento
	client, err := momento.NewSimpleCacheClient(&momento.SimpleCacheClientProps{
		Configuration:      config.LatestLaptopConfig(),
		CredentialProvider: credentialProvider,
		DefaultTTLSeconds:  itemDefaultTTLSeconds,
	})
	if err != nil {
		panic(err)
	}

	// Create Cache and check if CacheName exists
	err = client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		var momentoErr momento.MomentoError
		if errors.As(err, &momentoErr) {
			if momentoErr.Code() != momento.AlreadyExistsError {
				panic(err)
			}
		}
	}

	// List caches
	token := ""
	for {
		listCacheResp, err := client.ListCaches(ctx, &momento.ListCachesRequest{NextToken: token})
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
	err = client.Set(ctx, &momento.CacheSetRequest{
		CacheName: cacheName,
		Key:       key,
		Value:     value,
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Getting key: %s\n", key)
	resp, err := client.Get(ctx, &momento.CacheGetRequest{
		CacheName: cacheName,
		Key:       key,
	})
	if err != nil {
		panic(err)
	}
	if resp.IsHit() {
		log.Printf("Lookup resulted in cahce HIT. value=%s\n", resp.AsHit().ValueString())
	} else {
		log.Printf("Look up did not find a value key=%s", key)
	}

	// Permanently delete the cache
	err = client.DeleteCache(ctx, &momento.DeleteCacheRequest{CacheName: cacheName})
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
