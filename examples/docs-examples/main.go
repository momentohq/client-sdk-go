package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
	auth_resp "github.com/momentohq/client-sdk-go/responses/auth"
	"github.com/momentohq/client-sdk-go/utils"
)

var (
	ctx                context.Context
	client             momento.CacheClient
	database           map[string]string
	cacheName          string
	leaderboardClient  momento.PreviewLeaderboardClient
	leaderboard        momento.Leaderboard
	credentialProvider auth.CredentialProvider
	topicClient        momento.TopicClient
	authClient         momento.AuthClient
	err                error
	hashValue          string
)

func RetrieveApiKeyFromYourSecretsManager() string {
	return "your-api-key"
}

func example_API_CredentialProviderFromString() {
	apiKey := RetrieveApiKeyFromYourSecretsManager()
	credentialProvider, err = auth.NewStringMomentoTokenProvider(apiKey)
	if err != nil {
		fmt.Println("Error parsing API key:", err)
	}
}

func example_API_CredentialProviderFromEnvVar() {
	credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}
}

func example_API_CredentialProviderFromEnvVarV2() {
	credentialProvider, err = auth.NewEnvMomentoV2TokenProvider()
	if err != nil {
		panic(err)
	}
}

func example_API_CredentialProviderFromApiKeyV2() {
	apiKey := RetrieveApiKeyFromYourSecretsManager()
	endpoint := "https://api.cache.cell-4-us-west-2-1.prod.a.momentohq.com"
	props := auth.ApiKeyV2Props{ApiKey: apiKey, Endpoint: endpoint}
	credentialProvider, err = auth.NewApiKeyV2TokenProvider(props)
	if err != nil {
		panic(err)
	}
}

func example_API_InstantiateCacheClient() {
	credentialProvider, err = auth.NewEnvMomentoV2TokenProvider()
	if err != nil {
		panic(err)
	}
	defaultTtl := 60 * time.Second
	eagerConnectTimeout := 30 * time.Second

	client, err = momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatest(),
		credentialProvider,
		defaultTtl,
		eagerConnectTimeout,
	)
	if err != nil {
		panic(err)
	}

	client.Ping(ctx)
}

func example_API_InstantiateCacheClientWithReadConcern() {
	credentialProvider, err := auth.NewEnvMomentoV2TokenProvider()
	if err != nil {
		panic(err)
	}
	defaultTtl := time.Duration(9999)
	eagerConnectTimeout := 30 * time.Second

	client, err := momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatest().WithReadConcern(config.CONSISTENT),
		credentialProvider,
		defaultTtl,
		eagerConnectTimeout,
	)
	if err != nil {
		panic(err)
	}

	client.Ping(ctx)
}

func example_API_ListCaches() {
	resp, err := client.ListCaches(ctx, &momento.ListCachesRequest{})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *responses.ListCachesSuccess:
		log.Printf("Found caches %+v\n", r.Caches())
	}
}

func example_API_CreateCache() {
	_, err := client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}
}

func example_API_DeleteCache() {
	_, err := client.DeleteCache(ctx, &momento.DeleteCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}
}

func example_API_Get() {
	key := uuid.NewString()
	resp, err := client.Get(ctx, &momento.GetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
	})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *responses.GetHit:
		log.Printf("Lookup resulted in cache HIT. value=%s\n", r.ValueString())
	case *responses.GetMiss:
		log.Printf("Look up did not find a value key=%s\n", key)
	}
}

func example_API_Set() {
	key := uuid.NewString()
	value := uuid.NewString()
	log.Printf("Setting key: %s, value: %s\n", key, value)
	_, err := client.Set(ctx, &momento.SetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
		Value:     momento.String(value),
		Ttl:       time.Duration(9999),
	})
	if err != nil {
		var momentoErr momento.MomentoError
		if errors.As(err, &momentoErr) {
			if momentoErr.Code() != momento.TimeoutError {
				// this would represent a client-side timeout, and you could fall back to your original data source
			} else {
				panic(err)
			}
		}
	}
}

func example_API_Delete() {
	key := uuid.NewString()
	_, err := client.Delete(ctx, &momento.DeleteRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
	})
	if err != nil {
		panic(err)
	}
}

func example_API_InstantiateTopicClient() {
	credProvider, err := auth.NewEnvMomentoV2TokenProvider()
	if err != nil {
		panic(err)
	}

	topicClient, err = momento.NewTopicClient(
		config.TopicsDefault(),
		credProvider,
	)
	if err != nil {
		panic(err)
	}
}

func example_API_TopicPublish() {
	_, err := topicClient.Publish(ctx, &momento.TopicPublishRequest{
		CacheName: cacheName,
		TopicName: "test-topic",
		Value:     momento.String("test-message"),
	})
	if err != nil {
		panic(err)
	}
}

func example_API_TopicSubscribe() {
	// Instantiate subscriber
	sub, subErr := topicClient.Subscribe(ctx, &momento.TopicSubscribeRequest{
		CacheName: cacheName,
		TopicName: "test-topic",
	})
	if subErr != nil {
		panic(subErr)
	}

	time.Sleep(time.Second)
	_, pubErr := topicClient.Publish(ctx, &momento.TopicPublishRequest{
		CacheName: cacheName,
		TopicName: "test-topic",
		Value:     momento.String("test-message"),
	})
	if pubErr != nil {
		panic(pubErr)
	}
	time.Sleep(time.Second)

	// Receive only subscription items with messages
	item, err := sub.Item(ctx)
	if err != nil {
		panic(err)
	}
	switch msg := item.(type) {
	case momento.String:
		fmt.Printf("received message as string: '%v'\n", msg)
	case momento.Bytes:
		fmt.Printf("received message as bytes: '%v'\n", msg)
	}

	// Receive all subscription events (messages, discontinuities, heartbeats)
	event, err := sub.Event(ctx)
	if err != nil {
		panic(err)
	}
	switch e := event.(type) {
	case momento.TopicHeartbeat:
		fmt.Printf("received heartbeat\n")
	case momento.TopicDiscontinuity:
		fmt.Printf("received discontinuity\n")
	case momento.TopicItem:
		fmt.Printf(
			"received message with sequence number %d and publisher id %s: %v \n",
			e.GetTopicSequenceNumber(),
			e.GetPublisherId(),
			e.GetValue(),
		)
	}
}

func example_API_InstantiateAuthClient() {
	credentialProvider, err := auth.NewEnvMomentoV2TokenProvider()
	if err != nil {
		panic(err)
	}

	authClient, err = momento.NewAuthClient(config.AuthDefault(), credentialProvider)
	if err != nil {
		panic(err)
	}
}

func example_API_GenerateDisposableToken() {
	tokenId := "a token id"
	resp, err := authClient.GenerateDisposableToken(ctx, &momento.GenerateDisposableTokenRequest{
		ExpiresIn: utils.ExpiresInSeconds(10),
		Scope: momento.TopicSubscribeOnly(
			momento.CacheName{Name: "a cache"},
			momento.TopicName{Name: "a topic"},
		),
		Props: momento.DisposableTokenProps{
			TokenId: &tokenId,
		},
	})

	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *auth_resp.GenerateDisposableTokenSuccess:
		log.Printf("Successfully generated a disposable token for endpoint=%s with tokenId=%s\n", r.Endpoint, tokenId)
	}
}

func example_API_GenerateApiKey() {
	// Generate a token that allows all data plane APIs on all caches and topics.
	resp, err := authClient.GenerateApiKey(ctx, &momento.GenerateApiKeyRequest{
		ExpiresIn: utils.ExpiresInMinutes(30),
		Scope:     momento.AllDataReadWrite,
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *auth_resp.GenerateApiKeySuccess:
		log.Printf("Successfully generated an API key with AllDataReadWrite scope!\n")
		log.Printf("API key expires at: %d\n", r.ExpiresAt.Epoch())
	}

	// Generate a token that can only call read-only data plane APIs on a specific cache foo. No topic apis (publish/subscribe) are allowed.
	resp, err = authClient.GenerateApiKey(ctx, &momento.GenerateApiKeyRequest{
		ExpiresIn: utils.ExpiresInMinutes(30),
		Scope:     momento.CacheReadOnly(momento.CacheName{Name: "foo"}),
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *auth_resp.GenerateApiKeySuccess:
		log.Printf("Successfully generated an API key with read-only access to cache foo!\n")
		log.Printf("API key expires at: %d\n", r.ExpiresAt.Epoch())
	}

	// Generate a token that can call all data plane APIs on all caches. No topic apis (publish/subscribe) are allowed.
	resp, err = authClient.GenerateApiKey(ctx, &momento.GenerateApiKeyRequest{
		ExpiresIn: utils.ExpiresInMinutes(30),
		Scope:     momento.CacheReadWrite(momento.AllCaches{}),
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *auth_resp.GenerateApiKeySuccess:
		log.Printf("Successfully generated an API key with read-write access to all caches!\n")
		log.Printf("API key expires at: %d\n", r.ExpiresAt.Epoch())
	}

	// Generate a token that can call publish and subscribe on all topics within cache bar
	resp, err = authClient.GenerateApiKey(ctx, &momento.GenerateApiKeyRequest{
		ExpiresIn: utils.ExpiresInMinutes(30),
		Scope:     momento.TopicPublishSubscribe(momento.CacheName{Name: "bar"}, momento.AllTopics{}),
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *auth_resp.GenerateApiKeySuccess:
		log.Printf("Successfully generated an API key publish-subscribe access to all topics within cache bar!\n")
		log.Printf("API key expires at: %d\n", r.ExpiresAt.Epoch())
	}

	// Generate a token that can only call subscribe on topic where_is_mo within cache mo_nuts
	resp, err = authClient.GenerateApiKey(ctx, &momento.GenerateApiKeyRequest{
		ExpiresIn: utils.ExpiresInMinutes(30),
		Scope:     momento.TopicSubscribeOnly(momento.CacheName{Name: "mo_nuts"}, momento.TopicName{Name: "where_is_mo"}),
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *auth_resp.GenerateApiKeySuccess:
		log.Printf("Successfully generated an API key with subscribe-only access to topic where_is_mo within cache mo_nuts!\n")
		log.Printf("API key expires at: %d\n", r.ExpiresAt.Epoch())
	}

	// Generate a token with multiple permissions
	cachePermission1 := momento.CachePermission{
		Cache: momento.CacheName{Name: "acorns"}, // Scopes the access to a single cache named 'acorns'
		Role:  momento.ReadWrite,                 // Managed role that grants access to read as well as write apis on caches
	}
	cachePermission2 := momento.CachePermission{
		Cache: momento.AllCaches{}, // Built-in value for access to all caches in the account
		Role:  momento.ReadOnly,    // Managed role that grants access to only read data apis on caches
	}
	topicPermission1 := momento.TopicPermission{
		Cache: momento.CacheName{Name: "walnuts"},      // Scopes the access to a single cache named 'walnuts'
		Topic: momento.TopicName{Name: "mo_favorites"}, // Scopes the access to a single topic named 'mo_favorites' within cache 'walnuts'
		Role:  momento.PublishSubscribe,                // Managed role that grants access to subscribe as well as publish apis
	}
	topicPermission2 := momento.TopicPermission{
		Cache: momento.AllCaches{},   // Built-in value for all cache(s) in the account.
		Topic: momento.AllTopics{},   // Built-in value for access to all topics in the listed cache(s).
		Role:  momento.SubscribeOnly, // Managed role that grants access to only subscribe api
	}
	permissions := []momento.Permission{
		cachePermission1, cachePermission2, topicPermission1, topicPermission2,
	}

	resp, err = authClient.GenerateApiKey(ctx, &momento.GenerateApiKeyRequest{
		ExpiresIn: utils.ExpiresInMinutes(30),
		Scope: momento.Permissions{
			Permissions: permissions,
		},
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *auth_resp.GenerateApiKeySuccess:
		log.Printf("Successfully generated an API key with multiple cache and topic permissions!\n")
		log.Printf("API key expires at: %d\n", r.ExpiresAt.Epoch())
	}
}

func example_API_RefreshApiKey() {
	resp, err := authClient.GenerateApiKey(ctx, &momento.GenerateApiKeyRequest{
		ExpiresIn: utils.ExpiresInMinutes(30),
		Scope:     momento.AllDataReadWrite,
	})
	if err != nil {
		panic(err)
	}
	generateApiKeySuccess := resp.(*auth_resp.GenerateApiKeySuccess)

	newCredProvider, err := auth.FromString(generateApiKeySuccess.ApiKey)
	if err != nil {
		panic(err)
	}

	refreshAuthClient, err := momento.NewAuthClient(config.AuthDefault(), newCredProvider)
	if err != nil {
		panic(err)
	}

	refreshResp, err := refreshAuthClient.RefreshApiKey(ctx, &momento.RefreshApiKeyRequest{
		RefreshToken: generateApiKeySuccess.RefreshToken,
	})
	if err != nil {
		panic(err)
	}
	switch r := refreshResp.(type) {
	case *auth_resp.RefreshApiKeySuccess:
		log.Printf("Successfully refreshed API key!\n")
		log.Printf("Refreshed API key expires at: %d\n", r.ExpiresAt.Epoch())
	}
}

func example_API_SetIfPresent() {
	resp, err := client.SetIfPresent(ctx, &momento.SetIfPresentRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Value:     momento.String("value"),
	})
	if err != nil {
		panic(err)
	}
	switch resp.(type) {
	case *responses.SetIfPresentStored:
		log.Printf("Successfully set key in cache\n")
	}
}

func example_API_SetIfAbsent() {
	resp, err := client.SetIfAbsent(ctx, &momento.SetIfAbsentRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Value:     momento.String("value"),
	})
	if err != nil {
		panic(err)
	}
	switch resp.(type) {
	case *responses.SetIfAbsentStored:
		log.Printf("Successfully set key in cache\n")
	}
}

func example_API_SetIfEqual() {
	resp, err := client.SetIfEqual(ctx, &momento.SetIfEqualRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Value:     momento.String("value"),
		Equal:     momento.String("current-value"),
	})
	if err != nil {
		panic(err)
	}
	switch resp.(type) {
	case *responses.SetIfEqualStored:
		log.Printf("Successfully set key in cache\n")
	}
}

func example_API_SetIfNotEqual() {
	resp, err := client.SetIfNotEqual(ctx, &momento.SetIfNotEqualRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Value:     momento.String("value"),
		NotEqual:  momento.String("current-value"),
	})
	if err != nil {
		panic(err)
	}
	switch resp.(type) {
	case *responses.SetIfNotEqualStored:
		log.Printf("Successfully set key in cache\n")
	}
}

func example_API_SetIfPresentAndNotEqual() {
	resp, err := client.SetIfPresentAndNotEqual(ctx, &momento.SetIfPresentAndNotEqualRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Value:     momento.String("value"),
		NotEqual:  momento.String("current-value"),
	})
	if err != nil {
		panic(err)
	}
	switch resp.(type) {
	case *responses.SetIfPresentAndNotEqualStored:
		log.Printf("Successfully set key in cache\n")
	}
}

func example_API_SetIfAbsentOrEqual() {
	resp, err := client.SetIfAbsentOrEqual(ctx, &momento.SetIfAbsentOrEqualRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Value:     momento.String("value"),
		Equal:     momento.String("current-value"),
	})
	if err != nil {
		panic(err)
	}
	switch resp.(type) {
	case *responses.SetIfAbsentOrEqualStored:
		log.Printf("Successfully set key in cache\n")
	}
}

func example_API_GetWithHash() {
	resp, err := client.GetWithHash(ctx, &momento.GetWithHashRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *responses.GetWithHashHit:
		log.Printf("Successfully got value %s with hash %s\n", r.ValueString(), r.HashString())
	case *responses.GetWithHashMiss:
		log.Printf("Key does not exist in cache\n")
	}
}

func example_API_SetWithHash() {
	resp, err := client.SetWithHash(ctx, &momento.SetWithHashRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Value:     momento.String("value"),
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *responses.SetWithHashStored:
		log.Printf("Successfully set key in cache, item has new hash %s\n", r.HashString())
		hashValue = r.HashString()
	case *responses.SetWithHashNotStored:
		log.Printf("Unable to set key in cache\n")
	}
}

func example_API_SetIfPresentAndHashNotEqual() {
	resp, err := client.SetIfPresentAndHashNotEqual(ctx, &momento.SetIfPresentAndHashNotEqualRequest{
		CacheName:    cacheName,
		Key:          momento.String("key"),
		Value:        momento.String("value"),
		HashNotEqual: momento.String(hashValue),
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *responses.SetIfPresentAndHashNotEqualStored:
		log.Printf("Successfully set key in cache, item has new hash %s\n", r.HashString())
	case *responses.SetIfPresentAndHashNotEqualNotStored:
		log.Printf("Unable to set key in cache\n")
	}
}

func example_API_SetIfPresentAndHashEqual() {
	resp, err := client.SetIfPresentAndHashEqual(ctx, &momento.SetIfPresentAndHashEqualRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Value:     momento.String("value"),
		HashEqual: momento.String(hashValue),
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *responses.SetIfPresentAndHashEqualStored:
		log.Printf("Successfully set key in cache, item has new hash %s\n", r.HashString())
	case *responses.SetIfPresentAndHashEqualNotStored:
		log.Printf("Unable to set key in cache\n")
	}
}

func example_API_SetIfAbsentOrHashEqual() {
	resp, err := client.SetIfAbsentOrHashEqual(ctx, &momento.SetIfAbsentOrHashEqualRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Value:     momento.String("value"),
		HashEqual: momento.String(hashValue),
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *responses.SetIfAbsentOrHashEqualStored:
		log.Printf("Successfully set key in cache, item has new hash %s\n", r.HashString())
	case *responses.SetIfAbsentOrHashEqualNotStored:
		log.Printf("Unable to set key in cache\n")
	}
}

func example_API_SetIfAbsentOrHashNotEqual() {
	resp, err := client.SetIfAbsentOrHashNotEqual(ctx, &momento.SetIfAbsentOrHashNotEqualRequest{
		CacheName:    cacheName,
		Key:          momento.String("key"),
		Value:        momento.String("value"),
		HashNotEqual: momento.String(hashValue),
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *responses.SetIfAbsentOrHashNotEqualStored:
		log.Printf("Successfully set key in cache, item has new hash %s\n", r.HashString())
	case *responses.SetIfAbsentOrHashNotEqualNotStored:
		log.Printf("Unable to set key in cache\n")
	}
}

func example_API_KeysExist() {
	keys := []momento.Value{momento.String("key1"), momento.String("key2")}
	resp, err := client.KeysExist(ctx, &momento.KeysExistRequest{
		CacheName: cacheName,
		Keys:      keys,
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *responses.KeysExistSuccess:
		log.Printf("Does each key exist in cache?\n")
		for _, exists := range r.Exists() {
			log.Printf("key exists=%v\n", exists)
		}
	}
}

func example_API_ItemGetType() {
	resp, err := client.ItemGetType(ctx, &momento.ItemGetTypeRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *responses.ItemGetTypeHit:
		log.Printf("Type of item is %v\n", r.Type())
	case *responses.ItemGetTypeMiss:
		log.Printf("Item does not exist in cache\n")
	}
}

func example_API_UpdateTtl() {
	resp, err := client.UpdateTtl(ctx, &momento.UpdateTtlRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Ttl:       time.Duration(9999),
	})
	if err != nil {
		panic(err)
	}
	switch resp.(type) {
	case *responses.UpdateTtlSet:
		log.Printf("Successfully updated TTL for key\n")
	case *responses.UpdateTtlMiss:
		log.Printf("Key does not exist in cache\n")
	}
}

func example_API_IncreaseTtl() {
	resp, err := client.IncreaseTtl(ctx, &momento.IncreaseTtlRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Ttl:       time.Duration(9999),
	})
	if err != nil {
		panic(err)
	}
	switch resp.(type) {
	case *responses.IncreaseTtlSet:
		log.Printf("Successfully increased TTL for key\n")
	case *responses.IncreaseTtlMiss:
		log.Printf("Key does not exist in cache\n")
	}
}

func example_API_DecreaseTtl() {
	resp, err := client.DecreaseTtl(ctx, &momento.DecreaseTtlRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
		Ttl:       time.Duration(9999),
	})
	if err != nil {
		panic(err)
	}
	switch resp.(type) {
	case *responses.DecreaseTtlSet:
		log.Printf("Successfully decreased TTL for key\n")
	case *responses.DecreaseTtlMiss:
		log.Printf("Key does not exist in cache\n")
	}
}

func example_API_ItemGetTtl() {
	resp, err := client.ItemGetTtl(ctx, &momento.ItemGetTtlRequest{
		CacheName: cacheName,
		Key:       momento.String("key"),
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *responses.ItemGetTtlHit:
		log.Printf("TTL for key is %d\n", r.RemainingTtl().Milliseconds())
	case *responses.ItemGetTtlMiss:
		log.Printf("Key does not exist in cache\n")
	}
}

func example_API_Increment() {
	resp, err := client.Increment(ctx, &momento.IncrementRequest{
		CacheName: cacheName,
		Field:     momento.String("key"),
		Amount:    1,
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *responses.IncrementSuccess:
		log.Printf("Incremented value is %d\n", r.Value())
	}
}

func example_API_GetBatch() {
	resp, err := client.GetBatch(ctx, &momento.GetBatchRequest{
		CacheName: cacheName,
		Keys:      []momento.Value{momento.String("key1"), momento.String("key2")},
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *responses.GetBatchSuccess:
		log.Printf("Found values %+v\n", r.ValueMap())
	}
}

func example_API_SetBatch() {
	resp, err := client.SetBatch(ctx, &momento.SetBatchRequest{
		CacheName: cacheName,
		Items: []momento.BatchSetItem{
			{
				Key:   momento.String("key1"),
				Value: momento.String("value1"),
			},
			{
				Key:   momento.String("key2"),
				Value: momento.String("value2"),
			},
		},
	})
	if err != nil {
		panic(err)
	}
	switch resp.(type) {
	case *responses.SetBatchSuccess:
		log.Printf("Successfully set keys in cache\n")
	}
}

func example_API_InstantiateLeaderboardClient() {
	credentialProvider, err := auth.NewEnvMomentoV2TokenProvider()
	if err != nil {
		panic(err)
	}

	leaderboardClient, err = momento.NewPreviewLeaderboardClient(
		config.LeaderboardDefault(),
		credentialProvider,
	)
	if err != nil {
		panic(err)
	}
}

func example_API_CreateLeaderboard() momento.Leaderboard {
	leaderboard, err = leaderboardClient.Leaderboard(ctx, &momento.LeaderboardRequest{
		CacheName:       cacheName,
		LeaderboardName: "leaderboard",
	})
	if err != nil {
		panic(err)
	}
	return leaderboard
}

func example_API_LeaderboardUpsert() {
	upsertElements := []momento.LeaderboardUpsertElement{
		{Id: 123, Score: 10.33},
		{Id: 456, Score: 3333},
		{Id: 789, Score: 5678.9},
	}
	_, err := leaderboard.Upsert(ctx, momento.LeaderboardUpsertRequest{Elements: upsertElements})
	if err != nil {
		panic(err)
	}
}

func example_API_LeaderboardFetchByScore() {
	minScore := 150.0
	maxScore := 3000.0
	offset := uint32(1)
	count := uint32(2)
	fetchOrder := momento.ASCENDING
	fetchByScoreResponse, err := leaderboard.FetchByScore(ctx, momento.LeaderboardFetchByScoreRequest{
		MinScore: &minScore,
		MaxScore: &maxScore,
		Offset:   &offset,
		Count:    &count,
		Order:    &fetchOrder,
	})
	if err != nil {
		panic(err)
	}
	switch r := fetchByScoreResponse.(type) {
	case *responses.LeaderboardFetchSuccess:
		fmt.Printf("Successfully fetched elements by score:\n")
		for _, element := range r.Values() {
			fmt.Printf("ID: %d, Score: %f, Rank: %d\n", element.Id, element.Score, element.Rank)
		}
	}
}

func example_API_LeaderboardFetchByRank() {
	fetchOrder := momento.ASCENDING
	fetchByRankResponse, err := leaderboard.FetchByRank(ctx, momento.LeaderboardFetchByRankRequest{
		StartRank: 0,
		EndRank:   100,
		Order:     &fetchOrder,
	})
	if err != nil {
		panic(err)
	}
	switch r := fetchByRankResponse.(type) {
	case *responses.LeaderboardFetchSuccess:
		fmt.Printf("Successfully fetched elements by rank:\n")
		for _, element := range r.Values() {
			fmt.Printf("ID: %d, Score: %f, Rank: %d\n", element.Id, element.Score, element.Rank)
		}
	}
}

func example_API_LeaderboardGetRank() {
	getRankResponse, err := leaderboard.GetRank(ctx, momento.LeaderboardGetRankRequest{
		Ids: []uint32{123, 456},
	})
	if err != nil {
		panic(err)
	}
	switch r := getRankResponse.(type) {
	case *responses.LeaderboardFetchSuccess:
		fmt.Printf("Successfully fetched elements by ID:\n")
		for _, element := range r.Values() {
			fmt.Printf("ID: %d, Score: %f, Rank: %d\n", element.Id, element.Score, element.Rank)
		}
	}
}

func example_API_LeaderboardLength() {
	lengthResponse, err := leaderboard.Length(ctx)
	if err != nil {
		panic(err)
	}
	switch r := lengthResponse.(type) {
	case *responses.LeaderboardLengthSuccess:
		fmt.Printf("Leaderboard length: %d\n", r.Length())
	}
}

func example_API_LeaderboardRemoveElements() {
	_, err := leaderboard.RemoveElements(ctx, momento.LeaderboardRemoveElementsRequest{Ids: []uint32{123, 456}})
	if err != nil {
		panic(err)
	}
}

func example_API_LeaderboardDelete() {
	_, err := leaderboard.Delete(ctx)
	if err != nil {
		panic(err)
	}
}

func example_patterns_ReadAsideCaching() string {
	key := uuid.NewString()
	resp, err := client.Get(ctx, &momento.GetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
	})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	// cache hit
	case *responses.GetHit:
		return r.ValueString()
	}
	// lookup value in database
	val := database[key]
	client.Set(ctx, &momento.SetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
		Value:     momento.String(val),
	})
	return val
}

func example_patterns_WriteThroughCaching() {
	key := uuid.NewString()
	value := uuid.NewString()
	// set value in database
	database[key] = value
	// set value in cache
	client.Set(ctx, &momento.SetRequest{
		CacheName: cacheName,
		Key:       momento.String(key),
		Value:     momento.String(value),
	})
}

// Clean up any lingering resources even if something panics
func cleanup() {
	if leaderboard != nil {
		leaderboard.Delete(ctx)
	}

	if client != nil {
		client.DeleteCache(ctx, &momento.DeleteCacheRequest{
			CacheName: cacheName,
		})
	}
}

func main() {
	defer cleanup()

	ctx = context.Background()
	cacheName = fmt.Sprintf("golang-docs-examples-%s", uuid.NewString())
	database = make(map[string]string)

	example_API_CredentialProviderFromString()
	example_API_CredentialProviderFromEnvVar()
	example_API_CredentialProviderFromEnvVarV2()
	example_API_CredentialProviderFromApiKeyV2()
	example_API_InstantiateCacheClientWithReadConcern()
	example_API_InstantiateCacheClient()

	example_API_CreateCache()
	example_API_ListCaches()

	example_API_Set()
	example_API_Get()
	example_API_Delete()

	example_API_SetIfPresent()
	example_API_SetIfAbsent()
	example_API_SetIfEqual()
	example_API_SetIfNotEqual()
	example_API_SetIfPresentAndNotEqual()
	example_API_SetIfAbsentOrEqual()

	example_API_SetWithHash()
	example_API_GetWithHash()
	example_API_SetIfPresentAndHashNotEqual()
	example_API_SetIfPresentAndHashEqual()
	example_API_SetIfAbsentOrHashNotEqual()
	example_API_SetIfAbsentOrHashEqual()

	example_API_KeysExist()
	example_API_ItemGetType()
	example_API_UpdateTtl()
	example_API_IncreaseTtl()
	example_API_DecreaseTtl()
	example_API_ItemGetTtl()
	example_API_Increment()

	example_API_SetBatch()
	example_API_GetBatch()

	example_API_InstantiateTopicClient()
	example_API_TopicPublish()
	example_API_TopicSubscribe()

	example_API_InstantiateLeaderboardClient()
	example_API_CreateLeaderboard()
	example_API_LeaderboardUpsert()
	example_API_LeaderboardFetchByScore()
	example_API_LeaderboardFetchByRank()
	example_API_LeaderboardGetRank()
	example_API_LeaderboardLength()
	example_API_LeaderboardRemoveElements()
	example_API_LeaderboardDelete()

	example_patterns_ReadAsideCaching()
	example_patterns_WriteThroughCaching()

	example_API_DeleteCache()

	example_API_InstantiateAuthClient()
	example_API_GenerateDisposableToken()
	example_API_GenerateApiKey()
	example_API_RefreshApiKey()
}
