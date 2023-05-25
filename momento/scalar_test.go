package momento_test

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	. "github.com/momentohq/client-sdk-go/responses"
)

var _ = Describe("Scalar methods", func() {
	var sharedContext SharedContext
	BeforeEach(func() {
		sharedContext = NewSharedContext()
		sharedContext.CreateDefaultCaches()

		DeferCleanup(func() { sharedContext.Close() })
	})

	DescribeTable("Gets, Sets, and Deletes",
		func(clientType string, key Key, value Value, expectedString string, expectedBytes []byte) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			Expect(
				client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: cacheName,
					Key:       key,
					Value:     value,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			getResp, err := client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())

			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
			default:
				Fail("Unexpected type from Get")
			}

			Expect(
				client.Delete(sharedContext.Ctx, &DeleteRequest{
					CacheName: cacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&DeleteSuccess{}))

			Expect(
				client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: cacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		},
		Entry("when the key and value are strings", DefaultClient, String("key"), String("value"), "value", []byte("value")),
		Entry("when the key and value are bytes", DefaultClient, Bytes([]byte{1, 2, 3}), Bytes("string"), "string", []byte("string")),
		Entry("when the value is empty", DefaultClient, String("key"), String(""), "", []byte("")),
		Entry("when the value is blank", DefaultClient, String("key"), String("  "), "  ", []byte("  ")),
		Entry("with default cache name when the key and value are strings", WithDefaultCache, String("key"), String("value"), "value", []byte("value")),
		Entry("with default cache name when the key and value are bytes", WithDefaultCache, Bytes([]byte{1, 2, 3}), Bytes("string"), "string", []byte("string")),
		Entry("with default cache name when the value is empty", WithDefaultCache, String("key"), String(""), "", []byte("")),
		Entry("with default cache name when the value is blank", WithDefaultCache, String("key"), String("  "), "  ", []byte("  ")),
	)

	DescribeTable("errors when the cache is missing",
		func(clientType string) {
			client, _ := sharedContext.GetClientPrereqsForType(clientType)
			cacheName := uuid.NewString()
			key := String("key")
			value := String("value")

			getResp, err := client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(getResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(NotFoundError))

			setResp, err := client.Set(sharedContext.Ctx, &SetRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(setResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(NotFoundError))

			deleteResp, err := client.Delete(sharedContext.Ctx, &DeleteRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(deleteResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(NotFoundError))
		},
		Entry("with default client", DefaultClient),
		Entry("with client with default cache", WithDefaultCache),
	)

	DescribeTable("errors when the key is nil",
		func(clientType string) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			Expect(
				client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: cacheName,
					Key:       nil,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
			Expect(
				client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: cacheName,
					Key:       nil,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
			Expect(
				client.Delete(sharedContext.Ctx, &DeleteRequest{
					CacheName: cacheName,
					Key:       nil,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		},
		Entry("with default client", DefaultClient),
		Entry("with client with default cache", WithDefaultCache),
	)

	DescribeTable("returns a miss when the key doesn't exist",
		func(clientType string) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			Expect(
				client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: cacheName,
					Key:       String(uuid.NewString()),
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		},
		Entry("with default client", DefaultClient),
		Entry("with client with default cache", WithDefaultCache),
	)

	DescribeTable("invalid cache names and keys",
		func(clientType string, cacheName string, key Key, value Key) {
			client, _ := sharedContext.GetClientPrereqsForType(clientType)
			getResp, err := client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(getResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

			setResp, err := client.Set(sharedContext.Ctx, &SetRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(setResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

			deleteResp, err := client.Delete(sharedContext.Ctx, &DeleteRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(deleteResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
		},
		Entry("With default client and an empty cache name", DefaultClient, "", String("key"), String("value")),
		Entry("With default client and  an bank cache name", DefaultClient, "   ", String("key"), String("value")),
		Entry("With default client and  an empty key", DefaultClient, uuid.NewString(), String(""), String("value")),
		Entry("With client with default cache and an bank cache name", WithDefaultCache, "   ", String("key"), String("value")),
		Entry("With client with default cache and an empty key", WithDefaultCache, uuid.NewString(), String(""), String("value")),
	)

	Describe("Set", func() {
		It("Uses the default TTL", func() {
			key := String("key")
			value := String("value")

			Expect(
				sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Value:     value,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			time.Sleep(sharedContext.DefaultTtl / 2)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetHit{}))

			time.Sleep(sharedContext.DefaultTtl)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		})

		It("Overrides the default TTL", func() {
			key := String("key")
			value := String("value")

			Expect(
				sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Value:     value,
					Ttl:       sharedContext.DefaultTtl * 2,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			time.Sleep(sharedContext.DefaultTtl / 2)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetHit{}))

			time.Sleep(sharedContext.DefaultTtl)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetHit{}))

			time.Sleep(sharedContext.DefaultTtl)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		})

		It("returns an error for a nil value", func() {
			key := String("key")
			Expect(
				sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Value:     nil,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})
	})

	Describe("Keys exist", func() {
		BeforeEach(func() {
			for i := 0; i < 3; i++ {
				strKey := fmt.Sprintf("#%d", i)
				strVal := fmt.Sprintf("%sValue", strKey)
				Expect(
					sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
						CacheName: sharedContext.CacheName,
						Key:       String(strKey),
						Value:     String(strVal),
					}),
				).To(BeAssignableToTypeOf(&SetSuccess{}))
				Expect(
					sharedContext.ClientWithDefaultCacheName.Set(sharedContext.Ctx, &SetRequest{
						Key:   String(strKey),
						Value: String(strVal),
					}),
				).To(BeAssignableToTypeOf(&SetSuccess{}))
			}
		})

		DescribeTable("check for valid keys exist results",
			func(clientType string, toCheck []Key, expected []bool) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				resp, err := client.KeysExist(sharedContext.Ctx, &KeysExistRequest{
					CacheName: cacheName,
					Keys:      toCheck,
				})
				Expect(err).To(BeNil())
				switch result := resp.(type) {
				case *KeysExistSuccess:
					Expect(result.Exists()).To(Equal(expected))
				default:
					Fail(fmt.Sprintf("expected keys exist success but got %s", result))
				}
			},
			Entry("all hits", DefaultClient, []Key{String("#1"), String("#2")}, []bool{true, true}),
			Entry("all misses", DefaultClient, []Key{String("nope"), String("stillnope")}, []bool{false, false}),
			Entry("mixed", DefaultClient, []Key{String("nope"), String("#1")}, []bool{false, true}),
			Entry("all hits with default cache", WithDefaultCache, []Key{String("#1"), String("#2")}, []bool{true, true}),
			Entry("all misses with default cache", WithDefaultCache, []Key{String("nope"), String("stillnope")}, []bool{false, false}),
			Entry("mixed with default cache", WithDefaultCache, []Key{String("nope"), String("#1")}, []bool{false, true}),
		)
	})

	Describe("UpdateTtl", func() {
		DescribeTable("Overwrites Ttl",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				key := String("key")
				value := String("value")

				Expect(
					client.Set(sharedContext.Ctx, &SetRequest{
						CacheName: cacheName,
						Key:       key,
						Value:     value,
					}),
				).To(BeAssignableToTypeOf(&SetSuccess{}))

				Expect(
					client.UpdateTtl(sharedContext.Ctx, &UpdateTtlRequest{
						CacheName: cacheName,
						Key:       key,
						Ttl:       6 * time.Second,
					}),
				).To(BeAssignableToTypeOf(&UpdateTtlSet{}))

				time.Sleep(sharedContext.DefaultTtl)

				Expect(
					client.Get(sharedContext.Ctx, &GetRequest{
						CacheName: cacheName,
						Key:       key,
					}),
				).To(BeAssignableToTypeOf(&GetHit{}))
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)

		DescribeTable("Increases Ttl",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				key := String("key")
				value := String("value")

				Expect(
					client.Set(sharedContext.Ctx, &SetRequest{
						CacheName: cacheName,
						Key:       key,
						Value:     value,
						Ttl:       2 * time.Second,
					}),
				).To(BeAssignableToTypeOf(&SetSuccess{}))

				Expect(
					client.IncreaseTtl(sharedContext.Ctx, &IncreaseTtlRequest{
						CacheName: cacheName,
						Key:       key,
						Ttl:       3 * time.Second,
					}),
				).To(BeAssignableToTypeOf(&IncreaseTtlSet{}))

				time.Sleep(2 * time.Second)

				Expect(
					client.Get(sharedContext.Ctx, &GetRequest{
						CacheName: cacheName,
						Key:       key,
					}),
				).To(BeAssignableToTypeOf(&GetHit{}))
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)

		DescribeTable("Decreases Ttl",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)

				key := String("key")
				value := String("value")

				Expect(
					client.Set(sharedContext.Ctx, &SetRequest{
						CacheName: cacheName,
						Key:       key,
						Value:     value,
					}),
				).To(BeAssignableToTypeOf(&SetSuccess{}))

				Expect(
					client.DecreaseTtl(sharedContext.Ctx, &DecreaseTtlRequest{
						CacheName: cacheName,
						Key:       key,
						Ttl:       2 * time.Second,
					}),
				).To(BeAssignableToTypeOf(&DecreaseTtlSet{}))

				time.Sleep(2 * time.Second)

				Expect(
					client.Get(sharedContext.Ctx, &GetRequest{
						CacheName: cacheName,
						Key:       key,
					}),
				).To(BeAssignableToTypeOf(&GetMiss{}))
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)

		It("Returns InvalidArgumentError with negative or zero Ttl value for UpdateTtl", func() {
			key := String("key")
			value := String("value")

			Expect(
				sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Value:     value,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			Expect(
				sharedContext.Client.UpdateTtl(sharedContext.Ctx, &UpdateTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Ttl:       0,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.UpdateTtl(sharedContext.Ctx, &UpdateTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Ttl:       -1,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("Returns InvalidArgumentError with negative or zero Ttl value for IncreaseTtl", func() {
			key := String("key")
			value := String("value")

			Expect(
				sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Value:     value,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			Expect(
				sharedContext.Client.IncreaseTtl(sharedContext.Ctx, &IncreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Ttl:       0,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.IncreaseTtl(sharedContext.Ctx, &IncreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Ttl:       -1,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("Returns InvalidArgumentError with negative or zero Ttl value for DecreaseTtl", func() {
			key := String("key")
			value := String("value")

			Expect(
				sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Value:     value,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			Expect(
				sharedContext.Client.DecreaseTtl(sharedContext.Ctx, &DecreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Ttl:       0,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.DecreaseTtl(sharedContext.Ctx, &DecreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Ttl:       -1,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("Returns UpdateTtlMiss when a key doesn't exist", func() {
			Expect(
				sharedContext.Client.UpdateTtl(sharedContext.Ctx, &UpdateTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       String("test-update-ttl-miss"),
					Ttl:       sharedContext.DefaultTtl,
				}),
			).To(BeAssignableToTypeOf(&UpdateTtlMiss{}))
		})

		It("Returns IncreaseTtlMiss when a key doesn't exist", func() {
			Expect(
				sharedContext.Client.IncreaseTtl(sharedContext.Ctx, &IncreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       String("test-increase-ttl-miss"),
					Ttl:       sharedContext.DefaultTtl,
				}),
			).To(BeAssignableToTypeOf(&IncreaseTtlMiss{}))
		})

		It("Returns DecreaseTtlMiss when a key doesn't exist", func() {
			Expect(
				sharedContext.Client.DecreaseTtl(sharedContext.Ctx, &DecreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       String("test-decrease-ttl-miss"),
					Ttl:       sharedContext.DefaultTtl,
				}),
			).To(BeAssignableToTypeOf(&DecreaseTtlMiss{}))
		})
	})

	Describe("Increment", func() {
		It("Increments from 0 to expected amount with string field", func() {
			field := String("field")

			resp, err := sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    1,
			})
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			Expect(err).To(BeNil())
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(1)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}

			resp, err = sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    41,
			})
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			Expect(err).To(BeNil())
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(42)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}

			resp, err = sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    -1042,
			})
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			Expect(err).To(BeNil())
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(-1000)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}
		})

		It("Increments from 0 to expected amount with bytes field", func() {
			field := Bytes([]byte{1, 2, 3})

			resp, err := sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    1,
			})
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			Expect(err).To(BeNil())
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(1)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}

			resp, err = sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    41,
			})
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			Expect(err).To(BeNil())
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(42)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}

			resp, err = sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    -1042,
			})
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			Expect(err).To(BeNil())
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(-1000)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}
		})

		It("Increments with setting and resetting field", func() {
			field := String("field")
			value := String("10")
			Expect(sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{CacheName: sharedContext.CacheName, Key: field, Value: value})).To(BeAssignableToTypeOf(&SetSuccess{}))

			resp, err := sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    0,
			})
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			Expect(err).To(BeNil())
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(10)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}

			resp, err = sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    90,
			})
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			Expect(err).To(BeNil())
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(100)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}

			// Reset the field
			value = String("0")
			Expect(sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{CacheName: sharedContext.CacheName, Key: field, Value: value})).To(BeAssignableToTypeOf(&SetSuccess{}))
			resp, err = sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    0,
			})
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			Expect(err).To(BeNil())
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(0)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}
		})
	})

	Describe("ItemGetType", func() {
		BeforeEach(func() {
			_, err := sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
				CacheName: sharedContext.CacheName,
				Key:       String("SCALAR"),
				Value:     String("hi"),
			})
			if err != nil {
				Fail("failed trying to set key")
			}
			_, err = sharedContext.Client.DictionarySetField(sharedContext.Ctx, &DictionarySetFieldRequest{
				CacheName:      sharedContext.CacheName,
				DictionaryName: "DICTIONARY",
				Field:          String("hi"),
				Value:          String("there"),
			})
			if err != nil {
				Fail("failed trying to add dictionary")
			}
			_, err = sharedContext.Client.ListPushFront(sharedContext.Ctx, &ListPushFrontRequest{
				CacheName: sharedContext.CacheName,
				ListName:  "LIST",
				Value:     String("hi"),
			})
			if err != nil {
				Fail("failed trying to add list")
			}
			_, err = sharedContext.Client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
				CacheName: sharedContext.CacheName,
				SetName:   "SET",
				Element:   String("hi"),
			})
			if err != nil {
				Fail("failed trying to add set")
			}
			_, err = sharedContext.Client.SortedSetPutElement(sharedContext.Ctx, &SortedSetPutElementRequest{
				CacheName: sharedContext.CacheName,
				SetName:   "SORTED_SET",
				Value:     String("hi"),
				Score:     1.0,
			})
			if err != nil {
				Fail("failed trying to add sorted set")
			}
		})
		DescribeTable("returns the correct item type",
			func(key Key, expectedType ItemType) {
				typeResponse, err := sharedContext.Client.ItemGetType(sharedContext.Ctx, &ItemGetTypeRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				})
				Expect(err).To(BeNil())
				Expect(typeResponse).To(BeAssignableToTypeOf(&ItemGetTypeHit{}))
				switch result := typeResponse.(type) {
				case *ItemGetTypeHit:
					Expect(result.Type()).To(Equal(expectedType))
				default:
					Fail(fmt.Sprintf("expected ItemGetTypeHit but got %s", result))
				}
			},
			Entry("scalar", String("SCALAR"), ItemTypeScalar),
			Entry("dictionary", String("DICTIONARY"), ItemTypeDictionary),
			Entry("set", String("SET"), ItemTypeSet),
			Entry("list", String("LIST"), ItemTypeList),
			Entry("sorted set", String("SORTED_SET"), ItemTypeSortedSet),
		)
	})

	Describe("item get ttl", func() {
		It("accurately reports the remaining TTL for a key", func() {
			var ttl = time.Duration(time.Second * 60)
			_, err := sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
				CacheName: sharedContext.CacheName,
				Key:       String("hi"),
				Value:     String("there"),
				Ttl:       ttl,
			})
			Expect(err).To(BeNil())
			resp, err := sharedContext.Client.ItemGetTtl(sharedContext.Ctx, &ItemGetTtlRequest{
				CacheName: sharedContext.CacheName,
				Key:       String("hi"),
			})
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(&ItemGetTtlHit{}))
			switch result := resp.(type) {
			case *ItemGetTtlHit:
				Expect(ttl > result.RemainingTtl()).To(BeTrue())
				Expect(result.RemainingTtl() > (time.Second * 30)).To(BeTrue())
			default:
				Fail(fmt.Sprintf("expected ItemGetTtlHit but got %s", result))
			}
		})
	})
})
