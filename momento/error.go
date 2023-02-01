package momento

import (
	"fmt"
)

// Momento Error codes
const (
	InvalidArgumentError = "InvalidArgumentError"
	InternalServerError  = "InternalServerError"
	ClientSdkError       = "ClientSdkError"
	BadRequestError      = "BadRequestError"
	CanceledError        = "CanceledError"
	TimeoutError         = "TimeoutError"
	PermissionError      = "PermissionError"
	AuthenticationError  = "AuthenticationError"
	LimitExceededError   = "LimitExceededError"
	NotFoundError        = "NotFoundError"
	AlreadyExistsError   = "AlreadyExistsError"
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
