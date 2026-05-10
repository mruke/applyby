package application

import (
	"context"
	"errors"
	"testing"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestListDocumentsServiceReturnsDocuments
//
// Verifies that the list documents workflow returns document metadata for an application.
// -----------------------------------------------------------------------------
func TestListDocumentsServiceReturnsDocuments(t *testing.T) {
	repository := newFakeDocumentRepository()
	service := NewListDocumentsService(repository)

	document, err := domain.NewDocument("doc-001", "app-001", "Backend Resume", "resume", "documents/backend-resume.pdf")
	if err != nil {
		t.Fatalf("failed to create test document: %v", err)
	}

	repository.documents[document.ID] = document

	documents, err := service.Execute(context.Background(), ListDocumentsInput{
		ApplicationID: "app-001",
	})

	if err != nil {
		t.Fatalf("expected list documents workflow to succeed: %v", err)
	}

	if len(documents) != 1 {
		t.Fatalf("expected one document, got %d", len(documents))
	}
}

// -----------------------------------------------------------------------------
// TestListDocumentsServiceRejectsMissingApplicationID
//
// Verifies that listing documents requires an application id.
// -----------------------------------------------------------------------------
func TestListDocumentsServiceRejectsMissingApplicationID(t *testing.T) {
	service := NewListDocumentsService(newFakeDocumentRepository())

	_, err := service.Execute(context.Background(), ListDocumentsInput{})

	if err == nil {
		t.Fatal("expected missing application id to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestListDocumentsServiceReturnsRepositoryError
//
// Verifies that repository list errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestListDocumentsServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeDocumentRepository()
	repository.listErr = errors.New("list failed")
	service := NewListDocumentsService(repository)

	_, err := service.Execute(context.Background(), ListDocumentsInput{
		ApplicationID: "app-001",
	})

	if err == nil {
		t.Fatal("expected list error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestListDocumentsServiceRequiresRepository
//
// Verifies that the list documents workflow requires a repository boundary.
// -----------------------------------------------------------------------------
func TestListDocumentsServiceRequiresRepository(t *testing.T) {
	service := NewListDocumentsService(nil)

	_, err := service.Execute(context.Background(), ListDocumentsInput{
		ApplicationID: "app-001",
	})

	if err == nil {
		t.Fatal("expected missing repository to be rejected")
	}
}
