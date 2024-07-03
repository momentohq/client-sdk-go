package storageTypes

type Value interface {
	isValue()
}

// String type to store string values.
type String string

// Bytes type to store byte values.
type Bytes []byte

// Int type to store ints.
type Int int64

// Float type to store floats.
type Float float64

func (String) isValue() {}

func (Bytes) isValue() {}

func (Int) isValue() {}

func (Float) isValue() {}
