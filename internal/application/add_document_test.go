package application

import (
	"context"
	"errors"
	"testing"
)

// -----------------------------------------------------------------------------
// TestAddDocumentServiceSavesValidDocument
//
// Verifies that the add document workflow validates and saves document metadata.
// -----------------------------------------------------------------------------
func TestAddDocumentServiceSavesValidDocument(t *testing.T) {
	repository := newFakeDocumentRepository()
	service := NewAddDocumentService(repository)

	document, err := service.Execute(context.Background(), AddDocumentInput{
		ID:            "doc-001",
		ApplicationID: "app-001",
		Name:          "Backend Resume",
		Kind:          "resume",
		Path:          "documents/backend-resume.pdf",
	})

	if err != nil {
		t.Fatalf("expected add document workflow to succeed: %v", err)
	}

	if document.ID != "doc-001" {
		t.Fatalf("expected document id to be preserved")
	}

	if repository.saveCalls != 1 {
		t.Fatalf("expected repository save to be called once")
	}
}

// -----------------------------------------------------------------------------
// TestAddDocumentServiceRejectsInvalidDocument
//
// Verifies that invalid document metadata is rejected before saving.
// -----------------------------------------------------------------------------
func TestAddDocumentServiceRejectsInvalidDocument(t *testing.T) {
	repository := newFakeDocumentRepository()
	service := NewAddDocumentService(repository)

	_, err := service.Execute(context.Background(), AddDocumentInput{
		ID:            "",
		ApplicationID: "app-001",
		Name:          "Backend Resume",
		Kind:          "resume",
		Path:          "documents/backend-resume.pdf",
	})

	if err == nil {
		t.Fatal("expected invalid document to be rejected")
	}

	if repository.saveCalls != 0 {
		t.Fatal("expected invalid document not to be saved")
	}
}

// -----------------------------------------------------------------------------
// TestAddDocumentServiceReturnsRepositoryError
//
// Verifies that repository save errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestAddDocumentServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeDocumentRepository()
	repository.saveErr = errors.New("save failed")
	service := NewAddDocumentService(repository)

	_, err := service.Execute(context.Background(), AddDocumentInput{
		ID:            "doc-001",
		ApplicationID: "app-001",
		Name:          "Backend Resume",
		Kind:          "resume",
		Path:          "documents/backend-resume.pdf",
	})

	if err == nil {
		t.Fatal("expected repository error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestAddDocumentServiceRequiresRepository
//
// Verifies that the add document workflow requires a repository boundary.
// -----------------------------------------------------------------------------
func TestAddDocumentServiceRequiresRepository(t *testing.T) {
	service := NewAddDocumentService(nil)

	_, err := service.Execute(context.Background(), AddDocumentInput{
		ID:            "doc-001",
		ApplicationID: "app-001",
		Name:          "Backend Resume",
		Kind:          "resume",
		Path:          "documents/backend-resume.pdf",
	})

	if err == nil {
		t.Fatal("expected missing repository to be rejected")
	}
}
