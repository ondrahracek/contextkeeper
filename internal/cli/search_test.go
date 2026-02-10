package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ondrahracek/contextkeeper/internal/models"
	"github.com/ondrahracek/contextkeeper/internal/storage"
)

// createSearchTestStorage creates a temporary storage with test data for search tests.
func createSearchTestStorage(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "ck-search-test-*")
	if err != nil {
		t.Fatal(err)
	}
	storagePath := filepath.Join(tmpDir, "items.json")

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}

	stor := storage.NewStorage(storagePath)
	
	testItems := []models.ContextItem{
		{
			ID:        "bc2839b5-6a8b-4b2a-9e1e-7b5c4d3e2f1a",
			Content:   "change tokens something - find in whatsapp with",
			Project:   "carscoring-app",
			Tags:      []string{"whatsapp", "tokens"},
			CreatedAt: time.Date(2026, 2, 8, 21, 25, 0, 0, time.UTC),
		},
		{
			ID:        "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			Content:   "add rate limiting to auth middleware",
			Project:   "carscoring-app",
			Tags:      []string{"auth", "security"},
			CreatedAt: time.Date(2026, 2, 9, 10, 30, 0, 0, time.UTC),
		},
		{
			ID:        "def45678-1234-5678-90ab-cdef12345678",
			Content:   "implement user dashboard",
			Project:   "webapp",
			Tags:      []string{"ui", "feature"},
			CreatedAt: time.Date(2026, 2, 9, 14, 15, 0, 0, time.UTC),
			CompletedAt: func() *time.Time { t := time.Date(2026, 2, 10, 9, 0, 0, 0, time.UTC); return &t }(),
		},
		{
			ID:        "xyz11111-2222-3333-4444-555555555555",
			Content:   "API response handling for USER data",
			Project:   "api-service",
			Tags:      []string{"api", "user"},
			CreatedAt: time.Date(2026, 2, 10, 8, 0, 0, 0, time.UTC),
		},
	}

	for _, item := range testItems {
		stor.Add(item)
	}

	if err := stor.Save(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}

	cleanup := func() {
		os.Unsetenv("CK_STORAGE_PATH")
		os.RemoveAll(tmpDir)
		searchTagFilter = ""
		searchShowAll = false
		searchJsonOut = false
	}

	os.Setenv("CK_STORAGE_PATH", storagePath)
	return storagePath, cleanup
}

func TestSearchCommand(t *testing.T) {
	_, cleanup := createSearchTestStorage(t)
	defer cleanup()

	resetFlags := func() {
		searchTagFilter = ""
		searchShowAll = false
		searchJsonOut = false
	}

	runSearchTest := func(name string, args []string, expectedCount int, contentCheck func(map[string]interface{}) bool) {
		t.Run(name, func(t *testing.T) {
			resetFlags()
			buf := new(bytes.Buffer)
			RootCmd.SetOut(buf)
			RootCmd.SetArgs(args)

			if err := RootCmd.Execute(); err != nil {
				t.Fatalf("Execute failed: %v", err)
			}

			if buf.Len() == 0 {
				t.Fatal("Output buffer is empty")
			}

			var output []map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
				t.Fatalf("Failed to unmarshal JSON output: %v\nOutput: %s", err, buf.String())
			}

			if len(output) != expectedCount {
				t.Errorf("Expected %d items, got %d", expectedCount, len(output))
			}

			if expectedCount > 0 && contentCheck != nil {
				found := false
				for _, item := range output {
					if contentCheck(item) {
						found = true
						break
					}
				}
				if !found {
					t.Error("Expected content not found in results")
				}
			}
		})
	}

	runSearchTest(
		"search finds matching content",
		[]string{"search", "auth", "--json"},
		1,
		func(item map[string]interface{}) bool {
			content := strings.ToLower(item["content"].(string))
			return strings.Contains(content, "auth")
		},
	)

	runSearchTest(
		"search finds by tag",
		[]string{"search", "--tag", "whatsapp", "--json"},
		1,
		func(item map[string]interface{}) bool {
			tags := item["tags"].([]interface{})
			for _, tag := range tags {
				if tag.(string) == "whatsapp" {
					return true
				}
			}
			return false
		},
	)

	runSearchTest(
		"search excludes completed by default",
		[]string{"search", "dashboard", "--json"},
		0,
		nil,
	)

	runSearchTest(
		"search includes completed with --all",
		[]string{"search", "--all", "dashboard", "--json"},
		1,
		func(item map[string]interface{}) bool {
			content := strings.ToLower(item["content"].(string))
			return strings.Contains(content, "dashboard")
		},
	)

	runSearchTest(
		"search with no matches returns empty",
		[]string{"search", "nonexistent-query-xyz", "--json"},
		0,
		nil,
	)

	runSearchTest(
		"search is case insensitive",
		[]string{"search", "AUTH", "--json"},
		1,
		func(item map[string]interface{}) bool {
			content := strings.ToLower(item["content"].(string))
			return strings.Contains(content, "auth")
		},
	)

	runSearchTest(
		"search matches in content and tags",
		[]string{"search", "whatsapp", "--json"},
		1,
		func(item map[string]interface{}) bool {
			content := strings.ToLower(item["content"].(string))
			tags := item["tags"].([]interface{})
			if strings.Contains(content, "whatsapp") {
				return true
			}
			for _, tag := range tags {
				if tag.(string) == "whatsapp" {
					return true
				}
			}
			return false
		},
	)

	runSearchTest(
		"search with multiple words matches partial",
		[]string{"search", "rate limiting", "--json"},
		1,
		func(item map[string]interface{}) bool {
			content := strings.ToLower(item["content"].(string))
			return strings.Contains(content, "rate")
		},
	)

	runSearchTest(
		"search finds multiple items with tag",
		[]string{"search", "--tag", "api", "--json"},
		1,
		func(item map[string]interface{}) bool {
			tags := item["tags"].([]interface{})
			for _, tag := range tags {
				if tag.(string) == "api" {
					return true
				}
			}
			return false
		},
	)

	t.Run("search JSON output has correct structure", func(t *testing.T) {
		resetFlags()
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"search", "--json", "auth"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		var output []map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if len(output) == 0 {
			t.Fatal("Expected at least one item")
		}

		item := output[0]
		requiredFields := []string{"id", "fullId", "content", "project", "tags", "createdAt"}
		for _, field := range requiredFields {
			if _, ok := item[field]; !ok {
				t.Errorf("Missing required field: %s", field)
			}
		}
		if len(item["id"].(string)) != 8 {
			t.Errorf("Expected id length 8, got %d", len(item["id"].(string)))
		}
	})

	t.Run("search with empty query returns active items", func(t *testing.T) {
		resetFlags()
		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetArgs([]string{"search", "--json"})

		if err := RootCmd.Execute(); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		var output []map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if len(output) != 3 {
			t.Errorf("Expected 3 active items, got %d", len(output))
		}
	})

	runSearchTest(
		"search USER finds user data",
		[]string{"search", "USER", "--json"},
		1,
		func(item map[string]interface{}) bool {
			content := strings.ToLower(item["content"].(string))
			return strings.Contains(content, "user")
		},
	)
}
