package domain

import "testing"

// -----------------------------------------------------------------------------
// TestNewContactAcceptsValidContact
//
// Verifies that a contact with identity, application ownership, name, and email can be created.
// -----------------------------------------------------------------------------
func TestNewContactAcceptsValidContact(t *testing.T) {
	contact, err := NewContact("contact-001", "app-001", "Sam Recruiter", "sam@example.com", "Recruiter")

	if err != nil {
		t.Fatalf("expected contact to be valid: %v", err)
	}

	if contact.Name != "Sam Recruiter" {
		t.Fatalf("expected contact name to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestNewContactRejectsMissingID
//
// Verifies that a contact without an identity is rejected.
// -----------------------------------------------------------------------------
func TestNewContactRejectsMissingID(t *testing.T) {
	_, err := NewContact("", "app-001", "Sam Recruiter", "sam@example.com", "Recruiter")

	if err == nil {
		t.Fatal("expected contact without an id to be invalid")
	}
}

// -----------------------------------------------------------------------------
// TestNewContactRejectsMissingApplicationID
//
// Verifies that a contact without application ownership is rejected.
// -----------------------------------------------------------------------------
func TestNewContactRejectsMissingApplicationID(t *testing.T) {
	_, err := NewContact("contact-001", "", "Sam Recruiter", "sam@example.com", "Recruiter")

	if err == nil {
		t.Fatal("expected contact without an application id to be invalid")
	}
}

// -----------------------------------------------------------------------------
// TestNewContactRejectsMissingName
//
// Verifies that a contact without a name is rejected.
// -----------------------------------------------------------------------------
func TestNewContactRejectsMissingName(t *testing.T) {
	_, err := NewContact("contact-001", "app-001", "", "sam@example.com", "Recruiter")

	if err == nil {
		t.Fatal("expected contact without a name to be invalid")
	}
}

// -----------------------------------------------------------------------------
// TestNewContactRejectsInvalidEmailShape
//
// Verifies that a non-empty email must contain an at sign.
// -----------------------------------------------------------------------------
func TestNewContactRejectsInvalidEmailShape(t *testing.T) {
	_, err := NewContact("contact-001", "app-001", "Sam Recruiter", "not-an-email", "Recruiter")

	if err == nil {
		t.Fatal("expected contact with invalid email shape to be invalid")
	}
}
