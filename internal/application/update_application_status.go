package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// ApplicationStatusUpdaterRepository
//
// Defines the repository behavior required to update application status.
// -----------------------------------------------------------------------------
type ApplicationStatusUpdaterRepository interface {
	ApplicationFinder
	ApplicationSaver
}

// -----------------------------------------------------------------------------
// UpdateApplicationStatusInput
//
// Contains the data required to update an application lifecycle status.
// -----------------------------------------------------------------------------
type UpdateApplicationStatusInput struct {
	ID     domain.ApplicationID
	Status domain.ApplicationStatus
}

// -----------------------------------------------------------------------------
// UpdateApplicationStatusService
//
// Coordinates the workflow for changing an application status.
// -----------------------------------------------------------------------------
type UpdateApplicationStatusService struct {
	repository ApplicationStatusUpdaterRepository
}

// -----------------------------------------------------------------------------
// NewUpdateApplicationStatusService
//
// Creates a service for the update application status workflow.
// -----------------------------------------------------------------------------
func NewUpdateApplicationStatusService(repository ApplicationStatusUpdaterRepository) UpdateApplicationStatusService {
	return UpdateApplicationStatusService{
		repository: repository,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Loads an application, validates the lifecycle transition, and saves the change.
// -----------------------------------------------------------------------------
func (service UpdateApplicationStatusService) Execute(ctx context.Context, input UpdateApplicationStatusInput) (domain.Application, error) {
	if service.repository == nil {
		return domain.Application{}, fmt.Errorf("application status updater repository is required")
	}

	application, err := service.repository.FindApplicationByID(ctx, input.ID)
	if err != nil {
		return domain.Application{}, err
	}

	transition := domain.ApplicationStatusTransition{
		From: application.Status,
		To:   input.Status,
	}

	if err := domain.ValidateApplicationStatusTransition(transition); err != nil {
		return domain.Application{}, err
	}

	application.Status = input.Status

	if err := service.repository.SaveApplication(ctx, application); err != nil {
		return domain.Application{}, err
	}

	return application, nil
}
