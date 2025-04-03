package internal

import (
	"context"
	"runtime"
	"sync"

	"google.golang.org/grpc/metadata"
)

var FirstTimeHeadersSent sync.Map

func init() {
	FirstTimeHeadersSent.Store(Cache, false)
	FirstTimeHeadersSent.Store(Store, false)
	FirstTimeHeadersSent.Store(Leaderboard, false)
	FirstTimeHeadersSent.Store(Topic, false)
	FirstTimeHeadersSent.Store(Ping, false)
	FirstTimeHeadersSent.Store(Auth, false)
}

var Version = "1.34.0" // x-release-please-version

type ClientType string

const (
	Cache       ClientType = "cache"
	Store       ClientType = "store"
	Leaderboard ClientType = "leaderboard"
	Topic       ClientType = "topic"
	Ping        ClientType = "ping"
	Auth        ClientType = "auth"
)

func CreateMetadata(ctx context.Context, clientType ClientType, extraPairs ...string) context.Context {
	headers := extraPairs

	var ftHeadersSent, ok = FirstTimeHeadersSent.Load(clientType)

	if !ok || !ftHeadersSent.(bool) {
		FirstTimeHeadersSent.Store(clientType, true)
		headers = append(
			headers,
			"runtime-version", "golang:"+runtime.Version(),
			"agent", "golang:"+string(clientType)+":"+Version,
		)
	}

	return metadata.AppendToOutgoingContext(
		ctx, headers...,
	)
}

func metadataPairsToStrings(metadataPairs map[string]string) []string {
	pairs := make([]string, 0, len(metadataPairs)*2)
	for k, v := range metadataPairs {
		pairs = append(pairs, k, v)
	}
	return pairs
}

func CreateCacheRequestContextFromMetadataMap(ctx context.Context, cacheName string, metadataPairs map[string]string) context.Context {
	_, ok := metadataPairs["cache"]
	if !ok {
		metadataPairs["cache"] = cacheName
	}
	cacheMetadata := metadataPairsToStrings(metadataPairs)
	return CreateMetadata(ctx, Cache, cacheMetadata...)
}

func CreateTopicRequestContextFromMetadataMap(ctx context.Context, cacheName string, metadataPairs map[string]string) context.Context {
	_, ok := metadataPairs["cache"]
	if !ok {
		metadataPairs["cache"] = cacheName
	}
	cacheMetadata := metadataPairsToStrings(metadataPairs)
	return CreateMetadata(ctx, Topic, cacheMetadata...)
}

func CreateStoreMetadata(ctx context.Context, storeName string) context.Context {
	return CreateMetadata(ctx, Store, "store", storeName)
}

func CreateLeaderboardMetadata(ctx context.Context, cacheName string) context.Context {
	return CreateMetadata(ctx, Leaderboard, "cache", cacheName)
}
