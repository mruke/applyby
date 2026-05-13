import type { ReminderResponse } from "../types/application";

/**
 * ReminderListProps
 *
 * Defines the reminders to render and the complete action for incomplete reminders.
 */
type ReminderListProps = {
  isCompleting: boolean;
  onComplete: (reminderId: string) => Promise<void>;
  reminders: ReminderResponse[];
};

/**
 * formatDateTime
 *
 * Converts an API timestamp into a readable reminder date and time label.
 */
function formatDateTime(timestamp: string): string {
  return new Intl.DateTimeFormat(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "numeric",
    minute: "2-digit"
  }).format(new Date(timestamp));
}

/**
 * ReminderList
 *
 * Renders reminders in a readable list with a clear complete action
 * for incomplete reminders.
 */
export function ReminderList({ isCompleting, onComplete, reminders }: ReminderListProps) {
  if (reminders.length === 0) {
    return (
      <section className="state-card" aria-labelledby="reminders-heading">
        <h2 id="reminders-heading">Reminders</h2>
        <p>No reminders scheduled for this application.</p>
      </section>
    );
  }

  return (
    <section className="state-card" aria-labelledby="reminders-heading">
      <h2 id="reminders-heading">Reminders</h2>

      <ul className="stack-list">
        {reminders.map((reminder) => (
          <li key={reminder.id} className="stack-list__item">
            <div>
              <strong>{reminder.title}</strong>
              <p>
                Due {formatDateTime(reminder.due_at)} · {reminder.completed ? "Completed" : "Incomplete"}
              </p>
            </div>

            {!reminder.completed ? (
              <button type="button" disabled={isCompleting} onClick={() => void onComplete(reminder.id)}>
                {isCompleting ? "Completing..." : "Complete"}
              </button>
            ) : null}
          </li>
        ))}
      </ul>
    </section>
  );
}