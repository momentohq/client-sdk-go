package momento

// Value Interface to help users deal with passing us values as strings or bytes.
// Value: momento.Bytes([]bytes("abc"))
// Value: momento.String("abc")
type Value interface {
	asBytes() []byte
	asString() string
}

// Key Type alias to future proof passing in keys.
type Key = Value

// Bytes plain old []byte
type Bytes []byte

func (v Bytes) asBytes() []byte { return v }

func (v Bytes) asString() string {
	return string(v)
}

// String string type that will be converted to []byte
type String string

func (v String) asBytes() []byte {
	return []byte(v)
}

func (v String) asString() string {
	return string(v)
}

// Type to hold field/value elements in dictionaries for use with DictionarySetFields.
// Field and Value are both Value type which allows both Strings and Bytes.
type DictionaryElement struct {
	Field, Value Value
}
