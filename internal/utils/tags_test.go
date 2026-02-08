package utils

import (
	"strings"
	"testing"
)

// TestParseTags tests the ParseTags function with various input formats.
func TestParseTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "comma-separated tags",
			input:    "tag1,tag2",
			expected: []string{"tag1", "tag2"},
		},
		{
			name:     "comma and space separated tags",
			input:    "tag1, tag2, tag3",
			expected: []string{"tag1", "tag2", "tag3"},
		},
		{
			name:     "space-separated tags",
			input:    "tag1 tag2 tag3",
			expected: []string{"tag1", "tag2", "tag3"},
		},
		{
			name:     "single tag",
			input:    "single",
			expected: []string{"single"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: nil,
		},
		{
			name:     "duplicates removed",
			input:    "tag1, tag2, tag1, tag3",
			expected: []string{"tag1", "tag2", "tag3"},
		},
		{
			name:     "mixed separators",
			input:    "tag1,tag2 tag3, tag4",
			expected: []string{"tag1", "tag2", "tag3", "tag4"},
		},
		{
			name:     "tags with underscores",
			input:    "tag_1, tag-2",
			expected: []string{"tag_1", "tag-2"},
		},
		{
			name:     "tags with numbers",
			input:    "tag1, tag2tag3",
			expected: []string{"tag1", "tag2tag3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseTags(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("ParseTags(%q): got %d tags, want %d", tt.input, len(result), len(tt.expected))
				return
			}
			for i, tag := range result {
				if tag != tt.expected[i] {
					t.Errorf("ParseTags(%q)[%d]: got %q, want %q", tt.input, i, tag, tt.expected[i])
				}
			}
		})
	}
}

// TestValidateTags tests the ValidateTags function with valid and invalid tag combinations.
func TestValidateTags(t *testing.T) {
	tests := []struct {
		name        string
		tags        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid single tag",
			tags:        []string{"tag1"},
			expectError: false,
		},
		{
			name:        "valid multiple tags",
			tags:        []string{"tag1", "tag2", "tag_3"},
			expectError: false,
		},
		{
			name:        "valid tags with hyphens",
			tags:        []string{"my-tag", "another-tag"},
			expectError: false,
		},
		{
			name:        "valid tags with numbers",
			tags:        []string{"tag123", "tag-456"},
			expectError: false,
		},
		{
			name:        "empty slice is valid",
			tags:        []string{},
			expectError: false,
		},
		{
			name:        "nil slice is valid",
			tags:        nil,
			expectError: false,
		},
		{
			name:        "empty tag is invalid",
			tags:        []string{""},
			expectError: true,
			errorMsg:    "tag cannot be empty",
		},
		{
			name:        "tag too long",
			tags:        []string{"a" + strings.Repeat("b", maxTagLength)},
			expectError: true,
			errorMsg:    "tag too long: maximum 50 characters",
		},
		{
			name:        "tag with special characters",
			tags:        []string{"tag@here"},
			expectError: true,
			errorMsg:    "invalid tag format: must contain only alphanumeric characters, underscores, and hyphens",
		},
		{
			name:        "tag with spaces",
			tags:        []string{"tag here"},
			expectError: true,
			errorMsg:    "invalid tag format: must contain only alphanumeric characters, underscores, and hyphens",
		},
		{
			name:        "one invalid tag invalidates all",
			tags:        []string{"valid-tag", "invalid@tag"},
			expectError: true,
			errorMsg:    "invalid tag format: must contain only alphanumeric characters, underscores, and hyphens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTags(tt.tags)
			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateTags(%v): expected error, got nil", tt.tags)
				} else if err.Error() != tt.errorMsg {
					t.Errorf("ValidateTags(%v): got error %q, want %q", tt.tags, err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateTags(%v): unexpected error: %v", tt.tags, err)
				}
			}
		})
	}
}
