package momento_test

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
)

func HaveMomentoErrorCode(code string) types.GomegaMatcher {
	return WithTransform(
		func(err error) (string, error) {
			switch mErr := err.(type) {
			case MomentoError:
				return mErr.Code(), nil
			default:
				return "", fmt.Errorf("Expected MomentoError, but got %T", err)
			}
		}, Equal(code),
	)
}

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

			time.Sleep(sharedContext.DefaultTTL / 2)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetHit{}))

			time.Sleep(sharedContext.DefaultTTL)

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
					TTL:       sharedContext.DefaultTTL * 2,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			time.Sleep(sharedContext.DefaultTTL / 2)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetHit{}))

			time.Sleep(sharedContext.DefaultTTL)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetHit{}))

			time.Sleep(sharedContext.DefaultTTL)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		})
	})
})
