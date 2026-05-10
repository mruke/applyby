package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// ApplicationRepository
//
// Persists and retrieves ApplyBy applications using PostgreSQL.
// -----------------------------------------------------------------------------
type ApplicationRepository struct {
	db *sql.DB
}

// -----------------------------------------------------------------------------
// NewApplicationRepository
//
// Creates a PostgreSQL-backed application repository.
// -----------------------------------------------------------------------------
func NewApplicationRepository(db *sql.DB) ApplicationRepository {
	return ApplicationRepository{
		db: db,
	}
}

// -----------------------------------------------------------------------------
// SaveApplication
//
// Inserts or updates an application and its current company snapshot.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) SaveApplication(ctx context.Context, application domain.Application) error {
	if repository.db == nil {
		return fmt.Errorf("database connection is required")
	}

	if err := application.Validate(); err != nil {
		return err
	}

	companyID := companyIDForCompany(application.Company)

	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		`
        INSERT INTO companies (id, name, website)
        VALUES ($1, $2, $3)
        ON CONFLICT (id)
        DO UPDATE SET
            name = EXCLUDED.name,
            website = EXCLUDED.website
        `,
		companyID,
		application.Company.Name,
		application.Company.Website,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`
        INSERT INTO applications (
            id,
            title,
            company_id,
            company_name,
            company_website,
            status,
            source,
            notes,
            created_at,
            applied_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        ON CONFLICT (id)
        DO UPDATE SET
            title = EXCLUDED.title,
            company_id = EXCLUDED.company_id,
            company_name = EXCLUDED.company_name,
            company_website = EXCLUDED.company_website,
            status = EXCLUDED.status,
            source = EXCLUDED.source,
            notes = EXCLUDED.notes,
            created_at = EXCLUDED.created_at,
            applied_at = EXCLUDED.applied_at,
            updated_at = NOW()
        `,
		application.ID,
		application.Title,
		companyID,
		application.Company.Name,
		application.Company.Website,
		application.Status,
		application.Source,
		application.Notes,
		application.CreatedAt,
		application.AppliedAt,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// -----------------------------------------------------------------------------
// FindApplicationByID
//
// Retrieves one application by its stable domain identity.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) FindApplicationByID(ctx context.Context, id domain.ApplicationID) (domain.Application, error) {
	if repository.db == nil {
		return domain.Application{}, fmt.Errorf("database connection is required")
	}

	if err := id.Validate(); err != nil {
		return domain.Application{}, err
	}

	row := repository.db.QueryRowContext(
		ctx,
		`
        SELECT
            id,
            title,
            company_name,
            company_website,
            status,
            source,
            notes,
            created_at,
            applied_at
        FROM applications
        WHERE id = $1
        `,
		id,
	)

	application, err := scanApplication(row)
	if err != nil {
		return domain.Application{}, err
	}

	return application, nil
}

// -----------------------------------------------------------------------------
// ListApplications
//
// Retrieves tracked applications in stable creation order.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) ListApplications(ctx context.Context) ([]domain.Application, error) {
	if repository.db == nil {
		return nil, fmt.Errorf("database connection is required")
	}

	rows, err := repository.db.QueryContext(
		ctx,
		`
        SELECT
            id,
            title,
            company_name,
            company_website,
            status,
            source,
            notes,
            created_at,
            applied_at
        FROM applications
        ORDER BY created_at, id
        `,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applications := []domain.Application{}

	for rows.Next() {
		application, err := scanApplication(rows)
		if err != nil {
			return nil, err
		}

		applications = append(applications, application)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return applications, nil
}

// -----------------------------------------------------------------------------
// applicationScanner
//
// Defines the scan behavior shared by SQL rows and row results.
// -----------------------------------------------------------------------------
type applicationScanner interface {
	Scan(dest ...any) error
}

// -----------------------------------------------------------------------------
// scanApplication
//
// Converts a database row into a domain application.
// -----------------------------------------------------------------------------
func scanApplication(scanner applicationScanner) (domain.Application, error) {
	var id domain.ApplicationID
	var title string
	var companyName string
	var companyWebsite string
	var status domain.ApplicationStatus
	var source string
	var notes string
	var appliedAt sql.NullTime

	application := domain.Application{}

	err := scanner.Scan(
		&id,
		&title,
		&companyName,
		&companyWebsite,
		&status,
		&source,
		&notes,
		&application.CreatedAt,
		&appliedAt,
	)
	if err != nil {
		return domain.Application{}, err
	}

	application.ID = id
	application.Title = title
	application.Company = domain.Company{
		Name:    companyName,
		Website: companyWebsite,
	}
	application.Status = status
	application.Source = source
	application.Notes = notes

	if appliedAt.Valid {
		application.AppliedAt = &appliedAt.Time
	}

	if err := application.Validate(); err != nil {
		return domain.Application{}, err
	}

	return application, nil
}

// -----------------------------------------------------------------------------
// companyIDForCompany
//
// Creates a stable storage identity for the current company value object.
// -----------------------------------------------------------------------------
func companyIDForCompany(company domain.Company) string {
	normalizedName := strings.ToLower(strings.TrimSpace(company.Name))
	hash := sha256.Sum256([]byte(normalizedName))

	return "company-" + hex.EncodeToString(hash[:8])
}
