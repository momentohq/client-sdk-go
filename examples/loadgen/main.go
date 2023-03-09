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
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

const (
	CacheItemTtlSeconds = 60
	CacheName           = "momento-loadgen"
)

type loadGeneratorOptions struct {
	logLevel                   logger.LogLevel
	showStatsInterval          time.Duration
	cacheItemPayloadBytes      int
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
	unavailable   int
	timeout       int
	limitExceeded int
}

func newLoadGenerator(config config.Configuration, options loadGeneratorOptions) *loadGenerator {
	loggerFactory := config.GetLoggerFactory()
	lgLogger := loggerFactory.GetLogger("loadgen", options.logLevel)
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
	credentialProvider, err := auth.FromEnvironmentVariable("MOMENTO_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}
	client, err := momento.NewCacheClient(r.momentoClientConfig, credentialProvider, time.Second*CacheItemTtlSeconds)
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
			mErr.Code() == momento.LimitExceededError {
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
			getResp, err := client.Get(ctx, &momento.GetRequest{
				CacheName: CacheName,
				Key:       momento.String(cacheKey),
			})
			if err != nil {
				processError(err, errChan)
			} else {
				switch getResp.(type) {
				case *responses.GetHit:
					elapsed = hrtime.Since(getStart)
					getChan <- elapsed.Milliseconds()
				}
			}

			if elapsed.Milliseconds() < workerDelayBetweenRequests.Milliseconds() {
				time.Sleep(workerDelayBetweenRequests - elapsed)
			}
		}
	}
}

func printStats(gets []int64, sets []int64, errorCounter ErrorCounter) {
	fmt.Println("================\ncumulative stats:")
	fmt.Printf("\ttotal requests: %d\n", len(gets)+len(sets))

	fmt.Printf("\tsuccess: %d\n", len(gets)+len(sets))
	fmt.Printf("\tunavailable: %d\n", errorCounter.unavailable)
	fmt.Printf("\ttimeout exceeded: %d\n", errorCounter.timeout)
	fmt.Printf("\tlimit exceeded: %d\n\n", errorCounter.limitExceeded)

	getHisto := hdrhistogram.New(1, 1000, 1)
	for _, sample := range gets {
		if err := getHisto.RecordValue(sample); err != nil {
			panic(err)
		}
	}
	fmt.Printf(
		"cumulative get latencies:\n\ttotal requests: %d\n\tp50: %d\n\tp90: %d\n\tp99: %d\n\tp99.9: %d\n\tmax: %d\n\n",
		len(gets),
		getHisto.ValueAtQuantile(50.0),
		getHisto.ValueAtQuantile(90.0),
		getHisto.ValueAtQuantile(99.0),
		getHisto.ValueAtQuantile(99.9),
		getHisto.Max(),
	)

	setHisto := hdrhistogram.New(1, 1000, 1)
	for _, sample := range sets {
		if err := setHisto.RecordValue(sample); err != nil {
			panic(err)
		}
	}
	fmt.Printf(
		"cumulative set latencies:\n\ttotal requests: %d\n\tp50: %d\n\tp90: %d\n\tp99: %d\n\tp99.9: %d\n\tmax: %d\n\n",
		len(sets),
		setHisto.ValueAtQuantile(50.0),
		setHisto.ValueAtQuantile(90.0),
		setHisto.ValueAtQuantile(99.0),
		setHisto.ValueAtQuantile(99.9),
		setHisto.Max(),
	)
}

func (ec *ErrorCounter) updateErrors(err string) {
	if err == momento.ServerUnavailableError {
		ec.unavailable++
	} else if err == momento.TimeoutError {
		ec.timeout++
	} else if err == momento.LimitExceededError {
		ec.limitExceeded++
	}
}

func timer(ctx context.Context, getChan chan int64, setChan chan int64, errChan chan string, statsInterval time.Duration) {
	var getMessages []int64
	var setMessages []int64
	errorCounter := ErrorCounter{}

	startTime := hrtime.Now()
	for {

		if hrtime.Since(startTime) >= statsInterval {
			printStats(getMessages, setMessages, errorCounter)
			startTime = hrtime.Now()
		}

		select {
		case <-ctx.Done():
			fmt.Println("\n=====> run complete <====")
			printStats(getMessages, setMessages, errorCounter)
			return
		case getMessage := <-getChan:
			getMessages = append(getMessages, getMessage)
		case setMessage := <-setChan:
			setMessages = append(setMessages, setMessage)
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
	getChan := make(chan int64, r.options.numberOfConcurrentRequests)
	setChan := make(chan int64, r.options.numberOfConcurrentRequests)
	errChan := make(chan string, r.options.numberOfConcurrentRequests)

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
		logLevel:                   logger.DEBUG,
		showStatsInterval:          time.Second * 5,
		cacheItemPayloadBytes:      100,
		numberOfConcurrentRequests: 50,
		maxRequestsPerSecond:       100,
		howLongToRun:               time.Minute,
	}

	loggerFactory := logger.NewBuiltinMomentoLoggerFactory()
	loadGenerator := newLoadGenerator(config.LaptopLatest(loggerFactory), opts)
	client, workerDelayBetweenRequests := loadGenerator.init(ctx)

	runStart := time.Now()
	loadGenerator.run(ctx, client, workerDelayBetweenRequests)
	runTotal := time.Since(runStart)

	fmt.Printf("completed in %f seconds\n", runTotal.Seconds())
}
