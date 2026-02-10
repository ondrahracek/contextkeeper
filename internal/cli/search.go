// Package cli provides the command-line interface for ContextKeeper.
//
// This package implements the Cobra-based CLI for managing context and
// configuration. See the root.go file for the main command structure.
package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ondrahracek/contextkeeper/internal/config"
	"github.com/ondrahracek/contextkeeper/internal/models"
	"github.com/ondrahracek/contextkeeper/internal/storage"
	"github.com/ondrahracek/contextkeeper/internal/utils"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command for finding context items.
//
// The search command allows users to find context items by searching through
// their content and tags. Results can be filtered by various criteria and
// output in either human-readable or JSON format.
//
// # Usage
//
//	ck search <query>              Search by content
//	ck search --tag <tag>          Filter by tag
//	ck search --all                Include completed items
//	ck search --json               Output as JSON
//
// # Examples
//
//	# Search for items containing "auth"
//	ck search auth
//
//	# Search with JSON output for scripting
//	ck search auth --json
//
//	# Search by tag
//	ck search --tag bug
//
//	# Include completed items
//	ck search --all dashboard
//
// # Exit Codes
//
//	0 - Success
//	1 - Storage error or other failure
var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search context items",
	Long: `Search context items by content or tags.

The search is case-insensitive and matches against both the item content
and its tags. Use --all to include completed items in the results.
If no query is provided, returns all active (non-completed) items.`,
	Example: `  # Search for items containing "auth"
  ck search auth

  # Search with JSON output (for scripting)
  ck search auth --json

  # Search by tag
  ck search --tag bug

  # Include completed items
  ck search --all dashboard

  # List all active items (no query)
  ck search`,
	Args: cobra.ArbitraryArgs,
	RunE: runSearch,
}

// searchFlags holds the command-line flags for the search command.
// These are package-level variables to be set by Cobra during flag parsing.
var (
	searchTagFilter  string // -t, --tag: Filter by specific tags
	searchShowAll   bool   // -a, --all: Include completed items
	searchJsonOut   bool   // --json: Output as JSON
)

// searchResult represents the JSON structure returned by search --json.
type searchResult struct {
	ID          string     `json:"id"`           // 8-character ID prefix
	FullID      string     `json:"fullId"`       // Full UUID
	Content     string     `json:"content"`      // Item content
	Project     string     `json:"project"`      // Project name
	Tags        []string   `json:"tags"`         // Associated tags
	CompletedAt *time.Time `json:"completedAt"`  // Completion timestamp or nil
	CreatedAt   time.Time  `json:"createdAt"`   // Creation timestamp
}

// runSearch is the main execution function for the search command.
// It loads items from storage, applies filters, and outputs results.
func runSearch(cmd *cobra.Command, args []string) error {
	query := ""
	if len(args) > 0 {
		query = args[0]
	}

	stor := storage.NewStorage(config.FindStoragePath(""))
	if err := stor.Load(); err != nil {
		return fmt.Errorf("failed to load storage: %w", err)
	}

	items := stor.GetAll()
	items = applySearchFilters(items, query, searchTagFilter, searchShowAll)

	if searchJsonOut {
		return outputSearchJSON(cmd, items)
	}
	return outputSearchText(cmd, items, searchShowAll)
}

// applySearchFilters applies all search filters to the items slice.
// The order of operations is: completed filter -> tag filter -> query filter.
func applySearchFilters(items []models.ContextItem, query, tagFilter string, showAll bool) []models.ContextItem {
	// Filter completed items first (most restrictive)
	if !showAll {
		items = filterActive(items)
	}

	// Filter by tag if specified
	if tagFilter != "" {
		items = filterByTags(items, tagFilter)
	}

	// Filter by query if specified
	if query != "" {
		items = filterByQuery(items, query)
	}

	return items
}

// outputSearchJSON outputs search results as formatted JSON.
func outputSearchJSON(cmd *cobra.Command, items []models.ContextItem) error {
	results := make([]searchResult, 0, len(items))
	for _, item := range items {
		results = append(results, searchResult{
			ID:          item.ID[:8],
			FullID:      item.ID,
			Content:     item.Content,
			Project:     item.Project,
			Tags:        item.Tags,
			CompletedAt: item.CompletedAt,
			CreatedAt:   item.CreatedAt,
		})
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return nil
}

// outputSearchText outputs search results in human-readable format.
func outputSearchText(cmd *cobra.Command, items []models.ContextItem, showAll bool) error {
	fmt.Fprint(cmd.OutOrStdout(), utils.FormatItemList(items, showAll))
	return nil
}

// filterByQuery filters items by matching the query string
// against the content and tags (case-insensitive).
// Returns a new slice containing only matching items.
func filterByQuery(items []models.ContextItem, query string) []models.ContextItem {
	if query == "" {
		return items
	}

	lowerQuery := strings.ToLower(query)
	filtered := make([]models.ContextItem, 0, len(items))

	for _, item := range items {
		if matchesQuery(item, lowerQuery) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// matchesQuery checks if an item matches the query string
// by searching in both content and tags.
func matchesQuery(item models.ContextItem, query string) bool {
	// Check content first (most common case)
	if strings.Contains(strings.ToLower(item.Content), query) {
		return true
	}
	// Check tags
	for _, tag := range item.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

// init registers the search command with the root command.
func init() {
	// Register command flags with appropriate descriptions
	searchCmd.Flags().StringVarP(&searchTagFilter, "tag", "t", "",
		"Filter by tags (comma or space separated)")
	searchCmd.Flags().BoolVarP(&searchShowAll, "all", "a", false,
		"Include completed items in results")
	searchCmd.Flags().BoolVar(&searchJsonOut, "json", false,
		"Output results as JSON")

	// Add command to root
	RootCmd.AddCommand(searchCmd)
}