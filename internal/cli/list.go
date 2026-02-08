// Package cli provides the command-line interface for ContextKeeper.
//
// This package implements the Cobra-based CLI for managing project context
// and configuration. See the root.go file for the main command structure.
package cli

import (
	"encoding/json"
	"fmt"

	"github.com/ondrahracek/contextkeeper/internal/config"
	"github.com/ondrahracek/contextkeeper/internal/models"
	"github.com/ondrahracek/contextkeeper/internal/storage"
	"github.com/ondrahracek/contextkeeper/internal/utils"
	"github.com/spf13/cobra"
)

// listCmd lists context items from storage.
//
// The command supports filtering by project, tags, and completion status.
// Output can be displayed in a formatted table or as JSON.
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List context items",
	Long:    "List context items, optionally filtered by project or tags. Use --all to include completed items.",
	Example: `  # List all active items
  ck list

  # List items for a specific project
  ck list --project "my-project"

  # List items with specific tags
  ck list --tags "bug,urgent"

  # Include completed items
  ck list --all

  # Output as JSON
  ck list --json`,
	Args: cobra.NoArgs,
	RunE: listCommand,
}

// Command flags for the list command.
var (
	projectFilter string
	tagFilter     string
	showAll       bool
	jsonOutput    bool
)

// listCommand is the execution function for the list command.
// It retrieves and filters context items from storage.
func listCommand(cmd *cobra.Command, args []string) error {
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
	items := stor.GetAll()

	// Filter by project if specified
	if projectFilter != "" {
		filtered := make([]models.ContextItem, 0)
		for _, item := range items {
			if item.Project == projectFilter {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}

	// Filter by tags if specified
	if tagFilter != "" {
		filterTags := utils.ParseTags(tagFilter)
		filtered := make([]models.ContextItem, 0)
		for _, item := range items {
			if containsTags(item.Tags, filterTags) {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}

	// Filter out completed items unless --all is set
	if !showAll {
		active := make([]models.ContextItem, 0)
		for _, item := range items {
			if item.CompletedAt == nil {
				active = append(active, item)
			}
		}
		items = active
	}

	// Output in requested format
	if jsonOutput {
		data, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	} else {
		fmt.Print(utils.FormatItemList(items, showAll))
	}

	return nil
}

// containsTags checks if itemTags contains all filterTags.
func containsTags(itemTags, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}
	for _, ft := range filterTags {
		found := false
		for _, it := range itemTags {
			if it == ft {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// init registers the list command with the root command.
func init() {
	// Register command flags
	listCmd.Flags().StringVarP(&projectFilter, "project", "P", "", "Filter by project name")
	listCmd.Flags().StringVarP(&tagFilter, "tags", "t", "", "Filter by tags (comma or space separated)")
	listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all items including completed")
	listCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	// Add command to root
	RootCmd.AddCommand(listCmd)
}
