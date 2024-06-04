package momento

type StorePutRequest struct {
	// StoreName is the name of the store to put the value in.
	StoreName string
	// Key is the key to put the value for.
	Key string
	// Value is the `StoreValue` value to put in the store.
	Value StoreValue
}
