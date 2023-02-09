<head>
  <meta name="Momento Go Client Library Documentation" content="Go client software development kit for Momento Serverless Cache">
</head>
<img src="https://docs.momentohq.com/img/logo.svg" alt="logo" width="400"/>

[![project status](https://momentohq.github.io/standards-and-practices/badges/project-status-official.svg)](https://github.com/momentohq/standards-and-practices/blob/main/docs/momento-on-github.md)
[![project stability](https://momentohq.github.io/standards-and-practices/badges/project-stability-alpha.svg)](https://github.com/momentohq/standards-and-practices/blob/main/docs/momento-on-github.md) 

# Momento Go Client Library


Go client SDK for Momento Serverless Cache: a fast, simple, pay-as-you-go caching solution without
any of the operational overhead required by traditional caching solutions!



## Getting Started :running:

### Requirements

- [Go version 1.18.\*](https://go.dev/dl/)
- A Momento Auth Token is required, you can generate one using
  the [Momento CLI](https://github.com/momentohq/momento-cli)

### Examples

Check out full working code in the [examples](./examples/README.md) directory of this repository!

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
		cacheName             = "my-test-cache"
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

	// Sets key with default TTL and gets value with that key
	key := uuid.NewString()
	value := uuid.NewString()
	log.Printf("Setting key: %s, value: %s\n", key, value)
	err = client.Set(ctx, &momento.CacheSetRequest{
		CacheName: cacheName,
		Key:       &momento.StringBytes{Text: key},
		Value:     &momento.StringBytes{Text: value},
	})
	if err != nil {
		panic(err)
	}

	log.Printf("Getting key: %s\n", key)
	resp, err := client.Get(ctx, &momento.CacheGetRequest{
		CacheName: cacheName,
		Key:       &momento.StringBytes{Text: key},
	})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *momento.CacheGetHit:
		log.Printf("Lookup resulted in cahce HIT. value=%s\n", r.ValueString())
	case *momento.CacheGetMiss:
		log.Printf("Look up did not find a value key=%s", key)
	}

	// Permanently delete the cache
	if err = client.DeleteCache(ctx, &momento.DeleteCacheRequest{CacheName: cacheName}); err != nil {
		panic(err)
	}
	log.Printf("Cache named %s is deleted\n", cacheName)
}

```

### Error Handling

The preferred way of interpreting the return values from `ScsClient` methods is using a `switch` statement to match and handle the specific response type. 
Here's a quick example:

```go
switch r := resp.(type) {
case *momento.CacheGetHit:
    log.Printf("Lookup resulted in cahce HIT. value=%s\n", r.ValueString())
default: 
    // you can handle other cases via pattern matching in other `switch case`, or a default case
    // via the `default` block.  For each return value your IDE should be able to give you code 
    // completion indicating the other possible "case"; in this case, `*momento.CacheGetMiss`.
}
```

Using this approach, you get a type-safe `CacheGetHit` object in the case of a cache hit. 
But if the cache read results in a Miss, you'll also get a type-safe object that you can use to get more info about what happened.

In cases where you get an error response, it can be treated as `momentoErr` using `As` method and it always include an `momentoErr.Code` that you can use to check the error type:

```go
_, err := client.Get(ctx, &momento.CacheGetRequest{
    CacheName: cacheName,
    Key:       &momento.StringBytes{Text: key},
})

if err != nil {
    var momentoErr momento.MomentoError
    if errors.As(err, &momentoErr) {
        if momentoErr.Code() != momento.TimeoutError {
            // this would represent a client-side timeout, and you could fall back to your original data source
        }
    }
}
```

### Tuning

Coming soon...

----------------------------------------------------------------------------------------
For more info, visit our website at [https://gomomento.com](https://gomomento.com)!
