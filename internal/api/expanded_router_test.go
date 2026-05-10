package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// -----------------------------------------------------------------------------
// TestNewExpandedRouterRoutesApplicationSearch
//
// Verifies that the expanded router sends search requests to workflow handlers.
// -----------------------------------------------------------------------------
func TestNewExpandedRouterRoutesApplicationSearch(t *testing.T) {
	router := NewExpandedRouter(NewApplicationHandlers(nil, nil, nil), NewWorkflowHandlers(WorkflowHandlerDependencies{}))

	request := httptest.NewRequest(http.MethodGet, "/applications/search", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected search route to reach workflow handler, got %d", response.Code)
	}
}

// -----------------------------------------------------------------------------
// TestNewExpandedRouterRoutesApplicationWorkflow
//
// Verifies that the expanded router sends nested application workflow requests correctly.
// -----------------------------------------------------------------------------
func TestNewExpandedRouterRoutesApplicationWorkflow(t *testing.T) {
	router := NewExpandedRouter(NewApplicationHandlers(nil, nil, nil), NewWorkflowHandlers(WorkflowHandlerDependencies{}))

	request := httptest.NewRequest(http.MethodGet, "/applications/app-001/reminders", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected reminder route to reach workflow handler, got %d", response.Code)
	}
}

// -----------------------------------------------------------------------------
// TestNewExpandedRouterFallsBackToApplicationStatus
//
// Verifies that existing status routes still reach the application status handler.
// -----------------------------------------------------------------------------
func TestNewExpandedRouterFallsBackToApplicationStatus(t *testing.T) {
	router := NewExpandedRouter(NewApplicationHandlers(nil, nil, nil), NewWorkflowHandlers(WorkflowHandlerDependencies{}))

	request := httptest.NewRequest(http.MethodPatch, "/applications/app-001/status", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status route to reach existing application handler, got %d", response.Code)
	}
}

// -----------------------------------------------------------------------------
// TestNewExpandedRouterRoutesReminderComplete
//
// Verifies that the expanded router sends reminder completion requests correctly.
// -----------------------------------------------------------------------------
func TestNewExpandedRouterRoutesReminderComplete(t *testing.T) {
	router := NewExpandedRouter(NewApplicationHandlers(nil, nil, nil), NewWorkflowHandlers(WorkflowHandlerDependencies{}))

	request := httptest.NewRequest(http.MethodPatch, "/reminders/rem-001/complete", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected complete reminder route to reach workflow handler, got %d", response.Code)
	}
}
