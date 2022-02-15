package momento

import (
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

const (
	InvalidArgumentError = "InvalidArgumentError"
	InternalServerError  = "InternalServerError"
	ClientSdkError       = "ClientSdkError"
	FailedPrecondition   = "FailedPrecondition"
	Canceled             = "Canceled"
	DeadlineExceeded     = "DeadlineExceeded"
	PermissionDenied     = "PermissionDenied"
	Unauthenticated      = "Unauthenticated"
	ResourceExhausted    = "ResourceExhausted"
	NotFound             = "NotFound"
	AlreadyExists        = "AlreadyExists"
	Unavailable          = "Unavailable"
)

type MomentoError interface {
	momentoerrors.MomentoSvcErr
}

type momentoError struct {
	err momentoerrors.MomentoSvcErr
}

func newMomentoError(err momentoerrors.MomentoSvcErr) *momentoError {
	return &momentoError{
		err,
	}
}

func (momentoerror momentoError) Code() string {
	return momentoerror.err.Code()
}

func (momentoerror momentoError) Message() string {
	return momentoerror.err.Message()
}

func (momentoerror momentoError) Error() string {
	return fmt.Sprintf("%s: %s", momentoerror.err.Code(), momentoerror.err.Message())
}

func NewMomentoError(err momentoerrors.MomentoSvcErr) MomentoError {
	return newMomentoError(err)
}
