package api

import (
	"context"
	"net/http"

	"github.com/mruke/applyby/internal/application"
	"github.com/mruke/applyby/internal/domain"
)

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
// updateReminderExecutor
//
// Defines the application behavior needed by the update reminder handler.
// -----------------------------------------------------------------------------
type updateReminderExecutor interface {
	Execute(ctx context.Context, input application.UpdateReminderInput) (domain.Reminder, error)
}

// -----------------------------------------------------------------------------
// removeReminderExecutor
//
// Defines the application behavior needed by the remove reminder handler.
// -----------------------------------------------------------------------------
type removeReminderExecutor interface {
	Execute(ctx context.Context, input application.RemoveReminderInput) error
}

// -----------------------------------------------------------------------------
// ReminderWorkflowHandlers
//
// Groups HTTP handlers for reminder workflows.
// -----------------------------------------------------------------------------
type ReminderWorkflowHandlers struct {
	scheduleReminder scheduleReminderExecutor
	listReminders    listRemindersExecutor
	completeReminder completeReminderExecutor
	updateReminder   updateReminderExecutor
	removeReminder   removeReminderExecutor
}

// -----------------------------------------------------------------------------
// ReminderWorkflowDependencies
//
// Collects reminder workflow dependencies for API route construction.
// -----------------------------------------------------------------------------
type ReminderWorkflowDependencies struct {
	ScheduleReminder scheduleReminderExecutor
	ListReminders    listRemindersExecutor
	CompleteReminder completeReminderExecutor
	UpdateReminder   updateReminderExecutor
	RemoveReminder   removeReminderExecutor
}

// -----------------------------------------------------------------------------
// NewReminderWorkflowHandlers
//
// Creates reminder workflow handlers from application-layer dependencies.
// -----------------------------------------------------------------------------
func NewReminderWorkflowHandlers(dependencies ReminderWorkflowDependencies) ReminderWorkflowHandlers {
	return ReminderWorkflowHandlers{
		scheduleReminder: dependencies.ScheduleReminder,
		listReminders:    dependencies.ListReminders,
		completeReminder: dependencies.CompleteReminder,
		updateReminder:   dependencies.UpdateReminder,
		removeReminder:   dependencies.RemoveReminder,
	}
}

// -----------------------------------------------------------------------------
// HandleCollection
//
// Routes reminder collection requests for one application by HTTP method.
// -----------------------------------------------------------------------------
func (handlers ReminderWorkflowHandlers) HandleCollection(
	response http.ResponseWriter,
	request *http.Request,
	applicationID domain.ApplicationID,
) {
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
// HandleResource
//
// Routes item-level reminder maintenance requests by method.
// -----------------------------------------------------------------------------
func (handlers ReminderWorkflowHandlers) HandleResource(
	response http.ResponseWriter,
	request *http.Request,
	reminderID domain.ReminderID,
) {
	switch request.Method {
	case http.MethodPatch:
		handlers.handleUpdateReminder(response, request, reminderID)
	case http.MethodDelete:
		handlers.handleRemoveReminder(response, request, reminderID)
	default:
		writeJSON(response, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

// -----------------------------------------------------------------------------
// HandleComplete
//
// Handles reminder completion requests.
// -----------------------------------------------------------------------------
func (handlers ReminderWorkflowHandlers) HandleComplete(response http.ResponseWriter, request *http.Request, reminderID domain.ReminderID) {
	if request.Method != http.MethodPatch {
		writeJSON(response, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	if handlers.completeReminder == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "complete reminder service is not configured"})
		return
	}

	reminder, err := handlers.completeReminder.Execute(request.Context(), application.CompleteReminderInput{
		ID: reminderID,
	})
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusOK, reminderToResponse(reminder))
}

// -----------------------------------------------------------------------------
// handleScheduleReminder
//
// Decodes, validates, and executes the schedule reminder workflow.
// -----------------------------------------------------------------------------
func (handlers ReminderWorkflowHandlers) handleScheduleReminder(response http.ResponseWriter, request *http.Request, applicationID domain.ApplicationID) {
	if handlers.scheduleReminder == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "schedule reminder service is not configured"})
		return
	}

	var body reminderRequest

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
func (handlers ReminderWorkflowHandlers) handleListReminders(response http.ResponseWriter, request *http.Request, applicationID domain.ApplicationID) {
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
// handleUpdateReminder
//
// Decodes, validates, and executes the update reminder workflow.
// -----------------------------------------------------------------------------
func (handlers ReminderWorkflowHandlers) handleUpdateReminder(
	response http.ResponseWriter,
	request *http.Request,
	reminderID domain.ReminderID,
) {
	if handlers.updateReminder == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "update reminder service is not configured"})
		return
	}

	var body reminderRequest

	if err := decodeJSON(request, &body); err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	input, err := body.toUpdateInput(reminderID)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	reminder, err := handlers.updateReminder.Execute(request.Context(), input)
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(response, http.StatusOK, reminderToResponse(reminder))
}

// -----------------------------------------------------------------------------
// handleRemoveReminder
//
// Executes the remove reminder workflow.
// -----------------------------------------------------------------------------
func (handlers ReminderWorkflowHandlers) handleRemoveReminder(
	response http.ResponseWriter,
	request *http.Request,
	reminderID domain.ReminderID,
) {
	if handlers.removeReminder == nil {
		writeJSON(response, http.StatusInternalServerError, errorResponse{Error: "remove reminder service is not configured"})
		return
	}

	err := handlers.removeReminder.Execute(request.Context(), application.RemoveReminderInput{
		ID: reminderID,
	})
	if err != nil {
		writeJSON(response, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
