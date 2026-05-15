package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// SaveContact
//
// Inserts or updates a contact for an application.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) SaveContact(ctx context.Context, contact domain.Contact) error {
	if repository.db == nil {
		return fmt.Errorf("database connection is required")
	}

	if err := contact.Validate(); err != nil {
		return err
	}

	_, err := repository.db.ExecContext(
		ctx,
		`
        INSERT INTO contacts (
            id,
            application_id,
            name,
            email,
            role
        )
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id)
        DO UPDATE SET
            application_id = EXCLUDED.application_id,
            name = EXCLUDED.name,
            email = EXCLUDED.email,
            role = EXCLUDED.role,
            updated_at = NOW()
        `,
		contact.ID,
		contact.ApplicationID,
		contact.Name,
		contact.Email,
		contact.Role,
	)

	return err
}

// -----------------------------------------------------------------------------
// FindContactByID
//
// Retrieves one contact by application and contact identity.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) FindContactByID(ctx context.Context, applicationID domain.ApplicationID, contactID domain.ContactID) (domain.Contact, error) {
	if repository.db == nil {
		return domain.Contact{}, fmt.Errorf("database connection is required")
	}

	if err := applicationID.Validate(); err != nil {
		return domain.Contact{}, err
	}

	if err := contactID.Validate(); err != nil {
		return domain.Contact{}, err
	}

	row := repository.db.QueryRowContext(
		ctx,
		`
        SELECT
            id,
            application_id,
            name,
            email,
            role
        FROM contacts
        WHERE application_id = $1
            AND id = $2
        `,
		applicationID,
		contactID,
	)

	return scanContact(row)
}

// -----------------------------------------------------------------------------
// UpdateContact
//
// Updates an existing contact for an application.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) UpdateContact(ctx context.Context, contact domain.Contact) error {
	if repository.db == nil {
		return fmt.Errorf("database connection is required")
	}

	if err := contact.Validate(); err != nil {
		return err
	}

	result, err := repository.db.ExecContext(
		ctx,
		`
        UPDATE contacts
        SET
            name = $3,
            email = $4,
            role = $5,
            updated_at = NOW()
        WHERE application_id = $1
            AND id = $2
        `,
		contact.ApplicationID,
		contact.ID,
		contact.Name,
		contact.Email,
		contact.Role,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// -----------------------------------------------------------------------------
// RemoveContact
//
// Removes one contact from an application.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) RemoveContact(ctx context.Context, applicationID domain.ApplicationID, contactID domain.ContactID) error {
	if repository.db == nil {
		return fmt.Errorf("database connection is required")
	}

	if err := applicationID.Validate(); err != nil {
		return err
	}

	if err := contactID.Validate(); err != nil {
		return err
	}

	result, err := repository.db.ExecContext(
		ctx,
		`
        DELETE FROM contacts
        WHERE application_id = $1
            AND id = $2
        `,
		applicationID,
		contactID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// -----------------------------------------------------------------------------
// ListContactsForApplication
//
// Retrieves contacts for one application in stable name order.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) ListContactsForApplication(ctx context.Context, applicationID domain.ApplicationID) ([]domain.Contact, error) {
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
            email,
            role
        FROM contacts
        WHERE application_id = $1
        ORDER BY name, id
        `,
		applicationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts := []domain.Contact{}

	for rows.Next() {
		contact, err := scanContact(rows)
		if err != nil {
			return nil, err
		}

		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

// -----------------------------------------------------------------------------
// contactScanner
//
// Defines the scan behavior shared by SQL contact rows and row results.
// -----------------------------------------------------------------------------
type contactScanner interface {
	Scan(dest ...any) error
}

// -----------------------------------------------------------------------------
// scanContact
//
// Converts a database row into a domain contact.
// -----------------------------------------------------------------------------
func scanContact(scanner contactScanner) (domain.Contact, error) {
	var contact domain.Contact
	var email sql.NullString
	var role sql.NullString

	err := scanner.Scan(
		&contact.ID,
		&contact.ApplicationID,
		&contact.Name,
		&email,
		&role,
	)
	if err != nil {
		return domain.Contact{}, err
	}

	if email.Valid {
		contact.Email = email.String
	}

	if role.Valid {
		contact.Role = role.String
	}

	if err := contact.Validate(); err != nil {
		return domain.Contact{}, err
	}

	return contact, nil
}
