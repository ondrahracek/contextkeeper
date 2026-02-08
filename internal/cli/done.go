// Package cli provides the command-line interface for ContextKeeper.
//
// This package implements the Cobra-based CLI for managing context and
// configuration. See the root.go file for the main command structure.
package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ondrahracek/contextkeeper/internal/config"
	"github.com/ondrahracek/contextkeeper/internal/models"
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

  # Mark an item as done (partial ID prefix - at least 6 chars recommended)
  ck done abc12345`,
	Args: cobra.ExactArgs(1),
	RunE: doneCommand,
}

// doneCommand is the execution function for the done command.
// It finds and marks a context item as completed using prefix matching.
func doneCommand(cmd *cobra.Command, args []string) error {
	id := args[0]

	// Initialize storage
	stor := storage.NewStorage(config.FindStoragePath(""))
	if err := stor.Load(); err != nil {
		return err
	}

	// Try exact match first
	item, err := stor.GetByID(id)
	if err == nil {
		// Exact match found
		return markItemComplete(stor, cmd, item)
	}

	// Try prefix match
	if errors.Is(err, storage.ErrItemNotFound) || errors.Is(err, storage.ErrAmbiguousID) {
		item, err = stor.GetByPrefix(id)
		if err != nil {
			if errors.Is(err, storage.ErrItemNotFound) {
				return fmt.Errorf("item not found: %s", id)
			}
			if errors.Is(err, storage.ErrAmbiguousID) {
				return showAmbiguousMatches(stor, cmd, id)
			}
		}
		// Found unique match
		return markItemComplete(stor, cmd, item)
	}

	return err
}

// markItemComplete marks an item as completed.
func markItemComplete(stor storage.Storage, cmd *cobra.Command, item models.ContextItem) error {
	now := time.Now()
	item.CompletedAt = &now

	if err := stor.Update(item); err != nil {
		return err
	}

	cmd.Printf("Marked item as completed: %s\n", item.ID[:8])
	return nil
}

// showAmbiguousMatches shows all items matching the prefix.
func showAmbiguousMatches(stor storage.Storage, cmd *cobra.Command, prefix string) error {
	allItems := stor.GetAll()
	var matches []models.ContextItem
	for _, item := range allItems {
		if strings.HasPrefix(item.ID, prefix) {
			matches = append(matches, item)
		}
	}

	if len(matches) <= 1 {
		return fmt.Errorf("item not found: %s", prefix)
	}

	fmt.Fprintf(os.Stderr, "Error: %d items match %q:\n", len(matches), prefix)
	for _, item := range matches {
		preview := item.Content
		if len(preview) > 40 {
			preview = preview[:40] + "..."
		}
		fmt.Fprintf(os.Stderr, "  - %s: %s\n", item.ID[:6], preview)
	}
	fmt.Fprintf(os.Stderr, "\nUse more characters to disambiguate:\n")
	for _, item := range matches {
		fmt.Fprintf(os.Stderr, "  ck done %s\n", item.ID)
	}

	return fmt.Errorf("ambiguous ID: %s", prefix)
}

// init registers the done command with the root command.
func init() {
	// Add command to root
	RootCmd.AddCommand(doneCmd)
}
