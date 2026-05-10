package reminders

import (
	"container/heap"
	"time"

	"github.com/mruke/applyby/internal/domain"
)

// -----------------------------------------------------------------------------
// PrioritizeReminders
//
// Orders reminders by completion state and due date using a priority queue.
// -----------------------------------------------------------------------------
func PrioritizeReminders(reminders []domain.Reminder) []domain.Reminder {
	queue := reminderPriorityQueue{}

	for _, reminder := range reminders {
		heap.Push(&queue, reminder)
	}

	prioritizedReminders := make([]domain.Reminder, 0, len(reminders))

	for queue.Len() > 0 {
		reminder := heap.Pop(&queue).(domain.Reminder)
		prioritizedReminders = append(prioritizedReminders, reminder)
	}

	return prioritizedReminders
}

// -----------------------------------------------------------------------------
// reminderPriorityQueue
//
// Implements heap behavior for reminder prioritization.
// -----------------------------------------------------------------------------
type reminderPriorityQueue []domain.Reminder

// -----------------------------------------------------------------------------
// Len
//
// Returns the number of reminders in the priority queue.
// -----------------------------------------------------------------------------
func (queue reminderPriorityQueue) Len() int {
	return len(queue)
}

// -----------------------------------------------------------------------------
// Less
//
// Reports whether one reminder has higher priority than another.
// -----------------------------------------------------------------------------
func (queue reminderPriorityQueue) Less(leftIndex int, rightIndex int) bool {
	left := queue[leftIndex]
	right := queue[rightIndex]

	if left.Completed != right.Completed {
		return !left.Completed
	}

	if left.DueAt.Equal(right.DueAt) {
		return left.Title < right.Title
	}

	return left.DueAt.Before(right.DueAt)
}

// -----------------------------------------------------------------------------
// Swap
//
// Swaps two reminders in the priority queue.
// -----------------------------------------------------------------------------
func (queue reminderPriorityQueue) Swap(leftIndex int, rightIndex int) {
	queue[leftIndex], queue[rightIndex] = queue[rightIndex], queue[leftIndex]
}

// -----------------------------------------------------------------------------
// Push
//
// Adds a reminder to the priority queue.
// -----------------------------------------------------------------------------
func (queue *reminderPriorityQueue) Push(value any) {
	reminder := value.(domain.Reminder)
	*queue = append(*queue, reminder)
}

// -----------------------------------------------------------------------------
// Pop
//
// Removes and returns the highest-priority reminder.
// -----------------------------------------------------------------------------
func (queue *reminderPriorityQueue) Pop() any {
	oldQueue := *queue
	lastIndex := len(oldQueue) - 1
	reminder := oldQueue[lastIndex]
	*queue = oldQueue[:lastIndex]

	return reminder
}

// -----------------------------------------------------------------------------
// IsOverdue
//
// Reports whether a reminder is incomplete and due before the provided time.
// -----------------------------------------------------------------------------
func IsOverdue(reminder domain.Reminder, now time.Time) bool {
	return !reminder.Completed && reminder.DueAt.Before(now)
}
