package application

import (
	"context"
	"fmt"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

type ScheduleReminderInput struct {
	ID            domain.ReminderID
	ApplicationID domain.ApplicationID
	Title         string
	DueAt         time.Time
}

type ScheduleReminderService struct {
	repository       ReminderSaver
	activityRecorder ActivityEventRecorder
}

func NewScheduleReminderService(repository ReminderSaver, activityRecorder ActivityEventRecorder) ScheduleReminderService {
	return ScheduleReminderService{
		repository:       repository,
		activityRecorder: activityRecorder,
	}
}

func (service ScheduleReminderService) Execute(ctx context.Context, input ScheduleReminderInput) (domain.Reminder, error) {
	if service.repository == nil {
		return domain.Reminder{}, fmt.Errorf("reminder saver is required")
	}

	if service.activityRecorder == nil {
		return domain.Reminder{}, fmt.Errorf("activity recorder is required")
	}

	reminder, err := domain.NewReminder(input.ID, input.ApplicationID, input.Title, input.DueAt)
	if err != nil {
		return domain.Reminder{}, err
	}

	if err := service.repository.SaveReminder(ctx, reminder); err != nil {
		return domain.Reminder{}, err
	}

	event, err := domain.NewActivityEvent(
		reminder.ApplicationID,
		domain.ActivityReminderScheduled,
		time.Now().UTC(),
		fmt.Sprintf("Reminder scheduled: %s.", reminder.Title),
	)
	if err != nil {
		return domain.Reminder{}, err
	}

	if err := service.activityRecorder.RecordActivityEvent(ctx, event); err != nil {
		return domain.Reminder{}, err
	}

	return reminder, nil
}
