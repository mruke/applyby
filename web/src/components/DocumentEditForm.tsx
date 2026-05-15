import { type FormEvent, useEffect, useState } from "react";

import type { DocumentResponse, UpdateDocumentFormValues } from "../types/application";

// -----------------------------------------------------------------------------
// DocumentEditFormProps
//
// Defines the editable document metadata and submit behavior for document updates.
// -----------------------------------------------------------------------------
type DocumentEditFormProps = {
  document: DocumentResponse;
  isSubmitting: boolean;
  onSubmit: (values: UpdateDocumentFormValues) => Promise<void>;
};

// -----------------------------------------------------------------------------
// formValuesFromDocument
//
// Converts a document response into editable form values.
// -----------------------------------------------------------------------------
function formValuesFromDocument(document: DocumentResponse): UpdateDocumentFormValues {
  return {
    name: document.name,
    kind: document.kind,
    path: document.path
  };
}

// -----------------------------------------------------------------------------
// DocumentEditForm
//
// Renders document metadata editing controls with required-field validation.
// -----------------------------------------------------------------------------
export function DocumentEditForm({ document, isSubmitting, onSubmit }: DocumentEditFormProps) {
  const [values, setValues] = useState<UpdateDocumentFormValues>(formValuesFromDocument(document));
  const [validationMessage, setValidationMessage] = useState<string | null>(null);

  useEffect(() => {
    setValues(formValuesFromDocument(document));
  }, [document]);

  function updateValue<Field extends keyof UpdateDocumentFormValues>(
    field: Field,
    value: UpdateDocumentFormValues[Field]
  ) {
    setValues((currentValues) => ({
      ...currentValues,
      [field]: value
    }));
  }

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
  }

  return (
    <section className="form-card" aria-labelledby="edit-document-heading">
      <h2 id="edit-document-heading">Edit document metadata</h2>

      {validationMessage ? (
        <p className="form-message form-message--error" role="alert">
          {validationMessage}
        </p>
      ) : null}

      <form className="document-form" noValidate onSubmit={handleSubmit}>
        <div className="form-field">
          <label htmlFor="edit-document-name">Document name</label>
          <input
            id="edit-document-name"
            value={values.name}
            onChange={(event) => updateValue("name", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="edit-document-kind">Kind</label>
          <input
            id="edit-document-kind"
            value={values.kind}
            onChange={(event) => updateValue("kind", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="edit-document-path">Path or reference</label>
          <input
            id="edit-document-path"
            value={values.path}
            onChange={(event) => updateValue("path", event.target.value)}
          />
        </div>

        <div className="form-actions">
          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Saving..." : "Save document"}
          </button>
        </div>
      </form>
    </section>
  );
}