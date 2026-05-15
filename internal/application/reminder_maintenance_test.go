package application

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// fakeReminderMaintenanceRepository
//
// Provides reminder maintenance behavior for application-layer tests.
// -----------------------------------------------------------------------------
type fakeReminderMaintenanceRepository struct {
	reminders map[domain.ReminderID]domain.Reminder
	findErr   error
	updateErr error
	removeErr error
}

// -----------------------------------------------------------------------------
// newFakeReminderMaintenanceRepository
//
// Creates an empty reminder maintenance fake repository.
// -----------------------------------------------------------------------------
func newFakeReminderMaintenanceRepository() *fakeReminderMaintenanceRepository {
	return &fakeReminderMaintenanceRepository{
		reminders: make(map[domain.ReminderID]domain.Reminder),
	}
}

// -----------------------------------------------------------------------------
// FindReminderByID
//
// Finds a reminder by identity.
// -----------------------------------------------------------------------------
func (repository *fakeReminderMaintenanceRepository) FindReminderByID(ctx context.Context, id domain.ReminderID) (domain.Reminder, error) {
	if repository.findErr != nil {
		return domain.Reminder{}, repository.findErr
	}

	reminder, ok := repository.reminders[id]
	if !ok {
		return domain.Reminder{}, fmt.Errorf("reminder not found: %s", id)
	}

	return reminder, nil
}

// -----------------------------------------------------------------------------
// UpdateReminder
//
// Updates a reminder in memory.
// -----------------------------------------------------------------------------
func (repository *fakeReminderMaintenanceRepository) UpdateReminder(ctx context.Context, reminder domain.Reminder) error {
	if repository.updateErr != nil {
		return repository.updateErr
	}

	if _, ok := repository.reminders[reminder.ID]; !ok {
		return fmt.Errorf("reminder not found: %s", reminder.ID)
	}

	repository.reminders[reminder.ID] = reminder

	return nil
}

// -----------------------------------------------------------------------------
// RemoveReminder
//
// Removes a reminder from memory.
// -----------------------------------------------------------------------------
func (repository *fakeReminderMaintenanceRepository) RemoveReminder(ctx context.Context, id domain.ReminderID) error {
	if repository.removeErr != nil {
		return repository.removeErr
	}

	if _, ok := repository.reminders[id]; !ok {
		return fmt.Errorf("reminder not found: %s", id)
	}

	delete(repository.reminders, id)

	return nil
}

// -----------------------------------------------------------------------------
// TestUpdateReminderServiceUpdatesReminder
//
// Verifies that reminder details can be updated.
// -----------------------------------------------------------------------------
func TestUpdateReminderServiceUpdatesReminder(t *testing.T) {
	repository := newFakeReminderMaintenanceRepository()
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateReminderService(repository, activityRepository)

	reminder := newReminderMaintenanceTestReminder(t)
	reminder.Completed = true
	repository.reminders[reminder.ID] = reminder

	updatedDueAt := time.Date(2026, 5, 20, 9, 30, 0, 0, time.UTC)

	updatedReminder, err := service.Execute(context.Background(), UpdateReminderInput{
		ID:    reminder.ID,
		Title: "Send updated follow-up",
		DueAt: updatedDueAt,
	})
	if err != nil {
		t.Fatalf("expected reminder update to succeed: %v", err)
	}

	if updatedReminder.Title != "Send updated follow-up" {
		t.Fatalf("expected updated title, got %q", updatedReminder.Title)
	}

	if !updatedReminder.DueAt.Equal(updatedDueAt) {
		t.Fatalf("expected updated due date")
	}

	if !updatedReminder.Completed {
		t.Fatal("expected completed state to be preserved")
	}

	if updatedReminder.ApplicationID != reminder.ApplicationID {
		t.Fatalf("expected application id to remain unchanged")
	}

	if len(activityRepository.activityEvents) != 1 {
		t.Fatalf("expected one activity event, got %d", len(activityRepository.activityEvents))
	}

	if activityRepository.activityEvents[0].Type != domain.ActivityReminderUpdated {
		t.Fatalf("expected reminder updated activity event")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateReminderServiceRejectsInvalidReminder
//
// Verifies that invalid reminder updates are rejected.
// -----------------------------------------------------------------------------
func TestUpdateReminderServiceRejectsInvalidReminder(t *testing.T) {
	repository := newFakeReminderMaintenanceRepository()
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateReminderService(repository, activityRepository)

	reminder := newReminderMaintenanceTestReminder(t)
	repository.reminders[reminder.ID] = reminder

	_, err := service.Execute(context.Background(), UpdateReminderInput{
		ID:    reminder.ID,
		Title: "",
		DueAt: reminder.DueAt,
	})

	if err == nil {
		t.Fatal("expected invalid reminder update to fail")
	}
}

// -----------------------------------------------------------------------------
// TestUpdateReminderServiceReturnsRepositoryError
//
// Verifies update repository errors are returned.
// -----------------------------------------------------------------------------
func TestUpdateReminderServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeReminderMaintenanceRepository()
	repository.updateErr = errors.New("update failed")
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewUpdateReminderService(repository, activityRepository)

	reminder := newReminderMaintenanceTestReminder(t)
	repository.reminders[reminder.ID] = reminder

	_, err := service.Execute(context.Background(), UpdateReminderInput{
		ID:    reminder.ID,
		Title: "Send updated follow-up",
		DueAt: reminder.DueAt,
	})

	if err == nil {
		t.Fatal("expected update repository error")
	}
	if len(activityRepository.activityEvents) != 0 {
		t.Fatal("expected failed reminder update not to record activity")
	}
}

// -----------------------------------------------------------------------------
// TestRemoveReminderServiceRemovesReminder
//
// Verifies that a reminder can be removed.
// -----------------------------------------------------------------------------
func TestRemoveReminderServiceRemovesReminder(t *testing.T) {
	repository := newFakeReminderMaintenanceRepository()
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewRemoveReminderService(repository, activityRepository)

	reminder := newReminderMaintenanceTestReminder(t)
	repository.reminders[reminder.ID] = reminder

	err := service.Execute(context.Background(), RemoveReminderInput{
		ID: reminder.ID,
	})
	if err != nil {
		t.Fatalf("expected reminder removal to succeed: %v", err)
	}

	if _, ok := repository.reminders[reminder.ID]; ok {
		t.Fatal("expected reminder to be removed")
	}

	if len(activityRepository.activityEvents) != 1 {
		t.Fatalf("expected one activity event, got %d", len(activityRepository.activityEvents))
	}

	if activityRepository.activityEvents[0].Type != domain.ActivityReminderRemoved {
		t.Fatalf("expected reminder removed activity event")
	}
}

// -----------------------------------------------------------------------------
// TestRemoveReminderServiceReturnsRepositoryError
//
// Verifies remove repository errors are returned.
// -----------------------------------------------------------------------------
func TestRemoveReminderServiceReturnsRepositoryError(t *testing.T) {
	repository := newFakeReminderMaintenanceRepository()
	repository.removeErr = errors.New("remove failed")
	activityRepository := &fakeApplicationHistoryRepository{}
	service := NewRemoveReminderService(repository, activityRepository)

	reminder := newReminderMaintenanceTestReminder(t)
	repository.reminders[reminder.ID] = reminder

	err := service.Execute(context.Background(), RemoveReminderInput{
		ID: reminder.ID,
	})

	if err == nil {
		t.Fatal("expected remove repository error")
	}
	if len(activityRepository.activityEvents) != 0 {
		t.Fatal("expected failed reminder removal not to record activity")
	}
}

// -----------------------------------------------------------------------------
// newReminderMaintenanceTestReminder
//
// Creates a valid reminder for reminder maintenance tests.
// -----------------------------------------------------------------------------
func newReminderMaintenanceTestReminder(t *testing.T) domain.Reminder {
	t.Helper()

	reminder, err := domain.NewReminder(
		"reminder-001",
		"app-001",
		"Send follow-up",
		time.Date(2026, 5, 18, 9, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to create reminder: %v", err)
	}

	return reminder
}
