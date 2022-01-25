package utility

import (
	"fmt"
	"strings"

	pb "github.com/momentohq/client-sdk-go/protos"
)

func IsCacheNameValid(cacheName string) bool {
	return len(strings.TrimSpace(cacheName)) != 0
}

func ConvertEcacheResult(eCacheResult pb.ECacheResult, message string, opName string) error {
	return fmt.Errorf("CacheService returned an unexpected result: %v for operation: %s with message: %s", eCacheResult, opName, message)
}
