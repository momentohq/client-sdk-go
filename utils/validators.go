package utils

import (
	"fmt"
	"strings"

	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

func ValidateDisposableTokenExpiry(in ExpiresIn) momentoerrors.MomentoSvcErr {
	if !in.DoesExpire() {
		return momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"disposable tokens must have an expiry",
			nil,
		)
	}
	return nil
}

func ValidateApiKeyExpiry(in ExpiresIn) momentoerrors.MomentoSvcErr {
	if !in.DoesExpire() {
		return nil
	}
	if in.Seconds() < 0 {
		return momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"API key expiry must be a positive number of seconds",
			nil,
		)
	}
	return nil
}

func ValidateName(name string, label string) error {
	if len(strings.TrimSpace(name)) < 1 {
		errStr := fmt.Sprintf("%v cannot be empty or blank", label)
		return momentoerrors.NewMomentoSvcErr(momentoerrors.InvalidArgumentError, errStr, nil)
	}
	return nil
}
