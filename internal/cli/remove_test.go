package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ondrahracek/contextkeeper/internal/models"
	"github.com/ondrahracek/contextkeeper/internal/storage"
)

func TestRemoveCommandSyncFlag(t *testing.T) {
	// Ensure flags are reset after this test completes
	defer func() {
		addSyncFlag = false
		doneSyncFlag = false
		removeSyncFlag = false
		editSyncFlag = false
	}()

	t.Run("remove --force --sync creates sync files", func(t *testing.T) {
		removeSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-remove-sync-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		// Create agent directories
		claudeDir := filepath.Join(tmpDir, ".claude", "rules")
		cursorDir := filepath.Join(tmpDir, ".cursor", "rules")
		os.MkdirAll(claudeDir, 0755)
		os.MkdirAll(cursorDir, 0755)

		storagePath := filepath.Join(tmpDir, "items.json")
		stor := storage.NewStorage(storagePath)
		stor.Add(models.ContextItem{
			ID:      "remove-sync-item-12345",
			Content: "Item to remove with sync",
			Tags:    []string{"test"},
		})
		stor.Add(models.ContextItem{
			ID:      "another-remove-sync-67890",
			Content: "Another item to keep",
			Tags:    []string{"test"},
		})
		stor.Save()

		os.Setenv("CK_STORAGE_PATH", storagePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"remove", "remove-sync-item-12345", "--force", "--sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		output := buf.String()

		// Verify remove output
		if !bytes.Contains([]byte(output), []byte("Removed")) {
			t.Errorf("Expected 'Removed' status in output, got: %s", output)
		}

		// Verify sync output
		if !bytes.Contains([]byte(output), []byte("Synced")) {
			t.Errorf("Expected sync message in output, got: %s", output)
		}

		// Verify Claude file was created
		claudeFile := filepath.Join(tmpDir, ".claude", "rules", "ck-context.md")
		if _, err := os.Stat(claudeFile); os.IsNotExist(err) {
			t.Errorf("Expected Claude file to be created")
		}

		// Verify removed item is NOT in synced content
		content, _ := os.ReadFile(claudeFile)
		if bytes.Contains(content, []byte("Item to remove with sync")) {
			t.Errorf("Removed item should not be in synced content")
		}
	})

	t.Run("remove with --sync and --json flags", func(t *testing.T) {
		removeSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-remove-sync-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		// Create agent directories
		claudeDir := filepath.Join(tmpDir, ".claude", "rules")
		cursorDir := filepath.Join(tmpDir, ".cursor", "rules")
		os.MkdirAll(claudeDir, 0755)
		os.MkdirAll(cursorDir, 0755)

		storagePath := filepath.Join(tmpDir, "items.json")
		stor := storage.NewStorage(storagePath)
		stor.Add(models.ContextItem{
			ID:      "json-remove-sync-12345",
			Content: "JSON remove sync item",
			Tags:    []string{"test"},
		})
		stor.Add(models.ContextItem{
			ID:      "keep-item-67890",
			Content: "Item to keep",
			Tags:    []string{"test"},
		})
		stor.Save()

		os.Setenv("CK_STORAGE_PATH", storagePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetErr(buf)
		RootCmd.SetArgs([]string{"remove", "json-remove-sync-12345", "--force", "--sync", "--json"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		output := buf.String()

		// Should contain sync message
		if !bytes.Contains(buf.Bytes(), []byte("Synced")) {
			t.Errorf("Expected sync message, got: %s", output)
		}
	})

	t.Run("remove --force without --sync does not create sync files", func(t *testing.T) {
		removeSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-remove-sync-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		// Create agent directories
		claudeDir := filepath.Join(tmpDir, ".claude", "rules")
		cursorDir := filepath.Join(tmpDir, ".cursor", "rules")
		os.MkdirAll(claudeDir, 0755)
		os.MkdirAll(cursorDir, 0755)

		storagePath := filepath.Join(tmpDir, "items.json")
		stor := storage.NewStorage(storagePath)
		stor.Add(models.ContextItem{
			ID:      "remove-sync-item-12345",
			Content: "Item to remove with sync",
			Tags:    []string{"test"},
		})
		stor.Add(models.ContextItem{
			ID:      "another-remove-sync-67890",
			Content: "Another item to keep",
			Tags:    []string{"test"},
		})
		stor.Save()

		os.Setenv("CK_STORAGE_PATH", storagePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"remove", "another-remove-sync-67890", "--force"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		// Should NOT contain sync message
		if bytes.Contains(buf.Bytes(), []byte("Synced")) {
			t.Errorf("Did not expect sync message without --sync flag, got: %s", buf.String())
		}
	})

	t.Run("remove --sync with sync failure does not fail main operation", func(t *testing.T) {
		removeSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-remove-sync-test-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		// Create agent directories
		claudeDir := filepath.Join(tmpDir, ".claude", "rules")
		cursorDir := filepath.Join(tmpDir, ".cursor", "rules")
		os.MkdirAll(claudeDir, 0755)
		os.MkdirAll(cursorDir, 0755)

		storagePath := filepath.Join(tmpDir, "items.json")
		stor := storage.NewStorage(storagePath)
		stor.Add(models.ContextItem{
			ID:      "remove-sync-item-12345",
			Content: "Item to remove with sync",
			Tags:    []string{"test"},
		})
		stor.Add(models.ContextItem{
			ID:      "another-remove-sync-67890",
			Content: "Another item to keep",
			Tags:    []string{"test"},
		})
		stor.Save()

		os.Setenv("CK_STORAGE_PATH", storagePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		// Make agent directories read-only to cause sync failure
		os.Chmod(claudeDir, 0555)
		defer os.Chmod(claudeDir, 0755)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"remove", "another-remove-sync-67890", "--force", "--sync"})

		// This should NOT fail - sync failure should not fail the main operation
		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Expected no error when sync fails, got: %v", err)
		}

		// Item should still be removed
		stor = storage.NewStorage(storagePath)
		stor.Load()
		items := stor.GetAll()
		if len(items) != 1 {
			t.Errorf("Expected 1 item after removal, got %d", len(items))
		}
	})

	t.Run("remove --sync only includes remaining items", func(t *testing.T) {
		removeSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-remove-sync-test-*")
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
			ID:      "keep-item-sync-11111",
			Content: "Item that should be kept and synced",
			Tags:    []string{"keep"},
		})
		syncStor.Add(models.ContextItem{
			ID:      "remove-item-sync-22222",
			Content: "Item to be removed",
			Tags:    []string{"remove"},
		})
		syncStor.Save()

		os.Setenv("CK_STORAGE_PATH", syncStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		// Remove one item with sync
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"remove", "remove-item-sync-22222", "--force", "--sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		// Verify Claude file content
		claudeFile := filepath.Join(tmpDir, ".claude", "rules", "ck-context.md")
		content, _ := os.ReadFile(claudeFile)

		// Kept item should be in synced content
		if !bytes.Contains(content, []byte("Item that should be kept and synced")) {
			t.Errorf("Kept item should be in synced content")
		}

		// Removed item should NOT be in synced content
		if bytes.Contains(content, []byte("Item to be removed")) {
			t.Errorf("Removed item should not be in synced content")
		}
	})
}
