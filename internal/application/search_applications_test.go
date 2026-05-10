package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
	"github.com/mruke/applyby/internal/search"
)

// -----------------------------------------------------------------------------
// fakeApplicationSearcher
//
// Provides an in-memory application searcher for application-layer tests.
// -----------------------------------------------------------------------------
type fakeApplicationSearcher struct {
	applications []domain.Application
	err          error
	called       bool
}

// -----------------------------------------------------------------------------
// SearchApplications
//
// Searches in-memory applications for application-layer tests.
// -----------------------------------------------------------------------------
func (searcher *fakeApplicationSearcher) SearchApplications(ctx context.Context, criteria search.ApplicationCriteria) ([]domain.Application, error) {
	searcher.called = true

	if searcher.err != nil {
		return nil, searcher.err
	}

	matches := []domain.Application{}

	for _, application := range searcher.applications {
		if search.MatchesApplication(application, criteria) {
			matches = append(matches, application)
		}
	}

	return matches, nil
}

// -----------------------------------------------------------------------------
// TestSearchApplicationsServiceReturnsMatches
//
// Verifies that the search workflow returns matching applications.
// -----------------------------------------------------------------------------
func TestSearchApplicationsServiceReturnsMatches(t *testing.T) {
	application := newApplicationSearchTestApplication(t, "app-001", domain.StatusApplied)
	fakeSearcher := &fakeApplicationSearcher{
		applications: []domain.Application{application},
	}
	service := NewSearchApplicationsService(fakeSearcher)

	applications, err := service.Execute(context.Background(), search.ApplicationCriteria{
		Statuses: []domain.ApplicationStatus{domain.StatusApplied},
		Text:     "backend",
	})

	if err != nil {
		t.Fatalf("expected search workflow to succeed: %v", err)
	}

	if len(applications) != 1 {
		t.Fatalf("expected one matching application, got %d", len(applications))
	}

	if !fakeSearcher.called {
		t.Fatal("expected searcher to be called")
	}
}

// -----------------------------------------------------------------------------
// TestSearchApplicationsServiceRejectsInvalidCriteria
//
// Verifies that invalid search criteria are rejected before searching.
// -----------------------------------------------------------------------------
func TestSearchApplicationsServiceRejectsInvalidCriteria(t *testing.T) {
	fakeSearcher := &fakeApplicationSearcher{}
	service := NewSearchApplicationsService(fakeSearcher)

	_, err := service.Execute(context.Background(), search.ApplicationCriteria{
		Statuses: []domain.ApplicationStatus{domain.ApplicationStatus("paused")},
	})

	if err == nil {
		t.Fatal("expected invalid criteria to be rejected")
	}

	if fakeSearcher.called {
		t.Fatal("expected searcher not to be called for invalid criteria")
	}
}

// -----------------------------------------------------------------------------
// TestSearchApplicationsServiceReturnsSearchError
//
// Verifies that repository search errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestSearchApplicationsServiceReturnsSearchError(t *testing.T) {
	fakeSearcher := &fakeApplicationSearcher{
		err: errors.New("search failed"),
	}
	service := NewSearchApplicationsService(fakeSearcher)

	_, err := service.Execute(context.Background(), search.ApplicationCriteria{})

	if err == nil {
		t.Fatal("expected search error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestSearchApplicationsServiceRequiresSearcher
//
// Verifies that the search workflow requires a repository boundary.
// -----------------------------------------------------------------------------
func TestSearchApplicationsServiceRequiresSearcher(t *testing.T) {
	service := NewSearchApplicationsService(nil)

	_, err := service.Execute(context.Background(), search.ApplicationCriteria{})

	if err == nil {
		t.Fatal("expected missing searcher to be rejected")
	}
}

// -----------------------------------------------------------------------------
// newApplicationSearchTestApplication
//
// Creates a valid application for application search workflow tests.
// -----------------------------------------------------------------------------
func newApplicationSearchTestApplication(t *testing.T, id domain.ApplicationID, status domain.ApplicationStatus) domain.Application {
	t.Helper()

	application, err := domain.NewApplication(
		id,
		"Backend Developer",
		domain.Company{Name: "Example Studio", Website: "https://example.com"},
		status,
		time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to create test application: %v", err)
	}

	application.Source = "Company site"
	application.Notes = "Backend-focused role."

	return application
}
