package application

import (
	"context"
	"errors"
	"testing"
	"time"
)

// -----------------------------------------------------------------------------
// TestListRemindersServiceReturnsPrioritizedReminders
//
// Verifies that the list reminder workflow returns reminders in priority order.
// -----------------------------------------------------------------------------
func TestListRemindersServiceReturnsPrioritizedReminders(t *testing.T) {
	repository := newFakeReminderRepository()
	service := NewListRemindersService(repository)

	laterReminder := newApplicationReminderTestReminder(t, "rem-001", false)
	laterReminder.DueAt = time.Date(2026, 5, 20, 9, 0, 0, 0, time.UTC)

	earlierReminder := newApplicationReminderTestReminder(t, "rem-002", false)
	earlierReminder.DueAt = time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC)

	repository.reminders[laterReminder.ID] = laterReminder
	repository.reminders[earlierReminder.ID] = earlierReminder

	applicationReminders, err := service.Execute(context.Background(), ListRemindersInput{
		ApplicationID: "app-001",
	})

	if err != nil {
		t.Fatalf("expected list reminders workflow to succeed: %v", err)
	}

	if len(applicationReminders) != 2 {
		t.Fatalf("expected two reminders, got %d", len(applicationReminders))
	}

	if applicationReminders[0].ID != "rem-002" {
		t.Fatalf("expected earliest reminder first")
	}
}

// -----------------------------------------------------------------------------
// TestListRemindersServiceRejectsMissingApplicationID
//
// Verifies that listing reminders requires an application id.
// -----------------------------------------------------------------------------
func TestListRemindersServiceRejectsMissingApplicationID(t *testing.T) {
	service := NewListRemindersService(newFakeReminderRepository())

	_, err := service.Execute(context.Background(), ListRemindersInput{})

	if err == nil {
		t.Fatal("expected missing application id to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestListRemindersServiceReturnsRepositoryError
//
// Verifies that repository list errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestListRemindersServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeReminderRepository()
	repository.listErr = errors.New("list failed")
	service := NewListRemindersService(repository)

	_, err := service.Execute(context.Background(), ListRemindersInput{
		ApplicationID: "app-001",
	})

	if err == nil {
		t.Fatal("expected list error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestListRemindersServiceRequiresRepository
//
// Verifies that the list reminders workflow requires a repository boundary.
// -----------------------------------------------------------------------------
func TestListRemindersServiceRequiresRepository(t *testing.T) {
	service := NewListRemindersService(nil)

	_, err := service.Execute(context.Background(), ListRemindersInput{
		ApplicationID: "app-001",
	})

	if err == nil {
		t.Fatal("expected missing repository to be rejected")
	}
}
