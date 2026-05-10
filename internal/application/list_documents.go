package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// ListDocumentsInput
//
// Contains the data required to list document metadata for an application.
// -----------------------------------------------------------------------------
type ListDocumentsInput struct {
	ApplicationID domain.ApplicationID
}

// -----------------------------------------------------------------------------
// ListDocumentsService
//
// Coordinates the workflow for listing application document metadata.
// -----------------------------------------------------------------------------
type ListDocumentsService struct {
	repository DocumentLister
}

// -----------------------------------------------------------------------------
// NewListDocumentsService
//
// Creates a service for the list documents workflow.
// -----------------------------------------------------------------------------
func NewListDocumentsService(repository DocumentLister) ListDocumentsService {
	return ListDocumentsService{
		repository: repository,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Lists document metadata for an application through the repository boundary.
// -----------------------------------------------------------------------------
func (service ListDocumentsService) Execute(ctx context.Context, input ListDocumentsInput) ([]domain.Document, error) {
	if service.repository == nil {
		return nil, fmt.Errorf("document lister is required")
	}

	if err := input.ApplicationID.Validate(); err != nil {
		return nil, err
	}

	return service.repository.ListDocumentsForApplication(ctx, input.ApplicationID)
}
