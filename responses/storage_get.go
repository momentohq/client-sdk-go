package responses

type StorageValueType string

const (
	BYTES   StorageValueType = "BYTES"
	STRING  StorageValueType = "STRING"
	INTEGER StorageValueType = "INTEGER"
	DOUBLE  StorageValueType = "DOUBLE"
)

// StorageGetResponse is the base response type for a store get request.
type StorageGetResponse interface {
	isStoreGetResponse()
	ValueType() StorageValueType
	ValueString() (string, bool)
	ValueBytes() ([]byte, bool)
	ValueInteger() (int, bool)
	ValueFloat64() (float64, bool)
}

// StorageGetSuccess indicates a successful store get request.
type StorageGetSuccess struct {
	valueType    StorageValueType
	valueBytes   *[]byte
	valueString  *string
	valueFloat64 *float64
	valueInteger *int
}

func (StorageGetSuccess) isStoreGetResponse() {}

// ValueType returns the `StorageValueType` indicating the type of the value in the store.
func (resp StorageGetSuccess) ValueType() StorageValueType {
	return resp.valueType
}

// ValueString returns the value in the store as a string and a boolean `true` value if it was stored as a string. Otherwise, it returns a blank string and a boolean `false` value.
func (resp StorageGetSuccess) ValueString() (string, bool) {
	if resp.valueType == STRING {
		return *resp.valueString, true
	}
	return "", false
}

// ValueBytes returns the value in the store as a byte slice and a boolean `true` value if it was stored as a bytes. Otherwise, it returns a nil byte slice and a boolean `false` value.
func (resp StorageGetSuccess) ValueBytes() ([]byte, bool) {
	if resp.valueType == BYTES {
		return *resp.valueBytes, true
	}
	return nil, false
}

// ValueFloat64 returns the value in the store as a float64 and a boolean `true` value if it was stored as a double. Otherwise, it returns 0 and a boolean `false` value.
func (resp StorageGetSuccess) ValueFloat64() (float64, bool) {
	if resp.valueType == DOUBLE {
		return *resp.valueFloat64, true
	}
	return 0, false
}

// ValueInteger returns the value in the store as an int and a boolean `true` value if it was stored as an integer. Otherwise, it returns 0 and a boolean `false` value.
func (resp StorageGetSuccess) ValueInteger() (int, bool) {
	if resp.valueType == INTEGER {
		return *resp.valueInteger, true
	}
	return 0, false
}

// NewStoreGetSuccess_String returns a new StorageGetSuccess containing the supplied string value.
func NewStoreGetSuccess_String(valueType StorageValueType, value string) *StorageGetSuccess {
	return &StorageGetSuccess{
		valueType:   valueType,
		valueString: &value,
	}
}

// NewStoreGetSuccess_Bytes returns a new StorageGetSuccess containing the supplied byte slice value.
func NewStoreGetSuccess_Bytes(valueType StorageValueType, value []byte) *StorageGetSuccess {
	return &StorageGetSuccess{
		valueType:  valueType,
		valueBytes: &value,
	}
}

// NewStoreGetSuccess_Float64 returns a new StorageGetSuccess containing the supplied float64 value.
func NewStoreGetSuccess_Float64(valueType StorageValueType, value float64) *StorageGetSuccess {
	return &StorageGetSuccess{
		valueType:    valueType,
		valueFloat64: &value,
	}
}

// NewStoreGetSuccess_Integer returns a new StorageGetSuccess containing the supplied int value.
func NewStoreGetSuccess_Integer(valueType StorageValueType, value int) *StorageGetSuccess {
	return &StorageGetSuccess{
		valueType:    valueType,
		valueInteger: &value,
	}
}
