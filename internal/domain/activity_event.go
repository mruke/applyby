package domain

import (
	"fmt"
	"time"
)

// -----------------------------------------------------------------------------
// ActivityEventType
//
// Represents a typed activity history event category.
// -----------------------------------------------------------------------------
type ActivityEventType string

const (
	ActivityApplicationCreated ActivityEventType = "application_created"
	ActivityStatusChanged      ActivityEventType = "status_changed"
	ActivityNoteAdded          ActivityEventType = "note_added"
	ActivityReminderScheduled  ActivityEventType = "reminder_scheduled"
	ActivityReminderCompleted  ActivityEventType = "reminder_completed"
	ActivityContactAdded       ActivityEventType = "contact_added"
	ActivityDocumentAdded      ActivityEventType = "document_added"
)

// -----------------------------------------------------------------------------
// ActivityEvent
//
// Represents an append-only historical fact about an application.
// -----------------------------------------------------------------------------
type ActivityEvent struct {
	ApplicationID ApplicationID
	Type          ActivityEventType
	OccurredAt    time.Time
	Description   string
}

// -----------------------------------------------------------------------------
// AllActivityEventTypes
//
// Returns every supported activity event type in a stable order.
// -----------------------------------------------------------------------------
func AllActivityEventTypes() []ActivityEventType {
	return []ActivityEventType{
		ActivityApplicationCreated,
		ActivityStatusChanged,
		ActivityNoteAdded,
		ActivityReminderScheduled,
		ActivityReminderCompleted,
		ActivityContactAdded,
		ActivityDocumentAdded,
	}
}

// -----------------------------------------------------------------------------
// IsValidActivityEventType
//
// Reports whether a value is one of the supported activity event types.
// -----------------------------------------------------------------------------
func IsValidActivityEventType(eventType ActivityEventType) bool {
	for _, validType := range AllActivityEventTypes() {
		if eventType == validType {
			return true
		}
	}

	return false
}

// -----------------------------------------------------------------------------
// NewActivityEvent
//
// Creates an activity event after applying basic domain validation.
// -----------------------------------------------------------------------------
func NewActivityEvent(applicationID ApplicationID, eventType ActivityEventType, occurredAt time.Time, description string) (ActivityEvent, error) {
	event := ActivityEvent{
		ApplicationID: applicationID,
		Type:          eventType,
		OccurredAt:    occurredAt,
		Description:   description,
	}

	if err := event.Validate(); err != nil {
		return ActivityEvent{}, err
	}

	return event, nil
}

// -----------------------------------------------------------------------------
// Validate
//
// Verifies that an activity event has a valid application, type, timestamp, and description.
// -----------------------------------------------------------------------------
func (event ActivityEvent) Validate() error {
	if err := event.ApplicationID.Validate(); err != nil {
		return err
	}

	if !IsValidActivityEventType(event.Type) {
		return fmt.Errorf("invalid activity event type: %q", event.Type)
	}

	if err := requireNonZeroTime("activity event timestamp", event.OccurredAt); err != nil {
		return err
	}

	return requireNonEmptyField("activity event description", event.Description)
}
