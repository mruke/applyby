package application

import (
	"context"
	"fmt"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// CreateApplicationInput
//
// Contains the data required to create a tracked application.
// -----------------------------------------------------------------------------
type CreateApplicationInput struct {
	ID        domain.ApplicationID
	Title     string
	Company   domain.Company
	Status    domain.ApplicationStatus
	Source    string
	Notes     string
	CreatedAt time.Time
}

// -----------------------------------------------------------------------------
// CreateApplicationService
//
// Coordinates the workflow for creating a tracked application.
// -----------------------------------------------------------------------------
type CreateApplicationService struct {
	repository ApplicationSaver
}

// -----------------------------------------------------------------------------
// NewCreateApplicationService
//
// Creates a service for the create application workflow.
// -----------------------------------------------------------------------------
func NewCreateApplicationService(repository ApplicationSaver) CreateApplicationService {
	return CreateApplicationService{
		repository: repository,
	}
}

// -----------------------------------------------------------------------------
// Execute
//
// Validates and saves a new application through the repository boundary.
// -----------------------------------------------------------------------------
func (service CreateApplicationService) Execute(ctx context.Context, input CreateApplicationInput) (domain.Application, error) {
	if service.repository == nil {
		return domain.Application{}, fmt.Errorf("application saver is required")
	}

	application, err := domain.NewApplication(input.ID, input.Title, input.Company, input.Status, input.CreatedAt)
	if err != nil {
		return domain.Application{}, err
	}

	application.Source = input.Source
	application.Notes = input.Notes

	if err := service.repository.SaveApplication(ctx, application); err != nil {
		return domain.Application{}, err
	}

	return application, nil
}
