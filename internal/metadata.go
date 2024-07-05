package internal

import (
	"context"
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
	}
	return metadata.AppendToOutgoingContext(
		ctx, headers...,
	)
}
