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
// FindDocumentByID
//
// Retrieves one document metadata record by application and document identity.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) FindDocumentByID(ctx context.Context, applicationID domain.ApplicationID, documentID domain.DocumentID) (domain.Document, error) {
	if repository.db == nil {
		return domain.Document{}, fmt.Errorf("database connection is required")
	}

	if err := applicationID.Validate(); err != nil {
		return domain.Document{}, err
	}

	if err := documentID.Validate(); err != nil {
		return domain.Document{}, err
	}

	row := repository.db.QueryRowContext(
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
          AND id = $2
        `,
		applicationID,
		documentID,
	)

	document, err := scanDocument(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Document{}, fmt.Errorf("document not found: %s", documentID)
		}

		return domain.Document{}, err
	}

	return document, nil
}

// -----------------------------------------------------------------------------
// UpdateDocument
//
// Updates an existing document metadata record for an application.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) UpdateDocument(ctx context.Context, document domain.Document) error {
	if repository.db == nil {
		return fmt.Errorf("database connection is required")
	}

	if err := document.Validate(); err != nil {
		return err
	}

	result, err := repository.db.ExecContext(
		ctx,
		`
        UPDATE documents
        SET
            name = $3,
            kind = $4,
            path = $5,
            updated_at = NOW()
        WHERE application_id = $1
          AND id = $2
        `,
		document.ApplicationID,
		document.ID,
		document.Name,
		document.Kind,
		document.Path,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("document not found: %s", document.ID)
	}

	return nil
}

// -----------------------------------------------------------------------------
// RemoveDocument
//
// Removes one document metadata record from an application.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) RemoveDocument(ctx context.Context, applicationID domain.ApplicationID, documentID domain.DocumentID) error {
	if repository.db == nil {
		return fmt.Errorf("database connection is required")
	}

	if err := applicationID.Validate(); err != nil {
		return err
	}

	if err := documentID.Validate(); err != nil {
		return err
	}

	result, err := repository.db.ExecContext(
		ctx,
		`
        DELETE FROM documents
        WHERE application_id = $1
          AND id = $2
        `,
		applicationID,
		documentID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("document not found: %s", documentID)
	}

	return nil
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
