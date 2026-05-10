package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestUpdateApplicationStatusServiceSavesAllowedTransition
//
// Verifies that the status workflow saves a valid lifecycle transition.
// -----------------------------------------------------------------------------
func TestUpdateApplicationStatusServiceSavesAllowedTransition(t *testing.T) {
	repository := newFakeApplicationRepository()
	service := NewUpdateApplicationStatusService(repository)

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

	updatedApplication, err := service.Execute(context.Background(), UpdateApplicationStatusInput{
		ID:     "app-001",
		Status: domain.StatusInterviewing,
	})

	if err != nil {
		t.Fatalf("expected status update workflow to succeed: %v", err)
	}

	if updatedApplication.Status != domain.StatusInterviewing {
		t.Fatalf("expected application status to be updated")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateApplicationStatusServiceRejectsInvalidTransition
//
// Verifies that the status workflow rejects unsupported lifecycle movement.
// -----------------------------------------------------------------------------
func TestUpdateApplicationStatusServiceRejectsInvalidTransition(t *testing.T) {
	repository := newFakeApplicationRepository()
	service := NewUpdateApplicationStatusService(repository)

	application, err := domain.NewApplication(
		"app-001",
		"Backend Developer",
		domain.Company{Name: "Example Studio"},
		domain.StatusRejected,
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create test application: %v", err)
	}

	repository.applications[application.ID] = application

	_, err = service.Execute(context.Background(), UpdateApplicationStatusInput{
		ID:     "app-001",
		Status: domain.StatusInterviewing,
	})

	if err == nil {
		t.Fatal("expected invalid status transition to be rejected")
	}

	if repository.saveCalls != 0 {
		t.Fatal("expected invalid status transition not to be saved")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateApplicationStatusServiceReturnsFindError
//
// Verifies that repository lookup errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestUpdateApplicationStatusServiceReturnsFindError(t *testing.T) {
	repository := newFakeApplicationRepository()
	repository.findErr = errors.New("find failed")
	service := NewUpdateApplicationStatusService(repository)

	_, err := service.Execute(context.Background(), UpdateApplicationStatusInput{
		ID:     "app-001",
		Status: domain.StatusInterviewing,
	})

	if err == nil {
		t.Fatal("expected find error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateApplicationStatusServiceReturnsSaveError
//
// Verifies that repository save errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestUpdateApplicationStatusServiceReturnsSaveError(t *testing.T) {
	repository := newFakeApplicationRepository()
	repository.saveErr = errors.New("save failed")
	service := NewUpdateApplicationStatusService(repository)

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

	_, err = service.Execute(context.Background(), UpdateApplicationStatusInput{
		ID:     "app-001",
		Status: domain.StatusInterviewing,
	})

	if err == nil {
		t.Fatal("expected save error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateApplicationStatusServiceRequiresRepository
//
// Verifies that the status workflow requires a repository boundary.
// -----------------------------------------------------------------------------
func TestUpdateApplicationStatusServiceRequiresRepository(t *testing.T) {
	service := NewUpdateApplicationStatusService(nil)

	_, err := service.Execute(context.Background(), UpdateApplicationStatusInput{
		ID:     "app-001",
		Status: domain.StatusInterviewing,
	})

	if err == nil {
		t.Fatal("expected missing repository to be rejected")
	}
}
