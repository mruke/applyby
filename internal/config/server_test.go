package config

import "testing"

// -----------------------------------------------------------------------------
// TestLoadServerConfigUsesDefaultHTTPAddress
//
// Verifies that server config uses the default address when no environment value is set.
// -----------------------------------------------------------------------------
func TestLoadServerConfigUsesDefaultHTTPAddress(t *testing.T) {
	t.Setenv("APPLYBY_HTTP_ADDR", "")

	config := LoadServerConfig()

	if config.HTTPAddress != ":8080" {
		t.Fatalf("expected default HTTP address, got %q", config.HTTPAddress)
	}
}

// -----------------------------------------------------------------------------
// TestLoadServerConfigUsesEnvironmentHTTPAddress
//
// Verifies that server config uses the configured environment address.
// -----------------------------------------------------------------------------
func TestLoadServerConfigUsesEnvironmentHTTPAddress(t *testing.T) {
	t.Setenv("APPLYBY_HTTP_ADDR", ":9090")

	config := LoadServerConfig()

	if config.HTTPAddress != ":9090" {
		t.Fatalf("expected configured HTTP address, got %q", config.HTTPAddress)
	}
}
