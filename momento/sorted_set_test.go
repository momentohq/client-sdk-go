package momento_test

import (
	"time"

	"github.com/google/uuid"
	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	"github.com/momentohq/client-sdk-go/utils"
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

	// A convenience for adding elements to a sorted set.
	putElements := func(elements []*SortedSetPutElement) {
		Expect(
			sharedContext.Client.SortedSetPut(
				sharedContext.Ctx,
				&SortedSetPutRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Elements:  elements,
				},
			),
		).To(BeAssignableToTypeOf(&SortedSetPutSuccess{}))
	}

	// Convenience for fetching elements.
	fetch := func() (SortedSetFetchResponse, error) {
		return sharedContext.Client.SortedSetFetch(
			sharedContext.Ctx,
			&SortedSetFetchRequest{
				CacheName: sharedContext.CacheName,
				SetName:   sharedContext.CollectionName,
			},
		)
	}

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

			putElements := []*SortedSetPutElement{{
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

	DescribeTable(`Honors CollectionTTL`,
		func(
			changer func(SortedSetPutElement, *utils.CollectionTTL),
		) {
			value := "foo"
			element := SortedSetPutElement{
				Value: String(value),
				Score: 99,
			}

			expectedFetchHit := &SortedSetFetchHit{
				Elements: []*SortedSetElement{
					{Value: []byte(value), Score: element.Score},
				},
			}

			// It does nothing with no TTL
			putElements([]*SortedSetPutElement{{Value: String(value), Score: 0}})
			changer(element, nil)

			Expect(fetch()).To(Equal(expectedFetchHit))

			time.Sleep(sharedContext.DefaultTTL)

			Expect(fetch()).To(Equal(&SortedSetFetchMiss{}))

			// It does nothing without refresh TTL set.
			putElements([]*SortedSetPutElement{{Value: String(value), Score: 0}})
			changer(
				element,
				&utils.CollectionTTL{
					Ttl:        5 * time.Hour,
					RefreshTtl: false,
				},
			)

			Expect(fetch()).To(Equal(expectedFetchHit))

			time.Sleep(sharedContext.DefaultTTL)

			Expect(fetch()).To(Equal(&SortedSetFetchMiss{}))

			// It does nothing without refresh TTL set.
			putElements([]*SortedSetPutElement{{Value: String(value), Score: 0}})
			changer(
				element,
				&utils.CollectionTTL{
					Ttl:        sharedContext.DefaultTTL + 1*time.Second,
					RefreshTtl: true,
				},
			)

			Expect(fetch()).To(Equal(expectedFetchHit))

			time.Sleep(sharedContext.DefaultTTL)

			Expect(fetch()).To(Equal(expectedFetchHit))

			time.Sleep(1 * time.Second)

			Expect(fetch()).To(Equal(&SortedSetFetchMiss{}))
		},
		Entry(
			`SortedSetIncrementScore`,
			func(element SortedSetPutElement, ttl *utils.CollectionTTL) {
				request := &SortedSetIncrementScoreRequest{
					CacheName:   sharedContext.CacheName,
					SetName:     sharedContext.CollectionName,
					ElementName: element.Value,
					Amount:      element.Score,
				}
				if ttl != nil {
					request.CollectionTTL = *ttl
				}

				sharedContext.Client.SortedSetIncrementScore(sharedContext.Ctx, request)
			},
		),
		Entry(`SortedSetPut`,
			func(element SortedSetPutElement, ttl *utils.CollectionTTL) {
				request := &SortedSetPutRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Elements:  []*SortedSetPutElement{&element},
				}
				if ttl != nil {
					request.CollectionTTL = *ttl
				}

				sharedContext.Client.SortedSetPut(sharedContext.Ctx, request)
			},
		),
	)

	Describe("SortedSetFetch", func() {
		It(`Misses if the set does not exist`, func() {
			Expect(
				sharedContext.Client.SortedSetFetch(
					sharedContext.Ctx,
					&SortedSetFetchRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetFetchMiss{}))
		})

		It(`Fetches`, func() {
			putElements(
				[]*SortedSetPutElement{
					{Value: String("first"), Score: 9999},
					{Value: String("last"), Score: -9999},
					{Value: String("middle"), Score: 50},
				},
			)

			Expect(
				sharedContext.Client.SortedSetFetch(
					sharedContext.Ctx,
					&SortedSetFetchRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
					},
				),
			).To(Equal(
				&SortedSetFetchHit{
					Elements: []*SortedSetElement{
						{Value: []byte("last"), Score: -9999},
						{Value: []byte("middle"), Score: 50},
						{Value: []byte("first"), Score: 9999},
					},
				},
			))

			Expect(
				sharedContext.Client.SortedSetFetch(
					sharedContext.Ctx,
					&SortedSetFetchRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Order:     DESCENDING,
					},
				),
			).To(Equal(
				&SortedSetFetchHit{
					Elements: []*SortedSetElement{
						{Value: []byte("first"), Score: 9999},
						{Value: []byte("middle"), Score: 50},
						{Value: []byte("last"), Score: -9999},
					},
				},
			))

			Expect(
				sharedContext.Client.SortedSetFetch(
					sharedContext.Ctx,
					&SortedSetFetchRequest{
						CacheName:       sharedContext.CacheName,
						SetName:         sharedContext.CollectionName,
						Order:           DESCENDING,
						NumberOfResults: FetchLimitedElements{Limit: 2},
					},
				),
			).To(Equal(
				&SortedSetFetchHit{
					Elements: []*SortedSetElement{
						{Value: []byte("first"), Score: 9999},
						{Value: []byte("middle"), Score: 50},
					},
				},
			))
		})
	})

	Describe(`SortedSetGetRank`, func() {
		It(`Misses when the element does not exist`, func() {
			Expect(
				sharedContext.Client.SortedSetGetRank(
					sharedContext.Ctx,
					&SortedSetGetRankRequest{
						CacheName:   sharedContext.CacheName,
						SetName:     sharedContext.CollectionName,
						ElementName: String("foo"),
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetGetRankMiss{}))
		})

		It(`Gets the rank`, func() {
			putElements(
				[]*SortedSetPutElement{
					{Value: String("first"), Score: 9999},
					{Value: String("last"), Score: -9999},
					{Value: String("middle"), Score: 50},
				},
			)

			Expect(
				sharedContext.Client.SortedSetGetRank(
					sharedContext.Ctx,
					&SortedSetGetRankRequest{
						CacheName:   sharedContext.CacheName,
						SetName:     sharedContext.CollectionName,
						ElementName: String("first"),
					},
				),
			).To(Equal(&SortedSetGetRankHit{Rank: 2}))

			Expect(
				sharedContext.Client.SortedSetGetRank(
					sharedContext.Ctx,
					&SortedSetGetRankRequest{
						CacheName:   sharedContext.CacheName,
						SetName:     sharedContext.CollectionName,
						ElementName: String("last"),
					},
				),
			).To(Equal(&SortedSetGetRankHit{Rank: 0}))
		})
	})

	Describe(`SortedSetGetScore`, func() {
		It(`Misses when the element does not exist`, func() {
			Expect(
				sharedContext.Client.SortedSetGetScore(
					sharedContext.Ctx,
					&SortedSetGetScoreRequest{
						CacheName:    sharedContext.CacheName,
						SetName:      sharedContext.CollectionName,
						ElementNames: []Value{String("foo")},
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetGetScoreMiss{}))
		})

		It(`Gets the score`, func() {
			putElements(
				[]*SortedSetPutElement{
					{Value: String("first"), Score: 9999},
					{Value: String("last"), Score: -9999},
					{Value: String("middle"), Score: 50},
				},
			)

			Expect(
				sharedContext.Client.SortedSetGetScore(
					sharedContext.Ctx,
					&SortedSetGetScoreRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						ElementNames: []Value{
							String("first"), String("last"), String("dne"),
						},
					},
				),
			).To(Equal(
				&SortedSetGetScoreHit{
					Elements: []SortedSetScoreElement{
						&SortedSetScoreHit{Score: 9999},
						&SortedSetScoreHit{Score: -9999},
						&SortedSetScoreMiss{},
					},
				},
			))
		})
	})

	Describe(`SortedSetIncrementScore`, func() {
		It(`Increments if it does not exist`, func() {
			Expect(
				sharedContext.Client.SortedSetIncrementScore(
					sharedContext.Ctx,
					&SortedSetIncrementScoreRequest{
						CacheName:   sharedContext.CacheName,
						SetName:     sharedContext.CollectionName,
						ElementName: String("dne"),
						Amount:      99,
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetIncrementScoreSuccess{Value: 99}))
		})

		It(`Is invalid to increment by 0`, func() {
			Expect(
				sharedContext.Client.SortedSetIncrementScore(
					sharedContext.Ctx,
					&SortedSetIncrementScoreRequest{
						CacheName:   sharedContext.CacheName,
						SetName:     sharedContext.CollectionName,
						ElementName: String("dne"),
						Amount:      0,
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It(`Is invalid to not include the Amount`, func() {
			Expect(
				sharedContext.Client.SortedSetIncrementScore(
					sharedContext.Ctx,
					&SortedSetIncrementScoreRequest{
						CacheName:   sharedContext.CacheName,
						SetName:     sharedContext.CollectionName,
						ElementName: String("dne"),
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It(`Increments the score`, func() {
			putElements(
				[]*SortedSetPutElement{
					{Value: String("first"), Score: 9999},
					{Value: String("last"), Score: -9999},
					{Value: String("middle"), Score: 50},
				},
			)

			Expect(
				sharedContext.Client.SortedSetIncrementScore(
					sharedContext.Ctx,
					&SortedSetIncrementScoreRequest{
						CacheName:   sharedContext.CacheName,
						SetName:     sharedContext.CollectionName,
						ElementName: String("middle"),
						Amount:      42,
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetIncrementScoreSuccess{Value: 92}))

			Expect(
				sharedContext.Client.SortedSetIncrementScore(
					sharedContext.Ctx,
					&SortedSetIncrementScoreRequest{
						CacheName:   sharedContext.CacheName,
						SetName:     sharedContext.CollectionName,
						ElementName: String("middle"),
						Amount:      -42,
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetIncrementScoreSuccess{Value: 50}))
		})
	})

	Describe(`SortedSetRemove`, func() {
		It(`Succeeds when the element does not exist`, func() {
			Expect(
				sharedContext.Client.SortedSetRemove(
					sharedContext.Ctx,
					&SortedSetRemoveRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						ElementsToRemove: RemoveSomeElements{
							Elements: []Value{String("dne")},
						},
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetRemoveSuccess{}))
		})

		It(`Removes elements`, func() {
			putElements(
				[]*SortedSetPutElement{
					{Value: String("first"), Score: 9999},
					{Value: String("last"), Score: -9999},
					{Value: String("middle"), Score: 50},
				},
			)

			Expect(
				sharedContext.Client.SortedSetRemove(
					sharedContext.Ctx,
					&SortedSetRemoveRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						ElementsToRemove: RemoveSomeElements{
							Elements: []Value{
								String("first"), String("dne"),
							},
						},
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetRemoveSuccess{}))

			Expect(
				sharedContext.Client.SortedSetFetch(
					sharedContext.Ctx,
					&SortedSetFetchRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
					},
				),
			).To(Equal(
				&SortedSetFetchHit{
					Elements: []*SortedSetElement{
						{Value: []byte("last"), Score: -9999},
						{Value: []byte("middle"), Score: 50},
					},
				},
			))
		})
	})
})
