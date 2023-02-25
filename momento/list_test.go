package momento_test

import (
	"fmt"

	"github.com/google/uuid"
	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func getValueAndExpectedValueLists() ([]Value, []string) {
	var values []Value
	var expected []string
	for i := 0; i < 10; i++ {
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
				values, expected := getValueAndExpectedValueLists()
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

				fetchResp, err := sharedContext.Client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: sharedContext.CacheName,
					ListName:  sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveListLength(10))
				switch result := fetchResp.(type) {
				case *ListFetchHit:
					Expect(result.ValueList()).To(ConsistOf(expected))
					Expect(result.ValueList()).NotTo(Equal(expected))
				}
			})

		})

		When("pushing to the back of the list", func() {

			It("pushes strings and bytes on the happy path", func() {
				values, expected := getValueAndExpectedValueLists()
				for _, value := range values {
					Expect(
						sharedContext.Client.ListPushBack(sharedContext.Ctx, &ListPushBackRequest{
							CacheName:           sharedContext.CacheName,
							ListName:            sharedContext.CollectionName,
							Value:               value,
							TruncateFrontToSize: 0,
						}),
					).To(BeAssignableToTypeOf(&ListPushBackSuccess{}))
				}

				fetchResp, err := sharedContext.Client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: sharedContext.CacheName,
					ListName:  sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveListLength(10))
				switch result := fetchResp.(type) {
				case *ListFetchHit:
					Expect(result.ValueList()).To(ConsistOf(expected))
					Expect(result.ValueList()).To(Equal(expected))
				}
			})

		})

	})
})
