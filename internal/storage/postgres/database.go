package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/mruke/applyby/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// -----------------------------------------------------------------------------
// OpenDatabase
//
// Opens and verifies a PostgreSQL database connection.
// -----------------------------------------------------------------------------
func OpenDatabase(ctx context.Context, databaseConfig config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseConfig.URL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
