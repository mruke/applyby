package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// RemoveApplicationInput
//
// Contains the application identity required for removal.
// -----------------------------------------------------------------------------
type RemoveApplicationInput struct {
	ID domain.ApplicationID
}

// -----------------------------------------------------------------------------
// RemoveApplicationService
//
// Coordinates the workflow for removing an application.
// -----------------------------------------------------------------------------
type RemoveApplicationService struct {
	repository ApplicationRemover
}

// -----------------------------------------------------------------------------
// NewRemoveApplicationService
//
// Creates a service for the remove application workflow.
// -----------------------------------------------------------------------------
func NewRemoveApplicationService(repository ApplicationRemover) RemoveApplicationService {
	return RemoveApplicationService{
		repository: repository,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Removes one application. Related rows are removed by database cascade rules.
// -----------------------------------------------------------------------------
func (service RemoveApplicationService) Execute(ctx context.Context, input RemoveApplicationInput) error {
	if service.repository == nil {
		return fmt.Errorf("application remover repository is required")
	}

	if err := input.ID.Validate(); err != nil {
		return err
	}

	return service.repository.RemoveApplication(ctx, input.ID)
}
