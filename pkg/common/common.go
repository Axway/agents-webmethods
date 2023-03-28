package common

import "fmt"

const (
	AppID        = "appID"
	AttrAPIID    = "sampleApiId"
	AttrChecksum = "checksum"
	AttrAppID    = "webmethodsApplicationId"
)

// FormatAPICacheKey ensure consistent naming of the cache key for an API.
func FormatAPICacheKey(apiID, stageName string) string {
	return fmt.Sprintf("%s-%s", apiID, stageName)
}
