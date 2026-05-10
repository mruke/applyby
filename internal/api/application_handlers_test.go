package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/application"
	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// fakeCreateApplicationExecutor
//
// Provides a fake create workflow for API handler tests.
// -----------------------------------------------------------------------------
type fakeCreateApplicationExecutor struct {
	application domain.Application
	err         error
	called      bool
}

// -----------------------------------------------------------------------------
// Execute
//
// Records create workflow input and returns the configured fake result.
// -----------------------------------------------------------------------------
func (executor *fakeCreateApplicationExecutor) Execute(ctx context.Context, input application.CreateApplicationInput) (domain.Application, error) {
	executor.called = true

	if executor.err != nil {
		return domain.Application{}, executor.err
	}

	return executor.application, nil
}

// -----------------------------------------------------------------------------
// fakeListApplicationsExecutor
//
// Provides a fake list workflow for API handler tests.
// -----------------------------------------------------------------------------
type fakeListApplicationsExecutor struct {
	applications []domain.Application
	err          error
	called       bool
}

// -----------------------------------------------------------------------------
// Execute
//
// Records list workflow execution and returns the configured fake result.
// -----------------------------------------------------------------------------
func (executor *fakeListApplicationsExecutor) Execute(ctx context.Context) ([]domain.Application, error) {
	executor.called = true

	if executor.err != nil {
		return nil, executor.err
	}

	return executor.applications, nil
}

// -----------------------------------------------------------------------------
// fakeUpdateApplicationStatusExecutor
//
// Provides a fake update status workflow for API handler tests.
// -----------------------------------------------------------------------------
type fakeUpdateApplicationStatusExecutor struct {
	application domain.Application
	err         error
	called      bool
}

// -----------------------------------------------------------------------------
// Execute
//
// Records update status workflow execution and returns the configured fake result.
// -----------------------------------------------------------------------------
func (executor *fakeUpdateApplicationStatusExecutor) Execute(ctx context.Context, input application.UpdateApplicationStatusInput) (domain.Application, error) {
	executor.called = true

	if executor.err != nil {
		return domain.Application{}, executor.err
	}

	return executor.application, nil
}

// -----------------------------------------------------------------------------
// TestHandleApplicationsCreatesApplication
//
// Verifies that POST /applications executes the create workflow.
// -----------------------------------------------------------------------------
func TestHandleApplicationsCreatesApplication(t *testing.T) {
	createdApplication := newAPIHandlerTestApplication(t, "app-001", domain.StatusApplied)
	createExecutor := &fakeCreateApplicationExecutor{application: createdApplication}
	handlers := NewApplicationHandlers(createExecutor, nil, nil)

	requestBody := `{
        "id": "app-001",
        "title": "Backend Developer",
        "company_name": "Example Studio",
        "company_website": "https://example.com",
        "status": "applied",
        "source": "Company site",
        "notes": "Applied with backend resume.",
        "created_at": "2026-05-10T08:00:00Z"
    }`

	request := httptest.NewRequest(http.MethodPost, "/applications", strings.NewReader(requestBody))
	response := httptest.NewRecorder()

	handlers.HandleApplications(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	if !createExecutor.called {
		t.Fatal("expected create workflow to be called")
	}
}

// -----------------------------------------------------------------------------
// TestHandleApplicationsRejectsInvalidCreateJSON
//
// Verifies that POST /applications rejects invalid JSON.
// -----------------------------------------------------------------------------
func TestHandleApplicationsRejectsInvalidCreateJSON(t *testing.T) {
	handlers := NewApplicationHandlers(&fakeCreateApplicationExecutor{}, nil, nil)
	request := httptest.NewRequest(http.MethodPost, "/applications", bytes.NewBufferString("{"))
	response := httptest.NewRecorder()

	handlers.HandleApplications(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

// -----------------------------------------------------------------------------
// TestHandleApplicationsListsApplications
//
// Verifies that GET /applications executes the list workflow.
// -----------------------------------------------------------------------------
func TestHandleApplicationsListsApplications(t *testing.T) {
	application := newAPIHandlerTestApplication(t, "app-001", domain.StatusApplied)
	listExecutor := &fakeListApplicationsExecutor{applications: []domain.Application{application}}
	handlers := NewApplicationHandlers(nil, listExecutor, nil)

	request := httptest.NewRequest(http.MethodGet, "/applications", nil)
	response := httptest.NewRecorder()

	handlers.HandleApplications(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if !listExecutor.called {
		t.Fatal("expected list workflow to be called")
	}
}

// -----------------------------------------------------------------------------
// TestHandleApplicationsRejectsUnsupportedMethod
//
// Verifies that collection handler rejects unsupported HTTP methods.
// -----------------------------------------------------------------------------
func TestHandleApplicationsRejectsUnsupportedMethod(t *testing.T) {
	handlers := NewApplicationHandlers(nil, nil, nil)
	request := httptest.NewRequest(http.MethodDelete, "/applications", nil)
	response := httptest.NewRecorder()

	handlers.HandleApplications(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, response.Code)
	}
}

// -----------------------------------------------------------------------------
// TestHandleApplicationStatusUpdatesStatus
//
// Verifies that PATCH /applications/{id}/status executes the update workflow.
// -----------------------------------------------------------------------------
func TestHandleApplicationStatusUpdatesStatus(t *testing.T) {
	updatedApplication := newAPIHandlerTestApplication(t, "app-001", domain.StatusInterviewing)
	updateExecutor := &fakeUpdateApplicationStatusExecutor{application: updatedApplication}
	handlers := NewApplicationHandlers(nil, nil, updateExecutor)

	request := httptest.NewRequest(http.MethodPatch, "/applications/app-001/status", strings.NewReader(`{"status":"interviewing"}`))
	response := httptest.NewRecorder()

	handlers.HandleApplicationStatus(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if !updateExecutor.called {
		t.Fatal("expected update status workflow to be called")
	}
}

// -----------------------------------------------------------------------------
// TestHandleApplicationStatusRejectsInvalidRoute
//
// Verifies that invalid status update paths return not found.
// -----------------------------------------------------------------------------
func TestHandleApplicationStatusRejectsInvalidRoute(t *testing.T) {
	handlers := NewApplicationHandlers(nil, nil, &fakeUpdateApplicationStatusExecutor{})

	request := httptest.NewRequest(http.MethodPatch, "/applications/app-001", strings.NewReader(`{"status":"interviewing"}`))
	response := httptest.NewRecorder()

	handlers.HandleApplicationStatus(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

// -----------------------------------------------------------------------------
// TestHandleApplicationStatusRejectsWorkflowError
//
// Verifies that update workflow errors return a bad request response.
// -----------------------------------------------------------------------------
func TestHandleApplicationStatusRejectsWorkflowError(t *testing.T) {
	updateExecutor := &fakeUpdateApplicationStatusExecutor{err: errors.New("invalid transition")}
	handlers := NewApplicationHandlers(nil, nil, updateExecutor)

	request := httptest.NewRequest(http.MethodPatch, "/applications/app-001/status", strings.NewReader(`{"status":"interviewing"}`))
	response := httptest.NewRecorder()

	handlers.HandleApplicationStatus(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

// -----------------------------------------------------------------------------
// newAPIHandlerTestApplication
//
// Creates a valid application for API handler tests.
// -----------------------------------------------------------------------------
func newAPIHandlerTestApplication(t *testing.T, id domain.ApplicationID, status domain.ApplicationStatus) domain.Application {
	t.Helper()

	createdAt := time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC)

	application, err := domain.NewApplication(
		id,
		"Backend Developer",
		domain.Company{Name: "Example Studio", Website: "https://example.com"},
		status,
		createdAt,
	)
	if err != nil {
		t.Fatalf("failed to create test application: %v", err)
	}

	return application
}
