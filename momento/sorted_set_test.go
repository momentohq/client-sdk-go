package momento_test

import (
	"time"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	. "github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/utils"

	"github.com/google/uuid"
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
		func(cacheName string, collectionName string, expectedError string) {
			client := sharedContext.Client
			ctx := sharedContext.Ctx
			element := String(uuid.NewString())

			Expect(
				client.SortedSetFetch(ctx, &SortedSetFetchRequest{
					CacheName: cacheName, SetName: collectionName,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			Expect(
				client.SortedSetGetRank(ctx, &SortedSetGetRankRequest{
					CacheName: cacheName, SetName: collectionName, ElementValue: element,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			elements := []Value{element}
			Expect(
				client.SortedSetGetScore(ctx, &SortedSetGetScoreRequest{
					CacheName: cacheName, SetName: collectionName, ElementValues: elements,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			Expect(
				client.SortedSetIncrementScore(ctx, &SortedSetIncrementScoreRequest{
					CacheName: cacheName, SetName: collectionName, ElementValue: element, Amount: 1,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			putElements := []*SortedSetPutElement{{
				Value: element,
				Score: float64(1),
			}}
			Expect(
				client.SortedSetPut(ctx, &SortedSetPutRequest{
					CacheName: cacheName, SetName: collectionName, Elements: putElements,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			Expect(
				client.SortedSetRemove(ctx, &SortedSetRemoveRequest{
					CacheName: cacheName, SetName: collectionName, ElementsToRemove: &RemoveSomeElements{Elements: elements},
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))
		},
		Entry("Empty cache name", "", sharedContext.CollectionName, InvalidArgumentError),
		Entry("Blank cache name", "  ", sharedContext.CollectionName, InvalidArgumentError),
		Entry("Empty collection name", sharedContext.CacheName, "", InvalidArgumentError),
		Entry("Blank collection name", sharedContext.CacheName, "  ", InvalidArgumentError),
		Entry("Non-existent cache", uuid.NewString(), uuid.NewString(), NotFoundError),
	)

	DescribeTable(`Honors CollectionTtl  `,
		func(
			changer func(SortedSetPutElement, *utils.CollectionTtl),
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

			time.Sleep(sharedContext.DefaultTtl)

			Expect(fetch()).To(Equal(&SortedSetFetchMiss{}))

			// It does nothing without refresh TTL set.
			putElements([]*SortedSetPutElement{{Value: String(value), Score: 0}})
			changer(
				element,
				&utils.CollectionTtl{
					Ttl:        5 * time.Hour,
					RefreshTtl: false,
				},
			)

			Expect(fetch()).To(Equal(expectedFetchHit))

			time.Sleep(sharedContext.DefaultTtl)

			Expect(fetch()).To(Equal(&SortedSetFetchMiss{}))

			// It does nothing without refresh TTL set.
			putElements([]*SortedSetPutElement{{Value: String(value), Score: 0}})
			changer(
				element,
				&utils.CollectionTtl{
					Ttl:        sharedContext.DefaultTtl + 1*time.Second,
					RefreshTtl: true,
				},
			)

			Expect(fetch()).To(Equal(expectedFetchHit))

			time.Sleep(sharedContext.DefaultTtl)

			Expect(fetch()).To(Equal(expectedFetchHit))

			time.Sleep(1 * time.Second)

			Expect(fetch()).To(Equal(&SortedSetFetchMiss{}))
		},
		Entry(
			`SortedSetIncrementScore`,
			func(element SortedSetPutElement, ttl *utils.CollectionTtl) {
				request := &SortedSetIncrementScoreRequest{
					CacheName:    sharedContext.CacheName,
					SetName:      sharedContext.CollectionName,
					ElementValue: element.Value,
					Amount:       element.Score,
				}
				if ttl != nil {
					request.Ttl = ttl
				}

				Expect(
					sharedContext.Client.SortedSetIncrementScore(sharedContext.Ctx, request),
				).To(BeAssignableToTypeOf(&SortedSetIncrementScoreSuccess{}))
			},
		),
		Entry(`SortedSetPut`,
			func(element SortedSetPutElement, ttl *utils.CollectionTtl) {
				request := &SortedSetPutRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Elements:  []*SortedSetPutElement{&element},
				}
				if ttl != nil {
					request.Ttl = ttl
				}

				Expect(
					sharedContext.Client.SortedSetPut(sharedContext.Ctx, request),
				).To(BeAssignableToTypeOf(&SortedSetPutSuccess{}))
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

			// XXX This test needs to be changed for start/end.
			//
			// Expect(
			// 	sharedContext.Client.SortedSetFetch(
			// 		sharedContext.Ctx,
			// 		&SortedSetFetchRequest{
			// 			CacheName:       sharedContext.CacheName,
			// 			SetName:         sharedContext.CollectionName,
			// 			Order:           DESCENDING,
			// 			NumberOfResults: FetchLimitedElements{Limit: 2},
			// 		},
			// 	),
			// ).To(Equal(
			// 	&SortedSetFetchHit{
			// 		Elements: []*SortedSetElement{
			// 			{Value: []byte("first"), Score: 9999},
			// 			{Value: []byte("middle"), Score: 50},
			// 		},
			// 	},
			// ))
		})
	})

	Describe(`SortedSetGetRank`, func() {
		It(`Misses when the element does not exist`, func() {
			Expect(
				sharedContext.Client.SortedSetGetRank(
					sharedContext.Ctx,
					&SortedSetGetRankRequest{
						CacheName:    sharedContext.CacheName,
						SetName:      sharedContext.CollectionName,
						ElementValue: String("foo"),
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
						CacheName:    sharedContext.CacheName,
						SetName:      sharedContext.CollectionName,
						ElementValue: String("first"),
					},
				),
			).To(Equal(&SortedSetGetRankHit{Rank: 2}))

			Expect(
				sharedContext.Client.SortedSetGetRank(
					sharedContext.Ctx,
					&SortedSetGetRankRequest{
						CacheName:    sharedContext.CacheName,
						SetName:      sharedContext.CollectionName,
						ElementValue: String("last"),
					},
				),
			).To(Equal(&SortedSetGetRankHit{Rank: 0}))
		})

		It(`returns an error for a nil element value`, func() {
			Expect(
				sharedContext.Client.SortedSetGetRank(
					sharedContext.Ctx,
					&SortedSetGetRankRequest{
						CacheName:    sharedContext.CacheName,
						SetName:      sharedContext.CollectionName,
						ElementValue: nil,
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})
	})

	Describe(`SortedSetGetScore`, func() {
		It(`Misses when the element does not exist`, func() {
			Expect(
				sharedContext.Client.SortedSetGetScore(
					sharedContext.Ctx,
					&SortedSetGetScoreRequest{
						CacheName:     sharedContext.CacheName,
						SetName:       sharedContext.CollectionName,
						ElementValues: []Value{String("foo")},
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
						ElementValues: []Value{
							String("first"), String("last"), String("dne"),
						},
					},
				),
			).To(Equal(
				&SortedSetGetScoreHit{
					Elements: []SortedSetScoreElement{
						SortedSetScore(9999),
						SortedSetScore(-9999),
						&SortedSetScoreMiss{},
					},
				},
			))
		})

		It(`returns an error when element values are nil`, func() {
			Expect(
				sharedContext.Client.SortedSetGetScore(
					sharedContext.Ctx,
					&SortedSetGetScoreRequest{
						CacheName:     sharedContext.CacheName,
						SetName:       sharedContext.CollectionName,
						ElementValues: nil,
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.SortedSetGetScore(
					sharedContext.Ctx,
					&SortedSetGetScoreRequest{
						CacheName:     sharedContext.CacheName,
						SetName:       sharedContext.CollectionName,
						ElementValues: []Value{nil, String("aValue"), nil},
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})
	})

	Describe(`SortedSetIncrementScore`, func() {
		It(`Increments if it does not exist`, func() {
			Expect(
				sharedContext.Client.SortedSetIncrementScore(
					sharedContext.Ctx,
					&SortedSetIncrementScoreRequest{
						CacheName:    sharedContext.CacheName,
						SetName:      sharedContext.CollectionName,
						ElementValue: String("dne"),
						Amount:       99,
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetIncrementScoreSuccess{Value: 99}))
		})

		It(`Is invalid to increment by 0`, func() {
			Expect(
				sharedContext.Client.SortedSetIncrementScore(
					sharedContext.Ctx,
					&SortedSetIncrementScoreRequest{
						CacheName:    sharedContext.CacheName,
						SetName:      sharedContext.CollectionName,
						ElementValue: String("dne"),
						Amount:       0,
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It(`Is invalid to not include the Amount`, func() {
			Expect(
				sharedContext.Client.SortedSetIncrementScore(
					sharedContext.Ctx,
					&SortedSetIncrementScoreRequest{
						CacheName:    sharedContext.CacheName,
						SetName:      sharedContext.CollectionName,
						ElementValue: String("dne"),
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
						CacheName:    sharedContext.CacheName,
						SetName:      sharedContext.CollectionName,
						ElementValue: String("middle"),
						Amount:       42,
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetIncrementScoreSuccess{Value: 92}))

			Expect(
				sharedContext.Client.SortedSetIncrementScore(
					sharedContext.Ctx,
					&SortedSetIncrementScoreRequest{
						CacheName:    sharedContext.CacheName,
						SetName:      sharedContext.CollectionName,
						ElementValue: String("middle"),
						Amount:       -42,
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetIncrementScoreSuccess{Value: 50}))
		})

		It("returns an error when element value is nil", func() {
			Expect(
				sharedContext.Client.SortedSetIncrementScore(
					sharedContext.Ctx,
					&SortedSetIncrementScoreRequest{
						CacheName:    sharedContext.CacheName,
						SetName:      sharedContext.CollectionName,
						ElementValue: nil,
						Amount:       42,
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
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

		It("returns an error when elements are nil", func() {
			Expect(
				sharedContext.Client.SortedSetRemove(
					sharedContext.Ctx,
					&SortedSetRemoveRequest{
						CacheName:        sharedContext.CacheName,
						SetName:          sharedContext.CollectionName,
						ElementsToRemove: nil,
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.SortedSetRemove(
					sharedContext.Ctx,
					&SortedSetRemoveRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						ElementsToRemove: RemoveSomeElements{
							Elements: nil,
						},
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.SortedSetRemove(
					sharedContext.Ctx,
					&SortedSetRemoveRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						ElementsToRemove: RemoveSomeElements{
							Elements: []Value{nil, String("aValue"), nil},
						},
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

	})
})
