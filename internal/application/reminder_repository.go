package application

import (
	"context"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// ReminderSaver
//
// Defines persistence behavior required to save a reminder.
// -----------------------------------------------------------------------------
type ReminderSaver interface {
	SaveReminder(ctx context.Context, reminder domain.Reminder) error
}

// -----------------------------------------------------------------------------
// ReminderFinder
//
// Defines persistence behavior required to find a reminder by identity.
// -----------------------------------------------------------------------------
type ReminderFinder interface {
	FindReminderByID(ctx context.Context, id domain.ReminderID) (domain.Reminder, error)
}

// -----------------------------------------------------------------------------
// ReminderUpdater
//
// Defines persistence behavior required to update an existing reminder.
// -----------------------------------------------------------------------------
type ReminderUpdater interface {
	UpdateReminder(ctx context.Context, reminder domain.Reminder) error
}

// -----------------------------------------------------------------------------
// ReminderRemover
//
// Defines persistence behavior required to remove one reminder.
// -----------------------------------------------------------------------------
type ReminderRemover interface {
	RemoveReminder(ctx context.Context, id domain.ReminderID) error
}

// -----------------------------------------------------------------------------
// ReminderLister
//
// Defines persistence behavior required to list reminders for an application.
// -----------------------------------------------------------------------------
type ReminderLister interface {
	ListRemindersForApplication(ctx context.Context, applicationID domain.ApplicationID) ([]domain.Reminder, error)
}

// -----------------------------------------------------------------------------
// ReminderRepository
//
// Groups the full reminder repository behavior expected from persistence.
// -----------------------------------------------------------------------------
type ReminderRepository interface {
	ReminderSaver
	ReminderFinder
	ReminderUpdater
	ReminderRemover
	ReminderLister
}
