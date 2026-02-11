package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestInitCommand_WithPathFlag(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "ck-init-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	projectDir := filepath.Join(tmpDir, "my-project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	// Save original pathFlag and restore after
	origPathFlag := pathFlag
	defer func() { pathFlag = origPathFlag }()

	// Change to project directory
	origCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(origCwd)

	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	// Set pathFlag to the project directory (relative path)
	pathFlag = "."

	// Run init via RootCmd
	buf := new(bytes.Buffer)
	RootCmd.SetOut(buf)
	RootCmd.SetArgs([]string{"init"})

	if err := RootCmd.Execute(); err != nil {
		t.Errorf("initCommand failed: %v", err)
	}

	// Verify .contextkeeper directory was created
	contextDir := filepath.Join(projectDir, ".contextkeeper")
	if _, err := os.Stat(contextDir); os.IsNotExist(err) {
		t.Errorf(".contextkeeper directory was not created at %q", contextDir)
	}

	// Verify items.json was created
	itemsFile := filepath.Join(contextDir, "items.json")
	if _, err := os.Stat(itemsFile); os.IsNotExist(err) {
		t.Errorf("items.json was not created at %q", itemsFile)
	}
}

func TestInitCommand_WithoutPathFlag(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "ck-init-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to the temp directory
	origCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	defer os.Chdir(origCwd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	// Save original pathFlag and restore after
	origPathFlag := pathFlag
	defer func() { pathFlag = origPathFlag }()

	// pathFlag should be empty
	pathFlag = ""

	// Run init via RootCmd
	buf := new(bytes.Buffer)
	RootCmd.SetOut(buf)
	RootCmd.SetArgs([]string{"init"})

	if err := RootCmd.Execute(); err != nil {
		t.Errorf("initCommand failed: %v", err)
	}

	// Verify .contextkeeper directory was created in current directory
	contextDir := filepath.Join(tmpDir, ".contextkeeper")
	if _, err := os.Stat(contextDir); os.IsNotExist(err) {
		t.Errorf(".contextkeeper directory was not created at %q", contextDir)
	}

	// Verify items.json was created
	itemsFile := filepath.Join(contextDir, "items.json")
	if _, err := os.Stat(itemsFile); os.IsNotExist(err) {
		t.Errorf("items.json was not created at %q", itemsFile)
	}
}
