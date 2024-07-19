package momento_test

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/internal"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	responses "github.com/momentohq/client-sdk-go/responses"
	auth_responses "github.com/momentohq/client-sdk-go/responses/auth"
	"github.com/momentohq/client-sdk-go/utils"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
)

func authClientFromApiKey(ctx SharedContext, apiKey string, endpoint string) AuthClient {
	credProvider := credProviderFromApiKey(apiKey, endpoint)
	authClient, err := NewAuthClient(ctx.AuthConfiguration, credProvider)
	if err != nil {
		panic(err)
	}
	return authClient
}

func credProviderFromApiKey(apiKey string, endpoint string) auth.CredentialProvider {
	credProviderWithoutEndpoints, err := auth.NewStringMomentoTokenProvider(apiKey)
	if err != nil {
		panic(err)
	}

	credProviderWithEndpoints, err := credProviderWithoutEndpoints.WithEndpoints(
		auth.Endpoints{
			ControlEndpoint: fmt.Sprintf("control.%s", endpoint),
			CacheEndpoint:   fmt.Sprintf("cache.%s", endpoint),
			TokenEndpoint:   fmt.Sprintf("cache.%s", endpoint),
		},
	)
	if err != nil {
		panic(err)
	}

	return credProviderWithEndpoints
}

func credProviderFromDisposableToken(resp auth_responses.GenerateDisposableTokenResponse) auth.CredentialProvider {
	success := resp.(*auth_responses.GenerateDisposableTokenSuccess)
	credProviderWithoutEndpoints, err := auth.NewStringMomentoTokenProvider(success.ApiKey)

	if err != nil {
		panic(err)
	}
	credProviderWithEndpoints, err := credProviderWithoutEndpoints.WithEndpoints(
		auth.Endpoints{
			ControlEndpoint: fmt.Sprintf("control.%s", success.Endpoint),
			CacheEndpoint:   fmt.Sprintf("cache.%s", success.Endpoint),
			TokenEndpoint:   fmt.Sprintf("cache.%s", success.Endpoint),
		},
	)

	if err != nil {
		panic(err)
	}

	return credProviderWithEndpoints
}

func assertGetSuccess(cc CacheClient, key Value, cacheName string) {
	_, err := cc.Get(context.Background(), &GetRequest{Key: key, CacheName: cacheName})
	if err != nil {
		panic(err)
	}
}

func assertGetFailure(cc CacheClient, key Value, cacheName string) {
	_, err := cc.Get(context.Background(), &GetRequest{Key: key, CacheName: cacheName})
	if err == nil {
		panic("expected Get to fail but it succeeded")
	}
}

func assertSetSuccess(cc CacheClient, key Value, cacheName string) {
	_, err := cc.Set(context.Background(), &SetRequest{Key: key, Value: String(uuid.NewString()), CacheName: cacheName})
	if err != nil {
		panic(err)
	}
}

func assertSetFailure(cc CacheClient, key Value, cacheName string) {
	_, err := cc.Set(context.Background(), &SetRequest{Key: key, Value: String(uuid.NewString()), CacheName: cacheName})
	if err == nil {
		panic("expected Set to fail but it succeeded")
	}
}

func assertPublishSuccess(tc TopicClient, topicName string, cacheName string) {
	_, err := tc.Publish(context.Background(), &TopicPublishRequest{TopicName: topicName, Value: String(uuid.NewString()), CacheName: cacheName})
	if err != nil {
		panic(err)
	}
}

func assertPublishFailure(tc TopicClient, topicName string, cacheName string) {
	_, err := tc.Publish(context.Background(), &TopicPublishRequest{TopicName: topicName, Value: String(uuid.NewString()), CacheName: cacheName})
	if err == nil {
		panic("expected publish to fail but it succeeded")
	}
}

func assertSubscribeFailure(tc TopicClient, topicName string, cacheName string) {
	_, err := tc.Subscribe(context.Background(), &TopicSubscribeRequest{TopicName: topicName, CacheName: cacheName})
	if err == nil {
		panic("expected subscribe to fail but it succeeded")
	}
}

func assertSubscribeSuccess(tc TopicClient, topicName string, cacheName string) {
	_, err := tc.Subscribe(context.Background(), &TopicSubscribeRequest{TopicName: topicName, CacheName: cacheName})
	if err != nil {
		panic(err)
	}
}

func generateDisposableTokenSuccess(ctx SharedContext, scope DisposableTokenScope) auth_responses.GenerateDisposableTokenResponse {
	expiresIn := utils.ExpiresInMinutes(5)
	resp, err := ctx.AuthClient.GenerateDisposableToken(ctx.Ctx, &GenerateDisposableTokenRequest{
		ExpiresIn: expiresIn,
		Scope:     scope,
	})
	if err != nil {
		panic(err)
	}
	Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateDisposableTokenSuccess{}))
	return resp
}

func newCacheClient(ctx SharedContext, provider auth.CredentialProvider) CacheClient {
	cc, err := NewCacheClient(ctx.Configuration, provider, time.Second*60)
	if err != nil {
		panic(err)
	}
	return cc
}

func newTopicClient(ctx SharedContext, provider auth.CredentialProvider) TopicClient {
	tc, err := NewTopicClient(ctx.TopicConfiguration, provider)
	if err != nil {
		panic(err)
	}
	return tc
}

var _ = Describe("auth auth-client", func() {
	var sharedContext SharedContext

	BeforeEach(func() {
		sharedContext = NewSharedContext()
		sharedContext.CreateDefaultCaches()

		DeferCleanup(func() {
			sharedContext.Close()
		})
	})

	Describe("Generate disposable tokens", func() {
		Describe("CacheKeyReadOnly tokens", func() {
			It(`Generates disposable token CacheKeyReadOnly AllCaches, and validates its permissions`, func() {
				key := String(uuid.NewString())
				resp := generateDisposableTokenSuccess(sharedContext, CacheKeyReadOnly(AllCaches{}, key))
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we can read the key
				assertGetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we can read the key for another cache
				assertGetSuccess(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot read another key
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot read another key for another cache
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can not write the key
				assertSetFailure(cc, key, sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})

			It(`Generates disposable token CacheKeyReadOnly for a specific cache, and validates its permissions`, func() {
				key := String(uuid.NewString())
				resp := generateDisposableTokenSuccess(sharedContext, CacheKeyReadOnly(CacheName{Name: sharedContext.CacheName}, key))
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we can read the key
				assertGetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we cannot read the key for another cache
				assertGetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot read another key
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot read another key for another cache
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can not write the key
				assertSetFailure(cc, key, sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})
		})

		Describe("CacheKeyWriteOnly tokens", func() {
			It(`Generates disposable token CacheKeyWriteOnly for a specific cache, and validates its permissions`, func() {
				key := String(uuid.NewString())
				scope := CacheKeyWriteOnly(CacheName{Name: sharedContext.CacheName}, key)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot read the key for another cache
				assertGetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot read another key
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot read another key for another cache
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write the key
				assertSetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we cannot write another key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})

			It(`Generates disposable token CacheKeyWriteOnly for all caches, and validates its permissions`, func() {
				key := String(uuid.NewString())
				scope := CacheKeyWriteOnly(AllCaches{}, key)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot read the key for another cache
				assertGetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot read another key
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot read another key for another cache
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write the key
				assertSetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we can write the key in another cache
				assertSetSuccess(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write another key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})
		})

		Describe("CacheKeyReadWrite tokens", func() {
			It(`Generates disposable token CacheKeyReadWrite for all caches, and validates its permissions`, func() {
				key := String(uuid.NewString())
				scope := CacheKeyReadWrite(AllCaches{}, key)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cann read the key
				assertGetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we cann read the key for another cache
				assertGetSuccess(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot read another key
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot read another key for another cache
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write the key
				assertSetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we can write the key in another cache
				assertSetSuccess(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write another key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot write another key in another cache
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})

			It(`Generates disposable token CacheKeyReadWrite for a specific cache, and validates its permissions`, func() {
				key := String(uuid.NewString())
				scope := CacheKeyReadWrite(CacheName{Name: sharedContext.CacheName}, key)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we can read the key
				assertGetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we cannot read the key for another cache
				assertGetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot read another key
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot read another key for another cache
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write the key
				assertSetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write another key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot write another key in another cache
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})
		})

		Describe("CacheKeyPrefixReadWrite tokens", func() {
			It(`Generates disposable token CacheKeyPrefixReadWrite for a specific cache, and validates its permissions`, func() {
				key := String(uuid.NewString())
				scope := CacheKeyPrefixReadWrite(CacheName{Name: sharedContext.CacheName}, key)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we can read the key
				assertGetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we cannot read the key for another cache
				assertGetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can read a prefixed key
				assertGetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot read a prefixed key for another cache
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write the key
				assertSetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can write a prefixed key
				assertSetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot write a prefixed key in another cache
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})

			It(`Generates disposable token CacheKeyPrefixReadWrite for all caches, and validates its permissions`, func() {
				key := String(uuid.NewString())
				scope := CacheKeyPrefixReadWrite(AllCaches{}, key)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we can read the key
				assertGetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we can read the key for another cache
				assertGetSuccess(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can read a prefixed key
				assertGetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we can read a prefixed key for another cache
				assertGetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write the key
				assertSetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we can write the key in another cache
				assertSetSuccess(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can write a prefixed key
				assertSetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we can write a prefixed key in another cache
				assertSetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})
		})

		Describe("CacheKeyPrefixReadOnly tokens", func() {
			It(`Generates disposable token CacheKeyPrefixReadOnly for a specific cache, and validates its permissions`, func() {
				key := String(uuid.NewString())
				scope := CacheKeyPrefixReadOnly(CacheName{Name: sharedContext.CacheName}, key)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we can read the key
				assertGetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we cannot read the key for another cache
				assertGetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can read a prefixed key
				assertGetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot read a prefixed key for another cache
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot write a prefixed key in another cache
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})

			It(`Generates disposable token CacheKeyPrefixReadOnly for all caches, and validates its permissions`, func() {
				key := String(uuid.NewString())
				scope := CacheKeyPrefixReadOnly(AllCaches{}, key)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we can read the key
				assertGetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we can read the key for another cache
				assertGetSuccess(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can read a prefixed key
				assertGetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we can read a prefixed key for another cache
				assertGetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot write a prefixed key in another cache
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})
		})

		Describe("CacheKeyPrefixWriteOnly tokens", func() {
			It(`Generates disposable token CacheKeyPrefixWriteOnly for a specific cache, and validates its permissions`, func() {
				key := String(uuid.NewString())
				scope := CacheKeyPrefixWriteOnly(CacheName{Name: sharedContext.CacheName}, key)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot read the key for another cache
				assertGetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot read a prefixed key
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot read a prefixed key for another cache
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write the key
				assertSetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can write a prefixed key
				assertSetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot write a prefixed key in another cache
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we cannot write a key with another prefix
				assertSetFailure(cc, String("w00t"), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})

			It(`Generates disposable token CacheKeyPrefixWriteOnly for all caches, and validates its permissions`, func() {
				key := String(uuid.NewString())
				scope := CacheKeyPrefixWriteOnly(AllCaches{}, key)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot read the key for another cache
				assertGetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot read a prefixed key
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot read a prefixed key for another cache
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write the key
				assertSetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we can write the key in another cache
				assertSetSuccess(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can write a prefixed key
				assertSetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we can write a prefixed key in another cache
				assertSetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we cannot write a key with another prefix
				assertSetFailure(cc, String("w00t"), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})
		})

		Describe("TopicNamePrefixPublishSubscribe tokens", func() {
			It(`Generates disposable token TopicNamePrefixPublishSubscribe for a specific cache, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := TopicNamePrefixPublishSubscribe(CacheName{Name: sharedContext.CacheName}, topicName)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				// make sure can publish to topic
				assertPublishSuccess(tc, topicName, sharedContext.CacheName)
				// make sure can publish to prefixed topic
				assertPublishSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot publish to topic in another cache
				assertPublishFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot publish to a topic with a different prefix
				assertPublishFailure(tc, "w00t", sharedContext.CacheName)
				// make sure can subscribe to topic
				assertSubscribeSuccess(tc, topicName, sharedContext.CacheName)
				// make sure can subscribe to prefixed topic
				assertSubscribeSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot subscribe to topic in another cache
				assertSubscribeFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot subscribe to a topic with a different prefix
				assertSubscribeFailure(tc, "w00t", sharedContext.CacheName)
			})

			It(`Generates disposable token TopicNamePrefixPublishSubscribe for all caches, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := TopicNamePrefixPublishSubscribe(AllCaches{}, topicName)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				// make sure can publish to topic
				assertPublishSuccess(tc, topicName, sharedContext.CacheName)
				// make sure can publish to prefixed topic
				assertPublishSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure can publish to topic in another cache
				assertPublishSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot publish to a topic with a different prefix
				assertPublishFailure(tc, "w00t", sharedContext.CacheName)
				// make sure can subscribe to topic
				assertSubscribeSuccess(tc, topicName, sharedContext.CacheName)
				// make sure can subscribe to prefixed topic
				assertSubscribeSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure can subscribe to topic in another cache
				assertSubscribeSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot subscribe to a topic with a different prefix
				assertSubscribeFailure(tc, "w00t", sharedContext.CacheName)
			})
		})

		Describe("TopicNamePrefixPublishOnly tokens", func() {
			It(`Generates disposable token TopicNamePrefixPublishOnly for a specific cache, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := TopicNamePrefixPublishOnly(CacheName{Name: sharedContext.CacheName}, topicName)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				// make sure can publish to topic
				assertPublishSuccess(tc, topicName, sharedContext.CacheName)
				// make sure can publish to prefixed topic
				assertPublishSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot publish to topic in another cache
				assertPublishFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot publish to a topic with a different prefix
				assertPublishFailure(tc, "w00t", sharedContext.CacheName)
				// make sure cannot subscribe to topic
				assertSubscribeFailure(tc, topicName, sharedContext.CacheName)
				// make sure cannot subscribe to prefixed topic
				assertSubscribeFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot subscribe to topic in another cache
				assertSubscribeFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot subscribe to a topic with a different prefix
				assertSubscribeFailure(tc, "w00t", sharedContext.CacheName)
			})

			It(`Generates disposable token TopicNamePrefixPublishOnly for all caches, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := TopicNamePrefixPublishOnly(AllCaches{}, topicName)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				// make sure can publish to topic
				assertPublishSuccess(tc, topicName, sharedContext.CacheName)
				// make sure can publish to prefixed topic
				assertPublishSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure can publish to topic in another cache
				assertPublishSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot publish to a topic with a different prefix
				assertPublishFailure(tc, "w00t", sharedContext.CacheName)
				// make sure cannot subscribe to topic
				assertSubscribeFailure(tc, topicName, sharedContext.CacheName)
				// make sure cannot subscribe to prefixed topic
				assertSubscribeFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot subscribe to topic in another cache
				assertSubscribeFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot subscribe to a topic with a different prefix
				assertSubscribeFailure(tc, "w00t", sharedContext.CacheName)
			})
		})

		Describe("TopicNamePrefixSubscribeOnly tokens", func() {
			It(`Generates disposable token TopicNamePrefixSubscribeOnly for a specific cache, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := TopicNamePrefixSubscribeOnly(CacheName{Name: sharedContext.CacheName}, topicName)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				// make sure cannot publish to topic
				assertPublishFailure(tc, topicName, sharedContext.CacheName)
				// make sure cannot publish to prefixed topic
				assertPublishFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot publish to topic in another cache
				assertPublishFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot publish to a topic with a different prefix
				assertPublishFailure(tc, "w00t", sharedContext.CacheName)
				// make sure can subscribe to topic
				assertSubscribeSuccess(tc, topicName, sharedContext.CacheName)
				// make sure can subscribe to prefixed topic
				assertSubscribeSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot subscribe to topic in another cache
				assertSubscribeFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot subscribe to a topic with a different prefix
				assertSubscribeFailure(tc, "w00t", sharedContext.CacheName)
			})

			It(`Generates disposable token TopicNamePrefixSubscribeOnly for all caches, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := TopicNamePrefixSubscribeOnly(AllCaches{}, topicName)
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				// make sure cannot publish to topic
				assertPublishFailure(tc, topicName, sharedContext.CacheName)
				// make sure cannot publish to prefixed topic
				assertPublishFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot publish to topic in another cache
				assertPublishFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot publish to a topic with a different prefix
				assertPublishFailure(tc, "w00t", sharedContext.CacheName)
				// make sure can subscribe to topic
				assertSubscribeSuccess(tc, topicName, sharedContext.CacheName)
				// make sure can subscribe to prefixed topic
				assertSubscribeSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure can subscribe to topic in another cache
				assertSubscribeSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot subscribe to a topic with a different prefix
				assertSubscribeFailure(tc, "w00t", sharedContext.CacheName)
			})
		})

		Describe("TopicPublishOnly tokens", func() {
			It(`Generates disposable token TopicPublishOnly for a specific cache, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := TopicPublishOnly(CacheName{Name: sharedContext.CacheName}, TopicName{Name: topicName})
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				// make sure can publish to topic
				assertPublishSuccess(tc, topicName, sharedContext.CacheName)
				// make sure cannot publish to prefixed topic
				assertPublishFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot publish to topic in another cache
				assertPublishFailure(tc, topicName, sharedContext.DefaultCacheName)
				// make sure cannot publish to a topic with a different prefix
				assertPublishFailure(tc, "w00t", sharedContext.CacheName)
				// make sure cannot subscribe to topic
				assertSubscribeFailure(tc, topicName, sharedContext.CacheName)
				// make sure cannot subscribe to prefixed topic
				assertSubscribeFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot subscribe to topic in another cache
				assertSubscribeFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot subscribe to a topic with a different prefix
				assertSubscribeFailure(tc, "w00t", sharedContext.CacheName)
			})

			It(`Generates disposable token TopicPublishOnly for all caches and all topics, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := TopicPublishOnly(AllCaches{}, AllTopics{})
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				// make sure can publish to topic
				assertPublishSuccess(tc, topicName, sharedContext.CacheName)
				// make sure can publish to prefixed topic
				assertPublishSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure can publish to topic in another cache
				assertPublishSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure can publish to a topic with a different prefix
				assertPublishSuccess(tc, "w00t", sharedContext.CacheName)
				// make sure cannot subscribe to topic
				assertSubscribeFailure(tc, topicName, sharedContext.CacheName)
				// make sure cannot subscribe to prefixed topic
				assertSubscribeFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot subscribe to topic in another cache
				assertSubscribeFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot subscribe to a topic with a different prefix
				assertSubscribeFailure(tc, "w00t", sharedContext.CacheName)
			})
		})

		Describe("TopicSubscribeOnly tokens", func() {
			It(`Generates disposable token TopicSubscribeOnly for a specific cache, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := TopicSubscribeOnly(CacheName{Name: sharedContext.CacheName}, TopicName{Name: topicName})
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				// make sure cannot publish to topic
				assertPublishFailure(tc, topicName, sharedContext.CacheName)
				// make sure cannot publish to prefixed topic
				assertPublishFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot publish to topic in another cache
				assertPublishFailure(tc, topicName, sharedContext.DefaultCacheName)
				// make sure cannot publish to a topic with a different prefix
				assertPublishFailure(tc, "w00t", sharedContext.CacheName)
				// make sure can subscribe to topic
				assertSubscribeSuccess(tc, topicName, sharedContext.CacheName)
				// make sure cannot subscribe to prefixed topic
				assertSubscribeFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot subscribe to topic in another cache
				assertSubscribeFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot subscribe to a topic with a different prefix
				assertSubscribeFailure(tc, "w00t", sharedContext.CacheName)
			})

			It(`Generates disposable token TopicSubscribeOnly for all caches and all topics, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := TopicSubscribeOnly(AllCaches{}, AllTopics{})
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				// make sure cannot publish to topic
				assertPublishFailure(tc, topicName, sharedContext.CacheName)
				// make sure cannot publish to prefixed topic
				assertPublishFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot publish to topic in another cache
				assertPublishFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot publish to a topic with a different prefix
				assertPublishFailure(tc, "w00t", sharedContext.CacheName)
				// make sure can subscribe to topic
				assertSubscribeSuccess(tc, topicName, sharedContext.CacheName)
				// make sure can subscribe to prefixed topic
				assertSubscribeSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure can subscribe to topic in another cache
				assertSubscribeSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure can subscribe to a topic with a different prefix
				assertSubscribeSuccess(tc, "w00t", sharedContext.CacheName)
			})
		})

		Describe("TopicPublishSubscribe tokens", func() {
			It(`Generates disposable token TopicPublishSubscribe for a specific cache, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := TopicPublishSubscribe(CacheName{Name: sharedContext.CacheName}, TopicName{Name: topicName})
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				// make sure can publish to topic
				assertPublishSuccess(tc, topicName, sharedContext.CacheName)
				// make sure cannot publish to prefixed topic
				assertPublishFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot publish to topic in another cache
				assertPublishFailure(tc, topicName, sharedContext.DefaultCacheName)
				// make sure cannot publish to a topic with a different prefix
				assertPublishFailure(tc, "w00t", sharedContext.CacheName)
				// make sure can subscribe to topic
				assertSubscribeSuccess(tc, topicName, sharedContext.CacheName)
				// make sure cannot subscribe to prefixed topic
				assertSubscribeFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure cannot subscribe to topic in another cache
				assertSubscribeFailure(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure cannot subscribe to a topic with a different prefix
				assertSubscribeFailure(tc, "w00t", sharedContext.CacheName)
			})

			It(`Generates disposable token TopicPublishSubscribe for all caches and all topics, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := TopicPublishSubscribe(AllCaches{}, AllTopics{})
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				// make sure can publish to topic
				assertPublishSuccess(tc, topicName, sharedContext.CacheName)
				// make sure can publish to prefixed topic
				assertPublishSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure can publish to topic in another cache
				assertPublishSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure can publish to a topic with a different prefix
				assertPublishSuccess(tc, "w00t", sharedContext.CacheName)
				// make sure can subscribe to topic
				assertSubscribeSuccess(tc, topicName, sharedContext.CacheName)
				// make sure can subscribe to prefixed topic
				assertSubscribeSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.CacheName)
				// make sure can subscribe to topic in another cache
				assertSubscribeSuccess(tc, fmt.Sprintf("%sw00t", topicName), sharedContext.DefaultCacheName)
				// make sure can subscribe to a topic with a different prefix
				assertSubscribeSuccess(tc, "w00t", sharedContext.CacheName)
			})
		})

		Describe("CacheReadWrite tokens", func() {
			It(`Generates disposable token CacheReadWrite for a specific cache, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := CacheReadWrite(CacheName{Name: sharedContext.CacheName})
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we can read the key
				assertGetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we cannot read the key for another cache
				assertGetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can read a prefixed key
				assertGetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot read a prefixed key for another cache
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write the key
				assertSetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can write a prefixed key
				assertSetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot write a prefixed key in another cache
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write a key with another prefix
				assertSetSuccess(cc, String("w00t"), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})

			It(`Generates disposable token TopicPublishSubscribe for all caches, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := CacheReadWrite(AllCaches{})
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we can read the key
				assertGetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we can read the key for another cache
				assertGetSuccess(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can read a prefixed key
				assertGetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we can read a prefixed key for another cache
				assertGetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write the key
				assertSetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we can write the key in another cache
				assertSetSuccess(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can write a prefixed key
				assertSetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we can write a prefixed key in another cache
				assertSetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write a key with another prefix
				assertSetSuccess(cc, String("w00t"), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})
		})

		Describe("CacheReadOnly tokens", func() {
			It(`Generates disposable token CacheReadOnly for a specific cache, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := CacheReadOnly(CacheName{Name: sharedContext.CacheName})
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we can read the key
				assertGetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we cannot read the key for another cache
				assertGetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can read a prefixed key
				assertGetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot read a prefixed key for another cache
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot write a prefixed key in another cache
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we cannot write a key with another prefix
				assertSetFailure(cc, String("w00t"), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})

			It(`Generates disposable token CacheReadOnly for all caches, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := CacheReadOnly(AllCaches{})
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we can read the key
				assertGetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we can read the key for another cache
				assertGetSuccess(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can read a prefixed key
				assertGetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we can read a prefixed key for another cache
				assertGetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we cannot write the key
				assertSetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot write a prefixed key
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot write a prefixed key in another cache
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we cannot write a key with another prefix
				assertSetFailure(cc, String("w00t"), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})
		})

		Describe("CacheWriteOnly tokens", func() {
			It(`Generates disposable token CacheWriteOnly for a specific cache, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := CacheWriteOnly(CacheName{Name: sharedContext.CacheName})
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot read the key for another cache
				assertGetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot read a prefixed key
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot read a prefixed key for another cache
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write the key
				assertSetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we cannot write the key in another cache
				assertSetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can write a prefixed key
				assertSetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot write a prefixed key in another cache
				assertSetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write a key with another prefix
				assertSetSuccess(cc, String("w00t"), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})

			It(`Generates disposable token CacheWriteOnly for all caches, and validates its permissions`, func() {
				topicName := uuid.NewString()
				key := String(topicName)
				scope := CacheWriteOnly(AllCaches{})
				resp := generateDisposableTokenSuccess(sharedContext, scope)
				credProvider := credProviderFromDisposableToken(resp)

				// assert cache permissions
				cc := newCacheClient(sharedContext, credProvider)
				defer cc.Close()
				// make sure that we cannot read the key
				assertGetFailure(cc, key, sharedContext.CacheName)
				// make sure that we cannot read the key for another cache
				assertGetFailure(cc, key, sharedContext.DefaultCacheName)
				// make sure that we cannot read a prefixed key
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we cannot read a prefixed key for another cache
				assertGetFailure(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write the key
				assertSetSuccess(cc, key, sharedContext.CacheName)
				// make sure that we can write the key in another cache
				assertSetSuccess(cc, key, sharedContext.DefaultCacheName)
				// make sure that we can write a prefixed key
				assertSetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.CacheName)
				// make sure that we can write a prefixed key in another cache
				assertSetSuccess(cc, String(fmt.Sprintf("%sw00t", key)), sharedContext.DefaultCacheName)
				// make sure that we can write a key with another prefix
				assertSetSuccess(cc, String("w00t"), sharedContext.CacheName)

				// asserting topic permissions
				tc := newTopicClient(sharedContext, credProvider)
				defer tc.Close()
				assertPublishFailure(tc, uuid.NewString(), sharedContext.CacheName)
				assertSubscribeFailure(tc, uuid.NewString(), sharedContext.CacheName)
			})
		})
	})

	Describe("Generate api keys", func() {
		var sessionTokenClient AuthClient
		var authTestCache1 string
		var authTestCache2 string

		BeforeEach(func() {
			sessionCredsProvider, err := auth.NewEnvMomentoTokenProvider("TEST_SESSION_TOKEN")
			if err != nil {
				Fail(fmt.Sprintf("Failed to create session token credential provider: %v", err))
			}

			// session tokens don't include cache/control endpoints so we steal them from
			// the auth-token-based credential provider and override them here
			sessionCredsProvider, err = sessionCredsProvider.WithEndpoints(auth.Endpoints{
				ControlEndpoint: sharedContext.CredentialProvider.GetControlEndpoint(),
				CacheEndpoint:   sharedContext.CredentialProvider.GetCacheEndpoint(),
				TokenEndpoint:   sharedContext.CredentialProvider.GetTokenEndpoint(),
				StorageEndpoint: sharedContext.CredentialProvider.GetStorageEndpoint(),
			})
			if err != nil {
				Fail(fmt.Sprintf("Failed to override endpionts in session token credential provider: %v", err))
			}

			sessionTokenClient, err = NewAuthClient(sharedContext.AuthConfiguration, sessionCredsProvider)
			if err != nil {
				Fail(fmt.Sprintf("Failed to create session token auth client: %v", err))
			}

			resp, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
				ExpiresIn: utils.ExpiresInDays(1),
				Scope:     internal.InternalSuperUserPermissions{},
			})
			if err != nil {
				Fail(fmt.Sprintf("Failed to create superuser api key for creating auth testing cache client: %v", err))
			}
			successResponse := resp.(*auth_responses.GenerateApiKeySuccess)

			authTestingCacheClient := newCacheClient(sharedContext, credProviderFromApiKey(successResponse.ApiKey, successResponse.Endpoint))
			authTestCache1 = fmt.Sprintf("golang-auth-%s", uuid.NewString())
			authTestCache2 = fmt.Sprintf("golang-auth-%s", uuid.NewString())

			_, err = authTestingCacheClient.CreateCache(context.Background(), &CreateCacheRequest{
				CacheName: authTestCache1,
			})
			if err != nil {
				Fail(fmt.Sprintf("Failed to create cache 1 for auth client testing: %v", err))
			}

			_, err = authTestingCacheClient.CreateCache(context.Background(), &CreateCacheRequest{
				CacheName: authTestCache2,
			})
			if err != nil {
				Fail(fmt.Sprintf("Failed to create cache 2 for auth client testing: %v", err))
			}

			DeferCleanup(func() {
				_, err = authTestingCacheClient.DeleteCache(context.Background(), &DeleteCacheRequest{
					CacheName: authTestCache1,
				})
				_, err = authTestingCacheClient.DeleteCache(context.Background(), &DeleteCacheRequest{
					CacheName: authTestCache2,
				})
				authTestingCacheClient.Close()
				sessionTokenClient.Close()
			})
		})

		It("should successfully generate api key that expires", func() {
			secondsSinceEpoch := time.Now().Unix()

			resp, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
				ExpiresIn: utils.ExpiresInDays(1),
				Scope:     internal.InternalSuperUserPermissions{},
			})
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateApiKeySuccess{}))

			successResponse := resp.(*auth_responses.GenerateApiKeySuccess)

			oneHourInSeconds := 60 * 60
			oneDayInSeconds := 24 * oneHourInSeconds
			expiresIn := secondsSinceEpoch + int64(oneDayInSeconds)

			Expect(successResponse.ExpiresAt.DoesExpire()).To(BeTrue())
			Expect(successResponse.ExpiresAt.Epoch()).To(BeNumerically("~", expiresIn-int64(oneHourInSeconds), expiresIn+int64(oneHourInSeconds)))
		})

		It("should successfully generate api key that does not expire", func() {
			resp, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
				ExpiresIn: utils.ExpiresInNever(),
				Scope:     internal.InternalSuperUserPermissions{},
			})
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateApiKeySuccess{}))

			successResponse := resp.(*auth_responses.GenerateApiKeySuccess)
			Expect(successResponse.ExpiresAt.DoesExpire()).To(BeFalse())
		})

		It("should fail to generate an api key with invalid expiry", func() {
			_, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
				ExpiresIn: utils.ExpiresInSeconds(-100),
				Scope:     internal.InternalSuperUserPermissions{},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.InvalidArgumentError))
		})

		It("should fail to generate an api key with empty permission list", func() {
			_, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
				ExpiresIn: utils.ExpiresInSeconds(10),
				Scope:     Permissions{Permissions: []Permission{}},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.InvalidArgumentError))
		})

		It("cannot create token with duplicate/conflicting cache permissions - all caches", func() {
			_, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
				ExpiresIn: utils.ExpiresInSeconds(10),
				Scope: Permissions{
					Permissions: []Permission{
						CachePermission{Cache: AllCaches{}, Role: ReadOnly},
						CachePermission{Cache: AllCaches{}, Role: ReadWrite},
					},
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.InvalidArgumentError))
		})

		It("cannot create token with duplicate/conflicting cache permissions - cache name", func() {
			_, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
				ExpiresIn: utils.ExpiresInSeconds(10),
				Scope: Permissions{
					Permissions: []Permission{
						CachePermission{Cache: CacheName{Name: "cache-name"}, Role: ReadOnly},
						CachePermission{Cache: CacheName{Name: "cache-name"}, Role: ReadWrite},
					},
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.InvalidArgumentError))
		})

		It("cannot create token with duplicate/conflicting topic permissions - cache + topic name", func() {
			_, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
				ExpiresIn: utils.ExpiresInSeconds(10),
				Scope: Permissions{
					Permissions: []Permission{
						TopicPermission{Cache: CacheName{Name: "cache-name"}, Topic: TopicName{Name: "topic-name"}, Role: SubscribeOnly},
						TopicPermission{Cache: CacheName{Name: "cache-name"}, Topic: TopicName{Name: "topic-name"}, Role: PublishSubscribe},
					},
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.InvalidArgumentError))
		})

		Describe("Superuser scope", func() {
			It("expired token cannot create cache", func() {
				resp, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInSeconds(1),
					Scope:     internal.InternalSuperUserPermissions{},
				})
				Expect(err).To(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateApiKeySuccess{}))

				// Wait for token to expire
				time.Sleep(3 * time.Second)

				success := resp.(*auth_responses.GenerateApiKeySuccess)
				cacheClient := newCacheClient(sharedContext, credProviderFromApiKey(success.ApiKey, success.Endpoint))
				defer cacheClient.Close()

				_, err = cacheClient.CreateCache(sharedContext.Ctx, &CreateCacheRequest{
					CacheName: "cache-should-fail-to-create",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.AuthenticationError))
			})

			It("cannot generate superuser token from a superuser token", func() {
				resp, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInMinutes(5),
					Scope:     internal.InternalSuperUserPermissions{},
				})
				Expect(err).To(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateApiKeySuccess{}))

				success := resp.(*auth_responses.GenerateApiKeySuccess)
				superuserAuthClient := authClientFromApiKey(sharedContext, success.ApiKey, success.Endpoint)

				_, err = superuserAuthClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInSeconds(1),
					Scope:     internal.InternalSuperUserPermissions{},
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))
			})

			It("can generate AllDataReadWrite token from superuser token", func() {
				resp, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInMinutes(5),
					Scope:     internal.InternalSuperUserPermissions{},
				})
				Expect(err).To(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateApiKeySuccess{}))

				success := resp.(*auth_responses.GenerateApiKeySuccess)
				superuserAuthClient := authClientFromApiKey(sharedContext, success.ApiKey, success.Endpoint)

				resp, err = superuserAuthClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInSeconds(1),
					Scope:     AllDataReadWrite,
				})
				Expect(err).To(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateApiKeySuccess{}))
			})

		})

		Describe("AllDataReadWrite scope", func() {
			It("can complete only data plane requests", func() {
				resp, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInMinutes(10),
					Scope:     AllDataReadWrite,
				})
				Expect(err).To(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateApiKeySuccess{}))

				success := resp.(*auth_responses.GenerateApiKeySuccess)
				cacheClient := newCacheClient(sharedContext, credProviderFromApiKey(success.ApiKey, success.Endpoint))

				// Cannot create a cache
				_, err = cacheClient.CreateCache(sharedContext.Ctx, &CreateCacheRequest{
					CacheName: "cache-should-fail-to-create",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				// Cannot delete a cache
				_, err = cacheClient.DeleteCache(sharedContext.Ctx, &DeleteCacheRequest{
					CacheName: "cache-should-fail-to-delete",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				// Can set values in an existing cache
				setResp, err := cacheClient.Set(sharedContext.Ctx, &SetRequest{
					CacheName: authTestCache1,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(err).To(BeNil())
				Expect(setResp).To(BeAssignableToTypeOf(&responses.SetSuccess{}))

				// Can get values in an existing cache
				getResp, err := cacheClient.Get(sharedContext.Ctx, &GetRequest{
					CacheName: authTestCache1,
					Key:       String("key"),
				})
				Expect(err).To(BeNil())
				Expect(getResp).To(BeAssignableToTypeOf(&responses.GetHit{}))
			})

			It("cannot generate superuser token from an AllDataReadWrite token", func() {
				resp, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInMinutes(5),
					Scope:     AllDataReadWrite,
				})
				Expect(err).To(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateApiKeySuccess{}))

				success := resp.(*auth_responses.GenerateApiKeySuccess)
				superuserAuthClient := authClientFromApiKey(sharedContext, success.ApiKey, success.Endpoint)

				_, err = superuserAuthClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInSeconds(1),
					Scope:     internal.InternalSuperUserPermissions{},
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))
			})

			It("cannot generate AllDataReadWrite token from an AllDataReadWrite token", func() {
				resp, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInMinutes(5),
					Scope:     AllDataReadWrite,
				})
				Expect(err).To(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateApiKeySuccess{}))

				success := resp.(*auth_responses.GenerateApiKeySuccess)
				superuserAuthClient := authClientFromApiKey(sharedContext, success.ApiKey, success.Endpoint)

				_, err = superuserAuthClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInSeconds(1),
					Scope:     AllDataReadWrite,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))
			})

		})

		Describe("Fine-grained access scope", func() {
			It("CachePermission ReadOnly", func() {
				resp, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInMinutes(5),
					Scope: Permissions{Permissions: []Permission{
						CachePermission{Cache: AllCaches{}, Role: ReadOnly},
					}},
				})
				Expect(err).To(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateApiKeySuccess{}))

				success := resp.(*auth_responses.GenerateApiKeySuccess)
				cacheClient := newCacheClient(sharedContext, credProviderFromApiKey(success.ApiKey, success.Endpoint))
				topicClient := newTopicClient(sharedContext, credProviderFromApiKey(success.ApiKey, success.Endpoint))
				defer cacheClient.Close()
				defer topicClient.Close()

				// 1. Sets to both caches should fail
				_, err = cacheClient.Set(sharedContext.Ctx, &SetRequest{
					CacheName: authTestCache1,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				_, err = cacheClient.Set(sharedContext.Ctx, &SetRequest{
					CacheName: authTestCache2,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				// 2. Gets from both caches should succeed
				getResp, err := cacheClient.Get(sharedContext.Ctx, &GetRequest{
					CacheName: authTestCache1,
					Key:       String("key"),
				})
				Expect(err).To(BeNil())
				Expect(getResp).To(BeAssignableToTypeOf(&responses.GetMiss{}))

				getResp, err = cacheClient.Get(sharedContext.Ctx, &GetRequest{
					CacheName: authTestCache2,
					Key:       String("key"),
				})
				Expect(err).To(BeNil())
				Expect(getResp).To(BeAssignableToTypeOf(&responses.GetMiss{}))

				// 3. Publishes should fail
				_, err = topicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
					CacheName: authTestCache1,
					TopicName: "topic",
					Value:     String("Hello Mo"),
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				_, err = topicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
					CacheName: authTestCache2,
					TopicName: "topic",
					Value:     String("Hello Mo"),
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				// 4. Subscribes should fail
				_, err = topicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
					CacheName: authTestCache1,
					TopicName: "topic",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				_, err = topicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
					CacheName: authTestCache2,
					TopicName: "topic",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))
			})

			It("TopicPermission SubscribeOnly", func() {
				resp, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInMinutes(5),
					Scope: Permissions{Permissions: []Permission{
						TopicPermission{Cache: AllCaches{}, Topic: AllTopics{}, Role: SubscribeOnly},
					}},
				})
				Expect(err).To(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateApiKeySuccess{}))

				success := resp.(*auth_responses.GenerateApiKeySuccess)
				cacheClient := newCacheClient(sharedContext, credProviderFromApiKey(success.ApiKey, success.Endpoint))
				topicClient := newTopicClient(sharedContext, credProviderFromApiKey(success.ApiKey, success.Endpoint))
				defer cacheClient.Close()
				defer topicClient.Close()

				// 1. Sets to both caches should fail
				_, err = cacheClient.Set(sharedContext.Ctx, &SetRequest{
					CacheName: authTestCache1,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				_, err = cacheClient.Set(sharedContext.Ctx, &SetRequest{
					CacheName: authTestCache2,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				// 2. Gets from both caches should fail
				_, err = cacheClient.Get(sharedContext.Ctx, &GetRequest{
					CacheName: authTestCache1,
					Key:       String("key"),
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				_, err = cacheClient.Get(sharedContext.Ctx, &GetRequest{
					CacheName: authTestCache2,
					Key:       String("key"),
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				// 3. Publishes should fail
				_, err = topicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
					CacheName: authTestCache1,
					TopicName: "topic",
					Value:     String("Hello Mo"),
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				_, err = topicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
					CacheName: authTestCache2,
					TopicName: "topic",
					Value:     String("Hello Mo"),
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				// 4. Subscribes should succeed
				sub, err := topicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
					CacheName: authTestCache1,
					TopicName: "topic",
				})
				Expect(err).To(BeNil())
				Expect(sub).NotTo(BeNil())

				sub, err = topicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
					CacheName: authTestCache2,
					TopicName: "topic",
				})
				Expect(err).To(BeNil())
				Expect(sub).NotTo(BeNil())
			})

			It("Mixed Cache and Topic ReadWrite permissions", func() {
				resp, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInMinutes(5),
					Scope: Permissions{Permissions: []Permission{
						CachePermission{Cache: CacheName{Name: authTestCache1}, Role: ReadWrite},
						TopicPermission{Cache: CacheName{Name: authTestCache2}, Topic: AllTopics{}, Role: PublishSubscribe},
					}},
				})
				Expect(err).To(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateApiKeySuccess{}))

				success := resp.(*auth_responses.GenerateApiKeySuccess)
				cacheClient := newCacheClient(sharedContext, credProviderFromApiKey(success.ApiKey, success.Endpoint))
				topicClient := newTopicClient(sharedContext, credProviderFromApiKey(success.ApiKey, success.Endpoint))
				defer cacheClient.Close()
				defer topicClient.Close()

				// 1. Cache get/set on authTestCache1 should succeed
				setResp, err := cacheClient.Set(sharedContext.Ctx, &SetRequest{
					CacheName: authTestCache1,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(err).To(BeNil())
				Expect(setResp).To(BeAssignableToTypeOf(&responses.SetSuccess{}))

				getResp, err := cacheClient.Get(sharedContext.Ctx, &GetRequest{
					CacheName: authTestCache1,
					Key:       String("key"),
				})
				Expect(err).To(BeNil())
				Expect(getResp).To(BeAssignableToTypeOf(&responses.GetHit{}))

				// 2. Cache get/set on authTestCache2 should fail
				_, err = cacheClient.Set(sharedContext.Ctx, &SetRequest{
					CacheName: authTestCache2,
					Key:       String("key"),
					Value:     String("value"),
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				_, err = cacheClient.Get(sharedContext.Ctx, &GetRequest{
					CacheName: authTestCache2,
					Key:       String("key"),
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				// 3. Publish/subscribe on authTestCache1 should fail
				_, err = topicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
					CacheName: authTestCache1,
					TopicName: "topic",
					Value:     String("Hello Mo"),
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				_, err = topicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
					CacheName: authTestCache1,
					TopicName: "topic",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.PermissionError))

				// 4. Publish/subscribe on authTestCache2 should succeed
				pubResp, err := topicClient.Publish(sharedContext.Ctx, &TopicPublishRequest{
					CacheName: authTestCache2,
					TopicName: "topic",
					Value:     String("Hello Mo"),
				})
				Expect(err).To(BeNil())
				Expect(pubResp).To(BeAssignableToTypeOf(&responses.TopicPublishSuccess{}))

				sub, err := topicClient.Subscribe(sharedContext.Ctx, &TopicSubscribeRequest{
					CacheName: authTestCache2,
					TopicName: "topic",
				})
				Expect(err).To(BeNil())
				Expect(sub).NotTo(BeNil())
			})

		})

		Describe("Refresh api keys", func() {
			It("should successfully refresh api key with unexpired refresh token", func() {
				resp, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInSeconds(30),
					Scope:     internal.InternalSuperUserPermissions{},
				})
				Expect(err).To(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateApiKeySuccess{}))

				success := resp.(*auth_responses.GenerateApiKeySuccess)
				testingAuthClient := authClientFromApiKey(sharedContext, success.ApiKey, success.Endpoint)

				// we need to sleep for a bit here so that the timestamp on the refreshed token
				// will be different than the one on the original token
				delaySecondsBeforeRefresh := 2 * time.Second
				time.Sleep(delaySecondsBeforeRefresh)

				refreshResp, err := testingAuthClient.RefreshApiKey(sharedContext.Ctx, &RefreshApiKeyRequest{
					RefreshToken: success.RefreshToken,
				})
				Expect(err).To(BeNil())
				Expect(refreshResp).To(BeAssignableToTypeOf(&auth_responses.RefreshApiKeySuccess{}))

				refreshSuccess := refreshResp.(*auth_responses.RefreshApiKeySuccess)

				expiresAtDelta := refreshSuccess.ExpiresAt.Epoch() - success.ExpiresAt.Epoch()
				Expect(expiresAtDelta).To(BeNumerically("~", delaySecondsBeforeRefresh.Seconds(), delaySecondsBeforeRefresh.Seconds()+10))
			})

			It("should fail to refresh api key with expired refresh token", func() {
				resp, err := sessionTokenClient.GenerateApiKey(sharedContext.Ctx, &GenerateApiKeyRequest{
					ExpiresIn: utils.ExpiresInSeconds(1),
					Scope:     internal.InternalSuperUserPermissions{},
				})
				Expect(err).To(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(&auth_responses.GenerateApiKeySuccess{}))

				success := resp.(*auth_responses.GenerateApiKeySuccess)
				testingAuthClient := authClientFromApiKey(sharedContext, success.ApiKey, success.Endpoint)

				// wait for api key to expire
				time.Sleep(3 * time.Second)

				_, err = testingAuthClient.RefreshApiKey(sharedContext.Ctx, &RefreshApiKeyRequest{
					RefreshToken: success.RefreshToken,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.(MomentoError).Code()).To(Equal(momentoerrors.InvalidArgumentError))
			})
		})
	})

})
