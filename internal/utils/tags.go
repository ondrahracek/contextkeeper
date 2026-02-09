// Package utils provides utility functions for the contextkeeper application.
// It includes formatting helpers for output, tag parsing/validation, UUID generation,
// and time formatting utilities.
package utils

import (
	"errors"
	"regexp"
	"strings"
)

// Tag validation constants.
const (
	// maxTagLength is the maximum allowed length for a tag
	maxTagLength = 50
)

// tagRegex defines the valid characters for tags: alphanumeric, underscore, and hyphen
var tagRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ParseTags parses a space and/or comma-separated string into a slice of tags.
// The function normalizes the input by replacing commas with spaces and splitting,
// then filters out empty tags and removes duplicates while preserving order.
//
// Parameters:
//   - tagStr: A string containing tags separated by spaces, commas, or both
//
// Returns:
//
//	A slice of unique, trimmed tags in their original order
func ParseTags(tagStr string) []string {
	if strings.TrimSpace(tagStr) == "" {
		return nil
	}

	// Normalize separators: replace commas with spaces, then split on whitespace
	normalized := strings.ReplaceAll(tagStr, ",", " ")
	tags := strings.Fields(normalized)

	// Filter out empty tags and duplicates, preserving first occurrence order
	seen := make(map[string]bool)
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" && !seen[tag] {
			seen[tag] = true
			result = append(result, tag)
		}
	}

	return result
}

// ValidateTags validates a slice of tags against the tag format rules.
//
// Validation rules:
//   - Empty slice is considered valid
//   - Each tag must be non-empty
//   - Each tag must not exceed maxTagLength (50) characters
//   - Each tag must match the pattern: only alphanumeric chars, underscores, and hyphens
//
// Parameters:
//   - tags: A slice of tags to validate
//
// Returns:
//   - nil if all tags are valid
//   - An error describing the first validation failure
func ValidateTags(tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	for _, tag := range tags {
		if len(tag) == 0 {
			return errors.New("tag cannot be empty")
		}
		if len(tag) > maxTagLength {
			return errors.New("tag too long: maximum 50 characters")
		}
		if !tagRegex.MatchString(tag) {
			return errors.New("invalid tag format: must contain only alphanumeric characters, underscores, and hyphens")
		}
	}

	return nil
}
