// Package cli provides the command-line interface for ContextKeeper.
//
// This package implements the Cobra-based CLI for managing project context
// and configuration. See the root.go file for the main command structure.
package cli

import (
	"fmt"

	"github.com/ondrahracek/contextkeeper/internal/config"
	"github.com/ondrahracek/contextkeeper/internal/storage"
	"github.com/ondrahracek/contextkeeper/internal/utils"
	"github.com/spf13/cobra"
)

// editCmd edits an existing context item.
//
// The command opens the system editor with the current content,
// allowing modifications to the item's text.
var editCmd = &cobra.Command{
	Use:     "edit <id>",
	Short:   "Edit a context item",
	Long:    "Edit a context item using the system editor. Opens the current content for modification.",
	Example: `  # Edit an item
  ck edit abc12345`,
	Args: cobra.ExactArgs(1),
	RunE: editCommand,
}

// editCommand is the execution function for the edit command.
// It finds an item, opens it in the editor, and saves changes.
func editCommand(cmd *cobra.Command, args []string) error {
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

	// Get all items and find the item index
	allItems := stor.GetAll()
	found := false
	var itemIndex int
	var originalContent string

	for i, item := range allItems {
		// Match by prefix
		if len(item.ID) >= len(id) && item.ID[:len(id)] == id {
			found = true
			itemIndex = i
			originalContent = item.Content
			break
		}
	}

	if !found {
		return fmt.Errorf("item not found: %s", id)
	}

	// Open editor with current content
	newContent, err := utils.OpenEditor(originalContent)
	if err != nil {
		return err
	}

	// Update the item
	allItems[itemIndex].Content = newContent

	// Save the storage
	stor.SetItems(allItems)
	if err := stor.Save(); err != nil {
		return err
	}

	cmd.Printf("Updated item: %s\n", id[:8])
	return nil
}

// init registers the edit command with the root command.
func init() {
	// Add command to root
	RootCmd.AddCommand(editCmd)
}
