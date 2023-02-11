package internal_serve

import (
	"strings"
)

func isHeaderAllowed(name string) bool {
	lowerName := strings.ToLower(name)
	if lowerName == "content-length" ||
		lowerName == "connection" ||
		lowerName == "content-encoding" ||
		lowerName == "report-to" ||
		lowerName == "via" ||
		lowerName == "age" ||
		lowerName == ":authority" ||
		lowerName == "expect-ct" ||
		lowerName == "date" ||
		lowerName == "etag" ||
		lowerName == "last-modified" ||
		lowerName == "strict-transport-security" {
		return false
	}
	return true
}
