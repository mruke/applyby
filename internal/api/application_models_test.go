package api

import (
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestCreateApplicationRequestConvertsToInput
//
// Verifies that a valid create request converts into application-layer input.
// -----------------------------------------------------------------------------
func TestCreateApplicationRequestConvertsToInput(t *testing.T) {
	request := createApplicationRequest{
		ID:             "app-001",
		Title:          "Backend Developer",
		CompanyName:    "Example Studio",
		CompanyWebsite: "https://example.com",
		Status:         "applied",
		Source:         "Company site",
		Notes:          "Applied with backend resume.",
		CreatedAt:      "2026-05-10T08:00:00Z",
	}

	input, err := request.toInput()

	if err != nil {
		t.Fatalf("expected request conversion to succeed: %v", err)
	}

	if input.ID != "app-001" {
		t.Fatalf("expected application id to be preserved")
	}

	if input.Status != domain.StatusApplied {
		t.Fatalf("expected status to be parsed")
	}
}

// -----------------------------------------------------------------------------
// TestCreateApplicationRequestRejectsInvalidTimestamp
//
// Verifies that create request conversion rejects invalid timestamps.
// -----------------------------------------------------------------------------
func TestCreateApplicationRequestRejectsInvalidTimestamp(t *testing.T) {
	request := createApplicationRequest{
		ID:          "app-001",
		Title:       "Backend Developer",
		CompanyName: "Example Studio",
		Status:      "applied",
		CreatedAt:   "not-a-time",
	}

	_, err := request.toInput()

	if err == nil {
		t.Fatal("expected invalid timestamp to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestApplicationToResponse
//
// Verifies that a domain application converts into the expected API response.
// -----------------------------------------------------------------------------
func TestApplicationToResponse(t *testing.T) {
	createdAt := time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC)

	application, err := domain.NewApplication(
		"app-001",
		"Backend Developer",
		domain.Company{Name: "Example Studio", Website: "https://example.com"},
		domain.StatusApplied,
		createdAt,
	)
	if err != nil {
		t.Fatalf("failed to create application: %v", err)
	}

	response := applicationToResponse(application)

	if response.ID != "app-001" {
		t.Fatalf("expected response id to be preserved")
	}

	if response.CreatedAt != "2026-05-10T08:00:00Z" {
		t.Fatalf("expected created_at to be RFC3339")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateApplicationStatusRequestConvertsToStatus
//
// Verifies that an update status request converts into a domain status.
// -----------------------------------------------------------------------------
func TestUpdateApplicationStatusRequestConvertsToStatus(t *testing.T) {
	request := updateApplicationStatusRequest{Status: "interviewing"}

	status, err := request.toStatus()

	if err != nil {
		t.Fatalf("expected status conversion to succeed: %v", err)
	}

	if status != domain.StatusInterviewing {
		t.Fatalf("expected status to be interviewing")
	}
}
