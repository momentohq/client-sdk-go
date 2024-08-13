package momento_test

import (
	"github.com/momentohq/client-sdk-go/storageTypes"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/momentohq/client-sdk-go/momento"
)

var _ = Describe("storage-client scalar", func() {
	DescribeTable("Sets with correct StorageValueType",
		func(key string, value storageTypes.Value) {
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
			val := resp.Value()
			switch val.(type) {
			case storageTypes.String:
				Expect(val).To(BeAssignableToTypeOf(storageTypes.String("")))
			case storageTypes.Int:
				Expect(val).To(BeAssignableToTypeOf(storageTypes.Int(0)))
			case storageTypes.Float:
				Expect(val).To(BeAssignableToTypeOf(storageTypes.Float(0.0)))
			case storageTypes.Bytes:
				Expect(val).To(BeAssignableToTypeOf(storageTypes.Bytes([]byte{})))
			}
		},
		Entry("String", uuid.NewString(), storageTypes.String("string-value")),
		Entry("Int", uuid.NewString(), storageTypes.Int(42)),
		Entry("Float", uuid.NewString(), storageTypes.Float(3.14)),
		Entry("Bytes", uuid.NewString(), storageTypes.Bytes([]byte{0x01, 0x02, 0x03})),
	)

	It("handles a get without a switch for a known type", func() {
		key := uuid.NewString()
		value := storageTypes.String("string-value")
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
		Expect(resp.Value().(storageTypes.String)).To(Equal(value))
	})

	It("does the right thing on a miss", func() {
		key := uuid.NewString()
		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(resp.Value()).To(BeNil())
	})

	It("does the right thing on an incorrect cast", func() {
		key := uuid.NewString()
		value := storageTypes.String("string-value")
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
		val := resp.Value()
		_, ok := val.(storageTypes.Int)
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
		Expect(resp.Value()).To(BeNil())

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
			Value:     storageTypes.String("value"),
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
		Expect(resp.Value()).To(BeNil())

		_, err = sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			Key:   key,
			Value: storageTypes.String("value"),
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
			Value:     storageTypes.String("value"),
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       "",
		})
		Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
		Expect(resp.Value()).To(BeNil())

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
			Value:     storageTypes.String(value),
		})
		Expect(err).To(Succeed())

		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())

		respValue := resp.Value()
		storeValue, ok := respValue.(storageTypes.String)
		Expect(ok).To(BeTrue())
		Expect(string(storeValue)).To(Equal(value))
		_, ok = respValue.(storageTypes.Int)
		Expect(ok).To(BeFalse())
		_, ok = respValue.(storageTypes.Float)
		Expect(ok).To(BeFalse())
		_, ok = respValue.(storageTypes.Bytes)
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
			Value:     storageTypes.Int(value),
		})
		Expect(err).To(Succeed())

		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())

		respValue := resp.Value()
		storeValue, ok := respValue.(storageTypes.Int)
		Expect(ok).To(BeTrue())
		Expect(int(storeValue)).To(Equal(value))
		_, ok = respValue.(storageTypes.String)
		Expect(ok).To(BeFalse())
		_, ok = respValue.(storageTypes.Float)
		Expect(ok).To(BeFalse())
		_, ok = respValue.(storageTypes.Bytes)
		Expect(ok).To(BeFalse())
	})

	It("Handles floats/doubles", func() {
		key := uuid.NewString()
		value := 3.14
		_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     storageTypes.Float(value),
		})
		Expect(err).To(Succeed())

		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())

		respValue := resp.Value()
		storeValue, ok := respValue.(storageTypes.Float)
		Expect(ok).To(BeTrue())
		Expect(float64(storeValue)).To(Equal(value))
		_, ok = respValue.(storageTypes.String)
		Expect(ok).To(BeFalse())
		_, ok = respValue.(storageTypes.Int)
		Expect(ok).To(BeFalse())
		_, ok = respValue.(storageTypes.Bytes)
		Expect(ok).To(BeFalse())
	})

	It("Handles bytes", func() {
		key := uuid.NewString()
		value := []byte{0x01, 0x02, 0x03}
		_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     storageTypes.Bytes(value),
		})
		Expect(err).To(Succeed())

		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())

		respValue := resp.Value()
		storeValue, ok := respValue.(storageTypes.Bytes)
		Expect(ok).To(BeTrue())
		Expect([]byte(storeValue)).To(Equal(value))
		_, ok = respValue.(storageTypes.String)
		Expect(ok).To(BeFalse())
		_, ok = respValue.(storageTypes.Int)
		Expect(ok).To(BeFalse())
		_, ok = respValue.(storageTypes.Float)
		Expect(ok).To(BeFalse())
	})

	It("Handles a miss", func() {
		key := uuid.NewString()
		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())
		Expect(resp.Value()).To(BeNil())
		Expect(resp.Value()).To(BeNil())
	})

	It("reads directly from the response", func() {
		key := uuid.NewString()
		val := uuid.NewString()
		_, err := sharedContext.StorageClient.Put(sharedContext.Ctx, &StoragePutRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
			Value:     storageTypes.String(val),
		})
		Expect(err).To(Succeed())
		resp, err := sharedContext.StorageClient.Get(sharedContext.Ctx, &StorageGetRequest{
			StoreName: sharedContext.StoreName,
			Key:       key,
		})
		Expect(err).To(Succeed())

		// unwrap value with no switch at all
		_ = resp.Value().(storageTypes.String)
	})

})
