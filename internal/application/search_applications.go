package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
	"github.com/mruke/applyby/internal/search"
)

// -----------------------------------------------------------------------------
// ApplicationSearcher
//
// Defines the storage behavior required to search tracked applications.
// -----------------------------------------------------------------------------
type ApplicationSearcher interface {
	SearchApplications(ctx context.Context, criteria search.ApplicationCriteria) ([]domain.Application, error)
}

// -----------------------------------------------------------------------------
// SearchApplicationsService
//
// Coordinates the workflow for searching tracked applications.
// -----------------------------------------------------------------------------
type SearchApplicationsService struct {
	searcher ApplicationSearcher
}

// -----------------------------------------------------------------------------
// NewSearchApplicationsService
//
// Creates a service for the search applications workflow.
// -----------------------------------------------------------------------------
func NewSearchApplicationsService(searcher ApplicationSearcher) SearchApplicationsService {
	return SearchApplicationsService{
		searcher: searcher,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Validates criteria and searches applications through the repository boundary.
// -----------------------------------------------------------------------------
func (service SearchApplicationsService) Execute(ctx context.Context, criteria search.ApplicationCriteria) ([]domain.Application, error) {
	if service.searcher == nil {
		return nil, fmt.Errorf("application searcher is required")
	}

	criteria = criteria.Normalize()

	if err := criteria.Validate(); err != nil {
		return nil, err
	}

	return service.searcher.SearchApplications(ctx, criteria)
}
