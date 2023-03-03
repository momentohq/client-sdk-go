package responses

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

func NewDictionaryGetFieldHit(field []byte, body []byte) *DictionaryGetFieldHit {
	return &DictionaryGetFieldHit{
		field: field,
		body:  body,
	}
}

func NewDictionaryGetFieldMiss(field []byte) *DictionaryGetFieldMiss {
	return &DictionaryGetFieldMiss{field: field}
}

func NewDictionaryGetFieldHitFromFieldsHit(fieldHit *DictionaryGetFieldsHit) *DictionaryGetFieldHit {
	return &DictionaryGetFieldHit{
		field: fieldHit.fields[0],
		body:  fieldHit.elements[0].CacheBody,
	}
}
