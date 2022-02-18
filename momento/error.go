package momento

import (
	"fmt"
)

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
	error
	Code() string
	Message() string
	OriginalErr() error
}

type momentoError struct {
	code        string
	message     string
	originalErr error
}

func (err momentoError) Code() string {
	return err.code
}

func (err momentoError) Message() string {
	return err.message
}

func (err momentoError) OriginalErr() error {
	if err.originalErr != nil {
		return err.originalErr
	}
	return nil
}

func (err momentoError) Error() string {
	if err.originalErr != nil {
		return fmt.Sprintf("%s: %s\n%s", err.code, err.message, err.originalErr)
	}
	return fmt.Sprintf("%s: %s", err.code, err.message)
}

func NewMomentoError(code string, message string, originalErr error) MomentoError {
	return &momentoError{
		code,
		message,
		originalErr,
	}
}
