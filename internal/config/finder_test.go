package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindStoragePath_WithExplicitPath(t *testing.T) {
	f := NewFinder()

	// Test that explicit path gets .contextkeeper appended
	result := f.FindStoragePath("./project")

	expected := filepath.Join("./project", ".contextkeeper")
	if result != expected {
		t.Errorf("FindStoragePath(%q) = %q, want %q", "./project", result, expected)
	}
}

func TestFindStoragePath_WithAbsolutePath(t *testing.T) {
	f := NewFinder()

	// Test with absolute path
	absPath := "/tmp/test-project"
	result := f.FindStoragePath(absPath)

	expected := filepath.Join(absPath, ".contextkeeper")
	if result != expected {
		t.Errorf("FindStoragePath(%q) = %q, want %q", absPath, result, expected)
	}
}

func TestFindStoragePath_EmptyPath_UsesSearchStrategy(t *testing.T) {
	f := NewFinder()

	// Create a temporary directory with .contextkeeper
	tmpDir, err := os.MkdirTemp("", "ck-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	contextDir := filepath.Join(tmpDir, ".contextkeeper")
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		t.Fatalf("Failed to create .contextkeeper dir: %v", err)
	}

	// Save original cwd and restore after
	origCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(origCwd)

	// Change to the temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	// With empty path, should find the local .contextkeeper
	result := f.FindStoragePath("")
	// When run from tmpDir, the result should be ".contextkeeper" (relative)
	if result != ".contextkeeper" {
		t.Errorf("FindStoragePath(%q) = %q, want %q (local context)", "", result, ".contextkeeper")
	}
}

func TestFindStoragePath_EnvironmentVariable(t *testing.T) {
	f := NewFinder()

	// Set environment variable
	os.Setenv("CK_STORAGE_PATH", "/custom/path")
	defer os.Unsetenv("CK_STORAGE_PATH")

	result := f.FindStoragePath("")
	if result != "/custom/path" {
		t.Errorf("FindStoragePath(%q) = %q, want %q (from env)", "", result, "/custom/path")
	}
}
