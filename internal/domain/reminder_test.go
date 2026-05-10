package domain

import (
	"testing"
	"time"
)

// -----------------------------------------------------------------------------
// TestNewReminderAcceptsValidReminder
//
// Verifies that a reminder with identity, application, title, and due date can be created.
// -----------------------------------------------------------------------------
func TestNewReminderAcceptsValidReminder(t *testing.T) {
	dueAt := time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC)

	reminder, err := NewReminder("rem-001", "app-001", "Follow up", dueAt)

	if err != nil {
		t.Fatalf("expected reminder to be valid: %v", err)
	}

	if reminder.DueAt != dueAt {
		t.Fatalf("expected reminder due date to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestNewReminderRejectsMissingID
//
// Verifies that a reminder without an identity is rejected.
// -----------------------------------------------------------------------------
func TestNewReminderRejectsMissingID(t *testing.T) {
	_, err := NewReminder("", "app-001", "Follow up", time.Now())

	if err == nil {
		t.Fatal("expected reminder without an id to be invalid")
	}
}

// -----------------------------------------------------------------------------
// TestNewReminderRejectsMissingApplicationID
//
// Verifies that a reminder without an application identity is rejected.
// -----------------------------------------------------------------------------
func TestNewReminderRejectsMissingApplicationID(t *testing.T) {
	_, err := NewReminder("rem-001", "", "Follow up", time.Now())

	if err == nil {
		t.Fatal("expected reminder without an application id to be invalid")
	}
}

// -----------------------------------------------------------------------------
// TestNewReminderRejectsMissingTitle
//
// Verifies that a reminder without a title is rejected.
// -----------------------------------------------------------------------------
func TestNewReminderRejectsMissingTitle(t *testing.T) {
	_, err := NewReminder("rem-001", "app-001", "", time.Now())

	if err == nil {
		t.Fatal("expected reminder without a title to be invalid")
	}
}

// -----------------------------------------------------------------------------
// TestNewReminderRejectsMissingDueDate
//
// Verifies that a reminder without a due date is rejected.
// -----------------------------------------------------------------------------
func TestNewReminderRejectsMissingDueDate(t *testing.T) {
	_, err := NewReminder("rem-001", "app-001", "Follow up", time.Time{})

	if err == nil {
		t.Fatal("expected reminder without a due date to be invalid")
	}
}
