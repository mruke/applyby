package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// GetApplicationInput
//
// Represents the application identity needed for the detail lookup workflow.
// -----------------------------------------------------------------------------
type GetApplicationInput struct {
	ID domain.ApplicationID
}

// -----------------------------------------------------------------------------
// GetApplicationService
//
// Coordinates the workflow for retrieving one tracked application.
// -----------------------------------------------------------------------------
type GetApplicationService struct {
	repository ApplicationFinder
}

// -----------------------------------------------------------------------------
// NewGetApplicationService
//
// Creates a service for the get application workflow.
// -----------------------------------------------------------------------------
func NewGetApplicationService(repository ApplicationFinder) GetApplicationService {
	return GetApplicationService{
		repository: repository,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Retrieves one tracked application through the repository boundary.
// -----------------------------------------------------------------------------
func (service GetApplicationService) Execute(ctx context.Context, input GetApplicationInput) (domain.Application, error) {
	if service.repository == nil {
		return domain.Application{}, fmt.Errorf("application finder is required")
	}

	if err := input.ID.Validate(); err != nil {
		return domain.Application{}, err
	}

	return service.repository.FindApplicationByID(ctx, input.ID)
}
