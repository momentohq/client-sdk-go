package momento_test

import (
	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/utils"
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

	DescribeTable("Sets, gets, and deletes with correct StorageValueType",
		func(key string, value utils.StorageValue) {
			_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
				StoreName: sharedContext.StoreName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(Succeed())

			val, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
				StoreName: sharedContext.StoreName,
				Key:       key,
			})
			Expect(err).To(Succeed())
			switch val.(type) {
			case utils.StorageValueString:
				Expect(val).To(BeAssignableToTypeOf(utils.StorageValueString("")))
			case utils.StorageValueInt:
				Expect(val).To(BeAssignableToTypeOf(utils.StorageValueInt(0)))
			case utils.StorageValueFloat:
				Expect(val).To(BeAssignableToTypeOf(utils.StorageValueFloat(0.0)))
				Expect(val).To(Equal(value))
			case utils.StorageValueBytes:
				Expect(val).To(BeAssignableToTypeOf(utils.StorageValueBytes([]byte{})))
				Expect(val).To(Equal(value))
			}

			resp, err := sharedContext.StorageClient.Delete(sharedContext.Ctx, &StorageDeleteRequest{
				StoreName: sharedContext.StoreName,
				Key:       key,
			})
			Expect(err).To(Succeed())
			Expect(resp).To(BeAssignableToTypeOf(&StorageDeleteSuccess{}))

			val, err = sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
				StoreName: sharedContext.StoreName,
				Key:       key,
			})
			Expect(err).To(Succeed())
			Expect(val).To(BeNil())
		},
		Entry("StorageValueString", uuid.NewString(), utils.StorageValueString("string-value")),
		Entry("StorageValueInt", uuid.NewString(), utils.StorageValueInt(42)),
		Entry("StorageValueFloat", uuid.NewString(), utils.StorageValueFloat(3.14)),
		Entry("StorageValueBytes", uuid.NewString(), utils.StorageValueBytes([]byte{0x01, 0x02, 0x03})),
	)

	It("handles a get without a switch for a known type", func() {
		key := uuid.NewString()
		value := utils.StorageValueString("string-value")
		_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     value,
		})
		Expect(err).To(Succeed())

		val, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(val).To(BeAssignableToTypeOf(utils.StorageValueString("")))
		Expect(val.(utils.StorageValueString)).To(Equal(value))
	})

	It("does the right thing on an incorrect cast", func() {
		key := uuid.NewString()
		value := utils.StorageValueString("string-value")
		_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     value,
		})
		Expect(err).To(Succeed())

		val, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(val).To(BeAssignableToTypeOf(utils.StorageValueString("")))
		_, ok := val.(utils.StorageValueInt)
		Expect(ok).To(BeFalse())
	})

	It("errors on missing store or key", func() {
		key := uuid.NewString()
		store := uuid.NewString()

		// Missing key simply returns found=false
		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(BeNil())
		Expect(resp).To(BeNil())

		_, err = sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
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
			Value:     utils.StorageValueString("value"),
		})
		Expect(err).To(HaveMomentoErrorCode(StoreNotFoundError))
		Expect(err.Error()).To(ContainSubstring("Store not found"))
	})

	It("errors on missing store name", func() {
		key := uuid.NewString()
		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			Key: key,
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
		Expect(resp).To(BeNil())

		_, err = sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			Key:   key,
			Value: utils.StorageValueString("value"),
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
			Value:     utils.StorageValueString("value"),
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       "",
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
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
			Value:     utils.StorageValueString(value),
		})
		Expect(err).To(Succeed())

		val, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(val).To(BeAssignableToTypeOf(utils.StorageValueString("")))

		storeValue, ok := val.(utils.StorageValueString)
		Expect(ok).To(BeTrue())
		Expect(string(storeValue)).To(Equal(value))
		_, ok = val.(utils.StorageValueInt)
		Expect(ok).To(BeFalse())
		_, ok = val.(utils.StorageValueFloat)
		Expect(ok).To(BeFalse())
		_, ok = val.(utils.StorageValueBytes)
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
			Value:     utils.StorageValueInt(value),
		})
		Expect(err).To(Succeed())

		val, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(val).To(BeAssignableToTypeOf(utils.StorageValueInt(0)))

		storeValue, ok := val.(utils.StorageValueInt)
		Expect(ok).To(BeTrue())
		Expect(int(storeValue)).To(Equal(value))
		_, ok = val.(utils.StorageValueString)
		Expect(ok).To(BeFalse())
		_, ok = val.(utils.StorageValueFloat)
		Expect(ok).To(BeFalse())
		_, ok = val.(utils.StorageValueBytes)
		Expect(ok).To(BeFalse())
	})

	It("Handles floats/doubles", func() {
		key := uuid.NewString()
		value := 3.14
		_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     utils.StorageValueFloat(value),
		})
		Expect(err).To(Succeed())

		val, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(val).To(BeAssignableToTypeOf(utils.StorageValueFloat(0.0)))

		storeValue, ok := val.(utils.StorageValueFloat)
		Expect(ok).To(BeTrue())
		Expect(float64(storeValue)).To(Equal(value))
		_, ok = val.(utils.StorageValueString)
		Expect(ok).To(BeFalse())
		_, ok = val.(utils.StorageValueInt)
		Expect(ok).To(BeFalse())
		_, ok = val.(utils.StorageValueBytes)
		Expect(ok).To(BeFalse())
	})

	It("Handles bytes", func() {
		key := uuid.NewString()
		value := []byte{0x01, 0x02, 0x03}
		_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     utils.StorageValueBytes(value),
		})
		Expect(err).To(Succeed())

		val, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(val).To(BeAssignableToTypeOf(utils.StorageValueBytes{}))

		storeValue, ok := val.(utils.StorageValueBytes)
		Expect(ok).To(BeTrue())
		Expect([]byte(storeValue)).To(Equal(value))
		_, ok = val.(utils.StorageValueString)
		Expect(ok).To(BeFalse())
		_, ok = val.(utils.StorageValueInt)
		Expect(ok).To(BeFalse())
		_, ok = val.(utils.StorageValueFloat)
		Expect(ok).To(BeFalse())
	})

	It("Handles a miss", func() {
		key := uuid.NewString()
		val, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(val).To(BeNil())
		_, ok := val.(utils.StorageValueFloat)
		Expect(ok).To(BeFalse())
	})

})
