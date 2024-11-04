package momentoerrors

import "fmt"

type momentoSvcError struct {
	code           string
	message        string
	originalErr    error
	messageWrapper string
}

func newMomentoSvcErr(code string, message string, originalErr error, messageWrapper string) *momentoSvcError {
	return &momentoSvcError{
		code,
		message,
		originalErr,
		messageWrapper,
	}
}

func (err momentoSvcError) Code() string {
	return err.code
}

func (err momentoSvcError) Message() string {
	if err.messageWrapper == "" {
		return err.message
	}
	return err.messageWrapper + ": " + err.message
}

func (err momentoSvcError) OriginalErr() error {
	if err.originalErr != nil {
		return err.originalErr
	}
	return nil
}

func (err momentoSvcError) Error() string {
	if err.originalErr != nil {
		return fmt.Sprintf("%s: %s\n%s", err.code, err.message, err.originalErr.Error())
	}
	return fmt.Sprintf("%s: %s", err.code, err.message)
}
