package main

import "C"

import (
	"log"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/protosocket"

	"github.com/google/uuid"
)

const (
	cacheName             = "protosocket-loadgen"
	itemDefaultTTLSeconds = 60
)

func initProtosocket() {
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}

	// Initializes Momento protosocket client
	err = protosocket.NewProtosocketCacheClient(
		config.LaptopLatest(),
		credentialProvider,
		itemDefaultTTLSeconds*time.Second,
	)
	if err != nil {
		panic(err)
	}
}

func main() {
	initProtosocket()

	key := uuid.NewString()
	value := uuid.NewString()
	log.Printf("Setting key: %s, value: %s\n", key, value)
	protosocket.ProtosocketSet(cacheName, key, value)

	log.Printf("Getting key: %s\n", key)
	protosocket.ProtosocketGet(cacheName, key)

	// Make sure to close the client
	protosocket.CloseProtosocketCacheClient()
}
