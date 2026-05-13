import type { ApplicationStatus } from "../types/application";

/**
 * StatusOption
 *
 * Represents one readable application status option.
 */
export type StatusOption = {
  value: ApplicationStatus;
  label: string;
};

/**
 * SearchStatusOption
 *
 * Represents one readable search status option, including the empty any-status value.
 */
export type SearchStatusOption = {
  value: "" | ApplicationStatus;
  label: string;
};

/**
 * applicationStatusLabels
 *
 * Maps backend status values to readable frontend labels.
 */
export const applicationStatusLabels: Record<ApplicationStatus, string> = {
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
 * applicationStatusOptions
 *
 * Provides readable status options supported by the backend lifecycle.
 */
export const applicationStatusOptions: StatusOption[] = [
  { value: "draft", label: applicationStatusLabels.draft },
  { value: "interested", label: applicationStatusLabels.interested },
  { value: "applied", label: applicationStatusLabels.applied },
  { value: "interviewing", label: applicationStatusLabels.interviewing },
  { value: "offer", label: applicationStatusLabels.offer },
  { value: "rejected", label: applicationStatusLabels.rejected },
  { value: "withdrawn", label: applicationStatusLabels.withdrawn },
  { value: "archived", label: applicationStatusLabels.archived }
];

/**
 * applicationSearchStatusOptions
 *
 * Provides readable status options for the application search form.
 */
export const applicationSearchStatusOptions: SearchStatusOption[] = [
  { value: "", label: "Any status" },
  ...applicationStatusOptions
];