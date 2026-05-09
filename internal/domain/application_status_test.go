package domain

import "testing"

// -----------------------------------------------------------------------------
// TestApplicationStatusValidationAcceptsSupportedStatuses
//
// Verifies that every defined application status is treated as valid.
// -----------------------------------------------------------------------------
func TestApplicationStatusValidationAcceptsSupportedStatuses(t *testing.T) {
	for _, status := range AllApplicationStatuses() {
		if !IsValidApplicationStatus(status) {
			t.Fatalf("expected status %q to be valid", status)
		}
	}
}

// -----------------------------------------------------------------------------
// TestApplicationStatusValidationRejectsUnsupportedStatus
//
// Verifies that unknown application statuses are rejected.
// -----------------------------------------------------------------------------
func TestApplicationStatusValidationRejectsUnsupportedStatus(t *testing.T) {
	if IsValidApplicationStatus(ApplicationStatus("unknown")) {
		t.Fatal("expected unknown status to be invalid")
	}
}

// -----------------------------------------------------------------------------
// TestParseApplicationStatusNormalizesKnownStatus
//
// Verifies that status parsing trims whitespace and normalizes case.
// -----------------------------------------------------------------------------
func TestParseApplicationStatusNormalizesKnownStatus(t *testing.T) {
	status, err := ParseApplicationStatus(" Applied ")

	if err != nil {
		t.Fatalf("expected status to parse successfully: %v", err)
	}

	if status != StatusApplied {
		t.Fatalf("expected %q, got %q", StatusApplied, status)
	}
}

// -----------------------------------------------------------------------------
// TestParseApplicationStatusRejectsUnknownStatus
//
// Verifies that parsing fails for unsupported status text.
// -----------------------------------------------------------------------------
func TestParseApplicationStatusRejectsUnknownStatus(t *testing.T) {
	_, err := ParseApplicationStatus("paused")

	if err == nil {
		t.Fatal("expected parse error for unsupported status")
	}
}
