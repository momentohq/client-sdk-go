package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

func main() {
	context := context.Background()
	configuration := config.LaptopLatestWithLogger(logger.NewNoopMomentoLoggerFactory()).WithClientTimeout(15 * time.Second)
	credentialProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}
	defaultTtl := time.Duration(9999)

	client, err := momento.NewCacheClient(configuration, credentialProvider, defaultTtl)
	if err != nil {
		panic(err)
	}

	key := uuid.NewString()
	value := uuid.NewString()
	log.Printf("Setting key: %s, value: %s\n", key, value)
	_, _ = client.Set(context, &momento.SetRequest{
		CacheName: "cache-name",
		Key:       momento.String(key),
		Value:     momento.String(value),
		Ttl:       time.Duration(9999),
	})

	if err != nil {
		panic(err)
	}

	resp, err := client.Get(context, &momento.GetRequest{
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
