package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/config"
	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestApplicationRepositorySavesAndFindsApplication
//
// Verifies that PostgreSQL persistence can save and retrieve an application.
// -----------------------------------------------------------------------------
func TestApplicationRepositorySavesAndFindsApplication(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)

	application := newTestApplication(t, "app-001", "Backend Developer", domain.StatusApplied)

	if err := repository.SaveApplication(context.Background(), application); err != nil {
		t.Fatalf("expected application save to succeed: %v", err)
	}

	foundApplication, err := repository.FindApplicationByID(context.Background(), application.ID)
	if err != nil {
		t.Fatalf("expected application lookup to succeed: %v", err)
	}

	if foundApplication.ID != application.ID {
		t.Fatalf("expected application id %q, got %q", application.ID, foundApplication.ID)
	}

	if foundApplication.Company.Name != application.Company.Name {
		t.Fatalf("expected company name %q, got %q", application.Company.Name, foundApplication.Company.Name)
	}
}

// -----------------------------------------------------------------------------
// TestApplicationRepositoryUpdatesExistingApplication
//
// Verifies that saving an existing application updates current state.
// -----------------------------------------------------------------------------
func TestApplicationRepositoryUpdatesExistingApplication(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)

	application := newTestApplication(t, "app-001", "Backend Developer", domain.StatusApplied)

	if err := repository.SaveApplication(context.Background(), application); err != nil {
		t.Fatalf("expected initial save to succeed: %v", err)
	}

	application.Status = domain.StatusInterviewing
	application.Notes = "Moved to interview stage."

	if err := repository.SaveApplication(context.Background(), application); err != nil {
		t.Fatalf("expected update save to succeed: %v", err)
	}

	foundApplication, err := repository.FindApplicationByID(context.Background(), application.ID)
	if err != nil {
		t.Fatalf("expected application lookup to succeed: %v", err)
	}

	if foundApplication.Status != domain.StatusInterviewing {
		t.Fatalf("expected updated status %q, got %q", domain.StatusInterviewing, foundApplication.Status)
	}

	if foundApplication.Notes != "Moved to interview stage." {
		t.Fatalf("expected updated notes to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestApplicationRepositoryListsApplications
//
// Verifies that PostgreSQL persistence can list tracked applications.
// -----------------------------------------------------------------------------
func TestApplicationRepositoryListsApplications(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)

	firstApplication := newTestApplication(t, "app-001", "Backend Developer", domain.StatusApplied)
	secondApplication := newTestApplication(t, "app-002", "Frontend Developer", domain.StatusInterested)

	if err := repository.SaveApplication(context.Background(), firstApplication); err != nil {
		t.Fatalf("expected first application save to succeed: %v", err)
	}

	if err := repository.SaveApplication(context.Background(), secondApplication); err != nil {
		t.Fatalf("expected second application save to succeed: %v", err)
	}

	applications, err := repository.ListApplications(context.Background())
	if err != nil {
		t.Fatalf("expected application list to succeed: %v", err)
	}

	if len(applications) != 2 {
		t.Fatalf("expected two applications, got %d", len(applications))
	}
}

// -----------------------------------------------------------------------------
// openIntegrationDatabase
//
// Opens the PostgreSQL integration database and resets test data.
// -----------------------------------------------------------------------------
func openIntegrationDatabase(t *testing.T) *sql.DB {
	t.Helper()

	if os.Getenv("APPLYBY_INTEGRATION_TESTS") != "1" {
		t.Skip("set APPLYBY_INTEGRATION_TESTS=1 to run PostgreSQL integration tests")
	}

	db, err := OpenDatabase(context.Background(), config.LoadDatabaseConfig())
	if err != nil {
		t.Fatalf("failed to open integration database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	if err := RunMigrations(context.Background(), db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	_, err = db.ExecContext(context.Background(), "TRUNCATE TABLE applications, companies RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to reset integration test data: %v", err)
	}

	return db
}

// -----------------------------------------------------------------------------
// newTestApplication
//
// Creates a valid application for PostgreSQL repository tests.
// -----------------------------------------------------------------------------
func newTestApplication(t *testing.T, id domain.ApplicationID, title string, status domain.ApplicationStatus) domain.Application {
	t.Helper()

	createdAt := time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC)

	application, err := domain.NewApplication(
		id,
		title,
		domain.Company{
			Name:    "Example Studio",
			Website: "https://example.com",
		},
		status,
		createdAt,
	)
	if err != nil {
		t.Fatalf("failed to create test application: %v", err)
	}

	application.Source = "Company site"
	application.Notes = "Tracked during integration test."

	return application
}
