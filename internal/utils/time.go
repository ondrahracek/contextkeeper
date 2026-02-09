// Package utils provides utility functions for the contextkeeper application.
// It includes formatting helpers for output, tag parsing/validation, UUID generation,
// and time formatting utilities.
package utils

import "time"

// FormatTime formats a time.Time value using the specified format string.
//
// Parameters:
//   - t: The time value to format
//   - format: A format string compatible with time.Time.Format()
//
// Returns:
//
//	A formatted string representation of the time
func FormatTime(t time.Time, format string) string {
	return t.Format(format)
}

// ParseTime parses a string into a time.Time value using the specified format.
//
// Parameters:
//   - s: The string to parse
//   - format: A format string compatible with time.Parse()
//
// Returns:
//   - The parsed time value
//   - An error if parsing fails
func ParseTime(s string, format string) (time.Time, error) {
	return time.Parse(format, s)
}
