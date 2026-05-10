package domain

import (
	"testing"
	"time"
)

// -----------------------------------------------------------------------------
// TestNewApplicationStatusHistoryAcceptsValidHistory
//
// Verifies that valid status history can be created.
// -----------------------------------------------------------------------------
func TestNewApplicationStatusHistoryAcceptsValidHistory(t *testing.T) {
	changedAt := time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC)

	history, err := NewApplicationStatusHistory("app-001", StatusApplied, StatusInterviewing, changedAt)

	if err != nil {
		t.Fatalf("expected status history to be valid: %v", err)
	}

	if history.ApplicationID != "app-001" {
		t.Fatalf("expected application id to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestNewApplicationStatusHistoryRejectsMissingApplicationID
//
// Verifies that status history without an application id is rejected.
// -----------------------------------------------------------------------------
func TestNewApplicationStatusHistoryRejectsMissingApplicationID(t *testing.T) {
	_, err := NewApplicationStatusHistory("", StatusApplied, StatusInterviewing, time.Now())

	if err == nil {
		t.Fatal("expected status history without an application id to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestNewApplicationStatusHistoryRejectsInvalidTransition
//
// Verifies that status history cannot record unsupported lifecycle movement.
// -----------------------------------------------------------------------------
func TestNewApplicationStatusHistoryRejectsInvalidTransition(t *testing.T) {
	_, err := NewApplicationStatusHistory("app-001", StatusRejected, StatusInterviewing, time.Now())

	if err == nil {
		t.Fatal("expected invalid status transition to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestNewApplicationStatusHistoryRejectsMissingChangedAt
//
// Verifies that status history requires a timestamp.
// -----------------------------------------------------------------------------
func TestNewApplicationStatusHistoryRejectsMissingChangedAt(t *testing.T) {
	_, err := NewApplicationStatusHistory("app-001", StatusApplied, StatusInterviewing, time.Time{})

	if err == nil {
		t.Fatal("expected status history without changed date to be rejected")
	}
}
