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

func HaveSetLength(length int) types.GomegaMatcher {
	return WithTransform(
		func(fetchResp momento.SetFetchResponse) (int, error) {
			switch rtype := fetchResp.(type) {
			case *momento.SetFetchHit:
				return len(rtype.ValueString()), nil
			default:
				return 0, fmt.Errorf("expected set fetch hit but got %T", fetchResp)
			}
		}, Equal(length),
	)
}
