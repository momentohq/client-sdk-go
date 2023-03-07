package momento_test

import (
	"fmt"
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
		return sharedContext.Client.SortedSetFetchByIndex(
			sharedContext.Ctx,
			&SortedSetFetchByIndexRequest{
				CacheName: sharedContext.CacheName,
				SetName:   sharedContext.CollectionName,
			},
		)
	}

	DescribeTable(`Validates the names`,
		func(cacheName string, collectionName string, expectedError string) {
			client := sharedContext.Client
			ctx := sharedContext.Ctx
			value := String(uuid.NewString())

			Expect(
				client.SortedSetFetchByIndex(ctx, &SortedSetFetchByIndexRequest{
					CacheName: cacheName, SetName: collectionName,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			Expect(
				client.SortedSetGetRank(ctx, &SortedSetGetRankRequest{
					CacheName: cacheName, SetName: collectionName, Value: value,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			values := []Value{value}
			Expect(
				client.SortedSetGetScores(ctx, &SortedSetGetScoresRequest{
					CacheName: cacheName, SetName: collectionName, Values: values,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			Expect(
				client.SortedSetIncrementScore(ctx, &SortedSetIncrementScoreRequest{
					CacheName: cacheName, SetName: collectionName, Value: value, Amount: 1,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			putElements := []*SortedSetPutElement{{
				Value: value,
				Score: float64(1),
			}}
			Expect(
				client.SortedSetPut(ctx, &SortedSetPutRequest{
					CacheName: cacheName, SetName: collectionName, Elements: putElements,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			Expect(
				client.SortedSetRemove(ctx, &SortedSetRemoveRequest{
					CacheName: cacheName, SetName: collectionName, Values: values,
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
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Value:     element.Value,
					Amount:    element.Score,
				}
				if ttl != nil {
					request.Ttl = ttl
				}

				Expect(
					sharedContext.Client.SortedSetIncrementScore(sharedContext.Ctx, request),
				).To(BeAssignableToTypeOf(SortedSetIncrementScoreSuccess(0)))
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

	Describe(`SortedSetFetch common tests`, func() {
		Fetchers := []TableEntry{
			Entry(`SortedSetFetchByIndex`, func() (SortedSetFetchResponse, error) {
				return sharedContext.Client.SortedSetFetchByIndex(
					sharedContext.Ctx,
					&SortedSetFetchByIndexRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
					},
				)
			}),
			Entry(`SortedSetFetchByScore`, func() (SortedSetFetchResponse, error) {
				return sharedContext.Client.SortedSetFetchByScore(
					sharedContext.Ctx,
					&SortedSetFetchByScoreRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
					},
				)
			}),
		}

		DescribeTable(
			`Misses if the set does not exist`,
			func(fetcher func() (SortedSetFetchResponse, error)) {
				Expect(
					fetcher(),
				).To(BeAssignableToTypeOf(&SortedSetFetchMiss{}))
			},
			Fetchers,
		)

		DescribeTable(
			`Fetches`,
			func(fetcher func() (SortedSetFetchResponse, error)) {
				putElements(
					[]*SortedSetPutElement{
						{Value: String("one"), Score: 9999},
						{Value: String("two"), Score: 50},
						{Value: String("three"), Score: 0},
						{Value: String("four"), Score: -50},
					},
				)

				Expect(
					fetcher(),
				).To(Equal(
					&SortedSetFetchHit{
						Elements: []*SortedSetElement{
							{Value: []byte("four"), Score: -50},
							{Value: []byte("three"), Score: 0},
							{Value: []byte("two"), Score: 50},
							{Value: []byte("one"), Score: 9999},
						},
					},
				))
			},
			Fetchers,
		)
	})

	Describe("SortedSetFetchByIndex", func() {
		BeforeEach(func() {
			putElements(
				[]*SortedSetPutElement{
					{Value: String("one"), Score: 9999},
					{Value: String("two"), Score: 50},
					{Value: String("three"), Score: 0},
					{Value: String("four"), Score: -50},
				},
			)
		})

		It(`Orders`, func() {
			Expect(
				sharedContext.Client.SortedSetFetchByIndex(
					sharedContext.Ctx,
					&SortedSetFetchByIndexRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Order:     DESCENDING,
					},
				),
			).To(Equal(
				&SortedSetFetchHit{
					Elements: []*SortedSetElement{
						{Value: []byte("one"), Score: 9999},
						{Value: []byte("two"), Score: 50},
						{Value: []byte("three"), Score: 0},
						{Value: []byte("four"), Score: -50},
					},
				},
			))
		})

		It(`Constrains by start/end index`, func() {
			start := int32(1)
			end := int32(3)
			Expect(
				sharedContext.Client.SortedSetFetchByIndex(
					sharedContext.Ctx,
					&SortedSetFetchByIndexRequest{
						CacheName:  sharedContext.CacheName,
						SetName:    sharedContext.CollectionName,
						Order:      DESCENDING,
						StartIndex: &start,
						EndIndex:   &end,
					},
				),
			).To(Equal(
				&SortedSetFetchHit{
					Elements: []*SortedSetElement{
						{Value: []byte("two"), Score: 50},
						{Value: []byte("three"), Score: 0},
					},
				},
			))
		})

		It(`Uses negative start index`, func() {
			start := int32(-3)
			Expect(
				sharedContext.Client.SortedSetFetchByIndex(
					sharedContext.Ctx,
					&SortedSetFetchByIndexRequest{
						CacheName:  sharedContext.CacheName,
						SetName:    sharedContext.CollectionName,
						Order:      DESCENDING,
						StartIndex: &start,
					},
				),
			).To(Equal(
				&SortedSetFetchHit{
					Elements: []*SortedSetElement{
						{Value: []byte("two"), Score: 50},
						{Value: []byte("three"), Score: 0},
						{Value: []byte("four"), Score: -50},
					},
				},
			))
		})

		It(`Uses negative end index`, func() {
			end := int32(-2)
			Expect(
				sharedContext.Client.SortedSetFetchByIndex(
					sharedContext.Ctx,
					&SortedSetFetchByIndexRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Order:     DESCENDING,
						EndIndex:  &end,
					},
				),
			).To(Equal(
				&SortedSetFetchHit{
					Elements: []*SortedSetElement{
						{Value: []byte("one"), Score: 9999},
						{Value: []byte("two"), Score: 50},
					},
				},
			))
		})
	})

	Describe("SortedSetFetchByScore", func() {
		BeforeEach(func() {
			putElements(
				[]*SortedSetPutElement{
					{Value: String("one"), Score: 9999},
					{Value: String("two"), Score: 50},
					{Value: String("three"), Score: 0},
					{Value: String("four"), Score: -50},
					{Value: String("five"), Score: -500},
					{Value: String("six"), Score: -1000},
				},
			)
		})

		It(`orders`, func() {
			Expect(
				sharedContext.Client.SortedSetFetchByScore(
					sharedContext.Ctx,
					&SortedSetFetchByScoreRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Order:     DESCENDING,
					},
				),
			).To(Equal(
				&SortedSetFetchHit{
					Elements: []*SortedSetElement{
						{Value: []byte("one"), Score: 9999},
						{Value: []byte("two"), Score: 50},
						{Value: []byte("three"), Score: 0},
						{Value: []byte("four"), Score: -50},
						{Value: []byte("five"), Score: -500},
						{Value: []byte("six"), Score: -1000},
					},
				},
			))
		})

		It(`Constrains by score inclusive`, func() {
			minScore := float64(0)
			maxScore := float64(50)
			Expect(
				sharedContext.Client.SortedSetFetchByScore(
					sharedContext.Ctx,
					&SortedSetFetchByScoreRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Order:     DESCENDING,
						MinScore:  &minScore,
						MaxScore:  &maxScore,
					},
				),
			).To(Equal(
				&SortedSetFetchHit{
					Elements: []*SortedSetElement{
						{Value: []byte("two"), Score: 50},
						{Value: []byte("three"), Score: 0},
					},
				},
			))
		})

		It(`Limits and offsets`, func() {
			minScore := float64(-750)
			maxScore := float64(51)
			offet := uint32(1)
			count := uint32(2)
			Expect(
				sharedContext.Client.SortedSetFetchByScore(
					sharedContext.Ctx,
					&SortedSetFetchByScoreRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Order:     DESCENDING,
						MinScore:  &minScore,
						MaxScore:  &maxScore,
						Offset:    &offet,
						Count:     &count,
					},
				),
			).To(Equal(
				&SortedSetFetchHit{
					Elements: []*SortedSetElement{
						{Value: []byte("three"), Score: 0},
						{Value: []byte("four"), Score: -50},
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
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Value:     String("foo"),
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

			resp, err := sharedContext.Client.SortedSetGetRank(
				sharedContext.Ctx,
				&SortedSetGetRankRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Value:     String("first"),
				},
			)
			Expect(err).To(BeNil())
			Expect(resp).To(Equal(SortedSetGetRankHit(2)))
			switch r := resp.(type) {
			case SortedSetGetRankHit:
				Expect(r.Rank()).To(Equal(uint64(2)))
			default:
				Fail(fmt.Sprintf("Wrong type: %T", r))
			}

			Expect(
				sharedContext.Client.SortedSetGetRank(
					sharedContext.Ctx,
					&SortedSetGetRankRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Value:     String("last"),
					},
				),
			).To(Equal(SortedSetGetRankHit(0)))
		})

		It(`returns an error for a nil element value`, func() {
			Expect(
				sharedContext.Client.SortedSetGetRank(
					sharedContext.Ctx,
					&SortedSetGetRankRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Value:     nil,
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})
	})

	Describe(`SortedSetGetScores`, func() {
		It(`Misses when the element does not exist`, func() {
			Expect(
				sharedContext.Client.SortedSetGetScores(
					sharedContext.Ctx,
					&SortedSetGetScoresRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Values:    []Value{String("foo")},
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetGetScoresMiss{}))
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
				sharedContext.Client.SortedSetGetScores(
					sharedContext.Ctx,
					&SortedSetGetScoresRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Values: []Value{
							String("first"), String("last"), String("dne"),
						},
					},
				),
			).To(Equal(
				&SortedSetGetScoresHit{
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
				sharedContext.Client.SortedSetGetScores(
					sharedContext.Ctx,
					&SortedSetGetScoresRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Values:    nil,
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.SortedSetGetScores(
					sharedContext.Ctx,
					&SortedSetGetScoresRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Values:    []Value{nil, String("aValue"), nil},
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
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Value:     String("dne"),
						Amount:    99,
					},
				),
			).To(BeAssignableToTypeOf(SortedSetIncrementScoreSuccess(99)))
		})

		It(`Is invalid to increment by 0`, func() {
			Expect(
				sharedContext.Client.SortedSetIncrementScore(
					sharedContext.Ctx,
					&SortedSetIncrementScoreRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Value:     String("dne"),
						Amount:    0,
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It(`Is invalid to not include the Amount`, func() {
			Expect(
				sharedContext.Client.SortedSetIncrementScore(
					sharedContext.Ctx,
					&SortedSetIncrementScoreRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Value:     String("dne"),
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

			resp, err := sharedContext.Client.SortedSetIncrementScore(
				sharedContext.Ctx,
				&SortedSetIncrementScoreRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Value:     String("middle"),
					Amount:    42,
				},
			)
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(SortedSetIncrementScoreSuccess(92)))
			switch r := resp.(type) {
			case SortedSetIncrementScoreSuccess:
				Expect(r.Score()).To(Equal(float64(92)))
			default:
				Fail(fmt.Sprintf("Unexpected response type %T", r))
			}

			Expect(
				sharedContext.Client.SortedSetIncrementScore(
					sharedContext.Ctx,
					&SortedSetIncrementScoreRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Value:     String("middle"),
						Amount:    -42,
					},
				),
			).To(BeAssignableToTypeOf(SortedSetIncrementScoreSuccess(50)))
		})

		It("returns an error when element value is nil", func() {
			Expect(
				sharedContext.Client.SortedSetIncrementScore(
					sharedContext.Ctx,
					&SortedSetIncrementScoreRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Value:     nil,
						Amount:    42,
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
						Values:    []Value{String("dne")},
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
						Values: []Value{
							String("first"), String("dne"),
						},
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetRemoveSuccess{}))

			Expect(
				sharedContext.Client.SortedSetFetchByIndex(
					sharedContext.Ctx,
					&SortedSetFetchByIndexRequest{
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
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Values:    nil,
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.SortedSetRemove(
					sharedContext.Ctx,
					&SortedSetRemoveRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Values:    []Value{nil, String("aValue"), nil},
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})
	})
})
