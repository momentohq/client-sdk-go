package responses

// SetIfAbsentOrHashNotEqualResponse is the base response type
// for a set if absent or hash not equal request.
type SetIfAbsentOrHashNotEqualResponse interface {
	MomentoCacheResponse
	isSetIfAbsentOrHashNotEqualResponse()
}

// SetIfAbsentOrHashNotEqualNotStored NotStored response to a
// cache SetIfAbsentOrHashNotEqual api request.
type SetIfAbsentOrHashNotEqualNotStored struct{}

func (SetIfAbsentOrHashNotEqualNotStored) isSetIfAbsentOrHashNotEqualResponse() {}

// SetIfAbsentOrHashNotEqualStored Stored response to a
// cache SetIfAbsentOrHashNotEqual api request.
type SetIfAbsentOrHashNotEqualStored struct {
	hash []byte
}

func (SetIfAbsentOrHashNotEqualStored) isSetIfAbsentOrHashNotEqualResponse() {}

// HashString Returns hash of value stored in cache as string
// for SetIfAbsentOrHashNotEqualStored responses.
func (resp SetIfAbsentOrHashNotEqualStored) HashString() string {
	return string(resp.hash)
}

// HashByte Returns hash of value stored in cache as bytes
// for SetIfAbsentOrHashNotEqualStored responses.
func (resp SetIfAbsentOrHashNotEqualStored) HashByte() []byte {
	return resp.hash
}

// NewSetIfAbsentOrHashNotEqualStored returns a new
// SetIfAbsentOrHashNotEqualStored containing the supplied value.
func NewSetIfAbsentOrHashNotEqualStored(hash []byte) *SetIfAbsentOrHashNotEqualStored {
	return &SetIfAbsentOrHashNotEqualStored{hash: hash}
}
