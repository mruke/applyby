package domain

import (
	"fmt"
	"strings"
)

// -----------------------------------------------------------------------------
// ContactID
//
// Represents the stable domain identity for a contact.
// -----------------------------------------------------------------------------
type ContactID string

// -----------------------------------------------------------------------------
// Validate
//
// Verifies that a contact identity contains a non-empty value.
// -----------------------------------------------------------------------------
func (id ContactID) Validate() error {
	return requireNonEmptyField("contact id", string(id))
}

// -----------------------------------------------------------------------------
// String
//
// Returns the text representation of a contact identity.
// -----------------------------------------------------------------------------
func (id ContactID) String() string {
	return string(id)
}

// -----------------------------------------------------------------------------
// Contact
//
// Represents a person connected to an application or job-search process.
// -----------------------------------------------------------------------------
type Contact struct {
	ID            ContactID
	ApplicationID ApplicationID
	Name          string
	Email         string
	Role          string
}

// -----------------------------------------------------------------------------
// NewContact
//
// Creates a contact after applying basic domain validation.
// -----------------------------------------------------------------------------
func NewContact(id ContactID, applicationID ApplicationID, name string, email string, role string) (Contact, error) {
	contact := Contact{
		ID:            id,
		ApplicationID: applicationID,
		Name:          name,
		Email:         email,
		Role:          role,
	}

	if err := contact.Validate(); err != nil {
		return Contact{}, err
	}

	return contact, nil
}

// -----------------------------------------------------------------------------
// Validate
//
// Verifies that a contact has identity, application ownership, a name, and a valid optional email shape.
// -----------------------------------------------------------------------------
func (contact Contact) Validate() error {
	if err := contact.ID.Validate(); err != nil {
		return err
	}

	if err := contact.ApplicationID.Validate(); err != nil {
		return err
	}

	if err := requireNonEmptyField("contact name", contact.Name); err != nil {
		return err
	}

	if strings.TrimSpace(contact.Email) != "" && !strings.Contains(contact.Email, "@") {
		return fmt.Errorf("contact email must contain @")
	}

	return nil
}
