package momento_rpc_names

// MomentoRPCMethod represents the RPC methods used in Momento
type MomentoRPCMethod string

const (
	Get            MomentoRPCMethod = "/cache_client.Scs/Get"
	GetWithHash    MomentoRPCMethod = "/cache_client.Scs/GetWithHash"
	Set            MomentoRPCMethod = "/cache_client.Scs/Set"
	SetIfHash      MomentoRPCMethod = "/cache_client.Scs/SetIfHash"
	Delete         MomentoRPCMethod = "/cache_client.Scs/Delete"
	Increment      MomentoRPCMethod = "/cache_client.Scs/Increment"
	SetIf          MomentoRPCMethod = "/cache_client.Scs/SetIf"
	SetIfNotExists MomentoRPCMethod = "/cache_client.Scs/SetIfNotExists"
	GetBatch       MomentoRPCMethod = "/cache_client.Scs/GetBatch"
	SetBatch       MomentoRPCMethod = "/cache_client.Scs/SetBatch"
	KeysExist      MomentoRPCMethod = "/cache_client.Scs/KeysExist"
	UpdateTtl      MomentoRPCMethod = "/cache_client.Scs/UpdateTtl"
	ItemGetTtl     MomentoRPCMethod = "/cache_client.Scs/ItemGetTtl"
	ItemGetType    MomentoRPCMethod = "/cache_client.Scs/ItemGetType"

	// Dictionary operations
	DictionarySet       MomentoRPCMethod = "/cache_client.Scs/DictionarySet"
	DictionaryGet       MomentoRPCMethod = "/cache_client.Scs/DictionaryGet"
	DictionaryFetch     MomentoRPCMethod = "/cache_client.Scs/DictionaryFetch"
	DictionaryIncrement MomentoRPCMethod = "/cache_client.Scs/DictionaryIncrement"
	DictionaryDelete    MomentoRPCMethod = "/cache_client.Scs/DictionaryDelete"
	DictionaryLength    MomentoRPCMethod = "/cache_client.Scs/DictionaryLength"

	// Set operations
	SetFetch      MomentoRPCMethod = "/cache_client.Scs/SetFetch"
	SetSample     MomentoRPCMethod = "/cache_client.Scs/SetSample"
	SetUnion      MomentoRPCMethod = "/cache_client.Scs/SetUnion"
	SetDifference MomentoRPCMethod = "/cache_client.Scs/SetDifference"
	SetContains   MomentoRPCMethod = "/cache_client.Scs/SetContains"
	SetLength     MomentoRPCMethod = "/cache_client.Scs/SetLength"
	SetPop        MomentoRPCMethod = "/cache_client.Scs/SetPop"

	// List operations
	ListPushFront        MomentoRPCMethod = "/cache_client.Scs/ListPushFront"
	ListPushBack         MomentoRPCMethod = "/cache_client.Scs/ListPushBack"
	ListPopFront         MomentoRPCMethod = "/cache_client.Scs/ListPopFront"
	ListPopBack          MomentoRPCMethod = "/cache_client.Scs/ListPopBack"
	ListErase            MomentoRPCMethod = "/cache_client.Scs/ListErase"
	ListRemove           MomentoRPCMethod = "/cache_client.Scs/ListRemove"
	ListFetch            MomentoRPCMethod = "/cache_client.Scs/ListFetch"
	ListLength           MomentoRPCMethod = "/cache_client.Scs/ListLength"
	ListConcatenateFront MomentoRPCMethod = "/cache_client.Scs/ListConcatenateFront"
	ListConcatenateBack  MomentoRPCMethod = "/cache_client.Scs/ListConcatenateBack"
	ListRetain           MomentoRPCMethod = "/cache_client.Scs/ListRetain"

	// Sorted Set operations
	SortedSetPut           MomentoRPCMethod = "/cache_client.Scs/SortedSetPut"
	SortedSetFetch         MomentoRPCMethod = "/cache_client.Scs/SortedSetFetch"
	SortedSetGetScore      MomentoRPCMethod = "/cache_client.Scs/SortedSetGetScore"
	SortedSetRemove        MomentoRPCMethod = "/cache_client.Scs/SortedSetRemove"
	SortedSetIncrement     MomentoRPCMethod = "/cache_client.Scs/SortedSetIncrement"
	SortedSetGetRank       MomentoRPCMethod = "/cache_client.Scs/SortedSetGetRank"
	SortedSetLength        MomentoRPCMethod = "/cache_client.Scs/SortedSetLength"
	SortedSetLengthByScore MomentoRPCMethod = "/cache_client.Scs/SortedSetLengthByScore"

	// Topic operations
	TopicPublish   MomentoRPCMethod = "/cache_client.pubsub.Pubsub/Publish"
	TopicSubscribe MomentoRPCMethod = "/cache_client.pubsub.Pubsub/Subscribe"
)
