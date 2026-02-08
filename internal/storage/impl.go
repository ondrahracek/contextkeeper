package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ondrahracek/contextkeeper/internal/models"
)

// ErrItemNotFound is returned when an item with the specified ID doesn't exist.
var ErrItemNotFound = errors.New("item not found")

// ErrAmbiguousID is returned when multiple items match the given ID prefix.
var ErrAmbiguousID = errors.New("ambiguous ID: multiple items match")

const (
	// ItemsFileName is the default filename for storing items.
	ItemsFileName = "items.json"

	// DefaultDirPerms are the default permissions for created directories.
	DefaultDirPerms = 0755

	// DefaultFilePerms are the default permissions for created files.
	DefaultFilePerms = 0644
)

// Storage defines the interface for persisting context items.
// All implementations must be thread-safe.
type Storage interface {
	// Load reads all items from storage into memory.
	// Returns nil if the storage file doesn't exist (empty state).
	Load() error

	// Save writes all in-memory items to persistent storage.
	// Uses write-through semantics: data is immediately persisted to disk.
	Save() error

	// GetAll returns a copy of all stored items.
	GetAll() []models.ContextItem

	// GetByID retrieves a single item by its full ID.
	// Returns ErrItemNotFound if the item doesn't exist.
	GetByID(id string) (models.ContextItem, error)

	// GetByPrefix retrieves items by ID prefix.
	// Returns all items whose ID starts with the given prefix.
	// If no items match, returns ErrItemNotFound.
	// If multiple items match, returns ErrAmbiguousID.
	GetByPrefix(prefix string) (models.ContextItem, error)

	// Add inserts a new item into storage.
	Add(item models.ContextItem) error

	// Update modifies an existing item.
	// Returns ErrItemNotFound if the item doesn't exist.
	Update(item models.ContextItem) error

	// Archive marks an item as archived without deleting it.
	// Returns ErrItemNotFound if the item doesn't exist.
	Archive(id string) error

	// Delete removes an item from storage permanently.
	// Returns ErrItemNotFound if the item doesn't exist.
	Delete(id string) error

	// SetItems replaces all items with the provided slice.
	SetItems(items []models.ContextItem)
}

// storageImpl provides thread-safe JSON file storage for context items.
// All operations are protected by a sync.RWMutex for concurrent access.
type storageImpl struct {
	mu    sync.RWMutex // Protects all fields
	path  string       // Directory path for storage
	items []models.ContextItem
}

// NewStorage creates a new Storage instance that persists to the specified directory.
//
// Parameters:
//   - path: Directory path where items.json will be stored
//
// Returns:
//   - Storage interface for managing context items
func NewStorage(path string) Storage {
	// Ensure the path is the items.json file path
	if !strings.HasSuffix(path, ItemsFileName) {
		path = filepath.Join(path, ItemsFileName)
	}

	return &storageImpl{
		path:  path,
		items: make([]models.ContextItem, 0),
	}
}

// ensureDir creates the directory for the storage file if it doesn't exist.
func (s *storageImpl) ensureDir() error {
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, DefaultDirPerms); err != nil {
		return fmt.Errorf("failed to create storage directory %q: %w", dir, err)
	}
	return nil
}

// persistLocked saves the current items to the storage file.
// Caller must hold the write lock.
func (s *storageImpl) persistLocked() error {
	if err := s.ensureDir(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s.items, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal items to JSON: %w", err)
	}

	if err := os.WriteFile(s.path, data, DefaultFilePerms); err != nil {
		return fmt.Errorf("failed to write storage file %q: %w", s.path, err)
	}
	return nil
}

// Load reads all items from the storage file into memory.
func (s *storageImpl) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			s.items = make([]models.ContextItem, 0)
			return nil
		}
		return fmt.Errorf("failed to read storage file %q: %w", s.path, err)
	}

	if err := json.Unmarshal(data, &s.items); err != nil {
		return fmt.Errorf("failed to unmarshal JSON from storage file %q: %w", s.path, err)
	}

	return nil
}

// Save writes all in-memory items to the storage file.
func (s *storageImpl) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.persistLocked()
}

// GetAll returns a copy of all stored items.
func (s *storageImpl) GetAll() []models.ContextItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.ContextItem, len(s.items))
	copy(result, s.items)
	return result
}

// GetByID retrieves a single item by its ID.
func (s *storageImpl) GetByID(id string) (models.ContextItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.items {
		if item.ID == id {
			return item, nil
		}
	}

	return models.ContextItem{}, ErrItemNotFound
}

// GetByPrefix retrieves items by ID prefix.
// Returns the item if exactly one matches.
// Returns ErrItemNotFound if no items match.
// Returns ErrAmbiguousID if multiple items match.
func (s *storageImpl) GetByPrefix(prefix string) (models.ContextItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var matches []models.ContextItem
	for _, item := range s.items {
		if strings.HasPrefix(item.ID, prefix) {
			matches = append(matches, item)
		}
	}

	switch len(matches) {
	case 0:
		return models.ContextItem{}, ErrItemNotFound
	case 1:
		return matches[0], nil
	default:
		return models.ContextItem{}, ErrAmbiguousID
	}
}

// Add inserts a new item into storage.
func (s *storageImpl) Add(item models.ContextItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = append(s.items, item)
	return s.persistLocked()
}

// Update modifies an existing item.
func (s *storageImpl) Update(item models.ContextItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.items {
		if s.items[i].ID == item.ID {
			s.items[i] = item
			return s.persistLocked()
		}
	}

	return ErrItemNotFound
}

// Archive marks an item as archived without deleting it.
func (s *storageImpl) Archive(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].Archived = true
			return s.persistLocked()
		}
	}

	return ErrItemNotFound
}

// Delete removes an item from storage permanently.
func (s *storageImpl) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.items {
		if s.items[i].ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return s.persistLocked()
		}
	}

	return ErrItemNotFound
}

// SetItems replaces all items with the provided slice.
func (s *storageImpl) SetItems(items []models.ContextItem) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = items
	s.persistLocked()
}
