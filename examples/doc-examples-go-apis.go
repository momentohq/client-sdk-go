package main

import "github.com/momentohq/client-sdk-go/auth"

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