package momento

type CacheItemKey struct {
	Key []byte
}

func (CacheItemKey) IsCacheItemSelector() {}

type CacheItemKeyPrefix struct {
	KeyPrefix []byte
}

func (CacheItemKeyPrefix) IsCacheItemSelector() {}

type AllCacheItems struct{}

func (AllCacheItems) IsCacheItemSelector() {}

type CacheItemSelector interface {
	IsCacheItemSelector()
}

type DisposableTokenCachePermission struct {
	Role  CacheRole
	Cache CacheSelector
	Item  CacheItemSelector
}

func (DisposableTokenCachePermission) IsPermission() {}

type DisposableTokenCachePermissions struct {
	Permissions []DisposableTokenCachePermission
}

func (DisposableTokenCachePermissions) IsDisposableTokenScope() {}
func (Permissions) IsDisposableTokenScope()                     {}

type DisposableTokenScope interface {
	IsDisposableTokenScope()
}

type DisposableTokenProps struct {
	TokenId *string
}
