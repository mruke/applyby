import { useCallback, useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";

import { getActivityEvents } from "../api/activity";
import { getApplicationById, updateApplicationStatus } from "../api/applications";
import { addContact, getContacts, removeContact } from "../api/contacts";
import { addDocument, getDocuments, removeDocument } from "../api/documents";
import { completeReminder, getReminders, scheduleReminder } from "../api/reminders";
import { ActivityTimeline } from "../components/ActivityTimeline";
import { ContactSection } from "../components/ContactSection";
import { DocumentSection } from "../components/DocumentSection";
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
 * SectionErrorMessages
 *
 * Represents independently recoverable load failures for detail-page sections.
 */
type SectionErrorMessages = {
  activity: string | null;
  contacts: string | null;
  documents: string | null;
  reminders: string | null;
};

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
  isRemovingContact: boolean;
  isRemovingDocument: boolean;
  isSubmittingReminder: boolean;
  isSubmittingStatus: boolean;
  reminders: ReminderResponse[];
  sectionErrors: SectionErrorMessages;
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
  sectionErrors: SectionErrorMessages;
};

/**
 * emptySectionErrors
 *
 * Provides the default section-level error state.
 */
function emptySectionErrors(): SectionErrorMessages {
  return {
    activity: null,
    contacts: null,
    documents: null,
    reminders: null
  };
}

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
    reminders: [],
    sectionErrors: emptySectionErrors()
  };
}

/**
 * sectionErrorMessage
 *
 * Provides readable section-level load error messages.
 */
function sectionErrorMessage(section: keyof SectionErrorMessages): string {
  switch (section) {
    case "activity":
      return "Activity could not be loaded.";
    case "contacts":
      return "Contacts could not be loaded.";
    case "documents":
      return "Documents could not be loaded.";
    case "reminders":
      return "Reminders could not be loaded.";
    default:
      return "Section data could not be loaded.";
  }
}

/**
 * resultValueOrDefault
 *
 * Returns the fulfilled result value or a fallback when a section request fails.
 */
function resultValueOrDefault<TValue>(result: PromiseSettledResult<TValue>, fallback: TValue): TValue {
  return result.status === "fulfilled" ? result.value : fallback;
}

/**
 * sectionErrorFromResult
 *
 * Returns a section-specific error message when a section request fails.
 */
function sectionErrorFromResult<TValue>(
  section: keyof SectionErrorMessages,
  result: PromiseSettledResult<TValue>
): string | null {
  return result.status === "rejected" ? sectionErrorMessage(section) : null;
}

/**
 * fetchApplicationDetailData
 *
 * Loads one application and its related detail resources.
 * The core application must load successfully. Related sections are allowed
 * to fail independently so one supporting feature does not block the page.
 */
async function fetchApplicationDetailData(applicationId: string): Promise<ApplicationDetailData> {
  const application = await getApplicationById(applicationId);

  if (!application) {
    return emptyApplicationDetailData();
  }

  const [remindersResult, activityResult, contactsResult, documentsResult] = await Promise.allSettled([
    getReminders(applicationId),
    getActivityEvents(applicationId),
    getContacts(applicationId),
    getDocuments(applicationId)
  ]);

  const remindersResponse = resultValueOrDefault(remindersResult, { reminders: [] });
  const activityResponse = resultValueOrDefault(activityResult, { activity_events: [] });
  const contactsResponse = resultValueOrDefault(contactsResult, { contacts: [] });
  const documentsResponse = resultValueOrDefault(documentsResult, { documents: [] });

  return {
    activityEvents: activityResponse.activity_events,
    application,
    contacts: contactsResponse.contacts,
    documents: documentsResponse.documents,
    reminders: remindersResponse.reminders,
    sectionErrors: {
      activity: sectionErrorFromResult("activity", activityResult),
      contacts: sectionErrorFromResult("contacts", contactsResult),
      documents: sectionErrorFromResult("documents", documentsResult),
      reminders: sectionErrorFromResult("reminders", remindersResult)
    }
  };
}

/**
 * applyDetailData
 *
 * Applies loaded detail data to page state.
 */
function applyDetailData(
  currentState: ApplicationDetailPageState,
  detailData: ApplicationDetailData
): ApplicationDetailPageState {
  return {
    ...currentState,
    activityEvents: detailData.activityEvents,
    application: detailData.application,
    contacts: detailData.contacts,
    documents: detailData.documents,
    errorMessage: null,
    isLoading: false,
    reminders: detailData.reminders,
    sectionErrors: detailData.sectionErrors
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
    isRemovingContact: false,
    isRemovingDocument: false,
    isSubmittingReminder: false,
    isSubmittingStatus: false,
    reminders: [],
    sectionErrors: emptySectionErrors(),
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

    setState((currentState) => applyDetailData(currentState, detailData));
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

        setState((currentState) => applyDetailData(currentState, detailData));
      } catch {
        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          application: null,
          errorMessage: "Application could not be loaded. Check that the backend is running and try again.",
          isLoading: false,
          sectionErrors: emptySectionErrors()
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

  // ---------------------------------------------------------------------------
  // handleRemoveContact
  //
  // Removes a contact, refreshes detail data, and displays feedback.
  // ---------------------------------------------------------------------------
  async function handleRemoveContact(contactId: string) {
    if (!applicationId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isRemovingContact: true,
      successMessage: null
    }));

    try {
      await removeContact(applicationId, contactId);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isRemovingContact: false,
    isRemovingDocument: false,
        successMessage: "Contact removed."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Contact could not be removed. Try again.",
        isRemovingContact: false,
    isRemovingDocument: false,
        successMessage: null
      }));
    }
  }

  // ---------------------------------------------------------------------------
  // handleRemoveDocument
  //
  // Removes document metadata, refreshes detail data, and displays feedback.
  // ---------------------------------------------------------------------------
  async function handleRemoveDocument(documentId: string) {
    if (!applicationId) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isRemovingDocument: true,
      successMessage: null
    }));

    try {
      await removeDocument(applicationId, documentId);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isRemovingDocument: false,
        successMessage: "Document metadata removed."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Document metadata could not be removed. Try again.",
        isRemovingDocument: false,
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

        <Link className="secondary-button" to={`/applications/${state.application.id}/edit`}>
          Edit application
        </Link>

        <StatusUpdateForm
          currentStatus={state.application.status}
          isSubmitting={state.isSubmittingStatus}
          onSubmit={handleStatusUpdate}
        />

        <ReminderForm isSubmitting={state.isSubmittingReminder} onSubmit={handleScheduleReminder} />

        {state.sectionErrors.reminders ? (
          <section className="state-card" aria-labelledby="reminders-heading">
            <h2 id="reminders-heading">Reminders</h2>
            <p className="form-message form-message--error" role="alert">
              {state.sectionErrors.reminders}
            </p>
          </section>
        ) : (
          <ReminderList
            isCompleting={state.isCompletingReminder}
            onComplete={handleCompleteReminder}
            reminders={state.reminders}
          />
        )}

        {state.sectionErrors.activity ? (
          <section className="state-card" aria-labelledby="activity-heading">
            <h2 id="activity-heading">Activity</h2>
            <p className="form-message form-message--error" role="alert">
              {state.sectionErrors.activity}
            </p>
          </section>
        ) : (
          <ActivityTimeline events={state.activityEvents} />
        )}

        <ContactSection
          applicationId={state.application.id}
          contacts={state.contacts}
          errorMessage={state.sectionErrors.contacts}
          isAdding={state.isAddingContact}
          isRemoving={state.isRemovingContact}
          onAdd={handleAddContact}
          onRemove={handleRemoveContact}
        />

        <DocumentSection
          applicationId={state.application.id}
          documents={state.documents}
          errorMessage={state.sectionErrors.documents}
          isAdding={state.isAddingDocument}
          isRemoving={state.isRemovingDocument}
          onAdd={handleAddDocument}
          onRemove={handleRemoveDocument}
        />
      </section>
    </>
  );
}
