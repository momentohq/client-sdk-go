package momento

type StoreValue interface {
	isStoreValue()
}

type Integer int

type Double float64

// TODO: reusing String and Bytes types from value.go, but not sure I should be doing that.
func (String) isStoreValue() {}

func (Bytes) isStoreValue() {}

func (Integer) isStoreValue() {}

func (Double) isStoreValue() {}
