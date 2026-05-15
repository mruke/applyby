package application

import (
	"context"
	"fmt"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// UpdateReminderRepository
//
// Defines repository behavior required to update one reminder.
// -----------------------------------------------------------------------------
type UpdateReminderRepository interface {
	ReminderFinder
	ReminderUpdater
}

// -----------------------------------------------------------------------------
// UpdateReminderInput
//
// Contains editable reminder fields.
// -----------------------------------------------------------------------------
type UpdateReminderInput struct {
	ID    domain.ReminderID
	Title string
	DueAt time.Time
}

// -----------------------------------------------------------------------------
// UpdateReminderService
//
// Coordinates the workflow for editing reminder details.
// -----------------------------------------------------------------------------
type UpdateReminderService struct {
	repository       UpdateReminderRepository
	activityRecorder ActivityEventRecorder
}

// -----------------------------------------------------------------------------
// NewUpdateReminderService
//
// Creates a service for the update reminder workflow.
// -----------------------------------------------------------------------------
func NewUpdateReminderService(repository UpdateReminderRepository, activityRecorder ActivityEventRecorder) UpdateReminderService {
	return UpdateReminderService{
		repository:       repository,
		activityRecorder: activityRecorder,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Updates one reminder and records an activity event.
// -----------------------------------------------------------------------------
func (service UpdateReminderService) Execute(ctx context.Context, input UpdateReminderInput) (domain.Reminder, error) {
	if service.repository == nil {
		return domain.Reminder{}, fmt.Errorf("reminder updater repository is required")
	}

	if service.activityRecorder == nil {
		return domain.Reminder{}, fmt.Errorf("activity recorder is required")
	}

	if err := input.ID.Validate(); err != nil {
		return domain.Reminder{}, err
	}

	existingReminder, err := service.repository.FindReminderByID(ctx, input.ID)
	if err != nil {
		return domain.Reminder{}, err
	}

	updatedReminder, err := domain.NewReminder(existingReminder.ID, existingReminder.ApplicationID, input.Title, input.DueAt)
	if err != nil {
		return domain.Reminder{}, err
	}

	updatedReminder.Completed = existingReminder.Completed

	if err := service.repository.UpdateReminder(ctx, updatedReminder); err != nil {
		return domain.Reminder{}, err
	}
	if err := recordActivityEvent(
		ctx,
		service.activityRecorder,
		updatedReminder.ApplicationID,
		domain.ActivityReminderUpdated,
		fmt.Sprintf("Reminder updated: %s.", updatedReminder.Title),
	); err != nil {
		return domain.Reminder{}, err
	}

	return updatedReminder, nil
}
