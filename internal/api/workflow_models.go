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
// contactRequest
//
// Represents the JSON request body for adding a contact.
// -----------------------------------------------------------------------------
type contactRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// -----------------------------------------------------------------------------
// documentRequest
//
// Represents the JSON request body for adding document metadata.
// -----------------------------------------------------------------------------
type documentRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
	Path string `json:"path"`
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
// contactResponse
//
// Represents the JSON response shape for a contact.
// -----------------------------------------------------------------------------
type contactResponse struct {
	ID            string `json:"id"`
	ApplicationID string `json:"application_id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Role          string `json:"role"`
}

// -----------------------------------------------------------------------------
// contactsResponse
//
// Represents the JSON response shape for a contact collection.
// -----------------------------------------------------------------------------
type contactsResponse struct {
	Contacts []contactResponse `json:"contacts"`
}

// -----------------------------------------------------------------------------
// documentResponse
//
// Represents the JSON response shape for document metadata.
// -----------------------------------------------------------------------------
type documentResponse struct {
	ID            string `json:"id"`
	ApplicationID string `json:"application_id"`
	Name          string `json:"name"`
	Kind          string `json:"kind"`
	Path          string `json:"path"`
}

// -----------------------------------------------------------------------------
// documentsResponse
//
// Represents the JSON response shape for a document metadata collection.
// -----------------------------------------------------------------------------
type documentsResponse struct {
	Documents []documentResponse `json:"documents"`
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
// toInput
//
// Converts a contact request into an application-layer input model.
// -----------------------------------------------------------------------------
func (request contactRequest) toInput(applicationID domain.ApplicationID) application.AddContactInput {
	return application.AddContactInput{
		ID:            domain.ContactID(request.ID),
		ApplicationID: applicationID,
		Name:          request.Name,
		Email:         request.Email,
		Role:          request.Role,
	}
}

// -----------------------------------------------------------------------------
// toUpdateInput
//
// Converts a contact request into an update contact workflow input model.
// -----------------------------------------------------------------------------
func (request contactRequest) toUpdateInput(applicationID domain.ApplicationID, contactID domain.ContactID) application.UpdateContactInput {
	return application.UpdateContactInput{
		ApplicationID: applicationID,
		ContactID:     contactID,
		Name:          request.Name,
		Email:         request.Email,
		Role:          request.Role,
	}
}

// -----------------------------------------------------------------------------
// toInput
//
// Converts a document request into an application-layer input model.
// -----------------------------------------------------------------------------
func (request documentRequest) toInput(applicationID domain.ApplicationID) application.AddDocumentInput {
	return application.AddDocumentInput{
		ID:            domain.DocumentID(request.ID),
		ApplicationID: applicationID,
		Name:          request.Name,
		Kind:          request.Kind,
		Path:          request.Path,
	}
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

// -----------------------------------------------------------------------------
// contactToResponse
//
// Converts a domain contact into an API response model.
// -----------------------------------------------------------------------------
func contactToResponse(contact domain.Contact) contactResponse {
	return contactResponse{
		ID:            contact.ID.String(),
		ApplicationID: contact.ApplicationID.String(),
		Name:          contact.Name,
		Email:         contact.Email,
		Role:          contact.Role,
	}
}

// -----------------------------------------------------------------------------
// contactsToResponse
//
// Converts domain contacts into an API collection response model.
// -----------------------------------------------------------------------------
func contactsToResponse(contacts []domain.Contact) contactsResponse {
	response := contactsResponse{
		Contacts: make([]contactResponse, 0, len(contacts)),
	}

	for _, contact := range contacts {
		response.Contacts = append(response.Contacts, contactToResponse(contact))
	}

	return response
}

// -----------------------------------------------------------------------------
// documentToResponse
//
// Converts domain document metadata into an API response model.
// -----------------------------------------------------------------------------
func documentToResponse(document domain.Document) documentResponse {
	return documentResponse{
		ID:            document.ID.String(),
		ApplicationID: document.ApplicationID.String(),
		Name:          document.Name,
		Kind:          document.Kind,
		Path:          document.Path,
	}
}

// -----------------------------------------------------------------------------
// documentsToResponse
//
// Converts domain documents into an API collection response model.
// -----------------------------------------------------------------------------
func documentsToResponse(documents []domain.Document) documentsResponse {
	response := documentsResponse{
		Documents: make([]documentResponse, 0, len(documents)),
	}

	for _, document := range documents {
		response.Documents = append(response.Documents, documentToResponse(document))
	}

	return response
}
