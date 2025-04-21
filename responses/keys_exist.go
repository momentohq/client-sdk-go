package responses

// KeysExistResponse is the base response type for a key exist request.
type KeysExistResponse interface {
	MomentoCacheResponse
	isKeysExistResponse()
}

// KeysExistSuccess indicates a successful key exist request.
type KeysExistSuccess struct {
	values []bool
}

func (r *KeysExistSuccess) isKeysExistResponse() {}

// Exists returns an array of bool to indicate whether or not the given keys exist in the cache.
func (r *KeysExistSuccess) Exists() []bool {
	return r.values
}

// NewKeysExistSuccess returns a new KeysExistSuccess containing the supplied values.
func NewKeysExistSuccess(values []bool) *KeysExistSuccess {
	return &KeysExistSuccess{values: values}
}
