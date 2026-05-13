package config

import "os"

const defaultHTTPAddress = ":8080"

// -----------------------------------------------------------------------------
// ServerConfig
//
// Contains runtime settings for the HTTP API server.
// -----------------------------------------------------------------------------
type ServerConfig struct {
	HTTPAddress string
}

// -----------------------------------------------------------------------------
// LoadServerConfig
//
// Reads HTTP server configuration from environment variables.
// -----------------------------------------------------------------------------
func LoadServerConfig() ServerConfig {
	httpAddress := os.Getenv("APPLYBY_HTTP_ADDR")
	if httpAddress == "" {
		httpAddress = defaultHTTPAddress
	}

	return ServerConfig{
		HTTPAddress: httpAddress,
	}
}
