package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// ListContactsInput
//
// Contains the data required to list contacts for an application.
// -----------------------------------------------------------------------------
type ListContactsInput struct {
	ApplicationID domain.ApplicationID
}

// -----------------------------------------------------------------------------
// ListContactsService
//
// Coordinates the workflow for listing application contacts.
// -----------------------------------------------------------------------------
type ListContactsService struct {
	repository ContactLister
}

// -----------------------------------------------------------------------------
// NewListContactsService
//
// Creates a service for the list contacts workflow.
// -----------------------------------------------------------------------------
func NewListContactsService(repository ContactLister) ListContactsService {
	return ListContactsService{
		repository: repository,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Lists contacts for an application through the repository boundary.
// -----------------------------------------------------------------------------
func (service ListContactsService) Execute(ctx context.Context, input ListContactsInput) ([]domain.Contact, error) {
	if service.repository == nil {
		return nil, fmt.Errorf("contact lister is required")
	}

	if err := input.ApplicationID.Validate(); err != nil {
		return nil, err
	}

	return service.repository.ListContactsForApplication(ctx, input.ApplicationID)
}
