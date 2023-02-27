package momento_test

import (
	"errors"
	"time"

	"github.com/momentohq/client-sdk-go/config"
	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
)

var _ = Describe("CacheClient", func() {
	var sharedContext SharedContext

	BeforeEach(func() {
		sharedContext = NewSharedContext()
	})

	It(`errors on an invalid TTL`, func() {
		sharedContext.ClientProps.DefaultTTL = 0 * time.Second
		client, err := NewSimpleCacheClient(sharedContext.ClientProps)

		Expect(client).To(BeNil())
		Expect(err).NotTo(BeNil())
		var momentoErr MomentoError
		if errors.As(err, &momentoErr) {
			Expect(momentoErr.Code()).To(Equal(InvalidArgumentError))
		}
	})

	It(`errors on invalid timeout`, func() {
		badRequestTimeout := 0 * time.Second
		sharedContext.ClientProps.Configuration = config.LatestLaptopConfig().WithClientTimeout(badRequestTimeout)
		Expect(
			NewSimpleCacheClient(sharedContext.ClientProps),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
	})
})
