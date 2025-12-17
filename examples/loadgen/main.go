package main

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
	"github.com/loov/hrtime"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/momento"
)

const (
	CacheItemTtlSeconds = 60
	CacheName           = "momento-loadgen"
)

type loadGeneratorOptions struct {
	logLevel              momento_default_logger.LogLevel
	showStatsInterval     time.Duration
	cacheItemPayloadBytes int
	// Note: You are likely to see degraded performance if you increase this above 50
	// and observe elevated client-side latencies.
	numberOfConcurrentRequests int
	maxRequestsPerSecond       int
	howLongToRun               time.Duration
}

type loadGenerator struct {
	loggerFactory       logger.MomentoLoggerFactory
	logger              logger.MomentoLogger
	momentoClientConfig config.Configuration
	options             loadGeneratorOptions
	cacheValue          string
}

type ErrorCounter struct {
	unavailable   int64
	timeout       int64
	limitExceeded int64
	canceled      int64
}

func newLoadGenerator(config config.Configuration, options loadGeneratorOptions) *loadGenerator {
	loggerFactory := config.GetLoggerFactory()
	lgLogger := loggerFactory.GetLogger("loadgen")
	cacheValue := strings.Repeat("x", options.cacheItemPayloadBytes)
	return &loadGenerator{
		loggerFactory:       loggerFactory,
		logger:              lgLogger,
		momentoClientConfig: config,
		options:             options,
		cacheValue:          cacheValue,
	}
}

func (r *loadGenerator) init(ctx context.Context) (momento.CacheClient, time.Duration) {
	credentialProvider, err := auth.FromEnvironmentVariablesV2()
	if err != nil {
		panic(err)
	}
	client, err := momento.NewCacheClientWithEagerConnectTimeout(r.momentoClientConfig, credentialProvider, time.Second*CacheItemTtlSeconds, 30*time.Second)
	if err != nil {
		panic(err)
	}
	if _, err := client.CreateCache(ctx, &momento.CreateCacheRequest{CacheName: CacheName}); err != nil {
		panic(err)
	}

	workerDelayBetweenRequests := int32(math.Floor(
		(1000.0 * float64(r.options.numberOfConcurrentRequests)) / float64(r.options.maxRequestsPerSecond),
	))
	delay := time.Duration(workerDelayBetweenRequests) * time.Millisecond

	r.logger.Debug(
		fmt.Sprintf(
			"Targeting a max of %d requests per second (delay between requests: %d ms)",
			r.options.maxRequestsPerSecond,
			delay.Milliseconds(),
		),
	)
	r.logger.Debug(
		fmt.Sprintf(
			"Running %d concurrent requests for %d seconds",
			r.options.numberOfConcurrentRequests,
			int(r.options.howLongToRun.Seconds())),
	)

	return client, delay
}

func processError(err error, errChan chan string) {
	switch mErr := err.(type) {
	case momento.MomentoError:
		if mErr.Code() == momento.ServerUnavailableError ||
			mErr.Code() == momento.TimeoutError ||
			mErr.Code() == momento.LimitExceededError ||
			mErr.Code() == momento.CanceledError {
			errChan <- mErr.Code()
		} else {
			panic(fmt.Sprintf("unrecognized result: %T", mErr))
		}
	default:
		panic(fmt.Sprintf("unknown error type %T", err))
	}
}

func worker(
	ctx context.Context,
	id int,
	getChan chan int64,
	setChan chan int64,
	errChan chan string,
	client momento.CacheClient,
	workerDelayBetweenRequests time.Duration,
	cacheValue string,
) {
	i := 0
	var elapsed time.Duration
	for {
		select {
		case <-ctx.Done():
			return
		default:
			i++
			elapsed = 0
			cacheKey := fmt.Sprintf("worker%doperation%d", id, i)
			setStart := hrtime.Now()
			_, err := client.Set(ctx, &momento.SetRequest{
				CacheName: CacheName,
				Key:       momento.String(cacheKey),
				Value:     momento.String(cacheValue),
			})
			if err != nil {
				processError(err, errChan)
			} else {
				elapsed = hrtime.Since(setStart)
				setChan <- elapsed.Milliseconds()
			}

			if elapsed.Milliseconds() < workerDelayBetweenRequests.Milliseconds() {
				time.Sleep(workerDelayBetweenRequests - elapsed)
			}

			elapsed = 0
			getStart := hrtime.Now()
			_, err = client.Get(ctx, &momento.GetRequest{
				CacheName: CacheName,
				Key:       momento.String(cacheKey),
			})
			if err != nil {
				processError(err, errChan)
			} else {
				elapsed = hrtime.Since(getStart)
				getChan <- elapsed.Milliseconds()
			}

			if elapsed.Milliseconds() < workerDelayBetweenRequests.Milliseconds() {
				time.Sleep(workerDelayBetweenRequests - elapsed)
			}
		}
	}
}

func printStats(gets *hdrhistogram.Histogram, sets *hdrhistogram.Histogram, errorCounter ErrorCounter, startTime time.Duration) {
	totalRequests := gets.TotalCount() + sets.TotalCount() + errorCounter.timeout + errorCounter.unavailable + errorCounter.limitExceeded + errorCounter.canceled
	totalTps := int(math.Round(
		float64(totalRequests * 1000 / hrtime.Since(startTime).Milliseconds()),
	))
	fmt.Println("==============================\ncumulative stats:")
	fmt.Printf(
		"%20s: %d (%d tps)\n",
		"total requests",
		totalRequests,
		totalTps,
	)

	successfulRequests := gets.TotalCount() + sets.TotalCount()
	successTps := int(math.Round(
		float64(successfulRequests * 1000 / hrtime.Since(startTime).Milliseconds()),
	))
	fmt.Printf("%20s: %d (%d%%) (%d tps)\n", "success", successfulRequests, successfulRequests/totalRequests*100, successTps)
	fmt.Printf("%20s: %d (%d%%)\n", "unavailable", errorCounter.unavailable, errorCounter.unavailable/totalRequests*100)
	fmt.Printf("%20s: %d (%d%%)\n", "timeout exceeded", errorCounter.timeout, errorCounter.timeout/totalRequests*100)
	fmt.Printf("%20s: %d (%d%%)\n", "limit exceeded", errorCounter.limitExceeded, errorCounter.limitExceeded/totalRequests*100)
	fmt.Printf("%20s: %d (%d%%)\n\n", "canceled", errorCounter.canceled, errorCounter.canceled/totalRequests*100)

	fmt.Printf(
		"cumulative get latencies:\n%20s: %d\n%20s: %d\n%20s: %d\n%20s: %d\n%20s: %d\n%20s: %d\n\n",
		"total requests",
		gets.TotalCount(),
		"p50",
		gets.ValueAtQuantile(50.0),
		"p90",
		gets.ValueAtQuantile(90.0),
		"p99",
		gets.ValueAtQuantile(99.0),
		"p99.9",
		gets.ValueAtQuantile(99.9),
		"max",
		gets.Max(),
	)

	fmt.Printf(
		"cumulative set latencies:\n%20s: %d\n%20s: %d\n%20s: %d\n%20s: %d\n%20s: %d\n%20s: %d\n\n",
		"total requests",
		sets.TotalCount(),
		"p50",
		sets.ValueAtQuantile(50.0),
		"p90",
		sets.ValueAtQuantile(90.0),
		"p99",
		sets.ValueAtQuantile(99.0),
		"p99.9",
		sets.ValueAtQuantile(99.9),
		"max",
		sets.Max(),
	)
}

func (ec *ErrorCounter) updateErrors(err string) {
	if err == momento.ServerUnavailableError {
		ec.unavailable++
	} else if err == momento.TimeoutError {
		ec.timeout++
	} else if err == momento.LimitExceededError {
		ec.limitExceeded++
	} else if err == momento.CanceledError {
		ec.canceled++
	}
}

func timer(ctx context.Context, getChan chan int64, setChan chan int64, errChan chan string, statsInterval time.Duration) {
	getHistogram := hdrhistogram.New(1, 10000000000, 3)
	setHistogram := hdrhistogram.New(1, 10000000000, 3)
	errorCounter := ErrorCounter{}

	startTime := hrtime.Now()
	origStartTime := startTime

	for {

		if hrtime.Since(startTime) >= statsInterval {
			printStats(getHistogram, setHistogram, errorCounter, origStartTime)
			startTime = hrtime.Now()
		}

		select {
		case <-ctx.Done():
			fmt.Println("\n=====> run complete <====")
			printStats(getHistogram, setHistogram, errorCounter, origStartTime)
			return
		case getMessage := <-getChan:
			if err := getHistogram.RecordValue(getMessage); err != nil {
				panic(err)
			}
		case setMessage := <-setChan:
			if err := setHistogram.RecordValue(setMessage); err != nil {
				panic(err)
			}
		case errCode := <-errChan:
			errorCounter.updateErrors(errCode)
		default:
			time.Sleep(time.Millisecond * 25)
		}
	}
}

func (r *loadGenerator) run(ctx context.Context, client momento.CacheClient, workerDelayBetweenRequests time.Duration) {
	cancelCtx, cancelFunction := context.WithTimeout(ctx, r.options.howLongToRun)
	defer cancelFunction()

	var wg sync.WaitGroup

	// Setting the channel length to a max of number of concurrent requests was a bottleneck.
	//Hence, using a large number to ensure the channels were not getting clogged up.
	getChan := make(chan int64, 10_000)
	setChan := make(chan int64, 10_000)
	errChan := make(chan string, 10_000)

	wg.Add(1)
	go func() {
		defer wg.Done()
		timer(cancelCtx, getChan, setChan, errChan, r.options.showStatsInterval)
	}()

	// Launch and run workers
	for i := 1; i <= r.options.numberOfConcurrentRequests; i++ {
		wg.Add(1)

		// avoid reuse of the same i value in each closure
		i := i

		go func() {
			defer wg.Done()
			worker(cancelCtx, i, getChan, setChan, errChan, client, workerDelayBetweenRequests, r.cacheValue)
		}()
	}

	wg.Wait()

}

func main() {
	ctx := context.Background()

	opts := loadGeneratorOptions{
		logLevel:              momento_default_logger.DEBUG,
		showStatsInterval:     time.Second * 5,
		cacheItemPayloadBytes: 100,
		// Note: You are likely to see degraded performance if you increase this above 50
		// and observe elevated client-side latencies.
		numberOfConcurrentRequests: 50,
		maxRequestsPerSecond:       100,
		howLongToRun:               time.Minute,
	}

	momentoConfig := config.LaptopLatestWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(opts.logLevel))
	loadGenerator := newLoadGenerator(momentoConfig, opts)
	client, workerDelayBetweenRequests := loadGenerator.init(ctx)

	runStart := time.Now()
	loadGenerator.run(ctx, client, workerDelayBetweenRequests)
	runTotal := time.Since(runStart)

	fmt.Printf("completed in %f seconds\n", runTotal.Seconds())
}
