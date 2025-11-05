package service

import (
	"crypto/rand"
	"math/big"
	"time"
)

// GenerateID creates a ULID (Universally Unique Lexicographically Sortable Identifier)
// Format: 26 characters, Crockford Base32 encoded
// First 10 chars: timestamp (milliseconds since Unix epoch)
// Last 16 chars: random data
// Example: 01ARZ3NDEKTSV4RRFFQ69G5FAV
func GenerateID() string {
	// Crockford Base32 alphabet (no I, L, O, U to avoid confusion)
	alphabet := "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

	// Get current timestamp in milliseconds since Unix epoch
	timestampMs := time.Now().UnixMilli()

	// Convert timestamp to 10-character base32 string (48 bits)
	timestampPart := make([]byte, 10)
	tempVal := timestampMs
	for i := 9; i >= 0; i-- {
		timestampPart[i] = alphabet[tempVal%32]
		tempVal = tempVal / 32
	}

	// Generate 16 random characters (80 bits)
	randomPart := make([]byte, 16)
	for i := 0; i < 16; i++ {
		// Generate random number between 0-31
		n, err := rand.Int(rand.Reader, big.NewInt(32))
		if err != nil {
			// Fallback to timestamp-based randomness if crypto rand fails
			n = big.NewInt(time.Now().UnixNano() % 32)
		}
		randomPart[i] = alphabet[n.Int64()]
	}

	return string(timestampPart) + string(randomPart)
}
