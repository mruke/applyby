package application

import (
	"context"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// fakeContactRepository
//
// Provides an in-memory contact repository for application-layer unit tests.
// -----------------------------------------------------------------------------
type fakeContactRepository struct {
	contacts  map[domain.ContactID]domain.Contact
	saveErr   error
	listErr   error
	saveCalls int
}

// -----------------------------------------------------------------------------
// newFakeContactRepository
//
// Creates an empty fake contact repository for tests.
// -----------------------------------------------------------------------------
func newFakeContactRepository() *fakeContactRepository {
	return &fakeContactRepository{
		contacts: make(map[domain.ContactID]domain.Contact),
	}
}

// -----------------------------------------------------------------------------
// SaveContact
//
// Stores a contact in memory for application-layer tests.
// -----------------------------------------------------------------------------
func (repository *fakeContactRepository) SaveContact(ctx context.Context, contact domain.Contact) error {
	repository.saveCalls++

	if repository.saveErr != nil {
		return repository.saveErr
	}

	repository.contacts[contact.ID] = contact

	return nil
}

// -----------------------------------------------------------------------------
// ListContactsForApplication
//
// Lists in-memory contacts for one application.
// -----------------------------------------------------------------------------
func (repository *fakeContactRepository) ListContactsForApplication(ctx context.Context, applicationID domain.ApplicationID) ([]domain.Contact, error) {
	if repository.listErr != nil {
		return nil, repository.listErr
	}

	contacts := []domain.Contact{}

	for _, contact := range repository.contacts {
		if contact.ApplicationID == applicationID {
			contacts = append(contacts, contact)
		}
	}

	return contacts, nil
}
