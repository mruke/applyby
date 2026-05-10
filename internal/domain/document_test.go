package domain

import "testing"

// -----------------------------------------------------------------------------
// TestNewDocumentAcceptsValidDocument
//
// Verifies that document metadata with identity, application ownership, name, and kind can be created.
// -----------------------------------------------------------------------------
func TestNewDocumentAcceptsValidDocument(t *testing.T) {
	document, err := NewDocument("doc-001", "app-001", "Backend Resume", "resume", "documents/backend-resume.pdf")

	if err != nil {
		t.Fatalf("expected document to be valid: %v", err)
	}

	if document.Kind != "resume" {
		t.Fatalf("expected document kind to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestNewDocumentRejectsMissingID
//
// Verifies that document metadata without an identity is rejected.
// -----------------------------------------------------------------------------
func TestNewDocumentRejectsMissingID(t *testing.T) {
	_, err := NewDocument("", "app-001", "Backend Resume", "resume", "documents/backend-resume.pdf")

	if err == nil {
		t.Fatal("expected document without an id to be invalid")
	}
}

// -----------------------------------------------------------------------------
// TestNewDocumentRejectsMissingApplicationID
//
// Verifies that document metadata without application ownership is rejected.
// -----------------------------------------------------------------------------
func TestNewDocumentRejectsMissingApplicationID(t *testing.T) {
	_, err := NewDocument("doc-001", "", "Backend Resume", "resume", "documents/backend-resume.pdf")

	if err == nil {
		t.Fatal("expected document without an application id to be invalid")
	}
}

// -----------------------------------------------------------------------------
// TestNewDocumentRejectsMissingName
//
// Verifies that document metadata without a name is rejected.
// -----------------------------------------------------------------------------
func TestNewDocumentRejectsMissingName(t *testing.T) {
	_, err := NewDocument("doc-001", "app-001", "", "resume", "documents/backend-resume.pdf")

	if err == nil {
		t.Fatal("expected document without a name to be invalid")
	}
}

// -----------------------------------------------------------------------------
// TestNewDocumentRejectsMissingKind
//
// Verifies that document metadata without a kind is rejected.
// -----------------------------------------------------------------------------
func TestNewDocumentRejectsMissingKind(t *testing.T) {
	_, err := NewDocument("doc-001", "app-001", "Backend Resume", "", "documents/backend-resume.pdf")

	if err == nil {
		t.Fatal("expected document without a kind to be invalid")
	}
}
