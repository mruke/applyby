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
// scheduleReminderExecutor
//
// Defines the application behavior needed by the schedule reminder handler.
// -----------------------------------------------------------------------------
type scheduleReminderExecutor interface {
	Execute(ctx context.Context, input application.ScheduleReminderInput) (domain.Reminder, error)
}

// -----------------------------------------------------------------------------
// listRemindersExecutor
//
// Defines the application behavior needed by the list reminders handler.
// -----------------------------------------------------------------------------
type listRemindersExecutor interface {
	Execute(ctx context.Context, input application.ListRemindersInput) ([]domain.Reminder, error)
}

// -----------------------------------------------------------------------------
// completeReminderExecutor
//
// Defines the application behavior needed by the complete reminder handler.
// -----------------------------------------------------------------------------
type completeReminderExecutor interface {
	Execute(ctx context.Context, input application.CompleteReminderInput) (domain.Reminder, error)
}

// -----------------------------------------------------------------------------
// addDocumentExecutor
//
// Defines the application behavior needed by the add document handler.
// -----------------------------------------------------------------------------
type addDocumentExecutor interface {
	Execute(ctx context.Context, input application.AddDocumentInput) (domain.Document, error)
}

// -----------------------------------------------------------------------------
// listDocumentsExecutor
//
// Defines the application behavior needed by the list documents handler.
// -----------------------------------------------------------------------------
type listDocumentsExecutor interface {
	Execute(ctx context.Context, input application.ListDocumentsInput) ([]domain.Document, error)
}

// -----------------------------------------------------------------------------
// WorkflowHandlers
//
// Groups HTTP handlers for expanded application workflow API routes.
// -----------------------------------------------------------------------------
type WorkflowHandlers struct {
	searchApplications searchApplicationsExecutor
	listActivityEvents listActivityEventsExecutor
	scheduleReminder   scheduleReminderExecutor
	listReminders      listRemindersExecutor
	completeReminder   completeReminderExecutor
	contacts           ContactWorkflowHandlers
	addDocument        addDocumentExecutor
	listDocuments      listDocumentsExecutor
}

// -----------------------------------------------------------------------------
// WorkflowHandlerDependencies
//
// Collects expanded workflow dependencies for API route construction.
// -----------------------------------------------------------------------------
type WorkflowHandlerDependencies struct {
	SearchApplications searchApplicationsExecutor
	ListActivityEvents listActivityEventsExecutor
	ScheduleReminder   scheduleReminderExecutor
	ListReminders      listRemindersExecutor
	CompleteReminder   completeReminderExecutor
	Contacts           ContactWorkflowDependencies
	AddDocument        addDocumentExecutor
	ListDocuments      listDocumentsExecutor
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
		scheduleReminder:   dependencies.ScheduleReminder,
		listReminders:      dependencies.ListReminders,
		completeReminder:   dependencies.CompleteReminder,
		contacts:           NewContactWorkflowHandlers(dependencies.Contacts),
		addDocument:        dependencies.AddDocument,
		listDocuments:      dependencies.ListDocuments,
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
// HandleReminderComplete
//
// Handles reminder completion requests.
// -----------------------------------------------------------------------------
func (handlers WorkflowHandlers) HandleReminderComplete(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPatch {
		writeJSON(response, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	if handlers.completeReminder == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "complete reminder service is not configured"})
		return
	}

	id, ok := reminderIDFromCompletePath(request.URL.Path)
	if !ok {
		writeJSON(response, http.StatusNotFound, errorResponse{Error: "route not found"})
		return
	}

	reminder, err := handlers.completeReminder.Execute(request.Context(), application.CompleteReminderInput{
		ID: id,
	})
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusOK, reminderToResponse(reminder))
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

	applicationID, resource, ok := applicationWorkflowPathParts(request.URL.Path)
	if !ok {
		return false
	}

	switch resource {
	case "activity":
		handlers.handleActivityEvents(response, request, applicationID)
	case "reminders":
		handlers.handleReminders(response, request, applicationID)
	case "contacts":
		handlers.contacts.HandleCollection(response, request, applicationID)
	case "documents":
		handlers.handleDocuments(response, request, applicationID)
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
// handleReminders
//
// Handles reminder collection requests for one application.
// -----------------------------------------------------------------------------
func (handlers WorkflowHandlers) handleReminders(response http.ResponseWriter, request *http.Request, applicationID domain.ApplicationID) {
	switch request.Method {
	case http.MethodGet:
		handlers.handleListReminders(response, request, applicationID)
	case http.MethodPost:
		handlers.handleScheduleReminder(response, request, applicationID)
	default:
		writeJSON(response, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

// -----------------------------------------------------------------------------
// handleScheduleReminder
//
// Decodes, validates, and executes the schedule reminder workflow.
// -----------------------------------------------------------------------------
func (handlers WorkflowHandlers) handleScheduleReminder(response http.ResponseWriter, request *http.Request, applicationID domain.ApplicationID) {
	if handlers.scheduleReminder == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "schedule reminder service is not configured"})
		return
	}

	var body scheduleReminderRequest

	if err := decodeJSON(request, &body); err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	input, err := body.toInput(applicationID)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	reminder, err := handlers.scheduleReminder.Execute(request.Context(), input)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusCreated, reminderToResponse(reminder))
}

// -----------------------------------------------------------------------------
// handleListReminders
//
// Executes the list reminders workflow for one application.
// -----------------------------------------------------------------------------
func (handlers WorkflowHandlers) handleListReminders(response http.ResponseWriter, request *http.Request, applicationID domain.ApplicationID) {
	if handlers.listReminders == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "list reminders service is not configured"})
		return
	}

	reminders, err := handlers.listReminders.Execute(request.Context(), application.ListRemindersInput{
		ApplicationID: applicationID,
	})
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusOK, remindersToResponse(reminders))
}

// -----------------------------------------------------------------------------
// handleDocuments
//
// Handles document metadata collection requests for one application.
// -----------------------------------------------------------------------------
func (handlers WorkflowHandlers) handleDocuments(response http.ResponseWriter, request *http.Request, applicationID domain.ApplicationID) {
	switch request.Method {
	case http.MethodGet:
		handlers.handleListDocuments(response, request, applicationID)
	case http.MethodPost:
		handlers.handleAddDocument(response, request, applicationID)
	default:
		writeJSON(response, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

// -----------------------------------------------------------------------------
// handleAddDocument
//
// Decodes, validates, and executes the add document metadata workflow.
// -----------------------------------------------------------------------------
func (handlers WorkflowHandlers) handleAddDocument(response http.ResponseWriter, request *http.Request, applicationID domain.ApplicationID) {
	if handlers.addDocument == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "add document service is not configured"})
		return
	}

	var body documentRequest

	if err := decodeJSON(request, &body); err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	document, err := handlers.addDocument.Execute(request.Context(), body.toInput(applicationID))
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusCreated, documentToResponse(document))
}

// -----------------------------------------------------------------------------
// handleListDocuments
//
// Executes the list documents workflow for one application.
// -----------------------------------------------------------------------------
func (handlers WorkflowHandlers) handleListDocuments(response http.ResponseWriter, request *http.Request, applicationID domain.ApplicationID) {
	if handlers.listDocuments == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "list documents service is not configured"})
		return
	}

	documents, err := handlers.listDocuments.Execute(request.Context(), application.ListDocumentsInput{
		ApplicationID: applicationID,
	})
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusOK, documentsToResponse(documents))
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
