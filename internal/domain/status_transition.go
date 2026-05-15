package domain

import "fmt"

// -----------------------------------------------------------------------------
// ApplicationStatusTransition
//
// Represents a requested move from one application status to another.
// -----------------------------------------------------------------------------
type ApplicationStatusTransition struct {
	From ApplicationStatus
	To   ApplicationStatus
}

// -----------------------------------------------------------------------------
// AllowedNextStatuses
//
// Returns every valid status except the current one. ApplyBy is a tracker, not a
// strict workflow gate, so users can correct or revise an application's state.
// -----------------------------------------------------------------------------
func AllowedNextStatuses(status ApplicationStatus) []ApplicationStatus {
	if !IsValidApplicationStatus(status) {
		return []ApplicationStatus{}
	}

	allStatuses := AllApplicationStatuses()
	nextStatuses := make([]ApplicationStatus, 0, len(allStatuses)-1)

	for _, nextStatus := range allStatuses {
		if nextStatus != status {
			nextStatuses = append(nextStatuses, nextStatus)
		}
	}

	return nextStatuses
}

// -----------------------------------------------------------------------------
// CanTransitionApplicationStatus
//
// Reports whether an application can move from one valid status to another.
// -----------------------------------------------------------------------------
func CanTransitionApplicationStatus(from ApplicationStatus, to ApplicationStatus) bool {
	return IsValidApplicationStatus(from) && IsValidApplicationStatus(to) && from != to
}

// -----------------------------------------------------------------------------
// ValidateApplicationStatusTransition
//
// Verifies that a requested application status transition is valid.
// -----------------------------------------------------------------------------
func ValidateApplicationStatusTransition(transition ApplicationStatusTransition) error {
	if !IsValidApplicationStatus(transition.From) {
		return fmt.Errorf("invalid source application status: %q", transition.From)
	}

	if !IsValidApplicationStatus(transition.To) {
		return fmt.Errorf("invalid target application status: %q", transition.To)
	}

	if transition.From == transition.To {
		return fmt.Errorf("application status is already %q", transition.To)
	}

	return nil
}

// -----------------------------------------------------------------------------
// IsTerminalApplicationStatus
//
// Reports whether an application status has no forward lifecycle movement.
// -----------------------------------------------------------------------------
func IsTerminalApplicationStatus(status ApplicationStatus) bool {
	return false
}
