import type { ApplicationStatus } from "../types/application";

/**
 * StatusBadgeProps
 *
 * Defines the application status displayed by the status badge.
 */
type StatusBadgeProps = {
  status: ApplicationStatus;
};

/**
 * statusLabels
 *
 * Maps backend status values to readable frontend labels.
 * Text remains the primary status indicator so meaning does not depend on color.
 */
const statusLabels: Record<ApplicationStatus, string> = {
  draft: "Draft",
  interested: "Interested",
  applied: "Applied",
  interviewing: "Interviewing",
  offer: "Offer",
  rejected: "Rejected",
  withdrawn: "Withdrawn",
  archived: "Archived"
};

/**
 * StatusBadge
 *
 * Displays an application status as a readable text-first badge.
 */
export function StatusBadge({ status }: StatusBadgeProps) {
  return <span className={`status-badge status-badge--${status}`}>{statusLabels[status]}</span>;
}