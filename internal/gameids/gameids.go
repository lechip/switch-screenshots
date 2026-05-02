package gameids

import (
	"encoding/json"
	"strings"

	"github.com/lechip/switch-screenshots/assets"
)

var ids map[string]string

func init() {
	if err := json.Unmarshal(assets.GameIDsJSON, &ids); err != nil {
		panic("gameids: " + err.Error())
	}
}

// Lookup returns the game title for a Switch game ID.
// id is matched case-insensitively; ok is false for unknown IDs.
func Lookup(id string) (title string, ok bool) {
	title, ok = ids[strings.ToUpper(id)]
	return
}
