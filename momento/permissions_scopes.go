package momento

func CacheReadWrite(selector CacheSelector) Permissions {
	return Permissions{Permissions: []Permission{CachePermission{Cache: selector, Role: ReadWrite}}}
}

func CacheReadOnly(selector CacheSelector) Permissions {
	return Permissions{Permissions: []Permission{CachePermission{Cache: selector, Role: ReadOnly}}}
}

func CacheWriteOnly(selector CacheSelector) Permissions {
	return Permissions{Permissions: []Permission{CachePermission{Cache: selector, Role: WriteOnly}}}
}

func TopicSubscribeOnly(cacheSelector CacheSelector, topicSelector TopicSelector) Permissions {
	return Permissions{Permissions: []Permission{TopicPermission{Cache: cacheSelector, Topic: topicSelector, Role: SubscribeOnly}}}
}

func TopicPublishSubscribe(cacheSelector CacheSelector, topicSelector TopicSelector) Permissions {
	return Permissions{Permissions: []Permission{TopicPermission{Cache: cacheSelector, Topic: topicSelector, Role: PublishSubscribe}}}
}

func TopicPublishOnly(cacheSelector CacheSelector, topicSelector TopicSelector) Permissions {
	return Permissions{Permissions: []Permission{TopicPermission{Cache: cacheSelector, Topic: topicSelector, Role: PublishOnly}}}
}
