package responses

import "github.com/momentohq/client-sdk-go/utils"

// StorageGetResponse is the base response type for a store get request.
type StorageGetResponse interface {
	isStoreGetResponse()
}

type StorageGetNotFound struct{}

func (StorageGetNotFound) isStoreGetResponse() {}

func (StorageGetNotFound) Value() utils.StorageValue {
	return nil
}

func NewStoreGetNotFound() *StorageGetNotFound {
	return &StorageGetNotFound{}
}

// StorageGetFound indicates a successful store get request.
type StorageGetFound struct {
	value utils.StorageValue
}

func (StorageGetFound) isStoreGetResponse() {}

// ValueType returns the `StorageValueType` indicating the type of the value in the store.
func (resp StorageGetFound) Value() utils.StorageValue {
	return resp.value
}

// NewStoreGetFound_String returns a new StorageGetFound containing the supplied string value.
func NewStoreGetFound_String(value string) *StorageGetFound {
	return &StorageGetFound{
		value: utils.StorageValueString(value),
	}
}

// NewStoreGetFound_Bytes returns a new StorageGetFound containing the supplied byte slice value.
func NewStoreGetFound_Bytes(value []byte) *StorageGetFound {
	return &StorageGetFound{
		value: utils.StorageValueBytes(value),
	}
}

// NewStoreGetFound_Float returns a new StorageGetFound containing the supplied float64 value.
func NewStoreGetFound_Float(value float64) *StorageGetFound {
	return &StorageGetFound{
		value: utils.StorageValueFloat(value),
	}
}

// NewStoreGetFound_Integer returns a new StorageGetFound containing the supplied int value.
func NewStoreGetFound_Integer(value int) *StorageGetFound {
	return &StorageGetFound{
		value: utils.StorageValueInt(value),
	}
}
