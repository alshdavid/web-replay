package internal_serve

import (
	"fmt"
	"strings"
)

func transformDomains(sourceText string, domainMap map[string]string) string {
	result := sourceText

	for src, dest := range domainMap {
		result = strings.ReplaceAll(result, fmt.Sprintf("http://%s", src), dest)
		result = strings.ReplaceAll(result, fmt.Sprintf("https://%s", src), dest)
		result = strings.ReplaceAll(result, fmt.Sprintf("//%s", src), dest)
	}

	return result
}
