// Package utils provides utility functions for the contextkeeper application.
// It includes formatting helpers for output, tag parsing/validation, UUID generation,
// and time formatting utilities.
package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// OpenEditor opens the system editor with optional initial content.
// The function detects the user's preferred editor from environment variables
// or falls back to common editors (vim, vi, nano, code, notepad).
//
// Parameters:
//   - initialContent: Optional content to pre-populate in the editor
//
// Returns:
//   - The content edited by the user
//   - An error if no editor is found or if editing fails
func OpenEditor(initialContent string) (string, error) {
	editor := detectEditor()
	if editor == "" {
		return "", fmt.Errorf("no suitable editor found")
	}

	// Create a temporary file to hold the content
	tmpDir := os.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "contextkeeper-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write initial content if provided
	if initialContent != "" {
		if _, err := tmpFile.WriteString(initialContent); err != nil {
			return "", fmt.Errorf("failed to write initial content: %w", err)
		}
		tmpFile.Close()
	}

	// Open the editor with the temporary file
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor failed: %w", err)
	}

	// Read the edited content back from the file
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read editor result: %w", err)
	}

	return string(content), nil
}

// detectEditor determines the user's preferred editor by checking
// environment variables and common fallbacks in order of preference.
//
// The detection order is:
//  1. EDITOR environment variable
//  2. VISUAL environment variable
//  3. vim
//  4. vi
//  5. nano
//  6. code (VS Code)
//  7. notepad (Windows)
//
// Returns:
//
//	The absolute path to the detected editor, or empty string if none found
func detectEditor() string {
	// Define editors in order of preference with environment variables first
	// then common command-line editors as fallbacks
	editors := []string{
		os.Getenv("EDITOR"),
		os.Getenv("VISUAL"),
		"vim",
		"vi",
		"nano",
		"code",
		"notepad",
	}

	for _, editor := range editors {
		if editor != "" {
			// Check if the editor exists and is executable
			path, err := exec.LookPath(editor)
			if err == nil {
				return path
			}
		}
	}
	return ""
}
