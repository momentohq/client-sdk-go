package momento

type StoreDeleteRequest struct {
	// StoreName is the name of the store to delete.
	StoreName string
	// Key is the key to delete from the store.
	Key string
}
