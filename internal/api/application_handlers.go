package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mruke/applyby/internal/application"
	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// createApplicationExecutor
//
// Defines the application behavior needed by the create application handler.
// -----------------------------------------------------------------------------
type createApplicationExecutor interface {
	Execute(ctx context.Context, input application.CreateApplicationInput) (domain.Application, error)
}

// -----------------------------------------------------------------------------
// listApplicationsExecutor
//
// Defines the application behavior needed by the list applications handler.
// -----------------------------------------------------------------------------
type listApplicationsExecutor interface {
	Execute(ctx context.Context) ([]domain.Application, error)
}

// -----------------------------------------------------------------------------
// getApplicationExecutor
//
// Defines the application behavior needed by the get application handler.
// -----------------------------------------------------------------------------
type getApplicationExecutor interface {
	Execute(ctx context.Context, input application.GetApplicationInput) (domain.Application, error)
}

// -----------------------------------------------------------------------------
// updateApplicationStatusExecutor
//
// Defines the application behavior needed by the update status handler.
// -----------------------------------------------------------------------------
type updateApplicationStatusExecutor interface {
	Execute(ctx context.Context, input application.UpdateApplicationStatusInput) (domain.Application, error)
}

// -----------------------------------------------------------------------------
// ApplicationHandlers
//
// Groups HTTP handlers for application-related API workflows.
// -----------------------------------------------------------------------------
type ApplicationHandlers struct {
	createApplication       createApplicationExecutor
	listApplications        listApplicationsExecutor
	getApplication          getApplicationExecutor
	updateApplicationStatus updateApplicationStatusExecutor
}

// -----------------------------------------------------------------------------
// NewApplicationHandlers
//
// Creates application handlers from application-layer workflow dependencies.
// -----------------------------------------------------------------------------
func NewApplicationHandlers(
	createApplication createApplicationExecutor,
	listApplications listApplicationsExecutor,
	getApplication getApplicationExecutor,
	updateApplicationStatus updateApplicationStatusExecutor,
) ApplicationHandlers {
	return ApplicationHandlers{
		createApplication:       createApplication,
		listApplications:        listApplications,
		getApplication:          getApplication,
		updateApplicationStatus: updateApplicationStatus,
	}
}

// -----------------------------------------------------------------------------
// HandleApplications
//
// Routes collection-level application requests by HTTP method.
// -----------------------------------------------------------------------------
func (handlers ApplicationHandlers) HandleApplications(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		handlers.handleListApplications(response, request)
	case http.MethodPost:
		handlers.handleCreateApplication(response, request)
	default:
		writeJSON(response, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

// -----------------------------------------------------------------------------
// HandleApplicationResource
//
// Routes item-level application requests by path and HTTP method.
// -----------------------------------------------------------------------------
func (handlers ApplicationHandlers) HandleApplicationResource(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		id, ok := applicationIDFromDetailPath(request.URL.Path)
		if !ok {
			writeJSON(response, http.StatusNotFound, errorResponse{Error: "route not found"})
			return
		}

		handlers.handleGetApplication(response, request, id)
	case http.MethodPatch:
		id, ok := applicationIDFromStatusPath(request.URL.Path)
		if !ok {
			writeJSON(response, http.StatusNotFound, errorResponse{Error: "route not found"})
			return
		}

		handlers.handleUpdateApplicationStatus(response, request, id)
	default:
		writeJSON(response, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

// -----------------------------------------------------------------------------
// HandleApplicationStatus
//
// Handles status update requests for a single application.
// -----------------------------------------------------------------------------
func (handlers ApplicationHandlers) HandleApplicationStatus(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPatch {
		writeJSON(response, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	id, ok := applicationIDFromStatusPath(request.URL.Path)
	if !ok {
		writeJSON(response, http.StatusNotFound, errorResponse{Error: "route not found"})
		return
	}

	handlers.handleUpdateApplicationStatus(response, request, id)
}

// -----------------------------------------------------------------------------
// handleCreateApplication
//
// Decodes, validates, and executes the create application workflow.
// -----------------------------------------------------------------------------
func (handlers ApplicationHandlers) handleCreateApplication(response http.ResponseWriter, request *http.Request) {
	if handlers.createApplication == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "create application service is not configured"})
		return
	}

	var body createApplicationRequest

	if err := decodeJSON(request, &body); err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	input, err := body.toInput()
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	createdApplication, err := handlers.createApplication.Execute(request.Context(), input)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusCreated, applicationToResponse(createdApplication))
}

// -----------------------------------------------------------------------------
// handleListApplications
//
// Executes the list applications workflow and returns application responses.
// -----------------------------------------------------------------------------
func (handlers ApplicationHandlers) handleListApplications(response http.ResponseWriter, request *http.Request) {
	if handlers.listApplications == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "list applications service is not configured"})
		return
	}

	applications, err := handlers.listApplications.Execute(request.Context())
	if err != nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusOK, applicationsToResponse(applications))
}

// -----------------------------------------------------------------------------
// handleGetApplication
//
// Executes the get application workflow and returns one application response.
// -----------------------------------------------------------------------------
func (handlers ApplicationHandlers) handleGetApplication(response http.ResponseWriter, request *http.Request, id domain.ApplicationID) {
	if handlers.getApplication == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "get application service is not configured"})
		return
	}

	foundApplication, err := handlers.getApplication.Execute(request.Context(), application.GetApplicationInput{
		ID: id,
	})
	if err != nil {
		writeJSON(response, http.StatusNotFound, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusOK, applicationToResponse(foundApplication))
}

// -----------------------------------------------------------------------------
// handleUpdateApplicationStatus
//
// Decodes, validates, and executes the update application status workflow.
// -----------------------------------------------------------------------------
func (handlers ApplicationHandlers) handleUpdateApplicationStatus(response http.ResponseWriter, request *http.Request, id domain.ApplicationID) {
	if handlers.updateApplicationStatus == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "update application status service is not configured"})
		return
	}

	var body updateApplicationStatusRequest

	if err := decodeJSON(request, &body); err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	status, err := body.toStatus()
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	updatedApplication, err := handlers.updateApplicationStatus.Execute(request.Context(), application.UpdateApplicationStatusInput{
		ID:     id,
		Status: status,
	})
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusOK, applicationToResponse(updatedApplication))
}

// -----------------------------------------------------------------------------
// applicationIDFromDetailPath
//
// Extracts the application identity from the detail route path.
// -----------------------------------------------------------------------------
func applicationIDFromDetailPath(path string) (domain.ApplicationID, bool) {
	if !strings.HasPrefix(path, "/applications/") {
		return "", false
	}

	trimmedPath := strings.TrimPrefix(path, "/applications/")
	parts := strings.Split(strings.Trim(trimmedPath, "/"), "/")

	if len(parts) != 1 {
		return "", false
	}

	id, err := domain.NewApplicationID(parts[0])
	if err != nil {
		return "", false
	}

	return id, true
}

// -----------------------------------------------------------------------------
// applicationIDFromStatusPath
//
// Extracts the application identity from the status update route path.
// -----------------------------------------------------------------------------
func applicationIDFromStatusPath(path string) (domain.ApplicationID, bool) {
	if !strings.HasPrefix(path, "/applications/") || !strings.HasSuffix(path, "/status") {
		return "", false
	}

	trimmedPath := strings.TrimPrefix(path, "/applications/")
	rawID := strings.TrimSuffix(trimmedPath, "/status")
	rawID = strings.Trim(rawID, "/")

	id, err := domain.NewApplicationID(rawID)
	if err != nil {
		return "", false
	}

	return id, true
}

// -----------------------------------------------------------------------------
// decodeJSON
//
// Decodes a JSON request body into a destination model.
// -----------------------------------------------------------------------------
func decodeJSON(request *http.Request, destination any) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("invalid JSON request body: %w", err)
	}

	return nil
}

// -----------------------------------------------------------------------------
// writeJSON
//
// Writes a JSON response with the provided status code.
// -----------------------------------------------------------------------------
func writeJSON(response http.ResponseWriter, statusCode int, body any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)

	if err := json.NewEncoder(response).Encode(body); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}
}
