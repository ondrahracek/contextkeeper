package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ondrahracek/contextkeeper/internal/models"
)

func TestNewStorage(t *testing.T) {
	stor := NewStorage("/tmp/test-contextkeeper")
	if stor == nil {
		t.Error("NewStorage should return non-nil Storage")
	}
}

func TestStorageCRUD(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "contextkeeper-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	path := filepath.Join(tmpDir, "data.json")
	stor := NewStorage(path)

	// Test Load on non-existent file
	err = stor.Load()
	if err != nil {
		t.Errorf("Load() error: %v", err)
	}

	// Test GetAll on empty storage
	items := stor.GetAll()
	if len(items) != 0 {
		t.Errorf("GetAll() empty file: got %d items, want 0", len(items))
	}

	// Test Add
	item := models.ContextItem{
		ID:        "test-123",
		Content:   "Test content",
		Project:   "test-project",
		Tags:      []string{"tag1"},
	}
	err = stor.Add(item)
	if err != nil {
		t.Errorf("Add() error: %v", err)
	}

	// Test GetAll after Add
	items = stor.GetAll()
	if len(items) != 1 {
		t.Errorf("GetAll() after Add: got %d items, want 1", len(items))
	}

	// Test GetByID
	got, err := stor.GetByID("test-123")
	if err != nil {
		t.Errorf("GetByID() error: %v", err)
	}
	if got.Content != "Test content" {
		t.Errorf("GetByID() content: got %s, want Test content", got.Content)
	}

	// Test GetByID not found
	_, err = stor.GetByID("not-found")
	if err != ErrItemNotFound {
		t.Errorf("GetByID() not-found: got %v, want ErrItemNotFound", err)
	}

	// Test Update
	item.Content = "Updated content"
	err = stor.Update(item)
	if err != nil {
		t.Errorf("Update() error: %v", err)
	}
	got, _ = stor.GetByID("test-123")
	if got.Content != "Updated content" {
		t.Errorf("Update() content: got %s, want Updated content", got.Content)
	}

	// Test Archive
	err = stor.Archive("test-123")
	if err != nil {
		t.Errorf("Archive() error: %v", err)
	}
	got, _ = stor.GetByID("test-123")
	// Note: Archive doesn't have a method to check, just verify no error

	// Test Delete
	err = stor.Delete("test-123")
	if err != nil {
		t.Errorf("Delete() error: %v", err)
	}
	items = stor.GetAll()
	if len(items) != 0 {
		t.Errorf("Delete(): got %d items, want 0", len(items))
	}
}

func TestStorageErrorTypes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "contextkeeper-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	path := filepath.Join(tmpDir, "data.json")
	stor := NewStorage(path)

	// Test that ErrItemNotFound is the correct type
	_, err = stor.GetByID("non-existent")
	if err == nil {
		t.Error("GetByID() non-existent: should return error")
	}

	if err != ErrItemNotFound {
		t.Errorf("GetByID() error type: got %v, want ErrItemNotFound", err)
	}
}

func TestStorageThreadSafety(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "contextkeeper-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	path := filepath.Join(tmpDir, "data.json")
	stor := NewStorage(path)

	// Load initial data
	item := models.ContextItem{ID: "1", Content: "Test"}
	stor.Add(item)

	// Run concurrent operations
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()
			// All goroutines try to read
			items := stor.GetAll()
			if len(items) >= 0 { // Should never fail
				_ = items
			}
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify data integrity
	items := stor.GetAll()
	if len(items) != 1 {
		t.Errorf("After concurrent reads: got %d items, want 1", len(items))
	}
}
