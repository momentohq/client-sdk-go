package utility

import (
	"strings"
)

func IsCacheNameValid(cacheName string) bool {
	return len(strings.TrimSpace(cacheName)) != 0
}
