package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestApplicationRepositorySavesAndListsContacts
//
// Verifies that PostgreSQL persistence can save and list contacts for an application.
// -----------------------------------------------------------------------------
func TestApplicationRepositorySavesAndListsContacts(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)
	application := newContactDocumentRepositoryTestApplication(t)

	if err := repository.SaveApplication(context.Background(), application); err != nil {
		t.Fatalf("expected application save to succeed: %v", err)
	}

	contact, err := domain.NewContact("contact-001", application.ID, "Sam Recruiter", "sam@example.com", "Recruiter")
	if err != nil {
		t.Fatalf("failed to create contact: %v", err)
	}

	if err := repository.SaveContact(context.Background(), contact); err != nil {
		t.Fatalf("expected contact save to succeed: %v", err)
	}

	contacts, err := repository.ListContactsForApplication(context.Background(), application.ID)
	if err != nil {
		t.Fatalf("expected contact list to succeed: %v", err)
	}

	if len(contacts) != 1 {
		t.Fatalf("expected one contact, got %d", len(contacts))
	}

	if contacts[0].ID != contact.ID {
		t.Fatalf("expected contact id %q, got %q", contact.ID, contacts[0].ID)
	}
}

// -----------------------------------------------------------------------------
// TestApplicationRepositoryUpdatesExistingContact
//
// Verifies that PostgreSQL persistence updates existing contacts.
// -----------------------------------------------------------------------------
func TestApplicationRepositoryUpdatesExistingContact(t *testing.T) {
	db := openIntegrationDatabase(t)
	repository := NewApplicationRepository(db)
	application := newContactDocumentRepositoryTestApplication(t)

	if err := repository.SaveApplication(context.Background(), application); err != nil {
		t.Fatalf("expected application save to succeed: %v", err)
	}

	contact, err := domain.NewContact("contact-001", application.ID, "Sam Recruiter", "sam@example.com", "Recruiter")
	if err != nil {
		t.Fatalf("failed to create contact: %v", err)
	}

	if err := repository.SaveContact(context.Background(), contact); err != nil {
		t.Fatalf("expected contact save to succeed: %v", err)
	}

	contact.Role = "Hiring Manager"

	if err := repository.SaveContact(context.Background(), contact); err != nil {
		t.Fatalf("expected contact update to succeed: %v", err)
	}

	contacts, err := repository.ListContactsForApplication(context.Background(), application.ID)
	if err != nil {
		t.Fatalf("expected contact list to succeed: %v", err)
	}

	if contacts[0].Role != "Hiring Manager" {
		t.Fatalf("expected contact role to be updated")
	}
}

// -----------------------------------------------------------------------------
// newContactDocumentRepositoryTestApplication
//
// Creates a valid application for contact and document repository tests.
// -----------------------------------------------------------------------------
func newContactDocumentRepositoryTestApplication(t *testing.T) domain.Application {
	t.Helper()

	application, err := domain.NewApplication(
		"app-001",
		"Backend Developer",
		domain.Company{Name: "Example Studio", Website: "https://example.com"},
		domain.StatusApplied,
		time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to create contact/document repository test application: %v", err)
	}

	return application
}
