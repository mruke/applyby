package domain

import "strings"

// -----------------------------------------------------------------------------
// ApplicationID
//
// Represents the stable domain identity for a tracked job application.
// -----------------------------------------------------------------------------
type ApplicationID string

// -----------------------------------------------------------------------------
// NewApplicationID
//
// Creates an application identity after applying basic domain validation.
// -----------------------------------------------------------------------------
func NewApplicationID(value string) (ApplicationID, error) {
	id := ApplicationID(strings.TrimSpace(value))

	if err := id.Validate(); err != nil {
		return "", err
	}

	return id, nil
}

// -----------------------------------------------------------------------------
// Validate
//
// Verifies that an application identity contains a non-empty value.
// -----------------------------------------------------------------------------
func (id ApplicationID) Validate() error {
	return requireNonEmptyField("application id", string(id))
}

// -----------------------------------------------------------------------------
// String
//
// Returns the text representation of an application identity.
// -----------------------------------------------------------------------------
func (id ApplicationID) String() string {
	return string(id)
}
