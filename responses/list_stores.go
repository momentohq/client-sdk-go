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

// NextToken Next Page Token returned by Simple Store Service along with the list of stores.
// If nextToken is present, then this token must be provided in the next call to continue paginating through the list.
// This is done by setting this value in ListStoresRequest.
func (resp ListStoresSuccess) NextToken() string {
	return resp.nextToken
}

// Stores Returns all stores.
func (resp ListStoresSuccess) Stores() []StoreInfo {
	return resp.stores
}

// StoreInfo Information about a Store.
type StoreInfo struct {
	name string
}

// Name Returns cache's name.
func (si StoreInfo) Name() string {
	return si.name
}

// NewStoreInfo returns new CacheInfo with the supplied name.
func NewStoreInfo(name string) StoreInfo {
	return StoreInfo{name: name}
}
