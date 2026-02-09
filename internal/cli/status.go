// Package cli provides the command-line interface for ContextKeeper.
//
// This package implements the Cobra-based CLI for managing context and
// configuration. See the root.go file for the main command structure.
package cli

import (
	"encoding/json"
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
	Use:   "status",
	Short: "Show quick overview",
	Long:  "Show a quick overview of context items including counts and storage information.",
	Example: `  # Show status overview
  ck status`,
	Args: cobra.NoArgs,
	RunE: statusCommand,
}

// statusCommand is the execution function for the status command.
// It gathers and displays statistics about stored items.
func statusCommand(cmd *cobra.Command, args []string) error {
	// Get storage path
	storagePath := config.FindStoragePath("")

	// Initialize storage and load items
	stor := storage.NewStorage(storagePath)
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
	projectsMap := make(map[string]bool)
	tagsMap := make(map[string]bool)

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

		if item.Project != "" {
			projectsMap[item.Project] = true
		}
		for _, tag := range item.Tags {
			tagsMap[tag] = true
		}
	}

	active := total - completed

	if jsonOutput {
		projects := []string{}
		for p := range projectsMap {
			projects = append(projects, p)
		}
		tags := []string{}
		for t := range tagsMap {
			tags = append(tags, t)
		}

		status := map[string]interface{}{
			"totalItems":     total,
			"completedItems": completed,
			"activeItems":    active,
			"projects":       projects,
			"tags":           tags,
		}
		data, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal status to JSON: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	}

	// Print status
	fmt.Fprintln(cmd.OutOrStdout(), "ContextKeeper Status")
	fmt.Fprintln(cmd.OutOrStdout(), "===================")
	fmt.Fprintf(cmd.OutOrStdout(), "Storage Path: %s\n", storagePath)
	fmt.Fprintf(cmd.OutOrStdout(), "Total Items: %d\n", total)
	fmt.Fprintf(cmd.OutOrStdout(), "Active:      %d\n", active)
	fmt.Fprintf(cmd.OutOrStdout(), "Completed:   %d\n", completed)

	if oldestSet {
		daysAgo := int(time.Since(oldest).Hours() / 24)
		fmt.Fprintf(cmd.OutOrStdout(), "Oldest:      %d days ago\n", daysAgo)
	}

	return nil
}

// init registers the status command with the root command.
func init() {
	statusCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	// Add command to root
	RootCmd.AddCommand(statusCmd)
}
