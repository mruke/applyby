package domain

import "testing"

// -----------------------------------------------------------------------------
// TestNewContactAcceptsValidContact
//
// Verifies that a contact with a name and email can be created.
// -----------------------------------------------------------------------------
func TestNewContactAcceptsValidContact(t *testing.T) {
	contact, err := NewContact("Sam Recruiter", "sam@example.com", "Recruiter")

	if err != nil {
		t.Fatalf("expected contact to be valid: %v", err)
	}

	if contact.Name != "Sam Recruiter" {
		t.Fatalf("expected contact name to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestNewContactRejectsMissingName
//
// Verifies that a contact without a name is rejected.
// -----------------------------------------------------------------------------
func TestNewContactRejectsMissingName(t *testing.T) {
	_, err := NewContact("", "sam@example.com", "Recruiter")

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
	_, err := NewContact("Sam Recruiter", "not-an-email", "Recruiter")

	if err == nil {
		t.Fatal("expected contact with invalid email shape to be invalid")
	}
}
