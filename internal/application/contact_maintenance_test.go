package application

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// fakeContactMaintenanceRepository
//
// Provides contact maintenance behavior for application-layer tests.
// -----------------------------------------------------------------------------
type fakeContactMaintenanceRepository struct {
	contacts  map[domain.ContactID]domain.Contact
	findErr   error
	updateErr error
	removeErr error
}

// -----------------------------------------------------------------------------
// newFakeContactMaintenanceRepository
//
// Creates an empty contact maintenance fake repository.
// -----------------------------------------------------------------------------
func newFakeContactMaintenanceRepository() *fakeContactMaintenanceRepository {
	return &fakeContactMaintenanceRepository{
		contacts: make(map[domain.ContactID]domain.Contact),
	}
}

// -----------------------------------------------------------------------------
// FindContactByID
//
// Finds a contact by application and contact identity.
// -----------------------------------------------------------------------------
func (repository *fakeContactMaintenanceRepository) FindContactByID(ctx context.Context, applicationID domain.ApplicationID, contactID domain.ContactID) (domain.Contact, error) {
	if repository.findErr != nil {
		return domain.Contact{}, repository.findErr
	}

	contact, ok := repository.contacts[contactID]
	if !ok || contact.ApplicationID != applicationID {
		return domain.Contact{}, fmt.Errorf("contact not found: %s", contactID)
	}

	return contact, nil
}

// -----------------------------------------------------------------------------
// UpdateContact
//
// Updates a contact in memory.
// -----------------------------------------------------------------------------
func (repository *fakeContactMaintenanceRepository) UpdateContact(ctx context.Context, contact domain.Contact) error {
	if repository.updateErr != nil {
		return repository.updateErr
	}

	existingContact, ok := repository.contacts[contact.ID]
	if !ok || existingContact.ApplicationID != contact.ApplicationID {
		return fmt.Errorf("contact not found: %s", contact.ID)
	}

	repository.contacts[contact.ID] = contact

	return nil
}

// -----------------------------------------------------------------------------
// RemoveContact
//
// Removes a contact from memory.
// -----------------------------------------------------------------------------
func (repository *fakeContactMaintenanceRepository) RemoveContact(ctx context.Context, applicationID domain.ApplicationID, contactID domain.ContactID) error {
	if repository.removeErr != nil {
		return repository.removeErr
	}

	contact, ok := repository.contacts[contactID]
	if !ok || contact.ApplicationID != applicationID {
		return fmt.Errorf("contact not found: %s", contactID)
	}

	delete(repository.contacts, contactID)

	return nil
}

// -----------------------------------------------------------------------------
// TestUpdateContactServiceUpdatesContact
//
// Verifies that contact fields can be updated.
// -----------------------------------------------------------------------------
func TestUpdateContactServiceUpdatesContact(t *testing.T) {
	repository := newFakeContactMaintenanceRepository()
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateContactService(repository, activityRepository)

	contact := newContactMaintenanceTestContact(t)
	repository.contacts[contact.ID] = contact

	updatedContact, err := service.Execute(context.Background(), UpdateContactInput{
		ApplicationID: contact.ApplicationID,
		ContactID:     contact.ID,
		Name:          "Sam Hiring",
		Email:         "sam.hiring@example.com",
		Role:          "Hiring Manager",
	})
	if err != nil {
		t.Fatalf("expected contact update to succeed: %v", err)
	}

	if updatedContact.Name != "Sam Hiring" {
		t.Fatalf("expected updated name, got %q", updatedContact.Name)
	}

	if updatedContact.Email != "sam.hiring@example.com" {
		t.Fatalf("expected updated email, got %q", updatedContact.Email)
	}

	if updatedContact.ApplicationID != contact.ApplicationID {
		t.Fatalf("expected application id to remain unchanged")
	}

	if len(activityRepository.activityEvents) != 1 {
		t.Fatalf("expected one activity event, got %d", len(activityRepository.activityEvents))
	}

	if activityRepository.activityEvents[0].Type != domain.ActivityContactUpdated {
		t.Fatalf("expected contact updated activity event")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateContactServiceRejectsInvalidContact
//
// Verifies that invalid contact updates are rejected.
// -----------------------------------------------------------------------------
func TestUpdateContactServiceRejectsInvalidContact(t *testing.T) {
	repository := newFakeContactMaintenanceRepository()
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateContactService(repository, activityRepository)

	contact := newContactMaintenanceTestContact(t)
	repository.contacts[contact.ID] = contact

	_, err := service.Execute(context.Background(), UpdateContactInput{
		ApplicationID: contact.ApplicationID,
		ContactID:     contact.ID,
		Name:          "",
		Email:         "sam@example.com",
		Role:          "Recruiter",
	})

	if err == nil {
		t.Fatal("expected invalid contact update to fail")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateContactServiceReturnsRepositoryError
//
// Verifies update repository errors are returned.
// -----------------------------------------------------------------------------
func TestUpdateContactServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeContactMaintenanceRepository()
	repository.updateErr = errors.New("update failed")
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateContactService(repository, activityRepository)

	contact := newContactMaintenanceTestContact(t)
	repository.contacts[contact.ID] = contact

	_, err := service.Execute(context.Background(), UpdateContactInput{
		ApplicationID: contact.ApplicationID,
		ContactID:     contact.ID,
		Name:          "Sam Hiring",
		Email:         "sam.hiring@example.com",
		Role:          "Hiring Manager",
	})

	if err == nil {
		t.Fatal("expected update repository error")
	}
	if len(activityRepository.activityEvents) != 0 {
		t.Fatal("expected failed contact update not to record activity")
	}
}

// -----------------------------------------------------------------------------
// TestRemoveContactServiceRemovesContact
//
// Verifies that a contact can be removed.
// -----------------------------------------------------------------------------
func TestRemoveContactServiceRemovesContact(t *testing.T) {
	repository := newFakeContactMaintenanceRepository()
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewRemoveContactService(repository, activityRepository)

	contact := newContactMaintenanceTestContact(t)
	repository.contacts[contact.ID] = contact

	err := service.Execute(context.Background(), RemoveContactInput{
		ApplicationID: contact.ApplicationID,
		ContactID:     contact.ID,
	})
	if err != nil {
		t.Fatalf("expected contact removal to succeed: %v", err)
	}

	if _, ok := repository.contacts[contact.ID]; ok {
		t.Fatal("expected contact to be removed")
	}

	if len(activityRepository.activityEvents) != 1 {
		t.Fatalf("expected one activity event, got %d", len(activityRepository.activityEvents))
	}

	if activityRepository.activityEvents[0].Type != domain.ActivityContactRemoved {
		t.Fatalf("expected contact removed activity event")
	}
}

// -----------------------------------------------------------------------------
// TestRemoveContactServiceReturnsRepositoryError
//
// Verifies remove repository errors are returned.
// -----------------------------------------------------------------------------
func TestRemoveContactServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeContactMaintenanceRepository()
	repository.removeErr = errors.New("remove failed")
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewRemoveContactService(repository, activityRepository)

	contact := newContactMaintenanceTestContact(t)
	repository.contacts[contact.ID] = contact

	err := service.Execute(context.Background(), RemoveContactInput{
		ApplicationID: contact.ApplicationID,
		ContactID:     contact.ID,
	})

	if err == nil {
		t.Fatal("expected remove repository error")
	}
	if len(activityRepository.activityEvents) != 0 {
		t.Fatal("expected failed contact removal not to record activity")
	}
}

// -----------------------------------------------------------------------------
// newContactMaintenanceTestContact
//
// Creates a valid contact for contact maintenance tests.
// -----------------------------------------------------------------------------
func newContactMaintenanceTestContact(t *testing.T) domain.Contact {
	t.Helper()

	contact, err := domain.NewContact(
		"contact-001",
		"app-001",
		"Sam Recruiter",
		"sam@example.com",
		"Recruiter",
	)
	if err != nil {
		t.Fatalf("failed to create contact: %v", err)
	}

	return contact
}
