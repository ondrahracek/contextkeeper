// Package models provides data structures for ContextKeeper items.
//
// This package contains the core domain models used throughout the application,
// including the ContextItem struct which represents a single context entry.
package models

import "time"

// ContextItem represents a single context item stored in ContextKeeper.
//
// ContextItems are the fundamental units of storage, containing content along with
// metadata such as project association, tags, and completion status.
type ContextItem struct {
	// ID is the unique identifier for this context item
	ID string `json:"id"`

	// Content is the main text content of this context item
	Content string `json:"content"`

	// Project is the associated project name (optional)
	Project string `json:"project,omitempty"`

	// Tags is a list of tags for categorization (optional)
	Tags []string `json:"tags,omitempty"`

	// CreatedAt is the timestamp when this item was created
	CreatedAt time.Time `json:"created_at"`

	// CompletedAt is the timestamp when this item was completed (nil if not completed)
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	// Archived indicates whether this item has been archived
	Archived bool `json:"archived"`
}

// IsCompleted returns true if the context item has been completed.
//
// A completed item has a non-nil CompletedAt timestamp.
func (c *ContextItem) IsCompleted() bool {
	return c.CompletedAt != nil
}

// IsArchived returns true if the context item has been archived.
func (c *ContextItem) IsArchived() bool {
	return c.Archived
}
