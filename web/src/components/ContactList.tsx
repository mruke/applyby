import type { ContactResponse } from "../types/application";

/**
 * ContactListProps
 *
 * Defines the contacts rendered by the contact list component.
 */
type ContactListProps = {
  contacts: ContactResponse[];
};

/**
 * ContactList
 *
 * Renders application contacts in a readable list.
 */
export function ContactList({ contacts }: ContactListProps) {
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
          </li>
        ))}
      </ul>
    </section>
  );
}