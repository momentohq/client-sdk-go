package momento_test

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	responses "github.com/momentohq/client-sdk-go/responses/auth"
	"github.com/momentohq/client-sdk-go/utils"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
)

func credProviderFromDisposableToken(resp responses.GenerateDisposableTokenResponse) auth.CredentialProvider {
	success := resp.(*responses.GenerateDisposableTokenSuccess)
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

func generateDisposableTokenSuccess(ctx SharedContext, scope DisposableTokenScope) responses.GenerateDisposableTokenResponse {
	expiresIn := utils.ExpiresInMinutes(5)
	resp, err := ctx.AuthClient.GenerateDisposableToken(ctx.Ctx, &GenerateDisposableTokenRequest{
		ExpiresIn: expiresIn,
		Scope:     scope,
	})
	if err != nil {
		panic(err)
	}
	Expect(resp).To(BeAssignableToTypeOf(&responses.GenerateDisposableTokenSuccess{}))
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

	Describe("Disposable tokens", func() {
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

	Describe("Api keys", func() {})

})
