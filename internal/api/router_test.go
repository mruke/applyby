package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// -----------------------------------------------------------------------------
// TestNewRouterRoutesApplicationsCollection
//
// Verifies that the router sends collection requests to the application handler.
// -----------------------------------------------------------------------------
func TestNewRouterRoutesApplicationsCollection(t *testing.T) {
	router := NewRouter(NewApplicationHandlers(nil, nil, nil, nil))

	request := httptest.NewRequest(http.MethodGet, "/applications", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected applications route to reach handler, got %d", response.Code)
	}
}

// -----------------------------------------------------------------------------
// TestNewRouterRoutesApplicationDetail
//
// Verifies that the router sends detail requests to the application resource handler.
// -----------------------------------------------------------------------------
func TestNewRouterRoutesApplicationDetail(t *testing.T) {
	router := NewRouter(NewApplicationHandlers(nil, nil, nil, nil))

	request := httptest.NewRequest(http.MethodGet, "/applications/app-001", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected detail route to reach handler, got %d", response.Code)
	}
}

// -----------------------------------------------------------------------------
// TestNewRouterRoutesApplicationStatus
//
// Verifies that the router sends status update requests to the status handler.
// -----------------------------------------------------------------------------
func TestNewRouterRoutesApplicationStatus(t *testing.T) {
	router := NewRouter(NewApplicationHandlers(nil, nil, nil, nil))

	request := httptest.NewRequest(http.MethodPatch, "/applications/app-001/status", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status route to reach handler, got %d", response.Code)
	}
}
