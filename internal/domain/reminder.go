package domain

import "time"

// -----------------------------------------------------------------------------
// Reminder
//
// Represents a follow-up or deadline item in the job-search workflow.
// -----------------------------------------------------------------------------
type Reminder struct {
	Title     string
	DueAt     time.Time
	Completed bool
}

// -----------------------------------------------------------------------------
// NewReminder
//
// Creates a reminder after applying basic domain validation.
// -----------------------------------------------------------------------------
func NewReminder(title string, dueAt time.Time) (Reminder, error) {
	reminder := Reminder{
		Title: title,
		DueAt: dueAt,
	}

	if err := reminder.Validate(); err != nil {
		return Reminder{}, err
	}

	return reminder, nil
}

// -----------------------------------------------------------------------------
// Validate
//
// Verifies that a reminder has the minimum information required by the domain.
// -----------------------------------------------------------------------------
func (reminder Reminder) Validate() error {
	if err := requireNonEmptyField("reminder title", reminder.Title); err != nil {
		return err
	}

	return requireNonZeroTime("reminder due date", reminder.DueAt)
}
