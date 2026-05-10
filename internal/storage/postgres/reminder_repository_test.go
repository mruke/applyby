package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestApplicationRepositorySavesAndFindsReminder
//
// Verifies that PostgreSQL persistence can save and retrieve a reminder.
// -----------------------------------------------------------------------------
func TestApplicationRepositorySavesAndFindsReminder(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)
	application := newReminderRepositoryTestApplication(t)

	if err := repository.SaveApplication(context.Background(), application); err != nil {
		t.Fatalf("expected application save to succeed: %v", err)
	}

	reminder := newReminderRepositoryTestReminder(t, "rem-001", application.ID, false)

	if err := repository.SaveReminder(context.Background(), reminder); err != nil {
		t.Fatalf("expected reminder save to succeed: %v", err)
	}

	foundReminder, err := repository.FindReminderByID(context.Background(), reminder.ID)
	if err != nil {
		t.Fatalf("expected reminder lookup to succeed: %v", err)
	}

	if foundReminder.ID != reminder.ID {
		t.Fatalf("expected reminder id %q, got %q", reminder.ID, foundReminder.ID)
	}
}

// -----------------------------------------------------------------------------
// TestApplicationRepositoryUpdatesExistingReminder
//
// Verifies that saving an existing reminder updates current reminder state.
// -----------------------------------------------------------------------------
func TestApplicationRepositoryUpdatesExistingReminder(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)
	application := newReminderRepositoryTestApplication(t)

	if err := repository.SaveApplication(context.Background(), application); err != nil {
		t.Fatalf("expected application save to succeed: %v", err)
	}

	reminder := newReminderRepositoryTestReminder(t, "rem-001", application.ID, false)

	if err := repository.SaveReminder(context.Background(), reminder); err != nil {
		t.Fatalf("expected initial reminder save to succeed: %v", err)
	}

	reminder.Completed = true

	if err := repository.SaveReminder(context.Background(), reminder); err != nil {
		t.Fatalf("expected reminder update save to succeed: %v", err)
	}

	foundReminder, err := repository.FindReminderByID(context.Background(), reminder.ID)
	if err != nil {
		t.Fatalf("expected reminder lookup to succeed: %v", err)
	}

	if !foundReminder.Completed {
		t.Fatal("expected reminder to be completed")
	}
}

// -----------------------------------------------------------------------------
// TestApplicationRepositoryListsRemindersForApplication
//
// Verifies that PostgreSQL persistence lists reminders for one application.
// -----------------------------------------------------------------------------
func TestApplicationRepositoryListsRemindersForApplication(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)
	application := newReminderRepositoryTestApplication(t)

	if err := repository.SaveApplication(context.Background(), application); err != nil {
		t.Fatalf("expected application save to succeed: %v", err)
	}

	laterReminder := newReminderRepositoryTestReminder(t, "rem-001", application.ID, false)
	laterReminder.Title = "Later"
	laterReminder.DueAt = time.Date(2026, 5, 20, 9, 0, 0, 0, time.UTC)

	earlierReminder := newReminderRepositoryTestReminder(t, "rem-002", application.ID, false)
	earlierReminder.Title = "Earlier"
	earlierReminder.DueAt = time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC)

	if err := repository.SaveReminder(context.Background(), laterReminder); err != nil {
		t.Fatalf("expected later reminder save to succeed: %v", err)
	}

	if err := repository.SaveReminder(context.Background(), earlierReminder); err != nil {
		t.Fatalf("expected earlier reminder save to succeed: %v", err)
	}

	reminders, err := repository.ListRemindersForApplication(context.Background(), application.ID)
	if err != nil {
		t.Fatalf("expected reminder list to succeed: %v", err)
	}

	if len(reminders) != 2 {
		t.Fatalf("expected two reminders, got %d", len(reminders))
	}

	if reminders[0].ID != earlierReminder.ID {
		t.Fatalf("expected earliest reminder first")
	}
}

// -----------------------------------------------------------------------------
// newReminderRepositoryTestApplication
//
// Creates a valid application for PostgreSQL reminder repository tests.
// -----------------------------------------------------------------------------
func newReminderRepositoryTestApplication(t *testing.T) domain.Application {
	t.Helper()

	application, err := domain.NewApplication(
		"app-001",
		"Backend Developer",
		domain.Company{Name: "Example Studio", Website: "https://example.com"},
		domain.StatusApplied,
		time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to create reminder repository test application: %v", err)
	}

	return application
}

// -----------------------------------------------------------------------------
// newReminderRepositoryTestReminder
//
// Creates a valid reminder for PostgreSQL reminder repository tests.
// -----------------------------------------------------------------------------
func newReminderRepositoryTestReminder(t *testing.T, id domain.ReminderID, applicationID domain.ApplicationID, completed bool) domain.Reminder {
	t.Helper()

	reminder, err := domain.NewReminder(
		id,
		applicationID,
		"Follow up",
		time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to create reminder repository test reminder: %v", err)
	}

	reminder.Completed = completed

	return reminder
}
