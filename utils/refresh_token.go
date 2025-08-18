package utils

import (
	"crypto/rand"
	"encoding/base32"
)

// GenerateRefreshToken generates a random, almost-unique string of length 10.
func GenerateRefreshToken() string {
	b := make([]byte, 8) // 8 bytes = 64 bits
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // handle error as needed
	}
	// base32 encode and trim to 10 chars, remove padding
	return base32.StdEncoding.EncodeToString(b)[:10]
}
