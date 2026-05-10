package postgres

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"sort"
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

	migrationDir := filepath.Join(filepath.Dir(filename), "migrations")

	migrationPaths, err := filepath.Glob(filepath.Join(migrationDir, "*.sql"))
	if err != nil {
		return err
	}

	sort.Strings(migrationPaths)

	for _, migrationPath := range migrationPaths {
		migrationSQL, err := os.ReadFile(migrationPath)
		if err != nil {
			return err
		}

		if _, err := db.ExecContext(ctx, string(migrationSQL)); err != nil {
			return err
		}
	}

	return nil
}
