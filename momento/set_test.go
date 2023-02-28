package momento_test

import (
	"fmt"
	"time"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	"github.com/momentohq/client-sdk-go/utils"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func getElements(numElements int) []Value {
	var elements []Value
	for i := 0; i < numElements; i++ {
		elements = append(elements, String(fmt.Sprintf("#%d", i)))
	}
	return elements
}

var _ = Describe("Set methods", func() {
	var sharedContext SharedContext

	BeforeEach(func() {
		sharedContext = NewSharedContext()
		sharedContext.CreateDefaultCache()
		DeferCleanup(func() {
			sharedContext.Close()
		})
	})

	It("errors when the cache is missing", func() {
		cacheName := uuid.NewString()
		setName := uuid.NewString()

		Expect(
			sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
				CacheName: cacheName,
				SetName:   setName,
			}),
		).Error().To(HaveMomentoErrorCode(NotFoundError))

		Expect(
			sharedContext.Client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
				CacheName: cacheName,
				SetName:   setName,
				Element:   String("astring"),
			}),
		).Error().To(HaveMomentoErrorCode(NotFoundError))

		Expect(
			sharedContext.Client.SetAddElements(sharedContext.Ctx, &SetAddElementsRequest{
				CacheName: cacheName,
				SetName:   setName,
				Elements:  []Value{String("astring"), String("bstring")},
			}),
		).Error().To(HaveMomentoErrorCode(NotFoundError))

		Expect(
			sharedContext.Client.SetRemoveElement(sharedContext.Ctx, &SetRemoveElementRequest{
				CacheName: cacheName,
				SetName:   setName,
				Element:   String("astring"),
			}),
		).Error().To(HaveMomentoErrorCode(NotFoundError))
	})

	It("errors on invalid set name", func() {
		setName := ""
		Expect(
			sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
				CacheName: sharedContext.CacheName,
				SetName:   setName,
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

		Expect(
			sharedContext.Client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
				CacheName: sharedContext.CacheName,
				SetName:   setName,
				Element:   String("hi"),
			}),
		)

		Expect(
			sharedContext.Client.SetAddElements(sharedContext.Ctx, &SetAddElementsRequest{
				CacheName: sharedContext.CacheName,
				SetName:   setName,
				Elements:  []Value{String("hi")},
			}),
		)

		Expect(
			sharedContext.Client.SetRemoveElement(sharedContext.Ctx, &SetRemoveElementRequest{
				CacheName: sharedContext.CacheName,
				SetName:   setName,
				Element:   nil,
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

		Expect(
			sharedContext.Client.SetRemoveElements(sharedContext.Ctx, &SetRemoveElementsRequest{
				CacheName: sharedContext.CacheName,
				SetName:   setName,
				Elements:  nil,
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	Describe("add", func() {
		DescribeTable("add string and byte single elements happy path",
			func(element Value, expectedStrings []string, expectedBytes [][]byte) {
				Expect(
					sharedContext.Client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Element:   element,
					}),
				).To(BeAssignableToTypeOf(&SetAddElementSuccess{}))

				fetchResp, err := sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				switch result := fetchResp.(type) {
				case *SetFetchHit:
					Expect(result.ValueString()).To(Equal(expectedStrings))
					Expect(result.ValueByte()).To(Equal(expectedBytes))
				default:
					Fail("Unexpected result for Set Fetch")
				}
			},
			Entry("when element is a string", String("hello"), []string{"hello"}, [][]byte{[]byte("hello")}),
			Entry("when element is bytes", Bytes("hello"), []string{"hello"}, [][]byte{[]byte("hello")}),
			Entry("when element is a empty", String(""), []string{""}, [][]byte{[]byte("")}),
		)

		DescribeTable("add string and byte multiple elements happy path",
			func(elements []Value, expectedStrings []string, expectedBytes [][]byte) {
				Expect(
					sharedContext.Client.SetAddElements(sharedContext.Ctx, &SetAddElementsRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Elements:  elements,
					}),
				).To(BeAssignableToTypeOf(&SetAddElementsSuccess{}))
				fetchResp, err := sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				switch result := fetchResp.(type) {
				case *SetFetchHit:
					Expect(result.ValueString()).To(ConsistOf(expectedStrings))
					Expect(result.ValueByte()).To(ConsistOf(expectedBytes))
				default:
					Fail("Unexpected results for Set Fetch")
				}
			},
			Entry(
				"when elements are strings",
				[]Value{String("hello"), String("world"), String("!"), String("␆")},
				[]string{"hello", "world", "!", "␆"},
				[][]byte{[]byte("hello"), []byte("world"), []byte("!"), []byte("␆")},
			),
			Entry(
				"when elements are bytes",
				[]Value{Bytes([]byte("hello")), Bytes([]byte("world")), Bytes([]byte("!")), Bytes([]byte("␆"))},
				[]string{"hello", "world", "!", "␆"},
				[][]byte{[]byte("hello"), []byte("world"), []byte("!"), []byte("␆")},
			),
			Entry(
				"when elements are mixed",
				[]Value{Bytes([]byte("hello")), String([]byte("world")), Bytes([]byte("!")), String([]byte("␆"))},
				[]string{"hello", "world", "!", "␆"},
				[][]byte{[]byte("hello"), []byte("world"), []byte("!"), []byte("␆")},
			),
			Entry(
				"when elements are empty",
				[]Value{Bytes([]byte("")), Bytes([]byte(""))},
				[]string{""},
				[][]byte{[]byte("")},
			),
		)

		It("returns an error when trying to add nil elements", func() {
			Expect(
				sharedContext.Client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Element:   nil,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.SetAddElements(sharedContext.Ctx, &SetAddElementsRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Elements:  nil,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.SetAddElements(sharedContext.Ctx, &SetAddElementsRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Elements:  []Value{nil, String("aValue"), nil},
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

	})

	Describe("remove", func() {

		BeforeEach(func() {
			elements := getElements(10)
			Expect(
				sharedContext.Client.SetAddElements(sharedContext.Ctx, &SetAddElementsRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Elements:  elements,
				}),
			).Error().To(BeNil())
		})

		DescribeTable("single elements as strings and as bytes",
			func(toRemove Value, expectedLength int) {
				Expect(
					sharedContext.Client.SetRemoveElement(sharedContext.Ctx, &SetRemoveElementRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Element:   toRemove,
					}),
				).Error().To(BeNil())

				fetchResp, err := sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveSetLength(expectedLength))
				switch result := fetchResp.(type) {
				case *SetFetchHit:
					Expect(result.ValueString()).ToNot(ContainElement(toRemove))
				default:
					Fail("something really weird happened")
				}
			},
			Entry("as string", String("#5"), 9),
			Entry("as bytes", Bytes([]byte("#5")), 9),
			Entry("unmatched", String("notvalid"), 10),
		)

		DescribeTable("multiple elements as strings and bytes",
			func(toRemove []Value, expectedLength int) {
				Expect(
					sharedContext.Client.SetRemoveElements(sharedContext.Ctx, &SetRemoveElementsRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Elements:  toRemove,
					}),
				).Error().To(BeNil())

				fetchResp, err := sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveSetLength(expectedLength))
				switch result := fetchResp.(type) {
				case *SetFetchHit:
					Expect(result.ValueString()).ToNot(ContainElements(toRemove))
				default:
					Fail("something really weird happened")
				}
			},
			Entry("as strings", getElements(5), 5),
			Entry("as bytes", []Value{Bytes("#3"), Bytes("#4")}, 8),
			Entry("unmatched", []Value{String("notvalid")}, 10),
		)

		It("returns an error when trying to remove nil elements", func() {
			Expect(
				sharedContext.Client.SetRemoveElement(sharedContext.Ctx, &SetRemoveElementRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Element:   nil,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.SetRemoveElements(sharedContext.Ctx, &SetRemoveElementsRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Elements:  nil,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.SetRemoveElements(sharedContext.Ctx, &SetRemoveElementsRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Elements:  []Value{nil, String("aValue"), nil},
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

		})

	})

	Describe("using client default TTL", func() {
		Context("when the TTL is exceeded", func() {
			It("returns a miss for the collection", func() {
				Expect(
					sharedContext.Client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Element:   String("hello"),
					}),
				).Error().To(BeNil())

				fetchResp, err := sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveSetLength(1))

				time.Sleep(sharedContext.DefaultTtl)

				fetchResp, err = sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(BeAssignableToTypeOf(&SetFetchMiss{}))
			})
		})
	})

	Describe("using collection ttl", func() {
		Context("when collection ttl is longer than client default", func() {

			BeforeEach(func() {
				// Initialize the set. If the set isn't initialized, there's
				// nothing to refresh and it will use the passed in TTL.
				Expect(
					sharedContext.Client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Element:   String("goodbye"),
					}),
				).To(BeAssignableToTypeOf(&SetAddElementSuccess{}))
			})

			It("returns a hit after the client default has expired", func() {
				Expect(
					sharedContext.Client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Element:   String("hello"),
						Ttl: utils.CollectionTtl{
							Ttl:        sharedContext.DefaultTTL + time.Second*60,
							RefreshTtl: true,
						},
					}),
				).Error().To(BeNil())

				fetchResp, err := sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveSetLength(2))

				time.Sleep(sharedContext.DefaultTtl)

				fetchResp, err = sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(BeAssignableToTypeOf(&SetFetchHit{}))
			})

			It("returns a miss after the client default when refreshTTL is false", func() {
				Expect(
					sharedContext.Client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Element:   String("hello"),
						Ttl: utils.CollectionTtl{
							Ttl:        sharedContext.DefaultTTL + 1*time.Second,
							RefreshTtl: false,
						},
					}),
				).To(BeAssignableToTypeOf(&SetAddElementSuccess{}))

				Expect(
					sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
					}),
				).To(HaveSetLength(2))

				time.Sleep(sharedContext.DefaultTtl + 500*time.Millisecond)

				Expect(
					sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
					}),
				).To(BeAssignableToTypeOf(&SetFetchMiss{}))
			})

			It("ignores collection ttl when refresh ttl is false", func() {
				Expect(
					sharedContext.Client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Element:   String("hello"),
						Ttl: utils.CollectionTtl{
							Ttl:        time.Millisecond * 20,
							RefreshTtl: false,
						},
					}),
				).Error().To(BeNil())

				fetchResp, err := sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveSetLength(2))

				time.Sleep(sharedContext.DefaultTtl / 2)

				fetchResp, err = sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(BeAssignableToTypeOf(&SetFetchHit{}))
			})

			It("returns a miss after overriding the client timeout with a short duration", func() {
				Expect(
					sharedContext.Client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
						Element:   String("hello"),
						Ttl: utils.CollectionTtl{
							Ttl:        time.Millisecond * 200,
							RefreshTtl: true,
						},
					}),
				).To(BeAssignableToTypeOf(&SetAddElementSuccess{}))

				time.Sleep(time.Millisecond * 500)

				Expect(
					sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sharedContext.CollectionName,
					}),
				).To(BeAssignableToTypeOf(&SetFetchMiss{}))
			})
		})
	})
})
