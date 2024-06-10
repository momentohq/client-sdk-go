package momento

type StorageValue interface {
	isStorageValue()
}

// StorageValueString type to store string values.
type StorageValueString string

// StorageValueBytes type to store byte values.
type StorageValueBytes []byte

// StorageValueInteger type to store ints.
type StorageValueInteger int64

// StorageValueDouble backed by float64 as Go doesn't have a double type.
type StorageValueDouble float64

func (StorageValueString) isStorageValue() {}

func (StorageValueBytes) isStorageValue() {}

func (StorageValueInteger) isStorageValue() {}

func (StorageValueDouble) isStorageValue() {}
