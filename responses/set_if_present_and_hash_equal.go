package responses

// SetIfPresentAndHashEqualResponse is the base response type
// for a set if present and hash equal request.
type SetIfPresentAndHashEqualResponse interface {
	isSetIfPresentAndHashEqualResponse()
}

// SetIfPresentAndHashEqualNotStored NotStored response to a
// cache SetIfPresentAndHashEqual api request.
type SetIfPresentAndHashEqualNotStored struct{}

func (SetIfPresentAndHashEqualNotStored) isSetIfPresentAndHashEqualResponse() {}

// SetIfPresentAndHashEqualStored Stored response to a
// cache SetIfPresentAndHashEqual api request.
type SetIfPresentAndHashEqualStored struct {
	hash []byte
}

func (SetIfPresentAndHashEqualStored) isSetIfPresentAndHashEqualResponse() {}

// HashString Returns hash of value stored in cache as string for
// SetIfPresentAndHashEqualStored responses.
func (resp SetIfPresentAndHashEqualStored) HashString() string {
	return string(resp.hash)
}

// HashByte Returns hash of value stored in cache as bytes for
// SetIfPresentAndHashEqualStored responses.
func (resp SetIfPresentAndHashEqualStored) HashByte() []byte {
	return resp.hash
}

// NewSetIfPresentAndHashEqualStored returns a new
// SetIfPresentAndHashEqualStored containing the supplied value.
func NewSetIfPresentAndHashEqualStored(hash []byte) *SetIfPresentAndHashEqualStored {
	return &SetIfPresentAndHashEqualStored{hash: hash}
}
