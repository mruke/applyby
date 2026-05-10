package application

import (
	"context"
	"fmt"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// ScheduleReminderInput
//
// Contains the data required to schedule a reminder.
// -----------------------------------------------------------------------------
type ScheduleReminderInput struct {
	ID            domain.ReminderID
	ApplicationID domain.ApplicationID
	Title         string
	DueAt         time.Time
}

// -----------------------------------------------------------------------------
// ScheduleReminderService
//
// Coordinates the workflow for scheduling a reminder.
// -----------------------------------------------------------------------------
type ScheduleReminderService struct {
	repository ReminderSaver
}

// -----------------------------------------------------------------------------
// NewScheduleReminderService
//
// Creates a service for the schedule reminder workflow.
// -----------------------------------------------------------------------------
func NewScheduleReminderService(repository ReminderSaver) ScheduleReminderService {
	return ScheduleReminderService{
		repository: repository,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Validates and saves a new reminder through the repository boundary.
// -----------------------------------------------------------------------------
func (service ScheduleReminderService) Execute(ctx context.Context, input ScheduleReminderInput) (domain.Reminder, error) {
	if service.repository == nil {
		return domain.Reminder{}, fmt.Errorf("reminder saver is required")
	}

	reminder, err := domain.NewReminder(input.ID, input.ApplicationID, input.Title, input.DueAt)
	if err != nil {
		return domain.Reminder{}, err
	}

	if err := service.repository.SaveReminder(ctx, reminder); err != nil {
		return domain.Reminder{}, err
	}

	return reminder, nil
}
