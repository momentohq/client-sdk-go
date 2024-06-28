// TODO: move me to a more appropriate namespace
package utils

type StorageValue interface {
	isStorageValue()
}

// StorageValueString type to store string values.
type StorageValueString string

// StorageValueBytes type to store byte values.
type StorageValueBytes []byte

// StorageValueInt type to store ints.
type StorageValueInt int64

// StorageValueFloat type to store floats.
type StorageValueFloat float64

func (StorageValueString) isStorageValue() {}

func (StorageValueBytes) isStorageValue() {}

func (StorageValueInt) isStorageValue() {}

func (StorageValueFloat) isStorageValue() {}
