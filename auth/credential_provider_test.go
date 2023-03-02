package auth_test

import (
	"errors"
	"fmt"
	"os"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CredentialProvider", func() {
	It(`errors on a invalid auth token`, func() {
		badCredentialProvider, err := auth.NewStringMomentoTokenProvider("Invalid token")

		Expect(badCredentialProvider).To(BeNil())
		Expect(err).NotTo(BeNil())

		var momentoErr momentoerrors.MomentoSvcErr
		if errors.As(err, &momentoErr) {
			Expect(momentoErr.Code()).To(Equal(momentoerrors.InvalidArgumentError))
		}
	})

	It("returns a credential provider from an environment variable via constructor", func() {
		credentialProvider, err := auth.NewEnvMomentoTokenProvider("TEST_AUTH_TOKEN")
		Expect(err).To(BeNil())
		Expect(credentialProvider.GetAuthToken()).To(Equal(os.Getenv("TEST_AUTH_TOKEN")))
	})

	It("returns a credential provider from a string via constructor", func() {
		credentialProvider, err := auth.NewStringMomentoTokenProvider(os.Getenv("TEST_AUTH_TOKEN"))
		Expect(err).To(BeNil())
		Expect(credentialProvider.GetAuthToken()).To(Equal(os.Getenv("TEST_AUTH_TOKEN")))
	})

	It("returns a credential provider from an environment variable via method", func() {
		credentialProvider, err := auth.FromEnvironmentVariable("TEST_AUTH_TOKEN")
		Expect(err).To(BeNil())
		Expect(credentialProvider.GetAuthToken()).To(Equal(os.Getenv("TEST_AUTH_TOKEN")))
	})

	It("returns a credential provider from a string via method", func() {
		credentialProvider, err := auth.FromString(os.Getenv("TEST_AUTH_TOKEN"))
		Expect(err).To(BeNil())
		Expect(credentialProvider.GetAuthToken()).To(Equal(os.Getenv("TEST_AUTH_TOKEN")))
	})

	It("overrides endpoints", func() {
		credentialProvider, err := auth.NewEnvMomentoTokenProvider("TEST_AUTH_TOKEN")
		Expect(err).To(BeNil())
		controlEndpoint := credentialProvider.GetControlEndpoint()
		cacheEndpoint := credentialProvider.GetCacheEndpoint()
		Expect(controlEndpoint).ToNot(BeEmpty())
		Expect(cacheEndpoint).ToNot(BeEmpty())

		controlEndpoint = fmt.Sprintf("%s-overridden", controlEndpoint)
		cacheEndpoint = fmt.Sprintf("%s-overridden", cacheEndpoint)
		credentialProvider, err = credentialProvider.WithEndpoints(
			&auth.Endpoints{
				ControlEndpoint: controlEndpoint,
				CacheEndpoint:   cacheEndpoint,
			},
		)
		Expect(err).To(BeNil())
		Expect(credentialProvider.GetControlEndpoint()).To(Equal(controlEndpoint))
		Expect(credentialProvider.GetCacheEndpoint()).To(Equal(cacheEndpoint))
	})

})
