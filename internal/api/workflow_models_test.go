package api

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestScheduleReminderRequestConvertsToInput
//
// Verifies that a valid schedule reminder request converts into workflow input.
// -----------------------------------------------------------------------------
func TestScheduleReminderRequestConvertsToInput(t *testing.T) {
	request := reminderRequest{
		ID:    "rem-001",
		Title: "Follow up",
		DueAt: "2026-05-10T09:00:00Z",
	}

	input, err := request.toInput("app-001")

	if err != nil {
		t.Fatalf("expected reminder request conversion to succeed: %v", err)
	}

	if input.ID != "rem-001" {
		t.Fatalf("expected reminder id to be preserved")
	}

	if input.ApplicationID != "app-001" {
		t.Fatalf("expected application id to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestSearchCriteriaFromRequest
//
// Verifies that query parameters convert into explicit search criteria.
// -----------------------------------------------------------------------------
func TestSearchCriteriaFromRequest(t *testing.T) {
	request := httptest.NewRequest(
		"GET",
		"/applications/search?status=applied&company_name=Example+Studio&source=Company+site&text=backend&created_from=2026-05-01T00:00:00Z&created_to=2026-05-31T00:00:00Z",
		nil,
	)

	criteria, err := searchCriteriaFromRequest(request)

	if err != nil {
		t.Fatalf("expected search criteria conversion to succeed: %v", err)
	}

	if len(criteria.Statuses) != 1 || criteria.Statuses[0] != domain.StatusApplied {
		t.Fatalf("expected applied status criteria")
	}

	if criteria.CompanyName != "Example Studio" {
		t.Fatalf("expected company name criteria to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestSearchCriteriaFromRequestRejectsInvalidStatus
//
// Verifies that invalid status query parameters are rejected.
// -----------------------------------------------------------------------------
func TestSearchCriteriaFromRequestRejectsInvalidStatus(t *testing.T) {
	request := httptest.NewRequest("GET", "/applications/search?status=paused", nil)

	_, err := searchCriteriaFromRequest(request)

	if err == nil {
		t.Fatal("expected invalid status criteria to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestReminderToResponse
//
// Verifies that a domain reminder converts into the expected API response.
// -----------------------------------------------------------------------------
func TestReminderToResponse(t *testing.T) {
	reminder, err := domain.NewReminder(
		"rem-001",
		"app-001",
		"Follow up",
		time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to create reminder: %v", err)
	}

	response := reminderToResponse(reminder)

	if response.ID != "rem-001" {
		t.Fatalf("expected reminder id to be preserved")
	}

	if response.DueAt != "2026-05-10T09:00:00Z" {
		t.Fatalf("expected due_at to be RFC3339")
	}
}

// -----------------------------------------------------------------------------
// TestContactAndDocumentRequestsConvertToInput
//
// Verifies that contact and document requests preserve application ownership.
// -----------------------------------------------------------------------------
func TestContactAndDocumentRequestsConvertToInput(t *testing.T) {
	contactInput := contactRequest{
		ID:    "contact-001",
		Name:  "Sam Recruiter",
		Email: "sam@example.com",
		Role:  "Recruiter",
	}.toInput("app-001")

	if contactInput.ApplicationID != "app-001" {
		t.Fatalf("expected contact application id to be preserved")
	}

	contactUpdateInput := contactRequest{
		Name:  "Sam Hiring",
		Email: "sam.hiring@example.com",
		Role:  "Hiring Manager",
	}.toUpdateInput("app-001", "contact-001")

	if contactUpdateInput.ApplicationID != "app-001" || contactUpdateInput.ContactID != "contact-001" {
		t.Fatalf("expected contact update identities to be preserved")
	}

	documentInput := documentRequest{
		ID:   "doc-001",
		Name: "Backend Resume",
		Kind: "resume",
		Path: "documents/backend-resume.pdf",
	}.toInput("app-001")

	if documentInput.ApplicationID != "app-001" {
		t.Fatalf("expected document application id to be preserved")
	}
}
