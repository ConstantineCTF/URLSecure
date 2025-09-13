package api

import (
	"crypto/rand"
	"encoding/base64"
)

// generateCode returns a URL-safe random string of length n.
func generateCode(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)[:n]
}
