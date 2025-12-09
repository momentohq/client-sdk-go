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
	v2KeyEnvVar           = "MOMENTO_TEST_V2_API_KEY"
	testV2ApiKey          = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ0IjoiZyIsImlkIjoic29tZS1pZCJ9.WRhKpdh7cFCXO7lAaVojtQAxK6mxMdBrvXTJL1xu94S0d6V1YSstOObRlAIMA7i_yIxO1mWEF3rlF5UNc77VXQ"
	testEndpoint          = "testEndpoint"
	endpointEnvVar        = "MOMENTO_TEST_ENDPOINT"
	testPreV1ApiKey       = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyQHRlc3QuY29tIiwiY3AiOiJjb250cm9sLnRlc3QuY29tIiwiYyI6ImNhY2hlLnRlc3QuY29tIn0.c0Z8Ipetl6raCNHSHs7Mpq3qtWkFy4aLvGhIFR4CoR0OnBdGbdjN-4E58bAabrSGhRA8-B2PHzgDd4JF4clAzg"
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

	Context("V2 api keys", func() {
		It("returns a credential provider from an environment variables via constructor", func() {
			os.Setenv(endpointEnvVar, testEndpoint)
			os.Setenv(v2KeyEnvVar, testV2ApiKey)
			credentialProvider, err := auth.NewEnvVarV2TokenProvider(auth.FromEnvVarV2Props{
				ApiKeyEnvVar:   v2KeyEnvVar,
				EndpointEnvVar: endpointEnvVar,
			})
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testV2ApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal(fmt.Sprintf("cache.%s:443", testEndpoint)))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal(fmt.Sprintf("control.%s:443", testEndpoint)))
			Expect(credentialProvider.GetTokenEndpoint()).To(Equal(fmt.Sprintf("token.%s:443", testEndpoint)))
		})

		It("returns a credential provider from a string via constructor", func() {
			credentialProvider, err := auth.NewApiKeyV2TokenProvider(auth.FromApiKeyV2Props{
				ApiKey:   testV2ApiKey,
				Endpoint: testEndpoint,
			})
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testV2ApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal(fmt.Sprintf("cache.%s:443", testEndpoint)))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal(fmt.Sprintf("control.%s:443", testEndpoint)))
			Expect(credentialProvider.GetTokenEndpoint()).To(Equal(fmt.Sprintf("token.%s:443", testEndpoint)))
		})

		It("returns a credential provider from an environment variable via method", func() {
			os.Setenv(v2KeyEnvVar, testV2ApiKey)
			os.Setenv(endpointEnvVar, testEndpoint)
			credentialProvider, err := auth.FromEnvVarV2(auth.FromEnvVarV2Props{
				ApiKeyEnvVar:   v2KeyEnvVar,
				EndpointEnvVar: endpointEnvVar,
			})
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testV2ApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal(fmt.Sprintf("cache.%s:443", testEndpoint)))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal(fmt.Sprintf("control.%s:443", testEndpoint)))
			Expect(credentialProvider.GetTokenEndpoint()).To(Equal(fmt.Sprintf("token.%s:443", testEndpoint)))
		})

		It("returns a credential provider from a string via method", func() {
			credentialProvider, err := auth.FromApiKeyV2(auth.FromApiKeyV2Props{
				ApiKey:   testV2ApiKey,
				Endpoint: testEndpoint,
			})
			Expect(err).To(BeNil())
			Expect(credentialProvider.GetAuthToken()).To(Equal(testV2ApiKey))
			Expect(credentialProvider.GetCacheEndpoint()).To(Equal(fmt.Sprintf("cache.%s:443", testEndpoint)))
			Expect(credentialProvider.GetControlEndpoint()).To(Equal(fmt.Sprintf("control.%s:443", testEndpoint)))
			Expect(credentialProvider.GetTokenEndpoint()).To(Equal(fmt.Sprintf("token.%s:443", testEndpoint)))
		})

		DescribeTable("string method errors when missing data",
			func(apiKey string, endpoint string, expectedError string) {
				credentialProvider, err := auth.FromApiKeyV2(auth.FromApiKeyV2Props{
					ApiKey:   apiKey,
					Endpoint: endpoint,
				})
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
			Entry("empty endpoint", testV2ApiKey, "", "Endpoint is an empty string"),
		)

		DescribeTable("env var method errors when missing data",
			func(apiKeyEnvVar string, endpointEnvVar string, expectedError string) {
				credentialProvider, err := auth.FromEnvVarV2(auth.FromEnvVarV2Props{
					ApiKeyEnvVar:   apiKeyEnvVar,
					EndpointEnvVar: endpointEnvVar,
				})
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
			Entry("empty env var name", "", endpointEnvVar, "API key environment variable name is empty"),
			Entry("env var not set", "NON_EXISTENT_ENV_VAR", endpointEnvVar, "Missing required environment variable NON_EXISTENT_ENV_VAR"),
			Entry("empty endpoint", v2KeyEnvVar, "", "Endpoint environment variable name is empty"),
			Entry("env var not set", v2KeyEnvVar, "NON_EXISTENT_ENDPOINT", "Missing required environment variable NON_EXISTENT_ENDPOINT"),
		)

		It("errors when v2 api key is provided to FromDisposableToken", func() {
			credentialProvider, err := auth.FromDisposableToken(testV2ApiKey)
			Expect(credentialProvider).To(BeNil())
			Expect(err).To(Not(BeNil()))
			var momentoErr momentoerrors.MomentoSvcErr
			if errors.As(err, &momentoErr) {
				Expect(momentoErr.Code()).To(Equal(momentoerrors.InvalidArgumentError))
			} else {
				Fail(fmt.Sprintf("unknown error: %s", err.Error()))
			}
		})

		It("errors when v2 api key is provided to FromString", func() {
			credentialProvider, err := auth.FromString(testV2ApiKey)
			Expect(credentialProvider).To(BeNil())
			Expect(err).To(Not(BeNil()))
			var momentoErr momentoerrors.MomentoSvcErr
			if errors.As(err, &momentoErr) {
				Expect(momentoErr.Code()).To(Equal(momentoerrors.InvalidArgumentError))
			} else {
				Fail(fmt.Sprintf("unknown error: %s", err.Error()))
			}
		})

		It("errors when v2 api key is provided to FromEnvironmentVariable", func() {
			os.Setenv(v2KeyEnvVar, testV2ApiKey)
			credentialProvider, err := auth.FromEnvironmentVariable(v2KeyEnvVar)
			Expect(credentialProvider).To(BeNil())
			Expect(err).To(Not(BeNil()))
			var momentoErr momentoerrors.MomentoSvcErr
			if errors.As(err, &momentoErr) {
				Expect(momentoErr.Code()).To(Equal(momentoerrors.InvalidArgumentError))
			} else {
				Fail(fmt.Sprintf("unknown error: %s", err.Error()))
			}
		})

		It("errors when v1 api key is provided to FromApiKeyV2", func() {
			credentialProvider, err := auth.FromApiKeyV2(auth.FromApiKeyV2Props{
				ApiKey:   testV1ApiKey,
				Endpoint: testEndpoint,
			})
			Expect(credentialProvider).To(BeNil())
			Expect(err).To(Not(BeNil()))
			var momentoErr momentoerrors.MomentoSvcErr
			if errors.As(err, &momentoErr) {
				Expect(momentoErr.Code()).To(Equal(momentoerrors.InvalidArgumentError))
			} else {
				Fail(fmt.Sprintf("unknown error: %s", err.Error()))
			}
		})

		It("errors when v1 api key is provided to FromEnvVarV2", func() {
			os.Setenv(v2KeyEnvVar, testV1ApiKey)
			os.Setenv(endpointEnvVar, testEndpoint)
			credentialProvider, err := auth.FromEnvVarV2(auth.FromEnvVarV2Props{
				ApiKeyEnvVar:   v2KeyEnvVar,
				EndpointEnvVar: endpointEnvVar,
			})
			Expect(credentialProvider).To(BeNil())
			Expect(err).To(Not(BeNil()))
			var momentoErr momentoerrors.MomentoSvcErr
			if errors.As(err, &momentoErr) {
				Expect(momentoErr.Code()).To(Equal(momentoerrors.InvalidArgumentError))
			} else {
				Fail(fmt.Sprintf("unknown error: %s", err.Error()))
			}
		})

		It("errors when pre-v1 api key is provided to FromApiKeyV2", func() {
			credentialProvider, err := auth.FromApiKeyV2(auth.FromApiKeyV2Props{
				ApiKey:   testPreV1ApiKey,
				Endpoint: testEndpoint,
			})
			Expect(credentialProvider).To(BeNil())
			Expect(err).To(Not(BeNil()))
			var momentoErr momentoerrors.MomentoSvcErr
			if errors.As(err, &momentoErr) {
				Expect(momentoErr.Code()).To(Equal(momentoerrors.InvalidArgumentError))
			} else {
				Fail(fmt.Sprintf("unknown error: %s", err.Error()))
			}
		})

		It("errors when pre-v1 api key is provided to FromEnvVarV2", func() {
			os.Setenv(v2KeyEnvVar, testPreV1ApiKey)
			os.Setenv(endpointEnvVar, testEndpoint)
			credentialProvider, err := auth.FromEnvVarV2(auth.FromEnvVarV2Props{
				ApiKeyEnvVar:   v2KeyEnvVar,
				EndpointEnvVar: endpointEnvVar,
			})
			Expect(credentialProvider).To(BeNil())
			Expect(err).To(Not(BeNil()))
			var momentoErr momentoerrors.MomentoSvcErr
			if errors.As(err, &momentoErr) {
				Expect(momentoErr.Code()).To(Equal(momentoerrors.InvalidArgumentError))
			} else {
				Fail(fmt.Sprintf("unknown error: %s", err.Error()))
			}
		})
	})

})
