package main

import (
	"context"
	"fmt"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}

	const (
		cacheName             = "cache"
		itemDefaultTTLSeconds = 60
	)

	// Instantiate leaderboard client
	leaderboardClient, err := momento.NewPreviewLeaderboardClient(
		config.LeaderboardDefault(),
		credentialProvider,
	)
	if err != nil {
		panic(err)
	}

	// Create a new leaderboard
	leaderboard, err := leaderboardClient.Leaderboard(ctx, &momento.LeaderboardRequest{
		CacheName:       cacheName,
		LeaderboardName: "leaderboard",
	})
	if err != nil {
		panic(err)
	}

	// Upsert elements
	upsertElements := []momento.LeaderboardUpsertElement{
		{Id: 123, Score: 10.33},
		{Id: 234, Score: 123.4567},
		{Id: 345, Score: 20.4},
		{Id: 456, Score: 3333},
		{Id: 567, Score: 9876},
		{Id: 678, Score: 1.1111},
		{Id: 789, Score: 5678.9},
	}
	_, err = leaderboard.Upsert(ctx, momento.LeaderboardUpsertRequest{Elements: upsertElements})
	if err != nil {
		panic(err)
	}

	// Get leaderboard length
	lengthResponse, err := leaderboard.Length(ctx)
	if err != nil {
		panic(err)
	} else {
		switch r := lengthResponse.(type) {
		case *responses.LeaderboardLengthSuccess:
			fmt.Printf("Leaderboard length: %d\n", r.Length())
		}
	}

	// Fetch elements by rank
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

	// Fetch elements by score
	minScore := 150.0
	maxScore := 3000.0
	offset := uint32(1)
	count := uint32(2)
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

	// Get elements by ID
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

	// Remove elements from leaderboard
	_, err = leaderboard.RemoveElements(ctx, momento.LeaderboardRemoveElementsRequest{Ids: []uint32{123, 456}})
	if err != nil {
		panic(err)
	}

	// Delete all elements from the leaderboard
	_, err = leaderboard.Delete(ctx)
	if err != nil {
		panic(err)
	}

	// Close leaderboard client
	leaderboardClient.Close()
}
