package momento_test

import (
	"context"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers/topic_manager_lists"
	"github.com/momentohq/client-sdk-go/internal/models"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctx                 context.Context
	subscriptionRequest *pb.XSubscriptionRequest
	grpcConfig          *config.TopicsStaticGrpcConfiguration
	grpcManagerRequest  *models.TopicStreamGrpcManagerRequest
	log                 logger.MomentoLogger
)

// keepStreamAlive consumes heartbeats on its own goroutine until the stream
// ends. It tolerates the shutdown errors produced by pool Close (connection
// closing) and spec ctx cancellation; anything else fails the spec.
func keepStreamAlive(ctx context.Context, subscribeClient pb.Pubsub_SubscribeClient, streams *sync.WaitGroup) {
	streams.Add(1)
	go func() {
		defer GinkgoRecover()
		defer streams.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				item, err := subscribeClient.Recv()
				if err != nil {
					// The test is ending: the pool closed the connection or the
					// spec canceled its context.
					if strings.Contains(err.Error(), "the client connection is closing") ||
						strings.Contains(err.Error(), "context canceled") {
						return
					}
					Expect(err).ToNot(HaveOccurred())
				}
				Expect(item).NotTo(BeNil())
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

var _ = Describe("retry topic-grpc-managers", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
	BeforeEach(func() {
		ctx = context.Background()

		logFactory := momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.WARN)
		log = logFactory.GetLogger("grpcmanagers-test")
		// Same env override as retry_test.go so the port is not hardcoded.
		momentoLocalPort := os.Getenv("MOMENTO_PORT")
		if momentoLocalPort == "" {
			momentoLocalPort = "8080"
		}
		thePort, portErr := strconv.ParseUint(momentoLocalPort, 10, 32)
		Expect(portErr).To(BeNil())
		credProvider, err := auth.NewMomentoLocalProvider(&auth.MomentoLocalConfig{Port: uint(thePort)})
		Expect(err).ToNot(HaveOccurred())

		cacheName := uuid.New().String()
		cacheClient, err := NewCacheClient(config.LaptopLatestWithLogger(logFactory), credProvider, 30*time.Second)
		Expect(err).ToNot(HaveOccurred())
		createResponse, err := cacheClient.CreateCache(ctx, &CreateCacheRequest{
			CacheName: cacheName,
		})
		Expect(err).To(BeNil())
		Expect(createResponse).To(Not(BeNil()))

		subscriptionRequest = &pb.XSubscriptionRequest{
			CacheName: cacheName,
			Topic:     uuid.New().String(),
		}
		grpcConfig = config.NewTopicsStaticGrpcConfiguration(&config.TopicsGrpcConfigurationProps{})
		grpcManagerRequest = &models.TopicStreamGrpcManagerRequest{
			GrpcConfiguration:  grpcConfig,
			CredentialProvider: credProvider,
		}
	})

	// These specs leak reservations on purpose: they fill the pool to assert
	// its accounting. The pool is torn down via Close() at the end of each spec.
	Describe("StaticStreamManagerList", func() {
		It("Get one new stream at a time until max concurrent streams reached", func() {
			numGrpcChannels := uint32(2)
			maxConcurrentStreams := numGrpcChannels * uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
			staticList, err := topic_manager_lists.NewStaticStreamGrpcManagerPool(grpcManagerRequest, numGrpcChannels, log)
			Expect(err).ToNot(HaveOccurred())
			Expect(staticList).NotTo(BeNil())

			// Get one new stream at a time until max concurrent streams reached.
			ctx, cancel := context.WithCancel(ctx)
			waitGroup := sync.WaitGroup{}
			for i := 0; i < int(maxConcurrentStreams); i++ {
				streamManager, err := staticList.GetNextTopicGrpcManager()
				Expect(err).ToNot(HaveOccurred())
				Expect(streamManager).NotTo(BeNil())

				subscribeClient, subscribeErr := streamManager.Manager().StreamClient.Subscribe(ctx, subscriptionRequest)
				Expect(subscribeErr).ToNot(HaveOccurred())
				Expect(subscribeClient).NotTo(BeNil())

				// keep the stream alive until the spec tears down
				keepStreamAlive(ctx, subscribeClient, &waitGroup)
			}
			// Allow time for all streams to be established
			time.Sleep(500 * time.Millisecond)

			// Verify all managers are full of active subscriptions
			Expect(staticList.GetCurrentActiveStreamsCount()).To(Equal(uint64(maxConcurrentStreams)))

			// Get one more stream and expect an error.
			stream, err := staticList.GetNextTopicGrpcManager()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("ClientResourceExhaustedError"))
			Expect(stream).To(BeNil())

			staticList.Close()
			cancel()
			waitGroup.Wait()
		})

		It("Starts a burst of streams < max concurrent streams", func() {
			numGrpcChannels := uint32(2)
			maxConcurrentStreams := numGrpcChannels * uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
			staticList, err := topic_manager_lists.NewStaticStreamGrpcManagerPool(grpcManagerRequest, numGrpcChannels, log)
			Expect(err).ToNot(HaveOccurred())
			Expect(staticList).NotTo(BeNil())

			// Start a burst of streams to occupy just under half the max concurrent stream capacity.
			// Shadow the shared ctx so this spec owns its streams' lifetime and
			// can wait out its keep-alive goroutines before the next spec runs.
			ctx, cancel := context.WithCancel(ctx)
			streamsWaitGroup := sync.WaitGroup{}
			waitGroup := sync.WaitGroup{}
			for i := 0; i < int(maxConcurrentStreams/2-1); i++ {
				waitGroup.Add(1)
				go func() {
					defer GinkgoRecover()
					defer waitGroup.Done()
					streamManager, err := staticList.GetNextTopicGrpcManager()
					Expect(err).ToNot(HaveOccurred())
					Expect(streamManager).NotTo(BeNil())

					subscribeClient, subscribeErr := streamManager.Manager().StreamClient.Subscribe(ctx, subscriptionRequest)
					Expect(subscribeErr).ToNot(HaveOccurred())
					Expect(subscribeClient).NotTo(BeNil())

					// keep the stream alive until the spec tears down
					keepStreamAlive(ctx, subscribeClient, &streamsWaitGroup)
				}()
			}

			// Wait for the burst to complete.
			waitGroup.Wait()

			// Allow time for all streams to be established
			time.Sleep(500 * time.Millisecond)

			// Verify correct number of streams are active.
			Expect(staticList.GetCurrentActiveStreamsCount()).To(Equal(uint64(maxConcurrentStreams/2 - 1)))

			staticList.Close()
			cancel()
			streamsWaitGroup.Wait()
		})

		It("Starts a burst of streams == max concurrent streams", func() {
			numGrpcChannels := uint32(2)
			maxConcurrentStreams := numGrpcChannels * uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
			staticList, err := topic_manager_lists.NewStaticStreamGrpcManagerPool(grpcManagerRequest, numGrpcChannels, log)
			Expect(err).ToNot(HaveOccurred())
			Expect(staticList).NotTo(BeNil())

			// Start a burst of streams to occupy the max concurrent stream capacity.
			// Shadow the shared ctx so this spec owns its streams' lifetime and
			// can wait out its keep-alive goroutines before the next spec runs.
			ctx, cancel := context.WithCancel(ctx)
			streamsWaitGroup := sync.WaitGroup{}
			waitGroup := sync.WaitGroup{}
			for i := 0; i < int(maxConcurrentStreams); i++ {
				waitGroup.Add(1)
				go func() {
					defer GinkgoRecover()
					defer waitGroup.Done()
					streamManager, err := staticList.GetNextTopicGrpcManager()
					Expect(err).ToNot(HaveOccurred())
					Expect(streamManager).NotTo(BeNil())

					subscribeClient, subscribeErr := streamManager.Manager().StreamClient.Subscribe(ctx, subscriptionRequest)
					Expect(subscribeErr).ToNot(HaveOccurred())
					Expect(subscribeClient).NotTo(BeNil())

					// keep the stream alive until the spec tears down
					keepStreamAlive(ctx, subscribeClient, &streamsWaitGroup)
				}()
			}

			// Wait for the burst to complete.
			waitGroup.Wait()

			// Allow time for all streams to be established
			time.Sleep(500 * time.Millisecond)

			// Verify correct number of streams are active.
			Expect(staticList.GetCurrentActiveStreamsCount()).To(Equal(uint64(maxConcurrentStreams)))

			staticList.Close()
			cancel()
			streamsWaitGroup.Wait()
		})

		It("Starts a burst of streams > max concurrent streams", func() {
			numGrpcChannels := uint32(2)
			maxConcurrentStreams := numGrpcChannels * uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
			staticList, err := topic_manager_lists.NewStaticStreamGrpcManagerPool(grpcManagerRequest, numGrpcChannels, log)
			Expect(err).ToNot(HaveOccurred())
			Expect(staticList).NotTo(BeNil())

			// Start a burst of streams for 10 greater than the max concurrent stream capacity.
			// Shadow the shared ctx so this spec owns its streams' lifetime and
			// can wait out its keep-alive goroutines before the next spec runs.
			ctx, cancel := context.WithCancel(ctx)
			streamsWaitGroup := sync.WaitGroup{}
			waitGroup := sync.WaitGroup{}
			for i := 0; i < int(maxConcurrentStreams+10); i++ {
				waitGroup.Add(1)
				go func() {
					defer GinkgoRecover()
					defer waitGroup.Done()

					streamManager, err := staticList.GetNextTopicGrpcManager()
					if err != nil {
						Expect(err.Error()).To(ContainSubstring("ClientResourceExhaustedError"))
					} else {
						Expect(streamManager).NotTo(BeNil())

						subscribeClient, subscribeErr := streamManager.Manager().StreamClient.Subscribe(ctx, subscriptionRequest)
						Expect(subscribeErr).ToNot(HaveOccurred())
						Expect(subscribeClient).NotTo(BeNil())

						// keep the stream alive until the spec tears down
						keepStreamAlive(ctx, subscribeClient, &streamsWaitGroup)
					}
				}()
			}

			// Wait for the burst to complete.
			waitGroup.Wait()

			// Allow time for all streams to be established
			time.Sleep(500 * time.Millisecond)

			// Verify correct number of streams are active.
			Expect(staticList.GetCurrentActiveStreamsCount()).To(Equal(uint64(maxConcurrentStreams)))

			staticList.Close()
			cancel()
			streamsWaitGroup.Wait()
		})
	})

	Describe("DynamicStreamManagerList", func() {
		It("Get one new stream at a time until max concurrent streams reached", func() {
			numGrpcChannels := uint32(2)
			maxConcurrentStreams := numGrpcChannels * uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
			dynamicList, err := topic_manager_lists.NewDynamicStreamGrpcManagerPool(grpcManagerRequest, maxConcurrentStreams, log)
			Expect(err).ToNot(HaveOccurred())
			Expect(dynamicList).NotTo(BeNil())

			// Dynamic list always starts with only one grpc manager.
			Expect(dynamicList.GetCurrentNumberOfGrpcManagers()).To(Equal(1))

			// Get one new stream at a time until max concurrent streams reached.
			ctx, cancel := context.WithCancel(ctx)
			waitGroup := sync.WaitGroup{}
			for i := 0; i < int(maxConcurrentStreams); i++ {
				streamManager, err := dynamicList.GetNextTopicGrpcManager()
				Expect(err).ToNot(HaveOccurred())
				Expect(streamManager).NotTo(BeNil())

				subscribeClient, subscribeErr := streamManager.Manager().StreamClient.Subscribe(ctx, subscriptionRequest)
				Expect(subscribeErr).ToNot(HaveOccurred())
				Expect(subscribeClient).NotTo(BeNil())

				// keep the stream alive until the spec tears down
				keepStreamAlive(ctx, subscribeClient, &waitGroup)
			}
			// Allow time for all streams to be established
			time.Sleep(500 * time.Millisecond)

			// New managers should have been added as needed to support the max number of concurrent streams.
			Expect(dynamicList.GetCurrentNumberOfGrpcManagers()).To(Equal(int(numGrpcChannels)))

			// Verify all managers are full of active subscriptions
			Expect(dynamicList.GetCurrentActiveStreamsCount()).To(Equal(uint64(maxConcurrentStreams)))

			// Get one more stream and expect an error.
			stream, err := dynamicList.GetNextTopicGrpcManager()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("ClientResourceExhaustedError"))
			Expect(stream).To(BeNil())

			dynamicList.Close()
			cancel()
			waitGroup.Wait()
		})

		It("Starts a burst of streams < max concurrent streams", func() {
			numGrpcChannels := uint32(2)
			maxConcurrentStreams := numGrpcChannels * uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
			dynamicList, err := topic_manager_lists.NewDynamicStreamGrpcManagerPool(grpcManagerRequest, maxConcurrentStreams, log)
			Expect(err).ToNot(HaveOccurred())
			Expect(dynamicList).NotTo(BeNil())

			// Dynamic list always starts with only one grpc manager.
			Expect(dynamicList.GetCurrentNumberOfGrpcManagers()).To(Equal(1))

			// Start a burst of streams to occupy just under half the max concurrent stream capacity.
			// Shadow the shared ctx so this spec owns its streams' lifetime and
			// can wait out its keep-alive goroutines before the next spec runs.
			ctx, cancel := context.WithCancel(ctx)
			streamsWaitGroup := sync.WaitGroup{}
			waitGroup := sync.WaitGroup{}
			for i := 0; i < int(maxConcurrentStreams/2-1); i++ {
				waitGroup.Add(1)
				go func() {
					defer GinkgoRecover()
					defer waitGroup.Done()
					streamManager, err := dynamicList.GetNextTopicGrpcManager()
					Expect(err).ToNot(HaveOccurred())
					Expect(streamManager).NotTo(BeNil())

					subscribeClient, subscribeErr := streamManager.Manager().StreamClient.Subscribe(ctx, subscriptionRequest)
					Expect(subscribeErr).ToNot(HaveOccurred())
					Expect(subscribeClient).NotTo(BeNil())

					// keep the stream alive until the spec tears down
					keepStreamAlive(ctx, subscribeClient, &streamsWaitGroup)
				}()
			}

			// Wait for the burst to complete.
			waitGroup.Wait()

			// Allow time for all streams to be established
			time.Sleep(500 * time.Millisecond)

			// No new manager should have been added as we did not exceed a single channel's stream capacity.
			Expect(dynamicList.GetCurrentNumberOfGrpcManagers()).To(Equal(1))

			// Verify correct number of streams are active.
			Expect(dynamicList.GetCurrentActiveStreamsCount()).To(Equal(uint64(maxConcurrentStreams/2 - 1)))

			dynamicList.Close()
			cancel()
			streamsWaitGroup.Wait()
		})

		DescribeTable("Starts a burst of streams == max concurrent streams",
			func(numGrpcChannels uint32) {
				maxConcurrentStreams := numGrpcChannels * uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
				dynamicList, err := topic_manager_lists.NewDynamicStreamGrpcManagerPool(grpcManagerRequest, maxConcurrentStreams, log)
				Expect(err).ToNot(HaveOccurred())
				Expect(dynamicList).NotTo(BeNil())

				// Dynamic list always starts with only one grpc manager.
				Expect(dynamicList.GetCurrentNumberOfGrpcManagers()).To(Equal(1))

				// Start a burst of streams to occupy the max concurrent stream capacity.
				// Shadow the shared ctx so this spec owns its streams' lifetime and
				// can wait out its keep-alive goroutines before the next spec runs.
				ctx, cancel := context.WithCancel(ctx)
				streamsWaitGroup := sync.WaitGroup{}
				waitGroup := sync.WaitGroup{}
				for i := 0; i < int(maxConcurrentStreams); i++ {
					waitGroup.Add(1)
					go func() {
						defer GinkgoRecover()
						defer waitGroup.Done()
						streamManager, err := dynamicList.GetNextTopicGrpcManager()
						Expect(err).ToNot(HaveOccurred())
						Expect(streamManager).NotTo(BeNil())

						subscribeClient, subscribeErr := streamManager.Manager().StreamClient.Subscribe(ctx, subscriptionRequest)
						Expect(subscribeErr).ToNot(HaveOccurred())
						Expect(subscribeClient).NotTo(BeNil())

						// keep the stream alive until the spec tears down
						keepStreamAlive(ctx, subscribeClient, &streamsWaitGroup)
					}()
				}

				// Wait for the burst to complete.
				waitGroup.Wait()

				// Allow time for all streams to be established
				time.Sleep(500 * time.Millisecond)

				// New managers should have been added as needed to support the max number of concurrent streams.
				Expect(dynamicList.GetCurrentNumberOfGrpcManagers()).To(Equal(int(numGrpcChannels)))

				// Verify correct number of streams are active.
				Expect(dynamicList.GetCurrentActiveStreamsCount()).To(Equal(uint64(maxConcurrentStreams)))

				dynamicList.Close()
				cancel()
				streamsWaitGroup.Wait()
			},
			Entry("using max 2 channels", uint32(2)),
			Entry("using max 10 channels", uint32(10)),
			Entry("using max 20 channels", uint32(20)),
		)

		// Try different numbers of grpc channels to fuzz test for deadlocks and other concurrency issues.
		DescribeTable("Starts a burst of streams > max concurrent streams",
			func(numGrpcChannels uint32) {
				maxConcurrentStreams := numGrpcChannels * uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
				dynamicList, err := topic_manager_lists.NewDynamicStreamGrpcManagerPool(grpcManagerRequest, maxConcurrentStreams, log)
				Expect(err).ToNot(HaveOccurred())
				Expect(dynamicList).NotTo(BeNil())

				// Dynamic list always starts with only one grpc manager.
				Expect(dynamicList.GetCurrentNumberOfGrpcManagers()).To(Equal(1))

				// Start a burst of streams to occupy 10 greater than the max concurrent stream capacity.
				// Shadow the shared ctx so this spec owns its streams' lifetime and
				// can wait out its keep-alive goroutines before the next spec runs.
				ctx, cancel := context.WithCancel(ctx)
				streamsWaitGroup := sync.WaitGroup{}
				waitGroup := sync.WaitGroup{}
				for i := 0; i < int(maxConcurrentStreams+10); i++ {
					waitGroup.Add(1)
					go func() {
						defer GinkgoRecover()
						defer waitGroup.Done()

						streamManager, err := dynamicList.GetNextTopicGrpcManager()
						if err != nil {
							Expect(err.Error()).To(ContainSubstring("ClientResourceExhaustedError"))
						} else {
							Expect(streamManager).NotTo(BeNil())

							subscribeClient, subscribeErr := streamManager.Manager().StreamClient.Subscribe(ctx, subscriptionRequest)
							Expect(subscribeErr).ToNot(HaveOccurred())
							Expect(subscribeClient).NotTo(BeNil())

							// keep the stream alive until the spec tears down
							keepStreamAlive(ctx, subscribeClient, &streamsWaitGroup)
						}
					}()
				}

				// Wait for the burst to complete.
				waitGroup.Wait()

				// Allow time for all streams to be established
				time.Sleep(500 * time.Millisecond)

				// New managers should have been added as needed to support the max number of concurrent streams.
				Expect(dynamicList.GetCurrentNumberOfGrpcManagers()).To(Equal(int(numGrpcChannels)))

				// Verify correct number of streams are active.
				Expect(dynamicList.GetCurrentActiveStreamsCount()).To(Equal(uint64(maxConcurrentStreams)))

				dynamicList.Close()
				cancel()
				streamsWaitGroup.Wait()
			},
			Entry("using max 2 channels", uint32(2)),
			Entry("using max 10 channels", uint32(10)),
			Entry("using max 20 channels", uint32(20)),
		)
	})
})
