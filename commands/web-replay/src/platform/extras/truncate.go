package extras

import "unicode/utf8"

func TruncateString(input string, max int) string {
	if len(input) < max {
		return input
	}

	if utf8.ValidString(input[:max]) {
		return input[:max]
	}
	return input[:max+1]
}
