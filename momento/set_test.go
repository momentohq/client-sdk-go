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
		sharedContext.CreateDefaultCaches()
		DeferCleanup(func() {
			sharedContext.Close()
		})
	})

	DescribeTable("errors when the cache is missing",
		func(clientType string) {
			client, _ := sharedContext.GetClientPrereqsForType(clientType)
			cacheName := uuid.NewString()
			setName := uuid.NewString()

			Expect(
				client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: cacheName,
					SetName:   setName,
				}),
			).Error().To(HaveMomentoErrorCode(NotFoundError))

			Expect(
				client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
					CacheName: cacheName,
					SetName:   setName,
					Element:   String("astring"),
				}),
			).Error().To(HaveMomentoErrorCode(NotFoundError))

			Expect(
				client.SetAddElements(sharedContext.Ctx, &SetAddElementsRequest{
					CacheName: cacheName,
					SetName:   setName,
					Elements:  []Value{String("astring"), String("bstring")},
				}),
			).Error().To(HaveMomentoErrorCode(NotFoundError))

			Expect(
				client.SetLength(sharedContext.Ctx, &SetLengthRequest{
					CacheName: cacheName,
					SetName:   setName,
				}),
			).Error().To(HaveMomentoErrorCode(NotFoundError))

			Expect(
				client.SetRemoveElement(sharedContext.Ctx, &SetRemoveElementRequest{
					CacheName: cacheName,
					SetName:   setName,
					Element:   String("astring"),
				}),
			).Error().To(HaveMomentoErrorCode(NotFoundError))

			Expect(
				client.SetContainsElements(sharedContext.Ctx, &SetContainsElementsRequest{
					CacheName: cacheName,
					SetName:   setName,
					Elements:  []Value{String("hi")},
				}),
			).Error().To(HaveMomentoErrorCode(NotFoundError))
		},
		Entry("with default client", DefaultClient),
		Entry("with client with default cache", WithDefaultCache),
	)

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
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

		Expect(
			sharedContext.Client.SetAddElements(sharedContext.Ctx, &SetAddElementsRequest{
				CacheName: sharedContext.CacheName,
				SetName:   setName,
				Elements:  []Value{String("hi")},
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

		Expect(
			sharedContext.Client.SetLength(sharedContext.Ctx, &SetLengthRequest{
				CacheName: sharedContext.CacheName,
				SetName:   setName,
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

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

		Expect(
			sharedContext.Client.SetContainsElements(sharedContext.Ctx, &SetContainsElementsRequest{
				CacheName: sharedContext.CacheName,
				SetName:   setName,
				Elements:  []Value{String("hi")},
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

	})

	It("gets a miss trying to fetch a nonexistent set", func() {
		Expect(
			sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
				CacheName: sharedContext.CacheName,
				SetName:   uuid.NewString(),
			}),
		).To(BeAssignableToTypeOf(&SetFetchMiss{}))
	})

	Describe("add", func() {
		DescribeTable("add string and byte single elements happy path",
			func(clientType string, element Value, expectedStrings []string, expectedBytes [][]byte) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				Expect(
					client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
						CacheName: cacheName,
						SetName:   sharedContext.CollectionName,
						Element:   element,
					}),
				).To(BeAssignableToTypeOf(&SetAddElementSuccess{}))

				fetchResp, err := client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: cacheName,
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
			Entry("when element is a string", DefaultClient, String("hello"), []string{"hello"}, [][]byte{[]byte("hello")}),
			Entry("when element is bytes", DefaultClient, Bytes("hello"), []string{"hello"}, [][]byte{[]byte("hello")}),
			Entry("when element is a empty", DefaultClient, String(""), []string{""}, [][]byte{[]byte("")}),
			Entry("when element is a string", WithDefaultCache, String("hello"), []string{"hello"}, [][]byte{[]byte("hello")}),
			Entry("when element is bytes", WithDefaultCache, Bytes("hello"), []string{"hello"}, [][]byte{[]byte("hello")}),
			Entry("when element is a empty", WithDefaultCache, String(""), []string{""}, [][]byte{[]byte("")}),
		)

		DescribeTable("add string and byte multiple elements happy path",
			func(clientType string, elements []Value, expectedStrings []string, expectedBytes [][]byte) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				Expect(
					client.SetAddElements(sharedContext.Ctx, &SetAddElementsRequest{
						CacheName: cacheName,
						SetName:   sharedContext.CollectionName,
						Elements:  elements,
					}),
				).To(BeAssignableToTypeOf(&SetAddElementsSuccess{}))
				fetchResp, err := client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: cacheName,
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
				"with default client when elements are strings",
				DefaultClient,
				[]Value{String("hello"), String("world"), String("!"), String("␆")},
				[]string{"hello", "world", "!", "␆"},
				[][]byte{[]byte("hello"), []byte("world"), []byte("!"), []byte("␆")},
			),
			Entry(
				"with default client when elements are bytes",
				DefaultClient,
				[]Value{Bytes([]byte("hello")), Bytes([]byte("world")), Bytes([]byte("!")), Bytes([]byte("␆"))},
				[]string{"hello", "world", "!", "␆"},
				[][]byte{[]byte("hello"), []byte("world"), []byte("!"), []byte("␆")},
			),
			Entry(
				"with default client when elements are mixed",
				DefaultClient,
				[]Value{Bytes([]byte("hello")), String([]byte("world")), Bytes([]byte("!")), String([]byte("␆"))},
				[]string{"hello", "world", "!", "␆"},
				[][]byte{[]byte("hello"), []byte("world"), []byte("!"), []byte("␆")},
			),
			Entry(
				"with default client when elements are empty",
				DefaultClient,
				[]Value{Bytes([]byte("")), Bytes([]byte(""))},
				[]string{""},
				[][]byte{[]byte("")},
			),
			Entry(
				"with client with default cache when elements are strings",
				WithDefaultCache,
				[]Value{String("hello"), String("world"), String("!"), String("␆")},
				[]string{"hello", "world", "!", "␆"},
				[][]byte{[]byte("hello"), []byte("world"), []byte("!"), []byte("␆")},
			),
			Entry(
				"with client with default cache when elements are bytes",
				WithDefaultCache,
				[]Value{Bytes([]byte("hello")), Bytes([]byte("world")), Bytes([]byte("!")), Bytes([]byte("␆"))},
				[]string{"hello", "world", "!", "␆"},
				[][]byte{[]byte("hello"), []byte("world"), []byte("!"), []byte("␆")},
			),
			Entry(
				"with client with default cache when elements are mixed",
				WithDefaultCache,
				[]Value{Bytes([]byte("hello")), String([]byte("world")), Bytes([]byte("!")), String([]byte("␆"))},
				[]string{"hello", "world", "!", "␆"},
				[][]byte{[]byte("hello"), []byte("world"), []byte("!"), []byte("␆")},
			),
			Entry(
				"with client with default cache when elements are empty",
				WithDefaultCache,
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
			Expect(
				sharedContext.ClientWithDefaultCacheName.SetAddElements(sharedContext.Ctx, &SetAddElementsRequest{
					SetName:  sharedContext.CollectionName,
					Elements: elements,
				}),
			).Error().To(BeNil())
		})

		DescribeTable("single elements as strings and as bytes",
			func(clientType string, toRemove Value, expectedLength int) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				Expect(
					client.SetRemoveElement(sharedContext.Ctx, &SetRemoveElementRequest{
						CacheName: cacheName,
						SetName:   sharedContext.CollectionName,
						Element:   toRemove,
					}),
				).Error().To(BeNil())

				fetchResp, err := client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: cacheName,
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
			Entry("with default client as string", DefaultClient, String("#5"), 9),
			Entry("with default client as bytes", DefaultClient, Bytes([]byte("#5")), 9),
			Entry("with default client unmatched", DefaultClient, String("notvalid"), 10),
			Entry("with client with default cache as string", WithDefaultCache, String("#5"), 9),
			Entry("with client with default cache as bytes", WithDefaultCache, Bytes([]byte("#5")), 9),
			Entry("with client with default cache unmatched", WithDefaultCache, String("notvalid"), 10),
		)

		DescribeTable("multiple elements as strings and bytes",
			func(clientType string, toRemove []Value, expectedLength int) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				Expect(
					client.SetRemoveElements(sharedContext.Ctx, &SetRemoveElementsRequest{
						CacheName: cacheName,
						SetName:   sharedContext.CollectionName,
						Elements:  toRemove,
					}),
				).Error().To(BeNil())

				fetchResp, err := client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: cacheName,
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
			Entry("with default client as strings", DefaultClient, getElements(5), 5),
			Entry("with default client as bytes", DefaultClient, []Value{Bytes("#3"), Bytes("#4")}, 8),
			Entry("with default client unmatched", DefaultClient, []Value{String("notvalid")}, 10),
			Entry("with client with default cache as strings", WithDefaultCache, getElements(5), 5),
			Entry("with client with default cache as bytes", WithDefaultCache, []Value{Bytes("#3"), Bytes("#4")}, 8),
			Entry("with client with default cache unmatched", WithDefaultCache, []Value{String("notvalid")}, 10),
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

	It("Returns the correct Set length", func() {
		elements := getElements(7)

		_, err := sharedContext.Client.SetAddElements(sharedContext.Ctx, &SetAddElementsRequest{
			CacheName: sharedContext.CacheName,
			SetName:   sharedContext.CollectionName,
			Elements:  elements,
		})
		Expect(err).To(BeNil())
		resp, err := sharedContext.Client.SetLength(sharedContext.Ctx, &SetLengthRequest{
			CacheName: sharedContext.CacheName,
			SetName:   sharedContext.CollectionName,
		})
		Expect(err).To(BeNil())
		switch result := resp.(type) {
		case *SetLengthHit:
			Expect(result.Length()).To(Equal(uint32(len(elements))))
		default:
			Fail("expected a hit for set length but got a miss")
		}

		resp, err = sharedContext.Client.SetLength(sharedContext.Ctx, &SetLengthRequest{
			CacheName: sharedContext.CacheName,
			SetName:   "IdontExist",
		})
		Expect(err).To(BeNil())
		Expect(resp).To(BeAssignableToTypeOf(&SetLengthMiss{}))
	})

	Describe("contain elements", func() {
		BeforeEach(func() {
			elements := getElements(10)
			Expect(
				sharedContext.Client.SetAddElements(sharedContext.Ctx, &SetAddElementsRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
					Elements:  elements,
				}),
			).Error().To(BeNil())
			Expect(
				sharedContext.ClientWithDefaultCacheName.SetAddElements(sharedContext.Ctx, &SetAddElementsRequest{
					SetName:  sharedContext.CollectionName,
					Elements: elements,
				}),
			).Error().To(BeNil())
		})

		DescribeTable("check for various mixes of hits and misses",
			func(clientType string, toCheck []Value, expected []bool) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				containsResp, err := client.SetContainsElements(sharedContext.Ctx, &SetContainsElementsRequest{
					CacheName: cacheName,
					SetName:   sharedContext.CollectionName,
					Elements:  toCheck,
				})
				Expect(err).To(BeNil())
				Expect(containsResp).To(BeAssignableToTypeOf(&SetContainsElementsHit{}))
				switch result := containsResp.(type) {
				case *SetContainsElementsHit:
					Expect(result.ContainsElements()).To(Equal(expected))
				}
			},
			Entry("with default client all hits", DefaultClient, []Value{String("#1"), String("#2"), String("#3")}, []bool{true, true, true}),
			Entry("with default client all misses", DefaultClient, []Value{String("not"), String("this"), String("time")}, []bool{false, false, false}),
			Entry("with default client a mixture", DefaultClient, []Value{String("not"), String("#2"), String("time")}, []bool{false, true, false}),
			Entry("with client with default cache all hits", WithDefaultCache, []Value{String("#1"), String("#2"), String("#3")}, []bool{true, true, true}),
			Entry("with client with default cache all misses", WithDefaultCache, []Value{String("not"), String("this"), String("time")}, []bool{false, false, false}),
			Entry("with client with default cache a mixture", WithDefaultCache, []Value{String("not"), String("#2"), String("time")}, []bool{false, true, false}),
		)

		It("gets a miss on a nonexistent set", func() {
			Expect(
				sharedContext.Client.SetContainsElements(sharedContext.Ctx, &SetContainsElementsRequest{
					CacheName: sharedContext.CacheName,
					SetName:   uuid.NewString(),
					Elements:  []Value{String("hi")},
				}),
			).To(BeAssignableToTypeOf(&SetContainsElementsMiss{}))
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
						Ttl: &utils.CollectionTtl{
							Ttl:        sharedContext.DefaultTtl + time.Second*60,
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
						Ttl: &utils.CollectionTtl{
							Ttl:        sharedContext.DefaultTtl + 1*time.Second,
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
						Ttl: &utils.CollectionTtl{
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
						Ttl: &utils.CollectionTtl{
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

	Describe("set pop", func() {
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

		It("gets a miss on a nonexistent set", func() {
			resp, err := sharedContext.Client.SetPop(sharedContext.Ctx, &SetPopRequest{
				CacheName: sharedContext.CacheName,
				SetName:   uuid.NewString(),
			})
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(&SetPopMiss{}))
		})

		It("pops elements off until empty", func() {
			var count uint32

			// Pop one item from the set (1 is the default), 9 should remain
			resp1, err1 := sharedContext.Client.SetPop(sharedContext.Ctx, &SetPopRequest{
				CacheName: sharedContext.CacheName,
				SetName:   sharedContext.CollectionName,
			})
			Expect(err1).To(BeNil())
			Expect(resp1).To(BeAssignableToTypeOf(&SetPopHit{}))
			Expect(
				sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
				}),
			).To(HaveSetLength(9))

			// Pop 4 items from the set, 5 should remain
			count = 4
			resp2, err2 := sharedContext.Client.SetPop(sharedContext.Ctx, &SetPopRequest{
				CacheName: sharedContext.CacheName,
				SetName:   sharedContext.CollectionName,
				Count:     &count,
			})
			Expect(err2).To(BeNil())
			Expect(resp2).To(BeAssignableToTypeOf(&SetPopHit{}))
			Expect(
				sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
				}),
			).To(HaveSetLength(5))

			// Pop 5 items from the set, none should remain
			count = 5
			resp3, err3 := sharedContext.Client.SetPop(sharedContext.Ctx, &SetPopRequest{
				CacheName: sharedContext.CacheName,
				SetName:   sharedContext.CollectionName,
				Count:     &count,
			})
			Expect(err3).To(BeNil())
			Expect(resp3).To(BeAssignableToTypeOf(&SetPopHit{}))
			Expect(
				sharedContext.Client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: sharedContext.CacheName,
					SetName:   sharedContext.CollectionName,
				}),
			).To(BeAssignableToTypeOf(&SetFetchMiss{}))

			// Expect a miss when set is empty
			resp4, err4 := sharedContext.Client.SetPop(sharedContext.Ctx, &SetPopRequest{
				CacheName: sharedContext.CacheName,
				SetName:   sharedContext.CollectionName,
			})
			Expect(err4).To(BeNil())
			Expect(resp4).To(BeAssignableToTypeOf(&SetPopMiss{}))
		})

	})
})
