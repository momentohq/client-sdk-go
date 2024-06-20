package momento

import (
	"fmt"
)

// Momento Error codes
const (
	// InvalidArgumentError occurs when an invalid argument is passed to Momento client.
	InvalidArgumentError = "InvalidArgumentError"
	// InternalServerError occurs when an unexpected error is encountered while trying to fulfill the request.
	InternalServerError = "InternalServerError"
	// ClientSdkError occurs when a client side error happens.
	ClientSdkError = "ClientSdkError"
	// BadRequestError occurs when a request was invalid.
	BadRequestError = "BadRequestError"
	// CanceledError occurs when a request was cancelled by the server.
	CanceledError = "CanceledError"
	// TimeoutError occurs when an operation did not complete in time.
	TimeoutError = "TimeoutError"
	// PermissionError occurs when there are insufficient permissions to perform operation.
	PermissionError = "PermissionError"
	// AuthenticationError occurs when invalid authentication credentials to connect to cache service are provided.
	AuthenticationError = "AuthenticationError"
	// LimitExceededError occurs when request rate, bandwidth, or object size exceeded the limits for the account.
	LimitExceededError = "LimitExceededError"
	// CacheNotFoundError occurs when a cache with specified name doesn't exist.
	CacheNotFoundError = "CacheNotFoundError"
	// StoreNotFoundError occurs when a store with specified name doesn't exist.
	StoreNotFoundError = "StoreNotFoundError"
	// ItemNotFoundError occurs when an item with specified key doesn't exist.
	ItemNotFoundError = "ItemNotFoundError"
	// AlreadyExistsError occurs when a cache with specified name already exists.
	AlreadyExistsError = "AlreadyExistsError"
	// UnknownServiceError occurs when an unknown error has occurred.
	UnknownServiceError = "UnknownServiceError"
	// ServerUnavailableError occurs when the server was unable to handle the request.
	ServerUnavailableError = "ServerUnavailableError"
	// FailedPreconditionError occurs when the system is not in a state required for the operation's execution.
	FailedPreconditionError = "FailedPreconditionError"
	// ConnectionError occurs when there is an error connecting to Momento servers.
	ConnectionError = "ConnectionError"
)

type MomentoError interface {
	// Satisfy the generic error interface.
	error
	// Code Returns Momento Error codes.
	Code() string
	// Message Returns the error details message.
	Message() string
	// OriginalErr Returns the original error if one was set.  Nil is returned if not set.
	OriginalErr() error
}

type momentoError struct {
	code        string
	message     string
	originalErr error
}

// Code Returns Momento Error codes.
func (err momentoError) Code() string {
	return err.code
}

// Message Returns the error details message.
func (err momentoError) Message() string {
	return err.message
}

// OriginalErr Returns the original error if one was set.  Nil is returned if not set.
func (err momentoError) OriginalErr() error {
	if err.originalErr != nil {
		return err.originalErr
	}
	return nil
}

// Satisfies the generic error interface.
// Returns the error details message with code, message, original error if there is any.
func (err momentoError) Error() string {
	if err.originalErr != nil {
		return fmt.Sprintf("%s: %s\n%s", err.code, err.message, err.originalErr)
	}
	return fmt.Sprintf("%s: %s", err.code, err.message)
}

// NewMomentoError returns an initialized MomentoError wrapper.
func NewMomentoError(code string, message string, originalErr error) MomentoError {
	return &momentoError{
		code,
		message,
		originalErr,
	}
}
