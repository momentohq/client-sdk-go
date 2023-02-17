package momento

// DictionaryGetFieldResponse

type DictionaryGetFieldResponse interface {
	isDictionaryGetFieldResponse()
}

type DictionaryGetFieldHit struct {
	field []byte
	body  []byte
}

func (DictionaryGetFieldHit) isDictionaryGetFieldResponse() {}

func (resp DictionaryGetFieldHit) FieldString() string {
	return string(resp.field)
}

func (resp DictionaryGetFieldHit) FieldByte() []byte {
	return resp.field
}

func (resp DictionaryGetFieldHit) ValueString() string {
	return string(resp.body)
}

func (resp DictionaryGetFieldHit) ValueByte() []byte {
	return resp.body
}

type DictionaryGetFieldMiss struct {
	field []byte
}

func (DictionaryGetFieldMiss) isDictionaryGetFieldResponse() {}

func (resp DictionaryGetFieldMiss) FieldString() string {
	return string(resp.field)
}

func (resp DictionaryGetFieldMiss) FieldByte() []byte {
	return resp.field
}

// DictionaryGetFieldRequest

type DictionaryGetFieldRequest struct {
	CacheName      string
	DictionaryName string
	Field          Value
}
