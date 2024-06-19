package momento

type StorageDeleteRequest struct {
	// StoreName is the name of the store to delete the key from.
	StoreName string
	// Key is the key to delete from the store.
	Key string
}
