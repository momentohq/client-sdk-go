package momentoerrors

import (
	"fmt"
)

type momentoBaseError struct {
	code    string
	message string
}

func newMomentoBaseError(code string, message string) *momentoBaseError {
	return &momentoBaseError{
		code,
		message,
	}
}

func (err momentoBaseError) Code() string {
	return err.code
}

func (err momentoBaseError) Message() string {
	return err.message
}

func (err momentoBaseError) Error() string {
	return fmt.Sprintf("%s: %s", err.code, err.message)
}
