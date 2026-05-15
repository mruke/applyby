import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import { getApplicationById } from "../api/applications";
import { getContacts, updateContact } from "../api/contacts";
import { ContactEditForm } from "../components/ContactEditForm";
import { EmptyState } from "../components/EmptyState";
import { ErrorState } from "../components/ErrorState";
import { LoadingState } from "../components/LoadingState";
import type { ApplicationResponse, ContactResponse, UpdateContactFormValues } from "../types/application";

// -----------------------------------------------------------------------------
// ContactEditPageState
//
// Represents loading and submit state for editing one contact.
// -----------------------------------------------------------------------------
type ContactEditPageState = {
  application: ApplicationResponse | null;
  contact: ContactResponse | null;
  errorMessage: string | null;
  isLoading: boolean;
  isSubmitting: boolean;
};

// -----------------------------------------------------------------------------
// findContactById
//
// Finds one contact in an application contact collection.
// -----------------------------------------------------------------------------
function findContactById(contacts: ContactResponse[], contactId: string): ContactResponse | null {
  return contacts.find((contact) => contact.id === contactId) ?? null;
}

// -----------------------------------------------------------------------------
// ContactEditPage
//
// Loads an application contact and renders the contact edit workflow.
// -----------------------------------------------------------------------------
export function ContactEditPage() {
  const navigate = useNavigate();
  const { applicationId, contactId } = useParams<{ applicationId: string; contactId: string }>();

  const [state, setState] = useState<ContactEditPageState>({
    application: null,
    contact: null,
    errorMessage: null,
    isLoading: true,
    isSubmitting: false
  });

  useEffect(() => {
    let isCurrentRequest = true;

    // -------------------------------------------------------------------------
    // loadContact
    //
    // Loads the parent application and selects the requested contact.
    // -------------------------------------------------------------------------
    async function loadContact() {
      try {
        if (!applicationId || !contactId) {
          throw new Error("missing route ids");
        }

        const [application, contactsResponse] = await Promise.all([
          getApplicationById(applicationId),
          getContacts(applicationId)
        ]);

        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          application,
          contact: findContactById(contactsResponse.contacts, contactId),
          errorMessage: null,
          isLoading: false
        }));
      } catch {
        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          application: null,
          contact: null,
          errorMessage: "Contact could not be loaded. Check that the backend is running and try again.",
          isLoading: false
        }));
      }
    }

    void loadContact();

    return () => {
      isCurrentRequest = false;
    };
  }, [applicationId, contactId]);

  // ---------------------------------------------------------------------------
  // handleSubmit
  //
  // Updates the contact and returns to the application detail page.
  // ---------------------------------------------------------------------------
  async function handleSubmit(values: UpdateContactFormValues) {
    if (!applicationId || !contactId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isSubmitting: true
    }));

    try {
      await updateContact(applicationId, contactId, values);
      void navigate(`/applications/${applicationId}`);
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Contact could not be updated. Check the form and try again.",
        isSubmitting: false
      }));
    }
  }

  if (state.isLoading) {
    return <LoadingState message="Loading contact..." />;
  }

  if (state.errorMessage && !state.application && !state.contact) {
    return <ErrorState title="Contact could not be loaded" message={state.errorMessage} />;
  }

  if (!state.application) {
    return (
      <EmptyState
        title="Application not found"
        message="No application matched this route. Return to the applications list and choose an existing application."
      />
    );
  }

  if (!state.contact) {
    return (
      <EmptyState
        title="Contact not found"
        message="No contact matched this route. Return to the application detail page and choose an existing contact."
      />
    );
  }

  return (
    <>
      <header className="page-header">
        <Link to={`/applications/${state.application.id}`}>Back to application</Link>
        <h1>Edit contact</h1>
        <p>
          {state.contact.name} · {state.application.title}
        </p>
      </header>

      {state.errorMessage ? (
        <p className="form-message form-message--error" role="alert">
          {state.errorMessage}
        </p>
      ) : null}

      <ContactEditForm contact={state.contact} isSubmitting={state.isSubmitting} onSubmit={handleSubmit} />

      <p>
        <Link className="secondary-button" to={`/applications/${state.application.id}`}>
          Cancel
        </Link>
      </p>
    </>
  );
}
