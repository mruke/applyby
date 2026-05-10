package domain

import (
	"testing"
	"time"
)

// -----------------------------------------------------------------------------
// TestNewApplicationAcceptsValidApplication
//
// Verifies that a valid application can be created.
// -----------------------------------------------------------------------------
func TestNewApplicationAcceptsValidApplication(t *testing.T) {
	id := ApplicationID("app-001")
	company := Company{Name: "Example Studio"}
	createdAt := time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC)

	application, err := NewApplication(id, "Backend Developer", company, StatusApplied, createdAt)

	if err != nil {
		t.Fatalf("expected application to be valid: %v", err)
	}

	if application.Title != "Backend Developer" {
		t.Fatalf("expected application title to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestNewApplicationRejectsMissingID
//
// Verifies that an application without an identity is rejected.
// -----------------------------------------------------------------------------
func TestNewApplicationRejectsMissingID(t *testing.T) {
	company := Company{Name: "Example Studio"}

	_, err := NewApplication("", "Backend Developer", company, StatusApplied, time.Now())

	if err == nil {
		t.Fatal("expected application without an id to be invalid")
	}
}

// -----------------------------------------------------------------------------
// TestNewApplicationRejectsMissingTitle
//
// Verifies that an application without a title is rejected.
// -----------------------------------------------------------------------------
func TestNewApplicationRejectsMissingTitle(t *testing.T) {
	company := Company{Name: "Example Studio"}

	_, err := NewApplication("app-001", "", company, StatusApplied, time.Now())

	if err == nil {
		t.Fatal("expected application without a title to be invalid")
	}
}

// -----------------------------------------------------------------------------
// TestNewApplicationRejectsInvalidCompany
//
// Verifies that an application with an invalid company is rejected.
// -----------------------------------------------------------------------------
func TestNewApplicationRejectsInvalidCompany(t *testing.T) {
	company := Company{}

	_, err := NewApplication("app-001", "Backend Developer", company, StatusApplied, time.Now())

	if err == nil {
		t.Fatal("expected application with invalid company to be invalid")
	}
}

// -----------------------------------------------------------------------------
// TestNewApplicationRejectsInvalidStatus
//
// Verifies that an application with an unsupported status is rejected.
// -----------------------------------------------------------------------------
func TestNewApplicationRejectsInvalidStatus(t *testing.T) {
	company := Company{Name: "Example Studio"}

	_, err := NewApplication("app-001", "Backend Developer", company, ApplicationStatus("paused"), time.Now())

	if err == nil {
		t.Fatal("expected application with invalid status to be invalid")
	}
}

// -----------------------------------------------------------------------------
// TestApplicationValidateRejectsAppliedDateBeforeCreatedDate
//
// Verifies that an applied date cannot come before the application record exists.
// -----------------------------------------------------------------------------
func TestApplicationValidateRejectsAppliedDateBeforeCreatedDate(t *testing.T) {
	company := Company{Name: "Example Studio"}
	createdAt := time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC)
	appliedAt := time.Date(2026, 5, 9, 8, 0, 0, 0, time.UTC)

	application := Application{
		ID:        "app-001",
		Title:     "Backend Developer",
		Company:   company,
		Status:    StatusApplied,
		CreatedAt: createdAt,
		AppliedAt: &appliedAt,
	}

	err := application.Validate()

	if err == nil {
		t.Fatal("expected application with applied date before created date to be invalid")
	}
}
