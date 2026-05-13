package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestGetApplicationServiceReturnsApplication
//
// Verifies that the detail workflow returns one application from the repository.
// -----------------------------------------------------------------------------
func TestGetApplicationServiceReturnsApplication(t *testing.T) {
	repository := newFakeApplicationRepository()
	service := NewGetApplicationService(repository)

	application, err := domain.NewApplication(
		"app-001",
		"Backend Developer",
		domain.Company{Name: "Example Studio"},
		domain.StatusApplied,
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create test application: %v", err)
	}

	repository.applications[application.ID] = application

	foundApplication, err := service.Execute(context.Background(), GetApplicationInput{ID: application.ID})

	if err != nil {
		t.Fatalf("expected get application workflow to succeed: %v", err)
	}

	if foundApplication.ID != application.ID {
		t.Fatalf("expected application id %q, got %q", application.ID, foundApplication.ID)
	}
}

// -----------------------------------------------------------------------------
// TestGetApplicationServiceReturnsRepositoryError
//
// Verifies that repository lookup errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestGetApplicationServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeApplicationRepository()
	repository.findErr = errors.New("find failed")
	service := NewGetApplicationService(repository)

	_, err := service.Execute(context.Background(), GetApplicationInput{ID: "app-001"})

	if err == nil {
		t.Fatal("expected find error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestGetApplicationServiceRejectsInvalidID
//
// Verifies that invalid application identities are rejected before lookup.
// -----------------------------------------------------------------------------
func TestGetApplicationServiceRejectsInvalidID(t *testing.T) {
	repository := newFakeApplicationRepository()
	service := NewGetApplicationService(repository)

	_, err := service.Execute(context.Background(), GetApplicationInput{ID: ""})

	if err == nil {
		t.Fatal("expected invalid application id to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestGetApplicationServiceRequiresRepository
//
// Verifies that the detail workflow requires a repository boundary.
// -----------------------------------------------------------------------------
func TestGetApplicationServiceRequiresRepository(t *testing.T) {
	service := NewGetApplicationService(nil)

	_, err := service.Execute(context.Background(), GetApplicationInput{ID: "app-001"})

	if err == nil {
		t.Fatal("expected missing repository to be rejected")
	}
}
