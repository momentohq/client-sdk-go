package momentoerrors

import "fmt"

type momentoSvcError struct {
	code    string
	message string
}

func newMomentoSvcErr(code string, message string) *momentoSvcError {
	return &momentoSvcError{
		code,
		message,
	}
}

func (err momentoSvcError) Code() string {
	return err.code
}

func (err momentoSvcError) Message() string {
	return err.message
}

func (err momentoSvcError) Error() string {
	return fmt.Sprintf("%s: %s", err.code, err.message)
}
