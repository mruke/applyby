package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/mruke/applyby/internal/application"
	"github.com/mruke/applyby/internal/domain"
	"github.com/mruke/applyby/internal/search"
)

// -----------------------------------------------------------------------------
// searchApplicationsExecutor
//
// Defines the application behavior needed by the search applications handler.
// -----------------------------------------------------------------------------
type searchApplicationsExecutor interface {
	Execute(ctx context.Context, criteria search.ApplicationCriteria) ([]domain.Application, error)
}

// -----------------------------------------------------------------------------
// listActivityEventsExecutor
//
// Defines the application behavior needed by the activity timeline handler.
// -----------------------------------------------------------------------------
type listActivityEventsExecutor interface {
	Execute(ctx context.Context, input application.ListActivityEventsInput) ([]domain.ActivityEvent, error)
}

// -----------------------------------------------------------------------------
// WorkflowHandlers
//
// Groups HTTP handlers for expanded application workflow API routes.
// -----------------------------------------------------------------------------
type WorkflowHandlers struct {
	searchApplications searchApplicationsExecutor
	listActivityEvents listActivityEventsExecutor
	reminders          ReminderWorkflowHandlers
	contacts           ContactWorkflowHandlers
	documents          DocumentWorkflowHandlers
}

// -----------------------------------------------------------------------------
// WorkflowHandlerDependencies
//
// Collects expanded workflow dependencies for API route construction.
// -----------------------------------------------------------------------------
type WorkflowHandlerDependencies struct {
	SearchApplications searchApplicationsExecutor
	ListActivityEvents listActivityEventsExecutor
	Reminders          ReminderWorkflowDependencies
	Contacts           ContactWorkflowDependencies
	Documents          DocumentWorkflowDependencies
}

// -----------------------------------------------------------------------------
// NewWorkflowHandlers
//
// Creates expanded workflow handlers from application-layer dependencies.
// -----------------------------------------------------------------------------
func NewWorkflowHandlers(dependencies WorkflowHandlerDependencies) WorkflowHandlers {
	return WorkflowHandlers{
		searchApplications: dependencies.SearchApplications,
		listActivityEvents: dependencies.ListActivityEvents,
		reminders:          NewReminderWorkflowHandlers(dependencies.Reminders),
		contacts:           NewContactWorkflowHandlers(dependencies.Contacts),
		documents:          NewDocumentWorkflowHandlers(dependencies.Documents),
	}
}

// -----------------------------------------------------------------------------
// HandleApplicationSearch
//
// Handles application search requests.
// -----------------------------------------------------------------------------
func (handlers WorkflowHandlers) HandleApplicationSearch(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeJSON(response, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	if handlers.searchApplications == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "search applications service is not configured"})
		return
	}

	criteria, err := searchCriteriaFromRequest(request)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	applications, err := handlers.searchApplications.Execute(request.Context(), criteria)
	if err != nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusOK, applicationsToResponse(applications))
}

// -----------------------------------------------------------------------------
// HandleReminderResource
//
// Routes reminder item-level workflow requests.
// -----------------------------------------------------------------------------
func (handlers WorkflowHandlers) HandleReminderResource(response http.ResponseWriter, request *http.Request) {
	if reminderID, ok := reminderIDFromCompletePath(request.URL.Path); ok {
		handlers.reminders.HandleComplete(response, request, reminderID)
		return
	}

	if reminderID, ok := reminderIDFromResourcePath(request.URL.Path); ok {
		handlers.reminders.HandleResource(response, request, reminderID)
		return
	}

	writeJSON(response, http.StatusNotFound, errorResponse{Error: "route not found"})
}

// -----------------------------------------------------------------------------
// HandleReminderComplete
//
// Handles reminder completion requests.
//
// Deprecated: route through HandleReminderResource.
// -----------------------------------------------------------------------------
func (handlers WorkflowHandlers) HandleReminderComplete(response http.ResponseWriter, request *http.Request) {
	if reminderID, ok := reminderIDFromCompletePath(request.URL.Path); ok {
		handlers.reminders.HandleComplete(response, request, reminderID)
		return
	}

	writeJSON(response, http.StatusNotFound, errorResponse{Error: "route not found"})
}

// -----------------------------------------------------------------------------
// HandleApplicationWorkflow
//
// Routes application-owned workflow requests by path and method.
// -----------------------------------------------------------------------------
func (handlers WorkflowHandlers) HandleApplicationWorkflow(response http.ResponseWriter, request *http.Request) bool {
	if applicationID, contactID, ok := applicationContactResourcePathParts(request.URL.Path); ok {
		handlers.contacts.HandleResource(response, request, applicationID, contactID)
		return true
	}

	if applicationID, documentID, ok := applicationDocumentResourcePathParts(request.URL.Path); ok {
		handlers.documents.HandleResource(response, request, applicationID, documentID)
		return true
	}

	applicationID, resource, ok := applicationWorkflowPathParts(request.URL.Path)
	if !ok {
		return false
	}

	switch resource {
	case "activity":
		handlers.handleActivityEvents(response, request, applicationID)
	case "reminders":
		handlers.reminders.HandleCollection(response, request, applicationID)
	case "contacts":
		handlers.contacts.HandleCollection(response, request, applicationID)
	case "documents":
		handlers.documents.HandleCollection(response, request, applicationID)
	default:
		return false
	}

	return true
}

// -----------------------------------------------------------------------------
// handleActivityEvents
//
// Handles activity timeline requests for one application.
// -----------------------------------------------------------------------------
func (handlers WorkflowHandlers) handleActivityEvents(response http.ResponseWriter, request *http.Request, applicationID domain.ApplicationID) {
	if request.Method != http.MethodGet {
		writeJSON(response, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	if handlers.listActivityEvents == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "list activity events service is not configured"})
		return
	}

	events, err := handlers.listActivityEvents.Execute(request.Context(), application.ListActivityEventsInput{
		ApplicationID: applicationID,
	})
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusOK, activityEventsToResponse(events))
}

// -----------------------------------------------------------------------------
// applicationWorkflowPathParts
//
// Extracts the application identity and nested resource from an application route.
// -----------------------------------------------------------------------------
func applicationWorkflowPathParts(path string) (domain.ApplicationID, string, bool) {
	if !strings.HasPrefix(path, "/applications/") {
		return "", "", false
	}

	trimmedPath := strings.TrimPrefix(path, "/applications/")
	parts := strings.Split(strings.Trim(trimmedPath, "/"), "/")

	if len(parts) != 2 {
		return "", "", false
	}

	applicationID, err := domain.NewApplicationID(parts[0])
	if err != nil {
		return "", "", false
	}

	return applicationID, parts[1], true
}

// -----------------------------------------------------------------------------
// applicationContactResourcePathParts
//
// Extracts application and contact identities from a contact maintenance route.
// -----------------------------------------------------------------------------
func applicationContactResourcePathParts(path string) (domain.ApplicationID, domain.ContactID, bool) {
	if !strings.HasPrefix(path, "/applications/") {
		return "", "", false
	}

	trimmedPath := strings.TrimPrefix(path, "/applications/")
	parts := strings.Split(strings.Trim(trimmedPath, "/"), "/")

	if len(parts) != 3 || parts[1] != "contacts" {
		return "", "", false
	}

	applicationID, err := domain.NewApplicationID(parts[0])
	if err != nil {
		return "", "", false
	}

	contactID := domain.ContactID(parts[2])
	if err := contactID.Validate(); err != nil {
		return "", "", false
	}

	return applicationID, contactID, true
}

// -----------------------------------------------------------------------------
// applicationDocumentResourcePathParts
//
// Extracts application and document identities from a document maintenance route.
// -----------------------------------------------------------------------------
func applicationDocumentResourcePathParts(path string) (domain.ApplicationID, domain.DocumentID, bool) {
	if !strings.HasPrefix(path, "/applications/") {
		return "", "", false
	}

	trimmedPath := strings.TrimPrefix(path, "/applications/")
	parts := strings.Split(strings.Trim(trimmedPath, "/"), "/")

	if len(parts) != 3 || parts[1] != "documents" {
		return "", "", false
	}

	applicationID, err := domain.NewApplicationID(parts[0])
	if err != nil {
		return "", "", false
	}

	documentID := domain.DocumentID(parts[2])
	if err := documentID.Validate(); err != nil {
		return "", "", false
	}

	return applicationID, documentID, true
}

// -----------------------------------------------------------------------------
// reminderIDFromCompletePath
//
// Extracts the reminder identity from the complete reminder route path.
// -----------------------------------------------------------------------------
func reminderIDFromCompletePath(path string) (domain.ReminderID, bool) {
	if !strings.HasPrefix(path, "/reminders/") || !strings.HasSuffix(path, "/complete") {
		return "", false
	}

	trimmedPath := strings.TrimPrefix(path, "/reminders/")
	rawID := strings.TrimSuffix(trimmedPath, "/complete")
	rawID = strings.Trim(rawID, "/")

	id := domain.ReminderID(rawID)

	if err := id.Validate(); err != nil {
		return "", false
	}

	return id, true
}

// -----------------------------------------------------------------------------
// reminderIDFromResourcePath
//
// Extracts the reminder identity from an item-level reminder route path.
// -----------------------------------------------------------------------------
func reminderIDFromResourcePath(path string) (domain.ReminderID, bool) {
	if !strings.HasPrefix(path, "/reminders/") || strings.HasSuffix(path, "/complete") {
		return "", false
	}

	rawID := strings.TrimPrefix(path, "/reminders/")
	rawID = strings.Trim(rawID, "/")

	if rawID == "" || strings.Contains(rawID, "/") {
		return "", false
	}

	id := domain.ReminderID(rawID)

	if err := id.Validate(); err != nil {
		return "", false
	}

	return id, true
}
