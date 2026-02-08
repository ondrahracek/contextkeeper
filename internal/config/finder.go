// Package config provides configuration management for ContextKeeper.
//
// This subpackage contains the Finder type, which is responsible for locating
// the appropriate storage directory for context data. It implements a priority-based
// search strategy that checks multiple locations in order of specificity.
package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// Finder locates configuration and storage paths for ContextKeeper.
//
// The Finder implements a hierarchical search strategy, checking locations from
// most specific (explicit paths, local directories) to most general (global defaults).
type Finder struct{}

// NewFinder creates a new Finder instance.
//
// The Finder is safe for concurrent use and does not maintain internal state.
func NewFinder() *Finder {
	return &Finder{}
}

// FindStoragePath locates the storage path for context data.
//
// The search follows this priority order:
//   1. Explicit path: If a non-empty path is provided, it is used directly
//   2. Local context: Checks for .contextkeeper directory in current directory
//   3. Parent directories: Searches parent directories up to 10 levels for .contextkeeper
//   4. Global default: Falls back to OS-specific default location
//
// Parameters:
//   - explicitPath: A specific path to use; empty string triggers search strategy
//
// Returns:
//   - The resolved storage path as an absolute directory path
func (f *Finder) FindStoragePath(explicitPath string) string {
	// 1. If explicit path provided, use it
	if explicitPath != "" {
		return explicitPath
	}

	// 2. Check local context in current directory
	local := f.checkLocalContext(".")
	if local != "" {
		return local
	}

	// 3. Search parent directories
	cwd, err := os.Getwd()
	if err == nil {
		parents := f.searchParents(cwd)
		if parents != "" {
			return parents
		}
	}

	// 4. Fall back to global default
	return f.getGlobalDefault()
}

// checkLocalContext checks for a local context directory in the given directory.
//
// Looks for a directory named ".contextkeeper" within the specified directory.
// Returns the path if found, or an empty string if not present.
//
// Parameters:
//   - dir: The parent directory to check
//
// Returns:
//   - The full path to .contextkeeper if it exists, empty string otherwise
func (f *Finder) checkLocalContext(dir string) string {
	contextDir := filepath.Join(dir, ".contextkeeper")
	info, err := os.Stat(contextDir)
	if err == nil && info.IsDir() {
		return contextDir
	}
	return ""
}

// searchParents searches parent directories for a context directory.
//
// Starting from the given directory, this method checks each parent directory
// for a .contextkeeper directory. The search is limited to 10 levels to prevent
// infinite loops in edge cases like root filesystem traversal.
//
// Parameters:
//   - dir: The starting directory
//
// Returns:
//   - The first .contextkeeper path found, or empty string if none found
func (f *Finder) searchParents(dir string) string {
	current := dir

	// Limit search to avoid infinite loops (e.g., root filesystem)
	for i := 0; i < 10; i++ {
		local := f.checkLocalContext(current)
		if local != "" {
			return local
		}

		parent := filepath.Dir(current)
		if parent == current {
			break // Reached root
		}
		current = parent
	}

	return ""
}

// getGlobalDefault returns the global default storage path based on the OS.
//
// The path follows platform-specific conventions:
//   - Windows: %APPDATA%\ContextKeeper
//   - macOS: $HOME/Library/Application Support/ContextKeeper
//   - Linux/BSD: $HOME/.local/share/contextkeeper
//
// Returns:
//   - The platform-specific default storage path
func (f *Finder) getGlobalDefault() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "ContextKeeper")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "ContextKeeper")
	default: // linux, freebsd, etc.
		return filepath.Join(os.Getenv("HOME"), ".local", "share", "contextkeeper")
	}
}

// Standalone functions for package-level usage
// These functions create a temporary Finder instance and delegate to its methods.

// FindStoragePath locates the storage path using the default search strategy.
func FindStoragePath(explicitPath string) string {
	return NewFinder().FindStoragePath(explicitPath)
}