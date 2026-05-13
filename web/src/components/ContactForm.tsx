import { type FormEvent, useState } from "react";

import type { CreateContactFormValues } from "../types/application";

/**
 * ContactFormProps
 *
 * Defines submit behavior and pending state for the add contact form.
 */
type ContactFormProps = {
  isSubmitting: boolean;
  onSubmit: (values: CreateContactFormValues) => Promise<void>;
};

/**
 * ContactFormState
 *
 * Represents controlled form values for adding a contact.
 */
type ContactFormState = CreateContactFormValues;

/**
 * initialFormState
 *
 * Provides default form values for a new contact.
 */
const initialFormState: ContactFormState = {
  name: "",
  email: "",
  role: ""
};

/**
 * ContactForm
 *
 * Renders the add contact form with clear labels and React-owned validation.
 */
export function ContactForm({ isSubmitting, onSubmit }: ContactFormProps) {
  const [values, setValues] = useState<ContactFormState>(initialFormState);
  const [validationMessage, setValidationMessage] = useState<string | null>(null);

  /**
   * updateValue
   *
   * Updates one controlled contact form field.
   */
  function updateValue<Field extends keyof ContactFormState>(field: Field, value: ContactFormState[Field]) {
    setValues((currentValues) => ({
      ...currentValues,
      [field]: value
    }));
  }

  /**
   * handleSubmit
   *
   * Validates required contact fields and submits values to the parent page.
   */
  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (values.name.trim() === "") {
      setValidationMessage("Contact name is required.");
      return;
    }

    setValidationMessage(null);
    await onSubmit(values);
    setValues(initialFormState);
  }

  return (
    <section className="form-card" aria-labelledby="add-contact-heading">
      <h2 id="add-contact-heading">Add contact</h2>

      {validationMessage ? (
        <p className="form-message form-message--error" role="alert">
          {validationMessage}
        </p>
      ) : null}

      <form className="contact-form" noValidate onSubmit={handleSubmit}>
        <div className="form-field">
          <label htmlFor="contact-name">Name</label>
          <input id="contact-name" value={values.name} onChange={(event) => updateValue("name", event.target.value)} />
        </div>

        <div className="form-field">
          <label htmlFor="contact-email">Email</label>
          <input
            id="contact-email"
            type="email"
            value={values.email}
            onChange={(event) => updateValue("email", event.target.value)}
          />
        </div>

        <div className="form-field">
          <label htmlFor="contact-role">Role</label>
          <input id="contact-role" value={values.role} onChange={(event) => updateValue("role", event.target.value)} />
        </div>

        <div className="form-actions">
          <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Adding..." : "Add contact"}
          </button>
        </div>
      </form>
    </section>
  );
}