package momento

type StorageSetRequest struct {
	// StoreName is the name of the store to put the value in.
	StoreName string
	// Key is the key to put the value for.
	Key string
	// Value is the `StorageValue` value to put in the store.
	Value StorageValue
}
