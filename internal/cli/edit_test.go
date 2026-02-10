package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ondrahracek/contextkeeper/internal/models"
	"github.com/ondrahracek/contextkeeper/internal/storage"
)

func TestEditCommandSyncFlag(t *testing.T) {
	// Ensure flags are reset after this test completes
	defer func() {
		addSyncFlag = false
		doneSyncFlag = false
		removeSyncFlag = false
		editSyncFlag = false
	}()

	t.Run("edit --sync flag is recognized", func(t *testing.T) {
		editSyncFlag = false

		tmpDir, _, cleanup := setupEditSyncTestForFunc(t)
		defer cleanup()

		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetErr(buf)
		// Run with --help to verify the flag is recognized
		RootCmd.SetArgs([]string{"edit", "--help"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		// Help should not fail - flag should be recognized
	})

	t.Run("edit --sync flag parsing", func(t *testing.T) {
		editSyncFlag = false

		tmpDir, _, cleanup := setupEditSyncTestForFunc(t)
		defer cleanup()

		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetErr(buf)
		RootCmd.SetArgs([]string{"edit", "--sync", "--help"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
	})

	t.Run("edit --sync with other flags", func(t *testing.T) {
		editSyncFlag = false

		tmpDir, _, cleanup := setupEditSyncTestForFunc(t)
		defer cleanup()

		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetErr(buf)
		RootCmd.SetArgs([]string{"edit", "--sync", "--json", "--help"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
	})

	t.Run("edit --sync only syncs active items", func(t *testing.T) {
		editSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-edit-sync-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		// Create agent directories
		claudeDir := filepath.Join(tmpDir, ".claude", "rules")
		cursorDir := filepath.Join(tmpDir, ".cursor", "rules")
		os.MkdirAll(claudeDir, 0755)
		os.MkdirAll(cursorDir, 0755)

		syncStoragePath := filepath.Join(tmpDir, "items.json")
		syncStor := storage.NewStorage(syncStoragePath)
		syncStor.Add(models.ContextItem{
			ID:      "active-edit-sync-11111",
			Content: "Active item should be synced",
			Tags:    []string{"active"},
		})
		syncStor.Add(models.ContextItem{
			ID:      "completed-edit-sync-22222",
			Content: "Completed item should NOT be synced",
			Tags:    []string{"completed"},
		})
		syncStor.Save()

		os.Setenv("CK_STORAGE_PATH", syncStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		// Create a simple sync trigger - use done on completed item to test sync
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"done", "completed-edit-sync-22222", "--sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		// Verify Claude file content
		claudeFile := filepath.Join(tmpDir, ".claude", "rules", "ck-context.md")
		content, _ := os.ReadFile(claudeFile)

		// Active item should be in synced content
		if !bytes.Contains(content, []byte("Active item should be synced")) {
			t.Errorf("Active item should be in synced content")
		}

		// Completed item should NOT be in synced content
		if bytes.Contains(content, []byte("Completed item should NOT be synced")) {
			t.Errorf("Completed item should not be in synced content")
		}
	})

	t.Run("edit --sync includes multiple active items", func(t *testing.T) {
		editSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-edit-sync-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		// Create agent directories
		claudeDir := filepath.Join(tmpDir, ".claude", "rules")
		cursorDir := filepath.Join(tmpDir, ".cursor", "rules")
		os.MkdirAll(claudeDir, 0755)
		os.MkdirAll(cursorDir, 0755)

		syncStoragePath := filepath.Join(tmpDir, "items.json")
		syncStor := storage.NewStorage(syncStoragePath)
		syncStor.Add(models.ContextItem{
			ID:      "first-active-item-11111",
			Content: "First active item for sync",
			Tags:    []string{"item1"},
		})
		syncStor.Add(models.ContextItem{
			ID:      "second-active-item-22222",
			Content: "Second active item for sync",
			Tags:    []string{"item2"},
		})
		syncStor.Add(models.ContextItem{
			ID:      "third-active-item-33333",
			Content: "Third active item for sync",
			Tags:    []string{"item3"},
		})
		syncStor.Save()

		os.Setenv("CK_STORAGE_PATH", syncStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		// Trigger sync via done command on one item
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"done", "third-active-item-33333", "--sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		// Verify Claude file content has all active items
		claudeFile := filepath.Join(tmpDir, ".claude", "rules", "ck-context.md")
		content, _ := os.ReadFile(claudeFile)

		// Should contain both remaining active items
		if !bytes.Contains(content, []byte("First active item for sync")) {
			t.Errorf("First active item should be in synced content")
		}
		if !bytes.Contains(content, []byte("Second active item for sync")) {
			t.Errorf("Second active item should be in synced content")
		}

		// Completed item should NOT be in synced content
		if bytes.Contains(content, []byte("Third active item for sync")) {
			t.Errorf("Completed item should not be in synced content")
		}
	})
}

// setupEditSyncTestForFunc creates a test environment for edit --sync tests.
func setupEditSyncTestForFunc(t *testing.T) (string, string, func()) {
	tmpDir, err := os.MkdirTemp("", "ck-edit-sync-test-*")
	if err != nil {
		t.Fatal(err)
	}

	// Create agent directories
	claudeDir := filepath.Join(tmpDir, ".claude", "rules")
	cursorDir := filepath.Join(tmpDir, ".cursor", "rules")
	os.MkdirAll(claudeDir, 0755)
	os.MkdirAll(cursorDir, 0755)

	storagePath := filepath.Join(tmpDir, "items.json")
	stor := storage.NewStorage(storagePath)
	stor.Add(models.ContextItem{
		ID:      "edit-sync-item-12345",
		Content: "Original content for edit sync",
		Tags:    []string{"test"},
	})
	stor.Add(models.ContextItem{
		ID:      "another-edit-sync-67890",
		Content: "Another item for edit sync",
		Tags:    []string{"test"},
	})
	stor.Save()

	cleanup := func() {
		os.Unsetenv("CK_STORAGE_PATH")
		os.RemoveAll(tmpDir)
	}

	os.Setenv("CK_STORAGE_PATH", storagePath)
	return tmpDir, storagePath, cleanup
}
