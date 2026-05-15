import { Link } from "react-router-dom";

import type { ContactResponse } from "../types/application";

/**
 * ContactListProps
 *
 * Defines the contacts rendered by the contact list component.
 */
type ContactListProps = {
  applicationId: string;
  contacts: ContactResponse[];
  isRemoving?: boolean;
  onRemove?: (contactId: string) => Promise<void>;
};

/**
 * ContactList
 *
 * Renders application contacts in a readable list.
 */
export function ContactList({ applicationId, contacts, isRemoving = false, onRemove }: ContactListProps) {
  if (contacts.length === 0) {
    return (
      <section className="state-card" aria-labelledby="contacts-heading">
        <h2 id="contacts-heading">Contacts</h2>
        <p>No contacts added yet.</p>
      </section>
    );
  }

  return (
    <section className="state-card" aria-labelledby="contacts-heading">
      <h2 id="contacts-heading">Contacts</h2>

      <ul className="stack-list">
        {contacts.map((contact) => (
          <li key={contact.id} className="stack-list__item">
            <div>
              <strong>{contact.name}</strong>
              <p>
                {contact.role || "Role not specified"}
                {contact.email ? ` · ${contact.email}` : ""}
              </p>
            </div>

            <div className="stack-list__actions">
              <Link className="secondary-button" to={`/applications/${applicationId}/contacts/${contact.id}/edit`}>
                Edit
              </Link>

              {onRemove ? (
                <button type="button" disabled={isRemoving} onClick={() => void onRemove(contact.id)}>
                  {isRemoving ? "Removing..." : "Remove"}
                </button>
              ) : null}
            </div>
          </li>
        ))}
      </ul>
    </section>
  );
}
