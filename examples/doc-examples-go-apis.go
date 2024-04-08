package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/momentohq/client-sdk-go/utils"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
	auth_resp "github.com/momentohq/client-sdk-go/responses/auth"
)

var (
	ctx               context.Context
	client            momento.CacheClient
	cacheName         string
	leaderboardClient momento.PreviewLeaderboardClient
	leaderboard       momento.Leaderboard
)

func example_API_InstantiateCacheClient() {
	context := context.Background()
	credentialProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}
	defaultTtl := time.Duration(9999)
	eagerConnectTimeout := 30 * time.Second

	client, err := momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatest(),
		credentialProvider,
		defaultTtl,
		eagerConnectTimeout,
	)
	if err != nil {
		panic(err)
	}

	client.Ping(context)
}

func example_API_InstantiateCacheClientWithReadConcern() {
	context := context.Background()
	credentialProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
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

	client.Ping(context)
}

func example_API_ListCaches() {
	resp, err := client.ListCaches(ctx, &momento.ListCachesRequest{})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *responses.ListCachesSuccess:
		log.Printf("Found caches %+v", r.Caches())
	}
}

func example_API_CreateCache() {
	_, err := client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: "cache-name",
	})
	if err != nil {
		panic(err)
	}
}

func example_API_DeleteCache() {
	_, err := client.DeleteCache(ctx, &momento.DeleteCacheRequest{
		CacheName: "cache-name",
	})
	if err != nil {
		panic(err)
	}
}

func example_API_Get() {
	key := uuid.NewString()
	resp, err := client.Get(ctx, &momento.GetRequest{
		CacheName: "cache-name",
		Key:       momento.String(key),
	})
	if err != nil {
		panic(err)
	}

	switch r := resp.(type) {
	case *responses.GetHit:
		log.Printf("Lookup resulted in cache HIT. value=%s\n", r.ValueString())
	case *responses.GetMiss:
		log.Printf("Look up did not find a value key=%s", key)
	}
}

func example_API_Set() {
	key := uuid.NewString()
	value := uuid.NewString()
	log.Printf("Setting key: %s, value: %s\n", key, value)
	_, err := client.Set(ctx, &momento.SetRequest{
		CacheName: "cache-name",
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
		CacheName: "cache-name",
		Key:       momento.String(key),
	})
	if err != nil {
		panic(err)
	}
}

func example_API_InstantiateTopicClient() {
	credProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}

	topicClient, err := momento.NewTopicClient(
		config.TopicsDefault(),
		credProvider,
	)
	if err != nil {
		panic(err)
	}
}

func example_API_TopicPublish(client momento.TopicClient) {
	_, err := client.Publish(ctx, &momento.TopicPublishRequest{
		CacheName: "test-cache",
		TopicName: "test-topic",
		Value:     momento.String("test-message"),
	})
	if err != nil {
		panic(err)
	}
}

func example_API_TopicSubscribe(client momento.TopicClient) {
	// Instantiate subscriber
	sub, subErr := client.Subscribe(ctx, &momento.TopicSubscribeRequest{
		CacheName: "test-cache",
		TopicName: "test-topic",
	})
	if subErr != nil {
		panic(subErr)
	}

	time.Sleep(time.Second)
	_, pubErr := client.Publish(ctx, &momento.TopicPublishRequest{
		CacheName: "test-cache",
		TopicName: "test-topic",
		Value:     momento.String("test-message"),
	})
	if pubErr != nil {
		panic(pubErr)
	}
	time.Sleep(time.Second)

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
}

func example_API_GenerateDisposableToken(client momento.AuthClient) {
	tokenId := "a token id"
	resp, err := client.GenerateDisposableToken(ctx, &momento.GenerateDisposableTokenRequest{
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

func example_API_SetIfPresent() {
	resp, err := client.SetIfPresent(ctx, &momento.SetIfPresentRequest{
		CacheName: "cache-name",
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
		CacheName: "cache-name",
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
		CacheName: "cache-name",
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
		CacheName: "cache-name",
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
		CacheName: "cache-name",
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
		CacheName: "cache-name",
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

func example_API_KeysExist() {
	keys := []momento.Value{momento.String("key1"), momento.String("key2")}
	resp, err := client.KeysExist(ctx, &momento.KeysExistRequest{
		CacheName: "cache-name",
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
		CacheName: "cache-name",
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
		CacheName: "cache-name",
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
		CacheName: "cache-name",
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
		CacheName: "cache-name",
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
		CacheName: "cache-name",
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
		CacheName: "cache-name",
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
		CacheName: "cache-name",
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
		CacheName: "cache-name",
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
	credentialProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
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

func example_API_CreateLeaderboard() {
	leaderboard, err := leaderboardClient.Leaderboard(ctx, &momento.LeaderboardRequest{
		CacheName:       cacheName,
		LeaderboardName: "leaderboard",
	})
	if err != nil {
		panic(err)
	}
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
	} else {
		switch r := fetchByScoreResponse.(type) {
		case *responses.LeaderboardFetchSuccess:
			fmt.Printf("Successfully fetched elements by score:\n")
			for _, element := range r.Values() {
				fmt.Printf("ID: %d, Score: %f, Rank: %d\n", element.Id, element.Score, element.Rank)
			}
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
	} else {
		switch r := fetchByRankResponse.(type) {
		case *responses.LeaderboardFetchSuccess:
			fmt.Printf("Successfully fetched elements by rank:\n")
			for _, element := range r.Values() {
				fmt.Printf("ID: %d, Score: %f, Rank: %d\n", element.Id, element.Score, element.Rank)
			}
		}
	}
}

func example_API_LeaderboardGetRank() {
	getRankResponse, err := leaderboard.GetRank(ctx, momento.LeaderboardGetRankRequest{
		Ids: []uint32{123, 456},
	})
	if err != nil {
		panic(err)
	} else {
		switch r := getRankResponse.(type) {
		case *responses.LeaderboardFetchSuccess:
			fmt.Printf("Successfully fetched elements by ID:\n")
			for _, element := range r.Values() {
				fmt.Printf("ID: %d, Score: %f, Rank: %d\n", element.Id, element.Score, element.Rank)
			}
		}
	}
}

func example_API_LeaderboardLength() {
	lengthResponse, err := leaderboard.Length(ctx)
	if err != nil {
		panic(err)
	} else {
		switch r := lengthResponse.(type) {
		case *responses.LeaderboardLengthSuccess:
			fmt.Printf("Leaderboard length: %d\n", r.Length())
		}
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
