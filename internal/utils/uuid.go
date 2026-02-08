// Package utils provides utility functions for the contextkeeper application.
// It includes formatting helpers for output, tag parsing/validation, UUID generation,
// and time formatting utilities.
package utils

import (
	"crypto/rand"
	"fmt"
)

// UUID constants.
const (
	// uuidLength is the standard length of a UUID string (36 characters including dashes)
	uuidLength = 36
	// uuidVersion4Format indicates the UUID follows the version 4 random format
	uuidVersion4Format = 4
)

// GenerateUUID generates a random UUID version 4 using cryptographically secure
// random number generation (crypto/rand).
//
// Returns:
//   A string representation of the UUID in the format "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
func GenerateUUID() string {
	// Generate 16 random bytes (128 bits) for UUID v4
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to a simple random string if crypto/rand fails
		// This should rarely happen in practice
		return fmt.Sprintf("rand-%d", b[0])
	}

	// Set UUID version 4 bits (bits 4-7 of byte 6)
	b[6] = (b[6] & 0x0f) | (uuidVersion4Format << 4)
	// Set UUID variant to RFC 4122 (bits 6-7 of byte 8)
	b[8] = (b[8] & 0x3f) | 0x80

	// Format as hex groups: 8-4-4-4-12
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
