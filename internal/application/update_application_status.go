package application

import (
	"context"
	"fmt"
	"time"

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
	repository      ApplicationStatusUpdaterRepository
	historyRecorder ApplicationHistoryRecorder
}

// -----------------------------------------------------------------------------
// NewUpdateApplicationStatusService
//
// Creates a service for the update application status workflow.
// -----------------------------------------------------------------------------
func NewUpdateApplicationStatusService(
	repository ApplicationStatusUpdaterRepository,
	historyRecorder ApplicationHistoryRecorder,
) UpdateApplicationStatusService {
	return UpdateApplicationStatusService{
		repository:      repository,
		historyRecorder: historyRecorder,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Loads an application, validates the lifecycle transition, saves the change, and records history.
// -----------------------------------------------------------------------------
func (service UpdateApplicationStatusService) Execute(ctx context.Context, input UpdateApplicationStatusInput) (domain.Application, error) {
	if service.repository == nil {
		return domain.Application{}, fmt.Errorf("application status updater repository is required")
	}

	if service.historyRecorder == nil {
		return domain.Application{}, fmt.Errorf("application history recorder is required")
	}

	application, err := service.repository.FindApplicationByID(ctx, input.ID)
	if err != nil {
		return domain.Application{}, err
	}

	fromStatus := application.Status

	transition := domain.ApplicationStatusTransition{
		From: fromStatus,
		To:   input.Status,
	}

	if err := domain.ValidateApplicationStatusTransition(transition); err != nil {
		return domain.Application{}, err
	}

	changedAt := time.Now().UTC()

	statusHistory, err := domain.NewApplicationStatusHistory(application.ID, fromStatus, input.Status, changedAt)
	if err != nil {
		return domain.Application{}, err
	}

	activityEvent, err := domain.NewActivityEvent(
		application.ID,
		domain.ActivityStatusChanged,
		changedAt,
		fmt.Sprintf("Status changed from %s to %s.", fromStatus, input.Status),
	)
	if err != nil {
		return domain.Application{}, err
	}

	application.Status = input.Status

	if err := service.repository.SaveApplication(ctx, application); err != nil {
		return domain.Application{}, err
	}

	if err := service.historyRecorder.RecordApplicationStatusHistory(ctx, statusHistory); err != nil {
		return domain.Application{}, err
	}

	if err := service.historyRecorder.RecordActivityEvent(ctx, activityEvent); err != nil {
		return domain.Application{}, err
	}

	return application, nil
}
