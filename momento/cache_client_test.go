package momento_test

import (
	"errors"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"

	"github.com/momentohq/client-sdk-go/config"
	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CacheClient", func() {
	var sharedContext SharedContext

	BeforeEach(func() {
		sharedContext = NewSharedContext()
	})

	It(`errors on an invalid TTL`, func() {
		sharedContext.DefaultTtl = 0 * time.Second
		client, err := NewCacheClient(sharedContext.Configuration, sharedContext.CredentialProvider, sharedContext.DefaultTtl)

		Expect(client).To(BeNil())
		Expect(err).NotTo(BeNil())
		var momentoErr MomentoError
		if errors.As(err, &momentoErr) {
			Expect(momentoErr.Code()).To(Equal(InvalidArgumentError))
		}
	})

	It(`errors on invalid timeout`, func() {
		badRequestTimeout := 0 * time.Second
		sharedContext.Configuration = config.LaptopLatest().WithClientTimeout(badRequestTimeout)
		Expect(
			NewCacheClient(sharedContext.Configuration, sharedContext.CredentialProvider, sharedContext.DefaultTtl),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	It(`Supports constructing a laptop config with a logger`, func() {
		_, err := NewCacheClient(
			config.LaptopLatestWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO)),
			sharedContext.CredentialProvider,
			sharedContext.DefaultTtl,
		)
		if err != nil {
			panic(err)
		}
	})

	It(`Supports constructing an InRegion config with a logger`, func() {
		_, err := NewCacheClient(
			config.InRegionLatestWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO)),
			sharedContext.CredentialProvider,
			sharedContext.DefaultTtl,
		)
		if err != nil {
			panic(err)
		}
	})

	It(`Supports constructing a Lambda config with a logger`, func() {
		_, err := NewCacheClient(
			config.LambdaLatestWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO)),
			sharedContext.CredentialProvider,
			sharedContext.DefaultTtl,
		)
		if err != nil {
			panic(err)
		}
	})

	It(`Supports constructing a Lambda config with a logger with eager connections`, func() {
		_, err := NewCacheClientWithEagerConnectTimeout(
			config.LambdaLatestWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO)),
			sharedContext.CredentialProvider,
			sharedContext.DefaultTtl,
			30*time.Second,
		)
		if err != nil {
			panic(err)
		}
	})

	It(`Constructs Lambda config with keepalive turned off`, func() {
		config := config.LambdaLatest()
		grpcConfig := config.GetTransportStrategy().GetGrpcConfig()
		Expect(grpcConfig.GetKeepAlivePermitWithoutCalls()).To(BeFalse())
		Expect(grpcConfig.GetKeepAliveTime()).To(Equal(0 * time.Second))
		Expect(grpcConfig.GetKeepAliveTimeout()).To(Equal(0 * time.Second))
	})

	It(`Constructs InRegionLatest config with keepalive turned on`, func() {
		config := config.InRegionLatest()
		grpcConfig := config.GetTransportStrategy().GetGrpcConfig()
		Expect(grpcConfig.GetKeepAlivePermitWithoutCalls()).To(BeTrue())
		Expect(grpcConfig.GetKeepAliveTime()).To(Equal(5000 * time.Millisecond))
		Expect(grpcConfig.GetKeepAliveTimeout()).To(Equal(1000 * time.Millisecond))
	})

	It(`Constructs LaptopLatest config with keepalive turned on`, func() {
		config := config.LaptopLatest()
		grpcConfig := config.GetTransportStrategy().GetGrpcConfig()
		Expect(grpcConfig.GetKeepAlivePermitWithoutCalls()).To(BeTrue())
		Expect(grpcConfig.GetKeepAliveTime()).To(Equal(5000 * time.Millisecond))
		Expect(grpcConfig.GetKeepAliveTimeout()).To(Equal(1000 * time.Millisecond))
	})

	It(`Returns error when eager connection fails`, func() {
		client, err := NewCacheClientWithEagerConnectTimeout(
			config.LambdaLatestWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO)),
			sharedContext.CredentialProvider,
			sharedContext.DefaultTtl,
			1*time.Millisecond,
		)

		Expect(client).To(BeNil())
		Expect(err).NotTo(BeNil())
		var momentoErr MomentoError
		if errors.As(err, &momentoErr) {
			Expect(momentoErr.Code()).To(Equal(ConnectionError))
		}
	})
})
