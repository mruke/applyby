import { ReminderForm } from "./ReminderForm";
import { ReminderList } from "./ReminderList";
import type { CreateReminderFormValues, ReminderResponse } from "../types/application";

// -----------------------------------------------------------------------------
// ReminderSectionProps
//
// Defines reminder section data and workflow callbacks for the detail page.
// -----------------------------------------------------------------------------
type ReminderSectionProps = {
  applicationId: string;
  reminders: ReminderResponse[];
  errorMessage: string | null;
  isCompleting: boolean;
  isRemoving: boolean;
  isSubmitting: boolean;
  onAdd: (values: CreateReminderFormValues) => Promise<void>;
  onComplete: (reminderId: string) => Promise<void>;
  onRemove: (reminderId: string) => Promise<void>;
};

// -----------------------------------------------------------------------------
// ReminderSection
//
// Renders reminder creation, reminder load errors, and reminder maintenance actions.
// -----------------------------------------------------------------------------
export function ReminderSection({
  applicationId,
  reminders,
  errorMessage,
  isCompleting,
  isRemoving,
  isSubmitting,
  onAdd,
  onComplete,
  onRemove
}: ReminderSectionProps) {
  return (
    <>
      <ReminderForm isSubmitting={isSubmitting} onSubmit={onAdd} />

      {errorMessage ? (
        <section className="state-card" aria-labelledby="reminders-heading">
          <h2 id="reminders-heading">Reminders</h2>
          <p className="form-message form-message--error" role="alert">
            {errorMessage}
          </p>
        </section>
      ) : (
        <ReminderList
          applicationId={applicationId}
          reminders={reminders}
          isCompleting={isCompleting}
          isRemoving={isRemoving}
          onComplete={onComplete}
          onRemove={onRemove}
        />
      )}
    </>
  );
}