package helper

import (
	"crypto/sha1"
	"encoding/base64"
	"sort"
	"strings"
)

func CreateConversationID(userIDs ...string) string {
	sort.Strings(userIDs)
	s := strings.Join(userIDs, ";")
	h := sha1.New()

	h.Write([]byte(s))

	bs := h.Sum(nil)

	convID := base64.StdEncoding.EncodeToString(bs)

	return convID
}
