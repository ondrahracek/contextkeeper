// Package models provides data structures for ContextKeeper configuration.
//
// This package contains the Config struct which holds application settings
// and default constants for configuration values.
package models

// Default values for configuration options
const (
	// DefaultDateFormat is the standard date format used throughout the application
	DefaultDateFormat = "2006-01-02"

	// DefaultEditor is the default editor command to use
	DefaultEditor = "vim"
)

// Config represents the application configuration settings.
//
// This configuration controls how ContextKeeper stores and manages context items,
// including storage paths, default project names, and user preferences.
type Config struct {
	// StoragePath is the directory where context data is stored
	StoragePath string `json:"storagePath"`

	// DefaultProject is the project to use when none is specified (optional)
	DefaultProject string `json:"defaultProject,omitempty"`

	// DateFormat is the format string for displaying dates (optional)
	// Defaults to "2006-01-02" if empty
	DateFormat string `json:"dateFormat,omitempty"`

	// Editor is the command to launch for editing context items (optional)
	// Defaults to "vim" if empty
	Editor string `json:"editor,omitempty"`
}
