package application

import (
	"context"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// fakeApplicationHistoryRepository
//
// Provides in-memory history behavior for application-layer unit tests.
// -----------------------------------------------------------------------------
type fakeApplicationHistoryRepository struct {
	activityEvents  []domain.ActivityEvent
	statusHistory   []domain.ApplicationStatusHistory
	recordEventErr  error
	recordStatusErr error
	listErr         error
}

// -----------------------------------------------------------------------------
// RecordActivityEvent
//
// Records an activity event in memory for application-layer tests.
// -----------------------------------------------------------------------------
func (repository *fakeApplicationHistoryRepository) RecordActivityEvent(ctx context.Context, event domain.ActivityEvent) error {
	if repository.recordEventErr != nil {
		return repository.recordEventErr
	}

	repository.activityEvents = append(repository.activityEvents, event)

	return nil
}

// -----------------------------------------------------------------------------
// RecordApplicationStatusHistory
//
// Records status history in memory for application-layer tests.
// -----------------------------------------------------------------------------
func (repository *fakeApplicationHistoryRepository) RecordApplicationStatusHistory(ctx context.Context, history domain.ApplicationStatusHistory) error {
	if repository.recordStatusErr != nil {
		return repository.recordStatusErr
	}

	repository.statusHistory = append(repository.statusHistory, history)

	return nil
}

// -----------------------------------------------------------------------------
// ListActivityEventsForApplication
//
// Lists in-memory activity events for one application.
// -----------------------------------------------------------------------------
func (repository *fakeApplicationHistoryRepository) ListActivityEventsForApplication(ctx context.Context, applicationID domain.ApplicationID) ([]domain.ActivityEvent, error) {
	if repository.listErr != nil {
		return nil, repository.listErr
	}

	events := []domain.ActivityEvent{}

	for _, event := range repository.activityEvents {
		if event.ApplicationID == applicationID {
			events = append(events, event)
		}
	}

	return events, nil
}
