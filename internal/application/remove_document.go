package application

import (
	"context"
	"fmt"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// RemoveDocumentRepository
//
// Defines repository behavior required to remove one document metadata record.
// -----------------------------------------------------------------------------
type RemoveDocumentRepository interface {
	DocumentFinder
	DocumentRemover
}

// -----------------------------------------------------------------------------
// RemoveDocumentInput
//
// Contains identifiers required to remove application document metadata.
// -----------------------------------------------------------------------------
type RemoveDocumentInput struct {
	ApplicationID domain.ApplicationID
	DocumentID    domain.DocumentID
}

// -----------------------------------------------------------------------------
// RemoveDocumentService
//
// Coordinates the workflow for removing application document metadata.
// -----------------------------------------------------------------------------
type RemoveDocumentService struct {
	repository       RemoveDocumentRepository
	activityRecorder ActivityEventRecorder
}

// -----------------------------------------------------------------------------
// NewRemoveDocumentService
//
// Creates a service for the remove document metadata workflow.
// -----------------------------------------------------------------------------
func NewRemoveDocumentService(repository RemoveDocumentRepository, activityRecorder ActivityEventRecorder) RemoveDocumentService {
	return RemoveDocumentService{
		repository:       repository,
		activityRecorder: activityRecorder,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Removes one document metadata record and records an activity event.
// -----------------------------------------------------------------------------
func (service RemoveDocumentService) Execute(ctx context.Context, input RemoveDocumentInput) error {
	if service.repository == nil {
		return fmt.Errorf("document remover repository is required")
	}

	if service.activityRecorder == nil {
		return fmt.Errorf("activity recorder is required")
	}

	if err := input.ApplicationID.Validate(); err != nil {
		return err
	}

	if err := input.DocumentID.Validate(); err != nil {
		return err
	}

	document, err := service.repository.FindDocumentByID(ctx, input.ApplicationID, input.DocumentID)
	if err != nil {
		return err
	}

	if err := service.repository.RemoveDocument(ctx, input.ApplicationID, input.DocumentID); err != nil {
		return err
	}

	event, err := domain.NewActivityEvent(
		document.ApplicationID,
		domain.ActivityDocumentRemoved,
		time.Now().UTC(),
		fmt.Sprintf("Document metadata removed: %s.", document.Name),
	)
	if err != nil {
		return err
	}

	return service.activityRecorder.RecordActivityEvent(ctx, event)
}
