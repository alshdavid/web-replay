package extras

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func GetMD5Hash(text ...string) string {
	hash := md5.Sum([]byte(strings.Join(text, "")))
	return hex.EncodeToString(hash[:])
}
