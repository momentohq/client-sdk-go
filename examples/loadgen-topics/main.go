package main

import (
	"context"
	"fmt"
	"math"
	"strconv"
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
	CacheName           = "topics-loadgen"
)

type topicsLoadGeneratorOptions struct {
	logLevel                   momento_default_logger.LogLevel
	showStatsInterval          time.Duration
	messageBytes               int
	numberOfPublishers         int
	numberOfSubscribers        int
	subscriptionsPerSubscriber int
	maxPublishTps              int
	howLongToRun               time.Duration
}

type loadGenerator struct {
	loggerFactory     logger.MomentoLoggerFactory
	logger            logger.MomentoLogger
	topicClientConfig config.TopicsConfiguration
	options           topicsLoadGeneratorOptions
	messageValue      string
}

type ErrorCounter struct {
	unavailable   int64
	timeout       int64
	limitExceeded int64
}

func newLoadGenerator(config config.TopicsConfiguration, options topicsLoadGeneratorOptions) *loadGenerator {
	loggerFactory := config.GetLoggerFactory()
	lgLogger := loggerFactory.GetLogger("topic-loadgen")
	unixMilli := time.Now().UnixMilli()
	messageValue := strings.Repeat("x", options.messageBytes-len(strconv.FormatInt(unixMilli, 10)))
	return &loadGenerator{
		loggerFactory:     loggerFactory,
		logger:            lgLogger,
		topicClientConfig: config,
		options:           options,
		messageValue:      messageValue,
	}
}

func (r *loadGenerator) init(ctx context.Context) momento.TopicClient {
	credentialProvider, err := auth.FromEnvironmentVariable("MOMENTO_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}

	cacheClient, err := momento.NewCacheClient(config.LaptopLatest(), credentialProvider, time.Second*CacheItemTtlSeconds)
	if err != nil {
		panic(err)
	}

	if _, err := cacheClient.CreateCache(ctx, &momento.CreateCacheRequest{CacheName: CacheName}); err != nil {
		panic(err)
	}

	client, err := momento.NewTopicClient(r.topicClientConfig, credentialProvider)
	if err != nil {
		panic(err)
	}

	numberOfConcurrentRequests := r.options.numberOfSubscribers * r.options.subscriptionsPerSubscriber
	r.logger.Debug(
		fmt.Sprintf(
			"Running %d concurrent subscriptions for %d seconds",
			numberOfConcurrentRequests,
			int(r.options.howLongToRun.Seconds())),
	)

	return client
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

func worker(
	ctx context.Context,
	id int,
	subscribeChan chan int64,
	errChan chan string,
	client momento.TopicClient,
	subscriptionsPerSubscriber int,
) {
	for i := 0; i < subscriptionsPerSubscriber; i++ {
		subscription, err := client.Subscribe(ctx, &momento.TopicSubscribeRequest{
			CacheName: CacheName,
			TopicName: fmt.Sprintf("topic-%d", i),
		})
		if err != nil {
			panic(err)
		}
		go func() { pollForMessages(ctx, id, subscription, subscribeChan, errChan) }()
	}
}

func pollForMessages(ctx context.Context, id int, sub momento.TopicSubscription, subscribeChan chan int64, errChan chan string) {
	timestampLength := len(strconv.FormatInt(time.Now().UnixMilli(), 10))
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("--> #%d returning from poller\n", id)
			return
		default:
			//fmt.Printf("%d getting item\n", id)
			//itemStart := hrtime.Now()
			item, err := sub.Item(ctx)
			timestamp, err := strconv.ParseInt(fmt.Sprintf("%v", item)[0:timestampLength], 10, 64)
			elapsed := time.Now().UnixMilli() - timestamp
			//fmt.Printf("%d got item\n", id)
			if err != nil {
				processError(err, errChan)
			} else {
				//subscribeChan <- hrtime.Since(itemStart).Milliseconds()
				subscribeChan <- elapsed
			}
			//fmt.Printf("[%3d] received elapsed: '%d'\n", id, elapsed)
		}
	}
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

func printStats(subscribes *hdrhistogram.Histogram, publishes *hdrhistogram.Histogram, errorCounter ErrorCounter, startTime time.Duration) {
	totalSubscriptionRequests := subscribes.TotalCount() + errorCounter.timeout + errorCounter.unavailable + errorCounter.limitExceeded

	totalTps := int(math.Round(
		float64(totalSubscriptionRequests * 1000 / hrtime.Since(startTime).Milliseconds()),
	))
	fmt.Println("==============================\ncumulative stats:")
	fmt.Printf(
		"%20s: %d (%d tps)\n",
		"total subscription requests",
		totalSubscriptionRequests,
		totalTps,
	)

	successfulSubscribeRequests := subscribes.TotalCount()
	subscribeSuccessTps := int(math.Round(
		float64(successfulSubscribeRequests * 1000 / hrtime.Since(startTime).Milliseconds()),
	))
	// TODO: set up a separate error counter for publishes
	totalPublishRequests := publishes.TotalCount()
	successfulPublishRequests := publishes.TotalCount()
	publishSuccessTps := int(math.Round(
		float64(successfulPublishRequests * 1000 / hrtime.Since(startTime).Milliseconds()),
	))
	fmt.Printf("%20s: %d (%d%%) (%d tps)\n", "subscribe success", successfulSubscribeRequests, successfulSubscribeRequests/totalSubscriptionRequests*100, subscribeSuccessTps)
	fmt.Printf("%20s: %d (%d%%) (%d tps)\n", "publish success", successfulPublishRequests, successfulPublishRequests/totalPublishRequests*100, publishSuccessTps)
	fmt.Printf("%20s: %d (%d%%)\n", "unavailable", errorCounter.unavailable, errorCounter.unavailable/totalSubscriptionRequests*100)
	fmt.Printf("%20s: %d (%d%%)\n", "timeout exceeded", errorCounter.timeout, errorCounter.timeout/totalSubscriptionRequests*100)
	fmt.Printf("%20s: %d (%d%%)\n\n", "limit exceeded", errorCounter.limitExceeded, errorCounter.limitExceeded/totalSubscriptionRequests*100)

	fmt.Printf(
		"cumulative subscription latencies:\n%20s: %d\n%20s: %d\n%20s: %d\n%20s: %d\n%20s: %d\n%20s: %d\n\n",
		"total requests",
		subscribes.TotalCount(),
		"p50",
		subscribes.ValueAtQuantile(50.0),
		"p90",
		subscribes.ValueAtQuantile(90.0),
		"p99",
		subscribes.ValueAtQuantile(99.0),
		"p99.9",
		subscribes.ValueAtQuantile(99.9),
		"max",
		subscribes.Max(),
	)
	fmt.Printf(
		"cumulative publish latencies:\n%20s: %d\n%20s: %d\n%20s: %d\n%20s: %d\n%20s: %d\n%20s: %d\n\n",
		"total requests",
		publishes.TotalCount(),
		"p50",
		publishes.ValueAtQuantile(50.0),
		"p90",
		publishes.ValueAtQuantile(90.0),
		"p99",
		publishes.ValueAtQuantile(99.0),
		"p99.9",
		publishes.ValueAtQuantile(99.9),
		"max",
		publishes.Max(),
	)
}

func timer(
	ctx context.Context, subscribeChan chan int64, publishChan chan int64, errChan chan string, statsInterval time.Duration,
) {
	subscribeHistogram := hdrhistogram.New(1, 5000, 1)
	publishHistogram := hdrhistogram.New(1, 5000, 1)
	errorCounter := ErrorCounter{}

	startTime := hrtime.Now()
	origStartTime := startTime

	for {
		if hrtime.Since(startTime) > statsInterval {
			printStats(subscribeHistogram, publishHistogram, errorCounter, origStartTime)
			startTime = hrtime.Now()
		}

		select {
		case <-ctx.Done():
			fmt.Println("\n=====> run complete <=====")
			printStats(subscribeHistogram, publishHistogram, errorCounter, origStartTime)
			return
		case subscribeMessage := <-subscribeChan:
			if err := subscribeHistogram.RecordValue(subscribeMessage); err != nil {
				panic(err)
			}
		case publishMessage := <-publishChan:
			if err := publishHistogram.RecordValue(publishMessage); err != nil {
				panic(err)
			}
		case errCode := <-errChan:
			errorCounter.updateErrors(errCode)
		default:
			time.Sleep(time.Millisecond * 25)
		}
	}
}

func (r *loadGenerator) run(ctx context.Context, client momento.TopicClient) {
	cancelContext, cancelFunction := context.WithTimeout(ctx, r.options.howLongToRun)
	defer cancelFunction()

	var wg sync.WaitGroup
	subscribeChan := make(chan int64, r.options.numberOfSubscribers)
	publishChan := make(chan int64, r.options.numberOfSubscribers*r.options.numberOfPublishers)
	errChan := make(chan string, r.options.numberOfSubscribers)

	wg.Add(1)
	go func() {
		defer wg.Done()
		timer(cancelContext, subscribeChan, publishChan, errChan, r.options.showStatsInterval)
	}()

	// Launch and run subscriber workers
	for i := 1; i <= r.options.numberOfSubscribers; i++ {
		wg.Add(1)

		// avoid reuse of the same i value in each closure
		i := i

		go func() {
			defer wg.Done()
			worker(
				cancelContext,
				i,
				subscribeChan,
				errChan,
				client,
				r.options.subscriptionsPerSubscriber,
			)
		}()
	}

	for i := 0; i < r.options.numberOfPublishers; i++ {
		fmt.Printf("launching publisher #%d\n", i)
		// Launch publisher worker
		wg.Add(1)

		i := i

		go func() {
			defer wg.Done()
			tpsTimer := hrtime.Now()
			transactions := 0
			for {
				select {
				case <-cancelContext.Done():
					fmt.Printf("returning from publisher #%d\n", i)
					return
				default:
					for j := 0; j < r.options.subscriptionsPerSubscriber; j++ {
						publishStart := hrtime.Now()
						_, err := client.Publish(ctx, &momento.TopicPublishRequest{
							CacheName: CacheName,
							TopicName: fmt.Sprintf("topic-%d", i),
							Value: momento.String(
								fmt.Sprintf("%s%s", strconv.FormatInt(time.Now().UnixMilli(), 10), r.messageValue),
							),
						})
						publishChan <- hrtime.Since(publishStart).Milliseconds()
						// TODO: replace with error counter for publish
						if err != nil {
							panic(err)
						}
						transactions++
						if (transactions >= r.options.maxPublishTps/r.options.numberOfPublishers) && (hrtime.Since(tpsTimer) <= time.Second) {
							time.Sleep(time.Second - hrtime.Since(tpsTimer))
							tpsTimer = hrtime.Now()
							transactions = 0
						}
					}
				}
			}
		}()
	}

	wg.Wait()
}

func main() {
	ctx := context.Background()

	opts := topicsLoadGeneratorOptions{
		logLevel:                   momento_default_logger.DEBUG,
		showStatsInterval:          time.Second * 5,
		messageBytes:               500,
		numberOfPublishers:         5,
		numberOfSubscribers:        40,
		subscriptionsPerSubscriber: 5,
		maxPublishTps:              100,
		howLongToRun:               time.Second * 10,
	}

	maxSubscriptions := uint32(opts.numberOfSubscribers * opts.subscriptionsPerSubscriber)

	lgCfg := config.TopicsDefaultWithLogger(
		//momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.DEBUG),
		logger.NewNoopMomentoLoggerFactory(),
	).WithMaxSubscriptions(maxSubscriptions)

	loadGenerator := newLoadGenerator(lgCfg, opts)
	client := loadGenerator.init(ctx)

	runStart := time.Now()
	loadGenerator.run(ctx, client)
	runTotal := time.Since(runStart)
	fmt.Printf("completed in %f seconds\n", runTotal.Seconds())
	client.Close()
}
