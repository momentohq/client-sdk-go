package momento_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	. "github.com/momentohq/client-sdk-go/responses"
)

var _ = Describe("storage-client scalar", func() {
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
			Expect(resp).To(BeAssignableToTypeOf(&StorageGetFound{}))
			val := resp.(*StorageGetFound).Value()
			switch val.(type) {
			case utils.StorageValueString:
				Expect(val).To(BeAssignableToTypeOf(utils.StorageValueString("")))
			case utils.StorageValueInt:
				Expect(val).To(BeAssignableToTypeOf(utils.StorageValueInt(0)))
			case utils.StorageValueFloat:
				Expect(val).To(BeAssignableToTypeOf(utils.StorageValueFloat(0.0)))
			case utils.StorageValueBytes:
				Expect(val).To(BeAssignableToTypeOf(utils.StorageValueBytes([]byte{})))
			}
		},
		Entry("StorageValueString", uuid.NewString(), StorageValueString("string-value"), STRING),
		Entry("StorageValueInteger", uuid.NewString(), StorageValueInteger(42), INTEGER),
		Entry("StorageValueDouble", uuid.NewString(), StorageValueDouble(3.14), DOUBLE),
		Entry("StorageValueBytes", uuid.NewString(), StorageValueBytes([]byte{0x01, 0x02, 0x03}), BYTES),
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

		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(resp).To(BeAssignableToTypeOf(&StorageGetFound{}))
		Expect(resp.(*StorageGetFound).Value().(utils.StorageValueString)).To(Equal(value))
	})

	It("does the right thing on a miss", func() {
		key := uuid.NewString()
		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(resp).To(BeAssignableToTypeOf(&StorageGetNotFound{}))
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

		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(resp).To(BeAssignableToTypeOf(&StorageGetFound{}))
		val := resp.(*StorageGetFound).Value()
		_, ok := val.(utils.StorageValueInt)
		Expect(ok).To(BeFalse())
	})

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
		Expect(string(storeValue)).To(Equal(value))
		_, ok = hitValue.(utils.StorageValueInt)
		Expect(ok).To(BeFalse())
		_, ok = resp.ValueDouble()
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

		Expect(resp).To(BeAssignableToTypeOf(&StorageGetFound{}))

		hitResp := resp.(*StorageGetFound)
		hitValue := hitResp.Value()
		storeValue, ok := hitValue.(utils.StorageValueInt)
		Expect(ok).To(BeTrue())
		Expect(storeValue).To(Equal(value))
		_, ok = resp.ValueString()
		Expect(ok).To(BeFalse())
		_, ok = resp.ValueDouble()
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
			Value:     StorageValueDouble(value),
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

		storeValue, ok := resp.ValueDouble()
		Expect(ok).To(BeTrue())
		Expect(storeValue).To(Equal(value))
		_, ok = resp.ValueString()
		Expect(ok).To(BeFalse())
		_, ok = hitValue.(utils.StorageValueInt)
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
		_, ok = hitValue.(utils.StorageValueInt)
		Expect(ok).To(BeFalse())
		_, ok = resp.ValueDouble()
		Expect(ok).To(BeFalse())
	})

	It("Handles a miss", func() {
		key := uuid.NewString()
		resp, found, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(found).To(BeFalse())
		Expect(resp).To(BeNil())
	})

	It("reads directly from the response", func() {
		key := uuid.NewString()
		val := uuid.NewString()
		_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     utils.StorageValueString(val),
		})
		Expect(err).To(Succeed())
		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())

		// unwrap value with no switch at all
		_ = resp.(*StorageGetFound).Value().(utils.StorageValueString)
		// unwrap value with a switch for response type
		switch r := resp.(type) {
		case *StorageGetFound:
			fmt.Printf("resp as string: %s\n", r.Value().(utils.StorageValueString))
		case *StorageGetNotFound:
			fmt.Println("not found")
		}
	})

	//It("deleteallstores", func() {
	//	resp, err := sharedContext.StorageClient.ListStores(sharedContext.Ctx, &ListStoresRequest{})
	//	Expect(err).To(Succeed())
	//	for _, store := range resp.(*ListStoresSuccess).Stores() {
	//		name := store.Name()
	//		_, err := sharedContext.StorageClient.DeleteStore(sharedContext.Ctx, &DeleteStoreRequest{
	//			StoreName: name,
	//		})
	//		Expect(err).To(Succeed())
	//	}
	//})

})
