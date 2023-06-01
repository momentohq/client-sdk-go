package main

import (
	"context"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/momento"
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
	if (err != nil) {
		panic(err)
	}
}

func example_API_CreateCache() {
  _, err := client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: "cache-name",
	});
  if (err != nil) {
	panic(err)
  }
}
