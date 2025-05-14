package grpcmanagers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/config/logger"
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"github.com/momentohq/client-sdk-go/internal/models"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	testCtx             context.Context
	subscriptionRequest *pb.XSubscriptionRequest
	grpcConfig          *config.TopicsStaticGrpcConfiguration
	grpcManagerRequest  *models.TopicStreamGrpcManagerRequest
	log                 logger.MomentoLogger
)

var _ = Describe("TopicManager", Label("grpcmanagers"), func() {
	BeforeEach(func() {
		testCtx = context.Background()
		subscriptionRequest = &pb.XSubscriptionRequest{
			CacheName: "cache",
			Topic:     "topic",
		}
		grpcConfig = config.NewTopicsStaticGrpcConfiguration(&config.TopicsGrpcConfigurationProps{})
		credProvider, err := auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
		Expect(err).ToNot(HaveOccurred())
		grpcManagerRequest = &models.TopicStreamGrpcManagerRequest{
			GrpcConfiguration:  grpcConfig,
			CredentialProvider: credProvider,
		}
		log = momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.WARN).GetLogger("grpcmanagers-test")
	})

	Describe("StaticStreamManagerList", func() {
		It("Get one new stream at a time until max concurrent streams reached", func() {
			numGrpcChannels := uint32(2)
			maxConcurrentStreams := numGrpcChannels * config.MAX_CONCURRENT_STREAMS_PER_CHANNEL
			staticList, err := NewStaticStreamManagerList(grpcManagerRequest, numGrpcChannels, log)
			Expect(err).ToNot(HaveOccurred())
			Expect(staticList).NotTo(BeNil())

			// Get one new stream at a time until max concurrent streams reached.
			for i := 0; i < int(maxConcurrentStreams); i++ {
				streamManager, err := staticList.GetNextManager()
				Expect(err).ToNot(HaveOccurred())
				Expect(streamManager).NotTo(BeNil())

				subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(testCtx, subscriptionRequest)
				Expect(subscribeErr).ToNot(HaveOccurred())
				Expect(subscribeClient).NotTo(BeNil())

				// keep the stream alive using a goroutine
				go func() {
					subscribeClient.Recv()
					time.Sleep(1 * time.Second)
				}()
			}

			// Verify all managers are full of active subscriptions
			Expect(staticList.CountNumberOfActiveSubscriptions()).To(Equal(int64(maxConcurrentStreams)))

			// Get one more stream and expect an error.
			stream, err := staticList.GetNextManager()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("LimitExceededError"))
			Expect(stream).To(BeNil())
		})

		It("Starts a burst of streams < max concurrent streams", func() {
			numGrpcChannels := uint32(2)
			maxConcurrentStreams := numGrpcChannels * config.MAX_CONCURRENT_STREAMS_PER_CHANNEL
			staticList, err := NewStaticStreamManagerList(grpcManagerRequest, numGrpcChannels, log)
			Expect(err).ToNot(HaveOccurred())
			Expect(staticList).NotTo(BeNil())

			// Start a burst of streams to occupy half the max concurrent stream capacity.
			waitGroup := sync.WaitGroup{}
			for i := 0; i < int(maxConcurrentStreams/2); i++ {
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()
					streamManager, err := staticList.GetNextManager()
					Expect(err).ToNot(HaveOccurred())
					Expect(streamManager).NotTo(BeNil())

					subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(testCtx, subscriptionRequest)
					Expect(subscribeErr).ToNot(HaveOccurred())
					Expect(subscribeClient).NotTo(BeNil())

					// keep the stream alive using a goroutine
					go func() {
						subscribeClient.Recv()
						time.Sleep(1 * time.Second)
					}()
				}()
			}

			// Wait for the burst to complete.
			waitGroup.Wait()

			// Verify correct number of streams are active.
			Expect(staticList.CountNumberOfActiveSubscriptions()).To(Equal(int64(maxConcurrentStreams / 2)))
		})

		It("Starts a burst of streams == max concurrent streams", func() {
			numGrpcChannels := uint32(2)
			maxConcurrentStreams := numGrpcChannels * config.MAX_CONCURRENT_STREAMS_PER_CHANNEL
			staticList, err := NewStaticStreamManagerList(grpcManagerRequest, numGrpcChannels, log)
			Expect(err).ToNot(HaveOccurred())
			Expect(staticList).NotTo(BeNil())

			// Start a burst of streams to occupy the max concurrent stream capacity.
			waitGroup := sync.WaitGroup{}
			for i := 0; i < int(maxConcurrentStreams); i++ {
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()
					streamManager, err := staticList.GetNextManager()
					Expect(err).ToNot(HaveOccurred())
					Expect(streamManager).NotTo(BeNil())

					subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(testCtx, subscriptionRequest)
					Expect(subscribeErr).ToNot(HaveOccurred())
					Expect(subscribeClient).NotTo(BeNil())

					// keep the stream alive using a goroutine
					go func() {
						subscribeClient.Recv()
						time.Sleep(1 * time.Second)
					}()
				}()
			}

			// Wait for the burst to complete.
			waitGroup.Wait()

			// Verify correct number of streams are active.
			Expect(staticList.CountNumberOfActiveSubscriptions()).To(Equal(int64(maxConcurrentStreams)))
		})

		It("Starts a burst of streams > max concurrent streams", func() {
			numGrpcChannels := uint32(2)
			maxConcurrentStreams := numGrpcChannels * config.MAX_CONCURRENT_STREAMS_PER_CHANNEL
			staticList, err := NewStaticStreamManagerList(grpcManagerRequest, numGrpcChannels, log)
			Expect(err).ToNot(HaveOccurred())
			Expect(staticList).NotTo(BeNil())

			// Start a burst of streams for 10 greater than the max concurrent stream capacity.
			waitGroup := sync.WaitGroup{}
			for i := 0; i < int(maxConcurrentStreams+10); i++ {
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()
					streamManager, err := staticList.GetNextManager()
					// Expecting some errors to occur during the burst
					if err != nil {
						Expect(err.Error()).To(ContainSubstring("LimitExceededError"))
					} else {
						Expect(streamManager).NotTo(BeNil())

						subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(testCtx, subscriptionRequest)
						Expect(subscribeErr).ToNot(HaveOccurred())
						Expect(subscribeClient).NotTo(BeNil())

						// keep the stream alive using a goroutine
						go func() {
							subscribeClient.Recv()
							time.Sleep(1 * time.Second)
						}()
					}
				}()
			}

			// Wait for the burst to complete.
			waitGroup.Wait()

			// Verify correct number of streams are active.
			Expect(staticList.CountNumberOfActiveSubscriptions()).To(Equal(int64(maxConcurrentStreams)))
		})
	})

	Describe("DynamicStreamManagerList", func() {
		It("Get one new stream at a timeuntil max concurrent streams reached", func() {
			numGrpcChannels := uint32(2)
			maxConcurrentStreams := numGrpcChannels * config.MAX_CONCURRENT_STREAMS_PER_CHANNEL
			dynamicList, err := NewDynamicStreamManagerList(grpcManagerRequest, maxConcurrentStreams, log)
			Expect(err).ToNot(HaveOccurred())
			Expect(dynamicList).NotTo(BeNil())

			// Dynamic list always starts with only one grpc manager.
			Expect(len(dynamicList.grpcManagers)).To(Equal(1))

			// Get one new stream at a time until max concurrent streams reached.
			for i := 0; i < int(maxConcurrentStreams); i++ {
				streamManager, err := dynamicList.GetNextManager()
				Expect(err).ToNot(HaveOccurred())
				Expect(streamManager).NotTo(BeNil())

				subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(testCtx, subscriptionRequest)
				Expect(subscribeErr).ToNot(HaveOccurred())
				Expect(subscribeClient).NotTo(BeNil())

				// keep the stream alive using a goroutine
				go func() {
					subscribeClient.Recv()
					time.Sleep(1 * time.Second)
				}()
			}

			// New managers should have been added as needed to support the max number of concurrent streams.
			Expect(len(dynamicList.grpcManagers)).To(Equal(int(numGrpcChannels)))

			// Verify all managers are full of active subscriptions
			Expect(dynamicList.CountNumberOfActiveSubscriptions()).To(Equal(int64(maxConcurrentStreams)))

			// Get one more stream and expect an error.
			stream, err := dynamicList.GetNextManager()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("LimitExceededError"))
			Expect(stream).To(BeNil())
		})

		It("Starts a burst of streams < max concurrent streams", func() {
			numGrpcChannels := uint32(2)
			maxConcurrentStreams := numGrpcChannels * config.MAX_CONCURRENT_STREAMS_PER_CHANNEL
			dynamicList, err := NewDynamicStreamManagerList(grpcManagerRequest, maxConcurrentStreams, log)
			Expect(err).ToNot(HaveOccurred())
			Expect(dynamicList).NotTo(BeNil())

			// Dynamic list always starts with only one grpc manager.
			Expect(len(dynamicList.grpcManagers)).To(Equal(1))

			// Start a burst of streams to occupy half the max concurrent stream capacity.
			waitGroup := sync.WaitGroup{}
			for i := 0; i < int(maxConcurrentStreams/2); i++ {
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()
					streamManager, err := dynamicList.GetNextManager()
					Expect(err).ToNot(HaveOccurred())
					Expect(streamManager).NotTo(BeNil())

					subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(testCtx, subscriptionRequest)
					Expect(subscribeErr).ToNot(HaveOccurred())
					Expect(subscribeClient).NotTo(BeNil())

					// keep the stream alive using a goroutine
					go func() {
						subscribeClient.Recv()
						time.Sleep(1 * time.Second)
					}()
				}()
			}

			// Wait for the burst to complete.
			waitGroup.Wait()

			// No new manager should have been added as we did not exceed a single channel's stream capacity.
			Expect(len(dynamicList.grpcManagers)).To(Equal(1))

			// Verify correct number of streams are active.
			Expect(dynamicList.CountNumberOfActiveSubscriptions()).To(Equal(int64(maxConcurrentStreams / 2)))
		})

		DescribeTable("Starts a burst of streams == max concurrent streams",
			func(numGrpcChannels uint32) {
				maxConcurrentStreams := numGrpcChannels * config.MAX_CONCURRENT_STREAMS_PER_CHANNEL
				dynamicList, err := NewDynamicStreamManagerList(grpcManagerRequest, maxConcurrentStreams, log)
				Expect(err).ToNot(HaveOccurred())
				Expect(dynamicList).NotTo(BeNil())

				// Dynamic list always starts with only one grpc manager.
				Expect(len(dynamicList.grpcManagers)).To(Equal(1))

				// Start a burst of streams to occupy the max concurrent stream capacity.
				waitGroup := sync.WaitGroup{}
				for i := 0; i < int(maxConcurrentStreams); i++ {
					waitGroup.Add(1)
					go func() {
						defer waitGroup.Done()
						streamManager, err := dynamicList.GetNextManager()
						if err != nil {
							fmt.Println("error: ", err)
						}
						Expect(err).ToNot(HaveOccurred())
						Expect(streamManager).NotTo(BeNil())

						subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(testCtx, subscriptionRequest)
						Expect(subscribeErr).ToNot(HaveOccurred())
						Expect(subscribeClient).NotTo(BeNil())

						// keep the stream alive using a goroutine
						go func() {
							subscribeClient.Recv()
							time.Sleep(1 * time.Second)
						}()
					}()
				}

				// Wait for the burst to complete.
				waitGroup.Wait()

				// New managers should have been added as needed to support the max number of concurrent streams.
				Expect(len(dynamicList.grpcManagers)).To(Equal(int(numGrpcChannels)))

				// Verify correct number of streams are active.
				Expect(dynamicList.CountNumberOfActiveSubscriptions()).To(Equal(int64(maxConcurrentStreams)))
			},
			Entry("using max 2 channels", uint32(2)),
			Entry("using max 3 channels", uint32(5)),
			Entry("using max 4 channels", uint32(10)),
			Entry("using max 5 channels", uint32(20)),
			Entry("using max 5 channels", uint32(50)),
		)

		// Try different numbers of grpc channels to fuzz test for deadlocks and other concurrency issues.
		DescribeTable("Starts a burst of streams > max concurrent streams",
			func(numGrpcChannels uint32) {
				maxConcurrentStreams := numGrpcChannels * config.MAX_CONCURRENT_STREAMS_PER_CHANNEL
				dynamicList, err := NewDynamicStreamManagerList(grpcManagerRequest, maxConcurrentStreams, log)
				Expect(err).ToNot(HaveOccurred())
				Expect(dynamicList).NotTo(BeNil())

				// Dynamic list always starts with only one grpc manager.
				Expect(len(dynamicList.grpcManagers)).To(Equal(1))

				// Start a burst of streams to occupy 10 greater than the max concurrent stream capacity.
				waitGroup := sync.WaitGroup{}
				for i := 0; i < int(maxConcurrentStreams+10); i++ {
					waitGroup.Add(1)
					go func() {
						defer waitGroup.Done()
						streamManager, err := dynamicList.GetNextManager()
						// Expecting some errors to occur during the burst
						if err != nil {
							Expect(err.Error()).To(ContainSubstring("LimitExceededError"))
						} else {
							Expect(streamManager).NotTo(BeNil())

							subscribeClient, subscribeErr := streamManager.StreamClient.Subscribe(testCtx, subscriptionRequest)
							Expect(subscribeErr).ToNot(HaveOccurred())
							Expect(subscribeClient).NotTo(BeNil())

							// keep the stream alive using a goroutine
							go func() {
								subscribeClient.Recv()
								time.Sleep(1 * time.Second)
							}()
						}
					}()
				}

				// Wait for the burst to complete.
				waitGroup.Wait()

				// New managers should have been added as needed to support the max number of concurrent streams.
				Expect(len(dynamicList.grpcManagers)).To(Equal(int(numGrpcChannels)))

				// Verify correct number of streams are active.
				Expect(dynamicList.CountNumberOfActiveSubscriptions()).To(Equal(int64(maxConcurrentStreams)))
			},
			Entry("using max 2 channels", uint32(2)),
			Entry("using max 3 channels", uint32(5)),
			Entry("using max 4 channels", uint32(10)),
			Entry("using max 5 channels", uint32(20)),
			Entry("using max 5 channels", uint32(50)),
		)
	})
})
