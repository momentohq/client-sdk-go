package momento

type StoreValue interface {
	isStoreValue()
}

type Integer int

type Double float64

func (String) isStoreValue() {}

func (Bytes) isStoreValue() {}

func (Integer) isStoreValue() {}

func (Double) isStoreValue() {}
