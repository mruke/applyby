package domain

import (
	"fmt"
	"strings"
)

// -----------------------------------------------------------------------------
// ApplicationStatus
//
// Represents a typed job application lifecycle state.
// -----------------------------------------------------------------------------
type ApplicationStatus string

const (
	// -------------------------------------------------------------------------
	// StatusDraft
	//
	// Represents an application record that has not been submitted yet.
	// -------------------------------------------------------------------------
	StatusDraft ApplicationStatus = "draft"

	// -------------------------------------------------------------------------
	// StatusInterested
	//
	// Represents a role the user is interested in but has not applied to yet.
	// -------------------------------------------------------------------------
	StatusInterested ApplicationStatus = "interested"

	// -------------------------------------------------------------------------
	// StatusApplied
	//
	// Represents an application that has been submitted.
	// -------------------------------------------------------------------------
	StatusApplied ApplicationStatus = "applied"

	// -------------------------------------------------------------------------
	// StatusInterviewing
	//
	// Represents an application that has entered an interview process.
	// -------------------------------------------------------------------------
	StatusInterviewing ApplicationStatus = "interviewing"

	// -------------------------------------------------------------------------
	// StatusOffer
	//
	// Represents an application that resulted in an offer.
	// -------------------------------------------------------------------------
	StatusOffer ApplicationStatus = "offer"

	// -------------------------------------------------------------------------
	// StatusRejected
	//
	// Represents an application that was rejected.
	// -------------------------------------------------------------------------
	StatusRejected ApplicationStatus = "rejected"

	// -------------------------------------------------------------------------
	// StatusWithdrawn
	//
	// Represents an application the user chose to withdraw.
	// -------------------------------------------------------------------------
	StatusWithdrawn ApplicationStatus = "withdrawn"

	// -------------------------------------------------------------------------
	// StatusArchived
	//
	// Represents an application that is no longer active.
	// -------------------------------------------------------------------------
	StatusArchived ApplicationStatus = "archived"
)

// -----------------------------------------------------------------------------
// AllApplicationStatuses
//
// Returns every supported application status in a stable order.
// -----------------------------------------------------------------------------
func AllApplicationStatuses() []ApplicationStatus {
	return []ApplicationStatus{
		StatusDraft,
		StatusInterested,
		StatusApplied,
		StatusInterviewing,
		StatusOffer,
		StatusRejected,
		StatusWithdrawn,
		StatusArchived,
	}
}

// -----------------------------------------------------------------------------
// IsValidApplicationStatus
//
// Reports whether a value is one of the supported application statuses.
// -----------------------------------------------------------------------------
func IsValidApplicationStatus(status ApplicationStatus) bool {
	for _, validStatus := range AllApplicationStatuses() {
		if status == validStatus {
			return true
		}
	}

	return false
}

// -----------------------------------------------------------------------------
// ParseApplicationStatus
//
// Converts user or storage text into a typed application status.
// -----------------------------------------------------------------------------
func ParseApplicationStatus(value string) (ApplicationStatus, error) {
	status := ApplicationStatus(strings.ToLower(strings.TrimSpace(value)))

	if !IsValidApplicationStatus(status) {
		return "", fmt.Errorf("invalid application status: %q", value)
	}

	return status, nil
}

// -----------------------------------------------------------------------------
// String
//
// Returns the text representation of an application status.
// -----------------------------------------------------------------------------
func (status ApplicationStatus) String() string {
	return string(status)
}
