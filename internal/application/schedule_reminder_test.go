package application

import (
	"context"
	"errors"
	"testing"
	"time"
)

// -----------------------------------------------------------------------------
// TestScheduleReminderServiceSavesValidReminder
//
// Verifies that the schedule reminder workflow validates and saves a reminder.
// -----------------------------------------------------------------------------
func TestScheduleReminderServiceSavesValidReminder(t *testing.T) {
	repository := newFakeReminderRepository()
	service := NewScheduleReminderService(repository, &fakeApplicationHistoryRepository{})
	dueAt := time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC)

	reminder, err := service.Execute(context.Background(), ScheduleReminderInput{
		ID:            "rem-001",
		ApplicationID: "app-001",
		Title:         "Follow up",
		DueAt:         dueAt,
	})

	if err != nil {
		t.Fatalf("expected schedule reminder workflow to succeed: %v", err)
	}

	if reminder.ID != "rem-001" {
		t.Fatalf("expected reminder id to be preserved")
	}

	if repository.saveCalls != 1 {
		t.Fatalf("expected repository save to be called once")
	}
}

// -----------------------------------------------------------------------------
// TestScheduleReminderServiceRejectsInvalidReminder
//
// Verifies that invalid reminder data is rejected before saving.
// -----------------------------------------------------------------------------
func TestScheduleReminderServiceRejectsInvalidReminder(t *testing.T) {
	repository := newFakeReminderRepository()
	service := NewScheduleReminderService(repository, &fakeApplicationHistoryRepository{})

	_, err := service.Execute(context.Background(), ScheduleReminderInput{
		ID:            "",
		ApplicationID: "app-001",
		Title:         "Follow up",
		DueAt:         time.Now(),
	})

	if err == nil {
		t.Fatal("expected invalid reminder to be rejected")
	}

	if repository.saveCalls != 0 {
		t.Fatal("expected invalid reminder not to be saved")
	}
}

// -----------------------------------------------------------------------------
// TestScheduleReminderServiceReturnsRepositoryError
//
// Verifies that repository save errors are returned to the caller.
// -----------------------------------------------------------------------------
func TestScheduleReminderServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeReminderRepository()
	repository.saveErr = errors.New("save failed")
	service := NewScheduleReminderService(repository, &fakeApplicationHistoryRepository{})

	_, err := service.Execute(context.Background(), ScheduleReminderInput{
		ID:            "rem-001",
		ApplicationID: "app-001",
		Title:         "Follow up",
		DueAt:         time.Now(),
	})

	if err == nil {
		t.Fatal("expected repository error to be returned")
	}
}

// -----------------------------------------------------------------------------
// TestScheduleReminderServiceRequiresRepository
//
// Verifies that the schedule reminder workflow requires a repository boundary.
// -----------------------------------------------------------------------------
func TestScheduleReminderServiceRequiresRepository(t *testing.T) {
	service := NewScheduleReminderService(nil, &fakeApplicationHistoryRepository{})

	_, err := service.Execute(context.Background(), ScheduleReminderInput{
		ID:            "rem-001",
		ApplicationID: "app-001",
		Title:         "Follow up",
		DueAt:         time.Now(),
	})

	if err == nil {
		t.Fatal("expected missing repository to be rejected")
	}
}
