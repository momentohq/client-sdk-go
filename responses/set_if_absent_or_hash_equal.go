package responses

// SetIfAbsentOrHashEqualResponse is the base response type
// for a set if absent or hash equal request.
type SetIfAbsentOrHashEqualResponse interface {
	MomentoCacheResponse
	isSetIfAbsentOrHashEqualResponse()
}

// SetIfAbsentOrHashEqualNotStored NotStored response to a
// cache SetIfAbsentOrHashEqual api request.
type SetIfAbsentOrHashEqualNotStored struct{}

func (SetIfAbsentOrHashEqualNotStored) isSetIfAbsentOrHashEqualResponse() {}

// SetIfAbsentOrHashEqualStored Stored response to a
// cache SetIfAbsentOrHashEqual api request.
type SetIfAbsentOrHashEqualStored struct {
	hash []byte
}

func (SetIfAbsentOrHashEqualStored) isSetIfAbsentOrHashEqualResponse() {}

// HashString Returns hash of value stored in cache as string
// for SetIfAbsentOrHashEqualStored responses.
func (resp SetIfAbsentOrHashEqualStored) HashString() string {
	return string(resp.hash)
}

// HashByte Returns hash of value stored in cache as bytes
// for SetIfAbsentOrHashEqualStored responses.
func (resp SetIfAbsentOrHashEqualStored) HashByte() []byte {
	return resp.hash
}

// NewSetIfAbsentOrHashEqualStored returns a new
// SetIfAbsentOrHashEqualStored containing the supplied value.
func NewSetIfAbsentOrHashEqualStored(hash []byte) *SetIfAbsentOrHashEqualStored {
	return &SetIfAbsentOrHashEqualStored{hash: hash}
}
