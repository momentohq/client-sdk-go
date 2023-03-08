package momento_test

import (
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
		sharedContext.CreateDefaultCache()

		DeferCleanup(func() { sharedContext.Close() })
	})

	DescribeTable(`Gets, Sets, and Deletes`,
		func(key Key, value Value, expectedString string, expectedBytes []byte) {
			Expect(
				sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Value:     value,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			getResp, err := sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: sharedContext.CacheName,
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
				sharedContext.Client.Delete(sharedContext.Ctx, &DeleteRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&DeleteSuccess{}))

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		},
		Entry("when the key and value are strings", String("key"), String("value"), "value", []byte("value")),
		Entry("when the key and value are bytes", Bytes([]byte{1, 2, 3}), Bytes([]byte("string")), "string", []byte("string")),
		Entry("when the value is empty", String("key"), String(""), "", []byte("")),
		Entry("when the value is blank", String("key"), String("  "), "  ", []byte("  ")),
	)

	It(`errors when the cache is missing`, func() {
		cacheName := uuid.NewString()
		key := String("key")
		value := String("value")

		getResp, err := sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
			CacheName: cacheName,
			Key:       key,
		})
		Expect(getResp).To(BeNil())
		Expect(err).To(HaveMomentoErrorCode(NotFoundError))

		setResp, err := sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
			CacheName: cacheName,
			Key:       key,
			Value:     value,
		})
		Expect(setResp).To(BeNil())
		Expect(err).To(HaveMomentoErrorCode(NotFoundError))

		deleteResp, err := sharedContext.Client.Delete(sharedContext.Ctx, &DeleteRequest{
			CacheName: cacheName,
			Key:       key,
		})
		Expect(deleteResp).To(BeNil())
		Expect(err).To(HaveMomentoErrorCode(NotFoundError))
	})

	It(`errors when the key is nil`, func() {
		Expect(
			sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: sharedContext.CacheName,
				Key:       nil,
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		Expect(
			sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
				CacheName: sharedContext.CacheName,
				Key:       nil,
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		Expect(
			sharedContext.Client.Delete(sharedContext.Ctx, &DeleteRequest{
				CacheName: sharedContext.CacheName,
				Key:       nil,
			}),
		).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	It("returns a miss when the key doesn't exist", func() {
		Expect(
			sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: sharedContext.CacheName,
				Key:       String(uuid.NewString()),
			}),
		).To(BeAssignableToTypeOf(&GetMiss{}))
	})

	DescribeTable(`invalid cache names and keys`,
		func(cacheName string, key Key, value Key) {
			getResp, err := sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(getResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

			setResp, err := sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(setResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

			deleteResp, err := sharedContext.Client.Delete(sharedContext.Ctx, &DeleteRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(deleteResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
		},
		Entry(`With an empty cache name`, "", String("key"), String("value")),
		Entry(`With an bank cache name`, "   ", String("key"), String("value")),
		Entry(`With an empty key`, uuid.NewString(), String(""), String("value")),
	)

	Describe(`Set`, func() {
		It(`Uses the default TTL`, func() {
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

		It(`Overrides the default TTL`, func() {
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

	Describe(`UpdateTtl`, func() {
		It(`Overwrites Ttl`, func() {
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
					Ttl:       6 * time.Second,
				}),
			).To(BeAssignableToTypeOf(&UpdateTtlSet{}))

			time.Sleep(sharedContext.DefaultTtl)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetHit{}))
		})

		It(`Increases Ttl`, func() {
			key := String("key")
			value := String("value")

			Expect(
				sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Value:     value,
					Ttl:       2 * time.Second,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			Expect(
				sharedContext.Client.IncreaseTtl(sharedContext.Ctx, &IncreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Ttl:       3 * time.Second,
				}),
			).To(BeAssignableToTypeOf(&IncreaseTtlSet{}))

			time.Sleep(2 * time.Second)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetHit{}))
		})

		It(`Decreases Ttl`, func() {
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
					Ttl:       2 * time.Second,
				}),
			).To(BeAssignableToTypeOf(&DecreaseTtlSet{}))

			time.Sleep(sharedContext.DefaultTtl)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		})

		It(`Returns InvalidArgumentError with negative or zero Ttl value for UpdateTtl`, func() {
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

		It(`Returns InvalidArgumentError with negative or zero Ttl value for IncreaseTtl`, func() {
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

		It(`Returns InvalidArgumentError with negative or zero Ttl value for DecreaseTtl`, func() {
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

		It(`Returns UpdateTtlMiss when a key doesn't exist`, func() {
			Expect(
				sharedContext.Client.UpdateTtl(sharedContext.Ctx, &UpdateTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       String("test-update-ttl-miss"),
					Ttl:       sharedContext.DefaultTtl,
				}),
			).To(BeAssignableToTypeOf(&UpdateTtlMiss{}))
		})

		It(`Returns IncreaseTtlMiss when a key doesn't exist`, func() {
			Expect(
				sharedContext.Client.IncreaseTtl(sharedContext.Ctx, &IncreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       String("test-increase-ttl-miss"),
					Ttl:       sharedContext.DefaultTtl,
				}),
			).To(BeAssignableToTypeOf(&IncreaseTtlMiss{}))
		})

		It(`Returns DecreaseTtlMiss when a key doesn't exist`, func() {
			Expect(
				sharedContext.Client.DecreaseTtl(sharedContext.Ctx, &DecreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       String("test-decrease-ttl-miss"),
					Ttl:       sharedContext.DefaultTtl,
				}),
			).To(BeAssignableToTypeOf(&DecreaseTtlMiss{}))
		})
	})
})
