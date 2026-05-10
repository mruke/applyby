package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// ListApplicationsService
//
// Coordinates the workflow for listing tracked applications.
// -----------------------------------------------------------------------------
type ListApplicationsService struct {
	repository ApplicationLister
}

// -----------------------------------------------------------------------------
// NewListApplicationsService
//
// Creates a service for the list applications workflow.
// -----------------------------------------------------------------------------
func NewListApplicationsService(repository ApplicationLister) ListApplicationsService {
	return ListApplicationsService{
		repository: repository,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Retrieves tracked applications through the repository boundary.
// -----------------------------------------------------------------------------
func (service ListApplicationsService) Execute(ctx context.Context) ([]domain.Application, error) {
	if service.repository == nil {
		return nil, fmt.Errorf("application lister is required")
	}

	return service.repository.ListApplications(ctx)
}
