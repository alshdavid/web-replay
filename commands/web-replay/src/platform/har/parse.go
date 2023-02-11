package har

import (
	"encoding/json"
	"os"
)

func MustParseFile(targetFilepath string) []Entry {
	b, _ := os.ReadFile(targetFilepath)
	result := Model{}
	json.Unmarshal(b, &result)
	return result.Log.Entries
}
