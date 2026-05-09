package domain

import "testing"

// -----------------------------------------------------------------------------
// TestNewDocumentAcceptsValidDocument
//
// Verifies that document metadata with a name and kind can be created.
// -----------------------------------------------------------------------------
func TestNewDocumentAcceptsValidDocument(t *testing.T) {
	document, err := NewDocument("Backend Resume", "resume", "documents/backend-resume.pdf")

	if err != nil {
		t.Fatalf("expected document to be valid: %v", err)
	}

	if document.Kind != "resume" {
		t.Fatalf("expected document kind to be preserved")
	}
}

// -----------------------------------------------------------------------------
// TestNewDocumentRejectsMissingName
//
// Verifies that document metadata without a name is rejected.
// -----------------------------------------------------------------------------
func TestNewDocumentRejectsMissingName(t *testing.T) {
	_, err := NewDocument("", "resume", "documents/backend-resume.pdf")

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
	_, err := NewDocument("Backend Resume", "", "documents/backend-resume.pdf")

	if err == nil {
		t.Fatal("expected document without a kind to be invalid")
	}
}
