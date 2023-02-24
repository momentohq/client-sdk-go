package momento_test

import (
	"context"
	"reflect"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	. "github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/utils"
)

var _ = Describe("Dictionary methods", func() {
	var clientProps SimpleCacheClientProps
	var credentialProvider auth.CredentialProvider
	var configuration config.Configuration
	var client SimpleCacheClient
	var defaultTTL time.Duration
	var testCacheName string
	var testDictionaryName string
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		credentialProvider, err = auth.NewEnvMomentoTokenProvider("TEST_AUTH_TOKEN")
		if err != nil {
			panic(err)
		}
		configuration = config.LatestLaptopConfig()
		defaultTTL = 3 * time.Second

		clientProps = SimpleCacheClientProps{
			CredentialProvider: credentialProvider,
			Configuration:      configuration,
			DefaultTTL:         defaultTTL,
		}

		client, err = NewSimpleCacheClient(&clientProps)
		if err != nil {
			panic(err)
		}
		DeferCleanup(func() { client.Close() })

		testCacheName = uuid.NewString()
		testDictionaryName = uuid.NewString()
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

	DescribeTable("try using invalid cache and dictionary names",
		func(cacheName string, dictionaryName string, expectedErrorCode string) {
			Expect(
				client.DictionaryFetch(ctx, &DictionaryFetchRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				client.DictionaryIncrement(ctx, &DictionaryIncrementRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Field:          String("hi"),
					Amount:         1,
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				client.DictionaryRemoveField(ctx, &DictionaryRemoveFieldRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Field:          String("hi"),
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				client.DictionaryRemoveFields(ctx, &DictionaryRemoveFieldsRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Fields:         []Value{String("hi")},
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				client.DictionaryGetFields(ctx, &DictionaryGetFieldsRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Fields:         []Value{String("hi")},
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				client.DictionaryGetField(ctx, &DictionaryGetFieldRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Field:          String("hi"),
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				client.DictionarySetField(ctx, &DictionarySetFieldRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Field:          String("hi"),
					Value:          String("hi"),
					CollectionTTL:  utils.CollectionTTL{},
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))

			Expect(
				client.DictionarySetFields(ctx, &DictionarySetFieldsRequest{
					CacheName:      cacheName,
					DictionaryName: dictionaryName,
					Items:          nil,
					CollectionTTL:  utils.CollectionTTL{},
				}),
			).Error().To(HaveMomentoErrorCode(expectedErrorCode))
		},
		Entry("nonexistent cache name", uuid.NewString(), uuid.NewString(), NotFoundError),
		Entry("empty cache name", "", testDictionaryName, InvalidArgumentError),
		Entry("empty dictionary name", testCacheName, "", InvalidArgumentError),
	)

	DescribeTable("add string and bytes value for single field happy path",
		func(field Value, value Value, expectedFieldString string, expectedFieldBytes []byte, expectedValueString string, expectedValueBytes []byte) {
			Expect(
				client.DictionarySetField(ctx, &DictionarySetFieldRequest{
					CacheName:      testCacheName,
					DictionaryName: testDictionaryName,
					Field:          field,
					Value:          value,
				}),
			).Error().To(BeNil())
			getFieldResp, err := client.DictionaryGetField(ctx, &DictionaryGetFieldRequest{
				CacheName:      testCacheName,
				DictionaryName: testDictionaryName,
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

	It("returns an error when field is empty", func() {
		Expect(
			client.DictionarySetField(ctx, &DictionarySetFieldRequest{
				CacheName:      testCacheName,
				DictionaryName: testDictionaryName,
				Field:          String(""),
				Value:          String("myValue"),
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	DescribeTable("add string fields and string and bytes values for set fields happy path",
		func(items map[string]Value, expectedItemsStringValue map[string]string, expectedItemsByteValue map[string][]byte) {
			Expect(
				client.DictionarySetFields(ctx, &DictionarySetFieldsRequest{
					CacheName:      testCacheName,
					DictionaryName: testDictionaryName,
					Items:          items,
				}),
			).To(BeAssignableToTypeOf(&DictionarySetFieldsSuccess{}))
			fetchResp, err := client.DictionaryFetch(ctx, &DictionaryFetchRequest{
				CacheName:      testCacheName,
				DictionaryName: testDictionaryName,
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
			client.DictionarySetFields(ctx, &DictionarySetFieldsRequest{
				CacheName:      testCacheName,
				DictionaryName: testDictionaryName,
				Items:          map[string]Value{"myField": String("myValue"), "": String("myOtherValue")},
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	Describe("dictionary increment", func() {

		It("populates nonexistent field", func() {
			incrResp, err := client.DictionaryIncrement(ctx, &DictionaryIncrementRequest{
				CacheName:      testCacheName,
				DictionaryName: testDictionaryName,
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
				client.DictionarySetField(ctx, &DictionarySetFieldRequest{
					CacheName:      testCacheName,
					DictionaryName: testDictionaryName,
					Field:          String("notacounter"),
					Value:          String("notanumber"),
				}),
			).To(BeAssignableToTypeOf(&DictionarySetFieldSuccess{}))

			// TODO: update error code when additional error codes PR goes through
			Expect(
				client.DictionaryIncrement(ctx, &DictionaryIncrementRequest{
					CacheName:      testCacheName,
					DictionaryName: testDictionaryName,
					Field:          String("notacounter"),
					Amount:         1,
				}),
			).Error().To(HaveMomentoErrorCode(FailedPreconditionError))
		})

		It("returns an error when amount is zero", func() {
			_, err := client.DictionaryIncrement(ctx, &DictionaryIncrementRequest{
				CacheName:      testCacheName,
				DictionaryName: testDictionaryName,
				Field:          String("myField"),
				Amount:         0,
			})
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("increments on the happy path", func() {
			field := String("counter")
			for i := 0; i < 10; i++ {
				Expect(
					client.DictionaryIncrement(ctx, &DictionaryIncrementRequest{
						CacheName:      testCacheName,
						DictionaryName: testDictionaryName,
						Field:          field,
						Amount:         1,
					}),
				).To(BeAssignableToTypeOf(&DictionaryIncrementSuccess{}))
			}
			fetchResp, err := client.DictionaryGetField(ctx, &DictionaryGetFieldRequest{
				CacheName:      testCacheName,
				DictionaryName: testDictionaryName,
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

		When("getting single field", func() {

			BeforeEach(func() {
				Expect(
					client.DictionarySetFields(ctx, &DictionarySetFieldsRequest{
						CacheName:      testCacheName,
						DictionaryName: testDictionaryName,
						Items:          map[string]Value{"myField1": String("myValue1"), "myField2": Bytes("myValue2")},
					}),
				).To(BeAssignableToTypeOf(&DictionarySetFieldsSuccess{}))
			})

			It("returns the correct string and byte values", func() {
				expected := map[string]string{"myField1": "myValue1", "myField2": "myValue2"}

				for fieldName, valueStr := range expected {
					getResp, err := client.DictionaryGetField(ctx, &DictionaryGetFieldRequest{
						CacheName:      testCacheName,
						DictionaryName: testDictionaryName,
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
				getResp, err := client.DictionaryGetField(ctx, &DictionaryGetFieldRequest{
					CacheName:      testCacheName,
					DictionaryName: testDictionaryName,
					Field:          String("idontexist"),
				})
				Expect(err).To(BeNil())
				Expect(getResp).To(BeAssignableToTypeOf(&DictionaryGetFieldMiss{}))
			})

			It("returns a miss for a nonexistent dictionary", func() {
				getResp, err := client.DictionaryGetField(ctx, &DictionaryGetFieldRequest{
					CacheName:      testCacheName,
					DictionaryName: uuid.NewString(),
					Field:          String("idontexist"),
				})
				Expect(err).To(BeNil())
				Expect(getResp).To(BeAssignableToTypeOf(&DictionaryGetFieldMiss{}))
			})

		})

		When("getting multiple fields", func() {

			BeforeEach(func() {
				Expect(
					client.DictionarySetFields(ctx, &DictionarySetFieldsRequest{
						CacheName:      testCacheName,
						DictionaryName: testDictionaryName,
						Items:          map[string]Value{"myField1": String("myValue1"), "myField2": Bytes("myValue2")},
					}),
				).To(BeAssignableToTypeOf(&DictionarySetFieldsSuccess{}))
			})

			It("returns the correct string and byte values", func() {
				
			})

		})
	})

	Describe("dictionary fetch", func() {

	})

})
