package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

var (
	ctx    context.Context
	client momento.CacheClient
)

func retrieveAuthTokenFromSecretsManager() string {
	fakeTestV1ApiKey := ""

	return fakeTestV1ApiKey
}

func example_API_CredentialProviderFromEnvVar() {
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}
}

func example_API_CreateCache() {
	_, err := client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: "cache-name",
	})
	if err != nil {
		panic(err)
	}
}

func example_API_ListCaches() {
	_, err := client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: "cache-name",
	})
	if err != nil {
		panic(err)
	}
}

func example_API_Get() {
	key := uuid.NewString()
	resp, err := client.Get(ctx, &momento.GetRequest{
		CacheName: "cache-name",
		Key:       momento.String(key),
	})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *responses.GetHit:
		log.Printf("Lookup resulted in cache HIT. value=%s\n", r.ValueString())
	case *responses.GetMiss:
		log.Printf("Look up did not find a value key=%s", key)
	}
}

func example_API_Set() {
	key := uuid.NewString()
	value := uuid.NewString()
	log.Printf("Setting key: %s, value: %s\n", key, value)
	_, err := client.Set(ctx, &momento.SetRequest{
		CacheName: "cache-name",
		Key:       momento.String(key),
		Value: 	   momento.String(value),
		Ttl: 	   time.Duration(9999),
	})
	if err != nil {
		var momentoErr momento.MomentoError
    	if errors.As(err, &momentoErr) {
        if momentoErr.Code() != momento.TimeoutError {
            // this would represent a client-side timeout, and you could fall back to your original data source
        } else {
			panic(err)
		}
    }
	}
}

func example_API_Delete() {
	key := uuid.NewString()
	_, err := client.Delete(ctx, &momento.DeleteRequest{
		CacheName: "cache-name",
		Key:       momento.String(key),
	})
	if err != nil {
		panic(err)
	}
}
