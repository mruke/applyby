package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
	"github.com/mruke/applyby/internal/reminders"
)

// -----------------------------------------------------------------------------
// ListRemindersInput
//
// Contains the data required to list reminders for an application.
// -----------------------------------------------------------------------------
type ListRemindersInput struct {
	ApplicationID domain.ApplicationID
}

// -----------------------------------------------------------------------------
// ListRemindersService
//
// Coordinates the workflow for listing prioritized application reminders.
// -----------------------------------------------------------------------------
type ListRemindersService struct {
	repository ReminderLister
}

// -----------------------------------------------------------------------------
// NewListRemindersService
//
// Creates a service for the list reminders workflow.
// -----------------------------------------------------------------------------
func NewListRemindersService(repository ReminderLister) ListRemindersService {
	return ListRemindersService{
		repository: repository,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Lists reminders for an application and returns them in priority order.
// -----------------------------------------------------------------------------
func (service ListRemindersService) Execute(ctx context.Context, input ListRemindersInput) ([]domain.Reminder, error) {
	if service.repository == nil {
		return nil, fmt.Errorf("reminder lister is required")
	}

	if err := input.ApplicationID.Validate(); err != nil {
		return nil, err
	}

	applicationReminders, err := service.repository.ListRemindersForApplication(ctx, input.ApplicationID)
	if err != nil {
		return nil, err
	}

	return reminders.PrioritizeReminders(applicationReminders), nil
}
