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
	if in.Seconds() > 60*60 {
		return momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			"disposable tokens must expire within 1 hour",
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
