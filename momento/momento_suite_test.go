package momento_test

import (
	"fmt"
	"testing"

	"github.com/momentohq/client-sdk-go/momento"
	helpers "github.com/momentohq/client-sdk-go/momento/test_helpers"
	"github.com/momentohq/client-sdk-go/responses"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var sharedContext helpers.SharedContext
var AUTH_SERVICE_LABEL = "auth-service"
var CACHE_SERVICE_LABEL = "cache-service"
var LEADERBOARD_SERVICE_LABEL = "leaderboard-service"
var STORAGE_SERVICE_LABEL = "storage-service"
var TOPICS_SERVICE_LABEL = "topics-service"
var MOMENTO_LOCAL_LABEL = "momento-local"

func TestMomento(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Momento Suite")
}

var _ = BeforeSuite(func() {
	sharedContext = helpers.NewSharedContext()
	sharedContext.CreateDefaultCaches()

	if includesStorageTests() {
		sharedContext.CreateDefaultStores()
	}
})

var _ = AfterSuite(func() {
	sharedContext.Close()
})

// This assumes that when we narrow tests to a specific service, we are
// doing so with labels (as per the Makefile).
//
// If we want to focus tests based on test regex pattern, we will need to
// update this function to check the test regex pattern instead of labels.
func includesStorageTests() bool {
	// Case 1: No filter is set: run all tests, including storage tests
	// Case 2: Filter is set and it matches storage tests
	return Label("", STORAGE_SERVICE_LABEL).MatchesLabelFilter(GinkgoLabelFilter())
}

func HaveMomentoErrorCode(code string) types.GomegaMatcher {
	return WithTransform(
		func(err error) (string, error) {
			switch mErr := err.(type) {
			case momento.MomentoError:
				return mErr.Code(), nil
			default:
				return "", fmt.Errorf("expected MomentoError, but got %T", err)
			}
		}, Equal(code),
	)
}

func HaveSetLength(length int) types.GomegaMatcher {
	return WithTransform(
		func(fetchResp responses.SetFetchResponse) (int, error) {
			switch rtype := fetchResp.(type) {
			case *responses.SetFetchHit:
				return len(rtype.ValueString()), nil
			default:
				return 0, fmt.Errorf("expected set fetch hit but got %T", fetchResp)
			}
		}, Equal(length),
	)
}

func HaveListLength(length int) types.GomegaMatcher {
	return WithTransform(
		func(fetchResp responses.ListFetchResponse) (int, error) {
			switch rtype := fetchResp.(type) {
			case *responses.ListFetchHit:
				return len(rtype.ValueList()), nil
			default:
				return 0, fmt.Errorf("expected list fetch hit but got %T", fetchResp)
			}
		}, Equal(length),
	)
}

func HaveSortedSetElements(expected []responses.SortedSetBytesElement) types.GomegaMatcher {
	return WithTransform(
		func(fetchResp responses.SortedSetFetchResponse) ([]responses.SortedSetBytesElement, error) {
			switch rtype := fetchResp.(type) {
			case *responses.SortedSetFetchHit:
				return rtype.ValueBytesElements(), nil
			default:
				return nil, fmt.Errorf("expected SortedSetFetchHit, but got %T", fetchResp)
			}
		}, Equal(expected),
	)
}

func HaveSortedSetStringElements(expected []responses.SortedSetStringElement) types.GomegaMatcher {
	return WithTransform(
		func(fetchResp responses.SortedSetFetchResponse) ([]responses.SortedSetStringElement, error) {
			switch rtype := fetchResp.(type) {
			case *responses.SortedSetFetchHit:
				return rtype.ValueStringElements(), nil
			default:
				return nil, fmt.Errorf("expected SortedSetFetchHit, but got %T", fetchResp)
			}
		}, Equal(expected),
	)
}
