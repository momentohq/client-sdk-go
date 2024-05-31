package momento

type StorePutRequest struct {
	StoreName string
	Key       string
	Value     StoreValue
}
