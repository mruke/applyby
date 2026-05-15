package application

import (
	"context"
	"fmt"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// UpdateDocumentRepository
//
// Defines repository behavior required to update one document metadata record.
// -----------------------------------------------------------------------------
type UpdateDocumentRepository interface {
	DocumentFinder
	DocumentUpdater
}

// -----------------------------------------------------------------------------
// UpdateDocumentInput
//
// Contains fields required to update application document metadata.
// -----------------------------------------------------------------------------
type UpdateDocumentInput struct {
	ApplicationID domain.ApplicationID
	DocumentID    domain.DocumentID
	Name          string
	Kind          string
	Path          string
}

// -----------------------------------------------------------------------------
// UpdateDocumentService
//
// Coordinates the workflow for editing application document metadata.
// -----------------------------------------------------------------------------
type UpdateDocumentService struct {
	repository       UpdateDocumentRepository
	activityRecorder ActivityEventRecorder
}

// -----------------------------------------------------------------------------
// NewUpdateDocumentService
//
// Creates a service for the update document metadata workflow.
// -----------------------------------------------------------------------------
func NewUpdateDocumentService(repository UpdateDocumentRepository, activityRecorder ActivityEventRecorder) UpdateDocumentService {
	return UpdateDocumentService{
		repository:       repository,
		activityRecorder: activityRecorder,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Updates one document metadata record and records an activity event.
// -----------------------------------------------------------------------------
func (service UpdateDocumentService) Execute(ctx context.Context, input UpdateDocumentInput) (domain.Document, error) {
	if service.repository == nil {
		return domain.Document{}, fmt.Errorf("document updater repository is required")
	}

	if service.activityRecorder == nil {
		return domain.Document{}, fmt.Errorf("activity recorder is required")
	}

	if err := input.ApplicationID.Validate(); err != nil {
		return domain.Document{}, err
	}

	if err := input.DocumentID.Validate(); err != nil {
		return domain.Document{}, err
	}

	existingDocument, err := service.repository.FindDocumentByID(ctx, input.ApplicationID, input.DocumentID)
	if err != nil {
		return domain.Document{}, err
	}

	updatedDocument, err := domain.NewDocument(
		existingDocument.ID,
		existingDocument.ApplicationID,
		input.Name,
		input.Kind,
		input.Path,
	)
	if err != nil {
		return domain.Document{}, err
	}

	if err := service.repository.UpdateDocument(ctx, updatedDocument); err != nil {
		return domain.Document{}, err
	}

	event, err := domain.NewActivityEvent(
		updatedDocument.ApplicationID,
		domain.ActivityDocumentUpdated,
		time.Now().UTC(),
		fmt.Sprintf("Document metadata updated: %s.", updatedDocument.Name),
	)
	if err != nil {
		return domain.Document{}, err
	}

	if err := service.activityRecorder.RecordActivityEvent(ctx, event); err != nil {
		return domain.Document{}, err
	}

	return updatedDocument, nil
}
