package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// SaveDocument
//
// Inserts or updates document metadata for an application.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) SaveDocument(ctx context.Context, document domain.Document) error {
	if repository.db == nil {
		return fmt.Errorf("database connection is required")
	}

	if err := document.Validate(); err != nil {
		return err
	}

	_, err := repository.db.ExecContext(
		ctx,
		`
        INSERT INTO documents (
            id,
            application_id,
            name,
            kind,
            path
        )
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id)
        DO UPDATE SET
            application_id = EXCLUDED.application_id,
            name = EXCLUDED.name,
            kind = EXCLUDED.kind,
            path = EXCLUDED.path,
            updated_at = NOW()
        `,
		document.ID,
		document.ApplicationID,
		document.Name,
		document.Kind,
		document.Path,
	)

	return err
}

// -----------------------------------------------------------------------------
// ListDocumentsForApplication
//
// Retrieves document metadata for one application in stable name order.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) ListDocumentsForApplication(ctx context.Context, applicationID domain.ApplicationID) ([]domain.Document, error) {
	if repository.db == nil {
		return nil, fmt.Errorf("database connection is required")
	}

	if err := applicationID.Validate(); err != nil {
		return nil, err
	}

	rows, err := repository.db.QueryContext(
		ctx,
		`
        SELECT
            id,
            application_id,
            name,
            kind,
            path
        FROM documents
        WHERE application_id = $1
        ORDER BY name, id
        `,
		applicationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documents := []domain.Document{}

	for rows.Next() {
		document, err := scanDocument(rows)
		if err != nil {
			return nil, err
		}

		documents = append(documents, document)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return documents, nil
}

// -----------------------------------------------------------------------------
// documentScanner
//
// Defines the scan behavior shared by SQL document rows and row results.
// -----------------------------------------------------------------------------
type documentScanner interface {
	Scan(dest ...any) error
}

// -----------------------------------------------------------------------------
// scanDocument
//
// Converts a database row into domain document metadata.
// -----------------------------------------------------------------------------
func scanDocument(scanner documentScanner) (domain.Document, error) {
	var document domain.Document
	var path sql.NullString

	err := scanner.Scan(
		&document.ID,
		&document.ApplicationID,
		&document.Name,
		&document.Kind,
		&path,
	)
	if err != nil {
		return domain.Document{}, err
	}

	if path.Valid {
		document.Path = path.String
	}

	if err := document.Validate(); err != nil {
		return domain.Document{}, err
	}

	return document, nil
}
