package responses

type KeysExistResponse interface {
	isKeysExistResponse()
}

type KeysExistSuccess struct {
	values []bool
}

func (r *KeysExistSuccess) isKeysExistResponse() {}

func (r *KeysExistSuccess) Exists() []bool {
	return r.values
}

func NewKeysExistSuccess(values []bool) *KeysExistSuccess {
	return &KeysExistSuccess{values: values}
}
