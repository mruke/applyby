package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestListActivityEventsServiceReturnsEvents
//
// Verifies that the activity timeline workflow returns matching events.
// -----------------------------------------------------------------------------
func TestListActivityEventsServiceReturnsEvents(t *testing.T) {
	event, err := domain.NewActivityEvent(
		"app-001",
		domain.ActivityStatusChanged,
		time.Now(),
		"Status changed from applied to interviewing.",
	)
	if err != nil {
		t.Fatalf("failed to create activity event: %v", err)
	}

	historyRepository := &fakeApplicationHistoryRepository{
		activityEvents: []domain.ActivityEvent{event},
	}
	service := NewListActivityEventsService(historyRepository)

	events, err := service.Execute(context.Background(), ListActivityEventsInput{
		ApplicationID: "app-001",
	})

	if err != nil {
		t.Fatalf("expected activity list workflow to succeed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected one activity event, got %d", len(events))
	}
}

// -----------------------------------------------------------------------------
// TestListActivityEventsServiceRejectsMissingApplicationID
//
// Verifies that activity timeline lookup requires an application id.
// -----------------------------------------------------------------------------
func TestListActivityEventsServiceRejectsMissingApplicationID(t *testing.T) {
	service := NewListActivityEventsService(&fakeApplicationHistoryRepository{})

	_, err := service.Execute(context.Background(), ListActivityEventsInput{})

	if err == nil {
		t.Fatal("expected missing application id to be rejected")
	}
}

// -----------------------------------------------------------------------------
// TestListActivityEventsServiceReturnsListError
//
// Verifies that repository list errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestListActivityEventsServiceReturnsListError(t *testing.T) {
	historyRepository := &fakeApplicationHistoryRepository{
		listErr: errors.New("list failed"),
	}
	service := NewListActivityEventsService(historyRepository)

	_, err := service.Execute(context.Background(), ListActivityEventsInput{
		ApplicationID: "app-001",
	})

	if err == nil {
		t.Fatal("expected list error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestListActivityEventsServiceRequiresLister
//
// Verifies that the activity timeline workflow requires a lister boundary.
// -----------------------------------------------------------------------------
func TestListActivityEventsServiceRequiresLister(t *testing.T) {
	service := NewListActivityEventsService(nil)

	_, err := service.Execute(context.Background(), ListActivityEventsInput{
		ApplicationID: "app-001",
	})

	if err == nil {
		t.Fatal("expected missing lister to be rejected")
	}
}
