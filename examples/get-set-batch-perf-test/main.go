package main

import (
	"context"
	"fmt"
	"github.com/loov/hrtime"
	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/go-example/get-set-batch-perf-test/utils"
	"strings"
	"sync"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
)

const (
	cacheName             = "go-perf-test"
	itemDefaultTTLSeconds = 60
	requestTimeoutSeconds = 600 // 10 minutes
)

func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}

	perfTestOptions := utils.PerfTestOptions{
		RequestTimeoutSeconds: requestTimeoutSeconds,
	}

	client, err := initializeMomentoClient(ctx, credentialProvider, perfTestOptions)
	if err != nil {
		panic(err)
	}

	batchSizeOptions := []int{5, 10}
	itemSizeOptions := []int{10, 100}
	testConfiguration := utils.PerfTestConfiguration{
		MinimumRunDurationSecondsForTests: 5,
		Sets:                              generateConfigurations(batchSizeOptions, itemSizeOptions),
		Gets:                              generateConfigurations(batchSizeOptions, itemSizeOptions),
	}

	runAsyncSetRequests(ctx, client, testConfiguration)
	runAsyncGetRequests(ctx, client, testConfiguration)
	runSetBatchRequests(ctx, client, testConfiguration)
	runGetBatchRequests(ctx, client, testConfiguration)
}

func initializeMomentoClient(ctx context.Context, credentialProvider auth.CredentialProvider, options utils.PerfTestOptions) (momento.CacheClient, error) {
	client, err := momento.NewCacheClientWithEagerConnectTimeout(
		config.LaptopLatest(),
		credentialProvider,
		itemDefaultTTLSeconds*time.Second,
		options.RequestTimeoutSeconds*time.Second,
	)
	if err != nil {
		return nil, err
	}

	_, err = client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func runAsyncSetRequests(ctx context.Context, client momento.CacheClient, testConfiguration utils.PerfTestConfiguration) {
	for _, setConfig := range testConfiguration.Sets {
		perfTestContext := utils.InitiatePerfTestContext()
		for hrtime.Since(perfTestContext.StartTime).Seconds() < float64(testConfiguration.MinimumRunDurationSecondsForTests) {
			sendAsyncSetRequests(ctx, client, perfTestContext, setConfig)
		}
		fmt.Printf("Completed async set requests with batch size %d and item size %d\n", setConfig.BatchSize, setConfig.ItemSizeBytes)
		utils.CalculateSummary(perfTestContext, setConfig.BatchSize, setConfig.ItemSizeBytes, utils.AsyncSets)
	}
}

func runAsyncGetRequests(ctx context.Context, client momento.CacheClient, testConfiguration utils.PerfTestConfiguration) {
	for _, getConfig := range testConfiguration.Gets {
		ensureCacheIsPopulated(ctx, client, getConfig)
		perfTestContext := utils.InitiatePerfTestContext()
		for hrtime.Since(perfTestContext.StartTime).Seconds() < float64(testConfiguration.MinimumRunDurationSecondsForTests) {
			sendAsyncGetRequests(ctx, client, perfTestContext, getConfig)
		}
		fmt.Printf("Completed async get requests with batch size %d and item size %d\n", getConfig.BatchSize, getConfig.ItemSizeBytes)
		utils.CalculateSummary(perfTestContext, getConfig.BatchSize, getConfig.ItemSizeBytes, utils.AsyncGets)
	}
}

func runSetBatchRequests(ctx context.Context, client momento.CacheClient, testConfiguration utils.PerfTestConfiguration) {
	for _, setConfig := range testConfiguration.Sets {
		perfTestContext := utils.InitiatePerfTestContext()
		for hrtime.Since(perfTestContext.StartTime).Seconds() < float64(testConfiguration.MinimumRunDurationSecondsForTests) {
			sendSetBatchRequests(ctx, client, perfTestContext, setConfig)
		}
		fmt.Printf("Completed set batch requests with batch size %d and item size %d\n", setConfig.BatchSize, setConfig.ItemSizeBytes)
		utils.CalculateSummary(perfTestContext, setConfig.BatchSize, setConfig.ItemSizeBytes, utils.SetBatch)
	}
}

func runGetBatchRequests(ctx context.Context, client momento.CacheClient, testConfiguration utils.PerfTestConfiguration) {
	for _, getConfig := range testConfiguration.Gets {
		ensureCacheIsPopulated(ctx, client, getConfig)
		perfTestContext := utils.InitiatePerfTestContext()
		for hrtime.Since(perfTestContext.StartTime).Seconds() < float64(testConfiguration.MinimumRunDurationSecondsForTests) {
			sendGetBatchRequests(ctx, client, perfTestContext, getConfig)
		}
		fmt.Printf("Completed get batch requests with batch size %d and item size %d\n", getConfig.BatchSize, getConfig.ItemSizeBytes)
		utils.CalculateSummary(perfTestContext, getConfig.BatchSize, getConfig.ItemSizeBytes, utils.GetBatch)
	}
}

func sendAsyncSetRequests(ctx context.Context, client momento.CacheClient, context *utils.PerfTestContext, setConfig utils.GetSetConfig) {
	var wg sync.WaitGroup
	setResponses := make([]responses.SetResponse, setConfig.BatchSize)
	value := strings.Repeat("x", setConfig.ItemSizeBytes)
	setStartTime := hrtime.Now()

	for i := 0; i < setConfig.BatchSize; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i)
			setPromise, err := client.Set(ctx, &momento.SetRequest{
				CacheName: cacheName,
				Key:       momento.String(key),
				Value:     momento.String(value),
			})
			if err != nil {
				panic(err) // You might want to handle the error more gracefully
			}
			setResponses[i] = setPromise
			context.TotalItemSizeBytes += int64(setConfig.ItemSizeBytes)
		}(i)
	}
	// Wait for all goroutines to finish
	wg.Wait()
	// Calculate elapsed time and record it
	setDuration := hrtime.Since(setStartTime)
	err := context.AsyncSetLatencies.RecordValue(setDuration.Milliseconds())
	if err != nil {
		return
	}
}

func sendAsyncGetRequests(ctx context.Context, client momento.CacheClient, context *utils.PerfTestContext, getConfig utils.GetSetConfig) {
	var wg sync.WaitGroup
	getResponses := make([]responses.GetResponse, getConfig.BatchSize)
	getStartTime := hrtime.Now()

	for i := 0; i < getConfig.BatchSize; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i)
			getPromise, err := client.Get(ctx, &momento.GetRequest{
				CacheName: cacheName,
				Key:       momento.String(key),
			})
			if err != nil {
				panic(err) // You might want to handle the error more gracefully
			}
			getResponses[i] = getPromise
			context.TotalItemSizeBytes += int64(getConfig.ItemSizeBytes)
		}(i)
	}
	// Wait for all goroutines to finish
	wg.Wait()
	// Calculate elapsed time and record it
	getDuration := hrtime.Since(getStartTime)
	err := context.AsyncGetLatencies.RecordValue(getDuration.Milliseconds())
	if err != nil {
		return
	}
}

func sendSetBatchRequests(ctx context.Context, client momento.CacheClient, context *utils.PerfTestContext, setConfig utils.GetSetConfig) {
	keys := make([]string, setConfig.BatchSize)
	value := strings.Repeat("x", setConfig.ItemSizeBytes)
	items := make([]momento.BatchSetItem, setConfig.BatchSize)
	for i := 0; i < setConfig.BatchSize; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
		items[i] = momento.BatchSetItem{
			Key:   momento.String(keys[i]),
			Value: momento.String(value),
		}
	}
	setBatchStartTime := hrtime.Now()
	_, err := client.SetBatch(ctx, &momento.SetBatchRequest{
		CacheName: cacheName,
		Items:     items,
	})
	if err != nil {
		panic(err)
	}
	setBatchDuration := hrtime.Since(setBatchStartTime)
	err = context.SetBatchLatencies.RecordValue(setBatchDuration.Milliseconds())
	if err != nil {
		return
	}
	context.TotalItemSizeBytes += int64(setConfig.BatchSize * setConfig.ItemSizeBytes)
}

func sendGetBatchRequests(ctx context.Context, client momento.CacheClient, context *utils.PerfTestContext, getConfig utils.GetSetConfig) {
	keys := make([]momento.Value, getConfig.BatchSize)
	for i := 0; i < getConfig.BatchSize; i++ {
		keys[i] = momento.String(fmt.Sprintf("key-%d", i))
	}
	getBatchStartTime := hrtime.Now()
	_, err := client.GetBatch(ctx, &momento.GetBatchRequest{
		CacheName: cacheName,
		Keys:      keys,
	})
	if err != nil {
		panic(err)
	}
	getBatchDuration := hrtime.Since(getBatchStartTime)
	err = context.GetBatchLatencies.RecordValue(getBatchDuration.Milliseconds())
	if err != nil {
		return
	}
	context.TotalItemSizeBytes += int64(getConfig.BatchSize * getConfig.ItemSizeBytes)
}

func ensureCacheIsPopulated(ctx context.Context, client momento.CacheClient, getConfig utils.GetSetConfig) {
	keys := make([]string, getConfig.BatchSize)
	value := strings.Repeat("x", getConfig.ItemSizeBytes)

	for i := 0; i < getConfig.BatchSize; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
	}

	for i := 0; i < getConfig.BatchSize; i++ {
		_, err := client.SetIfAbsent(ctx, &momento.SetIfAbsentRequest{
			CacheName: cacheName,
			Key:       momento.String(keys[i]),
			Value:     momento.String(value),
		})
		if err != nil {
			panic(err)
		}
	}
}

func generateConfigurations(batchSizes []int, itemSizes []int) []utils.GetSetConfig {
	var configurations []utils.GetSetConfig
	for _, batchSize := range batchSizes {
		for _, itemSize := range itemSizes {
			configurations = append(configurations, utils.GetSetConfig{BatchSize: batchSize, ItemSizeBytes: itemSize})
		}
	}
	return configurations
}
