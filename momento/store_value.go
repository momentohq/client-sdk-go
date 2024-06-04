package momento

type StoreValue interface {
	isStoreValue()
}

// Integer plain old int.
type Integer int

// Double backed by float64 as Go doesn't have a double type.
type Double float64

func (String) isStoreValue() {}

func (Bytes) isStoreValue() {}

func (Integer) isStoreValue() {}

func (Double) isStoreValue() {}
