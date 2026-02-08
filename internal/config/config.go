package config

import (
	"github.com/ondrahracek/contextkeeper/internal/models"
	"os"
)

var cfg *models.Config

// Load loads the configuration
func Load() (*models.Config, error) {
	cfg = &models.Config{
		StoragePath: ".contextkeeper",
		Debug:       false,
	}
	return cfg, nil
}

// Get returns the current configuration
func Get() *models.Config {
	return cfg
}

// Save saves the configuration
func Save() error {
	// Implementation
	return nil
}
