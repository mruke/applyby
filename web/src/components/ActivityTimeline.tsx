import type { ActivityEventResponse } from "../types/application";

/**
 * ActivityTimelineProps
 *
 * Defines the activity events rendered by the timeline.
 */
type ActivityTimelineProps = {
  events: ActivityEventResponse[];
};

/**
 * formatDateTime
 *
 * Converts an API timestamp into a readable activity event date and time label.
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
 * ActivityTimeline
 *
 * Renders application activity events so the user can observe what changed recently.
 */
export function ActivityTimeline({ events }: ActivityTimelineProps) {
  if (events.length === 0) {
    return (
      <section className="state-card" aria-labelledby="activity-heading">
        <h2 id="activity-heading">Activity</h2>
        <p>No activity recorded yet.</p>
      </section>
    );
  }

  return (
    <section className="state-card" aria-labelledby="activity-heading">
      <h2 id="activity-heading">Activity</h2>

      <ol className="timeline-list">
        {events.map((event, index) => (
          <li key={`${event.occurred_at}-${index}`}>
            <strong>{event.description}</strong>
            <p>
              {event.type} · {formatDateTime(event.occurred_at)}
            </p>
          </li>
        ))}
      </ol>
    </section>
  );
}