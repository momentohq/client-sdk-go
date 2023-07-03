package responses

// SetIfNotExistsResponse is the base response type for a SetIfNotExists request
type SetIfNotExistsResponse interface {
	isSetIfNotExistsResponse()
}

// SetIfNotExistsNotStored indicates a successful set request where the value was already present.
type SetIfNotExistsNotStored struct{}

func (SetIfNotExistsNotStored) isSetIfNotExistsResponse() {}

// SetIfNotExistsStored indicates a successful set request where the value was stored.
type SetIfNotExistsStored struct {
	key   []byte
	value []byte
}

func (SetIfNotExistsStored) isSetIfNotExistsResponse() {}

// ValueString returns value stored in the cache.
func (resp SetIfNotExistsStored) ValueString() string {
	return string(resp.value)
}

// ValueByte Returns value stored in the cache as bytes.
func (resp SetIfNotExistsStored) ValueByte() []byte {
	return resp.value
}

// KeyString returns key stored in the cache.
func (resp SetIfNotExistsStored) KeyString() string {
	return string(resp.key)
}

// KeyByte Returns key stored in the cache as bytes.
func (resp SetIfNotExistsStored) KeyByte() []byte {
	return resp.key
}

// NewSetIfNotExistsStored returns a new SetIfNotExistsStored containing the supplied key and value.
func NewSetIfNotExistsStored(key, value []byte) *SetIfNotExistsStored {
	return &SetIfNotExistsStored{key: key, value: value}
}
