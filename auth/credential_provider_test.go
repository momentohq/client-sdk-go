package auth_test

import (
	"errors"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CredentialProvider", func() {
	It(`errors on a invalid auth token`, func() {
		badCredentialProvider, err := auth.NewStringMomentoTokenProvider("Invalid token")

		Expect(badCredentialProvider).To(BeNil())
		Expect(err).NotTo(BeNil())

		var momentoErr momentoerrors.MomentoSvcErr
		if errors.As(err, &momentoErr) {
			Expect(momentoErr.Code()).To(Equal(momentoerrors.InvalidArgumentError))
		}
	})
})
