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

const (
	testV1AuthToken       = "eyJhcGlfa2V5IjogImV5SjBlWEFpT2lKS1YxUWlMQ0poYkdjaU9pSklVekkxTmlKOS5leUpwYzNNaU9pSlBibXhwYm1VZ1NsZFVJRUoxYVd4a1pYSWlMQ0pwWVhRaU9qRTJOemd6TURVNE1USXNJbVY0Y0NJNk5EZzJOVFV4TlRReE1pd2lZWFZrSWpvaUlpd2ljM1ZpSWpvaWFuSnZZMnRsZEVCbGVHRnRjR3hsTG1OdmJTSjkuOEl5OHE4NExzci1EM1lDb19IUDRkLXhqSGRUOFVDSXV2QVljeGhGTXl6OCIsICJlbmRwb2ludCI6ICJ0ZXN0Lm1vbWVudG9ocS5jb20ifQ=="
	testV1ApiKey          = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2NzgzMDU4MTIsImV4cCI6NDg2NTUxNTQxMiwiYXVkIjoiIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSJ9.8Iy8q84Lsr-D3YCo_HP4d-xjHdT8UCIuvAYcxhFMyz8"
	testV1MissingEndpoint = "eyJhcGlfa2V5IjogImV5SmxibVJ3YjJsdWRDSTZJbU5sYkd3dE5DMTFjeTEzWlhOMExUSXRNUzV3Y205a0xtRXViVzl0Wlc1MGIyaHhMbU52YlNJc0ltRndhVjlyWlhraU9pSmxlVXBvWWtkamFVOXBTa2xWZWtreFRtbEtPUzVsZVVwNlpGZEphVTlwU25kYVdGSnNURzFrYUdSWVVuQmFXRXBCV2pJeGFHRlhkM1ZaTWpsMFNXbDNhV1J0Vm5sSmFtOTRabEV1VW5OMk9GazVkRE5KVEMwd1RHRjZiQzE0ZDNaSVZESmZZalJRZEhGTlVVMDVRV3hhVlVsVGFrbENieUo5In0="
	testV1MissingApiKey   = "eyJlbmRwb2ludCI6ICJhLmIuY29tIn0="
)

var _ = Describe("CredentialProvider", func() {
	envVar := "v1token"

	BeforeEach(func() {
		if err := os.Setenv(envVar, testV1AuthToken); err != nil {
			Fail(fmt.Sprintf("error setting env var %s: %s\n", envVar, err.Error()))
		}
	})

	AfterEach(func() {
		if err := os.Setenv(envVar, ""); err != nil {
			Fail(fmt.Sprintf("error resetting env var %s: %s\n", envVar, err.Error()))
		}
	})

	Context("V1 auth token", func() {

		It("errors on a invalid auth token", func() {
			badCredentialProvider, err := auth.NewStringMomentoTokenProvider("Invalid token")

			Expect(badCredentialProvider).To(BeNil())
			Expect(err).NotTo(BeNil())

			var momentoErr momentoerrors.MomentoSvcErr
			if errors.As(err, &momentoErr) {
				Expect(momentoErr.Code()).To(Equal(momentoerrors.InvalidArgumentError))
			}
		})

		It("returns a credential provider from an environment variable via constructor", func() {
			credentialProvider, err := auth.NewEnvMomentoTokenProvider(envVar)
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testV1ApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal("cache.test.momentohq.com"))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal("control.test.momentohq.com"))
		})

		It("returns a credential provider from a string via constructor", func() {
			credentialProvider, err := auth.NewStringMomentoTokenProvider(os.Getenv(envVar))
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testV1ApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal("cache.test.momentohq.com"))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal("control.test.momentohq.com"))
		})

		It("returns a credential provider from an environment variable via method", func() {
			credentialProvider, err := auth.FromEnvironmentVariable(envVar)
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testV1ApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal("cache.test.momentohq.com"))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal("control.test.momentohq.com"))
		})

		It("returns a credential provider from a string via method", func() {
			credentialProvider, err := auth.FromString(os.Getenv(envVar))
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testV1ApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal("cache.test.momentohq.com"))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal("control.test.momentohq.com"))
		})

		It("overrides endpoints", func() {
			credentialProvider, err := auth.NewEnvMomentoTokenProvider(envVar)
			Expect(err).To(BeNil())
			controlEndpoint := credentialProvider.GetControlEndpoint()
			cacheEndpoint := credentialProvider.GetCacheEndpoint()
			Expect(controlEndpoint).ToNot(BeEmpty())
			Expect(cacheEndpoint).ToNot(BeEmpty())

			controlEndpoint = fmt.Sprintf("%s-overridden", controlEndpoint)
			cacheEndpoint = fmt.Sprintf("%s-overridden", cacheEndpoint)
			credentialProvider, err = credentialProvider.WithEndpoints(
				auth.Endpoints{
					ControlEndpoint: controlEndpoint,
					CacheEndpoint:   cacheEndpoint,
				},
			)
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetControlEndpoint()).To(Equal(controlEndpoint))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal(cacheEndpoint))
		})

		DescribeTable("errors when v1 token is missing data",
			func(envVarValue string) {
				credentialProvider, err := auth.FromString(envVarValue)
				Expect(credentialProvider).To(BeNil())
				Expect(err).To(Not(BeNil()))
				var momentoErr momentoerrors.MomentoSvcErr
				if errors.As(err, &momentoErr) {
					Expect(momentoErr.Code()).To(Equal(momentoerrors.InvalidArgumentError))
				} else {
					Fail(fmt.Sprintf("unknown error: %s", err.Error()))
				}
			},
			Entry("missing endpoint", testV1MissingEndpoint),
			Entry("missing api key", testV1MissingApiKey),
		)
	})

})
