import { type FormEvent, useEffect, useState } from "react";

import type { ApplicationResponse, UpdateApplicationDetailsFormValues } from "../types/application";

/**
 * ApplicationEditFormProps
 *
 * Defines the editable application and submit behavior for detail updates.
 */
type ApplicationEditFormProps = {
  application: ApplicationResponse;
  isSubmitting: boolean;
  onSubmit: (values: UpdateApplicationDetailsFormValues) => Promise<void>;
};

/**
 * formValuesFromApplication
 *
 * Converts an application response into editable form values.
 */
function formValuesFromApplication(application: ApplicationResponse): UpdateApplicationDetailsFormValues {
  return {
    title: application.title,
    companyName: application.company_name,
    companyWebsite: application.company_website,
    source: application.source,
    notes: application.notes
  };
}

/**
 * ApplicationEditForm
 *
 * Renders non-status application detail editing.
 */
export function ApplicationEditForm({ application, isSubmitting, onSubmit }: ApplicationEditFormProps) {
  const [values, setValues] = useState<UpdateApplicationDetailsFormValues>(formValuesFromApplication(application));
  const [validationMessage, setValidationMessage] = useState<string | null>(null);

  useEffect(() => {
    setValues(formValuesFromApplication(application));
  }, [application]);

  /**
   * updateValue
   *
   * Updates one controlled edit field.
   */
  function updateValue<Field extends keyof UpdateApplicationDetailsFormValues>(
    field: Field,
    value: UpdateApplicationDetailsFormValues[Field]
  ) {
    setValues((currentValues) => ({
      ...currentValues,
      [field]: value
    }));
  }

  /**
   * handleSubmit
   *
   * Validates required fields and submits edited details.
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
  }

  return (
    <section className="form-card" aria-labelledby="edit-application-heading">
      <h2 id="edit-application-heading">Edit application details</h2>

      {validationMessage ? (
        <p className="form-message form-message--error" role="alert">
          {validationMessage}
        </p>
      ) : null}

      <form className="application-form" noValidate onSubmit={handleSubmit}>
        <div className="form-field">
          <label htmlFor="edit-application-title">Application title</label>
          <input
            id="edit-application-title"
            value={values.title}
            onChange={(event) => updateValue("title", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="edit-company-name">Company name</label>
          <input
            id="edit-company-name"
            value={values.companyName}
            onChange={(event) => updateValue("companyName", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="edit-company-website">Company website</label>
          <input
            id="edit-company-website"
            type="url"
            value={values.companyWebsite}
            onChange={(event) => updateValue("companyWebsite", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="edit-application-source">Source</label>
          <input
            id="edit-application-source"
            value={values.source}
            onChange={(event) => updateValue("source", event.target.value)}
          />
        </div>

        <div className="form-field form-field--full">
          <label htmlFor="edit-application-notes">Notes</label>
          <textarea
            id="edit-application-notes"
            rows={4}
            value={values.notes}
            onChange={(event) => updateValue("notes", event.target.value)}
          />
        </div>

        <div className="form-actions">
          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Saving..." : "Save application details"}
          </button>
        </div>
      </form>
    </section>
  );
}