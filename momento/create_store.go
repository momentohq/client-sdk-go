package momento

type CreateStoreRequest struct {
	// string used to create a store.
	StoreName string
}

func (c CreateStoreRequest) cacheName() string {
	return c.StoreName
}
