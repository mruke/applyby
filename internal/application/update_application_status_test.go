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
	historyRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateApplicationStatusService(repository, historyRepository)

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

	if len(historyRepository.statusHistory) != 1 {
		t.Fatalf("expected one status history record")
	}

	if len(historyRepository.activityEvents) != 1 {
		t.Fatalf("expected one activity event")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateApplicationStatusServiceAllowsCorrectiveTransition
//
// Verifies that the status workflow allows corrective lifecycle movement.
// -----------------------------------------------------------------------------
func TestUpdateApplicationStatusServiceAllowsCorrectiveTransition(t *testing.T) {
	repository := newFakeApplicationRepository()
	historyRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateApplicationStatusService(repository, historyRepository)

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

	updatedApplication, err := service.Execute(context.Background(), UpdateApplicationStatusInput{
		ID:     "app-001",
		Status: domain.StatusInterviewing,
	})
	if err != nil {
		t.Fatalf("expected corrective status transition to succeed: %v", err)
	}

	if updatedApplication.Status != domain.StatusInterviewing {
		t.Fatalf("expected application status to be updated")
	}

	if repository.saveCalls != 1 {
		t.Fatal("expected corrective status transition to be saved")
	}

	if len(historyRepository.statusHistory) != 1 {
		t.Fatal("expected corrective status transition to record status history")
	}

	if len(historyRepository.activityEvents) != 1 {
		t.Fatal("expected corrective status transition to record activity event")
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
	historyRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateApplicationStatusService(repository, historyRepository)

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
	historyRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateApplicationStatusService(repository, historyRepository)

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

	if len(historyRepository.statusHistory) != 0 {
		t.Fatal("expected save failure not to record status history")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateApplicationStatusServiceReturnsStatusHistoryError
//
// Verifies that status history recording errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestUpdateApplicationStatusServiceReturnsStatusHistoryError(t *testing.T) {
	repository := newFakeApplicationRepository()
	historyRepository := &fakeApplicationHistoryRepository{
		recordStatusErr: errors.New("status history failed"),
	}
	service := NewUpdateApplicationStatusService(repository, historyRepository)

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
		t.Fatal("expected status history error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateApplicationStatusServiceReturnsActivityEventError
//
// Verifies that activity event recording errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestUpdateApplicationStatusServiceReturnsActivityEventError(t *testing.T) {
	repository := newFakeApplicationRepository()
	historyRepository := &fakeApplicationHistoryRepository{
		recordEventErr: errors.New("activity event failed"),
	}
	service := NewUpdateApplicationStatusService(repository, historyRepository)

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
		t.Fatal("expected activity event error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateApplicationStatusServiceRequiresRepository
//
// Verifies that the status workflow requires a repository boundary.
// -----------------------------------------------------------------------------
func TestUpdateApplicationStatusServiceRequiresRepository(t *testing.T) {
	service := NewUpdateApplicationStatusService(nil, &fakeApplicationHistoryRepository{})

	_, err := service.Execute(context.Background(), UpdateApplicationStatusInput{
		ID:     "app-001",
		Status: domain.StatusInterviewing,
	})

	if err == nil {
		t.Fatal("expected missing repository to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateApplicationStatusServiceRequiresHistoryRecorder
//
// Verifies that the status workflow requires a history recorder boundary.
// -----------------------------------------------------------------------------
func TestUpdateApplicationStatusServiceRequiresHistoryRecorder(t *testing.T) {
	repository := newFakeApplicationRepository()
	service := NewUpdateApplicationStatusService(repository, nil)

	_, err := service.Execute(context.Background(), UpdateApplicationStatusInput{
		ID:     "app-001",
		Status: domain.StatusInterviewing,
	})

	if err == nil {
		t.Fatal("expected missing history recorder to be rejected")
	}
}
