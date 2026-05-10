package config

import (
	"os"
	"strings"
)

const defaultDatabaseURL = "postgres://applyby:applyby@localhost:5432/applyby?sslmode=disable"

// -----------------------------------------------------------------------------
// DatabaseConfig
//
// Contains the database settings required by persistence code.
// -----------------------------------------------------------------------------
type DatabaseConfig struct {
	URL string
}

// -----------------------------------------------------------------------------
// LoadDatabaseConfig
//
// Loads database configuration from the environment with development defaults.
// -----------------------------------------------------------------------------
func LoadDatabaseConfig() DatabaseConfig {
	url := strings.TrimSpace(os.Getenv("APPLYBY_DATABASE_URL"))

	if url == "" {
		url = defaultDatabaseURL
	}

	return DatabaseConfig{
		URL: url,
	}
}
