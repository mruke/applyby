package application

import (
	"context"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// ActivityEventRecorder
//
// Defines persistence behavior required to record activity events.
// -----------------------------------------------------------------------------
type ActivityEventRecorder interface {
	RecordActivityEvent(ctx context.Context, event domain.ActivityEvent) error
}

// -----------------------------------------------------------------------------
// ActivityEventLister
//
// Defines persistence behavior required to list activity events for an application.
// -----------------------------------------------------------------------------
type ActivityEventLister interface {
	ListActivityEventsForApplication(ctx context.Context, applicationID domain.ApplicationID) ([]domain.ActivityEvent, error)
}

// -----------------------------------------------------------------------------
// ApplicationStatusHistoryRecorder
//
// Defines persistence behavior required to record structured status history.
// -----------------------------------------------------------------------------
type ApplicationStatusHistoryRecorder interface {
	RecordApplicationStatusHistory(ctx context.Context, history domain.ApplicationStatusHistory) error
}

// -----------------------------------------------------------------------------
// ApplicationHistoryRecorder
//
// Groups the history recording behavior required by status-changing workflows.
// -----------------------------------------------------------------------------
type ApplicationHistoryRecorder interface {
	ActivityEventRecorder
	ApplicationStatusHistoryRecorder
}
