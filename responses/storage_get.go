package responses

import (
	"github.com/momentohq/client-sdk-go/storageTypes"
)

// StorageGetResponse is the base response type for a store get request.
type StorageGetResponse struct {
	value storageTypes.Value
}

func (r StorageGetResponse) Value() storageTypes.Value {
	return r.value
}

// NewStoreGetFound_String returns a new StorageGetResponse containing the supplied string value.
func NewStoreGetResponse_String(value string) *StorageGetResponse {
	return &StorageGetResponse{
		value: storageTypes.String(value),
	}
}

// NewStoreGetResponse_Bytes returns a new StorageGetResponse containing the supplied byte slice value.
func NewStoreGetResponse_Bytes(value []byte) *StorageGetResponse {
	return &StorageGetResponse{
		value: storageTypes.Bytes(value),
	}
}

// NewStoreGetResponse_Float returns a new StorageGetResponse containing the supplied float64 value.
func NewStoreGetResponse_Float(value float64) *StorageGetResponse {
	return &StorageGetResponse{
		value: storageTypes.Float(value),
	}
}

// NewStoreGetResponse_Integer returns a new StorageGetResponse containing the supplied int value.
func NewStoreGetResponse_Integer(value int) *StorageGetResponse {
	return &StorageGetResponse{
		value: storageTypes.Int(value),
	}
}

func NewStoreGetResponse_Nil() *StorageGetResponse {
	return &StorageGetResponse{
		value: nil,
	}
}
