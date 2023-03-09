package responses

// KeysExistResponse is a base response type for a key exist request.
type KeysExistResponse interface {
	isKeysExistResponse()
}

// KeysExistSuccess indicates a successful key exist request.
type KeysExistSuccess struct {
	values []bool
}

func (r *KeysExistSuccess) isKeysExistResponse() {}

// Exists returns an array of bool to indicate existence the given values.
func (r *KeysExistSuccess) Exists() []bool {
	return r.values
}

// NewKeysExistSuccess returns a new KeysExistSuccess contains values.
func NewKeysExistSuccess(values []bool) *KeysExistSuccess {
	return &KeysExistSuccess{values: values}
}
