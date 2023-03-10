package momento_test

import (
	"errors"
	"time"

	"github.com/momentohq/client-sdk-go/config/logger"

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
			config.LaptopLatestWithLogger(logger.NewBuiltinMomentoLoggerFactory()),
			sharedContext.CredentialProvider,
			sharedContext.DefaultTtl,
		)
		if err != nil {
			panic(err)
		}
	})

	It(`Supports constructing an InRegion config with a logger`, func() {
		_, err := NewCacheClient(
			config.InRegionLatestWithLogger(logger.NewBuiltinMomentoLoggerFactory()),
			sharedContext.CredentialProvider,
			sharedContext.DefaultTtl,
		)
		if err != nil {
			panic(err)
		}
	})
})
