package domain

import "fmt"

// -----------------------------------------------------------------------------
// ApplicationStatusTransition
//
// Represents a requested move from one application lifecycle status to another.
// -----------------------------------------------------------------------------
type ApplicationStatusTransition struct {
	From ApplicationStatus
	To   ApplicationStatus
}

var allowedApplicationStatusTransitions = map[ApplicationStatus][]ApplicationStatus{
	StatusDraft: {
		StatusInterested,
		StatusApplied,
		StatusWithdrawn,
		StatusArchived,
	},
	StatusInterested: {
		StatusApplied,
		StatusWithdrawn,
		StatusArchived,
	},
	StatusApplied: {
		StatusInterviewing,
		StatusOffer,
		StatusRejected,
		StatusWithdrawn,
		StatusArchived,
	},
	StatusInterviewing: {
		StatusOffer,
		StatusRejected,
		StatusWithdrawn,
		StatusArchived,
	},
	StatusOffer: {
		StatusArchived,
	},
	StatusRejected: {
		StatusArchived,
	},
	StatusWithdrawn: {
		StatusArchived,
	},
	StatusArchived: {},
}

// -----------------------------------------------------------------------------
// AllowedNextStatuses
//
// Returns the valid next statuses for a current application status.
// -----------------------------------------------------------------------------
func AllowedNextStatuses(status ApplicationStatus) []ApplicationStatus {
	nextStatuses, ok := allowedApplicationStatusTransitions[status]
	if !ok {
		return []ApplicationStatus{}
	}

	copiedStatuses := make([]ApplicationStatus, len(nextStatuses))
	copy(copiedStatuses, nextStatuses)

	return copiedStatuses
}

// -----------------------------------------------------------------------------
// CanTransitionApplicationStatus
//
// Reports whether an application can move from one status to another.
// -----------------------------------------------------------------------------
func CanTransitionApplicationStatus(from ApplicationStatus, to ApplicationStatus) bool {
	for _, nextStatus := range AllowedNextStatuses(from) {
		if nextStatus == to {
			return true
		}
	}

	return false
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

	if !CanTransitionApplicationStatus(transition.From, transition.To) {
		return fmt.Errorf("cannot transition application status from %q to %q", transition.From, transition.To)
	}

	return nil
}

// -----------------------------------------------------------------------------
// IsTerminalApplicationStatus
//
// Reports whether an application status has no forward lifecycle movement.
// -----------------------------------------------------------------------------
func IsTerminalApplicationStatus(status ApplicationStatus) bool {
	return IsValidApplicationStatus(status) && len(AllowedNextStatuses(status)) == 0
}
