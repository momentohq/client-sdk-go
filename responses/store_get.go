package responses

type StoreValueType string

const (
	BYTES   StoreValueType = "BYTES"
	STRING                 = "STRING"
	INTEGER                = "INTEGER"
	DOUBLE                 = "DOUBLE"
)

// StoreGetResponse is the base response type for a store get request.
type StoreGetResponse interface {
	isStoreGetResponse()
	ValueType() StoreValueType
	TryGetValueString() (string, bool)
	TryGetValueBytes() ([]byte, bool)
	TryGetValueInteger() (int, bool)
	TryGetValueDouble() (float64, bool)
}

// StoreGetSuccess indicates a successful store get request.
type StoreGetSuccess struct {
	valueType    StoreValueType
	valueBytes   *[]byte
	valueString  *string
	valueDouble  *float64
	valueInteger *int
}

func (StoreGetSuccess) isStoreGetResponse() {}

// ValueType returns the `StoreValueType` indicating the type of the value in the store.
func (resp StoreGetSuccess) ValueType() StoreValueType {
	return resp.valueType
}

// TryGetValueString returns the value in the store as a string and a boolean `true` value if it was stored as a string. Otherwise, it returns a blank string and a boolean `false` value.
func (resp StoreGetSuccess) TryGetValueString() (string, bool) {
	if resp.valueType == STRING {
		return *resp.valueString, true
	}
	// TODO: Should we return a blank string or return pointers instead so we can return nil?
	return "", false
}

// TryGetValueBytes returns the value in the store as a byte slice and a boolean `true` value if it was stored as a bytes. Otherwise, it returns a nil byte slice and a boolean `false` value.
func (resp StoreGetSuccess) TryGetValueBytes() ([]byte, bool) {
	if resp.valueType == BYTES {
		return *resp.valueBytes, true
	}
	return nil, false
}

// TryGetValueDouble returns the value in the store as a float64 and a boolean `true` value if it was stored as a double. Otherwise, it returns 0 and a boolean `false` value.
func (resp StoreGetSuccess) TryGetValueDouble() (float64, bool) {
	if resp.valueType == DOUBLE {
		return *resp.valueDouble, true
	}
	return 0, false
}

// TryGetValueInteger returns the value in the store as an int and a boolean `true` value if it was stored as an integer. Otherwise, it returns 0 and a boolean `false` value.
func (resp StoreGetSuccess) TryGetValueInteger() (int, bool) {
	if resp.valueType == INTEGER {
		return *resp.valueInteger, true
	}
	return 0, false
}

// NewStoreGetSuccess_String returns a new StoreGetSuccess containing the supplied string value.
func NewStoreGetSuccess_String(valueType StoreValueType, value string) *StoreGetSuccess {
	return &StoreGetSuccess{
		valueType:   valueType,
		valueString: &value,
	}
}

// NewStoreGetSuccess_Bytes returns a new StoreGetSuccess containing the supplied byte slice value.
func NewStoreGetSuccess_Bytes(valueType StoreValueType, value []byte) *StoreGetSuccess {
	return &StoreGetSuccess{
		valueType:  valueType,
		valueBytes: &value,
	}
}

// NewStoreGetSuccess_Double returns a new StoreGetSuccess containing the supplied float64 value.
func NewStoreGetSuccess_Double(valueType StoreValueType, value float64) *StoreGetSuccess {
	return &StoreGetSuccess{
		valueType:   valueType,
		valueDouble: &value,
	}
}

// NewStoreGetSuccess_Integer returns a new StoreGetSuccess containing the supplied int value.
func NewStoreGetSuccess_Integer(valueType StoreValueType, value int) *StoreGetSuccess {
	return &StoreGetSuccess{
		valueType:    valueType,
		valueInteger: &value,
	}
}
