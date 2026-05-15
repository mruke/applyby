package api

import (
	"net/http"
	"time"

	"github.com/mruke/applyby/internal/domain"
	"github.com/mruke/applyby/internal/search"
)

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
