package momento

type DeleteStoreRequest struct {
	// string store name to delete.
	StoreName string
}

func (c DeleteStoreRequest) cacheName() string {
	return c.StoreName
}
