package momento

// Bytes generic interface to help users deal with passing us bytes
type Bytes interface{ asBytes() []byte }

// RawBytes plain old []byte
type RawBytes struct {
	Bytes []byte
}

func (r RawBytes) asBytes() []byte {
	return r.Bytes
}

// StringBytes string type that will be converted to []byte
type StringBytes struct {
	Text string
}

func (r StringBytes) asBytes() []byte {
	return []byte(r.Text)
}
