package momento_test

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	"github.com/momentohq/client-sdk-go/responses"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// var sharedContext SharedContext

// var _ = BeforeSuite(func() {
// 	sharedContext = NewSharedContext(SharedContextProps{IsV2ApiKey: true})
// 	sharedContext.CreateDefaultCaches()
// })

// var _ = AfterSuite(func() {
// 	sharedContext.Close()
// })

var _ = Describe("v2 api key tests", Label(V2_API_LABEL), func() {
	// var sharedContext SharedContext

	// BeforeEach(func() {
	// 	if sharedContext.Client == nil {
	// 		sharedContext = NewSharedContext(SharedContextProps{IsV2ApiKey: true})
	// 		sharedContext.CreateDefaultCaches()
	// 	}
	// })

	Describe("BatchGetSet", func() {
		It("GetBatch happy path with all misses", func() {
			var batchSetKeys []Value
			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("Batch-key-%d", i)
				batchSetKeys = append(batchSetKeys, String(key))
			}

			getBatchResp, getBatchErr := sharedContext.Client.GetBatch(sharedContext.Ctx, &GetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Keys:      batchSetKeys,
			})
			Expect(getBatchErr).To(BeNil())
			Expect(getBatchResp).To(BeAssignableToTypeOf(responses.GetBatchSuccess{}))

			getResponses := getBatchResp.(responses.GetBatchSuccess).Results()
			Expect(len(getResponses)).To(Equal(len(batchSetKeys)))

			for _, getResp := range getResponses {
				Expect(getResp).To(BeAssignableToTypeOf(&responses.GetMiss{}))
			}
		})

		It("SetBatch happy path", func() {
			var items []BatchSetItem

			for i := 0; i < 10; i++ {
				item := BatchSetItem{
					Key:   String(fmt.Sprintf("Batch-key-%d", i)),
					Value: String(fmt.Sprintf("Batch-value-%d", i)),
				}
				items = append(items, item)
			}

			setBatchResp, setBatchErr := sharedContext.Client.SetBatch(sharedContext.Ctx, &SetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Items:     items,
			})
			Expect(setBatchErr).To(BeNil())
			Expect(setBatchResp).To(BeAssignableToTypeOf(responses.SetBatchSuccess{}))
			setResponses := setBatchResp.(responses.SetBatchSuccess).Results()
			Expect(len(setResponses)).To(Equal(len(items)))
			for _, setResp := range setResponses {
				Expect(setResp).To(BeAssignableToTypeOf(&responses.SetSuccess{}))
			}
		})

		It("SetBatch happy path with TTL", func() {
			var batchSetKeys []Value
			var items []BatchSetItem

			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("Batch-key-%d", i)
				batchSetKeys = append(batchSetKeys, String(key))
				item := BatchSetItem{
					Key:   String(key),
					Value: String(fmt.Sprintf("Batch-value-%d", i)),
				}
				items = append(items, item)
			}

			setBatchResp, setBatchErr := sharedContext.Client.SetBatch(sharedContext.Ctx, &SetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Items:     items,
				Ttl:       500 * time.Millisecond,
			})
			Expect(setBatchErr).To(BeNil())
			Expect(setBatchResp).To(BeAssignableToTypeOf(responses.SetBatchSuccess{}))
			setResponses := setBatchResp.(responses.SetBatchSuccess).Results()
			Expect(len(setResponses)).To(Equal(len(items)))
			for _, setResp := range setResponses {
				Expect(setResp).To(BeAssignableToTypeOf(&responses.SetSuccess{}))
			}

			time.Sleep(2 * time.Second)

			getBatchResp, getBatchErr := sharedContext.Client.GetBatch(sharedContext.Ctx, &GetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Keys:      batchSetKeys,
			})
			Expect(getBatchErr).To(BeNil())
			Expect(getBatchResp).To(BeAssignableToTypeOf(responses.GetBatchSuccess{}))
			getResponses := getBatchResp.(responses.GetBatchSuccess).Results()
			Expect(len(getResponses)).To(Equal(len(batchSetKeys)))
			for _, getResp := range getResponses {
				Expect(getResp).To(BeAssignableToTypeOf(&responses.GetMiss{}))
			}
		})

		It("GetBatch happy path with all hits", func() {
			var batchSetKeys []Value
			var batchSetKeysString []string
			var items []BatchSetItem

			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("All-hits-%d", i)
				batchSetKeysString = append(batchSetKeysString, key)
				batchSetKeys = append(batchSetKeys, String(key))
				item := BatchSetItem{
					Key:   String(key),
					Value: String(fmt.Sprintf("All-hits-%d", i)),
				}
				items = append(items, item)
			}

			setBatchResp, setBatchErr := sharedContext.Client.SetBatch(sharedContext.Ctx, &SetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Items:     items,
				Ttl:       10 * time.Second,
			})
			Expect(setBatchErr).To(BeNil())
			Expect(setBatchResp).To(BeAssignableToTypeOf(responses.SetBatchSuccess{}))
			setResponses := setBatchResp.(responses.SetBatchSuccess).Results()
			Expect(len(setResponses)).To(Equal(len(items)))
			for _, setResp := range setResponses {
				Expect(setResp).To(BeAssignableToTypeOf(&responses.SetSuccess{}))
			}

			getBatchResp, getBatchErr := sharedContext.Client.GetBatch(sharedContext.Ctx, &GetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Keys:      batchSetKeys,
			})
			Expect(getBatchErr).To(BeNil())
			Expect(getBatchResp).To(BeAssignableToTypeOf(responses.GetBatchSuccess{}))
			getResponses := getBatchResp.(responses.GetBatchSuccess).Results()
			Expect(len(getResponses)).To(Equal(len(batchSetKeys)))
			for _, getResp := range getResponses {
				Expect(getResp).To(BeAssignableToTypeOf(&responses.GetHit{}))
			}
			getValueMap := getBatchResp.(responses.GetBatchSuccess).ValueMap()
			for i := 0; i < len(batchSetKeysString); i++ {
				fetchedValue := getValueMap[batchSetKeysString[i]]
				Expect(fetchedValue).To(Equal(fmt.Sprintf("All-hits-%d", i)))
			}
		})

		It("GetBatch happy path with some hits and some misses", func() {
			var batchSetKeys []Value
			var batchSetKeysString []string
			var items []BatchSetItem

			for i := 0; i < 10; i++ {
				key := fmt.Sprintf("Some-hits-%d", i)

				if i < 5 {
					batchSetKeys = append(batchSetKeys, String(key))
					batchSetKeysString = append(batchSetKeysString, key)
				} else {
					differentKey := String(fmt.Sprintf("Some-hits-%d-miss", i))
					batchSetKeys = append(batchSetKeys, differentKey)
				}

				item := BatchSetItem{
					Key:   String(key),
					Value: String(fmt.Sprintf("Some-hits-%d", i)),
				}
				items = append(items, item)
			}

			setBatchResp, setBatchErr := sharedContext.Client.SetBatch(sharedContext.Ctx, &SetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Items:     items,
				Ttl:       10 * time.Second,
			})
			Expect(setBatchErr).To(BeNil())
			Expect(setBatchResp).To(BeAssignableToTypeOf(responses.SetBatchSuccess{}))
			setResponses := setBatchResp.(responses.SetBatchSuccess).Results()
			Expect(len(setResponses)).To(Equal(len(items)))
			for _, setResp := range setResponses {
				Expect(setResp).To(BeAssignableToTypeOf(&responses.SetSuccess{}))
			}

			getBatchResp, getBatchErr := sharedContext.Client.GetBatch(sharedContext.Ctx, &GetBatchRequest{
				CacheName: sharedContext.DefaultCacheName,
				Keys:      batchSetKeys,
			})
			Expect(getBatchErr).To(BeNil())
			Expect(getBatchResp).To(BeAssignableToTypeOf(responses.GetBatchSuccess{}))
			getResponses := getBatchResp.(responses.GetBatchSuccess).Results()
			Expect(len(getResponses)).To(Equal(len(batchSetKeys)))
			for i, getResp := range getResponses {
				if i < 5 {
					Expect(getResp).To(BeAssignableToTypeOf(&responses.GetHit{}))
					Expect(getResp.(*responses.GetHit).ValueString()).To(Equal(fmt.Sprintf("Some-hits-%d", i)))
				} else {
					Expect(getResp).To(BeAssignableToTypeOf(&responses.GetMiss{}))
				}
			}

			getValuesMap := getBatchResp.(responses.GetBatchSuccess).ValueMap()
			// for each key, check if the value is as expected
			for i, key := range batchSetKeysString {
				value, ok := getValuesMap[key]
				if i < 5 {
					Expect(ok).To(BeTrue())
					Expect(value).To(Equal(fmt.Sprintf("Some-hits-%d", i)))
				} else {
					Expect(ok).To(BeFalse())
				}
			}
		})
	})

	Describe("cache-client happy-path", Label(CACHE_SERVICE_LABEL), func() {
		It("creates, lists, and deletes caches", func() {
			cacheNames := []string{NewRandomString(), NewRandomString()}
			defer func() {
				for _, cacheName := range cacheNames {
					_, err := sharedContext.Client.DeleteCache(sharedContext.Ctx, &DeleteCacheRequest{CacheName: cacheName})
					if err != nil {
						panic(err)
					}
				}
			}()

			for _, cacheName := range cacheNames {
				Expect(
					sharedContext.Client.CreateCache(sharedContext.Ctx, &CreateCacheRequest{CacheName: cacheName}),
				).To(BeAssignableToTypeOf(&responses.CreateCacheSuccess{}))

				Expect(
					sharedContext.Client.CreateCache(sharedContext.Ctx, &CreateCacheRequest{CacheName: cacheName}),
				).To(BeAssignableToTypeOf(&responses.CreateCacheAlreadyExists{}))
			}

			resp, err := sharedContext.Client.ListCaches(sharedContext.Ctx, &ListCachesRequest{})
			Expect(err).To(Succeed())

			var listedCaches []string
			switch r := resp.(type) {
			case *responses.ListCachesSuccess:
				for _, info := range r.Caches() {
					listedCaches = append(listedCaches, info.Name())
				}
				Expect(listedCaches).To(ContainElements(cacheNames))
			default:
				Fail("Unexpected response type")
			}

			for _, cacheName := range cacheNames {
				Expect(
					sharedContext.Client.DeleteCache(sharedContext.Ctx, &DeleteCacheRequest{CacheName: cacheName}),
				).To(BeAssignableToTypeOf(&responses.DeleteCacheSuccess{}))
			}
			resp, err = sharedContext.Client.ListCaches(sharedContext.Ctx, &ListCachesRequest{})
			Expect(err).To(Succeed())
			Expect(resp).To(BeAssignableToTypeOf(&responses.ListCachesSuccess{}))
			switch r := resp.(type) {
			case *responses.ListCachesSuccess:
				Expect(r.Caches()).To(Not(ContainElements(cacheNames)))
			default:
				Fail("Unexpected response type")
			}
		})

		It("creates and deletes using a default cache", func() {
			// Create a separate client with a default cache name to be used only in this test
			// to avoid affecting the shared context when all tests run
			defaultCacheName := fmt.Sprintf("golang-default-%s", NewRandomString())
			clientWithDefaultCacheName, err := NewCacheClientWithDefaultCache(
				sharedContext.Configuration, sharedContext.CredentialProvider, sharedContext.DefaultTtl, defaultCacheName,
			)
			if err != nil {
				panic(err)
			}
			DeferCleanup(func() { clientWithDefaultCacheName.Close() })

			Expect(
				clientWithDefaultCacheName.CreateCache(sharedContext.Ctx, &CreateCacheRequest{}),
			).Error().NotTo(HaveOccurred())
			Expect(
				clientWithDefaultCacheName.DeleteCache(sharedContext.Ctx, &DeleteCacheRequest{}),
			).To(BeAssignableToTypeOf(&responses.DeleteCacheSuccess{}))
		})

	})

	Describe("dictionary tests", func() {
		DescribeTable("add string and bytes value for single field happy path",
			func(clientType string, field Value, value Value, expectedFieldString string, expectedFieldBytes []byte, expectedValueString string, expectedValueBytes []byte) {
				dictionaryName := NewRandomString()
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				Expect(
					client.DictionarySetField(sharedContext.Ctx, &DictionarySetFieldRequest{
						CacheName:      cacheName,
						DictionaryName: dictionaryName,
						Field:          field,
						Value:          value,
					}),
				).Error().To(BeNil())
				getFieldResp, err := client.DictionaryGetField(sharedContext.Ctx, &DictionaryGetFieldRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Field:          field,
				})
				Expect(err).To(BeNil())
				Expect(getFieldResp).To(BeAssignableToTypeOf(&responses.DictionaryGetFieldHit{}))
				switch result := getFieldResp.(type) {
				case *responses.DictionaryGetFieldHit:
					Expect(result.FieldString()).To(Equal(expectedFieldString))
					Expect(result.FieldByte()).To(Equal(expectedFieldBytes))
					Expect(result.ValueString()).To(Equal(expectedValueString))
					Expect(result.ValueByte()).To(Equal(expectedValueBytes))
				}
			},
			Entry("using string value and field", DefaultClient, String("myField"), String("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
			Entry("using string value and bytes field", DefaultClient, String("myField"), Bytes("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
			Entry("using bytes value and string field", DefaultClient, Bytes("myField"), String("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
			Entry("using bytes value and field", DefaultClient, Bytes("myField"), Bytes("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
			Entry("using string value and field with default cache", WithDefaultCache, String("myField"), String("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
			Entry("using string value and bytes field with default cache", WithDefaultCache, String("myField"), Bytes("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
			Entry("using bytes value and string field with default cache", WithDefaultCache, Bytes("myField"), String("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
			Entry("using bytes value and field with default cache", WithDefaultCache, Bytes("myField"), Bytes("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
		)

		DescribeTable("add string fields and string and bytes values for set fields happy path",
			func(clientType string, elements []DictionaryElement, expectedItemsStringValue map[string]string, expectedItemsByteValue map[string][]byte) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				dictionaryName := NewRandomString()
				Expect(
					client.DictionarySetFields(sharedContext.Ctx, &DictionarySetFieldsRequest{
						CacheName:      cacheName,
						DictionaryName: dictionaryName,
						Elements:       elements,
					}),
				).To(BeAssignableToTypeOf(&responses.DictionarySetFieldsSuccess{}))
				fetchResp, err := client.DictionaryFetch(sharedContext.Ctx, &DictionaryFetchRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
				})
				Expect(err).To(BeNil())
				switch result := fetchResp.(type) {
				case *responses.DictionaryFetchMiss:
					Fail("got a miss for a dictionary fetch that should have been a hit")
				case *responses.DictionaryFetchHit:
					i := 0
					keys := make([]string, len(result.ValueMap()))
					for k := range result.ValueMap() {
						keys[i] = k
						i++
					}
					Expect(len(result.ValueMap())).To(Equal(len(expectedItemsStringValue)))
					Expect(len(result.ValueMapStringString())).To(Equal(len(expectedItemsStringValue)))
					Expect(len(result.ValueMapStringByte())).To(Equal(len(expectedItemsByteValue)))
					for k, v := range result.ValueMap() {
						Expect(expectedItemsStringValue[k]).To(Equal(v))
					}
					for k, v := range result.ValueMapStringString() {
						Expect(expectedItemsStringValue[k]).To(Equal(v))
					}
					for k, v := range result.ValueMapStringByte() {
						Expect(expectedItemsByteValue[k]).To(Equal(v))
					}
				}
			},
			Entry(
				"with string values",
				DefaultClient,
				[]DictionaryElement{
					{Field: String("myField1"), Value: String("myValue1")},
					{Field: String("myField2"), Value: String("myValue2")},
				},
				map[string]string{"myField1": "myValue1", "myField2": "myValue2"},
				map[string][]byte{"myField1": []byte("myValue1"), "myField2": []byte("myValue2")},
			),
			Entry(
				"with byte values",
				DefaultClient,
				[]DictionaryElement{
					{Field: Bytes("myField1"), Value: Bytes("myValue1")},
					{Field: Bytes("myField2"), Value: Bytes("myValue2")},
				},
				map[string]string{"myField1": "myValue1", "myField2": "myValue2"},
				map[string][]byte{"myField1": []byte("myValue1"), "myField2": []byte("myValue2")},
			),
			Entry(
				"with mixed values",
				DefaultClient,
				[]DictionaryElement{
					{Field: Bytes("myField1"), Value: String("myValue1")},
					{Field: String("myField2"), Value: Bytes("myValue2")},
				},
				map[string]string{"myField1": "myValue1", "myField2": "myValue2"},
				map[string][]byte{"myField1": []byte("myValue1"), "myField2": []byte("myValue2")},
			),
			Entry(
				"with empty values",
				DefaultClient,
				[]DictionaryElement{
					{Field: Bytes("myField1"), Value: String("")},
					{Field: String("myField2"), Value: Bytes("")},
				},
				map[string]string{"myField1": "", "myField2": ""},
				map[string][]byte{"myField1": []byte(""), "myField2": []byte("")},
			),
			Entry(
				"with string values and default cache",
				WithDefaultCache,
				[]DictionaryElement{
					{Field: String("myField1"), Value: String("myValue1")},
					{Field: String("myField2"), Value: String("myValue2")},
				},
				map[string]string{"myField1": "myValue1", "myField2": "myValue2"},
				map[string][]byte{"myField1": []byte("myValue1"), "myField2": []byte("myValue2")},
			),
			Entry(
				"with byte values and default cache",
				WithDefaultCache,
				[]DictionaryElement{
					{Field: Bytes("myField1"), Value: Bytes("myValue1")},
					{Field: Bytes("myField2"), Value: Bytes("myValue2")},
				},
				map[string]string{"myField1": "myValue1", "myField2": "myValue2"},
				map[string][]byte{"myField1": []byte("myValue1"), "myField2": []byte("myValue2")},
			),
			Entry(
				"with mixed values and default cache",
				WithDefaultCache,
				[]DictionaryElement{
					{Field: Bytes("myField1"), Value: String("myValue1")},
					{Field: String("myField2"), Value: Bytes("myValue2")},
				},
				map[string]string{"myField1": "myValue1", "myField2": "myValue2"},
				map[string][]byte{"myField1": []byte("myValue1"), "myField2": []byte("myValue2")},
			),
			Entry(
				"with empty values and default cache",
				WithDefaultCache,
				[]DictionaryElement{
					{Field: Bytes("myField1"), Value: String("")},
					{Field: String("myField2"), Value: Bytes("")},
				},
				map[string]string{"myField1": "", "myField2": ""},
				map[string][]byte{"myField1": []byte(""), "myField2": []byte("")},
			),
		)

		DescribeTable("increments on the happy path",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				dictionaryName := NewRandomString()
				field := String("counter")
				for i := 0; i < 10; i++ {
					Expect(
						client.DictionaryIncrement(sharedContext.Ctx, &DictionaryIncrementRequest{
							CacheName:      cacheName,
							DictionaryName: dictionaryName,
							Field:          field,
							Amount:         1,
						}),
					).To(BeAssignableToTypeOf(&responses.DictionaryIncrementSuccess{}))
				}
				fetchResp, err := client.DictionaryGetField(sharedContext.Ctx, &DictionaryGetFieldRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Field:          field,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(BeAssignableToTypeOf(&responses.DictionaryGetFieldHit{}))
				switch result := fetchResp.(type) {
				case *responses.DictionaryGetFieldHit:
					Expect(result.ValueString()).To(Equal("10"))
				default:
					Fail("expected a hit for get field but got a miss")
				}
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)
		dictionaryName := NewRandomString()
		BeforeEach(func() {
			Expect(
				sharedContext.Client.DictionarySetFields(sharedContext.Ctx, &DictionarySetFieldsRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: dictionaryName,
					Elements: DictionaryElementsFromMapStringValue(
						map[string]Value{"myField1": String("myValue1"), "myField2": Bytes("myValue2")},
					),
				}),
			).To(BeAssignableToTypeOf(&responses.DictionarySetFieldsSuccess{}))
			Expect(
				sharedContext.ClientWithDefaultCacheName.DictionarySetFields(sharedContext.Ctx, &DictionarySetFieldsRequest{
					DictionaryName: dictionaryName,
					Elements: DictionaryElementsFromMapStringValue(
						map[string]Value{"myField1": String("myValue1"), "myField2": Bytes("myValue2")},
					),
				}),
			).To(BeAssignableToTypeOf(&responses.DictionarySetFieldsSuccess{}))
		})
		DescribeTable("fetches on the happy path",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				expected := map[string]string{"myField1": "myValue1", "myField2": "myValue2"}
				fetchResp, err := client.DictionaryFetch(sharedContext.Ctx, &DictionaryFetchRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(BeAssignableToTypeOf(&responses.DictionaryFetchHit{}))
				switch result := fetchResp.(type) {
				case *responses.DictionaryFetchHit:
					Expect(result.ValueMap()).To(Equal(expected))
				}
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)
		DescribeTable("add string and bytes value for single field happy path",
			func(clientType string, field Value, value Value, expectedFieldString string, expectedFieldBytes []byte, expectedValueString string, expectedValueBytes []byte) {
				dictionaryName := NewRandomString()
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				Expect(
					client.DictionarySetField(sharedContext.Ctx, &DictionarySetFieldRequest{
						CacheName:      cacheName,
						DictionaryName: dictionaryName,
						Field:          field,
						Value:          value,
					}),
				).Error().To(BeNil())
				getFieldResp, err := client.DictionaryGetField(sharedContext.Ctx, &DictionaryGetFieldRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Field:          field,
				})
				Expect(err).To(BeNil())
				Expect(getFieldResp).To(BeAssignableToTypeOf(&responses.DictionaryGetFieldHit{}))
				switch result := getFieldResp.(type) {
				case *responses.DictionaryGetFieldHit:
					Expect(result.FieldString()).To(Equal(expectedFieldString))
					Expect(result.FieldByte()).To(Equal(expectedFieldBytes))
					Expect(result.ValueString()).To(Equal(expectedValueString))
					Expect(result.ValueByte()).To(Equal(expectedValueBytes))
				}
			},
			Entry("using string value and field", DefaultClient, String("myField"), String("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
			Entry("using string value and bytes field", DefaultClient, String("myField"), Bytes("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
			Entry("using bytes value and string field", DefaultClient, Bytes("myField"), String("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
			Entry("using bytes value and field", DefaultClient, Bytes("myField"), Bytes("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
			Entry("using string value and field with default cache", WithDefaultCache, String("myField"), String("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
			Entry("using string value and bytes field with default cache", WithDefaultCache, String("myField"), Bytes("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
			Entry("using bytes value and string field with default cache", WithDefaultCache, Bytes("myField"), String("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
			Entry("using bytes value and field with default cache", WithDefaultCache, Bytes("myField"), Bytes("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
		)
	})
	Describe("list tests", func() {
		var listName string

		BeforeEach(func() {
			listName = uuid.NewString()
		})
		DescribeTable("pushing to the front of the list",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				numItems := 10
				values, expected := getValueAndExpectedValueLists(numItems)
				sort.Sort(sort.Reverse(sort.StringSlice(expected)))
				for _, value := range values {
					Expect(
						client.ListPushFront(sharedContext.Ctx, &ListPushFrontRequest{
							CacheName: cacheName,
							ListName:  listName,
							Value:     value,
						}),
					).To(BeAssignableToTypeOf(&responses.ListPushFrontSuccess{}))
				}
				fetchResp, err := client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: cacheName,
					ListName:  listName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(BeAssignableToTypeOf(&responses.ListFetchHit{}))
				Expect(fetchResp).To(HaveListLength(numItems))
				switch result := fetchResp.(type) {
				case *responses.ListFetchHit:
					Expect(result.ValueList()).To(Equal(expected))
				}
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)
		DescribeTable("pushing to the back of the list",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				numItems := 10
				values, expected := getValueAndExpectedValueLists(numItems)
				for _, value := range values {
					Expect(
						client.ListPushBack(sharedContext.Ctx, &ListPushBackRequest{
							CacheName: cacheName,
							ListName:  listName,
							Value:     value,
						}),
					).To(BeAssignableToTypeOf(&responses.ListPushBackSuccess{}))
				}

				fetchResp, err := client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: cacheName,
					ListName:  listName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveListLength(numItems))
				switch result := fetchResp.(type) {
				case *responses.ListFetchHit:
					Expect(result.ValueList()).To(Equal(expected))
				}
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)
		DescribeTable("concatenating to the front of the list",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				numItems := 10
				expected := populateList(sharedContext, listName, numItems)

				numConcatItems := 5
				concatValues, concatExpected := getValueAndExpectedValueLists(numConcatItems)
				concatResp, err := client.ListConcatenateFront(sharedContext.Ctx, &ListConcatenateFrontRequest{
					CacheName: cacheName,
					ListName:  listName,
					Values:    concatValues,
				})
				Expect(err).To(BeNil())
				Expect(concatResp).To(BeAssignableToTypeOf(&responses.ListConcatenateFrontSuccess{}))

				fetchResp, err := client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: cacheName,
					ListName:  listName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(BeAssignableToTypeOf(&responses.ListFetchHit{}))
				Expect(fetchResp).To(HaveListLength(numItems + numConcatItems))
				expected = append(concatExpected, expected...)
				switch result := fetchResp.(type) {
				case *responses.ListFetchHit:
					Expect(result.ValueList()).To(Equal(expected))
				}
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)
		DescribeTable("concatenating to the back of the list",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				numItems := 10
				expected := populateList(sharedContext, listName, numItems)

				numConcatItems := 5
				concatValues, concatExpected := getValueAndExpectedValueLists(numConcatItems)
				concatResp, err := client.ListConcatenateBack(sharedContext.Ctx, &ListConcatenateBackRequest{
					CacheName: cacheName,
					ListName:  listName,
					Values:    concatValues,
				})
				Expect(err).To(BeNil())
				Expect(concatResp).To(BeAssignableToTypeOf(&responses.ListConcatenateBackSuccess{}))

				fetchResp, err := client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: cacheName,
					ListName:  listName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(BeAssignableToTypeOf(&responses.ListFetchHit{}))
				Expect(fetchResp).To(HaveListLength(numItems + numConcatItems))
				expected = append(expected, concatExpected...)
				switch result := fetchResp.(type) {
				case *responses.ListFetchHit:
					Expect(result.ValueList()).To(Equal(expected))
				}
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)
		DescribeTable("popping from the front of the list",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				numItems := 5
				expected := populateList(sharedContext, listName, numItems)

				popResp, err := client.ListPopFront(sharedContext.Ctx, &ListPopFrontRequest{
					CacheName: cacheName,
					ListName:  listName,
				})
				Expect(err).To(BeNil())
				switch result := popResp.(type) {
				case *responses.ListPopFrontHit:
					Expect(result.ValueString()).To(Equal(string(expected[0])))
				default:
					Fail("expected a hit from list pop front but got a miss")
				}

				fetchResp, err := client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: cacheName,
					ListName:  listName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveListLength(numItems - 1))
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)
		DescribeTable("popping from the back of the list",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				numItems := 5
				expected := populateList(sharedContext, listName, numItems)

				popResp, err := client.ListPopBack(sharedContext.Ctx, &ListPopBackRequest{
					CacheName: cacheName,
					ListName:  listName,
				})
				Expect(err).To(BeNil())
				switch result := popResp.(type) {
				case *responses.ListPopBackHit:
					Expect(result.ValueString()).To(Equal(string(expected[numItems-1])))
				default:
					Fail("expected a hit from list pop front but got a miss")
				}

				fetchResp, err := client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: cacheName,
					ListName:  listName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveListLength(numItems - 1))
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)
		DescribeTable("provided no start and end to fetch fetches all results",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				numItems := 10
				values, expected := getValueAndExpectedValueLists(numItems)
				sort.Sort(sort.Reverse(sort.StringSlice(expected)))
				for _, value := range values {
					Expect(
						client.ListPushFront(sharedContext.Ctx, &ListPushFrontRequest{
							CacheName: cacheName,
							ListName:  listName,
							Value:     value,
						}),
					).To(BeAssignableToTypeOf(&responses.ListPushFrontSuccess{}))
				}
				fetchResp, err := client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: cacheName,
					ListName:  listName,
				})
				Expect(err).To(BeNil())
				Expect(fetchResp).To(BeAssignableToTypeOf(&responses.ListFetchHit{}))
				Expect(fetchResp).To(HaveListLength(numItems))
				switch result := fetchResp.(type) {
				case *responses.ListFetchHit:
					Expect(result.ValueList()).To(Equal(expected))
				}
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)
		DescribeTable("provides no start and 0 as end yields no result",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				numItems := 5
				values, _ := getValueAndExpectedValueLists(numItems)
				for _, value := range values {
					Expect(
						client.ListPushBack(sharedContext.Ctx, &ListPushBackRequest{
							CacheName: cacheName,
							ListName:  listName,
							Value:     value,
						}),
					).To(BeAssignableToTypeOf(&responses.ListPushBackSuccess{}))
				}

				endIndex := int32(0)
				fetchResp, err := client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: cacheName,
					ListName:  listName,
					EndIndex:  &endIndex,
				})

				Expect(err).To(BeNil())
				// start and end are 0; no result
				Expect(fetchResp).To(BeAssignableToTypeOf(&responses.ListFetchMiss{}))
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)
		DescribeTable("provides explicit start to list but nil end gets all results",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				numItems := 10
				values, _ := getValueAndExpectedValueLists(numItems)
				for _, value := range values {
					Expect(
						client.ListPushBack(sharedContext.Ctx, &ListPushBackRequest{
							CacheName: cacheName,
							ListName:  listName,
							Value:     value,
						}),
					).To(BeAssignableToTypeOf(&responses.ListPushBackSuccess{}))
				}

				startIndex := int32(1)
				fetchResp, err := client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName:  cacheName,
					ListName:   listName,
					StartIndex: &startIndex,
				})

				for _, vals := range fetchResp.(*responses.ListFetchHit).ValueList() {
					fmt.Println(`Values: `, vals)
				}

				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveListLength(numItems - 1))
				_, expectedVals := getValueAndExpectedValueListsRange(1, numItems)

				switch result := fetchResp.(type) {
				case *responses.ListFetchHit:
					Expect(result.ValueList()).To(Equal(expectedVals))
				}
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)

		DescribeTable("provides explicit end to list",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				numItems := 5
				values, _ := getValueAndExpectedValueLists(numItems)
				for _, value := range values {
					Expect(
						client.ListPushBack(sharedContext.Ctx, &ListPushBackRequest{
							CacheName: cacheName,
							ListName:  listName,
							Value:     value,
						}),
					).To(BeAssignableToTypeOf(&responses.ListPushBackSuccess{}))
				}

				endIndex := int32(3)
				fetchResp, err := client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName: cacheName,
					ListName:  listName,
					EndIndex:  &endIndex,
				})

				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveListLength(int(endIndex)))
				_, expectedVals := getValueAndExpectedValueListsRange(0, int(endIndex))

				switch result := fetchResp.(type) {
				case *responses.ListFetchHit:
					Expect(result.ValueList()).To(Equal(expectedVals))
				}
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)
		DescribeTable("provides explicit start and end to list",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				numItems := 5
				values, _ := getValueAndExpectedValueLists(numItems)
				for _, value := range values {
					Expect(
						client.ListPushBack(sharedContext.Ctx, &ListPushBackRequest{
							CacheName: cacheName,
							ListName:  listName,
							Value:     value,
						}),
					).To(BeAssignableToTypeOf(&responses.ListPushBackSuccess{}))
				}

				startIndex := int32(1)
				endIndex := int32(3)
				fetchResp, err := client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName:  cacheName,
					ListName:   listName,
					StartIndex: &startIndex,
					EndIndex:   &endIndex,
				})

				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveListLength(int(endIndex - startIndex)))
				_, expectedVals := getValueAndExpectedValueListsRange(int(startIndex), int(endIndex))

				switch result := fetchResp.(type) {
				case *responses.ListFetchHit:
					Expect(result.ValueList()).To(Equal(expectedVals))
				}
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)
		DescribeTable("provides negative start and unbounded end to list",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				numItems := 5
				values, _ := getValueAndExpectedValueLists(numItems)

				for _, value := range values {
					Expect(
						client.ListPushBack(sharedContext.Ctx, &ListPushBackRequest{
							CacheName: cacheName,
							ListName:  listName,
							Value:     value,
						}),
					).To(BeAssignableToTypeOf(&responses.ListPushBackSuccess{}))
				}

				startIndex := int32(-2)
				fetchResp, err := client.ListFetch(sharedContext.Ctx, &ListFetchRequest{
					CacheName:  cacheName,
					ListName:   listName,
					StartIndex: &startIndex,
				})

				Expect(err).To(BeNil())
				Expect(fetchResp).To(HaveListLength(2))
				_, expectedVals := getValueAndExpectedValueListsRange(3, 5)

				switch result := fetchResp.(type) {
				case *responses.ListFetchHit:
					Expect(result.ValueList()).To(Equal(expectedVals))
				}
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)
	})

	Describe("set tests", func() {
		var setName string

		BeforeEach(func() {
			setName = uuid.NewString()
		})
		DescribeTable("add string and byte single elements happy path",
			func(clientType string, element Value, expectedStrings []string, expectedBytes [][]byte) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				Expect(
					client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
						CacheName: cacheName,
						SetName:   setName,
						Element:   element,
					}),
				).To(BeAssignableToTypeOf(&responses.SetAddElementSuccess{}))

				fetchResp, err := client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: cacheName,
					SetName:   setName,
				})
				Expect(err).To(BeNil())
				switch result := fetchResp.(type) {
				case *responses.SetFetchHit:
					Expect(result.ValueString()).To(Equal(expectedStrings))
					Expect(result.ValueByte()).To(Equal(expectedBytes))
				default:
					fmt.Println(`Expected SetFetchHit but got:`, result)
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
						SetName:   setName,
						Elements:  elements,
					}),
				).To(BeAssignableToTypeOf(&responses.SetAddElementsSuccess{}))
				fetchResp, err := client.SetFetch(sharedContext.Ctx, &SetFetchRequest{
					CacheName: cacheName,
					SetName:   setName,
				})
				Expect(err).To(BeNil())
				switch result := fetchResp.(type) {
				case *responses.SetFetchHit:
					Expect(result.ValueString()).To(ConsistOf(expectedStrings))
					Expect(result.ValueByte()).To(ConsistOf(expectedBytes))
				default:
					Fail(fmt.Sprintf("Expected SetFetchHit, got: %T, %v", result, result))
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

	})
	Describe("sorted set tests", func() {
		var sortedSetName string

		BeforeEach(func() {
			sortedSetName = uuid.NewString()
		})
		// A convenience for adding elements to a sorted set.
		putElements := func(elements []SortedSetElement) {
			Expect(
				sharedContext.Client.SortedSetPutElements(
					sharedContext.Ctx,
					&SortedSetPutElementsRequest{
						CacheName: sharedContext.CacheName,
						SetName:   sortedSetName,
						Elements:  elements,
					},
				),
			).To(BeAssignableToTypeOf(&responses.SortedSetPutElementsSuccess{}))
			Expect(
				sharedContext.ClientWithDefaultCacheName.SortedSetPutElements(
					sharedContext.Ctx,
					&SortedSetPutElementsRequest{
						SetName:  sortedSetName,
						Elements: elements,
					},
				),
			).To(BeAssignableToTypeOf(&responses.SortedSetPutElementsSuccess{}))
		}
		DescribeTable("succeeds on the happy path",
			func(clientTYpe string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientTYpe)
				putElements(
					[]SortedSetElement{
						{Value: String("first"), Score: 9999},
						{Value: String("last"), Score: -9999},
						{Value: String("middle"), Score: 50},
					},
				)
				getResp, err := client.SortedSetGetScore(
					sharedContext.Ctx, &SortedSetGetScoreRequest{
						CacheName: cacheName,
						SetName:   sortedSetName,
						Value:     String("first"),
					},
				)
				Expect(err).To(BeNil())
				switch result := getResp.(type) {
				case *responses.SortedSetGetScoreHit:
					score := result.Score()
					Expect(score).To(Equal(9999.0))
				default:
					Fail("expected a sorted set get score hit but got a miss")
				}
			},
			Entry("with default cache", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)
	})
})
