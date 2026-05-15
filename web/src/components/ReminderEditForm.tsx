import { type FormEvent, useEffect, useState } from "react";

import type { ReminderResponse, UpdateReminderFormValues } from "../types/application";

// -----------------------------------------------------------------------------
// ReminderEditFormProps
//
// Defines editable reminder data and submit behavior for reminder updates.
// -----------------------------------------------------------------------------
type ReminderEditFormProps = {
  reminder: ReminderResponse;
  isSubmitting: boolean;
  onSubmit: (values: UpdateReminderFormValues) => Promise<void>;
};

// -----------------------------------------------------------------------------
// dateTimeLocalFromRFC3339
//
// Converts an API timestamp into a datetime-local input value.
// -----------------------------------------------------------------------------
function dateTimeLocalFromRFC3339(value: string): string {
  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "";
  }

  const offsetMs = date.getTimezoneOffset() * 60 * 1000;
  return new Date(date.getTime() - offsetMs).toISOString().slice(0, 16);
}

// -----------------------------------------------------------------------------
// rfc3339FromDateTimeLocal
//
// Converts a datetime-local input value into an API timestamp.
// -----------------------------------------------------------------------------
function rfc3339FromDateTimeLocal(value: string): string {
  return new Date(value).toISOString();
}

// -----------------------------------------------------------------------------
// formValuesFromReminder
//
// Converts a reminder response into editable form values.
// -----------------------------------------------------------------------------
function formValuesFromReminder(reminder: ReminderResponse): UpdateReminderFormValues {
  return {
    title: reminder.title,
    dueAt: dateTimeLocalFromRFC3339(reminder.due_at)
  };
}

// -----------------------------------------------------------------------------
// ReminderEditForm
//
// Renders reminder editing controls with required-field validation.
// -----------------------------------------------------------------------------
export function ReminderEditForm({ reminder, isSubmitting, onSubmit }: ReminderEditFormProps) {
  const [values, setValues] = useState<UpdateReminderFormValues>(formValuesFromReminder(reminder));
  const [validationMessage, setValidationMessage] = useState<string | null>(null);

  useEffect(() => {
    setValues(formValuesFromReminder(reminder));
  }, [reminder]);

  function updateValue<Field extends keyof UpdateReminderFormValues>(
    field: Field,
    value: UpdateReminderFormValues[Field]
  ) {
    setValues((currentValues) => ({
      ...currentValues,
      [field]: value
    }));
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (values.title.trim() === "") {
      setValidationMessage("Reminder title is required.");
      return;
    }

    if (values.dueAt.trim() === "") {
      setValidationMessage("Reminder due date is required.");
      return;
    }

    setValidationMessage(null);

    await onSubmit({
      title: values.title,
      dueAt: rfc3339FromDateTimeLocal(values.dueAt)
    });
  }

  return (
    <section className="form-card" aria-labelledby="edit-reminder-heading">
      <h2 id="edit-reminder-heading">Edit reminder</h2>

      {validationMessage ? (
        <p className="form-message form-message--error" role="alert">
          {validationMessage}
        </p>
      ) : null}

      <form className="reminder-form" noValidate onSubmit={handleSubmit}>
        <div className="form-field">
          <label htmlFor="edit-reminder-title">Reminder title</label>
          <input
            id="edit-reminder-title"
            value={values.title}
            onChange={(event) => updateValue("title", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="edit-reminder-due-at">Due date</label>
          <input
            id="edit-reminder-due-at"
            type="datetime-local"
            value={values.dueAt}
            onChange={(event) => updateValue("dueAt", event.target.value)}
          />
        </div>

        <div className="form-actions">
          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Saving..." : "Save reminder"}
          </button>
        </div>
      </form>
    </section>
  );
}