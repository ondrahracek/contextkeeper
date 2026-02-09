// Package cli provides the command-line interface for ContextKeeper.
//
// This package implements the Cobra-based CLI for managing context and
// configuration. See the root.go file for the main command structure.
package cli

import (
	"fmt"
	"strings"

	"github.com/ondrahracek/contextkeeper/internal/config"
	"github.com/ondrahracek/contextkeeper/internal/storage"
	"github.com/spf13/cobra"
)

// removeCmd removes a context item from storage.
//
// The command requires an item ID and can optionally skip confirmation
// with the --force flag.
var removeCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Remove a context item",
	Long:  "Remove a context item by its ID. Use --force to skip the confirmation prompt.",
	Example: `  # Remove with confirmation
  ck remove abc12345

  # Remove without confirmation
  ck remove abc12345 --force`,
	Args: cobra.ExactArgs(1),
	RunE: removeCommand,
}

// forceDelete skips the confirmation prompt when true.
var forceDelete bool

// removeCommand is the execution function for the remove command.
// It finds and removes a context item from storage.
func removeCommand(cmd *cobra.Command, args []string) error {
	id := args[0]

	// Initialize storage and load items
	stor := storage.NewStorage(config.FindStoragePath(""))
	if err := stor.Load(); err != nil {
		return fmt.Errorf("failed to load storage: %w", err)
	}

	// Get all items and find the full ID
	allItems := stor.GetAll()
	var itemID string

	for _, item := range allItems {
		// Match by prefix
		if strings.HasPrefix(item.ID, id) {
			itemID = item.ID
			break
		}
	}

	if itemID == "" {
		return fmt.Errorf("item not found: %s", id)
	}

	// Confirm removal unless --force is set
	if !forceDelete {
		cmd.Printf("Remove item: %s\n", itemID[:8])
		fmt.Print("Are you sure? (y/N): ")
		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Delete the item from storage
	if err := stor.Delete(itemID); err != nil {
		return fmt.Errorf("failed to delete item %q: %w", itemID, err)
	}

	// Display result
	displayID := id
	if len(displayID) > 8 {
		displayID = displayID[:8]
	}
	cmd.Printf("Removed item: %s\n", displayID)
	return nil
}

// init registers the remove command with the root command.
func init() {
	// Register command flags
	removeCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Skip confirmation and permanently delete")

	// Add command to root
	RootCmd.AddCommand(removeCmd)
}
