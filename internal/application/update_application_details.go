package application

import (
	"context"
	"fmt"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// ApplicationDetailsUpdaterRepository
//
// Defines repository behavior required to update non-status application details.
// -----------------------------------------------------------------------------
type ApplicationDetailsUpdaterRepository interface {
	ApplicationFinder
	ApplicationDetailsUpdater
}

// -----------------------------------------------------------------------------
// UpdateApplicationDetailsInput
//
// Contains fields that can be edited outside the application status lifecycle.
// -----------------------------------------------------------------------------
type UpdateApplicationDetailsInput struct {
	ID             domain.ApplicationID
	Title          string
	CompanyName    string
	CompanyWebsite string
	Source         string
	Notes          string
}

// -----------------------------------------------------------------------------
// UpdateApplicationDetailsService
//
// Coordinates the workflow for editing non-status application details.
// -----------------------------------------------------------------------------
type UpdateApplicationDetailsService struct {
	repository       ApplicationDetailsUpdaterRepository
	activityRecorder ActivityEventRecorder
}

// -----------------------------------------------------------------------------
// NewUpdateApplicationDetailsService
//
// Creates a service for the update application details workflow.
// -----------------------------------------------------------------------------
func NewUpdateApplicationDetailsService(
	repository ApplicationDetailsUpdaterRepository,
	activityRecorder ActivityEventRecorder,
) UpdateApplicationDetailsService {
	return UpdateApplicationDetailsService{
		repository:       repository,
		activityRecorder: activityRecorder,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Loads an application, updates editable details, saves the change, and records activity.
// -----------------------------------------------------------------------------
func (service UpdateApplicationDetailsService) Execute(ctx context.Context, input UpdateApplicationDetailsInput) (domain.Application, error) {
	if service.repository == nil {
		return domain.Application{}, fmt.Errorf("application details updater repository is required")
	}

	if service.activityRecorder == nil {
		return domain.Application{}, fmt.Errorf("activity recorder is required")
	}

	application, err := service.repository.FindApplicationByID(ctx, input.ID)
	if err != nil {
		return domain.Application{}, err
	}

	company, err := domain.NewCompany(input.CompanyName, input.CompanyWebsite)
	if err != nil {
		return domain.Application{}, err
	}

	application.Title = input.Title
	application.Company = company
	application.Source = input.Source
	application.Notes = input.Notes

	if err := application.Validate(); err != nil {
		return domain.Application{}, err
	}

	if err := service.repository.UpdateApplicationDetails(ctx, application); err != nil {
		return domain.Application{}, err
	}

	event, err := domain.NewActivityEvent(
		application.ID,
		domain.ActivityApplicationUpdated,
		time.Now().UTC(),
		fmt.Sprintf("Application details updated: %s at %s.", application.Title, application.Company.Name),
	)
	if err != nil {
		return domain.Application{}, err
	}

	if err := service.activityRecorder.RecordActivityEvent(ctx, event); err != nil {
		return domain.Application{}, err
	}

	return application, nil
}
