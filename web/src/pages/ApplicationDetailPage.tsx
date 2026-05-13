import { useCallback, useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";

import { getActivityEvents } from "../api/activity";
import { getApplicationById, updateApplicationStatus } from "../api/applications";
import { addContact, getContacts } from "../api/contacts";
import { addDocument, getDocuments } from "../api/documents";
import { completeReminder, getReminders, scheduleReminder } from "../api/reminders";
import { ActivityTimeline } from "../components/ActivityTimeline";
import { ContactForm } from "../components/ContactForm";
import { ContactList } from "../components/ContactList";
import { DocumentForm } from "../components/DocumentForm";
import { DocumentList } from "../components/DocumentList";
import { EmptyState } from "../components/EmptyState";
import { ErrorState } from "../components/ErrorState";
import { LoadingState } from "../components/LoadingState";
import { ReminderForm } from "../components/ReminderForm";
import { ReminderList } from "../components/ReminderList";
import { StatusBadge } from "../components/StatusBadge";
import { StatusUpdateForm } from "../components/StatusUpdateForm";
import type {
  ActivityEventResponse,
  ApplicationResponse,
  ApplicationStatus,
  ContactResponse,
  CreateContactFormValues,
  CreateDocumentFormValues,
  CreateReminderFormValues,
  DocumentResponse,
  ReminderResponse
} from "../types/application";
import { formatLongDate } from "../utils/dateFormatting";

/**
 * ApplicationDetailPageState
 *
 * Represents the loading and workflow state for the application detail page.
 */
type ApplicationDetailPageState = {
  activityEvents: ActivityEventResponse[];
  application: ApplicationResponse | null;
  contacts: ContactResponse[];
  documents: DocumentResponse[];
  errorMessage: string | null;
  isAddingContact: boolean;
  isAddingDocument: boolean;
  isCompletingReminder: boolean;
  isLoading: boolean;
  isSubmittingReminder: boolean;
  isSubmittingStatus: boolean;
  reminders: ReminderResponse[];
  successMessage: string | null;
};

/**
 * ApplicationDetailData
 *
 * Represents the data needed to render the application detail page.
 */
type ApplicationDetailData = {
  activityEvents: ActivityEventResponse[];
  application: ApplicationResponse | null;
  contacts: ContactResponse[];
  documents: DocumentResponse[];
  reminders: ReminderResponse[];
};

/**
 * emptyApplicationDetailData
 *
 * Provides empty detail data for a missing application.
 */
function emptyApplicationDetailData(): ApplicationDetailData {
  return {
    activityEvents: [],
    application: null,
    contacts: [],
    documents: [],
    reminders: []
  };
}

/**
 * fetchApplicationDetailData
 *
 * Loads one application and its related detail resources.
 */
async function fetchApplicationDetailData(applicationId: string): Promise<ApplicationDetailData> {
  const application = await getApplicationById(applicationId);

  if (!application) {
    return emptyApplicationDetailData();
  }

  const [remindersResponse, activityResponse, contactsResponse, documentsResponse] = await Promise.all([
    getReminders(applicationId),
    getActivityEvents(applicationId),
    getContacts(applicationId),
    getDocuments(applicationId)
  ]);

  return {
    activityEvents: activityResponse.activity_events,
    application,
    contacts: contactsResponse.contacts,
    documents: documentsResponse.documents,
    reminders: remindersResponse.reminders
  };
}

/**
 * ApplicationDetailPage
 *
 * Loads one application and exposes detail workflows for status, reminders,
 * contacts, document metadata, and activity history.
 */
export function ApplicationDetailPage() {
  const { applicationId } = useParams<{ applicationId: string }>();

  const [state, setState] = useState<ApplicationDetailPageState>({
    activityEvents: [],
    application: null,
    contacts: [],
    documents: [],
    errorMessage: null,
    isAddingContact: false,
    isAddingDocument: false,
    isCompletingReminder: false,
    isLoading: true,
    isSubmittingReminder: false,
    isSubmittingStatus: false,
    reminders: [],
    successMessage: null
  });

  /**
   * loadDetailData
   *
   * Loads current detail data and applies it to page state.
   */
  const loadDetailData = useCallback(async () => {
    if (!applicationId) {
      setState((currentState) => ({
        ...currentState,
        application: null,
        errorMessage: "Application id is missing from the route.",
        isLoading: false
      }));
      return;
    }

    const detailData = await fetchApplicationDetailData(applicationId);

    setState((currentState) => ({
      ...currentState,
      activityEvents: detailData.activityEvents,
      application: detailData.application,
      contacts: detailData.contacts,
      documents: detailData.documents,
      errorMessage: null,
      isLoading: false,
      reminders: detailData.reminders
    }));
  }, [applicationId]);

  useEffect(() => {
    let isCurrentRequest = true;

    async function loadInitialDetailData() {
      try {
        if (!applicationId) {
          throw new Error("missing application id");
        }

        const detailData = await fetchApplicationDetailData(applicationId);

        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          activityEvents: detailData.activityEvents,
          application: detailData.application,
          contacts: detailData.contacts,
          documents: detailData.documents,
          errorMessage: null,
          isLoading: false,
          reminders: detailData.reminders
        }));
      } catch {
        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          application: null,
          errorMessage: "Application could not be loaded. Check that the backend is running and try again.",
          isLoading: false
        }));
      }
    }

    void loadInitialDetailData();

    return () => {
      isCurrentRequest = false;
    };
  }, [applicationId]);

  /**
   * handleStatusUpdate
   *
   * Updates the application status, refreshes the detail state, and displays feedback.
   */
  async function handleStatusUpdate(status: ApplicationStatus) {
    if (!applicationId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isSubmittingStatus: true,
      successMessage: null
    }));

    try {
      await updateApplicationStatus(applicationId, status);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isSubmittingStatus: false,
        successMessage: "Status updated."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Status could not be updated. Check the selected status and try again.",
        isSubmittingStatus: false,
        successMessage: null
      }));
    }
  }

  /**
   * handleScheduleReminder
   *
   * Schedules a reminder, refreshes detail data, and displays feedback.
   */
  async function handleScheduleReminder(values: CreateReminderFormValues) {
    if (!applicationId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isSubmittingReminder: true,
      successMessage: null
    }));

    try {
      await scheduleReminder(applicationId, values);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isSubmittingReminder: false,
        successMessage: "Reminder scheduled."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Reminder could not be scheduled. Check the form and try again.",
        isSubmittingReminder: false,
        successMessage: null
      }));
    }
  }

  /**
   * handleCompleteReminder
   *
   * Completes a reminder, refreshes detail data, and displays feedback.
   */
  async function handleCompleteReminder(reminderId: string) {
    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isCompletingReminder: true,
      successMessage: null
    }));

    try {
      await completeReminder(reminderId);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isCompletingReminder: false,
        successMessage: "Reminder completed."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Reminder could not be completed. Try again.",
        isCompletingReminder: false,
        successMessage: null
      }));
    }
  }

  /**
   * handleAddContact
   *
   * Adds a contact, refreshes detail data, and displays feedback.
   */
  async function handleAddContact(values: CreateContactFormValues) {
    if (!applicationId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isAddingContact: true,
      successMessage: null
    }));

    try {
      await addContact(applicationId, values);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isAddingContact: false,
        successMessage: "Contact added."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Contact could not be added. Check the form and try again.",
        isAddingContact: false,
        successMessage: null
      }));
    }
  }

  /**
   * handleAddDocument
   *
   * Adds document metadata, refreshes detail data, and displays feedback.
   */
  async function handleAddDocument(values: CreateDocumentFormValues) {
    if (!applicationId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isAddingDocument: true,
      successMessage: null
    }));

    try {
      await addDocument(applicationId, values);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isAddingDocument: false,
        successMessage: "Document metadata added."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Document metadata could not be added. Check the form and try again.",
        isAddingDocument: false,
        successMessage: null
      }));
    }
  }

  if (state.isLoading) {
    return <LoadingState message="Loading application..." />;
  }

  if (state.errorMessage && !state.application) {
    return <ErrorState title="Application could not be loaded" message={state.errorMessage} />;
  }

  if (!state.application) {
    return (
      <EmptyState
        title="Application not found"
        message="No application matched this route. Return to the applications list and choose an existing application."
      />
    );
  }

  return (
    <>
      <header className="page-header">
        <Link to="/applications">Back to applications</Link>
        <h1>{state.application.title}</h1>
        <p>{state.application.company_name}</p>
      </header>

      {state.successMessage ? (
        <p className="form-message form-message--success" role="status">
          {state.successMessage}
        </p>
      ) : null}

      {state.errorMessage ? (
        <p className="form-message form-message--error" role="alert">
          {state.errorMessage}
        </p>
      ) : null}

      <section className="detail-grid" aria-label="Application details">
        <article className="state-card">
          <h2>Summary</h2>
          <dl className="detail-list">
            <div>
              <dt>Status</dt>
              <dd>
                <StatusBadge status={state.application.status} />
              </dd>
            </div>
            <div>
              <dt>Company website</dt>
              <dd>
                {state.application.company_website ? (
                  <a href={state.application.company_website}>{state.application.company_website}</a>
                ) : (
                  "Not specified"
                )}
              </dd>
            </div>
            <div>
              <dt>Source</dt>
              <dd>{state.application.source || "Not specified"}</dd>
            </div>
            <div>
              <dt>Created</dt>
              <dd>{formatLongDate(state.application.created_at)}</dd>
            </div>
            <div>
              <dt>Notes</dt>
              <dd>{state.application.notes || "No notes added yet."}</dd>
            </div>
          </dl>
        </article>

        <StatusUpdateForm
          currentStatus={state.application.status}
          isSubmitting={state.isSubmittingStatus}
          onSubmit={handleStatusUpdate}
        />

        <ReminderForm isSubmitting={state.isSubmittingReminder} onSubmit={handleScheduleReminder} />

        <ReminderList
          isCompleting={state.isCompletingReminder}
          onComplete={handleCompleteReminder}
          reminders={state.reminders}
        />

        <ActivityTimeline events={state.activityEvents} />

        <ContactForm isSubmitting={state.isAddingContact} onSubmit={handleAddContact} />

        <ContactList contacts={state.contacts} />

        <DocumentForm isSubmitting={state.isAddingDocument} onSubmit={handleAddDocument} />

        <DocumentList documents={state.documents} />
      </section>
    </>
  );
}