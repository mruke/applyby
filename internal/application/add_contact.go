package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// AddContactInput
//
// Contains the data required to attach a contact to an application.
// -----------------------------------------------------------------------------
type AddContactInput struct {
	ID            domain.ContactID
	ApplicationID domain.ApplicationID
	Name          string
	Email         string
	Role          string
}

// -----------------------------------------------------------------------------
// AddContactService
//
// Coordinates the workflow for attaching a contact to an application.
// -----------------------------------------------------------------------------
type AddContactService struct {
	repository ContactSaver
}

// -----------------------------------------------------------------------------
// NewAddContactService
//
// Creates a service for the add contact workflow.
// -----------------------------------------------------------------------------
func NewAddContactService(repository ContactSaver) AddContactService {
	return AddContactService{
		repository: repository,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Validates and saves a contact through the repository boundary.
// -----------------------------------------------------------------------------
func (service AddContactService) Execute(ctx context.Context, input AddContactInput) (domain.Contact, error) {
	if service.repository == nil {
		return domain.Contact{}, fmt.Errorf("contact saver is required")
	}

	contact, err := domain.NewContact(input.ID, input.ApplicationID, input.Name, input.Email, input.Role)
	if err != nil {
		return domain.Contact{}, err
	}

	if err := service.repository.SaveContact(ctx, contact); err != nil {
		return domain.Contact{}, err
	}

	return contact, nil
}
