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
// ContactFinder
//
// Defines persistence behavior required to find one contact for an application.
// -----------------------------------------------------------------------------
type ContactFinder interface {
	FindContactByID(ctx context.Context, applicationID domain.ApplicationID, contactID domain.ContactID) (domain.Contact, error)
}

// -----------------------------------------------------------------------------
// ContactUpdater
//
// Defines persistence behavior required to update an existing contact.
// -----------------------------------------------------------------------------
type ContactUpdater interface {
	UpdateContact(ctx context.Context, contact domain.Contact) error
}

// -----------------------------------------------------------------------------
// ContactRemover
//
// Defines persistence behavior required to remove one contact from an application.
// -----------------------------------------------------------------------------
type ContactRemover interface {
	RemoveContact(ctx context.Context, applicationID domain.ApplicationID, contactID domain.ContactID) error
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
	ContactFinder
	ContactUpdater
	ContactRemover
	ContactLister
}
