package momento_test

import (
	"github.com/google/uuid"
	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SortedSet", func() {
	var sharedContext SharedContext
	BeforeEach(func() {
		sharedContext = NewSharedContext()
		sharedContext.CreateDefaultCache()

		DeferCleanup(func() { sharedContext.Close() })
	})

	DescribeTable(`Validates the names`,
		func(badName string) {
			client := sharedContext.Client
			ctx := sharedContext.Ctx
			cacheName := uuid.NewString()
			collectionName := sharedContext.CollectionName
			element := String(uuid.NewString())

			Expect(
				client.SortedSetFetch(ctx, &SortedSetFetchRequest{
					CacheName: badName, SetName: collectionName,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
			Expect(
				client.SortedSetFetch(ctx, &SortedSetFetchRequest{
					CacheName: cacheName, SetName: badName,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				client.SortedSetGetRank(ctx, &SortedSetGetRankRequest{
					CacheName: badName, SetName: collectionName, ElementName: element,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
			Expect(
				client.SortedSetGetRank(ctx, &SortedSetGetRankRequest{
					CacheName: cacheName, SetName: badName, ElementName: element,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			elements := []Value{element}
			Expect(
				client.SortedSetGetScore(ctx, &SortedSetGetScoreRequest{
					CacheName: badName, SetName: collectionName, ElementNames: elements,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
			Expect(
				client.SortedSetGetScore(ctx, &SortedSetGetScoreRequest{
					CacheName: cacheName, SetName: badName, ElementNames: elements,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				client.SortedSetIncrementScore(ctx, &SortedSetIncrementScoreRequest{
					CacheName: badName, SetName: collectionName, ElementName: element,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
			Expect(
				client.SortedSetIncrementScore(ctx, &SortedSetIncrementScoreRequest{
					CacheName: cacheName, SetName: badName, ElementName: element,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			putElements := []*SortedSetScoreRequestElement{{
				Value: element,
				Score: float64(1),
			}}
			Expect(
				client.SortedSetPut(ctx, &SortedSetPutRequest{
					CacheName: badName, SetName: collectionName, Elements: putElements,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
			Expect(
				client.SortedSetPut(ctx, &SortedSetPutRequest{
					CacheName: cacheName, SetName: badName, Elements: putElements,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				client.SortedSetRemove(ctx, &SortedSetRemoveRequest{
					CacheName: badName, SetName: collectionName, ElementsToRemove: &RemoveSomeElements{elements},
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
			Expect(
				client.SortedSetRemove(ctx, &SortedSetRemoveRequest{
					CacheName: cacheName, SetName: badName, ElementsToRemove: &RemoveSomeElements{elements},
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		},
		Entry("Empty name", ""),
		Entry("Blank name", "  "),
	)
})
