package api

import (
	"net/http"
	"time"

	"github.com/mruke/applyby/internal/application"
	"github.com/mruke/applyby/internal/domain"
	"github.com/mruke/applyby/internal/search"
)

// -----------------------------------------------------------------------------
// scheduleReminderRequest
//
// Represents the JSON request body for scheduling a reminder.
// -----------------------------------------------------------------------------
type scheduleReminderRequest struct {
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
// activityEventResponse
//
// Represents the JSON response shape for an activity event.
// -----------------------------------------------------------------------------
type activityEventResponse struct {
	ApplicationID string `json:"application_id"`
	Type          string `json:"type"`
	OccurredAt    string `json:"occurred_at"`
	Description   string `json:"description"`
}

// -----------------------------------------------------------------------------
// activityEventsResponse
//
// Represents the JSON response shape for an activity event collection.
// -----------------------------------------------------------------------------
type activityEventsResponse struct {
	ActivityEvents []activityEventResponse `json:"activity_events"`
}

// -----------------------------------------------------------------------------
// toInput
//
// Converts a schedule reminder request into an application-layer input model.
// -----------------------------------------------------------------------------
func (request scheduleReminderRequest) toInput(applicationID domain.ApplicationID) (application.ScheduleReminderInput, error) {
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
// searchCriteriaFromRequest
//
// Converts search query parameters into explicit application search criteria.
// -----------------------------------------------------------------------------
func searchCriteriaFromRequest(request *http.Request) (search.ApplicationCriteria, error) {
	query := request.URL.Query()
	criteria := search.ApplicationCriteria{
		CompanyName: query.Get("company_name"),
		Source:      query.Get("source"),
		Text:        query.Get("text"),
	}

	for _, rawStatus := range query["status"] {
		status, err := domain.ParseApplicationStatus(rawStatus)
		if err != nil {
			return search.ApplicationCriteria{}, err
		}

		criteria.Statuses = append(criteria.Statuses, status)
	}

	if createdFrom := query.Get("created_from"); createdFrom != "" {
		parsedTime, err := time.Parse(time.RFC3339, createdFrom)
		if err != nil {
			return search.ApplicationCriteria{}, err
		}

		criteria.CreatedFrom = &parsedTime
	}

	if createdTo := query.Get("created_to"); createdTo != "" {
		parsedTime, err := time.Parse(time.RFC3339, createdTo)
		if err != nil {
			return search.ApplicationCriteria{}, err
		}

		criteria.CreatedTo = &parsedTime
	}

	criteria = criteria.Normalize()

	if err := criteria.Validate(); err != nil {
		return search.ApplicationCriteria{}, err
	}

	return criteria, nil
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

// -----------------------------------------------------------------------------
// activityEventToResponse
//
// Converts a domain activity event into an API response model.
// -----------------------------------------------------------------------------
func activityEventToResponse(event domain.ActivityEvent) activityEventResponse {
	return activityEventResponse{
		ApplicationID: event.ApplicationID.String(),
		Type:          string(event.Type),
		OccurredAt:    event.OccurredAt.Format(time.RFC3339),
		Description:   event.Description,
	}
}

// -----------------------------------------------------------------------------
// activityEventsToResponse
//
// Converts domain activity events into an API collection response model.
// -----------------------------------------------------------------------------
func activityEventsToResponse(events []domain.ActivityEvent) activityEventsResponse {
	response := activityEventsResponse{
		ActivityEvents: make([]activityEventResponse, 0, len(events)),
	}

	for _, event := range events {
		response.ActivityEvents = append(response.ActivityEvents, activityEventToResponse(event))
	}

	return response
}
