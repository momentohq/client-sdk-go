package internal

import (
	"context"
	"fmt"
	"runtime"

	"google.golang.org/grpc/metadata"
)

var FirstTimeHeadersSent = false
var Version = "1.24.0" // x-release-please-version

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

	if !FirstTimeHeadersSent {
		FirstTimeHeadersSent = true
		headers = append(headers, "runtime-version", "golang:"+runtime.Version())
		headers = append(headers, "agent", "golang:"+string(clientType)+":"+Version)
		fmt.Println("set first time headers")
	}
	newCtx := metadata.AppendToOutgoingContext(
		ctx, headers...,
	)
	updatedMetadata, ok := metadata.FromOutgoingContext(newCtx)
	if !ok {
		panic("metadata not found in newCtx")
	}
	fmt.Println("metadata function newCtx", updatedMetadata)
	return newCtx
}

func CreateMetadataHeaders(clientType ClientType, extraPairs ...string) []string {
	headers := extraPairs

	if !FirstTimeHeadersSent {
		FirstTimeHeadersSent = true
		headers = append(headers, "runtime-version", "golang:"+runtime.Version())
		headers = append(headers, "agent", "golang:"+string(clientType)+":"+Version)
		fmt.Println("set first time headers")
	}
	return headers
}
