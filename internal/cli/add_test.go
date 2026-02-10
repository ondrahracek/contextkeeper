package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ondrahracek/contextkeeper/internal/models"
	"github.com/ondrahracek/contextkeeper/internal/storage"
)

func TestAddCommandSyncFlag(t *testing.T) {
	// Ensure flags are reset after this test completes
	defer func() {
		addSyncFlag = false
		doneSyncFlag = false
		removeSyncFlag = false
		editSyncFlag = false
	}()

	t.Run("add --sync creates sync files", func(t *testing.T) {
		addSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-add-sync-test-*")
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
			ID:      "existing-item-12345",
			Content: "Existing item content",
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
		RootCmd.SetArgs([]string{"add", "New sync item", "--sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		// Verify sync output
		if !bytes.Contains(buf.Bytes(), []byte("Synced")) {
			t.Errorf("Expected sync message in output, got: %s", buf.String())
		}

		// Verify Claude file was created
		claudeFile := filepath.Join(tmpDir, ".claude", "rules", "ck-context.md")
		if _, err := os.Stat(claudeFile); os.IsNotExist(err) {
			t.Errorf("Expected Claude file to be created")
		}

		// Verify Cursor file was created
		cursorFile := filepath.Join(tmpDir, ".cursor", "rules", "ck-context.mdc")
		if _, err := os.Stat(cursorFile); os.IsNotExist(err) {
			t.Errorf("Expected Cursor file to be created")
		}
	})

	t.Run("add with --sync and --json flags", func(t *testing.T) {
		addSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-add-sync-test-*")
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
			ID:      "existing-item-12345",
			Content: "Existing item content",
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
		RootCmd.SetArgs([]string{"add", "JSON sync item", "--sync", "--json"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		output := buf.String()

		// Should contain JSON output for the add operation
		if !bytes.Contains(buf.Bytes(), []byte("added")) {
			t.Errorf("Expected 'added' status in JSON output, got: %s", output)
		}

		// Should contain sync message
		if !bytes.Contains(buf.Bytes(), []byte("Synced")) {
			t.Errorf("Expected sync message, got: %s", output)
		}
	})

	t.Run("add without --sync does not create sync files", func(t *testing.T) {
		addSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-add-sync-test-*")
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
			ID:      "existing-item-12345",
			Content: "Existing item content",
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
		RootCmd.SetArgs([]string{"add", "No sync item"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		// Should NOT contain sync message
		if bytes.Contains(buf.Bytes(), []byte("Synced")) {
			t.Errorf("Did not expect sync message without --sync flag, got: %s", buf.String())
		}
	})

	t.Run("add --sync with no agent directories", func(t *testing.T) {
		addSyncFlag = false

		noAgentDir, err := os.MkdirTemp("", "ck-no-agent-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(noAgentDir)

		noAgentStoragePath := filepath.Join(noAgentDir, "items.json")
		stor := storage.NewStorage(noAgentStoragePath)
		stor.Add(models.ContextItem{
			ID:      "no-agent-item-12345",
			Content: "Test item",
			Tags:    []string{"test"},
		})
		stor.Save()

		os.Setenv("CK_STORAGE_PATH", noAgentStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(noAgentDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"add", "No agent sync item", "--sync"})

		// This should NOT fail - sync failure should not fail the main operation
		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Expected no error when sync fails, got: %v", err)
		}

		// Item should still be added
		stor2 := storage.NewStorage(noAgentStoragePath)
		stor2.Load()
		items := stor2.GetAll()
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}
	})

	t.Run("add --sync with sync failure does not fail main operation", func(t *testing.T) {
		addSyncFlag = false

		tmpDir, err := os.MkdirTemp("", "ck-add-sync-test-*")
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
			ID:      "existing-item-12345",
			Content: "Existing item content",
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
		RootCmd.SetArgs([]string{"add", "Readonly sync item", "--sync"})

		// This should NOT fail - sync failure should not fail the main operation
		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Expected no error when sync fails, got: %v", err)
		}

		// Item should still be added even though sync failed
		stor = storage.NewStorage(storagePath)
		stor.Load()
		items := stor.GetAll()
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}
	})

	t.Run("add --sync with empty storage", func(t *testing.T) {
		addSyncFlag = false

		emptyDir, err := os.MkdirTemp("", "ck-empty-sync-*")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(emptyDir)

		// Create agent directories
		claudeDir := filepath.Join(emptyDir, ".claude", "rules")
		cursorDir := filepath.Join(emptyDir, ".cursor", "rules")
		os.MkdirAll(claudeDir, 0755)
		os.MkdirAll(cursorDir, 0755)

		emptyStoragePath := filepath.Join(emptyDir, "items.json")
		stor := storage.NewStorage(emptyStoragePath)
		stor.Save()

		os.Setenv("CK_STORAGE_PATH", emptyStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(emptyDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"add", "First item", "--sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		// Should sync with 1 item
		if !bytes.Contains(buf.Bytes(), []byte("Synced")) {
			t.Errorf("Expected sync message, got: %s", buf.String())
		}
	})
}
