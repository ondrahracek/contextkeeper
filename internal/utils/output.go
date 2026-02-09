// Package utils provides utility functions for the contextkeeper application.
// It includes formatting helpers for output, tag parsing/validation, UUID generation,
// and time formatting utilities.
package utils

import (
	"fmt"
	"strings"

	"github.com/ondrahracek/contextkeeper/internal/models"
)

// Color code constants for terminal output styling.
// These ANSI escape codes are used to add visual indicators to formatted output.
const (
	// colorReset resets all text attributes (color, bold, etc.)
	colorReset = "\033[0m"
	// colorGreen is used for positive indicators like completed items
	colorGreen = "\033[32m"
	// colorYellow is used for tag indicators
	colorYellow = "\033[33m"
	// colorCyan is used for project name indicators
	colorCyan = "\033[36m"
	// colorRed is used for error indicators
	colorRed = "\033[31m"
	// colorBold makes text appear bold
	colorBold = "\033[1m"
	// colorDim makes text appear dimmed
	colorDim = "\033[2m"
)

// Truncation limits for output formatting.
const (
	// maxContentLength is the maximum length for item content in list view
	maxContentLength = 50
	// contentTruncationIndicator is appended to truncated content
	contentTruncationIndicator = "..."
)

// Column width constants for table formatting.
const (
	// colWidthID is the fixed width for the ID column in table output
	colWidthID = 8
	// colWidthContent is the fixed width for the content column in table output
	colWidthContent = 60
	// colWidthProject is the fixed width for the project column in table output
	colWidthProject = 15
	// colWidthDate is the fixed width for the date column in table output
	colWidthDate = 16
)

// FormatItemList formats a list of ContextItem for display.
// Items can be optionally filtered to hide completed items.
// The output includes ID, status indicators, content, project info, tags, and creation date.
//
// Parameters:
//   - items: Slice of ContextItem to format
//   - showCompleted: If true, includes completed items in output
//
// Returns:
//
//	A formatted string representation of the items, or "No items found." if empty
func FormatItemList(items []models.ContextItem, showCompleted bool) string {
	if len(items) == 0 {
		return "No items found."
	}

	var sb strings.Builder
	for _, item := range items {
		if item.CompletedAt != nil && !showCompleted {
			continue
		}

		// Show first 6 characters of ID
		idDisplay := item.ID[:6]

		status := "[ ]"
		if item.CompletedAt != nil {
			status = fmt.Sprintf("%s[x]%s", colorGreen, colorReset)
		}

		projectInfo := ""
		if item.Project != "" {
			projectInfo = fmt.Sprintf("%s@%s%s", colorCyan, item.Project, colorReset)
		}

		tagsInfo := ""
		if len(item.Tags) > 0 {
			tagsInfo = fmt.Sprintf(" %s[%s]%s", colorYellow, strings.Join(item.Tags, ", "), colorReset)
		}

		createdAt := item.CreatedAt.Format("2006-01-02 15:04")
		truncatedContent := truncateString(item.Content, maxContentLength)

		// Format the output line.
		// Note: We use plain [ID] without bold ANSI codes to ensure cross-platform
		// compatibility, especially for Windows terminals that may not support
		// or have virtual terminal sequences enabled.
		fmt.Fprintf(&sb, "%s [%s] %s %s %s %s\n",
			status,
			idDisplay,
			truncatedContent,
			projectInfo,
			tagsInfo,
			createdAt,
		)
	}

	return sb.String()
}

// FormatTable formats items in a table-like format with aligned columns.
// Each row displays ID, content, project, and creation date.
//
// Parameters:
//   - items: Slice of ContextItem to format
//
// Returns:
//
//	A formatted table string, or "No items to display." if empty
func FormatTable(items []models.ContextItem) string {
	if len(items) == 0 {
		return "No items to display."
	}

	// Calculate column widths (using predefined constants)
	idWidth := colWidthID
	contentWidth := colWidthContent
	projectWidth := colWidthProject
	dateWidth := colWidthDate

	var sb strings.Builder

	// Header
	header := fmt.Sprintf("%-*s | %-*s | %-*s | %-*s",
		idWidth, "ID",
		contentWidth, "Content",
		projectWidth, "Project",
		dateWidth, "Created",
	)
	fmt.Fprintln(&sb, header)
	fmt.Fprintln(&sb, strings.Repeat("-", len(header)))

	// Rows
	for _, item := range items {
		id := truncateString(item.ID, idWidth)
		content := truncateString(item.Content, contentWidth)
		project := truncateString(item.Project, projectWidth)
		date := item.CreatedAt.Format("2006-01-02 15:04")

		row := fmt.Sprintf("%-*s | %-*s | %-*s | %-*s",
			idWidth, id,
			contentWidth, content,
			projectWidth, project,
			dateWidth, date,
		)
		fmt.Fprintln(&sb, row)
	}

	return sb.String()
}

// truncateString truncates a string to the specified maximum length.
// If the string is shorter than or equal to maxLen, it returns the original.
// If maxLen is 3 or less, returns the first maxLen characters without ellipsis.
// Otherwise, appends "..." to indicate truncation.
//
// Parameters:
//   - s: The string to truncate
//   - maxLen: Maximum length of the resulting string
//
// Returns:
//
//	The truncated string with ellipsis if applicable
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + contentTruncationIndicator
}
