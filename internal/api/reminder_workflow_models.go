package api

import (
	"time"

	"github.com/mruke/applyby/internal/application"
	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// reminderRequest
//
// Represents the JSON request body for scheduling or updating a reminder.
// -----------------------------------------------------------------------------
type reminderRequest struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	DueAt string `json:"due_at"`
}

// -----------------------------------------------------------------------------
// reminderResponse
//
// Represents the JSON response shape for a reminder.
// -----------------------------------------------------------------------------
type reminderResponse struct {
	ID            string `json:"id"`
	ApplicationID string `json:"application_id"`
	Title         string `json:"title"`
	DueAt         string `json:"due_at"`
	Completed     bool   `json:"completed"`
}

// -----------------------------------------------------------------------------
// remindersResponse
//
// Represents the JSON response shape for a reminder collection.
// -----------------------------------------------------------------------------
type remindersResponse struct {
	Reminders []reminderResponse `json:"reminders"`
}

// -----------------------------------------------------------------------------
// toInput
//
// Converts a reminder request into a schedule reminder workflow input model.
// -----------------------------------------------------------------------------
func (request reminderRequest) toInput(applicationID domain.ApplicationID) (application.ScheduleReminderInput, error) {
	reminderID := domain.ReminderID(request.ID)

	if err := reminderID.Validate(); err != nil {
		return application.ScheduleReminderInput{}, err
	}

	dueAt, err := parseRequiredRFC3339Time("due_at", request.DueAt)
	if err != nil {
		return application.ScheduleReminderInput{}, err
	}

	return application.ScheduleReminderInput{
		ID:            reminderID,
		ApplicationID: applicationID,
		Title:         request.Title,
		DueAt:         dueAt,
	}, nil
}

// -----------------------------------------------------------------------------
// toUpdateInput
//
// Converts a reminder request into an update reminder workflow input model.
// -----------------------------------------------------------------------------
func (request reminderRequest) toUpdateInput(reminderID domain.ReminderID) (application.UpdateReminderInput, error) {
	dueAt, err := parseRequiredRFC3339Time("due_at", request.DueAt)
	if err != nil {
		return application.UpdateReminderInput{}, err
	}

	return application.UpdateReminderInput{
		ID:    reminderID,
		Title: request.Title,
		DueAt: dueAt,
	}, nil
}

// -----------------------------------------------------------------------------
// reminderToResponse
//
// Converts a domain reminder into an API response model.
// -----------------------------------------------------------------------------
func reminderToResponse(reminder domain.Reminder) reminderResponse {
	return reminderResponse{
		ID:            reminder.ID.String(),
		ApplicationID: reminder.ApplicationID.String(),
		Title:         reminder.Title,
		DueAt:         reminder.DueAt.Format(time.RFC3339),
		Completed:     reminder.Completed,
	}
}

// -----------------------------------------------------------------------------
// remindersToResponse
//
// Converts domain reminders into an API collection response model.
// -----------------------------------------------------------------------------
func remindersToResponse(reminders []domain.Reminder) remindersResponse {
	response := remindersResponse{
		Reminders: make([]reminderResponse, 0, len(reminders)),
	}

	for _, reminder := range reminders {
		response.Reminders = append(response.Reminders, reminderToResponse(reminder))
	}

	return response
}
