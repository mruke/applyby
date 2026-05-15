import { type FormEvent, useEffect, useState } from "react";

import type { ContactResponse, UpdateContactFormValues } from "../types/application";

// -----------------------------------------------------------------------------
// ContactEditFormProps
//
// Defines the editable contact and submit behavior for contact updates.
// -----------------------------------------------------------------------------
type ContactEditFormProps = {
  contact: ContactResponse;
  isSubmitting: boolean;
  onSubmit: (values: UpdateContactFormValues) => Promise<void>;
};

// -----------------------------------------------------------------------------
// formValuesFromContact
//
// Converts a contact response into editable form values.
// -----------------------------------------------------------------------------
function formValuesFromContact(contact: ContactResponse): UpdateContactFormValues {
  return {
    name: contact.name,
    email: contact.email,
    role: contact.role
  };
}

// -----------------------------------------------------------------------------
// ContactEditForm
//
// Renders contact editing controls with simple required-field validation.
// -----------------------------------------------------------------------------
export function ContactEditForm({ contact, isSubmitting, onSubmit }: ContactEditFormProps) {
  const [values, setValues] = useState<UpdateContactFormValues>(formValuesFromContact(contact));
  const [validationMessage, setValidationMessage] = useState<string | null>(null);

  useEffect(() => {
    setValues(formValuesFromContact(contact));
  }, [contact]);

  // ---------------------------------------------------------------------------
  // updateValue
  //
  // Updates one controlled contact edit field.
  // ---------------------------------------------------------------------------
  function updateValue<Field extends keyof UpdateContactFormValues>(
    field: Field,
    value: UpdateContactFormValues[Field]
  ) {
    setValues((currentValues) => ({
      ...currentValues,
      [field]: value
    }));
  }

  // ---------------------------------------------------------------------------
  // handleSubmit
  //
  // Validates required fields and submits edited contact values.
  // ---------------------------------------------------------------------------
  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (values.name.trim() === "") {
      setValidationMessage("Contact name is required.");
      return;
    }

    setValidationMessage(null);
    await onSubmit(values);
  }

  return (
    <section className="form-card" aria-labelledby="edit-contact-heading">
      <h2 id="edit-contact-heading">Edit contact</h2>

      {validationMessage ? (
        <p className="form-message form-message--error" role="alert">
          {validationMessage}
        </p>
      ) : null}

      <form className="contact-form" noValidate onSubmit={handleSubmit}>
        <div className="form-field">
          <label htmlFor="edit-contact-name">Name</label>
          <input
            id="edit-contact-name"
            value={values.name}
            onChange={(event) => updateValue("name", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="edit-contact-email">Email</label>
          <input
            id="edit-contact-email"
            type="email"
            value={values.email}
            onChange={(event) => updateValue("email", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="edit-contact-role">Role</label>
          <input
            id="edit-contact-role"
            value={values.role}
            onChange={(event) => updateValue("role", event.target.value)}
          />
        </div>

        <div className="form-actions">
          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Saving..." : "Save contact"}
          </button>
        </div>
      </form>
    </section>
  );
}
