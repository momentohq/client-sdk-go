package utils

import (
	"fmt"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
)

func ValidateDisposableTokenExpiry(in ExpiresIn) momentoerrors.MomentoSvcErr {
	if !in.DoesExpire() {
		return momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			fmt.Sprintf("disposable tokens must have an expiry"),
			nil,
		)
	}
	if in.Seconds() > 60*60 {
		return momentoerrors.NewMomentoSvcErr(
			momentoerrors.InvalidArgumentError,
			fmt.Sprintf("disposable tokens must expire within 1 hour"),
			nil,
		)
	}
	return nil
}
