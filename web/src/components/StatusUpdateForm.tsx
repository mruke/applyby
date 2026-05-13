import { type FormEvent, useState } from "react";

import type { ApplicationStatus } from "../types/application";

/**
 * StatusUpdateFormProps
 *
 * Defines the current status and submit behavior for application status updates.
 */
type StatusUpdateFormProps = {
  currentStatus: ApplicationStatus;
  isSubmitting: boolean;
  onSubmit: (status: ApplicationStatus) => Promise<void>;
};

/**
 * statusOptions
 *
 * Provides readable status options supported by the backend lifecycle.
 */
const statusOptions: { value: ApplicationStatus; label: string }[] = [
  { value: "draft", label: "Draft" },
  { value: "interested", label: "Interested" },
  { value: "applied", label: "Applied" },
  { value: "interviewing", label: "Interviewing" },
  { value: "offer", label: "Offer" },
  { value: "rejected", label: "Rejected" },
  { value: "withdrawn", label: "Withdrawn" },
  { value: "archived", label: "Archived" }
];

/**
 * StatusUpdateForm
 *
 * Renders a focused status update control for the application detail page.
 * It shows all available statuses, while backend lifecycle rules remain authoritative.
 */
export function StatusUpdateForm({ currentStatus, isSubmitting, onSubmit }: StatusUpdateFormProps) {
  const [selectedStatus, setSelectedStatus] = useState<ApplicationStatus>(currentStatus);

  /**
   * handleSubmit
   *
   * Submits the selected status to the parent page.
   */
  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    await onSubmit(selectedStatus);
  }

  return (
    <section className="form-card" aria-labelledby="status-update-heading">
      <h2 id="status-update-heading">Update status</h2>

      <form className="status-update-form" onSubmit={handleSubmit}>
        <div className="form-field">
          <label htmlFor="application-status-update">Status</label>
          <select
            id="application-status-update"
            value={selectedStatus}
            onChange={(event) => setSelectedStatus(event.target.value as ApplicationStatus)}
          >
            {statusOptions.map((status) => (
              <option key={status.value} value={status.value}>
                {status.label}
              </option>
            ))}
          </select>
        </div>

        <div className="form-actions">
          <button type="submit" disabled={isSubmitting || selectedStatus === currentStatus}>
            {isSubmitting ? "Updating..." : "Update status"}
          </button>
        </div>
      </form>
    </section>
  );
}