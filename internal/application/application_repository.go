package application

import (
	"context"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// ApplicationSaver
//
// Defines the storage behavior required to save an application.
// -----------------------------------------------------------------------------
type ApplicationSaver interface {
	SaveApplication(ctx context.Context, application domain.Application) error
}

// -----------------------------------------------------------------------------
// ApplicationFinder
//
// Defines the storage behavior required to find an application by identity.
// -----------------------------------------------------------------------------
type ApplicationFinder interface {
	FindApplicationByID(ctx context.Context, id domain.ApplicationID) (domain.Application, error)
}

// -----------------------------------------------------------------------------
// ApplicationDetailsUpdater
//
// Defines the storage behavior required to update non-status application details.
// -----------------------------------------------------------------------------
type ApplicationDetailsUpdater interface {
	UpdateApplicationDetails(ctx context.Context, application domain.Application) error
}

// -----------------------------------------------------------------------------
// ApplicationRemover
//
// Defines the storage behavior required to remove an application.
// -----------------------------------------------------------------------------
type ApplicationRemover interface {
	RemoveApplication(ctx context.Context, id domain.ApplicationID) error
}

// -----------------------------------------------------------------------------
// ApplicationLister
//
// Defines the storage behavior required to list tracked applications.
// -----------------------------------------------------------------------------
type ApplicationLister interface {
	ListApplications(ctx context.Context) ([]domain.Application, error)
}

// -----------------------------------------------------------------------------
// ApplicationRepository
//
// Groups the full application repository behavior expected from persistence.
// -----------------------------------------------------------------------------
type ApplicationRepository interface {
	ApplicationSaver
	ApplicationFinder
	ApplicationDetailsUpdater
	ApplicationRemover
	ApplicationLister
}
