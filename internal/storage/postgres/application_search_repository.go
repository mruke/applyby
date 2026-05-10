package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/mruke/applyby/internal/domain"
	"github.com/mruke/applyby/internal/search"
)

// -----------------------------------------------------------------------------
// SearchApplications
//
// Retrieves applications matching explicit search and filter criteria.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) SearchApplications(ctx context.Context, criteria search.ApplicationCriteria) ([]domain.Application, error) {
	if repository.db == nil {
		return nil, fmt.Errorf("database connection is required")
	}

	criteria = criteria.Normalize()

	if err := criteria.Validate(); err != nil {
		return nil, err
	}

	query := `
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
    `

	whereClauses := []string{}
	args := []any{}

	if len(criteria.Statuses) > 0 {
		statusPlaceholders := []string{}

		for _, status := range criteria.Statuses {
			args = append(args, status)
			statusPlaceholders = append(statusPlaceholders, fmt.Sprintf("$%d", len(args)))
		}

		whereClauses = append(whereClauses, "status IN ("+strings.Join(statusPlaceholders, ", ")+")")
	}

	if criteria.CompanyName != "" {
		args = append(args, strings.ToLower(criteria.CompanyName))
		whereClauses = append(whereClauses, fmt.Sprintf("LOWER(company_name) = $%d", len(args)))
	}

	if criteria.Source != "" {
		args = append(args, strings.ToLower(criteria.Source))
		whereClauses = append(whereClauses, fmt.Sprintf("LOWER(source) = $%d", len(args)))
	}

	if criteria.CreatedFrom != nil {
		args = append(args, *criteria.CreatedFrom)
		whereClauses = append(whereClauses, fmt.Sprintf("created_at >= $%d", len(args)))
	}

	if criteria.CreatedTo != nil {
		args = append(args, *criteria.CreatedTo)
		whereClauses = append(whereClauses, fmt.Sprintf("created_at <= $%d", len(args)))
	}

	if criteria.Text != "" {
		args = append(args, "%"+strings.ToLower(criteria.Text)+"%")
		placeholder := fmt.Sprintf("$%d", len(args))
		whereClauses = append(whereClauses, "(LOWER(title) LIKE "+placeholder+" OR LOWER(company_name) LIKE "+placeholder+" OR LOWER(source) LIKE "+placeholder+" OR LOWER(notes) LIKE "+placeholder+")")
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " ORDER BY created_at, id"

	rows, err := repository.db.QueryContext(ctx, query, args...)
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
