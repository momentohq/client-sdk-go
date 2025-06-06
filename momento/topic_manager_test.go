package momento_test

import (
	"context"
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

var _ = Describe("retry topic-grpc-managers", Label(RETRY_LABEL, MOMENTO_LOCAL_LABEL), func() {
	BeforeEach(func() {
		ctx = context.Background()

		logFactory := momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.WARN)
		log = logFactory.GetLogger("grpcmanagers-test")
		credProvider, err := auth.NewMomentoLocalProvider(&auth.MomentoLocalConfig{Port: uint(8080)})
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

				subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(ctx, subscriptionRequest)
				Expect(subscribeErr).ToNot(HaveOccurred())
				Expect(subscribeClient).NotTo(BeNil())

				// keep the stream alive until end of test using a goroutine
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()
					for {
						select {
						case <-ctx.Done():
							return
						default:
							item, err := subscribeClient.Recv()

							if err != nil {
								// Test is ending if we get canceled error
								if strings.Contains(err.Error(), "the client connection is closing") {
									return
								}
								// Otherwise fail the test
								Expect(err).ToNot(HaveOccurred())
							}

							// Otherwise we expect to receive heartbeats
							Expect(item).NotTo(BeNil())
							time.Sleep(100 * time.Millisecond)
						}
					}
				}()
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
			waitGroup := sync.WaitGroup{}
			for i := 0; i < int(maxConcurrentStreams/2-1); i++ {
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()
					streamManager, err := staticList.GetNextTopicGrpcManager()
					Expect(err).ToNot(HaveOccurred())
					Expect(streamManager).NotTo(BeNil())

					subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(ctx, subscriptionRequest)
					Expect(subscribeErr).ToNot(HaveOccurred())
					Expect(subscribeClient).NotTo(BeNil())

					// keep the stream alive using a goroutine
					go func() {
						for {
							select {
							case <-ctx.Done():
								return
							default:
								item, err := subscribeClient.Recv()

								if err != nil {
									// Test is ending if we get canceled error
									if strings.Contains(err.Error(), "the client connection is closing") {
										return
									}
									// Otherwise fail the test
									Expect(err).ToNot(HaveOccurred())
								}

								// Otherwise we expect to receive heartbeats
								Expect(item).NotTo(BeNil())
								time.Sleep(100 * time.Millisecond)
							}
						}
					}()
				}()
			}

			// Wait for the burst to complete.
			waitGroup.Wait()

			// Allow time for all streams to be established
			time.Sleep(500 * time.Millisecond)

			// Verify correct number of streams are active.
			Expect(staticList.GetCurrentActiveStreamsCount()).To(Equal(uint64(maxConcurrentStreams / 2)))

			staticList.Close()
		})

		It("Starts a burst of streams == max concurrent streams", func() {
			numGrpcChannels := uint32(2)
			maxConcurrentStreams := numGrpcChannels * uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
			staticList, err := topic_manager_lists.NewStaticStreamGrpcManagerPool(grpcManagerRequest, numGrpcChannels, log)
			Expect(err).ToNot(HaveOccurred())
			Expect(staticList).NotTo(BeNil())

			// Start a burst of streams to occupy the max concurrent stream capacity.
			waitGroup := sync.WaitGroup{}
			for i := 0; i < int(maxConcurrentStreams); i++ {
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()
					streamManager, err := staticList.GetNextTopicGrpcManager()
					Expect(err).ToNot(HaveOccurred())
					Expect(streamManager).NotTo(BeNil())

					subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(ctx, subscriptionRequest)
					Expect(subscribeErr).ToNot(HaveOccurred())
					Expect(subscribeClient).NotTo(BeNil())

					// keep the stream alive using a goroutine
					go func() {
						for {
							select {
							case <-ctx.Done():
								return
							default:
								item, err := subscribeClient.Recv()

								if err != nil {
									// Test is ending if we get canceled error
									if strings.Contains(err.Error(), "the client connection is closing") {
										return
									}
									// Otherwise fail the test
									Expect(err).ToNot(HaveOccurred())
								}

								// Otherwise we expect to receive heartbeats
								Expect(item).NotTo(BeNil())
								time.Sleep(100 * time.Millisecond)
							}
						}
					}()
				}()
			}

			// Wait for the burst to complete.
			waitGroup.Wait()

			// Allow time for all streams to be established
			time.Sleep(500 * time.Millisecond)

			// Verify correct number of streams are active.
			Expect(staticList.GetCurrentActiveStreamsCount()).To(Equal(uint64(maxConcurrentStreams)))

			staticList.Close()
		})

		It("Starts a burst of streams > max concurrent streams", func() {
			numGrpcChannels := uint32(2)
			maxConcurrentStreams := numGrpcChannels * uint32(config.MAX_CONCURRENT_STREAMS_PER_CHANNEL)
			staticList, err := topic_manager_lists.NewStaticStreamGrpcManagerPool(grpcManagerRequest, numGrpcChannels, log)
			Expect(err).ToNot(HaveOccurred())
			Expect(staticList).NotTo(BeNil())

			// Start a burst of streams for 10 greater than the max concurrent stream capacity.
			waitGroup := sync.WaitGroup{}
			for i := 0; i < int(maxConcurrentStreams+10); i++ {
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()

					streamManager, err := staticList.GetNextTopicGrpcManager()
					if err != nil {
						Expect(err.Error()).To(ContainSubstring("ClientResourceExhaustedError"))
					} else {
						Expect(streamManager).NotTo(BeNil())

						subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(ctx, subscriptionRequest)
						Expect(subscribeErr).ToNot(HaveOccurred())
						Expect(subscribeClient).NotTo(BeNil())

						// keep the stream alive using a goroutine
						go func() {
							for {
								select {
								case <-ctx.Done():
									return
								default:
									item, err := subscribeClient.Recv()

									if err != nil {
										// Test is ending if we get canceled error
										if strings.Contains(err.Error(), "the client connection is closing") {
											return
										}
										// Otherwise fail the test
										Expect(err).ToNot(HaveOccurred())
									}

									// Otherwise we expect to receive heartbeats
									Expect(item).NotTo(BeNil())
									time.Sleep(100 * time.Millisecond)
								}
							}
						}()
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
		})
	})

	Describe("DynamicStreamManagerList", func() {
		It("Get one new stream at a timeuntil max concurrent streams reached", func() {
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

				subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(ctx, subscriptionRequest)
				Expect(subscribeErr).ToNot(HaveOccurred())
				Expect(subscribeClient).NotTo(BeNil())

				// keep the stream alive using a goroutine
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()
					for {
						select {
						case <-ctx.Done():
							return
						default:
							item, err := subscribeClient.Recv()

							if err != nil {
								// Test is ending if we get canceled error
								if strings.Contains(err.Error(), "the client connection is closing") {
									return
								}
								// Otherwise fail the test
								Expect(err).ToNot(HaveOccurred())
							}

							// Otherwise we expect to receive heartbeats
							Expect(item).NotTo(BeNil())
							time.Sleep(100 * time.Millisecond)
						}
					}
				}()
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
			waitGroup := sync.WaitGroup{}
			for i := 0; i < int(maxConcurrentStreams/2-1); i++ {
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()
					streamManager, err := dynamicList.GetNextTopicGrpcManager()
					Expect(err).ToNot(HaveOccurred())
					Expect(streamManager).NotTo(BeNil())

					subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(ctx, subscriptionRequest)
					Expect(subscribeErr).ToNot(HaveOccurred())
					Expect(subscribeClient).NotTo(BeNil())

					// keep the stream alive using a goroutine
					go func() {
						for {
							select {
							case <-ctx.Done():
								return
							default:
								item, err := subscribeClient.Recv()

								if err != nil {
									// Test is ending if we get canceled error
									if strings.Contains(err.Error(), "the client connection is closing") {
										return
									}
									// Otherwise fail the test
									Expect(err).ToNot(HaveOccurred())
								}

								// Otherwise we expect to receive heartbeats
								Expect(item).NotTo(BeNil())
								time.Sleep(100 * time.Millisecond)
							}
						}
					}()
				}()
			}

			// Wait for the burst to complete.
			waitGroup.Wait()

			// Allow time for all streams to be established
			time.Sleep(500 * time.Millisecond)

			// No new manager should have been added as we did not exceed a single channel's stream capacity.
			Expect(dynamicList.GetCurrentNumberOfGrpcManagers()).To(Equal(1))

			// Verify correct number of streams are active.
			Expect(dynamicList.GetCurrentActiveStreamsCount()).To(Equal(uint64(maxConcurrentStreams / 2)))

			dynamicList.Close()
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
				waitGroup := sync.WaitGroup{}
				for i := 0; i < int(maxConcurrentStreams); i++ {
					waitGroup.Add(1)
					go func() {
						defer waitGroup.Done()
						streamManager, err := dynamicList.GetNextTopicGrpcManager()
						Expect(err).ToNot(HaveOccurred())
						Expect(streamManager).NotTo(BeNil())

						subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(ctx, subscriptionRequest)
						Expect(subscribeErr).ToNot(HaveOccurred())
						Expect(subscribeClient).NotTo(BeNil())

						// keep the stream alive using a goroutine
						go func() {
							for {
								select {
								case <-ctx.Done():
									return
								default:
									item, err := subscribeClient.Recv()

									if err != nil {
										// Test is ending if we get canceled error
										if strings.Contains(err.Error(), "the client connection is closing") {
											return
										}
										// Otherwise fail the test
										Expect(err).ToNot(HaveOccurred())
									}

									// Otherwise we expect to receive heartbeats
									Expect(item).NotTo(BeNil())
									time.Sleep(100 * time.Millisecond)
								}
							}
						}()
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
				waitGroup := sync.WaitGroup{}
				for i := 0; i < int(maxConcurrentStreams+10); i++ {
					waitGroup.Add(1)
					go func() {
						defer waitGroup.Done()

						streamManager, err := dynamicList.GetNextTopicGrpcManager()
						if err != nil {
							Expect(err.Error()).To(ContainSubstring("ClientResourceExhaustedError"))
						} else {
							Expect(streamManager).NotTo(BeNil())

							subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(ctx, subscriptionRequest)
							Expect(subscribeErr).ToNot(HaveOccurred())
							Expect(subscribeClient).NotTo(BeNil())

							// keep the stream alive using a goroutine
							go func() {
								for {
									select {
									case <-ctx.Done():
										return
									default:
										item, err := subscribeClient.Recv()

										if err != nil {
											// Test is ending if we get canceled error
											if strings.Contains(err.Error(), "the client connection is closing") {
												return
											}
											// Otherwise fail the test
											Expect(err).ToNot(HaveOccurred())
										}

										// Otherwise we expect to receive heartbeats
										Expect(item).NotTo(BeNil())
										time.Sleep(100 * time.Millisecond)
									}
								}
							}()
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
			},
			Entry("using max 2 channels", uint32(2)),
			Entry("using max 10 channels", uint32(10)),
			Entry("using max 20 channels", uint32(20)),
		)
	})
})
