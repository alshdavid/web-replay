package fsExtras

import (
	"os"
)

func Exists(targetPath string) bool {
	result, err := os.Stat(targetPath)
	if result == nil || (err != nil && os.IsNotExist(err)) {
		return false
	}
	return true
}
