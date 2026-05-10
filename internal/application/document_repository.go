package application

import (
	"context"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// DocumentSaver
//
// Defines persistence behavior required to save document metadata.
// -----------------------------------------------------------------------------
type DocumentSaver interface {
	SaveDocument(ctx context.Context, document domain.Document) error
}

// -----------------------------------------------------------------------------
// DocumentLister
//
// Defines persistence behavior required to list document metadata for an application.
// -----------------------------------------------------------------------------
type DocumentLister interface {
	ListDocumentsForApplication(ctx context.Context, applicationID domain.ApplicationID) ([]domain.Document, error)
}

// -----------------------------------------------------------------------------
// DocumentRepository
//
// Groups the full document repository behavior expected from persistence.
// -----------------------------------------------------------------------------
type DocumentRepository interface {
	DocumentSaver
	DocumentLister
}
