package momento

type StoreGetRequest struct {
	// StoreName is the name of the store to get the value from.
	StoreName string
	// Key is the key to get the value for.
	Key string
}
