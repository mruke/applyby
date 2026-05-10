package postgres

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
)

// -----------------------------------------------------------------------------
// RunMigrations
//
// Applies the PostgreSQL schema required by the current persistence layer.
// -----------------------------------------------------------------------------
func RunMigrations(ctx context.Context, db *sql.DB) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return os.ErrInvalid
	}

	migrationPath := filepath.Join(filepath.Dir(filename), "migrations", "001_create_applications.sql")

	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, string(migrationSQL))

	return err
}
