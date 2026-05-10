package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// ListActivityEventsInput
//
// Contains the data required to list activity events for an application.
// -----------------------------------------------------------------------------
type ListActivityEventsInput struct {
	ApplicationID domain.ApplicationID
}

// -----------------------------------------------------------------------------
// ListActivityEventsService
//
// Coordinates the workflow for reading an application activity timeline.
// -----------------------------------------------------------------------------
type ListActivityEventsService struct {
	lister ActivityEventLister
}

// -----------------------------------------------------------------------------
// NewListActivityEventsService
//
// Creates a service for the list activity events workflow.
// -----------------------------------------------------------------------------
func NewListActivityEventsService(lister ActivityEventLister) ListActivityEventsService {
	return ListActivityEventsService{
		lister: lister,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Validates input and retrieves activity events through the repository boundary.
// -----------------------------------------------------------------------------
func (service ListActivityEventsService) Execute(ctx context.Context, input ListActivityEventsInput) ([]domain.ActivityEvent, error) {
	if service.lister == nil {
		return nil, fmt.Errorf("activity event lister is required")
	}

	if err := input.ApplicationID.Validate(); err != nil {
		return nil, err
	}

	return service.lister.ListActivityEventsForApplication(ctx, input.ApplicationID)
}
