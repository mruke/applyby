package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
	"github.com/mruke/applyby/internal/search"
)

// -----------------------------------------------------------------------------
// TestApplicationRepositorySearchesByStatus
//
// Verifies that PostgreSQL search can filter applications by status.
// -----------------------------------------------------------------------------
func TestApplicationRepositorySearchesByStatus(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)

	appliedApplication := newSearchRepositoryTestApplication(t, "app-001", "Backend Developer", domain.StatusApplied, "Company site")
	rejectedApplication := newSearchRepositoryTestApplication(t, "app-002", "Frontend Developer", domain.StatusRejected, "Company site")

	saveSearchRepositoryTestApplications(t, repository, appliedApplication, rejectedApplication)

	applications, err := repository.SearchApplications(context.Background(), search.ApplicationCriteria{
		Statuses: []domain.ApplicationStatus{domain.StatusApplied},
	})

	if err != nil {
		t.Fatalf("expected search to succeed: %v", err)
	}

	if len(applications) != 1 {
		t.Fatalf("expected one matching application, got %d", len(applications))
	}

	if applications[0].ID != appliedApplication.ID {
		t.Fatalf("expected applied application to match")
	}
}

// -----------------------------------------------------------------------------
// TestApplicationRepositorySearchesByCompanyAndSource
//
// Verifies that PostgreSQL search can filter by company and source.
// -----------------------------------------------------------------------------
func TestApplicationRepositorySearchesByCompanyAndSource(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)

	firstApplication := newSearchRepositoryTestApplication(t, "app-001", "Backend Developer", domain.StatusApplied, "Company site")
	secondApplication := newSearchRepositoryTestApplication(t, "app-002", "Frontend Developer", domain.StatusApplied, "Referral")
	secondApplication.Company = domain.Company{Name: "Other Studio", Website: "https://other.example.com"}

	saveSearchRepositoryTestApplications(t, repository, firstApplication, secondApplication)

	applications, err := repository.SearchApplications(context.Background(), search.ApplicationCriteria{
		CompanyName: "Example Studio",
		Source:      "Company site",
	})

	if err != nil {
		t.Fatalf("expected search to succeed: %v", err)
	}

	if len(applications) != 1 {
		t.Fatalf("expected one matching application, got %d", len(applications))
	}

	if applications[0].ID != firstApplication.ID {
		t.Fatalf("expected company and source match")
	}
}

// -----------------------------------------------------------------------------
// TestApplicationRepositorySearchesByText
//
// Verifies that PostgreSQL search can filter by text across searchable fields.
// -----------------------------------------------------------------------------
func TestApplicationRepositorySearchesByText(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)

	firstApplication := newSearchRepositoryTestApplication(t, "app-001", "Backend Developer", domain.StatusApplied, "Company site")
	firstApplication.Notes = "Backend-focused role."

	secondApplication := newSearchRepositoryTestApplication(t, "app-002", "Frontend Developer", domain.StatusApplied, "Company site")
	secondApplication.Notes = "Frontend-focused role."

	saveSearchRepositoryTestApplications(t, repository, firstApplication, secondApplication)

	applications, err := repository.SearchApplications(context.Background(), search.ApplicationCriteria{
		Text: "backend",
	})

	if err != nil {
		t.Fatalf("expected search to succeed: %v", err)
	}

	if len(applications) != 1 {
		t.Fatalf("expected one matching application, got %d", len(applications))
	}

	if applications[0].ID != firstApplication.ID {
		t.Fatalf("expected text match")
	}
}

// -----------------------------------------------------------------------------
// TestApplicationRepositorySearchesByCreatedDateRange
//
// Verifies that PostgreSQL search can filter by created date range.
// -----------------------------------------------------------------------------
func TestApplicationRepositorySearchesByCreatedDateRange(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)

	firstApplication := newSearchRepositoryTestApplication(t, "app-001", "Backend Developer", domain.StatusApplied, "Company site")
	secondApplication := newSearchRepositoryTestApplication(t, "app-002", "Frontend Developer", domain.StatusApplied, "Company site")
	secondApplication.CreatedAt = time.Date(2026, 6, 10, 8, 0, 0, 0, time.UTC)

	saveSearchRepositoryTestApplications(t, repository, firstApplication, secondApplication)

	from := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)

	applications, err := repository.SearchApplications(context.Background(), search.ApplicationCriteria{
		CreatedFrom: &from,
		CreatedTo:   &to,
	})

	if err != nil {
		t.Fatalf("expected search to succeed: %v", err)
	}

	if len(applications) != 1 {
		t.Fatalf("expected one matching application, got %d", len(applications))
	}

	if applications[0].ID != secondApplication.ID {
		t.Fatalf("expected date range match")
	}
}

// -----------------------------------------------------------------------------
// newSearchRepositoryTestApplication
//
// Creates a valid application for PostgreSQL search repository tests.
// -----------------------------------------------------------------------------
func newSearchRepositoryTestApplication(
	t *testing.T,
	id domain.ApplicationID,
	title string,
	status domain.ApplicationStatus,
	source string,
) domain.Application {
	t.Helper()

	application, err := domain.NewApplication(
		id,
		title,
		domain.Company{Name: "Example Studio", Website: "https://example.com"},
		status,
		time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to create search repository test application: %v", err)
	}

	application.Source = source

	return application
}

// -----------------------------------------------------------------------------
// saveSearchRepositoryTestApplications
//
// Saves test applications for PostgreSQL search repository tests.
// -----------------------------------------------------------------------------
func saveSearchRepositoryTestApplications(t *testing.T, repository ApplicationRepository, applications ...domain.Application) {
	t.Helper()

	for _, application := range applications {
		if err := repository.SaveApplication(context.Background(), application); err != nil {
			t.Fatalf("failed to save search repository test application: %v", err)
		}
	}
}
