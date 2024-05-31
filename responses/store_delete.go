package responses

type StoreDeleteResponse interface {
	isStoreDeleteResponse()
}

type StoreDeleteSuccess struct{}

func (StoreDeleteSuccess) isStoreDeleteResponse() {}
