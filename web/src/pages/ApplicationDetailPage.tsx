import { useCallback, useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";

import { getActivityEvents } from "../api/activity";
import { getApplicationById, updateApplicationStatus } from "../api/applications";
import { completeReminder, getReminders, scheduleReminder } from "../api/reminders";
import { ActivityTimeline } from "../components/ActivityTimeline";
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
  CreateReminderFormValues,
  ReminderResponse
} from "../types/application";

/**
 * ApplicationDetailPageState
 *
 * Represents the loading and workflow state for the application detail page.
 */
type ApplicationDetailPageState = {
  activityEvents: ActivityEventResponse[];
  application: ApplicationResponse | null;
  errorMessage: string | null;
  isCompletingReminder: boolean;
  isLoading: boolean;
  isSubmittingReminder: boolean;
  isSubmittingStatus: boolean;
  reminders: ReminderResponse[];
  successMessage: string | null;
};

/**
 * formatDate
 *
 * Converts an API timestamp into a readable date label for detail display.
 */
function formatDate(timestamp: string): string {
  return new Intl.DateTimeFormat(undefined, {
    year: "numeric",
    month: "long",
    day: "numeric"
  }).format(new Date(timestamp));
}

/**
 * ApplicationDetailPage
 *
 * Loads one application from the existing applications API and exposes summary,
 * status, reminder, and activity timeline workflows.
 */
export function ApplicationDetailPage() {
  const { applicationId } = useParams<{ applicationId: string }>();

  const [state, setState] = useState<ApplicationDetailPageState>({
    activityEvents: [],
    application: null,
    errorMessage: null,
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
   * Loads the current application, reminders, and activity timeline data.
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

    const application = await getApplicationById(applicationId);

    if (!application) {
      setState((currentState) => ({
        ...currentState,
        application: null,
        activityEvents: [],
        errorMessage: null,
        isLoading: false,
        reminders: []
      }));
      return;
    }

    const [remindersResponse, activityResponse] = await Promise.all([
      getReminders(applicationId),
      getActivityEvents(applicationId)
    ]);

    setState((currentState) => ({
      ...currentState,
      activityEvents: activityResponse.activity_events,
      application,
      errorMessage: null,
      isLoading: false,
      reminders: remindersResponse.reminders
    }));
  }, [applicationId]);

  useEffect(() => {
    let isCurrentRequest = true;

    async function loadInitialDetailData() {
      try {
        if (!applicationId) {
          throw new Error("missing application id");
        }

        const application = await getApplicationById(applicationId);

        if (!isCurrentRequest) {
          return;
        }

        if (!application) {
          setState((currentState) => ({
            ...currentState,
            application: null,
            activityEvents: [],
            errorMessage: null,
            isLoading: false,
            reminders: []
          }));
          return;
        }

        const [remindersResponse, activityResponse] = await Promise.all([
          getReminders(applicationId),
          getActivityEvents(applicationId)
        ]);

        if (!isCurrentRequest) {
          return;
        }

        setState((currentState) => ({
          ...currentState,
          activityEvents: activityResponse.activity_events,
          application,
          errorMessage: null,
          isLoading: false,
          reminders: remindersResponse.reminders
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
              <dd>{formatDate(state.application.created_at)}</dd>
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
      </section>
    </>
  );
}