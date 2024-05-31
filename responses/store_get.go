package responses

type StoreValueType string

// TODO: prefix these names? Probably don't want to squat on these for the entire responses namespace?
const (
	BYTES   StoreValueType = "BYTES"
	STRING                 = "STRING"
	INTEGER                = "INTEGER"
	DOUBLE                 = "DOUBLE"
)

type StoreGetResponse interface {
	isStoreGetResponse()
}

type StoreGetSuccess struct {
	ValueType    StoreValueType
	valueBytes   *[]byte
	valueString  *string
	valueDouble  *float64
	valueInteger *int
}

func (StoreGetSuccess) isStoreGetResponse() {}

func (resp StoreGetSuccess) TryGetValueString() (string, bool) {
	if resp.ValueType == STRING {
		return *resp.valueString, true
	}
	// TODO
	// If these returned pointers instead of values, we could return nil
	// for the first return value and get rid of the bool.
	return "", false
}

func (resp StoreGetSuccess) TryGetValueBytes() ([]byte, bool) {
	if resp.ValueType == BYTES {
		return *resp.valueBytes, true
	}
	return nil, false
}

func (resp StoreGetSuccess) TryGetValueDouble() (float64, bool) {
	if resp.ValueType == DOUBLE {
		return *resp.valueDouble, true
	}
	return 0, false
}

func (resp StoreGetSuccess) TryGetValueInteger() (int, bool) {
	if resp.ValueType == INTEGER {
		return *resp.valueInteger, true
	}
	return 0, false
}

func NewStoreGetSuccess_String(valueType StoreValueType, value string) *StoreGetSuccess {
	return &StoreGetSuccess{
		ValueType:   valueType,
		valueString: &value,
	}
}

func NewStoreGetSuccess_Bytes(valueType StoreValueType, value []byte) *StoreGetSuccess {
	return &StoreGetSuccess{
		ValueType:  valueType,
		valueBytes: &value,
	}
}

func NewStoreGetSuccess_Double(valueType StoreValueType, value float64) *StoreGetSuccess {
	return &StoreGetSuccess{
		ValueType:   valueType,
		valueDouble: &value,
	}
}

func NewStoreGetSuccess_Integer(valueType StoreValueType, value int) *StoreGetSuccess {
	return &StoreGetSuccess{
		ValueType:    valueType,
		valueInteger: &value,
	}
}
