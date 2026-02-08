// Package cli provides the command-line interface for ContextKeeper.
//
// This package implements the Cobra-based CLI for managing project context
// and configuration. See the root.go file for the main command structure.
package cli

import (
	"fmt"
	"time"

	"github.com/ondrahracek/contextkeeper/internal/config"
	"github.com/ondrahracek/contextkeeper/internal/storage"
	"github.com/spf13/cobra"
)

// doneCmd marks a context item as completed.
//
// The command requires an item ID (can be partial prefix) as an argument.
var doneCmd = &cobra.Command{
	Use:     "done <id>",
	Short:   "Mark a context item as completed",
	Long:    "Mark a context item as completed by its ID. Use 'ck list' to see item IDs.",
	Example: `  # Mark an item as done (full ID)
  ck done abc12345-def6-7890-1234-567890abcdef

  # Mark an item as done (partial ID prefix)
  ck done abc12345`,
	Args: cobra.ExactArgs(1),
	RunE: doneCommand,
}

// doneCommand is the execution function for the done command.
// It finds and marks a context item as completed.
func doneCommand(cmd *cobra.Command, args []string) error {
	id := args[0]

	// Load configuration to get storage path
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Initialize storage and load items
	stor := storage.NewStorage(cfg.StoragePath)
	if err := stor.Load(); err != nil {
		return err
	}

	// Get all items and find matching one
	allItems := stor.GetAll()
	found := false
	var itemID string

	for _, item := range allItems {
		// Match by prefix
		if len(item.ID) >= len(id) && item.ID[:len(id)] == id {
			found = true
			itemID = item.ID
			break
		}
	}

	if !found {
		return fmt.Errorf("item not found: %s", id)
	}

	// Get the item and update it
	item, err := stor.GetByID(itemID)
	if err != nil {
		return err
	}

	now := time.Now()
	item.CompletedAt = &now

	if err := stor.Update(item); err != nil {
		return err
	}

	cmd.Printf("Marked item as completed: %s\n", itemID[:8])
	return nil
}

// init registers the done command with the root command.
func init() {
	// Add command to root
	RootCmd.AddCommand(doneCmd)
}
