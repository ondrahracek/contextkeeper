package models

// Config represents the application configuration
type Config struct {
	StoragePath string `json:"storage_path"`
	Debug       bool   `json:"debug"`
}
