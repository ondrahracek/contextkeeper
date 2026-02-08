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

// statusCmd displays a quick overview of context items.
//
// The command shows storage path, total item counts, and the age of
// the oldest item.
var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Show quick overview",
	Long:    "Show a quick overview of context items including counts and storage information.",
	Example: `  # Show status overview
  ck status`,
	Args: cobra.NoArgs,
	RunE: statusCommand,
}

// statusCommand is the execution function for the status command.
// It gathers and displays statistics about stored items.
func statusCommand(cmd *cobra.Command, args []string) error {
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

	// Get all items
	allItems := stor.GetAll()

	// Calculate statistics
	total := len(allItems)
	completed := 0
	var oldest time.Time
	oldestSet := false

	for _, item := range allItems {
		// Count completed items
		if item.CompletedAt != nil {
			completed++
		}

		// Find oldest item
		if !item.CreatedAt.IsZero() {
			if !oldestSet || item.CreatedAt.Before(oldest) {
				oldest = item.CreatedAt
				oldestSet = true
			}
		}
	}

	active := total - completed

	// Print status
	fmt.Println("ContextKeeper Status")
	fmt.Println("===================")
	fmt.Printf("Storage Path: %s\n", cfg.StoragePath)
	fmt.Printf("Total Items: %d\n", total)
	fmt.Printf("Active:      %d\n", active)
	fmt.Printf("Completed:   %d\n", completed)

	if oldestSet {
		daysAgo := int(time.Since(oldest).Hours() / 24)
		fmt.Printf("Oldest:      %d days ago\n", daysAgo)
	}

	return nil
}

// init registers the status command with the root command.
func init() {
	// Add command to root
	RootCmd.AddCommand(statusCmd)
}
