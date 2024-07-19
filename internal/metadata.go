package internal

import (
	"context"
	"runtime"

	"google.golang.org/grpc/metadata"
)

var FirstTimeHeadersSent = map[ClientType]bool{
	Cache:       false,
	Store:       false,
	Leaderboard: false,
	Topic:       false,
	Ping:        false,
	Auth:        false,
}
var Version = "1.26.0" // x-release-please-version

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

	if !FirstTimeHeadersSent[clientType] {
		FirstTimeHeadersSent[clientType] = true
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
