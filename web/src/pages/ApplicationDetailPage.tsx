import { useCallback, useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";


import { removeApplication, updateApplicationStatus } from "../api/applications";
import { addContact, removeContact } from "../api/contacts";
import { addDocument, removeDocument } from "../api/documents";
import { completeReminder, removeReminder, scheduleReminder } from "../api/reminders";
import { ActivityTimeline } from "../components/ActivityTimeline";
import { ContactSection } from "../components/ContactSection";
import { DocumentSection } from "../components/DocumentSection";
import { EmptyState } from "../components/EmptyState";
import { ErrorState } from "../components/ErrorState";
import { LoadingState } from "../components/LoadingState";
import { ReminderSection } from "../components/ReminderSection";
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
import { emptySectionErrors, fetchApplicationDetailData } from "./applicationDetailData";
import type { ApplicationDetailData, SectionErrorMessages } from "./applicationDetailData";
import { formatLongDate } from "../utils/dateFormatting";

/**
 * ApplicationDetailPageState
 *
 * Represents the route-level UI state for the application detail page.
 */
type ApplicationDetailPageState = ApplicationDetailData & {
  errorMessage: string | null;
  isAddingContact: boolean;
  isAddingDocument: boolean;
  isCompletingReminder: boolean;
  isLoading: boolean;
  isRemovingApplication: boolean;
  isRemovingContact: boolean;
  isRemovingDocument: boolean;
  isRemovingReminder: boolean;
  isSchedulingReminder: boolean;
  isSubmittingReminder: boolean;
  isSubmittingStatus: boolean;
  successMessage: string | null;
};
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
  const navigate = useNavigate();
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
    isRemovingApplication: false,
    isRemovingContact: false,
    isRemovingDocument: false,
    isRemovingReminder: false,
    isSchedulingReminder: false,
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

  // ---------------------------------------------------------------------------
  // handleRemoveApplication
  //
  // Confirms and removes the current application, then returns to the application list.
  // ---------------------------------------------------------------------------
  async function handleRemoveApplication() {
    if (!state.application) {
      return;
    }

    const confirmed = window.confirm(
      "Remove this application? This also removes related reminders, contacts, documents, and activity history."
    );

    if (!confirmed) {
      return;
    }

    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      successMessage: null
    }));

    try {
      await removeApplication(state.application.id);
      void navigate("/applications");
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Application could not be removed. Try again.",
        successMessage: null
      }));
    }
  }

  // ---------------------------------------------------------------------------
  // handleRemoveReminder
  //
  // Removes a reminder, refreshes detail data, and displays feedback.
  // ---------------------------------------------------------------------------
  async function handleRemoveReminder(reminderId: string) {
    setState((currentState) => ({
      ...currentState,
      errorMessage: null,
      isRemovingReminder: true,
      successMessage: null
    }));

    try {
      await removeReminder(reminderId);
      await loadDetailData();

      setState((currentState) => ({
        ...currentState,
        isRemovingReminder: false,
        successMessage: "Reminder removed."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Reminder could not be removed. Try again.",
        isRemovingReminder: false,
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
    isRemovingReminder: false,
        successMessage: "Contact removed."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Contact could not be removed. Try again.",
        isRemovingContact: false,
    isRemovingDocument: false,
    isRemovingReminder: false,
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
    isRemovingReminder: false,
        successMessage: "Document metadata removed."
      }));
    } catch {
      setState((currentState) => ({
        ...currentState,
        errorMessage: "Document metadata could not be removed. Try again.",
        isRemovingDocument: false,
    isRemovingReminder: false,
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

        <div className="form-actions">
          <Link className="secondary-button" to={`/applications/${state.application.id}/edit`}>
            Edit application
          </Link>

          <button type="button" onClick={() => void handleRemoveApplication()}>
            Remove application
          </button>
        </div>

        <StatusUpdateForm
          currentStatus={state.application.status}
          isSubmitting={state.isSubmittingStatus}
          onSubmit={handleStatusUpdate}
        />

        <ReminderSection
          applicationId={state.application.id}
          reminders={state.reminders}
          errorMessage={state.sectionErrors.reminders}
          isCompleting={state.isCompletingReminder}
          isRemoving={state.isRemovingReminder}
          isSubmitting={state.isSubmittingReminder}
          onAdd={handleScheduleReminder}
          onComplete={handleCompleteReminder}
          onRemove={handleRemoveReminder}
        />

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
