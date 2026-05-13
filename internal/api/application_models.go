package api

import (
	"fmt"
	"time"

	"github.com/mruke/applyby/internal/application"
	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// createApplicationRequest
//
// Represents the JSON request body for creating an application.
// -----------------------------------------------------------------------------
type createApplicationRequest struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	CompanyName    string `json:"company_name"`
	CompanyWebsite string `json:"company_website"`
	Status         string `json:"status"`
	Source         string `json:"source"`
	Notes          string `json:"notes"`
	CreatedAt      string `json:"created_at"`
}

// -----------------------------------------------------------------------------
// updateApplicationDetailsRequest
//
// Represents the JSON request body for updating non-status application details.
// -----------------------------------------------------------------------------
type updateApplicationDetailsRequest struct {
	Title          string `json:"title"`
	CompanyName    string `json:"company_name"`
	CompanyWebsite string `json:"company_website"`
	Source         string `json:"source"`
	Notes          string `json:"notes"`
}

// -----------------------------------------------------------------------------
// updateApplicationStatusRequest
//
// Represents the JSON request body for updating an application status.
// -----------------------------------------------------------------------------
type updateApplicationStatusRequest struct {
	Status string `json:"status"`
}

// -----------------------------------------------------------------------------
// applicationResponse
//
// Represents the JSON response shape for an application.
// -----------------------------------------------------------------------------
type applicationResponse struct {
	ID             string  `json:"id"`
	Title          string  `json:"title"`
	CompanyName    string  `json:"company_name"`
	CompanyWebsite string  `json:"company_website"`
	Status         string  `json:"status"`
	Source         string  `json:"source"`
	Notes          string  `json:"notes"`
	CreatedAt      string  `json:"created_at"`
	AppliedAt      *string `json:"applied_at,omitempty"`
}

// -----------------------------------------------------------------------------
// applicationsResponse
//
// Represents the JSON response shape for a collection of applications.
// -----------------------------------------------------------------------------
type applicationsResponse struct {
	Applications []applicationResponse `json:"applications"`
}

// -----------------------------------------------------------------------------
// errorResponse
//
// Represents a simple JSON error response.
// -----------------------------------------------------------------------------
type errorResponse struct {
	Error string `json:"error"`
}

// -----------------------------------------------------------------------------
// toInput
//
// Converts a create application request into an application-layer input model.
// -----------------------------------------------------------------------------
func (request createApplicationRequest) toInput() (application.CreateApplicationInput, error) {
	id, err := domain.NewApplicationID(request.ID)
	if err != nil {
		return application.CreateApplicationInput{}, err
	}

	company, err := domain.NewCompany(request.CompanyName, request.CompanyWebsite)
	if err != nil {
		return application.CreateApplicationInput{}, err
	}

	status, err := domain.ParseApplicationStatus(request.Status)
	if err != nil {
		return application.CreateApplicationInput{}, err
	}

	createdAt, err := parseRequiredRFC3339Time("created_at", request.CreatedAt)
	if err != nil {
		return application.CreateApplicationInput{}, err
	}

	return application.CreateApplicationInput{
		ID:        id,
		Title:     request.Title,
		Company:   company,
		Status:    status,
		Source:    request.Source,
		Notes:     request.Notes,
		CreatedAt: createdAt,
	}, nil
}

// -----------------------------------------------------------------------------
// toInput
//
// Converts an update application details request into an application-layer input model.
// -----------------------------------------------------------------------------
func (request updateApplicationDetailsRequest) toInput(id domain.ApplicationID) application.UpdateApplicationDetailsInput {
	return application.UpdateApplicationDetailsInput{
		ID:             id,
		Title:          request.Title,
		CompanyName:    request.CompanyName,
		CompanyWebsite: request.CompanyWebsite,
		Source:         request.Source,
		Notes:          request.Notes,
	}
}

// -----------------------------------------------------------------------------
// toStatus
//
// Converts an update status request into a domain application status.
// -----------------------------------------------------------------------------
func (request updateApplicationStatusRequest) toStatus() (domain.ApplicationStatus, error) {
	return domain.ParseApplicationStatus(request.Status)
}

// -----------------------------------------------------------------------------
// applicationToResponse
//
// Converts a domain application into an API response model.
// -----------------------------------------------------------------------------
func applicationToResponse(application domain.Application) applicationResponse {
	response := applicationResponse{
		ID:             application.ID.String(),
		Title:          application.Title,
		CompanyName:    application.Company.Name,
		CompanyWebsite: application.Company.Website,
		Status:         application.Status.String(),
		Source:         application.Source,
		Notes:          application.Notes,
		CreatedAt:      application.CreatedAt.Format(time.RFC3339),
	}

	if application.AppliedAt != nil {
		appliedAt := application.AppliedAt.Format(time.RFC3339)
		response.AppliedAt = &appliedAt
	}

	return response
}

// -----------------------------------------------------------------------------
// applicationsToResponse
//
// Converts domain applications into an API collection response model.
// -----------------------------------------------------------------------------
func applicationsToResponse(applications []domain.Application) applicationsResponse {
	response := applicationsResponse{
		Applications: make([]applicationResponse, 0, len(applications)),
	}

	for _, application := range applications {
		response.Applications = append(response.Applications, applicationToResponse(application))
	}

	return response
}

// -----------------------------------------------------------------------------
// parseRequiredRFC3339Time
//
// Parses a required RFC3339 timestamp field from an API request.
// -----------------------------------------------------------------------------
func parseRequiredRFC3339Time(fieldName string, value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("%s is required", fieldName)
	}

	parsedTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s must be RFC3339 format", fieldName)
	}

	return parsedTime, nil
}
