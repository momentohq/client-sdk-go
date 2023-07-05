package responses

// DictionaryLengthResponse is the base response type for a dictionary length request.
type DictionaryLengthResponse interface {
	isDictionaryLengthResponse()
}

// DictionaryLengthHit indicates a dictionary length request was a hit.
type DictionaryLengthHit struct {
	value uint32
}

func (DictionaryLengthHit) isDictionaryLengthResponse() {}

// Length returns the length of the dictionary.
func (resp DictionaryLengthHit) Length() uint32 {
	return resp.value
}

// DictionaryLengthMiss indicates a dictionary length request was a miss.
type DictionaryLengthMiss struct{}

func (DictionaryLengthMiss) isDictionaryLengthResponse() {}

// NewDictionaryLengthHit returns a new DictionaryLengthHit containing the supplied value.
func NewDictionaryLengthHit(value uint32) *DictionaryLengthHit {
	return &DictionaryLengthHit{value: value}
}
