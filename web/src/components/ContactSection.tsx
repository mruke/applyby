import { ContactForm } from "./ContactForm";
import { ContactList } from "./ContactList";
import type { ContactResponse, CreateContactFormValues } from "../types/application";

// -----------------------------------------------------------------------------
// ContactSectionProps
//
// Defines contact section data and workflow callbacks for the detail page.
// -----------------------------------------------------------------------------
type ContactSectionProps = {
  applicationId: string;
  contacts: ContactResponse[];
  errorMessage: string | null;
  isAdding: boolean;
  isRemoving: boolean;
  onAdd: (values: CreateContactFormValues) => Promise<void>;
  onRemove: (contactId: string) => Promise<void>;
};

// -----------------------------------------------------------------------------
// ContactSection
//
// Renders contact creation, contact load errors, and contact maintenance actions.
// -----------------------------------------------------------------------------
export function ContactSection({
  applicationId,
  contacts,
  errorMessage,
  isAdding,
  isRemoving,
  onAdd,
  onRemove
}: ContactSectionProps) {
  return (
    <>
      <ContactForm isSubmitting={isAdding} onSubmit={onAdd} />

      {errorMessage ? (
        <section className="state-card" aria-labelledby="contacts-heading">
          <h2 id="contacts-heading">Contacts</h2>
          <p className="form-message form-message--error" role="alert">
            {errorMessage}
          </p>
        </section>
      ) : (
        <ContactList applicationId={applicationId} contacts={contacts} isRemoving={isRemoving} onRemove={onRemove} />
      )}
    </>
  );
}
