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

	DescribeTable("Sets with correct StorageValueType",
		func(key string, value StorageValue, expected StorageValueType) {
			_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
				StoreName: sharedContext.StoreName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(Succeed())

			resp, found, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
				StoreName: sharedContext.StoreName,
				Key:       key,
			})
			Expect(err).To(Succeed())
			Expect(found).To(BeTrue())
			Expect(resp).To(BeAssignableToTypeOf(&StorageGetSuccess{}))
			Expect(resp.ValueType()).To(Equal(expected))
		},
		Entry("StorageValueString", uuid.NewString(), StorageValueString("string-value"), STRING),
		Entry("StorageValueInteger", uuid.NewString(), StorageValueInteger(42), INTEGER),
		Entry("StorageValueFloat64", uuid.NewString(), StorageValueFloat64(3.14), DOUBLE),
		Entry("StorageValueBytes", uuid.NewString(), StorageValueBytes([]byte{0x01, 0x02, 0x03}), BYTES),
	)

	It("errors on missing store or key", func() {
		key := uuid.NewString()
		store := uuid.NewString()

		// Missing key simply returns found=false
		resp, found, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(found).To(BeFalse())
		Expect(err).To(BeNil())
		Expect(resp).To(BeNil())

		_, _, err = sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: store,
			Key:       key,
		})
		Expect(err).To(HaveMomentoErrorCode(StoreNotFoundError))
		Expect(err.Error()).To(ContainSubstring("Store not found"))

		_, err = sharedContext.StorageClient.Delete(sharedContext.Ctx, &StorageDeleteRequest{
			StoreName: store,
			Key:       key,
		})
		Expect(err).To(HaveMomentoErrorCode(StoreNotFoundError))
		Expect(err.Error()).To(ContainSubstring("Store not found"))

		_, err = sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: store,
			Key:       key,
			Value:     StorageValueString("value"),
		})
		Expect(err).To(HaveMomentoErrorCode(StoreNotFoundError))
		Expect(err.Error()).To(ContainSubstring("Store not found"))
	})

	It("errors on missing store name", func() {
		key := uuid.NewString()
		resp, found, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			Key: key,
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
		Expect(found).To(BeFalse())
		Expect(resp).To(BeNil())

		_, err = sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			Key:   key,
			Value: StorageValueString("value"),
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

		_, err = sharedContext.StorageClient.Delete(sharedContext.Ctx, &StorageDeleteRequest{
			Key: key,
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	It("errors on invalid key", func() {
		_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       "",
			Value:     StorageValueString("value"),
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

		resp, found, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       "",
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
		Expect(found).To(BeFalse())
		Expect(resp).To(BeNil())

		_, err = sharedContext.StorageClient.Delete(sharedContext.Ctx, &StorageDeleteRequest{
			StoreName: sharedContext.StoreName,
			Key:       "",
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	It("Handles strings", func() {
		key := uuid.NewString()
		value := "string-value"
		_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     StorageValueString(value),
		})
		Expect(err).To(Succeed())

		resp, found, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(found).To(BeTrue()) // item was found
		Expect(resp).To(BeAssignableToTypeOf(&StorageGetSuccess{}))
		Expect(resp.ValueType()).To(Equal(STRING))

		storeValue, ok := resp.ValueString()
		Expect(ok).To(BeTrue())
		Expect(storeValue).To(Equal(value))
		_, ok = resp.ValueInteger()
		Expect(ok).To(BeFalse())
		_, ok = resp.ValueFloat64()
		Expect(ok).To(BeFalse())
		_, ok = resp.ValueBytes()
		Expect(ok).To(BeFalse())
	})

	It("Handles nil values", func() {
		key := uuid.NewString()
		_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
	})

	It("Handles integers", func() {
		key := uuid.NewString()
		value := 42
		_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     StorageValueInteger(value),
		})
		Expect(err).To(Succeed())

		resp, found, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(found).To(BeTrue()) // item was found
		Expect(resp).To(BeAssignableToTypeOf(&StorageGetSuccess{}))
		Expect(resp.ValueType()).To(Equal(INTEGER))

		storeValue, ok := resp.ValueInteger()
		Expect(ok).To(BeTrue())
		Expect(storeValue).To(Equal(value))
		_, ok = resp.ValueString()
		Expect(ok).To(BeFalse())
		_, ok = resp.ValueFloat64()
		Expect(ok).To(BeFalse())
		_, ok = resp.ValueBytes()
		Expect(ok).To(BeFalse())
	})

	It("Handles doubles", func() {
		key := uuid.NewString()
		value := 3.14
		_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     StorageValueFloat64(value),
		})
		Expect(err).To(Succeed())

		resp, found, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(found).To(BeTrue()) // item was found
		Expect(resp).To(BeAssignableToTypeOf(&StorageGetSuccess{}))
		Expect(resp.ValueType()).To(Equal(DOUBLE))

		storeValue, ok := resp.ValueFloat64()
		Expect(ok).To(BeTrue())
		Expect(storeValue).To(Equal(value))
		_, ok = resp.ValueString()
		Expect(ok).To(BeFalse())
		_, ok = resp.ValueInteger()
		Expect(ok).To(BeFalse())
		_, ok = resp.ValueBytes()
		Expect(ok).To(BeFalse())
	})

	It("Handles bytes", func() {
		key := uuid.NewString()
		value := []byte{0x01, 0x02, 0x03}
		_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     StorageValueBytes(value),
		})
		Expect(err).To(Succeed())

		resp, found, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(found).To(BeTrue()) // item was found
		Expect(resp).To(BeAssignableToTypeOf(&StorageGetSuccess{}))
		Expect(resp.ValueType()).To(Equal(BYTES))

		storeValue, ok := resp.ValueBytes()
		Expect(ok).To(BeTrue())
		Expect(storeValue).To(Equal(value))
		_, ok = resp.ValueString()
		Expect(ok).To(BeFalse())
		_, ok = resp.ValueInteger()
		Expect(ok).To(BeFalse())
		_, ok = resp.ValueFloat64()
		Expect(ok).To(BeFalse())
	})
})
