package momento_test

import (
	"fmt"
	"testing"

	"github.com/momentohq/client-sdk-go/momento"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func TestMomento(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Momento Suite")
}

func HaveMomentoErrorCode(code string) types.GomegaMatcher {
	return WithTransform(
		func(err error) (string, error) {
			switch mErr := err.(type) {
			case momento.MomentoError:
				return mErr.Code(), nil
			default:
				return "", fmt.Errorf("Expected MomentoError, but got %T", err)
			}
		}, Equal(code),
	)
}
