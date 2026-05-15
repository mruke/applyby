package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// RemoveReminderRepository
//
// Defines repository behavior required to remove one reminder.
// -----------------------------------------------------------------------------
type RemoveReminderRepository interface {
	ReminderFinder
	ReminderRemover
}

// -----------------------------------------------------------------------------
// RemoveReminderInput
//
// Contains the reminder identity required for removal.
// -----------------------------------------------------------------------------
type RemoveReminderInput struct {
	ID domain.ReminderID
}

// -----------------------------------------------------------------------------
// RemoveReminderService
//
// Coordinates the workflow for removing a reminder.
// -----------------------------------------------------------------------------
type RemoveReminderService struct {
	repository       RemoveReminderRepository
	activityRecorder ActivityEventRecorder
}

// -----------------------------------------------------------------------------
// NewRemoveReminderService
//
// Creates a service for the remove reminder workflow.
// -----------------------------------------------------------------------------
func NewRemoveReminderService(repository RemoveReminderRepository, activityRecorder ActivityEventRecorder) RemoveReminderService {
	return RemoveReminderService{
		repository:       repository,
		activityRecorder: activityRecorder,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Removes one reminder and records an activity event.
// -----------------------------------------------------------------------------
func (service RemoveReminderService) Execute(ctx context.Context, input RemoveReminderInput) error {
	if service.repository == nil {
		return fmt.Errorf("reminder remover repository is required")
	}

	if service.activityRecorder == nil {
		return fmt.Errorf("activity recorder is required")
	}

	if err := input.ID.Validate(); err != nil {
		return err
	}

	reminder, err := service.repository.FindReminderByID(ctx, input.ID)
	if err != nil {
		return err
	}

	if err := service.repository.RemoveReminder(ctx, input.ID); err != nil {
		return err
	}
	return recordActivityEvent(
		ctx,
		service.activityRecorder,
		reminder.ApplicationID,
		domain.ActivityReminderRemoved,
		fmt.Sprintf("Reminder removed: %s.", reminder.Title),
	)
}
