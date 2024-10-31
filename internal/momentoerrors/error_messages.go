package momentoerrors

// LimitExceededError message wrappers to indicate the type of limit exceeded
const (
	TopicSubscriptionsLimitExceeded = "Topic subscriptions limit exceeded for this account"
	OperationsRateLimitExceeded     = "Request rate limit exceeded for this account"
	ThroughputRateLimitExceeded     = "Bandwidth limit exceeded for this account"
	RequestSizeLimitExceeded        = "Request size limit exceeded for this account"
	ItemSizeLimitExceeded           = "Item size limit exceeded for this account"
	ElementSizeLimitExceeded        = "Element size limit exceeded for this account"
	UnknownLimitExceeded            = "Limit exceeded for this account"
)

const (
	InvalidArgumentMessageWrapper     = "Invalid argument passed to Momento client"
	BadRequestMessageWrapper          = "The request was invalid; please contact us at support@momentohq.com"
	FailedPreconditionMessageWrapper  = "System is not in a state required for the operation's execution"
	CanceledMessageWrapper            = "The request was cancelled by the server; please contact us at support@momentohq.com"
	TimeoutMessageWrapper             = "The client's configured timeout was exceeded; you may need to use a Configuration with more lenient timeouts"
	PermissionMessageWrapper          = "Insufficient permissions to perform operation"
	AuthenticationMessageWrapper      = "Invalid authentication credentials to connect to Momento service"
	CacheNotFoundMessageWrapper       = "A cache with the specified name does not exist.  To resolve this error, make sure you have created the cache before attempting to use it"
	StoreNotFoundMessageWrapper       = "A store with the specified name does not exist.  To resolve this error, make sure you have created the store before attempting to use it"
	ItemNotFoundMessageWrapper        = "An item with the specified key does not exist"
	CacheAlreadyExistsMessageWrapper  = "A cache with the specified name already exists.  To resolve this error, either delete the existing cache and make a new one, or use a different name"
	StoreAlreadyExistsMessageWrapper  = "A store with the specified name already exists.  To resolve this error, either delete the existing store and make a new one, or use a different name"
	UnknownServiceErrorMessageWrapper = "Service returned an unknown response; please contact us at support@momentohq.com"
	InternalServerErrorMessageWrapper = "Unexpected error encountered while trying to fulfill the request; please contact us at support@momentohq.com"
	ServerUnavailableMessageWrapper   = "The server was unable to handle the request; consider retrying.  If the error persists, please contact us at support@momentohq.com"
)
