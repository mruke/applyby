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
// DocumentFinder
//
// Defines persistence behavior required to find one document metadata record.
// -----------------------------------------------------------------------------
type DocumentFinder interface {
	FindDocumentByID(ctx context.Context, applicationID domain.ApplicationID, documentID domain.DocumentID) (domain.Document, error)
}

// -----------------------------------------------------------------------------
// DocumentUpdater
//
// Defines persistence behavior required to update an existing document metadata record.
// -----------------------------------------------------------------------------
type DocumentUpdater interface {
	UpdateDocument(ctx context.Context, document domain.Document) error
}

// -----------------------------------------------------------------------------
// DocumentRemover
//
// Defines persistence behavior required to remove one document metadata record.
// -----------------------------------------------------------------------------
type DocumentRemover interface {
	RemoveDocument(ctx context.Context, applicationID domain.ApplicationID, documentID domain.DocumentID) error
}

// -----------------------------------------------------------------------------
// DocumentLister
//
// Defines persistence behavior required to list document metadata for an application.
// -----------------------------------------------------------------------------
type DocumentLister interface {
	ListDocumentsForApplication(ctx context.Context, applicationID domain.ApplicationID) ([]domain.Document, error)
}
