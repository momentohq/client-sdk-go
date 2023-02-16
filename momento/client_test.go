package momento_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	. "github.com/momentohq/client-sdk-go/momento"
)

var _ = Describe("SimpleCacheClient", func() {
	var clientProps SimpleCacheClientProps
	var credentialProvider auth.CredentialProvider
	var configuration config.Configuration

	BeforeEach(func() {
		credentialProvider, _ = auth.NewEnvMomentoTokenProvider("TEST_AUTH_TOKEN")
		configuration = config.LatestLaptopConfig()

		clientProps = SimpleCacheClientProps{
			CredentialProvider: credentialProvider,
			Configuration:      configuration,
			DefaultTTL:         100 * time.Second,
		}
	})

	It(`errors on an invalid TTL`, func() {
		clientProps.DefaultTTL = 0 * time.Second
		client, err := NewSimpleCacheClient(&clientProps)

		Expect(client).To(BeNil())
		Expect(err).NotTo(BeNil())
		var momentoErr MomentoError
		if errors.As(err, &momentoErr) {
			Expect(momentoErr.Code()).To(Equal(InvalidArgumentError))
		}
	})

	It(`errors on invalid timeout`, func() {
		badRequestTimeout := 0 * time.Second
		clientProps.Configuration = config.LatestLaptopConfig().WithClientTimeout(badRequestTimeout)
		client, err := NewSimpleCacheClient(&clientProps)

		Expect(client).To(BeNil())
		Expect(err).NotTo(BeNil())
		var momentoErr MomentoError
		if errors.As(err, &momentoErr) {
			Expect(momentoErr.Code()).To(Equal(InvalidArgumentError))
		}
	})
})
