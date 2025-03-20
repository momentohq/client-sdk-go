package responses

// SetWithHashResponse is the base response type for a set with hash request.
type SetWithHashResponse interface {
	isSetWithHashResponse()
}

// SetWithHashNotStored NotStored response to a cache SetWithHash api request.
type SetWithHashNotStored struct{}

func (SetWithHashNotStored) isSetWithHashResponse() {}

// SetWithHashStored Stored response to a cache SetWithHash api request.
type SetWithHashStored struct {
	hash []byte
}

func (SetWithHashStored) isSetWithHashResponse() {}

// HashString Returns hash of value stored in cache as string for SetWithHashStored responses.
func (resp SetWithHashStored) HashString() string {
	return string(resp.hash)
}

// HashByte Returns hash of value stored in cache as bytes for SetWithHashStored responses.
func (resp SetWithHashStored) HashByte() []byte {
	return resp.hash
}

// NewSetWithHashStored returns a new SetWithHashStored containing the supplied value.
func NewSetWithHashStored(hash []byte) *SetWithHashStored {
	return &SetWithHashStored{hash: hash}
}
