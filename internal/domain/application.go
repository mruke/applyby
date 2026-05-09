package domain

import (
	"fmt"
	"time"
)

// -----------------------------------------------------------------------------
// Application
//
// Represents a tracked job application in the ApplyBy domain.
// -----------------------------------------------------------------------------
type Application struct {
	Title          string
	Company        Company
	Status         ApplicationStatus
	Source         string
	CreatedAt      time.Time
	AppliedAt      *time.Time
	Notes          string
	Contacts       []Contact
	Reminders      []Reminder
	Documents      []Document
	ActivityEvents []ActivityEvent
}

// -----------------------------------------------------------------------------
// NewApplication
//
// Creates an application after applying basic domain validation.
// -----------------------------------------------------------------------------
func NewApplication(title string, company Company, status ApplicationStatus, createdAt time.Time) (Application, error) {
	application := Application{
		Title:     title,
		Company:   company,
		Status:    status,
		CreatedAt: createdAt,
	}

	if err := application.Validate(); err != nil {
		return Application{}, err
	}

	return application, nil
}

// -----------------------------------------------------------------------------
// Validate
//
// Verifies that an application has the minimum valid domain shape.
// -----------------------------------------------------------------------------
func (application Application) Validate() error {
	if err := requireNonEmptyField("application title", application.Title); err != nil {
		return err
	}

	if err := application.Company.Validate(); err != nil {
		return err
	}

	if !IsValidApplicationStatus(application.Status) {
		return fmt.Errorf("invalid application status: %q", application.Status)
	}

	if err := requireNonZeroTime("application created date", application.CreatedAt); err != nil {
		return err
	}

	if application.AppliedAt != nil && application.AppliedAt.Before(application.CreatedAt) {
		return fmt.Errorf("application applied date cannot be before created date")
	}

	return nil
}
