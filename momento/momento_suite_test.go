package momento_test

import (
	"fmt"
	"testing"

	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
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
				return "", fmt.Errorf("expected MomentoError, but got %T", err)
			}
		}, Equal(code),
	)
}

func HaveSetLength(length int) types.GomegaMatcher {
	return WithTransform(
		func(fetchResp responses.SetFetchResponse) (int, error) {
			switch rtype := fetchResp.(type) {
			case *responses.SetFetchHit:
				return len(rtype.ValueString()), nil
			default:
				return 0, fmt.Errorf("expected set fetch hit but got %T", fetchResp)
			}
		}, Equal(length),
	)
}

func HaveListLength(length int) types.GomegaMatcher {
	return WithTransform(
		func(fetchResp responses.ListFetchResponse) (int, error) {
			switch rtype := fetchResp.(type) {
			case *responses.ListFetchHit:
				return len(rtype.ValueList()), nil
			default:
				return 0, fmt.Errorf("expected list fetch hit but got %T", fetchResp)
			}
		}, Equal(length),
	)
}

func HaveSortedSetElements(expected []responses.SortedSetBytesElement) types.GomegaMatcher {
	return WithTransform(
		func(fetchResp responses.SortedSetFetchResponse) ([]responses.SortedSetBytesElement, error) {
			switch rtype := fetchResp.(type) {
			case *responses.SortedSetFetchHit:
				return rtype.ValueBytesElements(), nil
			default:
				return nil, fmt.Errorf("expected SortedSetFetchHit, but got %T", fetchResp)
			}
		}, Equal(expected),
	)
}

func HaveSortedSetStringElements(expected []responses.SortedSetStringElement) types.GomegaMatcher {
	return WithTransform(
		func(fetchResp responses.SortedSetFetchResponse) ([]responses.SortedSetStringElement, error) {
			switch rtype := fetchResp.(type) {
			case *responses.SortedSetFetchHit:
				return rtype.ValueStringElements(), nil
			default:
				return nil, fmt.Errorf("expected SortedSetFetchHit, but got %T", fetchResp)
			}
		}, Equal(expected),
	)
}
