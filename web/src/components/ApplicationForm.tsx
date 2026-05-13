import { type FormEvent, useState } from "react";

import type { ApplicationStatus, CreateApplicationFormValues } from "../types/application";

/**
 * ApplicationFormProps
 *
 * Defines the submit behavior and pending state for the create application form.
 */
type ApplicationFormProps = {
  isSubmitting: boolean;
  onSubmit: (values: CreateApplicationFormValues) => Promise<void>;
};

/**
 * ApplicationFormState
 *
 * Represents controlled form values for creating an application.
 */
type ApplicationFormState = CreateApplicationFormValues;

/**
 * initialFormState
 *
 * Provides default form values for a new application.
 */
const initialFormState: ApplicationFormState = {
  title: "",
  companyName: "",
  companyWebsite: "",
  status: "applied",
  source: "",
  notes: ""
};

/**
 * applicationStatuses
 *
 * Provides the status options exposed by the create application form.
 */
const applicationStatuses: { value: ApplicationStatus; label: string }[] = [
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
 * ApplicationForm
 *
 * Renders the create application form with accessible labels, required-field
 * validation, and clear submit behavior.
 */
export function ApplicationForm({ isSubmitting, onSubmit }: ApplicationFormProps) {
  const [values, setValues] = useState<ApplicationFormState>(initialFormState);
  const [validationMessage, setValidationMessage] = useState<string | null>(null);

  /**
   * updateValue
   *
   * Updates one controlled form field.
   */
  function updateValue<Field extends keyof ApplicationFormState>(field: Field, value: ApplicationFormState[Field]) {
    setValues((currentValues) => ({
      ...currentValues,
      [field]: value
    }));
  }

  /**
   * handleSubmit
   *
   * Validates required fields and submits the form values to the parent page.
   */
  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (values.title.trim() === "") {
      setValidationMessage("Application title is required.");
      return;
    }

    if (values.companyName.trim() === "") {
      setValidationMessage("Company name is required.");
      return;
    }

    setValidationMessage(null);
    await onSubmit(values);
    setValues(initialFormState);
  }

  return (
    <section className="form-card" aria-labelledby="create-application-heading">
      <h2 id="create-application-heading">Add application</h2>

      {validationMessage ? (
        <p className="form-message form-message--error" role="alert">
          {validationMessage}
        </p>
      ) : null}

      <form className="application-form" noValidate onSubmit={handleSubmit}>
        <div className="form-field">
          <label htmlFor="application-title">Application title</label>
          <input
            id="application-title"
            name="title"
            required
            value={values.title}
            onChange={(event) => updateValue("title", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="company-name">Company name</label>
          <input
            id="company-name"
            name="companyName"
            required
            value={values.companyName}
            onChange={(event) => updateValue("companyName", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="company-website">Company website</label>
          <input
            id="company-website"
            name="companyWebsite"
            type="url"
            value={values.companyWebsite}
            onChange={(event) => updateValue("companyWebsite", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="application-status">Status</label>
          <select
            id="application-status"
            name="status"
            value={values.status}
            onChange={(event) => updateValue("status", event.target.value as ApplicationStatus)}
          >
            {applicationStatuses.map((status) => (
              <option key={status.value} value={status.value}>
                {status.label}
              </option>
            ))}
          </select>
        </div>

        <div className="form-field">
          <label htmlFor="application-source">Source</label>
          <input
            id="application-source"
            name="source"
            value={values.source}
            onChange={(event) => updateValue("source", event.target.value)}
          />
        </div>

        <div className="form-field form-field--full">
          <label htmlFor="application-notes">Notes</label>
          <textarea
            id="application-notes"
            name="notes"
            rows={4}
            value={values.notes}
            onChange={(event) => updateValue("notes", event.target.value)}
          />
        </div>

        <div className="form-actions">
          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Adding..." : "Add application"}
          </button>
        </div>
      </form>
    </section>
  );
}