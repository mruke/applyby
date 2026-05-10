package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestCompleteReminderServiceMarksReminderComplete
//
// Verifies that the complete reminder workflow marks and saves a reminder.
// -----------------------------------------------------------------------------
func TestCompleteReminderServiceMarksReminderComplete(t *testing.T) {
	repository := newFakeReminderRepository()
	service := NewCompleteReminderService(repository)
	reminder := newApplicationReminderTestReminder(t, "rem-001", false)

	repository.reminders[reminder.ID] = reminder

	completedReminder, err := service.Execute(context.Background(), CompleteReminderInput{
		ID: "rem-001",
	})

	if err != nil {
		t.Fatalf("expected complete reminder workflow to succeed: %v", err)
	}

	if !completedReminder.Completed {
		t.Fatal("expected reminder to be completed")
	}

	if repository.saveCalls != 1 {
		t.Fatalf("expected repository save to be called once")
	}
}

// -----------------------------------------------------------------------------
// TestCompleteReminderServiceRejectsMissingID
//
// Verifies that the complete reminder workflow requires a reminder id.
// -----------------------------------------------------------------------------
func TestCompleteReminderServiceRejectsMissingID(t *testing.T) {
	repository := newFakeReminderRepository()
	service := NewCompleteReminderService(repository)

	_, err := service.Execute(context.Background(), CompleteReminderInput{})

	if err == nil {
		t.Fatal("expected missing reminder id to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestCompleteReminderServiceReturnsFindError
//
// Verifies that repository lookup errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestCompleteReminderServiceReturnsFindError(t *testing.T) {
	repository := newFakeReminderRepository()
	repository.findErr = errors.New("find failed")
	service := NewCompleteReminderService(repository)

	_, err := service.Execute(context.Background(), CompleteReminderInput{
		ID: "rem-001",
	})

	if err == nil {
		t.Fatal("expected find error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestCompleteReminderServiceReturnsSaveError
//
// Verifies that repository save errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestCompleteReminderServiceReturnsSaveError(t *testing.T) {
	repository := newFakeReminderRepository()
	repository.saveErr = errors.New("save failed")
	service := NewCompleteReminderService(repository)
	reminder := newApplicationReminderTestReminder(t, "rem-001", false)

	repository.reminders[reminder.ID] = reminder

	_, err := service.Execute(context.Background(), CompleteReminderInput{
		ID: "rem-001",
	})

	if err == nil {
		t.Fatal("expected save error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestCompleteReminderServiceRequiresRepository
//
// Verifies that the complete reminder workflow requires a repository boundary.
// -----------------------------------------------------------------------------
func TestCompleteReminderServiceRequiresRepository(t *testing.T) {
	service := NewCompleteReminderService(nil)

	_, err := service.Execute(context.Background(), CompleteReminderInput{
		ID: "rem-001",
	})

	if err == nil {
		t.Fatal("expected missing repository to be rejected")
	}
}

// -----------------------------------------------------------------------------
// newApplicationReminderTestReminder
//
// Creates a valid reminder for application reminder workflow tests.
// -----------------------------------------------------------------------------
func newApplicationReminderTestReminder(t *testing.T, id domain.ReminderID, completed bool) domain.Reminder {
	t.Helper()

	reminder, err := domain.NewReminder(
		id,
		"app-001",
		"Follow up",
		time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to create test reminder: %v", err)
	}

	reminder.Completed = completed

	return reminder
}
