package search

import (
	"fmt"
	"strings"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// ApplicationCriteria
//
// Represents search and filter criteria for tracked applications.
// -----------------------------------------------------------------------------
type ApplicationCriteria struct {
	Statuses    []domain.ApplicationStatus
	CompanyName string
	Source      string
	Text        string
	CreatedFrom *time.Time
	CreatedTo   *time.Time
}

// -----------------------------------------------------------------------------
// Normalize
//
// Returns a trimmed copy of the application search criteria.
// -----------------------------------------------------------------------------
func (criteria ApplicationCriteria) Normalize() ApplicationCriteria {
	return ApplicationCriteria{
		Statuses:    criteria.Statuses,
		CompanyName: strings.TrimSpace(criteria.CompanyName),
		Source:      strings.TrimSpace(criteria.Source),
		Text:        strings.TrimSpace(criteria.Text),
		CreatedFrom: criteria.CreatedFrom,
		CreatedTo:   criteria.CreatedTo,
	}
}

// -----------------------------------------------------------------------------
// Validate
//
// Verifies that application search criteria contain supported values.
// -----------------------------------------------------------------------------
func (criteria ApplicationCriteria) Validate() error {
	for _, status := range criteria.Statuses {
		if !domain.IsValidApplicationStatus(status) {
			return fmt.Errorf("invalid application search status: %q", status)
		}
	}

	if criteria.CreatedFrom != nil && criteria.CreatedTo != nil && criteria.CreatedTo.Before(*criteria.CreatedFrom) {
		return fmt.Errorf("created_to cannot be before created_from")
	}

	return nil
}

// -----------------------------------------------------------------------------
// MatchesApplication
//
// Reports whether an application matches the provided search criteria.
// -----------------------------------------------------------------------------
func MatchesApplication(application domain.Application, criteria ApplicationCriteria) bool {
	criteria = criteria.Normalize()

	if len(criteria.Statuses) > 0 && !matchesAnyStatus(application.Status, criteria.Statuses) {
		return false
	}

	if criteria.CompanyName != "" && !strings.EqualFold(application.Company.Name, criteria.CompanyName) {
		return false
	}

	if criteria.Source != "" && !strings.EqualFold(application.Source, criteria.Source) {
		return false
	}

	if criteria.CreatedFrom != nil && application.CreatedAt.Before(*criteria.CreatedFrom) {
		return false
	}

	if criteria.CreatedTo != nil && application.CreatedAt.After(*criteria.CreatedTo) {
		return false
	}

	if criteria.Text != "" && !containsApplicationText(application, criteria.Text) {
		return false
	}

	return true
}

// -----------------------------------------------------------------------------
// matchesAnyStatus
//
// Reports whether a status appears in a list of allowed statuses.
// -----------------------------------------------------------------------------
func matchesAnyStatus(status domain.ApplicationStatus, statuses []domain.ApplicationStatus) bool {
	for _, allowedStatus := range statuses {
		if status == allowedStatus {
			return true
		}
	}

	return false
}

// -----------------------------------------------------------------------------
// containsApplicationText
//
// Reports whether searchable application fields contain the search text.
// -----------------------------------------------------------------------------
func containsApplicationText(application domain.Application, text string) bool {
	normalizedText := strings.ToLower(strings.TrimSpace(text))

	return strings.Contains(strings.ToLower(application.Title), normalizedText) ||
		strings.Contains(strings.ToLower(application.Company.Name), normalizedText) ||
		strings.Contains(strings.ToLower(application.Source), normalizedText) ||
		strings.Contains(strings.ToLower(application.Notes), normalizedText)
}
