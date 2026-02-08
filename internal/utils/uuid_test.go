package utils

import (
	"strings"
	"testing"
)

// TestGenerateUUID tests the GenerateUUID function for format, uniqueness, and validity.
func TestGenerateUUID(t *testing.T) {
	tests := []struct {
		name   string
		action func() string
		check  func(t *testing.T, result string)
	}{
		{
			name:   "non-empty result",
			action: GenerateUUID,
			check: func(t *testing.T, result string) {
				if result == "" {
					t.Error("GenerateUUID should not return empty string")
				}
			},
		},
		{
			name:   "correct length",
			action: GenerateUUID,
			check: func(t *testing.T, result string) {
				if len(result) != uuidLength {
					t.Errorf("UUID length: got %d, want %d", len(result), uuidLength)
				}
			},
		},
		{
			name:   "correct format (8-4-4-4-12)",
			action: GenerateUUID,
			check: func(t *testing.T, result string) {
				parts := strings.Split(result, "-")
				if len(parts) != 5 {
					t.Errorf("UUID format: got %d parts, want 5", len(parts))
					return
				}
				// Check expected lengths of each part: 4, 2, 2, 2, 6 (in hex = 8, 4, 4, 4, 12 chars)
				expectedLengths := []int{8, 4, 4, 4, 12}
				for i, part := range parts {
					if len(part) != expectedLengths[i] {
						t.Errorf("UUID part %d: got length %d, want %d", i, len(part), expectedLengths[i])
					}
				}
			},
		},
		{
			name:   "uniqueness - different UUIDs generated",
			action: GenerateUUID,
			check: func(t *testing.T, result string) {
				// Note: This test has a theoretical chance of collision, but it's astronomically low
				result2 := GenerateUUID()
				if result == result2 {
					t.Error("GenerateUUID should return unique values")
				}
			},
		},
		{
			name:   "lowercase hex format",
			action: GenerateUUID,
			check: func(t *testing.T, result string) {
				// UUID should only contain lowercase hex characters and dashes
				for _, ch := range result {
					if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || ch == '-') {
						t.Errorf("UUID contains invalid character: %c", ch)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.action()
			tt.check(t, result)
		})
	}
}
