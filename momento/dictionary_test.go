package momento_test

import (
	"reflect"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	"github.com/momentohq/client-sdk-go/utils"
)

var _ = Describe("Dictionary methods", func() {
	var sharedContext SharedContext

	BeforeEach(func() {
		sharedContext = NewSharedContext()
		sharedContext.CreateDefaultCache()
		DeferCleanup(func() {
			sharedContext.Close()
		})
	})

	DescribeTable("try using invalid cache and dictionary names",
		func(cacheName string, dictionaryName string, expectedErrorCode string) {
			Expect(
				sharedContext.Client.DictionaryFetch(sharedContext.Ctx, &DictionaryFetchRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				sharedContext.Client.DictionaryIncrement(sharedContext.Ctx, &DictionaryIncrementRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Field:          String("hi"),
					Amount:         1,
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				sharedContext.Client.DictionaryRemoveField(sharedContext.Ctx, &DictionaryRemoveFieldRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Field:          String("hi"),
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				sharedContext.Client.DictionaryRemoveFields(sharedContext.Ctx, &DictionaryRemoveFieldsRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Fields:         []Value{String("hi")},
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				sharedContext.Client.DictionaryGetFields(sharedContext.Ctx, &DictionaryGetFieldsRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Fields:         []Value{String("hi")},
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				sharedContext.Client.DictionaryGetField(sharedContext.Ctx, &DictionaryGetFieldRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Field:          String("hi"),
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				sharedContext.Client.DictionarySetField(sharedContext.Ctx, &DictionarySetFieldRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Field:          String("hi"),
					Value:          String("hi"),
					CollectionTTL:  utils.CollectionTTL{},
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				sharedContext.Client.DictionarySetFields(sharedContext.Ctx, &DictionarySetFieldsRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Items:          nil,
					CollectionTTL:  utils.CollectionTTL{},
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))
		},
		Entry("nonexistent cache name", uuid.NewString(), uuid.NewString(), NotFoundError),
		Entry("empty cache name", "", sharedContext.CollectionName, InvalidArgumentError),
		Entry("empty dictionary name", sharedContext.CacheName, "", InvalidArgumentError),
	)

	DescribeTable("add string and bytes value for single field happy path",
		func(field Value, value Value, expectedFieldString string, expectedFieldBytes []byte, expectedValueString string, expectedValueBytes []byte) {
			Expect(
				sharedContext.Client.DictionarySetField(sharedContext.Ctx, &DictionarySetFieldRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Field:          field,
					Value:          value,
				}),
			).Error().To(BeNil())
			getFieldResp, err := sharedContext.Client.DictionaryGetField(sharedContext.Ctx, &DictionaryGetFieldRequest{
				CacheName:      sharedContext.CacheName,
				DictionaryName: sharedContext.CollectionName,
				Field:          field,
			})
			Expect(err).To(BeNil())
			Expect(getFieldResp).To(BeAssignableToTypeOf(&DictionaryGetFieldHit{}))
			switch result := getFieldResp.(type) {
			case *DictionaryGetFieldHit:
				Expect(result.FieldString()).To(Equal(expectedFieldString))
				Expect(result.FieldByte()).To(Equal(expectedFieldBytes))
				Expect(result.ValueString()).To(Equal(expectedValueString))
				Expect(result.ValueByte()).To(Equal(expectedValueBytes))
			}
		},
		Entry("using string value and field", String("myField"), String("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
		Entry("using string value and bytes field", String("myField"), Bytes("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
		Entry("using bytes value and string field", Bytes("myField"), String("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
		Entry("using bytes value and field", Bytes("myField"), Bytes("myValue"), "myField", []byte("myField"), "myValue", []byte("myValue")),
	)

	It("returns an error for set field when field is empty", func() {
		Expect(
			sharedContext.Client.DictionarySetField(sharedContext.Ctx, &DictionarySetFieldRequest{
				CacheName:      sharedContext.CacheName,
				DictionaryName: sharedContext.CollectionName,
				Field:          String(""),
				Value:          String("myValue"),
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	DescribeTable("add string fields and string and bytes values for set fields happy path",
		func(items map[string]Value, expectedItemsStringValue map[string]string, expectedItemsByteValue map[string][]byte) {
			Expect(
				sharedContext.Client.DictionarySetFields(sharedContext.Ctx, &DictionarySetFieldsRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Items:          items,
				}),
			).To(BeAssignableToTypeOf(&DictionarySetFieldsSuccess{}))
			fetchResp, err := sharedContext.Client.DictionaryFetch(sharedContext.Ctx, &DictionaryFetchRequest{
				CacheName:      sharedContext.CacheName,
				DictionaryName: sharedContext.CollectionName,
			})
			Expect(err).To(BeNil())
			switch result := fetchResp.(type) {
			case *DictionaryFetchMiss:
				Fail("got a miss for a dictionary fetch that should have been a hit")
			case *DictionaryFetchHit:
				Expect(reflect.DeepEqual(result.ValueMap(), expectedItemsStringValue)).To(BeTrue())
				Expect(reflect.DeepEqual(result.ValueMapStringString(), expectedItemsStringValue)).To(BeTrue())
				Expect(reflect.DeepEqual(result.ValueMapStringByte(), expectedItemsByteValue)).To(BeTrue())
			}
		},
		Entry(
			"with string values",
			map[string]Value{"myField1": String("myValue1"), "myField2": String("myValue2")},
			map[string]string{"myField1": "myValue1", "myField2": "myValue2"},
			map[string][]byte{"myField1": []byte("myValue1"), "myField2": []byte("myValue2")},
		),
		Entry(
			"with byte values",
			map[string]Value{"myField1": Bytes("myValue1"), "myField2": Bytes("myValue2")},
			map[string]string{"myField1": "myValue1", "myField2": "myValue2"},
			map[string][]byte{"myField1": []byte("myValue1"), "myField2": []byte("myValue2")},
		),
		Entry(
			"with mixed values",
			map[string]Value{"myField1": String("myValue1"), "myField2": Bytes("myValue2")},
			map[string]string{"myField1": "myValue1", "myField2": "myValue2"},
			map[string][]byte{"myField1": []byte("myValue1"), "myField2": []byte("myValue2")},
		),
	)

	It("returns an error if an item field is empty", func() {
		Expect(
			sharedContext.Client.DictionarySetFields(sharedContext.Ctx, &DictionarySetFieldsRequest{
				CacheName:      sharedContext.CacheName,
				DictionaryName: sharedContext.CollectionName,
				Items:          map[string]Value{"myField": String("myValue"), "": String("myOtherValue")},
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	Describe("dictionary increment", func() {

		It("populates nonexistent field", func() {
			incrResp, err := sharedContext.Client.DictionaryIncrement(sharedContext.Ctx, &DictionaryIncrementRequest{
				CacheName:      sharedContext.CacheName,
				DictionaryName: sharedContext.CollectionName,
				Field:          String("myField"),
				Amount:         3,
			})
			Expect(err).To(BeNil())
			Expect(incrResp).To(BeAssignableToTypeOf(&DictionaryIncrementSuccess{}))
			switch result := incrResp.(type) {
			case *DictionaryIncrementSuccess:
				Expect(result.Value()).To(Equal(int64(3)))
			}
		})

		It("returns an error when called on a non-integer field", func() {
			Expect(
				sharedContext.Client.DictionarySetField(sharedContext.Ctx, &DictionarySetFieldRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Field:          String("notacounter"),
					Value:          String("notanumber"),
				}),
			).To(BeAssignableToTypeOf(&DictionarySetFieldSuccess{}))

			Expect(
				sharedContext.Client.DictionaryIncrement(sharedContext.Ctx, &DictionaryIncrementRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Field:          String("notacounter"),
					Amount:         1,
				}),
			).Error().To(HaveMomentoErrorCode(FailedPreconditionError))
		})

		It("returns an error when amount is zero", func() {
			_, err := sharedContext.Client.DictionaryIncrement(sharedContext.Ctx, &DictionaryIncrementRequest{
				CacheName:      sharedContext.CacheName,
				DictionaryName: sharedContext.CollectionName,
				Field:          String("myField"),
				Amount:         0,
			})
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("increments on the happy path", func() {
			field := String("counter")
			for i := 0; i < 10; i++ {
				Expect(
					sharedContext.Client.DictionaryIncrement(sharedContext.Ctx, &DictionaryIncrementRequest{
						CacheName:      sharedContext.CacheName,
						DictionaryName: sharedContext.CollectionName,
						Field:          field,
						Amount:         1,
					}),
				).To(BeAssignableToTypeOf(&DictionaryIncrementSuccess{}))
			}
			fetchResp, err := sharedContext.Client.DictionaryGetField(sharedContext.Ctx, &DictionaryGetFieldRequest{
				CacheName:      sharedContext.CacheName,
				DictionaryName: sharedContext.CollectionName,
				Field:          field,
			})
			Expect(err).To(BeNil())
			Expect(fetchResp).To(BeAssignableToTypeOf(&DictionaryGetFieldHit{}))
			switch result := fetchResp.(type) {
			case *DictionaryGetFieldHit:
				Expect(result.ValueString()).To(Equal("10"))
			default:
				Fail("expected a hit for get field but got a miss")
			}
		})
	})

	Describe("dictionary get", func() {

		BeforeEach(func() {
			Expect(
				sharedContext.Client.DictionarySetFields(sharedContext.Ctx, &DictionarySetFieldsRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Items:          map[string]Value{"myField1": String("myValue1"), "myField2": Bytes("myValue2")},
				}),
			).To(BeAssignableToTypeOf(&DictionarySetFieldsSuccess{}))
		})

		When("getting single field", func() {

			It("returns the correct string and byte values", func() {
				expected := map[string]string{"myField1": "myValue1", "myField2": "myValue2"}

				for fieldName, valueStr := range expected {
					getResp, err := sharedContext.Client.DictionaryGetField(sharedContext.Ctx, &DictionaryGetFieldRequest{
						CacheName:      sharedContext.CacheName,
						DictionaryName: sharedContext.CollectionName,
						Field:          String(fieldName),
					})
					Expect(err).To(BeNil())
					Expect(getResp).To(BeAssignableToTypeOf(&DictionaryGetFieldHit{}))
					switch result := getResp.(type) {
					case *DictionaryGetFieldHit:
						Expect(result.FieldString()).To(Equal(fieldName))
						Expect(result.FieldByte()).To(Equal([]byte(fieldName)))
						Expect(result.ValueString()).To(Equal(valueStr))
						Expect(result.ValueByte()).To(Equal([]byte(valueStr)))
					default:
						Fail("something really weird happened")
					}
				}
			})

			It("returns a miss for a nonexistent field", func() {
				getResp, err := sharedContext.Client.DictionaryGetField(sharedContext.Ctx, &DictionaryGetFieldRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Field:          String("idontexist"),
				})
				Expect(err).To(BeNil())
				Expect(getResp).To(BeAssignableToTypeOf(&DictionaryGetFieldMiss{}))
			})

			It("returns a miss for a nonexistent dictionary", func() {
				getResp, err := sharedContext.Client.DictionaryGetField(sharedContext.Ctx, &DictionaryGetFieldRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: uuid.NewString(),
					Field:          String("idontexist"),
				})
				Expect(err).To(BeNil())
				Expect(getResp).To(BeAssignableToTypeOf(&DictionaryGetFieldMiss{}))
			})

		})

		When("getting multiple fields", func() {

			It("returns the correct string and byte values", func() {
				getResp, err := sharedContext.Client.DictionaryGetFields(sharedContext.Ctx, &DictionaryGetFieldsRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Fields:         []Value{String("myField1"), String("myField2")},
				})
				Expect(err).To(BeNil())
				Expect(getResp).To(BeAssignableToTypeOf(&DictionaryGetFieldsHit{}))

				expectedStrings := map[string]string{"myField1": "myValue1", "myField2": "myValue2"}
				expectedBytes := map[string][]byte{"myField1": []byte("myValue1"), "myField2": []byte("myValue2")}
				switch result := getResp.(type) {
				case *DictionaryGetFieldsHit:
					Expect(result.ValueMapStringString()).To(Equal(expectedStrings))
					Expect(result.ValueMap()).To(Equal(expectedStrings))
					Expect(result.ValueMapStringBytes()).To(Equal(expectedBytes))
				}
			})

			It("returns a miss for nonexistent dictionary", func() {
				Expect(
					sharedContext.Client.DictionaryGetFields(sharedContext.Ctx, &DictionaryGetFieldsRequest{
						CacheName:      sharedContext.CacheName,
						DictionaryName: uuid.NewString(),
						Fields:         []Value{String("myField1")},
					}),
				).To(BeAssignableToTypeOf(&DictionaryGetFieldsMiss{}))
			})

			It("returns misses for nonexistent fields", func() {
				getResp, err := sharedContext.Client.DictionaryGetFields(sharedContext.Ctx, &DictionaryGetFieldsRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Fields:         []Value{String("bogusField1"), String("bogusField2")},
				})
				Expect(err).To(BeNil())
				Expect(getResp).To(BeAssignableToTypeOf(&DictionaryGetFieldsHit{}))
				switch result := getResp.(type) {
				case *DictionaryGetFieldsHit:
					Expect(result.ValueMap()).To(BeEmpty())
					for _, value := range result.Responses() {
						switch value.(type) {
						case *DictionaryGetFieldHit:
							Fail("got a hit response for a field that should return a miss")
						}
					}
				}
			})

			It("filters missing fields out of response value maps", func() {
				getResp, err := sharedContext.Client.DictionaryGetFields(sharedContext.Ctx, &DictionaryGetFieldsRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Fields:         []Value{String("bogusField1"), String("myField2")},
				})
				Expect(err).To(BeNil())
				Expect(getResp).To(BeAssignableToTypeOf(&DictionaryGetFieldsHit{}))
				switch result := getResp.(type) {
				case *DictionaryGetFieldsHit:
					Expect(result.ValueMap()).To(Equal(map[string]string{"myField2": "myValue2"}))
					Expect(len(result.Responses())).To(Equal(2))
				}
			})

		})
	})

	Describe("dictionary fetch", func() {

		BeforeEach(func() {
			Expect(
				sharedContext.Client.DictionarySetFields(sharedContext.Ctx, &DictionarySetFieldsRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Items:          map[string]Value{"myField1": String("myValue1"), "myField2": Bytes("myValue2")},
				}),
			).To(BeAssignableToTypeOf(&DictionarySetFieldsSuccess{}))
		})

		It("fetches on the happy path", func() {
			expected := map[string]string{"myField1": "myValue1", "myField2": "myValue2"}
			fetchResp, err := sharedContext.Client.DictionaryFetch(sharedContext.Ctx, &DictionaryFetchRequest{
				CacheName:      sharedContext.CacheName,
				DictionaryName: sharedContext.CollectionName,
			})
			Expect(err).To(BeNil())
			Expect(fetchResp).To(BeAssignableToTypeOf(&DictionaryFetchHit{}))
			switch result := fetchResp.(type) {
			case *DictionaryFetchHit:
				Expect(result.ValueMap()).To(Equal(expected))
			}
		})

		It("returns a miss for nonexistent dictionary", func() {
			Expect(
				sharedContext.Client.DictionaryFetch(sharedContext.Ctx, &DictionaryFetchRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: uuid.NewString(),
				}),
			).To(BeAssignableToTypeOf(&DictionaryFetchMiss{}))
		})

	})

	Describe("dictionary remove", func() {

		BeforeEach(func() {
			Expect(
				sharedContext.Client.DictionarySetFields(sharedContext.Ctx, &DictionarySetFieldsRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Items: map[string]Value{
						"myField1": String("myValue1"),
						"myField2": Bytes("myValue2"),
						"myField3": String("myValue3"),
					},
				}),
			).To(BeAssignableToTypeOf(&DictionarySetFieldsSuccess{}))
		})

		When("removing a single field", func() {

			It("properly removes a field", func() {
				removeResp, err := sharedContext.Client.DictionaryRemoveField(sharedContext.Ctx, &DictionaryRemoveFieldRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Field:          String("myField1"),
				})
				Expect(err).To(BeNil())
				Expect(removeResp).To(BeAssignableToTypeOf(&DictionaryRemoveFieldSuccess{}))

				fetchResp, err := sharedContext.Client.DictionaryFetch(sharedContext.Ctx, &DictionaryFetchRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				switch result := fetchResp.(type) {
				case *DictionaryFetchHit:
					Expect(result.ValueMap()).To(Equal(map[string]string{
						"myField2": "myValue2",
						"myField3": "myValue3",
					}))
				default:
					Fail("expected a hit from dictionary fetch but got a miss")
				}
			})

			It("no-ops when attempting to remove a nonexistent field", func() {
				removeResp, err := sharedContext.Client.DictionaryRemoveField(sharedContext.Ctx, &DictionaryRemoveFieldRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Field:          String("bogusField1"),
				})
				Expect(err).To(BeNil())
				Expect(removeResp).To(BeAssignableToTypeOf(&DictionaryRemoveFieldSuccess{}))
			})

			It("no-ops when using a nonexistent dictionary", func() {
				removeResp, err := sharedContext.Client.DictionaryRemoveField(sharedContext.Ctx, &DictionaryRemoveFieldRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: uuid.NewString(),
					Field:          String("bogusField1"),
				})
				Expect(err).To(BeNil())
				Expect(removeResp).To(BeAssignableToTypeOf(&DictionaryRemoveFieldSuccess{}))
			})

		})

		When("removing multiple fields", func() {

			It("properly removes multiple fields", func() {
				removeResp, err := sharedContext.Client.DictionaryRemoveFields(sharedContext.Ctx, &DictionaryRemoveFieldsRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Fields:         []Value{String("myField1"), Bytes("myField2")},
				})
				Expect(err).To(BeNil())
				Expect(removeResp).To(BeAssignableToTypeOf(&DictionaryRemoveFieldsSuccess{}))

				fetchResp, err := sharedContext.Client.DictionaryFetch(sharedContext.Ctx, &DictionaryFetchRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
				})
				Expect(err).To(BeNil())
				switch result := fetchResp.(type) {
				case *DictionaryFetchHit:
					Expect(result.ValueMap()).To(Equal(map[string]string{
						"myField3": "myValue3",
					}))
				default:
					Fail("expected a hit from dictionary fetch but got a miss")
				}
			})

			It("no-ops when attempting to remove a nonexistent field", func() {
				removeResp, err := sharedContext.Client.DictionaryRemoveFields(sharedContext.Ctx, &DictionaryRemoveFieldsRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Fields:         []Value{String("bogusField1"), Bytes("bogusField2")},
				})
				Expect(err).To(BeNil())
				Expect(removeResp).To(BeAssignableToTypeOf(&DictionaryRemoveFieldsSuccess{}))
			})

			It("no-ops when using a nonexistent dictionary", func() {
				removeResp, err := sharedContext.Client.DictionaryRemoveFields(sharedContext.Ctx, &DictionaryRemoveFieldsRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: uuid.NewString(),
					Fields:         []Value{String("bogusField1"), Bytes("bogusField2")},
				})
				Expect(err).To(BeNil())
				Expect(removeResp).To(BeAssignableToTypeOf(&DictionaryRemoveFieldsSuccess{}))
			})

		})
	})

	Describe("client TTL", func() {

		When("client TTL is exceeded", func() {

			It("returns a miss for the collection", func() {
				Expect(
					sharedContext.Client.DictionarySetFields(sharedContext.Ctx, &DictionarySetFieldsRequest{
						CacheName:      sharedContext.CacheName,
						DictionaryName: sharedContext.CollectionName,
						Items:          map[string]Value{"myField1": String("myValue1"), "myField2": String("myValue2")},
					}),
				).Error().To(BeNil())

				Expect(
					sharedContext.Client.DictionaryFetch(sharedContext.Ctx, &DictionaryFetchRequest{
						CacheName:      sharedContext.CacheName,
						DictionaryName: sharedContext.CollectionName,
					}),
				).To(BeAssignableToTypeOf(&DictionaryFetchHit{}))

				time.Sleep(sharedContext.DefaultTTL)

				Expect(
					sharedContext.Client.DictionaryFetch(sharedContext.Ctx, &DictionaryFetchRequest{
						CacheName:      sharedContext.CacheName,
						DictionaryName: sharedContext.CollectionName,
					}),
				).To(BeAssignableToTypeOf(&DictionaryFetchMiss{}))
			})

		})

	})

	Describe("collection TTL", func() {

		BeforeEach(func() {
			Expect(
				sharedContext.Client.DictionarySetFields(sharedContext.Ctx, &DictionarySetFieldsRequest{
					CacheName:      sharedContext.CacheName,
					DictionaryName: sharedContext.CollectionName,
					Items:          map[string]Value{"myField1": String("myValue1"), "myField2": String("myValue2")},
				}),
			).Error().To(BeNil())
		})

		When("collection TTL is empty", func() {

			It("will have a false refreshTTL and fetch will miss after client default ttl", func() {
				time.Sleep(sharedContext.DefaultTTL / 2)
				Expect(
					sharedContext.Client.DictionarySetField(sharedContext.Ctx, &DictionarySetFieldRequest{
						CacheName:      sharedContext.CacheName,
						DictionaryName: sharedContext.CollectionName,
						Field:          String("foo"),
						Value:          String("bar"),
						CollectionTTL:  utils.CollectionTTL{},
					}),
				).To(BeAssignableToTypeOf(&DictionarySetFieldSuccess{}))

				time.Sleep(sharedContext.DefaultTTL / 2)

				Expect(
					sharedContext.Client.DictionaryFetch(sharedContext.Ctx, &DictionaryFetchRequest{
						CacheName:      sharedContext.CacheName,
						DictionaryName: sharedContext.CollectionName,
					}),
				).To(BeAssignableToTypeOf(&DictionaryFetchMiss{}))
			})

		})

		When("collection TTL is configured", func() {

			It("is ignored if refresh ttl is false", func() {
				Expect(
					sharedContext.Client.DictionarySetField(sharedContext.Ctx, &DictionarySetFieldRequest{
						CacheName:      sharedContext.CacheName,
						DictionaryName: sharedContext.CollectionName,
						Field:          String("myField3"),
						Value:          String("myValue3"),
						CollectionTTL: utils.CollectionTTL{
							Ttl:        sharedContext.DefaultTTL + time.Second*60,
							RefreshTtl: false,
						},
					}),
				).To(BeAssignableToTypeOf(&DictionarySetFieldSuccess{}))

				time.Sleep(sharedContext.DefaultTTL)

				Expect(
					sharedContext.Client.DictionaryFetch(sharedContext.Ctx, &DictionaryFetchRequest{
						CacheName:      sharedContext.CacheName,
						DictionaryName: sharedContext.CollectionName,
					}),
				).To(BeAssignableToTypeOf(&DictionaryFetchMiss{}))
			})

			It("is respected if refresh TTL is true", func() {
				Expect(
					sharedContext.Client.DictionarySetField(sharedContext.Ctx, &DictionarySetFieldRequest{
						CacheName:      sharedContext.CacheName,
						DictionaryName: sharedContext.CollectionName,
						Field:          String("myField3"),
						Value:          String("myValue3"),
						CollectionTTL: utils.CollectionTTL{
							Ttl:        sharedContext.DefaultTTL + time.Second*60,
							RefreshTtl: true,
						},
					}),
				).To(BeAssignableToTypeOf(&DictionarySetFieldSuccess{}))

				time.Sleep(sharedContext.DefaultTTL)

				Expect(
					sharedContext.Client.DictionaryFetch(sharedContext.Ctx, &DictionaryFetchRequest{
						CacheName:      sharedContext.CacheName,
						DictionaryName: sharedContext.CollectionName,
					}),
				).To(BeAssignableToTypeOf(&DictionaryFetchHit{}))
			})

		})

	})

})
