package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestContextItemJSON(t *testing.T) {
	now := time.Now()
	item := ContextItem{
		ID:        "test-123",
		Content:   "Test content",
		Project:   "test-project",
		Tags:      []string{"tag1", "tag2"},
		CreatedAt: now,
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded ContextItem
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.ID != item.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, item.ID)
	}
	if decoded.Content != item.Content {
		t.Errorf("Content mismatch: got %s, want %s", decoded.Content, item.Content)
	}
	if decoded.Project != item.Project {
		t.Errorf("Project mismatch: got %s, want %s", decoded.Project, item.Project)
	}
}

func TestContextItemOptionalFields(t *testing.T) {
	item := ContextItem{
		ID:      "test-123",
		Content: "Test without optional fields",
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Check that omitempty fields are not present when empty
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	if _, exists := m["project"]; exists {
		t.Error("project should be omitted when empty")
	}
	if _, exists := m["tags"]; exists {
		t.Error("tags should be omitted when empty")
	}
}
