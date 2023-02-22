package momento_test

import (
	"context"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/utils"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	. "github.com/momentohq/client-sdk-go/momento"
)

func HaveLength(length int) types.GomegaMatcher {
	return WithTransform(
		func(fetchResp SetFetchResponse) (int, error) {
			switch rtype := fetchResp.(type) {
			case *SetFetchHit:
				return len(rtype.ValueString()), nil
			default:
				return 0, fmt.Errorf("expected set fetch hit but got %T", fetchResp)
			}
		}, Equal(length),
	)
}

func getElements(numElements int) []Value {
	var elements []Value
	for i := 0; i < numElements; i++ {
		elements = append(elements, String(fmt.Sprintf("#%d", i)))
	}
	return elements
}

var _ = Describe("Set methods", func() {
	var clientProps SimpleCacheClientProps
	var credentialProvider auth.CredentialProvider
	var configuration config.Configuration
	var client SimpleCacheClient
	var defaultTTL time.Duration
	var testCacheName string
	var testSetName string
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		credentialProvider, _ = auth.NewEnvMomentoTokenProvider("TEST_AUTH_TOKEN")
		configuration = config.LatestLaptopConfig()
		defaultTTL = 3 * time.Second

		clientProps = SimpleCacheClientProps{
			CredentialProvider: credentialProvider,
			Configuration:      configuration,
			DefaultTTL:         defaultTTL,
		}

		var err error
		client, err = NewSimpleCacheClient(&clientProps)
		if err != nil {
			panic(err)
		}
		DeferCleanup(func() { client.Close() })

		testCacheName = uuid.NewString()
		testSetName = uuid.NewString()
		Expect(
			client.CreateCache(ctx, &CreateCacheRequest{CacheName: testCacheName}),
		).To(BeAssignableToTypeOf(&CreateCacheSuccess{}))
		DeferCleanup(func() {
			_, err := client.DeleteCache(ctx, &DeleteCacheRequest{CacheName: testCacheName})
			if err != nil {
				panic(err)
			}
		})
	})

	It("errors when the cache is missing", func() {
		cacheName := uuid.NewString()
		setName := uuid.NewString()

		fetchResp, err := client.SetFetch(ctx, &SetFetchRequest{
			CacheName: cacheName,
			SetName:   setName,
		})
		Expect(fetchResp).To(BeNil())
		Expect(err).To(HaveMomentoErrorCode(NotFoundError))

		addOneResp, err := client.SetAddElement(ctx, &SetAddElementRequest{
			CacheName: cacheName,
			SetName:   setName,
			Element:   String("astring"),
		})
		Expect(addOneResp).To(BeNil())
		Expect(err).To(HaveMomentoErrorCode(NotFoundError))

		addManyResp, err := client.SetAddElements(ctx, &SetAddElementsRequest{
			CacheName: cacheName,
			SetName:   setName,
			Elements:  []Value{String("astring"), String("bstring")},
		})
		Expect(addManyResp).To(BeNil())
		Expect(err).To(HaveMomentoErrorCode(NotFoundError))

		removeOneResp, err := client.SetRemoveElement(ctx, &SetRemoveElementRequest{
			CacheName: cacheName,
			SetName:   setName,
			Element:   String("astring"),
		})
		Expect(removeOneResp).To(BeNil())
		Expect(err).To(HaveMomentoErrorCode(NotFoundError))

		removeManyResp, err := client.SetRemoveElements(ctx, &SetRemoveElementsRequest{
			CacheName: cacheName,
			SetName:   setName,
			Elements:  nil,
		})
		Expect(removeManyResp).To(BeNil())
		Expect(err).To(HaveMomentoErrorCode(NotFoundError))
	})

	It("errors on invalid set name", func() {
		setName := ""
		fetchResp, err := client.SetFetch(ctx, &SetFetchRequest{
			CacheName: testCacheName,
			SetName:   setName,
		})
		Expect(fetchResp).To(BeNil())
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

		removeManyResp, err := client.SetRemoveElements(ctx, &SetRemoveElementsRequest{
			CacheName: testCacheName,
			SetName:   setName,
			Elements:  nil,
		})
		Expect(removeManyResp).To(BeNil())
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	DescribeTable("add string and byte single elements happy path",
		func(element Value, expectedStrings []string, expectedBytes [][]byte) {
			Expect(
				client.SetAddElement(ctx, &SetAddElementRequest{
					CacheName: testCacheName,
					SetName:   testSetName,
					Element:   element,
				}),
			).To(BeAssignableToTypeOf(&SetAddElementSuccess{}))

			fetchResp, err := client.SetFetch(ctx, &SetFetchRequest{
				CacheName: testCacheName,
				SetName:   testSetName,
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
				client.SetAddElements(ctx, &SetAddElementsRequest{
					CacheName: testCacheName,
					SetName:   testSetName,
					Elements:  elements,
				}),
			).To(BeAssignableToTypeOf(&SetAddElementsSuccess{}))
			fetchResp, err := client.SetFetch(ctx, &SetFetchRequest{
				CacheName: testCacheName,
				SetName:   testSetName,
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

	Describe("remove", func() {

		BeforeEach(func() {
			elements := getElements(10)
			_, err := client.SetAddElements(ctx, &SetAddElementsRequest{
				CacheName: testCacheName,
				SetName:   testSetName,
				Elements:  elements,
			})
			Expect(err).To(BeNil())
		})

		DescribeTable("single elements as strings and as bytes",
			func(toRemove Value, expectedLength int) {
				_, err := client.SetRemoveElement(ctx, &SetRemoveElementRequest{
					CacheName: testCacheName,
					SetName:   testSetName,
					Element:   toRemove,
				})
				Expect(err).To(BeNil())
				fetchResp, err := client.SetFetch(ctx, &SetFetchRequest{
					CacheName: testCacheName,
					SetName:   testSetName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveLength(expectedLength))
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
				_, err := client.SetRemoveElements(ctx, &SetRemoveElementsRequest{
					CacheName: testCacheName,
					SetName:   testSetName,
					Elements:  toRemove,
				})
				Expect(err).To(BeNil())
				fetchResp, err := client.SetFetch(ctx, &SetFetchRequest{
					CacheName: testCacheName,
					SetName:   testSetName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveLength(expectedLength))
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
	})

	Describe("using client default TTL", func() {
		Context("when the TTL is exceeded", func() {
			It("returns a miss for the collection", func() {
				_, err := client.SetAddElement(ctx, &SetAddElementRequest{
					CacheName: testCacheName,
					SetName:   testSetName,
					Element:   String("hello"),
				})
				Expect(err).To(BeNil())

				fetchResp, err := client.SetFetch(ctx, &SetFetchRequest{
					CacheName: testCacheName,
					SetName:   testSetName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveLength(1))

				time.Sleep(defaultTTL)

				fetchResp, err = client.SetFetch(ctx, &SetFetchRequest{
					CacheName: testCacheName,
					SetName:   testSetName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(BeAssignableToTypeOf(&SetFetchMiss{}))
			})
		})
	})

	Describe("using collection ttl", func() {
		Context("when collection ttl is longer than client default", func() {
			It("returns a hit after the client default has expired", func() {
				_, err := client.SetAddElement(ctx, &SetAddElementRequest{
					CacheName: testCacheName,
					SetName:   testSetName,
					Element:   String("hello"),
					CollectionTTL: utils.CollectionTTL{
						Ttl:        time.Second * 10,
						RefreshTtl: true,
					},
				})
				Expect(err).To(BeNil())

				fetchResp, err := client.SetFetch(ctx, &SetFetchRequest{
					CacheName: testCacheName,
					SetName:   testSetName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveLength(1))

				time.Sleep(defaultTTL)

				fetchResp, err = client.SetFetch(ctx, &SetFetchRequest{
					CacheName: testCacheName,
					SetName:   testSetName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(BeAssignableToTypeOf(&SetFetchHit{}))
			})

			It("returns a miss after the client default when refreshTTL is false", func() {
				// Initialize the set. If the set isn't initialized, there's
				// nothing to refresh and it will use the passed in TTL.
				Expect(
					client.SetAddElement(ctx, &SetAddElementRequest{
						CacheName: testCacheName,
						SetName:   testSetName,
						Element:   String("goodbye"),
					}),
				).To(BeAssignableToTypeOf(&SetAddElementSuccess{}))

				Expect(
					client.SetAddElement(ctx, &SetAddElementRequest{
						CacheName: testCacheName,
						SetName:   testSetName,
						Element:   String("hello"),
						CollectionTTL: utils.CollectionTTL{
							Ttl:        defaultTTL + 1*time.Second,
							RefreshTtl: false,
						},
					}),
				).To(BeAssignableToTypeOf(&SetAddElementSuccess{}))

				Expect(
					client.SetFetch(ctx, &SetFetchRequest{
						CacheName: testCacheName,
						SetName:   testSetName,
					}),
				).To(HaveLength(2))

				time.Sleep(defaultTTL + 500*time.Millisecond)

				Expect(
					client.SetFetch(ctx, &SetFetchRequest{
						CacheName: testCacheName,
						SetName:   testSetName,
					}),
				).To(BeAssignableToTypeOf(&SetFetchMiss{}))
			})

			It("returns a miss after overriding the client timeout with a short duration", func() {
				Expect(
					client.SetAddElement(ctx, &SetAddElementRequest{
						CacheName: testCacheName,
						SetName:   testSetName,
						Element:   String("hello"),
						CollectionTTL: utils.CollectionTTL{
							Ttl:        time.Millisecond * 200,
							RefreshTtl: true,
						},
					}),
				).To(BeAssignableToTypeOf(&SetAddElementSuccess{}))

				time.Sleep(time.Millisecond * 500)

				Expect(
					client.SetFetch(ctx, &SetFetchRequest{
						CacheName: testCacheName,
						SetName:   testSetName,
					}),
				).To(BeAssignableToTypeOf(&SetFetchMiss{}))
			})
		})
	})
})
