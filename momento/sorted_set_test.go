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
		sharedContext.CreateDefaultCaches()

		DeferCleanup(func() { sharedContext.Close() })
	})

	// A convenience for adding elements to a sorted set.
	putElements := func(elements []SortedSetElement) {
		Expect(
			sharedContext.Client.SortedSetPutElements(
				sharedContext.Ctx,
				&SortedSetPutElementsRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Elements:  elements,
				},
			),
		).To(BeAssignableToTypeOf(&SortedSetPutElementsSuccess{}))
	}

	// Convenience for fetching elements.
	fetch := func() (SortedSetFetchResponse, error) {
		return sharedContext.Client.SortedSetFetchByRank(
			sharedContext.Ctx,
			&SortedSetFetchByRankRequest{
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
				client.SortedSetFetchByRank(ctx, &SortedSetFetchByRankRequest{
					CacheName: cacheName, SetName: collectionName,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			Expect(
				client.SortedSetFetchByScore(ctx, &SortedSetFetchByScoreRequest{
					CacheName: cacheName, SetName: collectionName,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			Expect(
				client.SortedSetGetRank(ctx, &SortedSetGetRankRequest{
					CacheName: cacheName, SetName: collectionName, Value: value,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			Expect(
				client.SortedSetGetScore(ctx, &SortedSetGetScoreRequest{
					CacheName: cacheName, SetName: collectionName, Value: String("hi"),
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

			Expect(
				client.SortedSetPutElement(ctx, &SortedSetPutElementRequest{
					CacheName: cacheName, SetName: collectionName, Value: value, Score: float64(1),
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			putElements := []SortedSetElement{{
				Value: value,
				Score: float64(1),
			}}
			Expect(
				client.SortedSetPutElements(ctx, &SortedSetPutElementsRequest{
					CacheName: cacheName, SetName: collectionName, Elements: putElements,
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			Expect(
				client.SortedSetRemoveElement(ctx, &SortedSetRemoveElementRequest{
					CacheName: cacheName, SetName: collectionName, Value: String("hi"),
				}),
			).Error().To(HaveMomentoErrorCode(expectedError))

			Expect(
				client.SortedSetRemoveElements(ctx, &SortedSetRemoveElementsRequest{
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
			changer func(SortedSetElement, *utils.CollectionTtl),
		) {
			value := "foo"
			element := SortedSetElement{
				Value: String(value),
				Score: 99,
			}

			expectedElements := []SortedSetBytesElement{
				{Value: []byte(value), Score: element.Score},
			}

			// It does nothing with no TTL
			putElements([]SortedSetElement{{Value: String(value), Score: 0}})
			changer(element, nil)

			Expect(fetch()).To(HaveSortedSetElements(expectedElements))

			time.Sleep(sharedContext.DefaultTtl)

			Expect(fetch()).To(Equal(&SortedSetFetchMiss{}))

			// It does nothing without refresh TTL set.
			putElements([]SortedSetElement{{Value: String(value), Score: 0}})
			changer(
				element,
				&utils.CollectionTtl{
					Ttl:        5 * time.Hour,
					RefreshTtl: false,
				},
			)

			Expect(fetch()).To(HaveSortedSetElements(expectedElements))

			time.Sleep(sharedContext.DefaultTtl)

			Expect(fetch()).To(Equal(&SortedSetFetchMiss{}))

			// It does nothing without refresh TTL set.
			putElements([]SortedSetElement{{Value: String(value), Score: 0}})
			changer(
				element,
				&utils.CollectionTtl{
					Ttl:        sharedContext.DefaultTtl + 1*time.Second,
					RefreshTtl: true,
				},
			)

			Expect(fetch()).To(HaveSortedSetElements(expectedElements))

			time.Sleep(sharedContext.DefaultTtl)

			Expect(fetch()).To(HaveSortedSetElements(expectedElements))

			time.Sleep(1 * time.Second)

			Expect(fetch()).To(Equal(&SortedSetFetchMiss{}))
		},
		Entry(
			`SortedSetIncrementScore`,
			func(element SortedSetElement, ttl *utils.CollectionTtl) {
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
		Entry(`SortedSetPutElement`,
			func(element SortedSetElement, ttl *utils.CollectionTtl) {
				request := &SortedSetPutElementRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Value:     element.Value,
					Score:     element.Score,
				}
				if ttl != nil {
					request.Ttl = ttl
				}

				Expect(
					sharedContext.Client.SortedSetPutElement(sharedContext.Ctx, request),
				).To(BeAssignableToTypeOf(&SortedSetPutElementSuccess{}))
			},
		),
		Entry(`SortedSetPutElements`,
			func(element SortedSetElement, ttl *utils.CollectionTtl) {
				request := &SortedSetPutElementsRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Elements:  []SortedSetElement{element},
				}
				if ttl != nil {
					request.Ttl = ttl
				}

				Expect(
					sharedContext.Client.SortedSetPutElements(sharedContext.Ctx, request),
				).To(BeAssignableToTypeOf(&SortedSetPutElementsSuccess{}))
			},
		),
	)

	Describe(`SortedSetFetchByRank`, func() {
		It(`Misses if the set does not exist`, func() {
			Expect(
				sharedContext.Client.SortedSetFetchByRank(
					sharedContext.Ctx,
					&SortedSetFetchByRankRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetFetchMiss{}))
		})

		Context(`With a populated SortedSet`, func() {
			// We'll populate the SortedSet with these elements.
			sortedSetElements := []SortedSetElement{
				{Value: String("one"), Score: 9999},
				{Value: String("two"), Score: 50},
				{Value: String("three"), Score: 0},
				{Value: String("four"), Score: -50},
				{Value: String("five"), Score: -500},
				{Value: String("six"), Score: -1000},
			}

			BeforeEach(func() {
				putElements(sortedSetElements)
			})

			It(`With no extra args it fetches everything in ascending order`, func() {
				resp, err := sharedContext.Client.SortedSetFetchByRank(
					sharedContext.Ctx,
					&SortedSetFetchByRankRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
					},
				)
				Expect(err).To(BeNil())
				Expect(resp).To(HaveSortedSetElements(
					[]SortedSetBytesElement{
						{Value: []byte("six"), Score: -1000},
						{Value: []byte("five"), Score: -500},
						{Value: []byte("four"), Score: -50},
						{Value: []byte("three"), Score: 0},
						{Value: []byte("two"), Score: 50},
						{Value: []byte("one"), Score: 9999},
					},
				))
				Expect(resp).To(HaveSortedSetStringElements(
					[]SortedSetStringElement{
						{Value: "six", Score: -1000},
						{Value: "five", Score: -500},
						{Value: "four", Score: -50},
						{Value: "three", Score: 0},
						{Value: "two", Score: 50},
						{Value: "one", Score: 9999},
					},
				))
			})

			It(`Orders`, func() {
				Expect(
					sharedContext.Client.SortedSetFetchByRank(
						sharedContext.Ctx,
						&SortedSetFetchByRankRequest{
							CacheName: sharedContext.CacheName,
							SetName:   sharedContext.CollectionName,
							Order:     DESCENDING,
						},
					),
				).To(HaveSortedSetElements(
					[]SortedSetBytesElement{
						{Value: []byte("one"), Score: 9999},
						{Value: []byte("two"), Score: 50},
						{Value: []byte("three"), Score: 0},
						{Value: []byte("four"), Score: -50},
						{Value: []byte("five"), Score: -500},
						{Value: []byte("six"), Score: -1000},
					},
				))
			})

			It(`Constrains by start/end rank`, func() {
				start := int32(1)
				end := int32(4)
				Expect(
					sharedContext.Client.SortedSetFetchByRank(
						sharedContext.Ctx,
						&SortedSetFetchByRankRequest{
							CacheName: sharedContext.CacheName,
							SetName:   sharedContext.CollectionName,
							Order:     DESCENDING,
							StartRank: &start,
							EndRank:   &end,
						},
					),
				).To(HaveSortedSetElements(
					[]SortedSetBytesElement{
						{Value: []byte("two"), Score: 50},
						{Value: []byte("three"), Score: 0},
						{Value: []byte("four"), Score: -50},
					},
				))
			})

			It(`Counts negative start rank inclusive from the end`, func() {
				start := int32(-3)
				Expect(
					sharedContext.Client.SortedSetFetchByRank(
						sharedContext.Ctx,
						&SortedSetFetchByRankRequest{
							CacheName: sharedContext.CacheName,
							SetName:   sharedContext.CollectionName,
							Order:     DESCENDING,
							StartRank: &start,
						},
					),
				).To(HaveSortedSetElements(
					[]SortedSetBytesElement{
						{Value: []byte("four"), Score: -50},
						{Value: []byte("five"), Score: -500},
						{Value: []byte("six"), Score: -1000},
					},
				))
			})

			It("returns an empty list when start is after end", func() {
				start := int32(3)
				end := int32(-5)
				fetchResp, err := sharedContext.Client.SortedSetFetchByRank(
					sharedContext.Ctx,
					&SortedSetFetchByRankRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Order:     DESCENDING,
						StartRank: &start,
						EndRank:   &end,
					},
				)
				Expect(err).To(BeNil())
				switch fetchResp := fetchResp.(type) {
				case *SortedSetFetchHit:
					Expect(fetchResp.ValueBytesElements()).To(Equal([]SortedSetBytesElement{}))
					Expect(fetchResp.ValueStringElements()).To(Equal([]SortedSetStringElement{}))
				}
			})

			It(`Counts negative end rank exclusively from the end`, func() {
				end := int32(-3)
				Expect(
					sharedContext.Client.SortedSetFetchByRank(
						sharedContext.Ctx,
						&SortedSetFetchByRankRequest{
							CacheName: sharedContext.CacheName,
							SetName:   sharedContext.CollectionName,
							Order:     DESCENDING,
							EndRank:   &end,
						},
					),
				).To(HaveSortedSetElements(
					[]SortedSetBytesElement{
						{Value: []byte("one"), Score: 9999},
						{Value: []byte("two"), Score: 50},
						{Value: []byte("three"), Score: 0},
					},
				))
			})

			DescribeTable("we get an error for detectable invalid ranges",
				func(start int32, end int32) {
					Expect(
						sharedContext.Client.SortedSetFetchByRank(
							sharedContext.Ctx,
							&SortedSetFetchByRankRequest{
								CacheName: sharedContext.CacheName,
								SetName:   sharedContext.CollectionName,
								StartRank: &start,
								EndRank:   &end,
							},
						),
					).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

					start = int32(-5)
					end = int32(-3)

				},
				Entry("positive values", int32(5), int32(3)),
				Entry("negative values", int32(-5), int32(-8)),
				Entry("equal positives", int32(5), int32(5)),
				Entry("equal negatives", int32(-5), int32(-5)),
			)
		})
	})

	Describe(`SortedSetFetchByScore`, func() {
		It(`Misses if the set does not exist`, func() {
			Expect(
				sharedContext.Client.SortedSetFetchByScore(
					sharedContext.Ctx,
					&SortedSetFetchByScoreRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetFetchMiss{}))
		})

		Context(`With a populated SortedSet`, func() {
			// We'll populate the SortedSet with these elements.
			sortedSetElements := []SortedSetElement{
				{Value: String("one"), Score: 9999},
				{Value: String("two"), Score: 50},
				{Value: String("three"), Score: 0},
				{Value: String("four"), Score: -50},
				{Value: String("five"), Score: -500},
				{Value: String("six"), Score: -1000},
			}

			BeforeEach(func() {
				putElements(sortedSetElements)
			})

			It(`With no extra args it fetches everything in ascending order`, func() {
				resp, err := sharedContext.Client.SortedSetFetchByScore(
					sharedContext.Ctx,
					&SortedSetFetchByScoreRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
					},
				)
				Expect(err).To(BeNil())
				Expect(resp).To(HaveSortedSetElements(
					[]SortedSetBytesElement{
						{Value: []byte("six"), Score: -1000},
						{Value: []byte("five"), Score: -500},
						{Value: []byte("four"), Score: -50},
						{Value: []byte("three"), Score: 0},
						{Value: []byte("two"), Score: 50},
						{Value: []byte("one"), Score: 9999},
					},
				))
				Expect(resp).To(HaveSortedSetStringElements(
					[]SortedSetStringElement{
						{Value: "six", Score: -1000},
						{Value: "five", Score: -500},
						{Value: "four", Score: -50},
						{Value: "three", Score: 0},
						{Value: "two", Score: 50},
						{Value: "one", Score: 9999},
					},
				))
			})

			It(`Orders`, func() {
				Expect(
					sharedContext.Client.SortedSetFetchByScore(
						sharedContext.Ctx,
						&SortedSetFetchByScoreRequest{
							CacheName: sharedContext.CacheName,
							SetName:   sharedContext.CollectionName,
							Order:     DESCENDING,
						},
					),
				).To(HaveSortedSetElements(
					[]SortedSetBytesElement{
						{Value: []byte("one"), Score: 9999},
						{Value: []byte("two"), Score: 50},
						{Value: []byte("three"), Score: 0},
						{Value: []byte("four"), Score: -50},
						{Value: []byte("five"), Score: -500},
						{Value: []byte("six"), Score: -1000},
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
				).To(HaveSortedSetElements(
					[]SortedSetBytesElement{
						{Value: []byte("two"), Score: 50},
						{Value: []byte("three"), Score: 0},
					},
				))
			})

			It(`Limits and offsets`, func() {
				minScore := float64(-750)
				maxScore := float64(51)
				offset := uint32(1)
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
							Offset:    &offset,
							Count:     &count,
						},
					),
				).To(HaveSortedSetElements(
					[]SortedSetBytesElement{
						{Value: []byte("three"), Score: 0},
						{Value: []byte("four"), Score: -50},
					},
				))
			})
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

		It("Misses when the sorted set doesn't exist", func() {
			getResp, err := sharedContext.Client.SortedSetGetRank(
				sharedContext.Ctx, &SortedSetGetRankRequest{
					CacheName: sharedContext.CacheName,
					SetName:   uuid.NewString(),
					Value:     String("idontexist"),
				},
			)
			Expect(err).To(BeNil())
			Expect(getResp).To(BeAssignableToTypeOf(&SortedSetGetRankMiss{}))
		})

		It(`Gets the rank`, func() {
			putElements(
				[]SortedSetElement{
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

			Expect(
				sharedContext.Client.SortedSetGetRank(
					sharedContext.Ctx,
					&SortedSetGetRankRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Order:     DESCENDING,
						Value:     String("last"),
					},
				),
			).To(Equal(SortedSetGetRankHit(2)))
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

		It(`accepts an empty value`, func() {
			putElements([]SortedSetElement{
				{Value: String(""), Score: 0},
			})

			Expect(
				sharedContext.Client.SortedSetGetRank(
					sharedContext.Ctx,
					&SortedSetGetRankRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Value:     String(""),
					},
				),
			).To(Equal(SortedSetGetRankHit(0)))
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

		It("misses when the sorted set doesn't exist", func() {
			getResp, err := sharedContext.Client.SortedSetGetScores(
				sharedContext.Ctx, &SortedSetGetScoresRequest{
					CacheName: sharedContext.CacheName,
					SetName:   uuid.NewString(),
					Values:    []Value{String("idontexist")},
				},
			)
			Expect(err).To(BeNil())
			Expect(getResp).To(BeAssignableToTypeOf(&SortedSetGetScoresMiss{}))
		})

		It(`Gets the correct set of scores`, func() {
			putElements(
				[]SortedSetElement{
					{Value: String("first"), Score: 9999},
					{Value: String("last"), Score: -9999},
					{Value: String("middle"), Score: 50},
				},
			)

			getResp, err := sharedContext.Client.SortedSetGetScores(
				sharedContext.Ctx,
				&SortedSetGetScoresRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Values: []Value{
						String("first"), String("last"), String("dne"),
					},
				},
			)
			Expect(err).To(BeNil())
			switch resp := getResp.(type) {
			case *SortedSetGetScoresHit:
				Expect(resp.Responses()).To(Equal(
					[]SortedSetGetScoreResponse{
						NewSortedSetGetScoreHit(9999),
						NewSortedSetGetScoreHit(-9999),
						&SortedSetGetScoreMiss{},
					},
				))

				Expect(resp.ScoresArray()).To(Equal([]float64{9999, -9999}))

				Expect(resp.ScoresMap()).To(Equal(map[string]float64{"first": 9999, "last": -9999}))
			}
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

	Describe("sorted set get score", func() {
		It("succeeds on the happy path", func() {
			putElements(
				[]SortedSetElement{
					{Value: String("first"), Score: 9999},
					{Value: String("last"), Score: -9999},
					{Value: String("middle"), Score: 50},
				},
			)
			getResp, err := sharedContext.Client.SortedSetGetScore(
				sharedContext.Ctx, &SortedSetGetScoreRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Value:     String("first"),
				},
			)
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *SortedSetGetScoreHit:
				score := result.Score()
				Expect(score).To(Equal(9999.0))
			default:
				Fail("expected a sorted set get score hit but got a miss")
			}
		})

		It("misses when the element doesn't exist", func() {
			putElements([]SortedSetElement{
				{Value: Bytes("value1"), Score: 0},
				{Value: String("value2"), Score: 10},
			})
			getResp, err := sharedContext.Client.SortedSetGetScore(
				sharedContext.Ctx, &SortedSetGetScoreRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Value:     String("idontexist"),
				},
			)
			Expect(err).To(BeNil())
			Expect(getResp).To(BeAssignableToTypeOf(&SortedSetGetScoreMiss{}))
		})

		It("misses when the sorted set doesn't exist", func() {
			getResp, err := sharedContext.Client.SortedSetGetScore(
				sharedContext.Ctx, &SortedSetGetScoreRequest{
					CacheName: sharedContext.CacheName,
					SetName:   uuid.NewString(),
					Value:     String("idontexist"),
				},
			)
			Expect(err).To(BeNil())
			Expect(getResp).To(BeAssignableToTypeOf(&SortedSetGetScoreMiss{}))
		})

		It("returns an error for a nil value", func() {
			Expect(
				sharedContext.Client.SortedSetPutElement(sharedContext.Ctx, &SortedSetPutElementRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Value:     nil,
					Score:     10,
				}),
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
				[]SortedSetElement{
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

		It(`accepts an empty value`, func() {
			putElements([]SortedSetElement{
				{Value: String(""), Score: 50},
			})

			resp, err := sharedContext.Client.SortedSetIncrementScore(
				sharedContext.Ctx,
				&SortedSetIncrementScoreRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Value:     String(""),
					Amount:    10,
				},
			)
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(SortedSetIncrementScoreSuccess(60)))
		})
	})

	Describe("SortedSetPutElement", func() {

		DescribeTable("put an element with each of string and byte values",
			func(inputValue Value, inputScore float64, expected []SortedSetBytesElement) {
				resp, err := sharedContext.Client.SortedSetPutElement(
					sharedContext.Ctx,
					&SortedSetPutElementRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Value:     inputValue,
						Score:     inputScore,
					})
				Expect(err).To(BeNil())
				Expect(resp).To(BeAssignableToTypeOf(&SortedSetPutElementSuccess{}))

				fetchResp, fetchErr := sharedContext.Client.SortedSetFetchByRank(
					sharedContext.Ctx,
					&SortedSetFetchByRankRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
					},
				)
				Expect(fetchErr).To(BeNil())
				switch fetchResp := fetchResp.(type) {
				case *SortedSetFetchHit:
					Expect(fetchResp.ValueBytesElements()).To(Equal(expected))
				}
			},
			Entry("string value", String("aString"), 42.0, []SortedSetBytesElement{{Value: []byte("aString"), Score: 42}}),
			Entry("bytes value", Bytes("aString"), 42.0, []SortedSetBytesElement{{Value: []byte("aString"), Score: 42}}),
		)

		It("overwrites an existing element's score", func() {
			strValue := String("aValue")
			putElements([]SortedSetElement{{Value: strValue, Score: 5}})
			putElements([]SortedSetElement{{Value: strValue, Score: 500}})
			fetchResp, err := sharedContext.Client.SortedSetGetScore(sharedContext.Ctx, &SortedSetGetScoreRequest{
				CacheName: sharedContext.CacheName,
				SetName:   sharedContext.CollectionName,
				Value:     strValue,
			})
			Expect(err).To(BeNil())
			switch fetchResp := fetchResp.(type) {
			case *SortedSetGetScoreHit:
				Expect(fetchResp.Score()).To(Equal(500.0))
			default:
				Fail("expected a hit from sorted set get score but got a miss")
			}
		})

		It("creates the sorted set if it doesn't exist", func() {
			newSetName := uuid.NewString()
			Expect(
				sharedContext.Client.SortedSetPutElement(sharedContext.Ctx, &SortedSetPutElementRequest{
					CacheName: sharedContext.CacheName,
					SetName:   newSetName,
					Value:     String("aValue"),
					Score:     42,
				}),
			).To(BeAssignableToTypeOf(&SortedSetPutElementSuccess{}))

			fetchResp, err := sharedContext.Client.SortedSetGetScore(sharedContext.Ctx, &SortedSetGetScoreRequest{
				CacheName: sharedContext.CacheName,
				SetName:   newSetName,
				Value:     String("aValue"),
			})
			Expect(err).To(BeNil())
			switch fetchResp := fetchResp.(type) {
			case *SortedSetGetScoreHit:
				Expect(fetchResp.Score()).To(Equal(42.0))
			default:
				Fail("expected a hit from sorted set get score but got a miss")
			}
		})

		It("returns an error for a nil value", func() {
			Expect(
				sharedContext.Client.SortedSetPutElement(sharedContext.Ctx, &SortedSetPutElementRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Value:     nil,
					Score:     0,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})
	})

	Describe(`SortedSetPutElements`, func() {
		It("puts multiple elements", func() {
			elems := []SortedSetElement{{Value: String("val1"), Score: 0}, {Value: Bytes("val2"), Score: 10}}
			Expect(
				sharedContext.Client.SortedSetPutElements(
					sharedContext.Ctx,
					&SortedSetPutElementsRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Elements:  elems,
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetPutElementsSuccess{}))

			fetchResp, err := sharedContext.Client.SortedSetFetchByRank(sharedContext.Ctx, &SortedSetFetchByRankRequest{
				CacheName: sharedContext.CacheName,
				SetName:   sharedContext.CollectionName,
				Order:     ASCENDING,
			})
			Expect(err).To(BeNil())
			switch result := fetchResp.(type) {
			case *SortedSetFetchHit:
				Expect(
					result.ValueStringElements(),
				).To(Equal([]SortedSetStringElement{{Value: "val1", Score: 0}, {Value: "val2", Score: 10}}))
				Expect(
					result.ValueBytesElements(),
				).To(Equal([]SortedSetBytesElement{{Value: []byte("val1"), Score: 0}, {Value: []byte("val2"), Score: 10}}))
			default:
				Fail("expected a hit for sorted set fetch but got a miss")
			}
		})

		It("overwrites multiple elements", func() {
			elems := []SortedSetElement{{Value: String("val1"), Score: 0}, {Value: Bytes("val2"), Score: 10}}
			putElements(elems)
			elems = []SortedSetElement{{Value: String("val1"), Score: -999}, {Value: Bytes("val2"), Score: -100}}
			putElements(elems)
			fetchResp, err := fetch()
			Expect(err).To(BeNil())
			switch fetchResp := fetchResp.(type) {
			case *SortedSetFetchHit:
				Expect(fetchResp.ValueBytesElements()).To(Equal(
					[]SortedSetBytesElement{{Value: []byte("val1"), Score: -999}, {Value: []byte("val2"), Score: -100}},
				))
				Expect(fetchResp.ValueStringElements()).To(Equal(
					[]SortedSetStringElement{{Value: "val1", Score: -999}, {Value: "val2", Score: -100}},
				))
			}
		})

		It("returns an error for nil elements", func() {
			Expect(sharedContext.Client.SortedSetPutElements(sharedContext.Ctx, &SortedSetPutElementsRequest{
				CacheName: sharedContext.CacheName,
				SetName:   sharedContext.CollectionName,
				Elements:  nil,
			})).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("returns an error for nil values", func() {
			Expect(sharedContext.Client.SortedSetPutElements(sharedContext.Ctx, &SortedSetPutElementsRequest{
				CacheName: sharedContext.CacheName,
				SetName:   sharedContext.CollectionName,
				Elements:  []SortedSetElement{{Value: String("hi"), Score: 10}, {Value: nil, Score: 500}},
			})).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})
	})

	Describe("Sorted set remove element", func() {
		It("removes an element", func() {
			elems := []SortedSetElement{{Value: String("val1"), Score: 0}, {Value: String("val2"), Score: 10}}
			Expect(
				sharedContext.Client.SortedSetPutElements(
					sharedContext.Ctx,
					&SortedSetPutElementsRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Elements:  elems,
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetPutElementsSuccess{}))

			Expect(
				sharedContext.Client.SortedSetRemoveElement(sharedContext.Ctx, &SortedSetRemoveElementRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Value:     String("val1"),
				}),
			).To(BeAssignableToTypeOf(&SortedSetRemoveElementSuccess{}))

			fetchResp, err := sharedContext.Client.SortedSetFetchByRank(sharedContext.Ctx, &SortedSetFetchByRankRequest{
				CacheName: sharedContext.CacheName,
				SetName:   sharedContext.CollectionName,
				Order:     ASCENDING,
			})
			Expect(err).To(BeNil())
			switch result := fetchResp.(type) {
			case *SortedSetFetchHit:
				Expect(
					result.ValueStringElements(),
				).To(Equal([]SortedSetStringElement{{Value: "val2", Score: 10}}))
			default:
				Fail("expected a hit for sorted set fetch but got a miss")
			}

		})

		It("succeeds when the element doesn't exist", func() {
			Expect(
				sharedContext.Client.SortedSetRemoveElement(sharedContext.Ctx, &SortedSetRemoveElementRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Value:     String("idontexist"),
				}),
			).To(BeAssignableToTypeOf(&SortedSetRemoveElementSuccess{}))
		})

		It("returns an error when value is nil", func() {
			Expect(
				sharedContext.Client.SortedSetRemoveElement(sharedContext.Ctx, &SortedSetRemoveElementRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Value:     nil,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.SortedSetRemoveElement(sharedContext.Ctx, &SortedSetRemoveElementRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})
	})

	Describe(`SortedSetRemoveElements`, func() {
		It(`Succeeds when the element does not exist`, func() {
			Expect(
				sharedContext.Client.SortedSetRemoveElements(
					sharedContext.Ctx,
					&SortedSetRemoveElementsRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Values:    []Value{String("dne")},
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetRemoveElementsSuccess{}))
		})

		It(`Removes elements`, func() {
			putElements(
				[]SortedSetElement{
					{Value: String("first"), Score: 9999},
					{Value: String("last"), Score: -9999},
					{Value: String("middle"), Score: 50},
				},
			)

			Expect(
				sharedContext.Client.SortedSetRemoveElements(
					sharedContext.Ctx,
					&SortedSetRemoveElementsRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Values: []Value{
							String("first"), String("dne"),
						},
					},
				),
			).To(BeAssignableToTypeOf(&SortedSetRemoveElementsSuccess{}))

			Expect(
				sharedContext.Client.SortedSetFetchByRank(
					sharedContext.Ctx,
					&SortedSetFetchByRankRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
					},
				),
			).To(HaveSortedSetElements(
				[]SortedSetBytesElement{
					{Value: []byte("last"), Score: -9999},
					{Value: []byte("middle"), Score: 50},
				},
			))
		})

		It("returns an error when elements are nil", func() {
			Expect(
				sharedContext.Client.SortedSetRemoveElements(
					sharedContext.Ctx,
					&SortedSetRemoveElementsRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Values:    nil,
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.SortedSetRemoveElements(
					sharedContext.Ctx,
					&SortedSetRemoveElementsRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Values:    []Value{nil, String("aValue"), nil},
					},
				),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})
	})
})
