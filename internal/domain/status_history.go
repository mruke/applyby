package domain

import (
	"fmt"
	"time"
)

// -----------------------------------------------------------------------------
// ApplicationStatusHistory
//
// Represents a structured append-only record of an application status change.
// -----------------------------------------------------------------------------
type ApplicationStatusHistory struct {
	ApplicationID ApplicationID
	FromStatus    ApplicationStatus
	ToStatus      ApplicationStatus
	ChangedAt     time.Time
}

// -----------------------------------------------------------------------------
// NewApplicationStatusHistory
//
// Creates status history after applying lifecycle and domain validation.
// -----------------------------------------------------------------------------
func NewApplicationStatusHistory(
	applicationID ApplicationID,
	fromStatus ApplicationStatus,
	toStatus ApplicationStatus,
	changedAt time.Time,
) (ApplicationStatusHistory, error) {
	history := ApplicationStatusHistory{
		ApplicationID: applicationID,
		FromStatus:    fromStatus,
		ToStatus:      toStatus,
		ChangedAt:     changedAt,
	}

	if err := history.Validate(); err != nil {
		return ApplicationStatusHistory{}, err
	}

	return history, nil
}

// -----------------------------------------------------------------------------
// Validate
//
// Verifies that status history records a valid lifecycle transition.
// -----------------------------------------------------------------------------
func (history ApplicationStatusHistory) Validate() error {
	if err := history.ApplicationID.Validate(); err != nil {
		return err
	}

	transition := ApplicationStatusTransition{
		From: history.FromStatus,
		To:   history.ToStatus,
	}

	if err := ValidateApplicationStatusTransition(transition); err != nil {
		return err
	}

	if err := requireNonZeroTime("status history changed date", history.ChangedAt); err != nil {
		return err
	}

	if history.FromStatus == history.ToStatus {
		return fmt.Errorf("status history must change status")
	}

	return nil
}
