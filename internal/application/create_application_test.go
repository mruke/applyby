package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestCreateApplicationServiceSavesValidApplication
//
// Verifies that the create workflow validates and saves a valid application.
// -----------------------------------------------------------------------------
func TestCreateApplicationServiceSavesValidApplication(t *testing.T) {
	repository := newFakeApplicationRepository()
	service := NewCreateApplicationService(repository)
	createdAt := time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC)

	application, err := service.Execute(context.Background(), CreateApplicationInput{
		ID:        "app-001",
		Title:     "Backend Developer",
		Company:   domain.Company{Name: "Example Studio"},
		Status:    domain.StatusApplied,
		Source:    "Company site",
		Notes:     "Applied with backend resume.",
		CreatedAt: createdAt,
	})

	if err != nil {
		t.Fatalf("expected create application workflow to succeed: %v", err)
	}

	if application.ID != "app-001" {
		t.Fatalf("expected application id to be preserved")
	}

	if repository.saveCalls != 1 {
		t.Fatalf("expected repository save to be called once")
	}
}

// -----------------------------------------------------------------------------
// TestCreateApplicationServiceRejectsInvalidApplication
//
// Verifies that invalid domain data is rejected before saving.
// -----------------------------------------------------------------------------
func TestCreateApplicationServiceRejectsInvalidApplication(t *testing.T) {
	repository := newFakeApplicationRepository()
	service := NewCreateApplicationService(repository)

	_, err := service.Execute(context.Background(), CreateApplicationInput{
		ID:        "",
		Title:     "Backend Developer",
		Company:   domain.Company{Name: "Example Studio"},
		Status:    domain.StatusApplied,
		CreatedAt: time.Now(),
	})

	if err == nil {
		t.Fatal("expected invalid application to be rejected")
	}

	if repository.saveCalls != 0 {
		t.Fatal("expected invalid application not to be saved")
	}
}

// -----------------------------------------------------------------------------
// TestCreateApplicationServiceReturnsRepositoryError
//
// Verifies that repository save errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestCreateApplicationServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeApplicationRepository()
	repository.saveErr = errors.New("save failed")
	service := NewCreateApplicationService(repository)

	_, err := service.Execute(context.Background(), CreateApplicationInput{
		ID:        "app-001",
		Title:     "Backend Developer",
		Company:   domain.Company{Name: "Example Studio"},
		Status:    domain.StatusApplied,
		CreatedAt: time.Now(),
	})

	if err == nil {
		t.Fatal("expected repository error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestCreateApplicationServiceRequiresRepository
//
// Verifies that the create workflow requires a repository boundary.
// -----------------------------------------------------------------------------
func TestCreateApplicationServiceRequiresRepository(t *testing.T) {
	service := NewCreateApplicationService(nil)

	_, err := service.Execute(context.Background(), CreateApplicationInput{
		ID:        "app-001",
		Title:     "Backend Developer",
		Company:   domain.Company{Name: "Example Studio"},
		Status:    domain.StatusApplied,
		CreatedAt: time.Now(),
	})

	if err == nil {
		t.Fatal("expected missing repository to be rejected")
	}
}
