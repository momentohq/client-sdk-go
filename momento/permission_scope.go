package momento

type CacheRole int64

const (
	ReadWrite CacheRole = iota
	ReadOnly
	WriteOnly
)

type CacheSelector interface {
	IsAllCaches() bool
	CacheName() string
}

type AllCaches struct{}

func (AllCaches) IsAllCaches() bool {
	return true
}
func (AllCaches) CacheName() string {
	return ""
}

type CacheName struct {
	Name string
}

func (CacheName) IsAllCaches() bool {
	return false
}
func (c CacheName) CacheName() string {
	return c.Name
}

type CachePermission struct {
	Role  CacheRole
	Cache CacheSelector
}

func (CachePermission) IsPermission() {}

type TopicRole int64

const (
	PublishSubscribe TopicRole = iota
	SubscribeOnly
	PublishOnly
)

type TopicSelector interface {
	IsAllTopics() bool
	TopicName() string
}

type AllTopics struct{}

func (AllTopics) IsAllTopics() bool {
	return true
}
func (AllTopics) TopicName() string {
	return ""
}

type TopicName struct {
	Name string
}

func (TopicName) IsAllTopics() bool {
	return false
}
func (t TopicName) TopicName() string {
	return t.Name
}

type TopicNamePrefix struct {
	NamePrefix string
}

func (TopicNamePrefix) IsAllTopics() bool {
	return false
}
func (t TopicNamePrefix) TopicName() string {
	return t.NamePrefix
}

type TopicPermission struct {
	Role  TopicRole
	Cache CacheSelector
	Topic TopicSelector
}

func (TopicPermission) IsPermission() {}

type Permission interface {
	IsPermission()
}

type Permissions struct {
	Permissions []Permission
}

func (Permissions) IsPermissionScope() {}

var AllDataReadWrite = &Permissions{
	Permissions: []Permission{
		TopicPermission{Topic: AllTopics{}, Cache: AllCaches{}, Role: PublishSubscribe},
		CachePermission{Cache: AllCaches{}, Role: ReadWrite},
	},
}

type PredefinedScopeInterface interface {
	IsPredefinedScope()
}

type PredefinedScope struct{}

func (PredefinedScope) IsPredefinedScope() {}

func (PredefinedScope) IsPermissionScope() {}

type PermissionScope interface {
	IsPermissionScope()
}
