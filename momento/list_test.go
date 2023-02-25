package momento_test

import (
	"fmt"

	"github.com/google/uuid"
	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func getValueAndExpectedValueLists(numItems int) ([]Value, []string) {
	var values []Value
	var expected []string
	for i := 0; i < numItems; i++ {
		strVal := fmt.Sprintf("#%d", i)
		var value Value
		if i%2 == 0 {
			value = String(strVal)
		} else {
			value = Bytes(strVal)
		}
		values = append(values, value)
		expected = append(expected, strVal)
	}
	return values, expected
}

func populateList(sharedContext SharedContext, numItems int) []string {
	values, expected := getValueAndExpectedValueLists(numItems)
	for _, value := range values {
		Expect(
			sharedContext.Client.ListPushFront(sharedContext.Ctx, &ListPushFrontRequest{
				CacheName:          sharedContext.CacheName,
				ListName:           sharedContext.CollectionName,
				Value:              value,
				TruncateBackToSize: 0,
			}),
		).To(BeAssignableToTypeOf(&ListPushFrontSuccess{}))
	}
	return expected
}

var _ = Describe("List methods", func() {
	var sharedContext SharedContext

	BeforeEach(func() {
		sharedContext = NewSharedContext()
		sharedContext.CreateDefaultCache()
		DeferCleanup(func() {
			sharedContext.Close()
		})
	})

	DescribeTable("try using invalid cache and list names",
		func(cacheName string, listName string, expectedErrorCode string) {
			Expect(
				sharedContext.Client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: cacheName,
					ListName:  listName,
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				sharedContext.Client.ListLength(sharedContext.Ctx, &ListLengthRequest{
					CacheName: cacheName,
					ListName:  listName,
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

		},
		Entry("nonexistent cache name", uuid.NewString(), uuid.NewString(), NotFoundError),
		Entry("empty cache name", "", sharedContext.CollectionName, InvalidArgumentError),
		Entry("empty list name", sharedContext.CacheName, "", InvalidArgumentError),
	)

	Describe("list push", func() {

		When("pushing to the front of the list", func() {

			It("pushes strings and bytes on the happy path", func() {
				numItems := 10
				expected := populateList(sharedContext, numItems)

				fetchResp, err := sharedContext.Client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: sharedContext.CacheName,
					ListName:  sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveListLength(numItems))
				switch result := fetchResp.(type) {
				case *ListFetchHit:
					Expect(result.ValueList()).To(ConsistOf(expected))
					Expect(result.ValueList()).NotTo(Equal(expected))
				}
			})

			It("truncates the list properly", func() {
				numItems := 10
				truncateTo := 5
				populateList(sharedContext, numItems)
				Expect(
					sharedContext.Client.ListPushFront(sharedContext.Ctx, &ListPushFrontRequest{
						CacheName:          sharedContext.CacheName,
						ListName:           sharedContext.CollectionName,
						Value:              String("andherlittledogtoo"),
						TruncateBackToSize: uint32(truncateTo),
					}),
				).Error().To(BeNil())
				fetchResp, err := sharedContext.Client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: sharedContext.CacheName,
					ListName:  sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveListLength(truncateTo))
			})

			It("returns invalid argument for a nil value", func() {
				Expect(
					sharedContext.Client.ListPushBack(sharedContext.Ctx, &ListPushBackRequest{
						CacheName: sharedContext.CacheName,
						ListName:  sharedContext.CollectionName,
						Value:     nil,
					}),
				).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
			})

		})

		When("pushing to the back of the list", func() {

			It("pushes strings and bytes on the happy path", func() {
				numItems := 10
				values, expected := getValueAndExpectedValueLists(numItems)
				for _, value := range values {
					Expect(
						sharedContext.Client.ListPushBack(sharedContext.Ctx, &ListPushBackRequest{
							CacheName: sharedContext.CacheName,
							ListName:  sharedContext.CollectionName,
							Value:     value,
						}),
					).To(BeAssignableToTypeOf(&ListPushBackSuccess{}))
				}

				fetchResp, err := sharedContext.Client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: sharedContext.CacheName,
					ListName:  sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveListLength(numItems))
				switch result := fetchResp.(type) {
				case *ListFetchHit:
					Expect(result.ValueList()).To(ConsistOf(expected))
					Expect(result.ValueList()).To(Equal(expected))
				}
			})

			It("truncates the list properly", func() {
				numItems := 10
				truncateTo := 5
				populateList(sharedContext, numItems)
				Expect(
					sharedContext.Client.ListPushBack(sharedContext.Ctx, &ListPushBackRequest{
						CacheName:           sharedContext.CacheName,
						ListName:            sharedContext.CollectionName,
						Value:               String("andherlittledogtoo"),
						TruncateFrontToSize: uint32(truncateTo),
					}),
				).Error().To(BeNil())
				fetchResp, err := sharedContext.Client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: sharedContext.CacheName,
					ListName:  sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveListLength(truncateTo))
			})

			It("returns invalid argument for a nil value", func() {
				Expect(
					sharedContext.Client.ListPushBack(sharedContext.Ctx, &ListPushBackRequest{
						CacheName: sharedContext.CacheName,
						ListName:  sharedContext.CollectionName,
						Value:     nil,
					}),
				).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
			})

		})

	})

	Describe("list concatenate", func() {

		When("concatenating to the front of the list", func() {

			It("pushes strings and bytes on the happy path", func() {
				numItems := 10
				expected := populateList(sharedContext, numItems)
				// items are in reverse order from expected because of how they're pushed.
				for i, j := 0, len(expected)-1; i < j; i, j = i+1, j-1 {
					expected[i], expected[j] = expected[j], expected[i]
				}

				numConcatItems := 5
				concatValues, concatExpected := getValueAndExpectedValueLists(numConcatItems)
				concatResp, err := sharedContext.Client.ListConcatenateFront(sharedContext.Ctx, &ListConcatenateFrontRequest{
					CacheName: sharedContext.CacheName,
					ListName:  sharedContext.CollectionName,
					Values:    concatValues,
				})
				Expect(err).To(BeNil())
				Expect(concatResp).To(BeAssignableToTypeOf(&ListConcatenateFrontSuccess{}))

				fetchResp, err := sharedContext.Client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: sharedContext.CacheName,
					ListName:  sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(BeAssignableToTypeOf(&ListFetchHit{}))
				Expect(fetchResp).To(HaveListLength(numItems + numConcatItems))
				expected = append(concatExpected, expected...)
				switch result := fetchResp.(type) {
				case *ListFetchHit:
					Expect(result.ValueList()).To(Equal(expected))
				}
			})

			It("truncates the list properly", func() {
				populateList(sharedContext, 5)
				concatValues := []Value{String("100"), String("101"), String("102")}
				concatResp, err := sharedContext.Client.ListConcatenateFront(sharedContext.Ctx, &ListConcatenateFrontRequest{
					CacheName:          sharedContext.CacheName,
					ListName:           sharedContext.CollectionName,
					Values:             concatValues,
					TruncateBackToSize: 3,
				})
				Expect(err).To(BeNil())
				Expect(concatResp).To(BeAssignableToTypeOf(&ListConcatenateFrontSuccess{}))

				fetchResp, err := sharedContext.Client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: sharedContext.CacheName,
					ListName:  sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(BeAssignableToTypeOf(&ListFetchHit{}))
				Expect(fetchResp).To(HaveListLength(3))
				switch result := fetchResp.(type) {
				case *ListFetchHit:
					Expect(result.ValueList()).To(Equal([]string{"100", "101", "102"}))
				}
			})

			It("returns an invalid argument for a nil value", func() {
				populateList(sharedContext, 5)
				concatValues := []Value{nil, nil}
				Expect(
					sharedContext.Client.ListConcatenateFront(sharedContext.Ctx, &ListConcatenateFrontRequest{
						CacheName:          sharedContext.CacheName,
						ListName:           sharedContext.CollectionName,
						Values:             concatValues,
						TruncateBackToSize: 3,
					}),
				).Error().To(BeAssignableToTypeOf(InvalidArgumentError))
			})

		})

	})
})
