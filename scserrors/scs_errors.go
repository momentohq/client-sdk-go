package errors

import "errors"

func InvalidInputError(errMessage string) error {
	return errors.New("InvalidInputError: " + errMessage)
}
