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

	event, err := NewActivityEvent(ActivityApplicationCreated, occurredAt, "Application created.")

	if err != nil {
		t.Fatalf("expected activity event to be valid: %v", err)
	}

	if event.Type != ActivityApplicationCreated {
		t.Fatalf("expected activity event type to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestNewActivityEventRejectsInvalidType
//
// Verifies that unsupported activity event types are rejected.
// -----------------------------------------------------------------------------
func TestNewActivityEventRejectsInvalidType(t *testing.T) {
	_, err := NewActivityEvent(ActivityEventType("unknown"), time.Now(), "Unknown event.")

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
	_, err := NewActivityEvent(ActivityApplicationCreated, time.Time{}, "Application created.")

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
	_, err := NewActivityEvent(ActivityApplicationCreated, time.Now(), "")

	if err == nil {
		t.Fatal("expected activity event without a description to be rejected")
	}
}
