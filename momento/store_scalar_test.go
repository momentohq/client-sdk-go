package momento_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	. "github.com/momentohq/client-sdk-go/responses"
)

var _ = Describe("Store scalar methods", func() {
	var sharedContext SharedContext
	BeforeEach(func() {
		sharedContext = NewSharedContext()
		sharedContext.CreateDefaultStores()

		DeferCleanup(func() { sharedContext.Close() })
	})

	DescribeTable("Puts with correct StoreValueType",
		func(key string, value StoreValue, expected StoreValueType) {
			_, err := sharedContext.StoreClient.Put(sharedContext.Ctx, &StorePutRequest{
				StoreName: sharedContext.StoreName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(Succeed())

			resp, err := sharedContext.StoreClient.Get(sharedContext.Ctx, &StoreGetRequest{Key: key})
			Expect(err).To(Succeed())
			Expect(resp).To(BeAssignableToTypeOf(&StoreGetSuccess{}))
			Expect(resp.ValueType()).To(Equal(expected))
		},
		Entry(uuid.NewString(), String("string-value"), STRING),
		Entry(uuid.NewString(), Integer(42), INTEGER),
		Entry(uuid.NewString(), Double(3.14), DOUBLE),
		Entry(uuid.NewString(), Bytes([]byte{0x01, 0x02, 0x03}), BYTES),
	)

	It("errors on missing store name", func() {
		key := uuid.NewString()
		_, err := sharedContext.StoreClient.Get(sharedContext.Ctx, &StoreGetRequest{
			Key: key,
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
		_, err = sharedContext.StoreClient.Put(sharedContext.Ctx, &StorePutRequest{
			Key:   key,
			Value: String("value"),
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

		_, err = sharedContext.StoreClient.Delete(sharedContext.Ctx, &StoreDeleteRequest{
			Key: key,
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	It("errors on invalid key", func() {
		_, err := sharedContext.StoreClient.Put(sharedContext.Ctx, &StorePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       "",
			Value:     String("value"),
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
		_, err = sharedContext.StoreClient.Get(sharedContext.Ctx, &StoreGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       "",
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
		_, err = sharedContext.StoreClient.Delete(sharedContext.Ctx, &StoreDeleteRequest{
			StoreName: sharedContext.StoreName,
			Key:       "",
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	It("Handles strings", func() {
		key := uuid.NewString()
		value := "string-value"
		_, err := sharedContext.StoreClient.Put(sharedContext.Ctx, &StorePutRequest{
			Key:   key,
			Value: String(value),
		})
		Expect(err).To(Succeed())

		resp, err := sharedContext.StoreClient.Get(sharedContext.Ctx, &StoreGetRequest{Key: key})
		Expect(err).To(Succeed())
		Expect(resp).To(BeAssignableToTypeOf(&StoreGetSuccess{}))
		Expect(resp.ValueType()).To(Equal(STRING))
		storeValue, ok := resp.TryGetValueString()
		Expect(ok).To(BeTrue())
		Expect(storeValue).To(Equal(value))
		_, ok = resp.TryGetValueInteger()
		Expect(ok).To(BeFalse())
		_, ok = resp.TryGetValueDouble()
		Expect(ok).To(BeFalse())
		_, ok = resp.TryGetValueBytes()
		Expect(ok).To(BeFalse())
	})

	It("Handles nil values", func() {
		key := uuid.NewString()
		_, err := sharedContext.StoreClient.Put(sharedContext.Ctx, &StorePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		if err != nil {
			panic(err)
		}
	})

	It("Handles integers", func() {
		key := uuid.NewString()
		value := 42
		_, err := sharedContext.StoreClient.Put(sharedContext.Ctx, &StorePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     Integer(value),
		})
		Expect(err).To(Succeed())

		resp, err := sharedContext.StoreClient.Get(sharedContext.Ctx, &StoreGetRequest{Key: key})
		Expect(err).To(Succeed())
		Expect(resp).To(BeAssignableToTypeOf(&StoreGetSuccess{}))
		Expect(resp.ValueType()).To(Equal(INTEGER))
		storeValue, ok := resp.TryGetValueInteger()
		Expect(ok).To(BeTrue())
		Expect(storeValue).To(Equal(value))
		_, ok = resp.TryGetValueString()
		Expect(ok).To(BeFalse())
		_, ok = resp.TryGetValueDouble()
		Expect(ok).To(BeFalse())
		_, ok = resp.TryGetValueBytes()
		Expect(ok).To(BeFalse())
	})

	It("Handles doubles", func() {
		key := uuid.NewString()
		value := 3.14
		_, err := sharedContext.StoreClient.Put(sharedContext.Ctx, &StorePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     Double(value),
		})
		Expect(err).To(Succeed())

		resp, err := sharedContext.StoreClient.Get(sharedContext.Ctx, &StoreGetRequest{Key: key})
		Expect(err).To(Succeed())
		Expect(resp).To(BeAssignableToTypeOf(&StoreGetSuccess{}))
		Expect(resp.ValueType()).To(Equal(DOUBLE))
		storeValue, ok := resp.TryGetValueDouble()
		Expect(ok).To(BeTrue())
		Expect(storeValue).To(Equal(value))
		_, ok = resp.TryGetValueString()
		Expect(ok).To(BeFalse())
		_, ok = resp.TryGetValueInteger()
		Expect(ok).To(BeFalse())
		_, ok = resp.TryGetValueBytes()
		Expect(ok).To(BeFalse())
	})

	It("Handles bytes", func() {
		key := uuid.NewString()
		value := []byte{0x01, 0x02, 0x03}
		_, err := sharedContext.StoreClient.Put(sharedContext.Ctx, &StorePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     Bytes(value),
		})
		Expect(err).To(Succeed())

		resp, err := sharedContext.StoreClient.Get(sharedContext.Ctx, &StoreGetRequest{Key: key})
		Expect(err).To(Succeed())
		Expect(resp).To(BeAssignableToTypeOf(&StoreGetSuccess{}))
		Expect(resp.ValueType()).To(Equal(BYTES))
		storeValue, ok := resp.TryGetValueBytes()
		Expect(ok).To(BeTrue())
		Expect(storeValue).To(Equal(value))
		_, ok = resp.TryGetValueString()
		Expect(ok).To(BeFalse())
		_, ok = resp.TryGetValueInteger()
		Expect(ok).To(BeFalse())
		_, ok = resp.TryGetValueDouble()
		Expect(ok).To(BeFalse())
	})
})
