import { Link } from "react-router-dom";

import type { ReminderResponse } from "../types/application";
import { formatLongDate } from "../utils/dateFormatting";

// -----------------------------------------------------------------------------
// ReminderListProps
//
// Defines reminders rendered by the reminder list component.
// -----------------------------------------------------------------------------
type ReminderListProps = {
  applicationId: string;
  reminders: ReminderResponse[];
  isCompleting: boolean;
  isRemoving?: boolean;
  onComplete: (reminderId: string) => Promise<void>;
  onRemove?: (reminderId: string) => Promise<void>;
};

// -----------------------------------------------------------------------------
// ReminderList
//
// Renders reminders with completion and maintenance actions.
// -----------------------------------------------------------------------------
export function ReminderList({
  applicationId,
  reminders,
  isCompleting,
  isRemoving = false,
  onComplete,
  onRemove
}: ReminderListProps) {
  if (reminders.length === 0) {
    return (
      <section className="state-card" aria-labelledby="reminders-heading">
        <h2 id="reminders-heading">Reminders</h2>
        <p>No reminders scheduled yet.</p>
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
                Due {formatLongDate(reminder.due_at)}
                {reminder.completed ? " · Completed" : ""}
              </p>
            </div>

            <div className="stack-list__actions">
              <Link className="secondary-button" to={`/applications/${applicationId}/reminders/${reminder.id}/edit`}>
                Edit
              </Link>

              {!reminder.completed ? (
                <button type="button" disabled={isCompleting} onClick={() => void onComplete(reminder.id)}>
                  {isCompleting ? "Completing..." : "Complete"}
                </button>
              ) : null}

              {onRemove ? (
                <button type="button" disabled={isRemoving} onClick={() => void onRemove(reminder.id)}>
                  {isRemoving ? "Removing..." : "Remove"}
                </button>
              ) : null}
            </div>
          </li>
        ))}
      </ul>
    </section>
  );
}