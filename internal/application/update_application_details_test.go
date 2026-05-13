package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestUpdateApplicationDetailsServiceUpdatesEditableFields
//
// Verifies that non-status application details can be updated.
// -----------------------------------------------------------------------------
func TestUpdateApplicationDetailsServiceUpdatesEditableFields(t *testing.T) {
	repository := newFakeApplicationRepository()
	historyRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateApplicationDetailsService(repository, historyRepository)

	application := newUpdateDetailsTestApplication(t)
	repository.applications[application.ID] = application

	updatedApplication, err := service.Execute(context.Background(), UpdateApplicationDetailsInput{
		ID:             application.ID,
		Title:          "Senior Backend Developer",
		CompanyName:    "Example Labs",
		CompanyWebsite: "https://labs.example.com",
		Source:         "Referral",
		Notes:          "Updated application details.",
	})
	if err != nil {
		t.Fatalf("expected update details workflow to succeed: %v", err)
	}

	if updatedApplication.Title != "Senior Backend Developer" {
		t.Fatalf("expected updated title, got %q", updatedApplication.Title)
	}

	if updatedApplication.Company.Name != "Example Labs" {
		t.Fatalf("expected updated company, got %q", updatedApplication.Company.Name)
	}

	if updatedApplication.Status != application.Status {
		t.Fatalf("expected status to remain %q, got %q", application.Status, updatedApplication.Status)
	}

	if len(historyRepository.activityEvents) != 1 {
		t.Fatalf("expected one activity event, got %d", len(historyRepository.activityEvents))
	}

	if historyRepository.activityEvents[0].Type != domain.ActivityApplicationUpdated {
		t.Fatalf("expected application updated activity event")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateApplicationDetailsServiceRejectsInvalidDetails
//
// Verifies that invalid edited details are rejected.
// -----------------------------------------------------------------------------
func TestUpdateApplicationDetailsServiceRejectsInvalidDetails(t *testing.T) {
	repository := newFakeApplicationRepository()
	historyRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateApplicationDetailsService(repository, historyRepository)

	application := newUpdateDetailsTestApplication(t)
	repository.applications[application.ID] = application

	_, err := service.Execute(context.Background(), UpdateApplicationDetailsInput{
		ID:          application.ID,
		Title:       "",
		CompanyName: "Example Labs",
	})

	if err == nil {
		t.Fatal("expected invalid details to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateApplicationDetailsServiceReturnsFindError
//
// Verifies that lookup errors are returned.
// -----------------------------------------------------------------------------
func TestUpdateApplicationDetailsServiceReturnsFindError(t *testing.T) {
	repository := newFakeApplicationRepository()
	repository.findErr = errors.New("find failed")
	historyRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateApplicationDetailsService(repository, historyRepository)

	_, err := service.Execute(context.Background(), UpdateApplicationDetailsInput{
		ID:          "app-001",
		Title:       "Backend Developer",
		CompanyName: "Example Studio",
	})

	if err == nil {
		t.Fatal("expected find error")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateApplicationDetailsServiceReturnsUpdateError
//
// Verifies that repository update errors are returned.
// -----------------------------------------------------------------------------
func TestUpdateApplicationDetailsServiceReturnsUpdateError(t *testing.T) {
	repository := newFakeApplicationRepository()
	repository.updateDetailsErr = errors.New("update failed")
	historyRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateApplicationDetailsService(repository, historyRepository)

	application := newUpdateDetailsTestApplication(t)
	repository.applications[application.ID] = application

	_, err := service.Execute(context.Background(), UpdateApplicationDetailsInput{
		ID:          application.ID,
		Title:       "Backend Developer",
		CompanyName: "Example Studio",
	})

	if err == nil {
		t.Fatal("expected update error")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateApplicationDetailsServiceReturnsActivityError
//
// Verifies that activity recording errors are returned.
// -----------------------------------------------------------------------------
func TestUpdateApplicationDetailsServiceReturnsActivityError(t *testing.T) {
	repository := newFakeApplicationRepository()
	historyRepository := &fakeApplicationHistoryRepository{recordEventErr: errors.New("activity failed")}
	service := NewUpdateApplicationDetailsService(repository, historyRepository)

	application := newUpdateDetailsTestApplication(t)
	repository.applications[application.ID] = application

	_, err := service.Execute(context.Background(), UpdateApplicationDetailsInput{
		ID:          application.ID,
		Title:       "Backend Developer",
		CompanyName: "Example Studio",
	})

	if err == nil {
		t.Fatal("expected activity error")
	}
}

// -----------------------------------------------------------------------------
// newUpdateDetailsTestApplication
//
// Creates a valid application for update details tests.
// -----------------------------------------------------------------------------
func newUpdateDetailsTestApplication(t *testing.T) domain.Application {
	t.Helper()

	application, err := domain.NewApplication(
		"app-001",
		"Backend Developer",
		domain.Company{Name: "Example Studio", Website: "https://example.com"},
		domain.StatusApplied,
		time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to create application: %v", err)
	}

	application.Source = "Company site"
	application.Notes = "Original notes."

	return application
}
