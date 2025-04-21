package responses

// DictionaryGetFieldResponse is the base response type for a dictionary get field request.
type DictionaryGetFieldResponse interface {
	MomentoCacheResponse
	isDictionaryGetFieldResponse()
}

// DictionaryGetFieldHit Indicates that the requested data was successfully retrieved from the cache.  Provides
// `Value*` accessors to retrieve the data in the appropriate format.
type DictionaryGetFieldHit struct {
	field []byte
	body  []byte
}

func (DictionaryGetFieldHit) isDictionaryGetFieldResponse() {}

// FieldString returns the field name for the retrieved element, as an utf-8 string decoded from the underlying byte array.
func (resp DictionaryGetFieldHit) FieldString() string {
	return string(resp.field)
}

// FieldByte returns the field name for the retrieved element, as a byte array.
func (resp DictionaryGetFieldHit) FieldByte() []byte {
	return resp.field
}

// ValueString returns the data as a utf-8 string, decoded from the underlying byte array.
func (resp DictionaryGetFieldHit) ValueString() string {
	return string(resp.body)
}

// ValueByte returns the data as a byte array.
func (resp DictionaryGetFieldHit) ValueByte() []byte {
	return resp.body
}

// DictionaryGetFieldMiss indicates that the requested data was not available in the cache.
type DictionaryGetFieldMiss struct {
	field []byte
}

func (DictionaryGetFieldMiss) isDictionaryGetFieldResponse() {}

// FieldString returns the field name for the retrieved element, as an utf-8 string decoded from the underlying byte array.
func (resp DictionaryGetFieldMiss) FieldString() string {
	return string(resp.field)
}

// FieldByte returns the field name for the retrieved element, as a byte array.
func (resp DictionaryGetFieldMiss) FieldByte() []byte {
	return resp.field
}

// NewDictionaryGetFieldHit returns a new DictionaryGetFieldHit which contains field and body.
func NewDictionaryGetFieldHit(field []byte, body []byte) *DictionaryGetFieldHit {
	return &DictionaryGetFieldHit{
		field: field,
		body:  body,
	}
}

// NewDictionaryGetFieldMiss returns a new DictionaryGetFieldMiss which contains the requested field.
func NewDictionaryGetFieldMiss(field []byte) *DictionaryGetFieldMiss {
	return &DictionaryGetFieldMiss{field: field}
}

// NewDictionaryGetFieldHitFromFieldsHit returns a new DictionaryGetFieldHit containing the first field and element's body.
// This is used specifically to provide a hit response for getting a single dictionary field.
func NewDictionaryGetFieldHitFromFieldsHit(fieldHit *DictionaryGetFieldsHit) *DictionaryGetFieldHit {
	return &DictionaryGetFieldHit{
		field: fieldHit.fields[0],
		body:  fieldHit.elements[0].CacheBody,
	}
}
