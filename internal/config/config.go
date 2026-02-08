// Package config provides configuration management for ContextKeeper.
//
// This package handles loading, saving, and accessing application configuration
// including storage paths and user preferences. It integrates with the Finder
// package to locate storage locations across different platforms.
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ondrahracek/contextkeeper/internal/models"
)

var cfg *models.Config

// Load reads and parses the configuration from the storage path.
//
// If no configuration file exists, a new Config with default settings is returned.
// The storage path is automatically determined using the Finder.
//
// Returns:
//   - (*models.Config): The loaded or default configuration
//   - (error): An error if reading or parsing fails
func Load() (*models.Config, error) {
	finder := NewFinder()
	storagePath := finder.FindStoragePath("")

	configPath := filepath.Join(storagePath, "config.json")

	cfg = &models.Config{
		StoragePath: storagePath,
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, wrapFileError(err, configPath, "read")
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, wrapConfigError(err, "parse")
	}

	return cfg, nil
}

// Get returns the currently loaded configuration.
//
// Returns nil if Load() has not been called.
func Get() *models.Config {
	return cfg
}

// Save writes the current configuration to the storage path.
//
// If the storage directory does not exist, it is created with appropriate permissions.
// Returns an error if writing fails.
func Save() error {
	if cfg == nil {
		return nil
	}

	finder := NewFinder()
	storagePath := finder.FindStoragePath("")

	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return wrapFileError(err, storagePath, "create directory")
	}

	configPath := filepath.Join(storagePath, "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return wrapConfigError(err, "serialize")
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return wrapFileError(err, configPath, "write")
	}

	return nil
}

// SetConfig sets the current configuration to the provided Config instance.
//
// This is primarily used for testing or when loading configuration from non-standard sources.
func SetConfig(c *models.Config) {
	cfg = c
}

// GetDefaultProject returns the default project name from the configuration.
//
// Returns an empty string if no configuration is loaded or if no default project is set.
func GetDefaultProject() string {
	if cfg == nil {
		return ""
	}
	return cfg.DefaultProject
}

// GetGlobalDefault returns the global default storage path based on the current OS.
//
// The path follows platform-specific conventions:
//   - Windows: %APPDATA%\ContextKeeper
//   - macOS: $HOME/Library/Application Support/ContextKeeper
//   - Linux/BSD: $HOME/.local/share/contextkeeper
func GetGlobalDefault() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "ContextKeeper")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "ContextKeeper")
	default: // linux, freebsd, etc.
		return filepath.Join(os.Getenv("HOME"), ".local", "share", "contextkeeper")
	}
}

// wrapFileError creates a descriptive error message for file operations.
func wrapFileError(err error, path string, operation string) error {
	return &fileError{
		operation: operation,
		path:      path,
		original:  err,
	}
}

// wrapConfigError creates a descriptive error message for configuration errors.
func wrapConfigError(err error, operation string) error {
	return &configError{
		operation: operation,
		original:  err,
	}
}

// fileError represents an error that occurred during a file operation.
type fileError struct {
	operation string
	path      string
	original  error
}

func (e *fileError) Error() string {
	return "failed to " + e.operation + " " + e.path + ": " + e.original.Error()
}

func (e *fileError) Unwrap() error {
	return e.original
}

// configError represents an error that occurred during configuration processing.
type configError struct {
	operation string
	original  error
}

func (e *configError) Error() string {
	return "configuration " + e.operation + " error: " + e.original.Error()
}

func (e *configError) Unwrap() error {
	return e.original
}

// IsConfigError checks if an error is a configuration-related error.
func IsConfigError(err error) bool {
	var configErr *configError
	return errors.As(err, &configErr)
}

// IsFileError checks if an error is a file-related error.
func IsFileError(err error) bool {
	var fileErr *fileError
	return errors.As(err, &fileErr)
}
