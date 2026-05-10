package config

import "testing"

// -----------------------------------------------------------------------------
// TestLoadDatabaseConfigUsesDefaultURL
//
// Verifies that database configuration has a local development default.
// -----------------------------------------------------------------------------
func TestLoadDatabaseConfigUsesDefaultURL(t *testing.T) {
	t.Setenv("APPLYBY_DATABASE_URL", "")

	config := LoadDatabaseConfig()

	if config.URL != defaultDatabaseURL {
		t.Fatalf("expected default database URL, got %q", config.URL)
	}
}

// -----------------------------------------------------------------------------
// TestLoadDatabaseConfigUsesEnvironmentURL
//
// Verifies that database configuration can be provided through the environment.
// -----------------------------------------------------------------------------
func TestLoadDatabaseConfigUsesEnvironmentURL(t *testing.T) {
	expectedURL := "postgres://user:pass@localhost:5432/custom?sslmode=disable"
	t.Setenv("APPLYBY_DATABASE_URL", expectedURL)

	config := LoadDatabaseConfig()

	if config.URL != expectedURL {
		t.Fatalf("expected environment database URL, got %q", config.URL)
	}
}
