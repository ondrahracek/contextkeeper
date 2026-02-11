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
	Use:   "list",
	Short: "List context items",
	Long:  "List context items, optionally filtered by project or tags. Use --all to include completed items.",
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
	// Initialize storage and load items
	stor := storage.NewStorage(config.FindStoragePath(pathFlag))
	if err := stor.Load(); err != nil {
		return fmt.Errorf("failed to load storage: %w", err)
	}

	// Get all items
	items := stor.GetAll()

	// Filter by project if specified
	if projectFilter != "" {
		items = filterByProject(items, projectFilter)
	}

	// Filter by tags if specified
	if tagFilter != "" {
		items = filterByTags(items, tagFilter)
	}

	// Filter out completed items unless --all is set
	if !showAll {
		items = filterActive(items)
	}

	// Output in requested format
	if jsonOutput {
		type jsonItem struct {
			ID          string     `json:"id"`
			FullID      string     `json:"fullId"`
			Content     string     `json:"content"`
			Project     string     `json:"project"`
			Tags        []string   `json:"tags"`
			CompletedAt *time.Time `json:"completedAt"`
			CreatedAt   time.Time  `json:"createdAt"`
		}

		jsonItems := make([]jsonItem, 0, len(items))
		for _, item := range items {
			jsonItems = append(jsonItems, jsonItem{
				ID:          item.ID[:8],
				FullID:      item.ID,
				Content:     item.Content,
				Project:     item.Project,
				Tags:        item.Tags,
				CompletedAt: item.CompletedAt,
				CreatedAt:   item.CreatedAt,
			})
		}

		data, err := json.MarshalIndent(jsonItems, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal items to JSON: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
	} else {
		fmt.Fprint(cmd.OutOrStdout(), utils.FormatItemList(items, showAll))
	}

	return nil
}

// filterByProject filters items by the specified project name.
func filterByProject(items []models.ContextItem, project string) []models.ContextItem {
	filtered := make([]models.ContextItem, 0)
	for _, item := range items {
		if item.Project == project {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// filterByTags filters items by the specified tags.
func filterByTags(items []models.ContextItem, tags string) []models.ContextItem {
	filterTags := utils.ParseTags(tags)
	filtered := make([]models.ContextItem, 0)
	for _, item := range items {
		if containsTags(item.Tags, filterTags) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// filterActive filters out completed items.
func filterActive(items []models.ContextItem) []models.ContextItem {
	active := make([]models.ContextItem, 0)
	for _, item := range items {
		if item.CompletedAt == nil {
			active = append(active, item)
		}
	}
	return active
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
