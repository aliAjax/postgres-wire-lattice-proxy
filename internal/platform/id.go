package platform

import (
	"crypto/rand"
	"encoding/hex"
)

func ID() string {
	b := make([]byte, 16)
	if _, e := rand.Read(b); e != nil {
		return "fallback-id"
	}
	return hex.EncodeToString(b)
}
