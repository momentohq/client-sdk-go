package responses

type StorageValueType string

const (
	BYTES   StorageValueType = "BYTES"
	STRING  StorageValueType = "STRING"
	INTEGER StorageValueType = "INTEGER"
	DOUBLE  StorageValueType = "DOUBLE"
)

type StorageGetResponse interface {
	isStorageGetResponse()
}

// StorageGetNotFound indicates the item was not found in the store
type StorageGetNotFound struct{}

func (StorageGetNotFound) isStorageGetResponse() {}

// StorageGetFound indicates the item was found in the store
type StorageGetFound struct {
	valueType    StorageValueType
	valueBytes   *[]byte
	valueString  *string
	valueFloat   *float64
	valueInteger *int
}

func (StorageGetFound) isStorageGetResponse() {}

// ValueType returns the `StorageValueType` indicating the type of the value in the store.
func (resp StorageGetFound) ValueType() StorageValueType {
	return resp.valueType
}

// ValueString returns the value in the store as a string and a boolean `true` value if it was stored as a string. Otherwise, it returns a blank string and a boolean `false` value.
func (resp StorageGetFound) ValueString() (string, bool) {
	if resp.valueType == STRING {
		return *resp.valueString, true
	}
	return "", false
}

// ValueBytes returns the value in the store as a byte slice and a boolean `true` value if it was stored as a bytes. Otherwise, it returns a nil byte slice and a boolean `false` value.
func (resp StorageGetFound) ValueBytes() ([]byte, bool) {
	if resp.valueType == BYTES {
		return *resp.valueBytes, true
	}
	return nil, false
}

// ValueFloat returns the value in the store as a float64 and a boolean `true` value if it was stored as a double. Otherwise, it returns 0 and a boolean `false` value.
func (resp StorageGetFound) ValueFloat() (float64, bool) {
	if resp.valueType == DOUBLE {
		return *resp.valueFloat, true
	}
	return 0, false
}

// ValueInteger returns the value in the store as an int and a boolean `true` value if it was stored as an integer. Otherwise, it returns 0 and a boolean `false` value.
func (resp StorageGetFound) ValueInteger() (int, bool) {
	if resp.valueType == INTEGER {
		return *resp.valueInteger, true
	}
	return 0, false
}

// NewStorageGetFound_String returns a new StorageGetFound containing the supplied string value.
func NewStorageGetFound_String(valueType StorageValueType, value string) *StorageGetFound {
	return &StorageGetFound{
		valueType:   valueType,
		valueString: &value,
	}
}

// NewStorageGetFound_Bytes returns a new StorageGetFound containing the supplied byte slice value.
func NewStorageGetFound_Bytes(valueType StorageValueType, value []byte) *StorageGetFound {
	return &StorageGetFound{
		valueType:  valueType,
		valueBytes: &value,
	}
}

// NewStorageGetFound_Float returns a new StorageGetFound containing the supplied float64 value.
func NewStorageGetFound_Float(valueType StorageValueType, value float64) *StorageGetFound {
	return &StorageGetFound{
		valueType:  valueType,
		valueFloat: &value,
	}
}

// NewStorageGetFound_Integer returns a new StorageGetFound containing the supplied int value.
func NewStorageGetFound_Integer(valueType StorageValueType, value int) *StorageGetFound {
	return &StorageGetFound{
		valueType:    valueType,
		valueInteger: &value,
	}
}
