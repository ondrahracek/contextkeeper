package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ondrahracek/contextkeeper/internal/models"
	"github.com/ondrahracek/contextkeeper/internal/storage"
)

func TestDoneCommandSyncFlag(t *testing.T) {
	// Ensure flags are reset after this test completes
	defer func() {
		addSyncFlag = false
		doneSyncFlag = false
		removeSyncFlag = false
		editSyncFlag = false
	}()

	t.Run("done --sync creates sync files", func(t *testing.T) {
		doneSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-done-sync-test-*")
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
			ID:      "done-sync-item-12345",
			Content: "Item to mark done with sync",
			Tags:    []string{"test"},
		})
		stor.Add(models.ContextItem{
			ID:      "another-sync-item-67890",
			Content: "Another item",
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
		RootCmd.SetArgs([]string{"done", "done-sync", "--sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		output := buf.String()

		// Verify done output
		if !bytes.Contains([]byte(output), []byte("completed")) {
			t.Errorf("Expected 'completed' status in output, got: %s", output)
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

		// Verify completed item is NOT in synced content
		content, _ := os.ReadFile(claudeFile)
		if bytes.Contains(content, []byte("Item to mark done with sync")) {
			t.Errorf("Completed item should not be in synced content")
		}
	})

	t.Run("done with --sync and --json flags", func(t *testing.T) {
		doneSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-done-sync-test-*")
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
			ID:      "done-sync-item-12345",
			Content: "Item to mark done with sync",
			Tags:    []string{"test"},
		})
		stor.Add(models.ContextItem{
			ID:      "json-sync-item-12345",
			Content: "JSON sync item",
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
		RootCmd.SetArgs([]string{"done", "json-sync", "--sync", "--json"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		output := buf.String()

		// Should contain JSON output for the done operation
		if !bytes.Contains(buf.Bytes(), []byte("completed")) {
			t.Errorf("Expected 'completed' status in JSON output, got: %s", output)
		}

		// Should contain sync message
		if !bytes.Contains(buf.Bytes(), []byte("Synced")) {
			t.Errorf("Expected sync message, got: %s", output)
		}
	})

	t.Run("done without --sync does not create sync files", func(t *testing.T) {
		doneSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-done-sync-test-*")
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
			ID:      "done-sync-item-12345",
			Content: "Item to mark done with sync",
			Tags:    []string{"test"},
		})
		stor.Add(models.ContextItem{
			ID:      "another-sync-item-67890",
			Content: "Another item",
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
		RootCmd.SetArgs([]string{"done", "another-sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		// Should NOT contain sync message
		if bytes.Contains(buf.Bytes(), []byte("Synced")) {
			t.Errorf("Did not expect sync message without --sync flag, got: %s", buf.String())
		}
	})

	t.Run("done --sync with sync failure does not fail main operation", func(t *testing.T) {
		doneSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-done-sync-test-*")
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
			ID:      "done-sync-item-12345",
			Content: "Item to mark done with sync",
			Tags:    []string{"test"},
		})
		stor.Add(models.ContextItem{
			ID:      "another-sync-item-67890",
			Content: "Another item",
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
		RootCmd.SetArgs([]string{"done", "another-sync", "--sync"})

		// This should NOT fail - sync failure should not fail the main operation
		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Expected no error when sync fails, got: %v", err)
		}

		// Item should still be marked as done
		stor = storage.NewStorage(storagePath)
		stor.Load()
		items := stor.GetAll()
		for _, item := range items {
			if item.ID == "another-sync-item-67890" {
				if item.CompletedAt == nil {
					t.Errorf("Expected item to be marked as done")
				}
			}
		}
	})

	t.Run("done --sync only includes active items", func(t *testing.T) {
		doneSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-done-sync-test-*")
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
			ID:      "active-item-sync-11111",
			Content: "Active item should be synced",
			Tags:    []string{"active"},
		})
		syncStor.Add(models.ContextItem{
			ID:      "completed-item-sync-22222",
			Content: "Completed item should NOT be synced",
			Tags:    []string{"completed"},
		})
		syncStor.Save()

		os.Setenv("CK_STORAGE_PATH", syncStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(oldWd)

		// Mark one item as done with sync
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"done", "completed-item-sync", "--sync"})

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
}
