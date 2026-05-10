package postgres

import (
	"context"
	"testing"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestApplicationRepositorySavesAndListsDocuments
//
// Verifies that PostgreSQL persistence can save and list document metadata for an application.
// -----------------------------------------------------------------------------
func TestApplicationRepositorySavesAndListsDocuments(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)
	application := newContactDocumentRepositoryTestApplication(t)

	if err := repository.SaveApplication(context.Background(), application); err != nil {
		t.Fatalf("expected application save to succeed: %v", err)
	}

	document, err := domain.NewDocument("doc-001", application.ID, "Backend Resume", "resume", "documents/backend-resume.pdf")
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}

	if err := repository.SaveDocument(context.Background(), document); err != nil {
		t.Fatalf("expected document save to succeed: %v", err)
	}

	documents, err := repository.ListDocumentsForApplication(context.Background(), application.ID)
	if err != nil {
		t.Fatalf("expected document list to succeed: %v", err)
	}

	if len(documents) != 1 {
		t.Fatalf("expected one document, got %d", len(documents))
	}

	if documents[0].ID != document.ID {
		t.Fatalf("expected document id %q, got %q", document.ID, documents[0].ID)
	}
}

// -----------------------------------------------------------------------------
// TestApplicationRepositoryUpdatesExistingDocument
//
// Verifies that PostgreSQL persistence updates existing document metadata.
// -----------------------------------------------------------------------------
func TestApplicationRepositoryUpdatesExistingDocument(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)
	application := newContactDocumentRepositoryTestApplication(t)

	if err := repository.SaveApplication(context.Background(), application); err != nil {
		t.Fatalf("expected application save to succeed: %v", err)
	}

	document, err := domain.NewDocument("doc-001", application.ID, "Backend Resume", "resume", "documents/backend-resume.pdf")
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}

	if err := repository.SaveDocument(context.Background(), document); err != nil {
		t.Fatalf("expected document save to succeed: %v", err)
	}

	document.Kind = "cover_letter"

	if err := repository.SaveDocument(context.Background(), document); err != nil {
		t.Fatalf("expected document update to succeed: %v", err)
	}

	documents, err := repository.ListDocumentsForApplication(context.Background(), application.ID)
	if err != nil {
		t.Fatalf("expected document list to succeed: %v", err)
	}

	if documents[0].Kind != "cover_letter" {
		t.Fatalf("expected document kind to be updated")
	}
}
