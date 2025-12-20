package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/loov/hrtime"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/go-example/get-set-batch-perf-test/utils"
)

const (
	cacheName             = "go-perf-test"
	itemDefaultTTLSeconds = 60
	requestTimeoutSeconds = 600 // 10 minutes
	maxRequestsPerSecond  = 10_000
)

func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoV2TokenProvider()
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

	batchSizeOptions := []int{5, 10, 100, 500, 1_000, 5_000, 10_000}
	itemSizeOptions := []int{10, 100, 1024, 1024 * 10, 1024 * 100, 1024 * 1024}
	testConfiguration := utils.PerfTestConfiguration{
		MinimumRunDurationSecondsForTests: 60,
		Sets:                              generateConfigurations(batchSizeOptions, itemSizeOptions),
		Gets:                              generateConfigurations(batchSizeOptions, itemSizeOptions),
	}

	runAsyncSetRequests(ctx, client, testConfiguration)
	runAsyncGetRequests(ctx, client, testConfiguration)
	runSetBatchRequests(ctx, client, testConfiguration)
	runGetBatchRequests(ctx, client, testConfiguration)
}

func initializeMomentoClient(ctx context.Context, credentialProvider auth.CredentialProvider, options utils.PerfTestOptions) (momento.CacheClient, error) {
	client, err := momento.NewCacheClient(
		config.LaptopLatest().WithClientTimeout(options.RequestTimeoutSeconds*time.Second),
		credentialProvider,
		itemDefaultTTLSeconds*time.Second,
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
		fmt.Println("Beginning run for ASYNC_SETS, batch size:", setConfig.BatchSize, "item size:", setConfig.ItemSizeBytes)
		numLoops := 0
		perfTestContext := utils.InitiatePerfTestContext()
		for hrtime.Since(perfTestContext.StartTime).Seconds() < float64(testConfiguration.MinimumRunDurationSecondsForTests) {
			numLoops++
			workerDelay := (time.Second / time.Duration(maxRequestsPerSecond)) * time.Duration(setConfig.BatchSize)
			sendAsyncSetRequests(ctx, client, perfTestContext, setConfig, workerDelay)
		}
		utils.CalculateSummary(perfTestContext, setConfig.BatchSize, setConfig.ItemSizeBytes, utils.AsyncSets)
		// calculate tps and print
		tps := float64(perfTestContext.TotalNumberOfRequests) / hrtime.Since(perfTestContext.StartTime).Seconds()
		fmt.Println("Total TPS:", tps)
		fmt.Println("Completed ASYNC_SETS requests with batch size:", setConfig.BatchSize, "item size:", setConfig.ItemSizeBytes, "numLoops:", numLoops, "elapsedTimeMillis:", hrtime.Since(perfTestContext.StartTime).Milliseconds())
	}
}

func runAsyncGetRequests(ctx context.Context, client momento.CacheClient, testConfiguration utils.PerfTestConfiguration) {
	for _, getConfig := range testConfiguration.Gets {
		fmt.Println("Populating cache for ASYNC_GETS, batch size:", getConfig.BatchSize, "item size:", getConfig.ItemSizeBytes)
		var cachePopulationStartTime = hrtime.Now()
		ensureCacheIsPopulated(ctx, client, getConfig)
		fmt.Printf("Populated cache with batch size %d and item size %d in %d ms\n", getConfig.BatchSize, getConfig.ItemSizeBytes, hrtime.Since(cachePopulationStartTime).Milliseconds())

		fmt.Println("Beginning run for ASYNC_GETS, batch size:", getConfig.BatchSize, "item size:", getConfig.ItemSizeBytes)
		numLoops := 0
		perfTestContext := utils.InitiatePerfTestContext()
		for hrtime.Since(perfTestContext.StartTime).Seconds() < float64(testConfiguration.MinimumRunDurationSecondsForTests) {
			numLoops++
			workerDelay := (time.Second / time.Duration(maxRequestsPerSecond)) * time.Duration(getConfig.BatchSize)
			sendAsyncGetRequests(ctx, client, perfTestContext, getConfig, workerDelay)
		}
		utils.CalculateSummary(perfTestContext, getConfig.BatchSize, getConfig.ItemSizeBytes, utils.AsyncGets)
		// calculate tps and print
		tps := float64(perfTestContext.TotalNumberOfRequests) / hrtime.Since(perfTestContext.StartTime).Seconds()
		fmt.Println("Total TPS:", tps)
		fmt.Println("Completed ASYNC_GETS requests with batch size:", getConfig.BatchSize, "item size:", getConfig.ItemSizeBytes, "numLoops:", numLoops, "elapsedTimeMillis:", hrtime.Since(perfTestContext.StartTime).Milliseconds())
	}
}

func runSetBatchRequests(ctx context.Context, client momento.CacheClient, testConfiguration utils.PerfTestConfiguration) {
	for _, setConfig := range testConfiguration.Sets {
		if setConfig.BatchSize*setConfig.ItemSizeBytes >= 5*1024*1024 {
			fmt.Printf("Skipping run for SET_BATCH with batch size %d and item size %d\n", setConfig.BatchSize, setConfig.ItemSizeBytes)
			continue
		}
		fmt.Println("Beginning run for SET_BATCH, batch size:", setConfig.BatchSize, "item size:", setConfig.ItemSizeBytes)
		numLoops := 0
		perfTestContext := utils.InitiatePerfTestContext()
		for hrtime.Since(perfTestContext.StartTime).Seconds() < float64(testConfiguration.MinimumRunDurationSecondsForTests) {
			numLoops++
			sendSetBatchRequests(ctx, client, perfTestContext, setConfig)
		}
		utils.CalculateSummary(perfTestContext, setConfig.BatchSize, setConfig.ItemSizeBytes, utils.SetBatch)
		// calculate tps and print
		tps := float64(perfTestContext.TotalNumberOfRequests) / hrtime.Since(perfTestContext.StartTime).Seconds()
		fmt.Println("Total TPS:", tps)
		fmt.Printf("Completed SET_BATCH requests with batch size %d and item size %d\n", setConfig.BatchSize, setConfig.ItemSizeBytes)
	}
}

func runGetBatchRequests(ctx context.Context, client momento.CacheClient, testConfiguration utils.PerfTestConfiguration) {
	for _, getConfig := range testConfiguration.Gets {
		fmt.Println("Populating cache for GET_BATCH, batch size:", getConfig.BatchSize, "item size:", getConfig.ItemSizeBytes)
		var cachePopulationStartTime = hrtime.Now()
		ensureCacheIsPopulated(ctx, client, getConfig)
		fmt.Printf("Populated cache with batch size %d and item size %d in %d ms\n", getConfig.BatchSize, getConfig.ItemSizeBytes, hrtime.Since(cachePopulationStartTime).Milliseconds())

		fmt.Println("Beginning run for GET_BATCH, batch size:", getConfig.BatchSize, "item size:", getConfig.ItemSizeBytes)
		numLoops := 0
		perfTestContext := utils.InitiatePerfTestContext()
		for hrtime.Since(perfTestContext.StartTime).Seconds() < float64(testConfiguration.MinimumRunDurationSecondsForTests) {
			numLoops++
			sendGetBatchRequests(ctx, client, perfTestContext, getConfig)
		}
		utils.CalculateSummary(perfTestContext, getConfig.BatchSize, getConfig.ItemSizeBytes, utils.GetBatch)
		// calculate tps and print
		tps := float64(perfTestContext.TotalNumberOfRequests) / hrtime.Since(perfTestContext.StartTime).Seconds()
		fmt.Println("Total TPS:", tps)
		fmt.Println("Completed GET_BATCH requests with batch size:", getConfig.BatchSize, "item size:", getConfig.ItemSizeBytes, "numLoops:", numLoops, "elapsedTimeMillis:", hrtime.Since(perfTestContext.StartTime).Milliseconds())
	}
}

func sendAsyncSetRequests(ctx context.Context, client momento.CacheClient, context *utils.PerfTestContext, setConfig utils.GetSetConfig, workerDelay time.Duration) {
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
				panic(err)
			}
			setResponses[i] = setPromise
			context.TotalItemSizeBytes += int64(setConfig.ItemSizeBytes)
		}(i)
	}
	// Wait for all goroutines to finish
	wg.Wait()

	// Calculate total number of requests
	context.TotalNumberOfRequests += int64(setConfig.BatchSize)

	// Calculate elapsed time and record it
	setDuration := hrtime.Since(setStartTime)
	err := context.AsyncSetLatencies.RecordValue(setDuration.Milliseconds())

	if setDuration.Milliseconds() < workerDelay.Milliseconds() {
		fmt.Println("Sleeping for", workerDelay.Milliseconds()-setDuration.Milliseconds(), "ms")
		time.Sleep(workerDelay - setDuration)
	}

	if err != nil {
		return
	}
}

func sendAsyncGetRequests(ctx context.Context, client momento.CacheClient, context *utils.PerfTestContext, getConfig utils.GetSetConfig, workerDelay time.Duration) {
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
	context.TotalNumberOfRequests += int64(getConfig.BatchSize)

	// Calculate elapsed time and record it
	getDuration := hrtime.Since(getStartTime)
	err := context.AsyncGetLatencies.RecordValue(getDuration.Milliseconds())

	if getDuration.Milliseconds() < workerDelay.Milliseconds() {
		fmt.Println("Sleeping for", workerDelay.Milliseconds()-getDuration.Milliseconds(), "ms")
		time.Sleep(workerDelay - getDuration)
	}

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
	context.TotalNumberOfRequests += 1
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
	context.TotalNumberOfRequests += 1
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
			// exclude permutations where total payload is greater than 1GB, they are not realistic and will cause OOM
			if batchSize*itemSize < 1024*1024*1024 {
				configurations = append(configurations, utils.GetSetConfig{BatchSize: batchSize, ItemSizeBytes: itemSize})
			}
		}
	}
	return configurations
}
