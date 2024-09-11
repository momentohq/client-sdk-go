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

var Version = "1.27.5" // x-release-please-version

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
		headers = append(headers, "runtime-version", "golang:"+runtime.Version())
		headers = append(headers, "agent", "golang:"+string(clientType)+":"+Version)
	}
	return metadata.AppendToOutgoingContext(
		ctx, headers...,
	)
}

func CreateCacheMetadata(ctx context.Context, cacheName string) context.Context {
	return CreateMetadata(ctx, Cache, "cache", cacheName)
}

func CreateStoreMetadata(ctx context.Context, storeName string) context.Context {
	return CreateMetadata(ctx, Store, "store", storeName)
}

func CreateLeaderboardMetadata(ctx context.Context, cacheName string) context.Context {
	return CreateMetadata(ctx, Leaderboard, "cache", cacheName)
}
