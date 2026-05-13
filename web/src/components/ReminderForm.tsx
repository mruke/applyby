import { type FormEvent, useState } from "react";

import type { CreateReminderFormValues } from "../types/application";

/**
 * ReminderFormProps
 *
 * Defines submit behavior and pending state for the schedule reminder form.
 */
type ReminderFormProps = {
  isSubmitting: boolean;
  onSubmit: (values: CreateReminderFormValues) => Promise<void>;
};

/**
 * ReminderFormState
 *
 * Represents controlled form values for scheduling a reminder.
 */
type ReminderFormState = CreateReminderFormValues;

/**
 * initialFormState
 *
 * Provides default form values for a new reminder.
 */
const initialFormState: ReminderFormState = {
  title: "",
  dueAt: ""
};

/**
 * ReminderForm
 *
 * Renders a focused reminder scheduling form with clear labels and
 * React-owned validation messages.
 */
export function ReminderForm({ isSubmitting, onSubmit }: ReminderFormProps) {
  const [values, setValues] = useState<ReminderFormState>(initialFormState);
  const [validationMessage, setValidationMessage] = useState<string | null>(null);

  /**
   * updateValue
   *
   * Updates one controlled reminder form field.
   */
  function updateValue<Field extends keyof ReminderFormState>(field: Field, value: ReminderFormState[Field]) {
    setValues((currentValues) => ({
      ...currentValues,
      [field]: value
    }));
  }

  /**
   * handleSubmit
   *
   * Validates reminder form values and submits them to the parent page.
   */
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
    await onSubmit(values);
    setValues(initialFormState);
  }

  return (
    <section className="form-card" aria-labelledby="schedule-reminder-heading">
      <h2 id="schedule-reminder-heading">Schedule reminder</h2>

      {validationMessage ? (
        <p className="form-message form-message--error" role="alert">
          {validationMessage}
        </p>
      ) : null}

      <form className="reminder-form" noValidate onSubmit={handleSubmit}>
        <div className="form-field">
          <label htmlFor="reminder-title">Reminder title</label>
          <input
            id="reminder-title"
            name="title"
            value={values.title}
            onChange={(event) => updateValue("title", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="reminder-due-at">Due date and time</label>
          <input
            id="reminder-due-at"
            name="dueAt"
            type="datetime-local"
            value={values.dueAt}
            onChange={(event) => updateValue("dueAt", event.target.value)}
          />
        </div>

        <div className="form-actions">
          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Scheduling..." : "Schedule reminder"}
          </button>
        </div>
      </form>
    </section>
  );
}