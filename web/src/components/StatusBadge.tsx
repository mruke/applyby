import { applicationStatusLabels } from "../constants/applicationStatuses";
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
 * StatusBadge
 *
 * Displays an application status as a readable text-first badge.
 */
export function StatusBadge({ status }: StatusBadgeProps) {
  return <span className={`status-badge status-badge--${status}`}>{applicationStatusLabels[status]}</span>;
}