/**
 * formatShortDate
 *
 * Converts an API timestamp into a compact readable date label.
 */
export function formatShortDate(timestamp: string): string {
  return new Intl.DateTimeFormat(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric"
  }).format(new Date(timestamp));
}

/**
 * formatLongDate
 *
 * Converts an API timestamp into a full readable date label.
 */
export function formatLongDate(timestamp: string): string {
  return new Intl.DateTimeFormat(undefined, {
    year: "numeric",
    month: "long",
    day: "numeric"
  }).format(new Date(timestamp));
}

/**
 * formatDateTime
 *
 * Converts an API timestamp into a readable date and time label.
 */
export function formatDateTime(timestamp: string): string {
  return new Intl.DateTimeFormat(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "numeric",
    minute: "2-digit"
  }).format(new Date(timestamp));
}