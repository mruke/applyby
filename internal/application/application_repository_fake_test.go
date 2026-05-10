package application

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// fakeApplicationRepository
//
// Provides an in-memory repository for application-layer unit tests.
// -----------------------------------------------------------------------------
type fakeApplicationRepository struct {
	applications map[domain.ApplicationID]domain.Application
	saveErr      error
	findErr      error
	listErr      error
	saveCalls    int
}

// -----------------------------------------------------------------------------
// newFakeApplicationRepository
//
// Creates an empty fake application repository for tests.
// -----------------------------------------------------------------------------
func newFakeApplicationRepository() *fakeApplicationRepository {
	return &fakeApplicationRepository{
		applications: make(map[domain.ApplicationID]domain.Application),
	}
}

// -----------------------------------------------------------------------------
// SaveApplication
//
// Stores an application in memory for application-layer tests.
// -----------------------------------------------------------------------------
func (repository *fakeApplicationRepository) SaveApplication(ctx context.Context, application domain.Application) error {
	repository.saveCalls++

	if repository.saveErr != nil {
		return repository.saveErr
	}

	repository.applications[application.ID] = application

	return nil
}

// -----------------------------------------------------------------------------
// FindApplicationByID
//
// Finds an application by identity in memory for application-layer tests.
// -----------------------------------------------------------------------------
func (repository *fakeApplicationRepository) FindApplicationByID(ctx context.Context, id domain.ApplicationID) (domain.Application, error) {
	if repository.findErr != nil {
		return domain.Application{}, repository.findErr
	}

	application, ok := repository.applications[id]
	if !ok {
		return domain.Application{}, fmt.Errorf("application not found: %s", id)
	}

	return application, nil
}

// -----------------------------------------------------------------------------
// ListApplications
//
// Returns all in-memory applications for application-layer tests.
// -----------------------------------------------------------------------------
func (repository *fakeApplicationRepository) ListApplications(ctx context.Context) ([]domain.Application, error) {
	if repository.listErr != nil {
		return nil, repository.listErr
	}

	applications := make([]domain.Application, 0, len(repository.applications))

	for _, application := range repository.applications {
		applications = append(applications, application)
	}

	return applications, nil
}
