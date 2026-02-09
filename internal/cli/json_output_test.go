package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ondrahracek/contextkeeper/internal/models"
	"github.com/ondrahracek/contextkeeper/internal/storage"
)

func setupTestStorage(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "ck-test-*")
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(tmpDir, "items.json")
}

func TestJSONOutput(t *testing.T) {
	storagePath := setupTestStorage(t)
	defer os.RemoveAll(filepath.Dir(storagePath))

	// Set the storage path in config if possible or use a mock
	// For these tests, we'll need to ensure the commands use this path.
	// Since the CLI commands use config.FindStoragePath(""), we might need to set an env var.
	os.Setenv("CK_STORAGE_PATH", storagePath)
	defer os.Unsetenv("CK_STORAGE_PATH")

	// Ensure the directory exists
	os.MkdirAll(filepath.Dir(storagePath), 0755)

	stor := storage.NewStorage(storagePath)
	item := models.ContextItem{
		ID:        "bc2839b5-6a8b-4b2a-9e1e-7b5c4d3e2f1a",
		Content:   "change tokens something",
		Project:   "",
		Tags:      []string{"whatsapp"},
		CreatedAt: time.Date(2026, 2, 8, 21, 25, 0, 0, time.UTC),
	}
	stor.Add(item)
	stor.Save()

	t.Run("list --json", func(t *testing.T) {
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"list", "--json"})

		err := RootCmd.Execute()
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if buf.Len() == 0 {
			t.Fatal("Output buffer is empty")
		}

		var output []map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
			t.Fatalf("Failed to unmarshal JSON output: %v\nOutput: %s", err, buf.String())
		}

		if len(output) != 1 {
			t.Errorf("Expected 1 item, got %d", len(output))
		}
		if output[0]["id"] != "bc2839b5" {
			t.Errorf("Expected id bc2839b5, got %v", output[0]["id"])
		}
	})

	t.Run("status --json", func(t *testing.T) {
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"status", "--json"})

		err := RootCmd.Execute()
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if buf.Len() == 0 {
			t.Fatal("Output buffer is empty")
		}

		var output map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
			t.Fatalf("Failed to unmarshal JSON output: %v\nOutput: %s", err, buf.String())
		}

		if output["totalItems"].(float64) != 1 {
			t.Errorf("Expected 1 total item, got %v", output["totalItems"])
		}
	})

	t.Run("done --json", func(t *testing.T) {
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"done", "bc2839", "--json"})

		err := RootCmd.Execute()
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if buf.Len() == 0 {
			t.Fatal("Output buffer is empty")
		}

		var output map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
			t.Fatalf("Failed to unmarshal JSON output: %v\nOutput: %s", err, buf.String())
		}

		if output["id"] != "bc2839b5" {
			t.Errorf("Expected id bc2839b5, got %v", output["id"])
		}
		if output["status"] != "completed" {
			t.Errorf("Expected status completed, got %v", output["status"])
		}
	})

	t.Run("add --json", func(t *testing.T) {
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"add", "new item", "--json"})

		err := RootCmd.Execute()
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if buf.Len() == 0 {
			t.Fatal("Output buffer is empty")
		}

		var output map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
			t.Fatalf("Failed to unmarshal JSON output: %v\nOutput: %s", err, buf.String())
		}

		if output["status"] != "added" {
			t.Errorf("Expected status added, got %v", output["status"])
		}
		if output["id"] == "" {
			t.Errorf("Expected non-empty id")
		}
	})
}
