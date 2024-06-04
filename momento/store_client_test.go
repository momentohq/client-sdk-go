package momento_test

import (
	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"
	"time"

	"github.com/momentohq/client-sdk-go/config"
	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StoreClient", func() {
	var sharedContext = NewSharedContext()

	BeforeEach(func() {
		sharedContext = NewSharedContext()
	})

	It("errors on invalid timeout", func() {
		badRequestTimeout := 0 * time.Second
		sharedContext.StoreConfiguration = config.StoreDefault().WithClientTimeout(badRequestTimeout)
		Expect(
			NewPreviewStoreClient(sharedContext.StoreConfiguration, sharedContext.CredentialProvider),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	It("supports constructing a default config with a logger", func() {
		_, err := NewPreviewStoreClient(
			config.StoreDefaultWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO)),
			sharedContext.CredentialProvider,
		)
		if err != nil {
			panic(err)
		}
	})
})
