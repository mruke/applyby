package search

import (
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestApplicationCriteriaValidateAcceptsValidCriteria
//
// Verifies that valid application search criteria are accepted.
// -----------------------------------------------------------------------------
func TestApplicationCriteriaValidateAcceptsValidCriteria(t *testing.T) {
	from := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC)

	criteria := ApplicationCriteria{
		Statuses:    []domain.ApplicationStatus{domain.StatusApplied},
		CompanyName: "Example Studio",
		Source:      "Company site",
		Text:        "backend",
		CreatedFrom: &from,
		CreatedTo:   &to,
	}

	if err := criteria.Validate(); err != nil {
		t.Fatalf("expected criteria to be valid: %v", err)
	}
}

// -----------------------------------------------------------------------------
// TestApplicationCriteriaValidateRejectsInvalidStatus
//
// Verifies that unsupported status filters are rejected.
// -----------------------------------------------------------------------------
func TestApplicationCriteriaValidateRejectsInvalidStatus(t *testing.T) {
	criteria := ApplicationCriteria{
		Statuses: []domain.ApplicationStatus{domain.ApplicationStatus("paused")},
	}

	if err := criteria.Validate(); err == nil {
		t.Fatal("expected invalid status criteria to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestApplicationCriteriaValidateRejectsInvalidDateRange
//
// Verifies that created_to cannot come before created_from.
// -----------------------------------------------------------------------------
func TestApplicationCriteriaValidateRejectsInvalidDateRange(t *testing.T) {
	from := time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)

	criteria := ApplicationCriteria{
		CreatedFrom: &from,
		CreatedTo:   &to,
	}

	if err := criteria.Validate(); err == nil {
		t.Fatal("expected invalid date range to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestMatchesApplicationAcceptsMatchingApplication
//
// Verifies that an application matching all criteria is accepted.
// -----------------------------------------------------------------------------
func TestMatchesApplicationAcceptsMatchingApplication(t *testing.T) {
	application := newSearchTestApplication(t)
	from := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC)

	criteria := ApplicationCriteria{
		Statuses:    []domain.ApplicationStatus{domain.StatusApplied},
		CompanyName: "Example Studio",
		Source:      "Company site",
		Text:        "backend",
		CreatedFrom: &from,
		CreatedTo:   &to,
	}

	if !MatchesApplication(application, criteria) {
		t.Fatal("expected application to match criteria")
	}
}

// -----------------------------------------------------------------------------
// TestMatchesApplicationRejectsNonMatchingApplication
//
// Verifies that an application failing criteria is rejected.
// -----------------------------------------------------------------------------
func TestMatchesApplicationRejectsNonMatchingApplication(t *testing.T) {
	application := newSearchTestApplication(t)

	criteria := ApplicationCriteria{
		Statuses: []domain.ApplicationStatus{domain.StatusRejected},
	}

	if MatchesApplication(application, criteria) {
		t.Fatal("expected application not to match criteria")
	}
}

// -----------------------------------------------------------------------------
// newSearchTestApplication
//
// Creates a valid application for search criteria tests.
// -----------------------------------------------------------------------------
func newSearchTestApplication(t *testing.T) domain.Application {
	t.Helper()

	application, err := domain.NewApplication(
		"app-001",
		"Backend Developer",
		domain.Company{Name: "Example Studio", Website: "https://example.com"},
		domain.StatusApplied,
		time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to create test application: %v", err)
	}

	application.Source = "Company site"
	application.Notes = "Backend-focused role."

	return application
}
