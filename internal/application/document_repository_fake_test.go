package application

import (
	"context"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// fakeDocumentRepository
//
// Provides an in-memory document repository for application-layer unit tests.
// -----------------------------------------------------------------------------
type fakeDocumentRepository struct {
	documents map[domain.DocumentID]domain.Document
	saveErr   error
	listErr   error
	saveCalls int
}

// -----------------------------------------------------------------------------
// newFakeDocumentRepository
//
// Creates an empty fake document repository for tests.
// -----------------------------------------------------------------------------
func newFakeDocumentRepository() *fakeDocumentRepository {
	return &fakeDocumentRepository{
		documents: make(map[domain.DocumentID]domain.Document),
	}
}

// -----------------------------------------------------------------------------
// SaveDocument
//
// Stores document metadata in memory for application-layer tests.
// -----------------------------------------------------------------------------
func (repository *fakeDocumentRepository) SaveDocument(ctx context.Context, document domain.Document) error {
	repository.saveCalls++

	if repository.saveErr != nil {
		return repository.saveErr
	}

	repository.documents[document.ID] = document

	return nil
}

// -----------------------------------------------------------------------------
// ListDocumentsForApplication
//
// Lists in-memory document metadata for one application.
// -----------------------------------------------------------------------------
func (repository *fakeDocumentRepository) ListDocumentsForApplication(ctx context.Context, applicationID domain.ApplicationID) ([]domain.Document, error) {
	if repository.listErr != nil {
		return nil, repository.listErr
	}

	documents := []domain.Document{}

	for _, document := range repository.documents {
		if document.ApplicationID == applicationID {
			documents = append(documents, document)
		}
	}

	return documents, nil
}
