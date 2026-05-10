package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// fakeReminderRepository
//
// Provides an in-memory reminder repository for application-layer unit tests.
// -----------------------------------------------------------------------------
type fakeReminderRepository struct {
	reminders map[domain.ReminderID]domain.Reminder
	saveErr   error
	findErr   error
	listErr   error
	saveCalls int
}

// -----------------------------------------------------------------------------
// newFakeReminderRepository
//
// Creates an empty fake reminder repository for tests.
// -----------------------------------------------------------------------------
func newFakeReminderRepository() *fakeReminderRepository {
	return &fakeReminderRepository{
		reminders: make(map[domain.ReminderID]domain.Reminder),
	}
}

// -----------------------------------------------------------------------------
// SaveReminder
//
// Stores a reminder in memory for application-layer tests.
// -----------------------------------------------------------------------------
func (repository *fakeReminderRepository) SaveReminder(ctx context.Context, reminder domain.Reminder) error {
	repository.saveCalls++

	if repository.saveErr != nil {
		return repository.saveErr
	}

	repository.reminders[reminder.ID] = reminder

	return nil
}

// -----------------------------------------------------------------------------
// FindReminderByID
//
// Finds a reminder by identity in memory for application-layer tests.
// -----------------------------------------------------------------------------
func (repository *fakeReminderRepository) FindReminderByID(ctx context.Context, id domain.ReminderID) (domain.Reminder, error) {
	if repository.findErr != nil {
		return domain.Reminder{}, repository.findErr
	}

	reminder, ok := repository.reminders[id]
	if !ok {
		return domain.Reminder{}, fmt.Errorf("reminder not found: %s", id)
	}

	return reminder, nil
}

// -----------------------------------------------------------------------------
// ListRemindersForApplication
//
// Lists in-memory reminders for one application.
// -----------------------------------------------------------------------------
func (repository *fakeReminderRepository) ListRemindersForApplication(ctx context.Context, applicationID domain.ApplicationID) ([]domain.Reminder, error) {
	if repository.listErr != nil {
		return nil, repository.listErr
	}

	reminders := []domain.Reminder{}

	for _, reminder := range repository.reminders {
		if reminder.ApplicationID == applicationID {
			reminders = append(reminders, reminder)
		}
	}

	return reminders, nil
}
