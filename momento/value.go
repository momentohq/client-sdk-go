package momento

// Interface to help users deal with passing us values as strings or bytes.
// Value: momento.Bytes([]bytes("abc"))
// Value: momento.String("abc")
type Value interface{ asBytes() []byte }

// Bytes plain old []byte
type Bytes struct {
	Bytes []byte
}

func (r Bytes) asBytes() []byte {
	return r.Bytes
}

// String string type that will be converted to []byte
type String struct {
	Text string
}

func (r String) asBytes() []byte {
	return []byte(r.Text)
}
