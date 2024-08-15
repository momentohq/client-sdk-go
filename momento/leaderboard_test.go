package momento_test

import (
	"fmt"
	"time"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/responses"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("leaderboard-client", func() {
	// Convenience method for creating temporary leaderboard
	createLeaderboard := func() Leaderboard {
		leaderboard, err := sharedContext.LeaderboardClient.Leaderboard(sharedContext.Ctx, &LeaderboardRequest{
			CacheName:       sharedContext.CacheName,
			LeaderboardName: fmt.Sprintf("golang-test-%s", uuid.NewString()),
		})
		if err != nil {
			panic(fmt.Sprintf("Failed to create leaderboard before a test: %v", err))
		}
		return leaderboard
	}

	// Convenience method for deleting temporary leaderboard
	deleteLeaderboard := func(leaderboard Leaderboard) LeaderboardDeleteResponse {
		response, err := leaderboard.Delete(sharedContext.Ctx)
		if err != nil {
			panic(fmt.Sprintf("Failed to delete leaderboard after a test: %v", err))
		}
		return response
	}

	// Convenience method for adding elements to a leaderboard
	upsertElements := func(leaderboard Leaderboard, elements []LeaderboardUpsertElement) LeaderboardUpsertResponse {
		response, err := leaderboard.Upsert(sharedContext.Ctx, LeaderboardUpsertRequest{
			Elements: elements,
		})
		if err != nil {
			panic(fmt.Sprintf("Failed to upsert elements to leaderboard: %v", err))
		}
		return response
	}

	// Convenience method for fetching all elements from a leaderboard
	fetchAllElements := func(leaderboard Leaderboard) []LeaderboardElement {
		response, err := leaderboard.FetchByScore(sharedContext.Ctx, LeaderboardFetchByScoreRequest{})
		if err != nil {
			panic(fmt.Sprintf("Failed to fetch elements from leaderboard: %v", err))
		}
		switch response := response.(type) {
		case *LeaderboardFetchSuccess:
			return response.Values()
		default:
			panic(fmt.Sprintf("Unexpected fetch all elements response type: %T", response))
		}
	}

	Describe("Leaderboard", func() {
		var testLeaderboard Leaderboard
		BeforeEach(func() {
			time.Sleep(100 * time.Millisecond)
			testLeaderboard = createLeaderboard()
			DeferCleanup(func() { deleteLeaderboard(testLeaderboard) })
		})

		It("Validates cache name before creating a leaderboard", func() {
			Expect(
				sharedContext.LeaderboardClient.Leaderboard(sharedContext.Ctx, &LeaderboardRequest{
					CacheName:       "   ",
					LeaderboardName: "test-leaderboard",
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("Validates leaderboard name before creating a leaderboard", func() {
			Expect(
				sharedContext.LeaderboardClient.Leaderboard(sharedContext.Ctx, &LeaderboardRequest{
					CacheName:       "test-cache",
					LeaderboardName: "   ",
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("Creates an empty leaderboard", func() {
			Expect(createLeaderboard()).ToNot(BeNil())
		})
	})

	Describe("Upsert", func() {
		var testLeaderboard Leaderboard
		BeforeEach(func() {
			testLeaderboard = createLeaderboard()
			DeferCleanup(func() { deleteLeaderboard(testLeaderboard) })
		})

		It("Validates number of elements", func() {
			Expect(testLeaderboard.Upsert(sharedContext.Ctx, LeaderboardUpsertRequest{
				Elements: []LeaderboardUpsertElement{},
			})).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("Inserts and updates elements", func() {
			upsert1 := []LeaderboardUpsertElement{
				{Id: 123, Score: 100.0},
				{Id: 456, Score: 200.0},
				{Id: 789, Score: 300.0},
			}
			Expect(upsertElements(testLeaderboard, upsert1)).To(BeAssignableToTypeOf(&LeaderboardUpsertSuccess{}))

			fetch1 := fetchAllElements(testLeaderboard)
			Expect(fetch1).To(HaveLen(3))
			Expect(fetch1).To(ContainElements(
				LeaderboardElement{Id: 123, Score: 100.0, Rank: 0},
				LeaderboardElement{Id: 456, Score: 200.0, Rank: 1},
				LeaderboardElement{Id: 789, Score: 300.0, Rank: 2},
			))

			upsert2 := []LeaderboardUpsertElement{
				{Id: 456, Score: 50.0},
				{Id: 1011, Score: 500.0},
			}
			Expect(upsertElements(testLeaderboard, upsert2)).To(BeAssignableToTypeOf(&LeaderboardUpsertSuccess{}))

			fetch2 := fetchAllElements(testLeaderboard)
			Expect(fetch2).To(HaveLen(4))
			Expect(fetch2).To(ContainElements(
				LeaderboardElement{Id: 456, Score: 50.0, Rank: 0},
				LeaderboardElement{Id: 123, Score: 100.0, Rank: 1},
				LeaderboardElement{Id: 789, Score: 300.0, Rank: 2},
				LeaderboardElement{Id: 1011, Score: 500.0, Rank: 3},
			))
		})

		It("Can work with double precision floats", func() {
			upsert := []LeaderboardUpsertElement{
				{Id: 123, Score: 27.004309862363737},
				{Id: 456, Score: 16777217.0},
				{Id: 789, Score: 300.5},
			}
			Expect(upsertElements(testLeaderboard, upsert)).To(BeAssignableToTypeOf(&LeaderboardUpsertSuccess{}))

			fetchedElements := fetchAllElements(testLeaderboard)
			Expect(fetchedElements).To(HaveLen(3))
			Expect(fetchedElements).To(ContainElements(
				LeaderboardElement{Id: 123, Score: 27.004309862363737, Rank: 0},
				LeaderboardElement{Id: 789, Score: 300.5, Rank: 1},
				LeaderboardElement{Id: 456, Score: 16777217.0, Rank: 2},
			))
		})
	})

	Describe("FetchByScore", func() {
		var testLeaderboard Leaderboard
		BeforeEach(func() {
			testLeaderboard = createLeaderboard()
			DeferCleanup(func() { deleteLeaderboard(testLeaderboard) })
		})

		It("Validates the score range", func() {
			minScore := 100.0
			maxScore := 50.0
			Expect(testLeaderboard.FetchByScore(sharedContext.Ctx, LeaderboardFetchByScoreRequest{
				MinScore: &minScore,
				MaxScore: &maxScore,
			})).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("Returns Success response with no elements when leaderboard is empty", func() {
			response := fetchAllElements(testLeaderboard)
			Expect(response).To(BeEmpty())
		})

		It("Fetches elements by score given different optional arguments", func() {
			upsert := []LeaderboardUpsertElement{
				{Id: 123, Score: 10.0},
				{Id: 234, Score: 100.0},
				{Id: 345, Score: 250.0},
				{Id: 456, Score: 500.0},
				{Id: 567, Score: 750.0},
				{Id: 678, Score: 1000.0},
				{Id: 789, Score: 1500.0},
				{Id: 890, Score: 2000.0},
			}
			Expect(upsertElements(testLeaderboard, upsert)).To(BeAssignableToTypeOf(&LeaderboardUpsertSuccess{}))

			// Fetch using unbounded score range and specifc offset and count
			offset := uint32(2)
			count := uint32(2)
			fetch1, err1 := testLeaderboard.FetchByScore(sharedContext.Ctx, LeaderboardFetchByScoreRequest{
				Offset: &offset,
				Count:  &count,
			})
			Expect(err1).To(BeNil())
			Expect(fetch1).To(BeAssignableToTypeOf(&LeaderboardFetchSuccess{}))
			Expect(fetch1.(*LeaderboardFetchSuccess).Values()).To(HaveLen(2))
			Expect(fetch1.(*LeaderboardFetchSuccess).Values()).To(ContainElements(
				LeaderboardElement{Id: 345, Score: 250.0, Rank: 2},
				LeaderboardElement{Id: 456, Score: 500.0, Rank: 3},
			))

			// Fetch using score range
			minScore := 750.0
			maxScore := 1500.0
			fetch2, err2 := testLeaderboard.FetchByScore(sharedContext.Ctx, LeaderboardFetchByScoreRequest{
				MinScore: &minScore,
				MaxScore: &maxScore,
			})
			Expect(err2).To(BeNil())
			Expect(fetch2).To(BeAssignableToTypeOf(&LeaderboardFetchSuccess{}))
			Expect(fetch2.(*LeaderboardFetchSuccess).Values()).To(HaveLen(2))
			Expect(fetch2.(*LeaderboardFetchSuccess).Values()).To(ContainElements(
				LeaderboardElement{Id: 567, Score: 750.0, Rank: 4},
				LeaderboardElement{Id: 678, Score: 1000.0, Rank: 5},
			))

			// Fetch using all options
			minScore = 0.0
			maxScore = 800.0
			offset = uint32(2)
			count = uint32(10)
			order := DESCENDING
			fetch3, err3 := testLeaderboard.FetchByScore(sharedContext.Ctx, LeaderboardFetchByScoreRequest{
				MinScore: &minScore,
				MaxScore: &maxScore,
				Offset:   &offset,
				Count:    &count,
				Order:    &order,
			})
			Expect(err3).To(BeNil())
			Expect(fetch3).To(BeAssignableToTypeOf(&LeaderboardFetchSuccess{}))
			Expect(fetch3.(*LeaderboardFetchSuccess).Values()).To(HaveLen(3))
			Expect(fetch3.(*LeaderboardFetchSuccess).Values()).To(ContainElements(
				LeaderboardElement{Id: 345, Score: 250.0, Rank: 5},
				LeaderboardElement{Id: 234, Score: 100.0, Rank: 6},
				LeaderboardElement{Id: 123, Score: 10.0, Rank: 7},
			))
		})
	})

	Describe("FetchByRank", func() {
		var testLeaderboard Leaderboard
		BeforeEach(func() {
			testLeaderboard = createLeaderboard()
			DeferCleanup(func() { deleteLeaderboard(testLeaderboard) })
		})

		It("Returns Success response with no elements when leaderboard is empty", func() {
			response, err := testLeaderboard.FetchByRank(sharedContext.Ctx, LeaderboardFetchByRankRequest{
				StartRank: uint32(0),
				EndRank:   uint32(10),
			})
			Expect(err).To(BeNil())
			elements := response.(*LeaderboardFetchSuccess).Values()
			Expect(elements).To(BeEmpty())
		})

		It("Validates the rank range", func() {
			// No range provided
			startRank := uint32(0)
			endRank := uint32(0)
			Expect(testLeaderboard.FetchByRank(sharedContext.Ctx, LeaderboardFetchByRankRequest{
				StartRank: startRank,
				EndRank:   endRank,
			})).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			// Range with start > end
			startRank = uint32(10)
			endRank = uint32(5)
			Expect(testLeaderboard.FetchByRank(sharedContext.Ctx, LeaderboardFetchByRankRequest{
				StartRank: startRank,
				EndRank:   endRank,
			})).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			// Range over the 8192 limit
			startRank = uint32(0)
			endRank = uint32(8193)
			Expect(testLeaderboard.FetchByRank(sharedContext.Ctx, LeaderboardFetchByRankRequest{
				StartRank: startRank,
				EndRank:   endRank,
			})).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("Fetches elements by rank given different arguments", func() {
			upsert := []LeaderboardUpsertElement{
				{Id: 123, Score: 100.0},
				{Id: 456, Score: 200.0},
				{Id: 789, Score: 300.0},
			}
			Expect(upsertElements(testLeaderboard, upsert)).To(BeAssignableToTypeOf(&LeaderboardUpsertSuccess{}))

			// Fetch in ascending order
			ascendingOrder := ASCENDING
			response1, err1 := testLeaderboard.FetchByRank(sharedContext.Ctx, LeaderboardFetchByRankRequest{
				StartRank: uint32(0),
				EndRank:   uint32(10),
				Order:     &ascendingOrder,
			})
			Expect(err1).To(BeNil())
			elements1 := response1.(*LeaderboardFetchSuccess).Values()
			Expect(elements1).To(HaveLen(3))
			Expect(elements1).To(ContainElements(
				LeaderboardElement{Id: 123, Score: 100.0, Rank: 0},
				LeaderboardElement{Id: 456, Score: 200.0, Rank: 1},
				LeaderboardElement{Id: 789, Score: 300.0, Rank: 2},
			))

			// Fetch in descending order
			descendingOrder := DESCENDING
			response2, err2 := testLeaderboard.FetchByRank(sharedContext.Ctx, LeaderboardFetchByRankRequest{
				StartRank: uint32(0),
				EndRank:   uint32(10),
				Order:     &descendingOrder,
			})
			Expect(err2).To(BeNil())
			elements2 := response2.(*LeaderboardFetchSuccess).Values()
			Expect(elements2).To(HaveLen(3))
			Expect(elements2).To(ContainElements(
				LeaderboardElement{Id: 789, Score: 300.0, Rank: 0},
				LeaderboardElement{Id: 456, Score: 200.0, Rank: 1},
				LeaderboardElement{Id: 123, Score: 100.0, Rank: 2},
			))

			// Fetch the top two
			response3, err3 := testLeaderboard.FetchByRank(sharedContext.Ctx, LeaderboardFetchByRankRequest{
				StartRank: uint32(0),
				EndRank:   uint32(2),
			})
			Expect(err3).To(BeNil())
			elements3 := response3.(*LeaderboardFetchSuccess).Values()
			Expect(elements3).To(HaveLen(2))
			Expect(elements3).To(ContainElements(
				LeaderboardElement{Id: 123, Score: 100.0, Rank: 0},
				LeaderboardElement{Id: 456, Score: 200.0, Rank: 1},
			))
		})
	})

	Describe("GetRank", func() {
		var testLeaderboard Leaderboard
		BeforeEach(func() {
			testLeaderboard = createLeaderboard()
			DeferCleanup(func() { deleteLeaderboard(testLeaderboard) })
		})

		It("Returns Success response with no elements when leaderboard is empty", func() {
			response, err := testLeaderboard.GetRank(sharedContext.Ctx, LeaderboardGetRankRequest{
				Ids: []uint32{123, 456, 789},
			})
			Expect(err).To(BeNil())
			Expect(response).To(BeAssignableToTypeOf(&LeaderboardFetchSuccess{}))
			Expect(response.(*LeaderboardFetchSuccess).Values()).To(BeEmpty())
		})

		It("Fetches elements given list of IDs", func() {
			upsert := []LeaderboardUpsertElement{
				{Id: 123, Score: 100.0},
				{Id: 456, Score: 200.0},
				{Id: 789, Score: 300.0},
			}
			Expect(upsertElements(testLeaderboard, upsert)).To(BeAssignableToTypeOf(&LeaderboardUpsertSuccess{}))

			// Fetch given empty list of IDs
			response1, err1 := testLeaderboard.GetRank(sharedContext.Ctx, LeaderboardGetRankRequest{
				Ids: []uint32{},
			})
			Expect(err1).To(BeNil())
			Expect(response1).To(BeAssignableToTypeOf(&LeaderboardFetchSuccess{}))
			Expect(response1.(*LeaderboardFetchSuccess).Values()).To(BeEmpty())

			// Fetch elements with ascending order
			ascendingOrder := ASCENDING
			response2, err2 := testLeaderboard.GetRank(sharedContext.Ctx, LeaderboardGetRankRequest{
				Ids:   []uint32{123, 456, 789},
				Order: &ascendingOrder,
			})
			Expect(err2).To(BeNil())
			Expect(response2).To(BeAssignableToTypeOf(&LeaderboardFetchSuccess{}))
			Expect(response2.(*LeaderboardFetchSuccess).Values()).To(HaveLen(3))
			Expect(response2.(*LeaderboardFetchSuccess).Values()).To(ContainElements(
				LeaderboardElement{Id: 123, Score: 100.0, Rank: 0},
				LeaderboardElement{Id: 456, Score: 200.0, Rank: 1},
				LeaderboardElement{Id: 789, Score: 300.0, Rank: 2},
			))

			// Fetch elements with descending order
			descendingOrder := DESCENDING
			response3, err3 := testLeaderboard.GetRank(sharedContext.Ctx, LeaderboardGetRankRequest{
				Ids:   []uint32{123, 456, 789},
				Order: &descendingOrder,
			})
			Expect(err3).To(BeNil())
			Expect(response3).To(BeAssignableToTypeOf(&LeaderboardFetchSuccess{}))
			Expect(response3.(*LeaderboardFetchSuccess).Values()).To(HaveLen(3))
			Expect(response3.(*LeaderboardFetchSuccess).Values()).To(ContainElements(
				LeaderboardElement{Id: 789, Score: 300.0, Rank: 0},
				LeaderboardElement{Id: 456, Score: 200.0, Rank: 1},
				LeaderboardElement{Id: 123, Score: 100.0, Rank: 2},
			))
		})
	})

	Describe("Length", func() {
		var testLeaderboard Leaderboard
		BeforeEach(func() {
			testLeaderboard = createLeaderboard()
			DeferCleanup(func() { deleteLeaderboard(testLeaderboard) })
		})

		It("Returns Success response with zero length when leaderboard is empty", func() {
			response, err := testLeaderboard.Length(sharedContext.Ctx)
			Expect(err).To(BeNil())
			Expect(response).To(BeAssignableToTypeOf(&LeaderboardLengthSuccess{}))
			Expect(response.(*LeaderboardLengthSuccess).Length()).To(BeZero())
		})

		It("Returns Success response with nonzero length when leaderboard has elements", func() {
			upsert := []LeaderboardUpsertElement{
				{Id: 123, Score: 100.0},
				{Id: 456, Score: 200.0},
				{Id: 789, Score: 300.0},
			}
			Expect(upsertElements(testLeaderboard, upsert)).To(BeAssignableToTypeOf(&LeaderboardUpsertSuccess{}))

			response, err := testLeaderboard.Length(sharedContext.Ctx)
			Expect(err).To(BeNil())
			Expect(response).To(BeAssignableToTypeOf(&LeaderboardLengthSuccess{}))
			Expect(response.(*LeaderboardLengthSuccess).Length()).To(Equal(uint32(3)))
		})
	})

	Describe("RemoveElements", func() {
		var testLeaderboard Leaderboard
		BeforeEach(func() {
			testLeaderboard = createLeaderboard()
			DeferCleanup(func() { deleteLeaderboard(testLeaderboard) })
		})

		It("Returns Successful no-op response when leaderboard is empty", func() {
			response, err := testLeaderboard.RemoveElements(sharedContext.Ctx, LeaderboardRemoveElementsRequest{
				Ids: []uint32{123, 456, 789},
			})
			Expect(err).To(BeNil())
			Expect(response).To(BeAssignableToTypeOf(&LeaderboardRemoveElementsSuccess{}))
		})

		It("Validates number of arguments", func() {
			Expect(testLeaderboard.RemoveElements(sharedContext.Ctx, LeaderboardRemoveElementsRequest{
				Ids: []uint32{},
			})).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("Removes elements when leaderboard has elements", func() {
			upsert := []LeaderboardUpsertElement{
				{Id: 123, Score: 100.0},
				{Id: 456, Score: 200.0},
				{Id: 789, Score: 300.0},
			}
			Expect(upsertElements(testLeaderboard, upsert)).To(BeAssignableToTypeOf(&LeaderboardUpsertSuccess{}))

			response, err := testLeaderboard.RemoveElements(sharedContext.Ctx, LeaderboardRemoveElementsRequest{
				Ids: []uint32{123, 789},
			})
			Expect(err).To(BeNil())
			Expect(response).To(BeAssignableToTypeOf(&LeaderboardRemoveElementsSuccess{}))

			// Check remaining elements are as expected
			fetchedElements := fetchAllElements(testLeaderboard)
			Expect(fetchedElements).To(HaveLen(1))
			Expect(fetchedElements).To(ContainElements(
				LeaderboardElement{Id: 456, Score: 200.0, Rank: 0},
			))
		})
	})

	Describe("Delete", func() {
		var testLeaderboard Leaderboard
		BeforeEach(func() {
			testLeaderboard = createLeaderboard()
			DeferCleanup(func() { deleteLeaderboard(testLeaderboard) })
		})

		It("Returns Successful no-op response when leaderboard is empty", func() {
			response, err := testLeaderboard.Delete(sharedContext.Ctx)
			Expect(err).To(BeNil())
			Expect(response).To(BeAssignableToTypeOf(&LeaderboardDeleteSuccess{}))
		})

		It("Deletes non-empty leaderboard", func() {
			upsert := []LeaderboardUpsertElement{
				{Id: 123, Score: 100.0},
				{Id: 456, Score: 200.0},
				{Id: 789, Score: 300.0},
			}
			Expect(upsertElements(testLeaderboard, upsert)).To(BeAssignableToTypeOf(&LeaderboardUpsertSuccess{}))

			response, err := testLeaderboard.Delete(sharedContext.Ctx)
			Expect(err).To(BeNil())
			Expect(response).To(BeAssignableToTypeOf(&LeaderboardDeleteSuccess{}))

			// Check remaining elements are as expected
			fetchedElements := fetchAllElements(testLeaderboard)
			Expect(fetchedElements).To(HaveLen(0))
			Expect(fetchedElements).To(BeEmpty())
		})
	})

})
