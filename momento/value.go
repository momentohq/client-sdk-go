package momento

// Value Interface to help users deal with passing us values as strings or bytes.
// Value: momento.Bytes([]bytes("abc"))
// Value: momento.String("abc")
type Value interface {
	asBytes() []byte
	asString() string
}

// Key Type alias to future proof passing in keys.
type Key = Value

// Bytes plain old []byte
type Bytes []byte

func (v Bytes) asBytes() []byte { return v }

func (v Bytes) asString() string {
	return string(v)
}

// String string type that will be converted to []byte
type String string

func (v String) asBytes() []byte {
	return []byte(v)
}

func (v String) asString() string {
	return string(v)
}

// Type to hold field/value elements in dictionaries for use with DictionarySetFields.
// Field and Value are both Value type which allows both Strings and Bytes.
type DictionaryElement struct {
	Field, Value Value
}

// DictionaryElementsFromMap converts a map[string]string to an array of momento DictionaryElements.
//
//	DictionaryElements are used as input to DictionarySetFields.
func DictionaryElementsFromMap(theMap map[string]string) []DictionaryElement {
	return DictionaryElementsFromMapStringString(theMap)
}

// DictionaryElementsFromMapStringString converts a map[string]string to an array of momento DictionaryElements.
//
//	DictionaryElements are used as input to DictionarySetFields.
func DictionaryElementsFromMapStringString(theMap map[string]string) []DictionaryElement {
	var elements []DictionaryElement
	for k, v := range theMap {
		elements = append(elements, DictionaryElement{
			Field: String(k),
			Value: String(v),
		})
	}
	return elements
}

// DictionaryElementsFromMapStringBytes converts a map[string][]byte to an array of momento DictionaryElements.
//
//	DictionaryElements are used as input to DictionarySetFields.
func DictionaryElementsFromMapStringBytes(theMap map[string][]byte) []DictionaryElement {
	var elements []DictionaryElement
	for k, v := range theMap {
		elements = append(elements, DictionaryElement{
			Field: String(k),
			Value: Bytes(v),
		})
	}
	return elements
}

// DictionaryElementsFromMapStringValue converts a map[string]momento.Value to an array of momento DictionaryElements.
//
//	DictionaryElements are used as input to DictionarySetFields.
func DictionaryElementsFromMapStringValue(theMap map[string]Value) []DictionaryElement {
	var elements []DictionaryElement
	for k, v := range theMap {
		elements = append(elements, DictionaryElement{
			Field: String(k),
			Value: v,
		})
	}
	return elements
}

type SortedSetElement struct {
	Value Value
	Score float64
}

// SortedSetElementsFromMap converts a map[string]float64 to an array of momento SortedSetElements.
//
//	SortedSetElements are used as input to SortedSetPutElements.
func SortedSetElementsFromMap(theMap map[string]float64) []SortedSetElement {
	var elements []SortedSetElement
	for k, v := range theMap {
		elements = append(elements, SortedSetElement{
			Value: String(k),
			Score: v,
		})
	}
	return elements
}
