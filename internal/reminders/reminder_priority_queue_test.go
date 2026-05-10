package reminders

import (
	"testing"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// TestPrioritizeRemindersOrdersIncompleteRemindersByDueDate
//
// Verifies that incomplete reminders are ordered by earliest due date.
// -----------------------------------------------------------------------------
func TestPrioritizeRemindersOrdersIncompleteRemindersByDueDate(t *testing.T) {
	laterReminder := newReminderPriorityTestReminder(t, "Later", time.Date(2026, 5, 20, 9, 0, 0, 0, time.UTC), false)
	earlierReminder := newReminderPriorityTestReminder(t, "Earlier", time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC), false)

	prioritizedReminders := PrioritizeReminders([]domain.Reminder{laterReminder, earlierReminder})

	if prioritizedReminders[0].Title != "Earlier" {
		t.Fatalf("expected earliest incomplete reminder first")
	}
}

// -----------------------------------------------------------------------------
// TestPrioritizeRemindersPlacesCompletedRemindersLast
//
// Verifies that completed reminders are lower priority than incomplete reminders.
// -----------------------------------------------------------------------------
func TestPrioritizeRemindersPlacesCompletedRemindersLast(t *testing.T) {
	completedReminder := newReminderPriorityTestReminder(t, "Completed", time.Date(2026, 5, 1, 9, 0, 0, 0, time.UTC), true)
	incompleteReminder := newReminderPriorityTestReminder(t, "Incomplete", time.Date(2026, 5, 20, 9, 0, 0, 0, time.UTC), false)

	prioritizedReminders := PrioritizeReminders([]domain.Reminder{completedReminder, incompleteReminder})

	if prioritizedReminders[0].Title != "Incomplete" {
		t.Fatalf("expected incomplete reminder before completed reminder")
	}
}

// -----------------------------------------------------------------------------
// TestPrioritizeRemindersUsesTitleAsStableTieBreaker
//
// Verifies that reminders with equal due dates have deterministic ordering.
// -----------------------------------------------------------------------------
func TestPrioritizeRemindersUsesTitleAsStableTieBreaker(t *testing.T) {
	dueAt := time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC)
	secondReminder := newReminderPriorityTestReminder(t, "Beta", dueAt, false)
	firstReminder := newReminderPriorityTestReminder(t, "Alpha", dueAt, false)

	prioritizedReminders := PrioritizeReminders([]domain.Reminder{secondReminder, firstReminder})

	if prioritizedReminders[0].Title != "Alpha" {
		t.Fatalf("expected title to provide stable tie-breaker")
	}
}

// -----------------------------------------------------------------------------
// TestIsOverdue
//
// Verifies that only incomplete reminders due before now are overdue.
// -----------------------------------------------------------------------------
func TestIsOverdue(t *testing.T) {
	now := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	overdueReminder := newReminderPriorityTestReminder(t, "Overdue", time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC), false)
	completedReminder := newReminderPriorityTestReminder(t, "Completed", time.Date(2026, 5, 10, 9, 0, 0, 0, time.UTC), true)

	if !IsOverdue(overdueReminder, now) {
		t.Fatal("expected incomplete reminder due before now to be overdue")
	}

	if IsOverdue(completedReminder, now) {
		t.Fatal("expected completed reminder not to be overdue")
	}
}

// -----------------------------------------------------------------------------
// newReminderPriorityTestReminder
//
// Creates a valid reminder for reminder priority tests.
// -----------------------------------------------------------------------------
func newReminderPriorityTestReminder(t *testing.T, title string, dueAt time.Time, completed bool) domain.Reminder {
	t.Helper()

	reminder, err := domain.NewReminder(domain.ReminderID(title), "app-001", title, dueAt)
	if err != nil {
		t.Fatalf("failed to create reminder: %v", err)
	}

	reminder.Completed = completed

	return reminder
}
