package domain

import (
	"fmt"
	"strings"
)

// -----------------------------------------------------------------------------
// Contact
//
// Represents a person connected to a company, application, or job-search process.
// -----------------------------------------------------------------------------
type Contact struct {
	Name  string
	Email string
	Role  string
}

// -----------------------------------------------------------------------------
// NewContact
//
// Creates a contact after applying basic domain validation.
// -----------------------------------------------------------------------------
func NewContact(name string, email string, role string) (Contact, error) {
	contact := Contact{
		Name:  name,
		Email: email,
		Role:  role,
	}

	if err := contact.Validate(); err != nil {
		return Contact{}, err
	}

	return contact, nil
}

// -----------------------------------------------------------------------------
// Validate
//
// Verifies that a contact has a name and a valid optional email shape.
// -----------------------------------------------------------------------------
func (contact Contact) Validate() error {
	if err := requireNonEmptyField("contact name", contact.Name); err != nil {
		return err
	}

	if strings.TrimSpace(contact.Email) != "" && !strings.Contains(contact.Email, "@") {
		return fmt.Errorf("contact email must contain @")
	}

	return nil
}
