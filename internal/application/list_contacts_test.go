package application

import (
	"context"
	"errors"
	"testing"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestListContactsServiceReturnsContacts
//
// Verifies that the list contacts workflow returns contacts for an application.
// -----------------------------------------------------------------------------
func TestListContactsServiceReturnsContacts(t *testing.T) {
	repository := newFakeContactRepository()
	service := NewListContactsService(repository)

	contact, err := domain.NewContact("contact-001", "app-001", "Sam Recruiter", "sam@example.com", "Recruiter")
	if err != nil {
		t.Fatalf("failed to create test contact: %v", err)
	}

	repository.contacts[contact.ID] = contact

	contacts, err := service.Execute(context.Background(), ListContactsInput{
		ApplicationID: "app-001",
	})

	if err != nil {
		t.Fatalf("expected list contacts workflow to succeed: %v", err)
	}

	if len(contacts) != 1 {
		t.Fatalf("expected one contact, got %d", len(contacts))
	}
}

// -----------------------------------------------------------------------------
// TestListContactsServiceRejectsMissingApplicationID
//
// Verifies that listing contacts requires an application id.
// -----------------------------------------------------------------------------
func TestListContactsServiceRejectsMissingApplicationID(t *testing.T) {
	service := NewListContactsService(newFakeContactRepository())

	_, err := service.Execute(context.Background(), ListContactsInput{})

	if err == nil {
		t.Fatal("expected missing application id to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestListContactsServiceReturnsRepositoryError
//
// Verifies that repository list errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestListContactsServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeContactRepository()
	repository.listErr = errors.New("list failed")
	service := NewListContactsService(repository)

	_, err := service.Execute(context.Background(), ListContactsInput{
		ApplicationID: "app-001",
	})

	if err == nil {
		t.Fatal("expected list error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestListContactsServiceRequiresRepository
//
// Verifies that the list contacts workflow requires a repository boundary.
// -----------------------------------------------------------------------------
func TestListContactsServiceRequiresRepository(t *testing.T) {
	service := NewListContactsService(nil)

	_, err := service.Execute(context.Background(), ListContactsInput{
		ApplicationID: "app-001",
	})

	if err == nil {
		t.Fatal("expected missing repository to be rejected")
	}
}
