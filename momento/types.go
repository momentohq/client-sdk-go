package momento

// Bytes generic interface to help users deal with passing us bytes
type Bytes interface{ AsBytes() []byte }

// RawBytes plain old []byte
type RawBytes struct {
	Bytes []byte
}

func (r RawBytes) AsBytes() []byte {
	return r.Bytes
}

// StringBytes string type that will be converted to []byte
type StringBytes struct {
	Text string
}

func (r StringBytes) AsBytes() []byte {
	return []byte(r.Text)
}

type CacheResult int32

const (
	Hit  CacheResult = 2
	Miss CacheResult = 3
)
