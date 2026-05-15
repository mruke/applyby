package application

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// fakeDocumentMaintenanceRepository
//
// Provides document maintenance behavior for application-layer tests.
// -----------------------------------------------------------------------------
type fakeDocumentMaintenanceRepository struct {
	documents map[domain.DocumentID]domain.Document
	findErr   error
	updateErr error
	removeErr error
}

// -----------------------------------------------------------------------------
// newFakeDocumentMaintenanceRepository
//
// Creates an empty document maintenance fake repository.
// -----------------------------------------------------------------------------
func newFakeDocumentMaintenanceRepository() *fakeDocumentMaintenanceRepository {
	return &fakeDocumentMaintenanceRepository{
		documents: make(map[domain.DocumentID]domain.Document),
	}
}

// -----------------------------------------------------------------------------
// FindDocumentByID
//
// Finds document metadata by application and document identity.
// -----------------------------------------------------------------------------
func (repository *fakeDocumentMaintenanceRepository) FindDocumentByID(ctx context.Context, applicationID domain.ApplicationID, documentID domain.DocumentID) (domain.Document, error) {
	if repository.findErr != nil {
		return domain.Document{}, repository.findErr
	}

	document, ok := repository.documents[documentID]
	if !ok || document.ApplicationID != applicationID {
		return domain.Document{}, fmt.Errorf("document not found: %s", documentID)
	}

	return document, nil
}

// -----------------------------------------------------------------------------
// UpdateDocument
//
// Updates document metadata in memory.
// -----------------------------------------------------------------------------
func (repository *fakeDocumentMaintenanceRepository) UpdateDocument(ctx context.Context, document domain.Document) error {
	if repository.updateErr != nil {
		return repository.updateErr
	}

	repository.documents[document.ID] = document

	return nil
}

// -----------------------------------------------------------------------------
// RemoveDocument
//
// Removes document metadata from memory.
// -----------------------------------------------------------------------------
func (repository *fakeDocumentMaintenanceRepository) RemoveDocument(ctx context.Context, applicationID domain.ApplicationID, documentID domain.DocumentID) error {
	if repository.removeErr != nil {
		return repository.removeErr
	}

	document, ok := repository.documents[documentID]
	if !ok || document.ApplicationID != applicationID {
		return fmt.Errorf("document not found: %s", documentID)
	}

	delete(repository.documents, documentID)

	return nil
}

// -----------------------------------------------------------------------------
// TestUpdateDocumentServiceUpdatesDocument
//
// Verifies that document metadata fields can be updated.
// -----------------------------------------------------------------------------
func TestUpdateDocumentServiceUpdatesDocument(t *testing.T) {
	repository := newFakeDocumentMaintenanceRepository()
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateDocumentService(repository, activityRepository)

	document := newDocumentMaintenanceTestDocument(t)
	repository.documents[document.ID] = document

	updatedDocument, err := service.Execute(context.Background(), UpdateDocumentInput{
		ApplicationID: document.ApplicationID,
		DocumentID:    document.ID,
		Name:          "Backend Resume v2",
		Kind:          "resume",
		Path:          "documents/backend-resume-v2.pdf",
	})
	if err != nil {
		t.Fatalf("expected document update to succeed: %v", err)
	}

	if updatedDocument.Name != "Backend Resume v2" {
		t.Fatalf("expected updated name, got %q", updatedDocument.Name)
	}

	if updatedDocument.Path != "documents/backend-resume-v2.pdf" {
		t.Fatalf("expected updated path, got %q", updatedDocument.Path)
	}

	if updatedDocument.ApplicationID != document.ApplicationID {
		t.Fatalf("expected application id to remain unchanged")
	}

	if len(activityRepository.activityEvents) != 1 {
		t.Fatalf("expected one activity event, got %d", len(activityRepository.activityEvents))
	}

	if activityRepository.activityEvents[0].Type != domain.ActivityDocumentUpdated {
		t.Fatalf("expected document updated activity event")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateDocumentServiceRejectsInvalidDocument
//
// Verifies that invalid document updates are rejected.
// -----------------------------------------------------------------------------
func TestUpdateDocumentServiceRejectsInvalidDocument(t *testing.T) {
	repository := newFakeDocumentMaintenanceRepository()
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateDocumentService(repository, activityRepository)

	document := newDocumentMaintenanceTestDocument(t)
	repository.documents[document.ID] = document

	_, err := service.Execute(context.Background(), UpdateDocumentInput{
		ApplicationID: document.ApplicationID,
		DocumentID:    document.ID,
		Name:          "",
		Kind:          "resume",
		Path:          "documents/backend-resume.pdf",
	})

	if err == nil {
		t.Fatal("expected invalid document update to fail")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateDocumentServiceReturnsRepositoryError
//
// Verifies update repository errors are returned.
// -----------------------------------------------------------------------------
func TestUpdateDocumentServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeDocumentMaintenanceRepository()
	repository.updateErr = errors.New("update failed")
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateDocumentService(repository, activityRepository)

	document := newDocumentMaintenanceTestDocument(t)
	repository.documents[document.ID] = document

	_, err := service.Execute(context.Background(), UpdateDocumentInput{
		ApplicationID: document.ApplicationID,
		DocumentID:    document.ID,
		Name:          "Backend Resume v2",
		Kind:          "resume",
		Path:          "documents/backend-resume-v2.pdf",
	})

	if err == nil {
		t.Fatal("expected update repository error")
	}
}

// -----------------------------------------------------------------------------
// TestRemoveDocumentServiceRemovesDocument
//
// Verifies that document metadata can be removed.
// -----------------------------------------------------------------------------
func TestRemoveDocumentServiceRemovesDocument(t *testing.T) {
	repository := newFakeDocumentMaintenanceRepository()
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewRemoveDocumentService(repository, activityRepository)

	document := newDocumentMaintenanceTestDocument(t)
	repository.documents[document.ID] = document

	err := service.Execute(context.Background(), RemoveDocumentInput{
		ApplicationID: document.ApplicationID,
		DocumentID:    document.ID,
	})
	if err != nil {
		t.Fatalf("expected document removal to succeed: %v", err)
	}

	if _, ok := repository.documents[document.ID]; ok {
		t.Fatal("expected document to be removed")
	}

	if len(activityRepository.activityEvents) != 1 {
		t.Fatalf("expected one activity event, got %d", len(activityRepository.activityEvents))
	}

	if activityRepository.activityEvents[0].Type != domain.ActivityDocumentRemoved {
		t.Fatalf("expected document removed activity event")
	}
}

// -----------------------------------------------------------------------------
// TestRemoveDocumentServiceReturnsRepositoryError
//
// Verifies remove repository errors are returned.
// -----------------------------------------------------------------------------
func TestRemoveDocumentServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeDocumentMaintenanceRepository()
	repository.removeErr = errors.New("remove failed")
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewRemoveDocumentService(repository, activityRepository)

	document := newDocumentMaintenanceTestDocument(t)
	repository.documents[document.ID] = document

	err := service.Execute(context.Background(), RemoveDocumentInput{
		ApplicationID: document.ApplicationID,
		DocumentID:    document.ID,
	})

	if err == nil {
		t.Fatal("expected remove repository error")
	}
}

// -----------------------------------------------------------------------------
// newDocumentMaintenanceTestDocument
//
// Creates valid document metadata for document maintenance tests.
// -----------------------------------------------------------------------------
func newDocumentMaintenanceTestDocument(t *testing.T) domain.Document {
	t.Helper()

	document, err := domain.NewDocument(
		"document-001",
		"app-001",
		"Backend Resume",
		"resume",
		"documents/backend-resume.pdf",
	)
	if err != nil {
		t.Fatalf("failed to create document: %v", err)
	}

	return document
}
