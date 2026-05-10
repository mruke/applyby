package postgres

import (
	"context"
	"fmt"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// RecordApplicationStatusHistory
//
// Persists a structured status history record for an application.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) RecordApplicationStatusHistory(ctx context.Context, history domain.ApplicationStatusHistory) error {
	if repository.db == nil {
		return fmt.Errorf("database connection is required")
	}

	if err := history.Validate(); err != nil {
		return err
	}

	_, err := repository.db.ExecContext(
		ctx,
		`
        INSERT INTO application_status_history (
            application_id,
            from_status,
            to_status,
            changed_at
        )
        VALUES ($1, $2, $3, $4)
        `,
		history.ApplicationID,
		history.FromStatus,
		history.ToStatus,
		history.ChangedAt,
	)

	return err
}

// -----------------------------------------------------------------------------
// RecordActivityEvent
//
// Persists an append-only activity event for an application.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) RecordActivityEvent(ctx context.Context, event domain.ActivityEvent) error {
	if repository.db == nil {
		return fmt.Errorf("database connection is required")
	}

	if err := event.Validate(); err != nil {
		return err
	}

	_, err := repository.db.ExecContext(
		ctx,
		`
        INSERT INTO activity_events (
            application_id,
            event_type,
            description,
            occurred_at
        )
        VALUES ($1, $2, $3, $4)
        `,
		event.ApplicationID,
		event.Type,
		event.Description,
		event.OccurredAt,
	)

	return err
}

// -----------------------------------------------------------------------------
// ListActivityEventsForApplication
//
// Retrieves activity events for one application in timeline order.
// -----------------------------------------------------------------------------
func (repository ApplicationRepository) ListActivityEventsForApplication(ctx context.Context, applicationID domain.ApplicationID) ([]domain.ActivityEvent, error) {
	if repository.db == nil {
		return nil, fmt.Errorf("database connection is required")
	}

	if err := applicationID.Validate(); err != nil {
		return nil, err
	}

	rows, err := repository.db.QueryContext(
		ctx,
		`
        SELECT
            application_id,
            event_type,
            occurred_at,
            description
        FROM activity_events
        WHERE application_id = $1
        ORDER BY occurred_at, id
        `,
		applicationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []domain.ActivityEvent{}

	for rows.Next() {
		var event domain.ActivityEvent

		if err := rows.Scan(&event.ApplicationID, &event.Type, &event.OccurredAt, &event.Description); err != nil {
			return nil, err
		}

		if err := event.Validate(); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
