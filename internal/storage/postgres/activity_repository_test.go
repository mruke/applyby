package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestApplicationRepositoryRecordsStatusHistory
//
// Verifies that PostgreSQL persistence records structured status history.
// -----------------------------------------------------------------------------
func TestApplicationRepositoryRecordsStatusHistory(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)
	application := newActivityRepositoryTestApplication(t, "app-001", domain.StatusApplied)

	if err := repository.SaveApplication(context.Background(), application); err != nil {
		t.Fatalf("expected application save to succeed: %v", err)
	}

	history, err := domain.NewApplicationStatusHistory(
		application.ID,
		domain.StatusApplied,
		domain.StatusInterviewing,
		time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to create status history: %v", err)
	}

	if err := repository.RecordApplicationStatusHistory(context.Background(), history); err != nil {
		t.Fatalf("expected status history record to succeed: %v", err)
	}

	var count int
	err = db.QueryRowContext(
		context.Background(),
		"SELECT COUNT(*) FROM application_status_history WHERE application_id = $1",
		application.ID,
	).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count status history records: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected one status history record, got %d", count)
	}
}

// -----------------------------------------------------------------------------
// TestApplicationRepositoryRecordsAndListsActivityEvents
//
// Verifies that PostgreSQL persistence records and lists activity timeline events.
// -----------------------------------------------------------------------------
func TestApplicationRepositoryRecordsAndListsActivityEvents(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)
	application := newActivityRepositoryTestApplication(t, "app-001", domain.StatusApplied)

	if err := repository.SaveApplication(context.Background(), application); err != nil {
		t.Fatalf("expected application save to succeed: %v", err)
	}

	event, err := domain.NewActivityEvent(
		application.ID,
		domain.ActivityStatusChanged,
		time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC),
		"Status changed from applied to interviewing.",
	)
	if err != nil {
		t.Fatalf("failed to create activity event: %v", err)
	}

	if err := repository.RecordActivityEvent(context.Background(), event); err != nil {
		t.Fatalf("expected activity event record to succeed: %v", err)
	}

	events, err := repository.ListActivityEventsForApplication(context.Background(), application.ID)
	if err != nil {
		t.Fatalf("expected activity event list to succeed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected one activity event, got %d", len(events))
	}

	if events[0].Description != "Status changed from applied to interviewing." {
		t.Fatalf("expected activity event description to be preserved")
	}
}

// -----------------------------------------------------------------------------
// newActivityRepositoryTestApplication
//
// Creates a valid application for PostgreSQL activity repository tests.
// -----------------------------------------------------------------------------
func newActivityRepositoryTestApplication(t *testing.T, id domain.ApplicationID, status domain.ApplicationStatus) domain.Application {
	t.Helper()

	application, err := domain.NewApplication(
		id,
		"Backend Developer",
		domain.Company{Name: "Example Studio", Website: "https://example.com"},
		status,
		time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to create activity repository test application: %v", err)
	}

	return application
}
