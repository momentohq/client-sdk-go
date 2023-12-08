package momento

func CacheKeyReadWrite(selector CacheSelector, key Value) DisposableTokenScope {
	return Permissions{
		Permissions: []Permission{DisposableTokenCachePermission{Cache: selector, Item: CacheItemKey{Key: key.asBytes()}, Role: ReadWrite}},
	}
}

func CacheKeyPrefixReadWrite(selector CacheSelector, keyPrefix Value) DisposableTokenScope {
	return Permissions{
		Permissions: []Permission{DisposableTokenCachePermission{Cache: selector, Item: CacheItemKeyPrefix{KeyPrefix: keyPrefix.asBytes()}, Role: ReadWrite}},
	}
}

func CacheKeyReadOnly(selector CacheSelector, key Value) DisposableTokenScope {
	return Permissions{
		Permissions: []Permission{DisposableTokenCachePermission{Cache: selector, Item: CacheItemKey{Key: key.asBytes()}, Role: ReadOnly}},
	}
}

func CacheKeyPrefixReadOnly(selector CacheSelector, keyPrefix Value) DisposableTokenScope {
	return Permissions{
		Permissions: []Permission{DisposableTokenCachePermission{Cache: selector, Item: CacheItemKeyPrefix{KeyPrefix: keyPrefix.asBytes()}, Role: ReadOnly}},
	}
}

func CacheKeyWriteOnly(selector CacheSelector, key Value) DisposableTokenScope {
	return Permissions{
		Permissions: []Permission{DisposableTokenCachePermission{Cache: selector, Item: CacheItemKey{Key: key.asBytes()}, Role: WriteOnly}},
	}
}

func CacheKeyPrefixWriteOnly(selector CacheSelector, keyPrefix Value) DisposableTokenScope {
	return Permissions{
		Permissions: []Permission{DisposableTokenCachePermission{Cache: selector, Item: CacheItemKeyPrefix{KeyPrefix: keyPrefix.asBytes()}, Role: WriteOnly}},
	}
}

func TopicNamePrefixPublishSubscribe(selector CacheSelector, topicNamePrefix string) DisposableTokenScope {
	return Permissions{
		Permissions: []Permission{TopicPermission{Cache: selector, Topic: TopicNamePrefix{NamePrefix: topicNamePrefix}, Role: PublishSubscribe}},
	}
}

func TopicNamePrefixPublishOnly(selector CacheSelector, topicNamePrefix string) DisposableTokenScope {
	return Permissions{
		Permissions: []Permission{TopicPermission{Cache: selector, Topic: TopicNamePrefix{NamePrefix: topicNamePrefix}, Role: PublishOnly}},
	}
}

func TopicNamePrefixSubscribeOnly(selector CacheSelector, topicNamePrefix string) DisposableTokenScope {
	return Permissions{
		Permissions: []Permission{TopicPermission{Cache: selector, Topic: TopicNamePrefix{NamePrefix: topicNamePrefix}, Role: SubscribeOnly}},
	}
}
