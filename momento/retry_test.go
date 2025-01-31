package momento_test

import (
	// . "github.com/momentohq/client-sdk-go/momento"
	// . "github.com/momentohq/client-sdk-go/momento/test_helpers"
	// . "github.com/momentohq/client-sdk-go/responses"
	. "github.com/onsi/ginkgo/v2"
	// . "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
)

var _ = Describe("cache-client retry eligilibity-strategy", func() {
	DescribeTable(
		"DefaultEligibilityStrategy -- determine retry eligibility given grpc status code and request method",
		func(grpcStatus codes.Code, requestMethod string, expected bool) {},
		Entry("name", codes.Internal, "/cache_client.Scs/Get", true),
		Entry("name", codes.Internal, "/cache_client.Scs/Set", true),
		Entry("name", codes.Internal, "/cache_client.Scs/DictionaryIncrement", false),
		Entry("name", codes.Unknown, "/cache_client.Scs/Get", false),
		Entry("name", codes.Unknown, "/cache_client.Scs/Set", false),
		Entry("name", codes.Unknown, "/cache_client.Scs/DictionaryIncrement", false),
		Entry("name", codes.Unavailable, "/cache_client.Scs/Get", true),
		Entry("name", codes.Unavailable, "/cache_client.Scs/Set", true),
		Entry("name", codes.Unavailable, "/cache_client.Scs/DictionaryIncrement", false),
		Entry("name", codes.Canceled, "/cache_client.Scs/Get", true),
		Entry("name", codes.Canceled, "/cache_client.Scs/Set", true),
		Entry("name", codes.Canceled, "/cache_client.Scs/DictionaryIncrement", false),
		Entry("name", codes.DeadlineExceeded, "/cache_client.Scs/Get", false),
		Entry("name", codes.DeadlineExceeded, "/cache_client.Scs/Set", false),
		Entry("name", codes.DeadlineExceeded, "/cache_client.Scs/DictionaryIncrement", false),
	)
})
