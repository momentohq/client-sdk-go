package utility

import (
	"fmt"
	"strings"

	internalRequests "github.com/momentohq/client-sdk-go/internal/requests"
)

func IsCacheNameValid(cacheName string) bool {
	return len(strings.TrimSpace(cacheName)) != 0
}

func ConvertEcacheResult(resultRequest internalRequests.ConvertEcacheResultRequest) error {
	return fmt.Errorf("CacheService returned an unexpected result: %v for operation: %s with message: %s", resultRequest.ECacheResult, resultRequest.OpName, resultRequest.Message)
}
