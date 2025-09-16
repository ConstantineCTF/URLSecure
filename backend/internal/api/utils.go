package api

import (
	"crypto/rand"
	"encoding/base64"
)

// generateCode returns a URL-safe random string of length n.
// Used to generate new unique short codes for URLs
func generateCode(n int) string {
	// Create a buffer of n random bytes
	b := make([]byte, n)

	// Fill the buffer with random data, panic on failure (should rarely happen)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	// Encode bytes to URL-safe base64, truncate to requested length
	return base64.URLEncoding.EncodeToString(b)[:n]
}
