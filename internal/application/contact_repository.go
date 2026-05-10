package application

import (
	"context"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// ContactSaver
//
// Defines persistence behavior required to save a contact.
// -----------------------------------------------------------------------------
type ContactSaver interface {
	SaveContact(ctx context.Context, contact domain.Contact) error
}

// -----------------------------------------------------------------------------
// ContactLister
//
// Defines persistence behavior required to list contacts for an application.
// -----------------------------------------------------------------------------
type ContactLister interface {
	ListContactsForApplication(ctx context.Context, applicationID domain.ApplicationID) ([]domain.Contact, error)
}

// -----------------------------------------------------------------------------
// ContactRepository
//
// Groups the full contact repository behavior expected from persistence.
// -----------------------------------------------------------------------------
type ContactRepository interface {
	ContactSaver
	ContactLister
}
