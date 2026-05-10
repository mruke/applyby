package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// ReminderCompleterRepository
//
// Defines repository behavior required to complete a reminder.
// -----------------------------------------------------------------------------
type ReminderCompleterRepository interface {
	ReminderFinder
	ReminderSaver
}

// -----------------------------------------------------------------------------
// CompleteReminderInput
//
// Contains the data required to mark a reminder complete.
// -----------------------------------------------------------------------------
type CompleteReminderInput struct {
	ID domain.ReminderID
}

// -----------------------------------------------------------------------------
// CompleteReminderService
//
// Coordinates the workflow for completing a reminder.
// -----------------------------------------------------------------------------
type CompleteReminderService struct {
	repository ReminderCompleterRepository
}

// -----------------------------------------------------------------------------
// NewCompleteReminderService
//
// Creates a service for the complete reminder workflow.
// -----------------------------------------------------------------------------
func NewCompleteReminderService(repository ReminderCompleterRepository) CompleteReminderService {
	return CompleteReminderService{
		repository: repository,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Loads a reminder, marks it complete, and saves the update.
// -----------------------------------------------------------------------------
func (service CompleteReminderService) Execute(ctx context.Context, input CompleteReminderInput) (domain.Reminder, error) {
	if service.repository == nil {
		return domain.Reminder{}, fmt.Errorf("reminder completer repository is required")
	}

	if err := input.ID.Validate(); err != nil {
		return domain.Reminder{}, err
	}

	reminder, err := service.repository.FindReminderByID(ctx, input.ID)
	if err != nil {
		return domain.Reminder{}, err
	}

	reminder.Completed = true

	if err := service.repository.SaveReminder(ctx, reminder); err != nil {
		return domain.Reminder{}, err
	}

	return reminder, nil
}
