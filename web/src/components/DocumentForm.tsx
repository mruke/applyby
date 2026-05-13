import { type FormEvent, useState } from "react";

import type { CreateDocumentFormValues } from "../types/application";

/**
 * DocumentFormProps
 *
 * Defines submit behavior and pending state for the add document metadata form.
 */
type DocumentFormProps = {
  isSubmitting: boolean;
  onSubmit: (values: CreateDocumentFormValues) => Promise<void>;
};

/**
 * DocumentFormState
 *
 * Represents controlled form values for adding document metadata.
 */
type DocumentFormState = CreateDocumentFormValues;

/**
 * initialFormState
 *
 * Provides default form values for a new document metadata record.
 */
const initialFormState: DocumentFormState = {
  name: "",
  kind: "",
  path: ""
};

/**
 * DocumentForm
 *
 * Renders the add document metadata form. This records metadata only;
 * actual file upload/storage is intentionally deferred.
 */
export function DocumentForm({ isSubmitting, onSubmit }: DocumentFormProps) {
  const [values, setValues] = useState<DocumentFormState>(initialFormState);
  const [validationMessage, setValidationMessage] = useState<string | null>(null);

  /**
   * updateValue
   *
   * Updates one controlled document form field.
   */
  function updateValue<Field extends keyof DocumentFormState>(field: Field, value: DocumentFormState[Field]) {
    setValues((currentValues) => ({
      ...currentValues,
      [field]: value
    }));
  }

  /**
   * handleSubmit
   *
   * Validates required document metadata fields and submits values to the parent page.
   */
  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (values.name.trim() === "") {
      setValidationMessage("Document name is required.");
      return;
    }

    if (values.kind.trim() === "") {
      setValidationMessage("Document kind is required.");
      return;
    }

    setValidationMessage(null);
    await onSubmit(values);
    setValues(initialFormState);
  }

  return (
    <section className="form-card" aria-labelledby="add-document-heading">
      <h2 id="add-document-heading">Add document metadata</h2>

      {validationMessage ? (
        <p className="form-message form-message--error" role="alert">
          {validationMessage}
        </p>
      ) : null}

      <form className="document-form" noValidate onSubmit={handleSubmit}>
        <div className="form-field">
          <label htmlFor="document-name">Document name</label>
          <input id="document-name" value={values.name} onChange={(event) => updateValue("name", event.target.value)} />
        </div>

        <div className="form-field">
          <label htmlFor="document-kind">Kind</label>
          <input id="document-kind" value={values.kind} onChange={(event) => updateValue("kind", event.target.value)} />
        </div>

        <div className="form-field">
          <label htmlFor="document-path">Path or reference</label>
          <input id="document-path" value={values.path} onChange={(event) => updateValue("path", event.target.value)} />
        </div>

        <div className="form-actions">
          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Adding..." : "Add document"}
          </button>
        </div>
      </form>
    </section>
  );
}