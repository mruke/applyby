package application

import (
	"context"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// recordActivityEvent
//
// Creates and records one activity event for an application workflow.
// -----------------------------------------------------------------------------
func recordActivityEvent(
	ctx context.Context,
	recorder ActivityEventRecorder,
	applicationID domain.ApplicationID,
	eventType domain.ActivityEventType,
	description string,
) error {
	event, err := domain.NewActivityEvent(
		applicationID,
		eventType,
		time.Now().UTC(),
		description,
	)
	if err != nil {
		return err
	}

	return recorder.RecordActivityEvent(ctx, event)
}
