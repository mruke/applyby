import type { ReminderResponse } from "../types/application";
import { formatDateTime } from "../utils/dateFormatting";

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