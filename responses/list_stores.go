package responses

// ListStoresResponse is the base response type for a list stores request.
type ListStoresResponse interface {
	isListStoresResponse()
}

// ListStoresSuccess Output of the List stores operation.
type ListStoresSuccess struct {
	nextToken string
	stores    []StoreInfo
}

func (ListStoresSuccess) isListStoresResponse() {}

// NewListStoresSuccess returns a new ListStoresSuccess which indicates a successful list stores request.
func NewListStoresSuccess(nextToken string, stores []StoreInfo) *ListStoresSuccess {
	return &ListStoresSuccess{
		nextToken: nextToken,
		stores:    stores,
	}
}

// Stores Returns all stores.
func (resp ListStoresSuccess) Stores() []StoreInfo {
	return resp.stores
}

// StoreInfo Information about a Store.
type StoreInfo struct {
	name string
}

// Name Returns store's name.
func (si StoreInfo) Name() string {
	return si.name
}

// NewStoreInfo returns new StoreInfo with the supplied name.
func NewStoreInfo(name string) StoreInfo {
	return StoreInfo{name: name}
}
