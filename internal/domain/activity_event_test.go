package domain

import (
	"testing"
	"time"
)

// -----------------------------------------------------------------------------
// TestNewActivityEventAcceptsValidEvent
//
// Verifies that a valid activity event can be created.
// -----------------------------------------------------------------------------
func TestNewActivityEventAcceptsValidEvent(t *testing.T) {
	occurredAt := time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC)

	event, err := NewActivityEvent("app-001", ActivityApplicationCreated, occurredAt, "Application created.")

	if err != nil {
		t.Fatalf("expected activity event to be valid: %v", err)
	}

	if event.ApplicationID != "app-001" {
		t.Fatalf("expected activity event application id to be preserved")
	}

	if event.Type != ActivityApplicationCreated {
		t.Fatalf("expected activity event type to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestNewActivityEventRejectsMissingApplicationID
//
// Verifies that an activity event without an application id is rejected.
// -----------------------------------------------------------------------------
func TestNewActivityEventRejectsMissingApplicationID(t *testing.T) {
	_, err := NewActivityEvent("", ActivityApplicationCreated, time.Now(), "Application created.")

	if err == nil {
		t.Fatal("expected activity event without an application id to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestNewActivityEventRejectsInvalidType
//
// Verifies that unsupported activity event types are rejected.
// -----------------------------------------------------------------------------
func TestNewActivityEventRejectsInvalidType(t *testing.T) {
	_, err := NewActivityEvent("app-001", ActivityEventType("unknown"), time.Now(), "Unknown event.")

	if err == nil {
		t.Fatal("expected activity event with invalid type to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestNewActivityEventRejectsMissingTimestamp
//
// Verifies that an activity event without a timestamp is rejected.
// -----------------------------------------------------------------------------
func TestNewActivityEventRejectsMissingTimestamp(t *testing.T) {
	_, err := NewActivityEvent("app-001", ActivityApplicationCreated, time.Time{}, "Application created.")

	if err == nil {
		t.Fatal("expected activity event without a timestamp to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestNewActivityEventRejectsMissingDescription
//
// Verifies that an activity event without a description is rejected.
// -----------------------------------------------------------------------------
func TestNewActivityEventRejectsMissingDescription(t *testing.T) {
	_, err := NewActivityEvent("app-001", ActivityApplicationCreated, time.Now(), "")

	if err == nil {
		t.Fatal("expected activity event without a description to be rejected")
	}
}
