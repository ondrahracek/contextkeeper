// Package cli provides the command-line interface for ContextKeeper.
//
// This package implements the Cobra-based CLI for managing project context
// and configuration. See the root.go file for the main command structure.
package cli

import (
	"fmt"

	"github.com/ondrahracek/contextkeeper/internal/config"
	"github.com/ondrahracek/contextkeeper/internal/models"
	"github.com/spf13/cobra"
)

// configCmd manages ContextKeeper configuration.
//
// The command supports showing current settings, getting specific values,
// setting values, and resetting to defaults.
var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Manage configuration",
	Long:    "Show, get, set, or reset configuration values.",
	Example: `  # Show current configuration
  ck config --show

  # Get a specific value
  ck config --get defaultProject

  # Set a value
  ck config --set defaultProject "my-project"

  # Reset to defaults
  ck config --reset`,
	Args: cobra.NoArgs,
	RunE: configCommand,
}

// Command flags for the config command.
var (
	showConfig     bool
	resetConfig    bool
	getConfigKey   string
	setConfigKey   string
	setConfigValue string
)

// configCommand is the execution function for the config command.
// It handles configuration display, retrieval, modification, and reset.
func configCommand(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Handle flags
	switch {
	case showConfig:
		return configShowValues(cfg)
	case resetConfig:
		return configResetValues(cfg)
	case getConfigKey != "":
		return configGetValue(cfg, getConfigKey)
	case setConfigKey != "" && setConfigValue != "":
		return configSetValue(cfg, setConfigKey, setConfigValue)
	default:
		return configShowValues(cfg)
	}
}

// configShowValues prints all configuration values.
func configShowValues(cfg *models.Config) error {
	fmt.Println("Current Configuration:")
	fmt.Printf("  StoragePath:    %s\n", cfg.StoragePath)
	fmt.Printf("  DefaultProject: %s\n", cfg.DefaultProject)
	fmt.Printf("  DateFormat:     %s\n", cfg.DateFormat)
	fmt.Printf("  Editor:         %s\n", cfg.Editor)
	return nil
}

// configGetValue prints a specific configuration value.
func configGetValue(cfg *models.Config, key string) error {
	switch key {
	case "storagePath":
		fmt.Println(cfg.StoragePath)
	case "defaultProject":
		fmt.Println(cfg.DefaultProject)
	case "dateFormat":
		fmt.Println(cfg.DateFormat)
	case "editor":
		fmt.Println(cfg.Editor)
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	return nil
}

// configSetValue modifies a configuration value.
func configSetValue(cfg *models.Config, key, value string) error {
	switch key {
	case "storagePath":
		cfg.StoragePath = value
	case "defaultProject":
		cfg.DefaultProject = value
	case "dateFormat":
		cfg.DateFormat = value
	case "editor":
		cfg.Editor = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	if err := config.Save(); err != nil {
		return err
	}

	fmt.Printf("Set %s to: %s\n", key, value)
	return nil
}

// configResetValues restores configuration to defaults.
func configResetValues(cfg *models.Config) error {
	cfg.DefaultProject = ""
	cfg.DateFormat = "2006-01-02 15:04"
	cfg.Editor = ""

	if err := config.Save(); err != nil {
		return err
	}

	fmt.Println("Configuration reset to defaults.")
	return nil
}

// init registers the config command with the root command.
func init() {
	// Register command flags
	configCmd.Flags().BoolVar(&showConfig, "show", false, "Show current configuration")
	configCmd.Flags().BoolVar(&resetConfig, "reset", false, "Reset configuration to defaults")
	configCmd.Flags().StringVar(&getConfigKey, "get", "", "Get a specific configuration value")
	configCmd.Flags().StringVar(&setConfigKey, "set", "", "Set a configuration key")
	configCmd.Flags().StringVar(&setConfigValue, "value", "", "The value to set")

	// Add command to root
	RootCmd.AddCommand(configCmd)
}
