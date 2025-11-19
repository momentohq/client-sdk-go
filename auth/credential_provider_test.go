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

var _ = Describe("auth credential-provider", func() {
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
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal("cache.test.momentohq.com:443"))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal("control.test.momentohq.com:443"))
		})

		It("returns a credential provider from a string via constructor", func() {
			credentialProvider, err := auth.NewStringMomentoTokenProvider(os.Getenv(envVar))
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testV1ApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal("cache.test.momentohq.com:443"))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal("control.test.momentohq.com:443"))
		})

		It("returns a credential provider from an environment variable via method", func() {
			credentialProvider, err := auth.FromEnvironmentVariable(envVar)
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testV1ApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal("cache.test.momentohq.com:443"))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal("control.test.momentohq.com:443"))
		})

		It("returns a credential provider from a string via method", func() {
			credentialProvider, err := auth.FromString(os.Getenv(envVar))
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testV1ApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal("cache.test.momentohq.com:443"))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal("control.test.momentohq.com:443"))
		})

		It("overrides endpoints", func() {
			credentialProvider, err := auth.NewEnvMomentoTokenProvider(envVar)
			Expect(err).To(BeNil())
			controlEndpoint := credentialProvider.GetControlEndpoint()
			cacheEndpoint := credentialProvider.GetCacheEndpoint()
			Expect(controlEndpoint).ToNot(BeEmpty())
			Expect(cacheEndpoint).ToNot(BeEmpty())
			Expect(credentialProvider.IsControlEndpointSecure()).To(BeTrue())
			Expect(credentialProvider.IsCacheEndpointSecure()).To(BeTrue())

			controlEndpoint = fmt.Sprintf("%s-overridden", controlEndpoint)
			cacheEndpoint = fmt.Sprintf("%s-overridden", cacheEndpoint)
			credentialProvider, err = credentialProvider.WithEndpoints(
				auth.AllEndpoints{
					ControlEndpoint: auth.Endpoint{Endpoint: controlEndpoint},
					CacheEndpoint:   auth.Endpoint{Endpoint: cacheEndpoint},
				},
			)
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetControlEndpoint()).To(Equal(controlEndpoint))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal(cacheEndpoint))
			Expect(credentialProvider.IsControlEndpointSecure()).To(BeTrue())
			Expect(credentialProvider.IsCacheEndpointSecure()).To(BeTrue())
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

		It("correctly sets Momento Local endpoints", func() {
			// Using default config
			credentialProvider, err := auth.NewMomentoLocalProvider(nil)
			defaultEndpoint := "127.0.0.1:8080"
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(""))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal(defaultEndpoint))
			Expect(credentialProvider.IsCacheEndpointSecure()).To(BeFalse())
			Expect(credentialProvider.GetControlEndpoint()).To(Equal(defaultEndpoint))
			Expect(credentialProvider.IsControlEndpointSecure()).To(BeFalse())
			Expect(credentialProvider.GetTokenEndpoint()).To(Equal(defaultEndpoint))
			Expect(credentialProvider.IsTokenEndpointSecure()).To(BeFalse())

			// Using provided config
			credentialProvider, err = auth.NewMomentoLocalProvider(&auth.MomentoLocalConfig{
				Hostname: "localhost",
				Port:     9090,
			})
			nonDefaultEndpoint := "localhost:9090"
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(""))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal(nonDefaultEndpoint))
			Expect(credentialProvider.IsCacheEndpointSecure()).To(BeFalse())
			Expect(credentialProvider.GetControlEndpoint()).To(Equal(nonDefaultEndpoint))
			Expect(credentialProvider.IsControlEndpointSecure()).To(BeFalse())
			Expect(credentialProvider.GetTokenEndpoint()).To(Equal(nonDefaultEndpoint))
			Expect(credentialProvider.IsTokenEndpointSecure()).To(BeFalse())
		})
	})

	Context("Global api keys", func() {
		testEnvVar := "MOMENTO_TEST_GLOBAL_API_KEY"
		testApiKey := "testToken"
		testEndpoint := "testEndpoint"

		It("returns a credential provider from an environment variable via constructor", func() {
			os.Setenv(testEnvVar, testApiKey)
			credentialProvider, err := auth.NewGlobalEnvMomentoTokenProvider(testEnvVar, testEndpoint)
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal(fmt.Sprintf("cache.%s:443", testEndpoint)))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal(fmt.Sprintf("control.%s:443", testEndpoint)))
			Expect(credentialProvider.GetTokenEndpoint()).To(Equal(fmt.Sprintf("token.%s:443", testEndpoint)))
		})

		It("returns a credential provider from a string via constructor", func() {
			credentialProvider, err := auth.NewGlobalStringMomentoTokenProvider(testApiKey, testEndpoint)
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal(fmt.Sprintf("cache.%s:443", testEndpoint)))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal(fmt.Sprintf("control.%s:443", testEndpoint)))
			Expect(credentialProvider.GetTokenEndpoint()).To(Equal(fmt.Sprintf("token.%s:443", testEndpoint)))
		})

		It("returns a credential provider from an environment variable via method", func() {
			credentialProvider, err := auth.GlobalKeyFromEnvironmentVariable(testEnvVar, testEndpoint)
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal(fmt.Sprintf("cache.%s:443", testEndpoint)))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal(fmt.Sprintf("control.%s:443", testEndpoint)))
			Expect(credentialProvider.GetTokenEndpoint()).To(Equal(fmt.Sprintf("token.%s:443", testEndpoint)))
		})

		It("returns a credential provider from a string via method", func() {
			credentialProvider, err := auth.GlobalKeyFromString(testApiKey, testEndpoint)
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal(fmt.Sprintf("cache.%s:443", testEndpoint)))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal(fmt.Sprintf("control.%s:443", testEndpoint)))
			Expect(credentialProvider.GetTokenEndpoint()).To(Equal(fmt.Sprintf("token.%s:443", testEndpoint)))
		})

		DescribeTable("string method errors when missing data",
			func(apiKey string, endpoint string, expectedError string) {
				credentialProvider, err := auth.GlobalKeyFromString(apiKey, endpoint)
				Expect(credentialProvider).To(BeNil())
				Expect(err).To(Not(BeNil()))
				var momentoErr momentoerrors.MomentoSvcErr
				if errors.As(err, &momentoErr) {
					Expect(momentoErr.Code()).To(Equal(momentoerrors.InvalidArgumentError))
					Expect(momentoErr.Error()).To(ContainSubstring(expectedError))
				} else {
					Fail(fmt.Sprintf("unknown error: %s", err.Error()))
				}
			},
			Entry("empty key", "", testEndpoint, "Auth token is an empty string"),
			Entry("empty endpoint", testApiKey, "", "Endpoint is an empty string"),
		)

		DescribeTable("env var method errors when missing data",
			func(envVarName string, endpoint string, expectedError string) {
				credentialProvider, err := auth.GlobalKeyFromEnvironmentVariable(envVarName, endpoint)
				Expect(credentialProvider).To(BeNil())
				Expect(err).To(Not(BeNil()))
				var momentoErr momentoerrors.MomentoSvcErr
				if errors.As(err, &momentoErr) {
					Expect(momentoErr.Code()).To(Equal(momentoerrors.InvalidArgumentError))
					Expect(momentoErr.Error()).To(ContainSubstring(expectedError))
				} else {
					Fail(fmt.Sprintf("unknown error: %s", err.Error()))
				}
			},
			Entry("empty env var name", "", testEndpoint, "Environment variable name is empty"),
			Entry("env var not set", "NON_EXISTENT_ENV_VAR", testEndpoint, "Missing required environment variable NON_EXISTENT_ENV_VAR"),
			Entry("empty endpoint", testEnvVar, "", "Endpoint is an empty string"),
		)
	})

})
