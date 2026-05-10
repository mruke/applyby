package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// AddDocumentInput
//
// Contains the data required to attach document metadata to an application.
// -----------------------------------------------------------------------------
type AddDocumentInput struct {
	ID            domain.DocumentID
	ApplicationID domain.ApplicationID
	Name          string
	Kind          string
	Path          string
}

// -----------------------------------------------------------------------------
// AddDocumentService
//
// Coordinates the workflow for attaching document metadata to an application.
// -----------------------------------------------------------------------------
type AddDocumentService struct {
	repository DocumentSaver
}

// -----------------------------------------------------------------------------
// NewAddDocumentService
//
// Creates a service for the add document workflow.
// -----------------------------------------------------------------------------
func NewAddDocumentService(repository DocumentSaver) AddDocumentService {
	return AddDocumentService{
		repository: repository,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Validates and saves document metadata through the repository boundary.
// -----------------------------------------------------------------------------
func (service AddDocumentService) Execute(ctx context.Context, input AddDocumentInput) (domain.Document, error) {
	if service.repository == nil {
		return domain.Document{}, fmt.Errorf("document saver is required")
	}

	document, err := domain.NewDocument(input.ID, input.ApplicationID, input.Name, input.Kind, input.Path)
	if err != nil {
		return domain.Document{}, err
	}

	if err := service.repository.SaveDocument(ctx, document); err != nil {
		return domain.Document{}, err
	}

	return document, nil
}
