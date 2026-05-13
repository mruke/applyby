package application

import (
	"context"
	"fmt"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

type ReminderCompleterRepository interface {
	ReminderFinder
	ReminderSaver
}

type CompleteReminderInput struct {
	ID domain.ReminderID
}

type CompleteReminderService struct {
	repository       ReminderCompleterRepository
	activityRecorder ActivityEventRecorder
}

func NewCompleteReminderService(repository ReminderCompleterRepository, activityRecorder ActivityEventRecorder) CompleteReminderService {
	return CompleteReminderService{
		repository:       repository,
		activityRecorder: activityRecorder,
	}
}

func (service CompleteReminderService) Execute(ctx context.Context, input CompleteReminderInput) (domain.Reminder, error) {
	if service.repository == nil {
		return domain.Reminder{}, fmt.Errorf("reminder completer repository is required")
	}

	if service.activityRecorder == nil {
		return domain.Reminder{}, fmt.Errorf("activity recorder is required")
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

	event, err := domain.NewActivityEvent(
		reminder.ApplicationID,
		domain.ActivityReminderCompleted,
		time.Now().UTC(),
		fmt.Sprintf("Reminder completed: %s.", reminder.Title),
	)
	if err != nil {
		return domain.Reminder{}, err
	}

	if err := service.activityRecorder.RecordActivityEvent(ctx, event); err != nil {
		return domain.Reminder{}, err
	}

	return reminder, nil
}
