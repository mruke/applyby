package domain

import "testing"

// -----------------------------------------------------------------------------
// TestNewApplicationIDAcceptsValue
//
// Verifies that a non-empty application identity can be created.
// -----------------------------------------------------------------------------
func TestNewApplicationIDAcceptsValue(t *testing.T) {
	id, err := NewApplicationID("app-001")

	if err != nil {
		t.Fatalf("expected application id to be valid: %v", err)
	}

	if id != ApplicationID("app-001") {
		t.Fatalf("expected application id to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestNewApplicationIDTrimsWhitespace
//
// Verifies that application identity input is normalized before storage.
// -----------------------------------------------------------------------------
func TestNewApplicationIDTrimsWhitespace(t *testing.T) {
	id, err := NewApplicationID(" app-001 ")

	if err != nil {
		t.Fatalf("expected application id to be valid: %v", err)
	}

	if id != ApplicationID("app-001") {
		t.Fatalf("expected application id to be trimmed")
	}
}

// -----------------------------------------------------------------------------
// TestNewApplicationIDRejectsEmptyValue
//
// Verifies that an empty application identity is rejected.
// -----------------------------------------------------------------------------
func TestNewApplicationIDRejectsEmptyValue(t *testing.T) {
	_, err := NewApplicationID(" ")

	if err == nil {
		t.Fatal("expected empty application id to be invalid")
	}
}
