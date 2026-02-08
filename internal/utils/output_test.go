package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/ondrahracek/contextkeeper/internal/models"
)

// TestFormatItemList tests the formatting of the item list output.
// It ensures that:
// 1. The output contains the 6-character ID prefix.
// 2. The output does not contain raw ANSI bold escape sequences (cross-platform compatibility).
// 3. Completed items are marked with [x] and uncompleted with [ ].
func TestFormatItemList(t *testing.T) {
	now := time.Now()
	items := []models.ContextItem{
		{
			ID:        "bc2839b5-6a8b-4b2a-9e1e-7b5c4d3e2f1a",
			Content:   "Task 1",
			CreatedAt: now,
		},
		{
			ID:          "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
			Content:     "Task 2",
			CreatedAt:   now,
			CompletedAt: &now,
		},
	}

	output := FormatItemList(items, true)

	// 1. Test that output contains the 6-char ID prefix
	if !strings.Contains(output, "bc2839") {
		t.Errorf("Expected output to contain ID prefix 'bc2839', but it didn't")
	}
	if !strings.Contains(output, "a1b2c3") {
		t.Errorf("Expected output to contain ID prefix 'a1b2c3', but it didn't")
	}

	// 2. Test that output does NOT contain raw ANSI escape sequences like \x1b[1m or ←[1m
	// Specifically checking for the problematic colorBold ("\033[1m") around the ID
	if strings.Contains(output, "\033[1m") {
		t.Errorf("Output contains raw ANSI escape sequence for bold (\\033[1m), which causes issues on Windows")
	}

	// 3. Test that completed items show [x] and uncompleted show [ ]
	if !strings.Contains(output, "[ ]") {
		t.Errorf("Expected output to contain '[ ]' for uncompleted item")
	}
	if !strings.Contains(output, "[x]") {
		t.Errorf("Expected output to contain '[x]' for completed item")
	}
}
