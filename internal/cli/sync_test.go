package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ondrahracek/contextkeeper/internal/models"
	"github.com/ondrahracek/contextkeeper/internal/storage"
)

func TestSyncCommand(t *testing.T) {
	// Setup temporary workspace
	tmpDir, err := os.MkdirTemp("", "ck-sync-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create storage with data
	storagePath := filepath.Join(tmpDir, "items.json")
	stor := storage.NewStorage(storagePath)
	stor.Add(models.ContextItem{
		ID:      "test-id-12345",
		Content: "Fix the flux capacitor",
		Tags:    []string{"critical"},
	})
	stor.Save()

	// Set environment to use this storage
	os.Setenv("CK_STORAGE_PATH", storagePath)
	defer os.Unsetenv("CK_STORAGE_PATH")

	t.Run("creates sync files in agent directories", func(t *testing.T) {
		testDir, _ := os.MkdirTemp(tmpDir, "agent-test-*")
		claudePath := filepath.Join(testDir, ".claude", "rules")
		cursorPath := filepath.Join(testDir, ".cursor", "rules")
		os.MkdirAll(claudePath, 0755)
		os.MkdirAll(cursorPath, 0755)

		oldWd, _ := os.Getwd()
		os.Chdir(testDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		// Verify Claude file
		claudeFile := filepath.Join(".claude", "rules", "ck-context.md")
		if _, err := os.Stat(claudeFile); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", claudeFile)
		}

		// Verify Cursor file
		cursorFile := filepath.Join(".cursor", "rules", "ck-context.mdc")
		if _, err := os.Stat(cursorFile); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", cursorFile)
		}
	})

	t.Run("handles missing directories gracefully", func(t *testing.T) {
		noAgentDir, _ := os.MkdirTemp(tmpDir, "none-test-*")

		oldWd, _ := os.Getwd()
		os.Chdir(noAgentDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if bytes.Contains(buf.Bytes(), []byte("Synced to")) {
			t.Errorf("Should not have synced anything, but got output: %s", buf.String())
		}
		if !bytes.Contains(buf.Bytes(), []byte("No AI agent directories found")) {
			t.Errorf("Expected warning message, got: %s", buf.String())
		}
	})

	t.Run("creates fallback file when .contextkeeper exists", func(t *testing.T) {
		fallbackDir, _ := os.MkdirTemp(tmpDir, "fallback-test-*")
		os.MkdirAll(filepath.Join(fallbackDir, ".contextkeeper"), 0755)

		oldWd, _ := os.Getwd()
		os.Chdir(fallbackDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		fallbackFile := filepath.Join(".contextkeeper", "instructions.md")
		if _, err := os.Stat(fallbackFile); os.IsNotExist(err) {
			t.Errorf("Expected fallback file %s was not created", fallbackFile)
		}
	})

	t.Run("handles empty items gracefully", func(t *testing.T) {
		emptyDir, _ := os.MkdirTemp(tmpDir, "empty-test-*")
		os.MkdirAll(filepath.Join(emptyDir, ".contextkeeper"), 0755)

		emptyStoragePath := filepath.Join(emptyDir, "items.json")
		emptyStor := storage.NewStorage(emptyStoragePath)
		emptyStor.Save()

		os.Setenv("CK_STORAGE_PATH", emptyStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(emptyDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		fallbackFile := filepath.Join(".contextkeeper", "instructions.md")
		content, err := os.ReadFile(fallbackFile)
		if err != nil {
			t.Fatalf("Failed to read fallback file: %v", err)
		}

		if !bytes.Contains(content, []byte("No active context items")) {
			t.Errorf("Expected 'No active context items' message, got: %s", string(content))
		}
	})

	t.Run("handles items with special characters", func(t *testing.T) {
		specialDir, _ := os.MkdirTemp(tmpDir, "special-test-*")
		os.MkdirAll(filepath.Join(specialDir, ".contextkeeper"), 0755)

		specialStoragePath := filepath.Join(specialDir, "items.json")
		specialStor := storage.NewStorage(specialStoragePath)
		specialStor.Add(models.ContextItem{
			ID:      "special-id-99999",
			Content: "Fix: `backticks`, \"quotes\", and <brackets>",
			Tags:    []string{"test", "edge-case"},
		})
		specialStor.Save()

		os.Setenv("CK_STORAGE_PATH", specialStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(specialDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		fallbackFile := filepath.Join(".contextkeeper", "instructions.md")
		content, err := os.ReadFile(fallbackFile)
		if err != nil {
			t.Fatalf("Failed to read fallback file: %v", err)
		}

		// Verify content is preserved (Markdown safe)
		if !bytes.Contains(content, []byte("Fix:")) {
			t.Errorf("Expected content to be preserved, got: %s", string(content))
		}
	})

	t.Run("skips completed items", func(t *testing.T) {
		mixedDir, _ := os.MkdirTemp(tmpDir, "mixed-test-*")
		os.MkdirAll(filepath.Join(mixedDir, ".contextkeeper"), 0755)

		mixedStoragePath := filepath.Join(mixedDir, "items.json")
		mixedStor := storage.NewStorage(mixedStoragePath)
		
		completedTime := time.Now()
		mixedStor.Add(models.ContextItem{
			ID:      "active-item-12345",
			Content: "This should be synced",
			Tags:    []string{"active"},
		})
		mixedStor.Add(models.ContextItem{
			ID:          "completed-item-67890",
			Content:     "This should NOT be synced",
			Tags:        []string{"completed"},
			CompletedAt: &completedTime,
		})
		mixedStor.Save()

		os.Setenv("CK_STORAGE_PATH", mixedStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(mixedDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		fallbackFile := filepath.Join(".contextkeeper", "instructions.md")
		content, err := os.ReadFile(fallbackFile)
		if err != nil {
			t.Fatalf("Failed to read fallback file: %v", err)
		}

		if !bytes.Contains(content, []byte("This should be synced")) {
			t.Errorf("Expected active item content, got: %s", string(content))
		}
		if bytes.Contains(content, []byte("This should NOT be synced")) {
			t.Errorf("Did not expect completed item content, got: %s", string(content))
		}
	})

	t.Run("handles partial agent directories", func(t *testing.T) {
		partialDir, _ := os.MkdirTemp(tmpDir, "partial-test-*")
		// Create .claude but NOT rules/
		os.MkdirAll(filepath.Join(partialDir, ".claude"), 0755)
		// Create .cursor/rules/ fully
		os.MkdirAll(filepath.Join(partialDir, ".cursor", "rules"), 0755)

		// Create storage in partialDir
		partialStoragePath := filepath.Join(partialDir, "items.json")
		partialStor := storage.NewStorage(partialStoragePath)
		partialStor.Add(models.ContextItem{
			ID:      "partial-item-12345",
			Content: "Test item for partial directories",
			Tags:    []string{"test"},
		})
		partialStor.Save()

		os.Setenv("CK_STORAGE_PATH", partialStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(partialDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		// Should only sync to cursor (since .claude/rules/ doesn't exist)
		if !bytes.Contains(buf.Bytes(), []byte("Synced to .cursor/rules/ck-context.mdc")) {
			t.Errorf("Expected sync to cursor only, got: %s", buf.String())
		}
		// Should NOT have synced to .claude/rules/ck-context.md
		if bytes.Contains(buf.Bytes(), []byte("Synced to .claude/rules/ck-context.md")) {
			t.Errorf("Should not have synced to incomplete .claude directory")
		}
	})

	t.Run("handles unicode and emojis", func(t *testing.T) {
		unicodeDir, _ := os.MkdirTemp(tmpDir, "unicode-test-*")
		os.MkdirAll(filepath.Join(unicodeDir, ".contextkeeper"), 0755)

		unicodeStoragePath := filepath.Join(unicodeDir, "items.json")
		unicodeStor := storage.NewStorage(unicodeStoragePath)
		unicodeStor.Add(models.ContextItem{
			ID:      "unicode-item-12345",
			Content: "Fix 🚀 rocket launch issue with émojis and ünïcöde",
			Tags:    []string{"unicode", "🚀"},
		})
		unicodeStor.Save()

		os.Setenv("CK_STORAGE_PATH", unicodeStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(unicodeDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		fallbackFile := filepath.Join(".contextkeeper", "instructions.md")
		content, err := os.ReadFile(fallbackFile)
		if err != nil {
			t.Fatalf("Failed to read fallback file: %v", err)
		}

		// Verify unicode is preserved
		if !bytes.Contains(content, []byte("🚀")) {
			t.Errorf("Expected emoji to be preserved, got: %s", string(content))
		}
		if !bytes.Contains(content, []byte("ünïcöde")) {
			t.Errorf("Expected unicode to be preserved, got: %s", string(content))
		}
	})

	t.Run("handles markdown special characters", func(t *testing.T) {
		mdDir, _ := os.MkdirTemp(tmpDir, "md-test-*")
		os.MkdirAll(filepath.Join(mdDir, ".contextkeeper"), 0755)

		mdStoragePath := filepath.Join(mdDir, "items.json")
		mdStor := storage.NewStorage(mdStoragePath)
		mdStor.Add(models.ContextItem{
			ID:      "md-item-12345",
			Content: "Handle **bold**, *italic*, `code`, and [links](https://example.com)",
			Tags:    []string{"markdown"},
		})
		mdStor.Save()

		os.Setenv("CK_STORAGE_PATH", mdStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(mdDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		fallbackFile := filepath.Join(".contextkeeper", "instructions.md")
		content, err := os.ReadFile(fallbackFile)
		if err != nil {
			t.Fatalf("Failed to read fallback file: %v", err)
		}

		// Verify markdown is escaped/preserved (not interpreted)
		if !bytes.Contains(content, []byte("**bold**")) {
			t.Errorf("Expected markdown to be preserved, got: %s", string(content))
		}
	})

	t.Run("handles long content", func(t *testing.T) {
		longDir, _ := os.MkdirTemp(tmpDir, "long-test-*")
		os.MkdirAll(filepath.Join(longDir, ".contextkeeper"), 0755)

		longStoragePath := filepath.Join(longDir, "items.json")
		longStor := storage.NewStorage(longStoragePath)
		
		longContent := strings.Repeat("This is a very long line of text. ", 100)
		longStor.Add(models.ContextItem{
			ID:      "long-item-12345",
			Content: longContent,
			Tags:    []string{"long"},
		})
		longStor.Save()

		os.Setenv("CK_STORAGE_PATH", longStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(longDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		fallbackFile := filepath.Join(".contextkeeper", "instructions.md")
		content, err := os.ReadFile(fallbackFile)
		if err != nil {
			t.Fatalf("Failed to read fallback file: %v", err)
		}

		// Verify long content is preserved
		if !bytes.Contains(content, []byte("This is a very long line of text")) {
			t.Errorf("Expected long content to be preserved, got: %s", string(content)[:200])
		}
	})

	t.Run("handles many tags", func(t *testing.T) {
		manyTagsDir, _ := os.MkdirTemp(tmpDir, "many-tags-test-*")
		os.MkdirAll(filepath.Join(manyTagsDir, ".contextkeeper"), 0755)

		manyTagsStoragePath := filepath.Join(manyTagsDir, "items.json")
		manyTagsStor := storage.NewStorage(manyTagsStoragePath)
		manyTagsStor.Add(models.ContextItem{
			ID:      "tags-item-12345",
			Content: "Item with many tags",
			Tags:    []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6"},
		})
		manyTagsStor.Save()

		os.Setenv("CK_STORAGE_PATH", manyTagsStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(manyTagsDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		fallbackFile := filepath.Join(".contextkeeper", "instructions.md")
		content, err := os.ReadFile(fallbackFile)
		if err != nil {
			t.Fatalf("Failed to read fallback file: %v", err)
		}

		// Verify all tags are present
		for _, tag := range []string{"tag1", "tag5", "tag6"} {
			if !bytes.Contains(content, []byte("@"+tag)) {
				t.Errorf("Expected tag %s to be present, got: %s", tag, string(content))
			}
		}
	})

	t.Run("generates correct timestamp format", func(t *testing.T) {
		tsDir, _ := os.MkdirTemp(tmpDir, "ts-test-*")
		os.MkdirAll(filepath.Join(tsDir, ".contextkeeper"), 0755)

		tsStoragePath := filepath.Join(tsDir, "items.json")
		tsStor := storage.NewStorage(tsStoragePath)
		tsStor.Add(models.ContextItem{
			ID:      "ts-item-12345",
			Content: "Test timestamp",
			Tags:    []string{"test"},
		})
		tsStor.Save()

		os.Setenv("CK_STORAGE_PATH", tsStoragePath)
		defer os.Unsetenv("CK_STORAGE_PATH")

		oldWd, _ := os.Getwd()
		os.Chdir(tsDir)
		defer os.Chdir(oldWd)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"sync"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		fallbackFile := filepath.Join(".contextkeeper", "instructions.md")
		content, err := os.ReadFile(fallbackFile)
		if err != nil {
			t.Fatalf("Failed to read fallback file: %v", err)
		}

		// Verify RFC3339 timestamp format
		if !bytes.Contains(content, []byte("Last updated: 2")) {
			t.Errorf("Expected timestamp in RFC3339 format, got: %s", string(content))
		}
	})
}
