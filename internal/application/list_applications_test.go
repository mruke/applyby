package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestListApplicationsServiceReturnsApplications
//
// Verifies that the list workflow returns applications from the repository.
// -----------------------------------------------------------------------------
func TestListApplicationsServiceReturnsApplications(t *testing.T) {
	repository := newFakeApplicationRepository()
	service := NewListApplicationsService(repository)

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

	applications, err := service.Execute(context.Background())

	if err != nil {
		t.Fatalf("expected list applications workflow to succeed: %v", err)
	}

	if len(applications) != 1 {
		t.Fatalf("expected one application, got %d", len(applications))
	}
}

// -----------------------------------------------------------------------------
// TestListApplicationsServiceReturnsRepositoryError
//
// Verifies that repository list errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestListApplicationsServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeApplicationRepository()
	repository.listErr = errors.New("list failed")
	service := NewListApplicationsService(repository)

	_, err := service.Execute(context.Background())

	if err == nil {
		t.Fatal("expected list error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestListApplicationsServiceRequiresRepository
//
// Verifies that the list workflow requires a repository boundary.
// -----------------------------------------------------------------------------
func TestListApplicationsServiceRequiresRepository(t *testing.T) {
	service := NewListApplicationsService(nil)

	_, err := service.Execute(context.Background())

	if err == nil {
		t.Fatal("expected missing repository to be rejected")
	}
}
