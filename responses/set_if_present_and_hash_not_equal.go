package responses

// SetIfPresentAndHashNotEqualResponse is the base response type
// for a set if present and hash not equal request.
type SetIfPresentAndHashNotEqualResponse interface {
	MomentoCacheResponse
	isSetIfPresentAndHashNotEqualResponse()
}

// SetIfPresentAndHashNotEqualNotStored NotStored response to a
// cache SetIfPresentAndHashNotEqual api request.
type SetIfPresentAndHashNotEqualNotStored struct{}

func (SetIfPresentAndHashNotEqualNotStored) isSetIfPresentAndHashNotEqualResponse() {}

// SetIfPresentAndHashNotEqualStored Stored response to a
// cache SetIfPresentAndHashNotEqual api request.
type SetIfPresentAndHashNotEqualStored struct {
	hash []byte
}

func (SetIfPresentAndHashNotEqualStored) isSetIfPresentAndHashNotEqualResponse() {}

// HashString Returns hash of value stored in cache as string for
// SetIfPresentAndHashNotEqualStored responses.
func (resp SetIfPresentAndHashNotEqualStored) HashString() string {
	return string(resp.hash)
}

// HashByte Returns hash of value stored in cache as bytes for
// SetIfPresentAndHashNotEqualStored responses.
func (resp SetIfPresentAndHashNotEqualStored) HashByte() []byte {
	return resp.hash
}

// NewSetIfPresentAndHashNotEqualStored returns a new
// SetIfPresentAndHashNotEqualStored containing the supplied value.
func NewSetIfPresentAndHashNotEqualStored(hash []byte) *SetIfPresentAndHashNotEqualStored {
	return &SetIfPresentAndHashNotEqualStored{hash: hash}
}
