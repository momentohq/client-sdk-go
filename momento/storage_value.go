package momento

type StorageValue interface {
	isStorageValue()
}

// Integer plain old int.
type Integer int64

// Double backed by float64 as Go doesn't have a double type.
type Double float64

func (String) isStorageValue() {}

func (Bytes) isStorageValue() {}

func (Integer) isStorageValue() {}

func (Double) isStorageValue() {}
