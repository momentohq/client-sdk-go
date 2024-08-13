package momento_test

import (
	. "github.com/momentohq/client-sdk-go/responses"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ping-client", func() {
	It("receives successful ping responses", func() {
		for i := 0; i < 25; i++ {
			Expect(
				sharedContext.Client.Ping(sharedContext.Ctx),
			).To(BeAssignableToTypeOf(&PingSuccess{}))
		}
	})

})
