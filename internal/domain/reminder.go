package domain

import "time"

// -----------------------------------------------------------------------------
// ReminderID
//
// Represents the stable domain identity for a follow-up reminder.
// -----------------------------------------------------------------------------
type ReminderID string

// -----------------------------------------------------------------------------
// Validate
//
// Verifies that a reminder identity contains a non-empty value.
// -----------------------------------------------------------------------------
func (id ReminderID) Validate() error {
	return requireNonEmptyField("reminder id", string(id))
}

// -----------------------------------------------------------------------------
// String
//
// Returns the text representation of a reminder identity.
// -----------------------------------------------------------------------------
func (id ReminderID) String() string {
	return string(id)
}

// -----------------------------------------------------------------------------
// Reminder
//
// Represents a follow-up or deadline item in the job-search workflow.
// -----------------------------------------------------------------------------
type Reminder struct {
	ID            ReminderID
	ApplicationID ApplicationID
	Title         string
	DueAt         time.Time
	Completed     bool
}

// -----------------------------------------------------------------------------
// NewReminder
//
// Creates a reminder after applying basic domain validation.
// -----------------------------------------------------------------------------
func NewReminder(id ReminderID, applicationID ApplicationID, title string, dueAt time.Time) (Reminder, error) {
	reminder := Reminder{
		ID:            id,
		ApplicationID: applicationID,
		Title:         title,
		DueAt:         dueAt,
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
	if err := reminder.ID.Validate(); err != nil {
		return err
	}

	if err := reminder.ApplicationID.Validate(); err != nil {
		return err
	}

	if err := requireNonEmptyField("reminder title", reminder.Title); err != nil {
		return err
	}

	return requireNonZeroTime("reminder due date", reminder.DueAt)
}
