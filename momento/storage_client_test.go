package momento_test

import (
	"time"

	"github.com/momentohq/client-sdk-go/config/logger/momento_default_logger"

	"github.com/momentohq/client-sdk-go/config"
	. "github.com/momentohq/client-sdk-go/momento"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("storage-client misc", func() {
	It("errors on invalid timeout", func() {
		badRequestTimeout := 0 * time.Second
		sharedContext.StorageConfiguration = config.StorageLaptopLatest().WithClientTimeout(badRequestTimeout)
		Expect(
			NewPreviewStorageClient(sharedContext.StorageConfiguration, sharedContext.CredentialProvider),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	It("supports constructing a default config with a logger", func() {
		_, err := NewPreviewStorageClient(
			config.StorageLaptopLatestWithLogger(momento_default_logger.NewDefaultMomentoLoggerFactory(momento_default_logger.INFO)),
			sharedContext.CredentialProvider,
		)
		if err != nil {
			panic(err)
		}
	})
})
