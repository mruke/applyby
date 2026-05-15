package application

import (
	"context"
	"errors"
	"testing"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// fakeRemoveApplicationRepository
//
// Provides application removal behavior for application-layer tests.
// -----------------------------------------------------------------------------
type fakeRemoveApplicationRepository struct {
	removedID domain.ApplicationID
	err       error
	called    bool
}

// -----------------------------------------------------------------------------
// RemoveApplication
//
// Records the removed application id and returns the configured fake result.
// -----------------------------------------------------------------------------
func (repository *fakeRemoveApplicationRepository) RemoveApplication(ctx context.Context, id domain.ApplicationID) error {
	repository.called = true
	repository.removedID = id

	if repository.err != nil {
		return repository.err
	}

	return nil
}

// -----------------------------------------------------------------------------
// TestRemoveApplicationServiceRemovesApplication
//
// Verifies that an application can be removed.
// -----------------------------------------------------------------------------
func TestRemoveApplicationServiceRemovesApplication(t *testing.T) {
	repository := &fakeRemoveApplicationRepository{}
	service := NewRemoveApplicationService(repository)

	err := service.Execute(context.Background(), RemoveApplicationInput{
		ID: "app-001",
	})
	if err != nil {
		t.Fatalf("expected application removal to succeed: %v", err)
	}

	if !repository.called {
		t.Fatal("expected repository to be called")
	}

	if repository.removedID != "app-001" {
		t.Fatalf("expected application id to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestRemoveApplicationServiceRejectsInvalidID
//
// Verifies that invalid application ids are rejected.
// -----------------------------------------------------------------------------
func TestRemoveApplicationServiceRejectsInvalidID(t *testing.T) {
	repository := &fakeRemoveApplicationRepository{}
	service := NewRemoveApplicationService(repository)

	err := service.Execute(context.Background(), RemoveApplicationInput{
		ID: "",
	})

	if err == nil {
		t.Fatal("expected invalid application id to fail")
	}

	if repository.called {
		t.Fatal("expected repository not to be called")
	}
}

// -----------------------------------------------------------------------------
// TestRemoveApplicationServiceReturnsRepositoryError
//
// Verifies repository errors are returned.
// -----------------------------------------------------------------------------
func TestRemoveApplicationServiceReturnsRepositoryError(t *testing.T) {
	repository := &fakeRemoveApplicationRepository{err: errors.New("remove failed")}
	service := NewRemoveApplicationService(repository)

	err := service.Execute(context.Background(), RemoveApplicationInput{
		ID: "app-001",
	})

	if err == nil {
		t.Fatal("expected repository error")
	}
}
